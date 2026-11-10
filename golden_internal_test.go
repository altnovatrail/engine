package engine

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The golden vectors of the choice (docs/architecture.md, sections 3.5 and 3.7; R07, R10): inputs
// on the graphs of the three modules of the first slice, from the packages of the release of
// modules, read by the learner model's real Read; and the choice with its code, and the bits of
// the P the engine computed of a task it gives. The files in testdata/choose are the same on every
// target, and change only with a version of engine or learner-model; ENGINE_UPDATE=1 writes them.
//
// Each scenario names the branch of table 3.5 it is about, and a test holds the pair of branch
// and code to the table: so the vectors check the codes of every branch, and that the rule gives
// each of the nine codes.

const (
	updateEnv  = "ENGINE_UPDATE"
	goldenTime = int64(1_800_000_000) // 2027-01-15, UTC
	day        = int64(86_400)
)

// branchCodes is table 3.5: the code of each branch, by its name in the scenarios. "nothing" is
// nothing to give, with no code.
type branchCode struct {
	branch string
	code   ReasonCode
}

var branchCodes = []branchCode{
	{"1", ReasonOpenRun}, {"2", ReasonSkillChosen}, {"3", ReasonConsolidation}, {"4", ReasonProbe},
	{"5", ReasonRepetition}, {"6", ReasonNewSkill}, {"6a", ReasonAuthorOrder}, {"7", ReasonQueue},
	{"7-mastered", ReasonRotation}, {"8", ReasonRotation}, {"9-new", ReasonNewSkill}, {"9", ReasonRotation},
	{"nothing", ReasonUnspecified},
}

// reasonNames are the names of the codes in cause.proto, by number.
var reasonNames = []string{
	"REASON_CODE_UNSPECIFIED", "REASON_CODE_CONSOLIDATION", "REASON_CODE_NEW_SKILL", "REASON_CODE_QUEUE",
	"REASON_CODE_REPETITION", "REASON_CODE_PROBE", "REASON_CODE_SKILL_CHOSEN", "REASON_CODE_AUTHOR_ORDER",
	"REASON_CODE_OPEN_RUN", "REASON_CODE_ROTATION",
}

// golden is the input of one vector: the modules of the engine, the learner model's state and the
// learner's grade, and the rest of the input. The reading is Read's, at goldenTime.
type golden struct {
	modules []Module
	state   learnermodel.State
	grade   learnermodel.Grade
	in      Input
}

func (g golden) input(t *testing.T) Input {
	t.Helper()
	graphs := learnermodel.Graphs{}
	var ids []string
	for _, m := range g.modules {
		graphs[m.Graph.Module] = learnermodel.ModuleGraph{Graph: m.Graph, Current: true}
		ids = append(ids, m.Graph.Module)
	}
	st := g.state
	st.Kind = learnermodel.StateKind
	st.Modules = slices.Clone(st.Modules)
	slices.SortFunc(st.Modules, func(a, b learnermodel.ModuleState) int { return strings.Compare(a.Module, b.Module) })
	for i := range st.Modules {
		st.Modules[i].Skills = slices.Clone(st.Modules[i].Skills)
		slices.SortFunc(st.Modules[i].Skills, func(a, b learnermodel.SkillState) int { return strings.Compare(a.Skill, b.Skill) })
	}
	r, err := learnermodel.Read(st, graphs, goldenTime, g.grade, ids)
	if err != nil {
		t.Fatal(err)
	}
	in := g.in
	in.Time, in.Readout, in.Modules = goldenTime, r, g.modules
	return in
}

var learnerGrade3 = learnermodel.Grade{Value: 3, Known: true}

// module is the state of a module with the level given, past the trial series at a weight of 10.
func module(id string, level float64, start learnermodel.BandRef, skills ...learnermodel.SkillState) learnermodel.ModuleState {
	return learnermodel.ModuleState{Module: id, Level: level, AnswerWeight: 10, StartBand: start, Skills: skills}
}

// practised is the state of a skill with evidence: a correction, a count, and the last practice,
// with a stability of 30 days, so that forgetting moves its P little.
func practised(skill string, correction float64, count int64, last int64) learnermodel.SkillState {
	return learnermodel.SkillState{Skill: skill, Correction: correction, AnswerWeight: float64(count), EvidenceCount: count, LastPractice: last, Stability: 30}
}

// masteredState is a skill the model has mastered at a band, with a stability in days.
func masteredState(skill string, band int32, stability float64, last int64) learnermodel.SkillState {
	s := practised(skill, 1, 12, last)
	s.Mastered, s.MasteryBand, s.Stability = true, band, stability
	return s
}

func issued(at int64, skills ...string) []SkillIssue {
	var out []SkillIssue
	for _, s := range skills {
		out = append(out, SkillIssue{Skill: s, Issued: true, LastIssued: at})
	}
	return out
}

func passed(l Lesson, at int64) Run {
	return Run{Lesson: l.ID, Skills: l.Skills, End: RunPassed, EndedAt: at}
}

func lessonNamed(m Module, name string) Lesson {
	return m.Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(m.Graph.Module, name) })]
}

func failedOn(task, skill string) Last {
	return Last{Task: task, Results: []SkillResult{{Skill: skill, Answered: 1}}}
}

type scenario struct {
	name, branch string
	input        func(rel map[string]release) golden
}

// goldenSets are the scenarios, by the file of their vectors.
var goldenSets = []struct {
	file      string
	scenarios []scenario
}{
	{"olympiad", olympiadScenarios},
	{"grade-3", grade3Scenarios},
	{"essays", essayScenarios},
	{"modules", moduleScenarios},
}

var (
	oEnum   = sid(releaseOlympiad, "combinatorics-enumeration")
	oGaps   = sid(releaseOlympiad, "counting-gaps")
	oPigeon = sid(releaseOlympiad, "pigeonhole-basic")
	band34  = learnermodel.BandRef{Place: 2, ID: "grades-3-4"}
)

// Olympiad maths for a learner of grade 3 (sections 3.1–3.3).
var olympiadScenarios = []scenario{
	{"new-learner", "6", func(rel map[string]release) golden {
		return golden{modules: []Module{rel[releaseOlympiad].forAge("under_13")}, grade: learnerGrade3}
	}},
	{"queue", "7", olympiadQueue},
	{"consolidation", "3", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.in.Last, g.in.HasLast = failedOn(tid(releaseOlympiad, "gaps-34-d3-1"), oGaps), true
		return g
	}},
	// The one task of counting-gaps in the corridor had its solution seen two days ago: it comes
	// before the unseen ones outside the corridor (the amendment of 2026-11-10 to R03, candidate B).
	{"consolidation-seen-in-the-corridor", "3", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.in.Last, g.in.HasLast = failedOn(tid(releaseOlympiad, "gaps-34-d3-1"), oGaps), true
		g.in.Summary.Tasks = []TaskSolution{{Task: tid(releaseOlympiad, "gaps-34-d1-1"), SolutionSentAt: goldenTime - 2*day}}
		return g
	}},
	{"probe", "4", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.state.Modules[0].Skills[1] = practised(oGaps, -3, 6, goldenTime-2*day)
		g.in.Summary.Modules = []ModuleChoices{{Module: releaseOlympiad, Choices: 11}}
		return g
	}},
	{"lesson-of-a-new-band", "7", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.state.Modules[0].Skills[0] = practised(oEnum, 3, 8, goldenTime-3*day)
		return g
	}},
	{"rotation-above-mastery-band", "7-mastered", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.state.Modules[0].Skills[0] = masteredState(oEnum, 2, 400, goldenTime-3*day)
		g.state.Modules[0].Skills[0].Correction = 3
		g.in.Summary.Skills = append(g.in.Summary.Skills, issued(goldenTime-day, oPigeon)...)
		return g
	}},
	{"map-closed-skill", "2", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		g.in.MapSkill = oPigeon
		return g
	}},
	{"open-run", "1", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		l := lessonNamed(g.modules[0], "pigeonhole-basic-grades-3-4")
		g.in.Runs = append(g.in.Runs, Run{Lesson: l.ID, Skills: l.Skills})
		return g
	}},
	{"bottom-of-the-ladder-new", "9-new", olympiadLadder},
	{"bottom-of-the-ladder", "9", func(rel map[string]release) golden {
		g := olympiadLadder(rel)
		g.in.Summary.Skills = slices.Concat(issued(goldenTime-3*day, oEnum), issued(goldenTime-2*day, oGaps))
		return g
	}},
}

var (
	gTimes     = sid(releaseGrade3, "times-tables")
	gDivision  = sid(releaseGrade3, "division-with-remainder")
	gFractions = sid(releaseGrade3, "fractions-of-a-whole")
	gComparing = sid(releaseGrade3, "comparing-fractions")
	grade3Band = learnermodel.BandRef{Place: 1, ID: "grade-3"}
)

// School maths for grade 3 (sections 3.1, 3.3 and 3.4).
var grade3Scenarios = []scenario{
	{"times-tables-lesson", "6", func(rel map[string]release) golden {
		return golden{modules: []Module{rel[releaseGrade3].forAge("under_13")}, grade: learnerGrade3}
	}},
	{"after-the-lesson-second-pass", "3", func(rel map[string]release) golden {
		m := rel[releaseGrade3].forAge("under_13")
		return golden{modules: []Module{m}, grade: learnerGrade3,
			state: learnermodel.State{Modules: []learnermodel.ModuleState{module(releaseGrade3, 0, grade3Band, practised(gTimes, 0, 6, goldenTime-60))}},
			in: Input{Runs: []Run{passed(lessonNamed(m, "times-tables"), goldenTime-60)},
				Last: Last{Task: tid(releaseGrade3, "times-tables-d5"), InLesson: true, Lesson: lid(releaseGrade3, "times-tables"),
					Results: []SkillResult{{Skill: gTimes, Answered: 1}}},
				HasLast: true,
				Summary: Summary{Skills: issued(goldenTime-3600, gTimes), Tasks: []TaskSolution{
					{Task: tid(releaseGrade3, "times-tables-d1"), SolutionSentAt: goldenTime - 300},
					{Task: tid(releaseGrade3, "times-tables-d3"), SolutionSentAt: goldenTime - 200},
					{Task: tid(releaseGrade3, "times-tables-d5"), SolutionSentAt: goldenTime - 100},
				}}}}
	}},
	{"repetition-on-its-turn", "5", func(rel map[string]release) golden {
		g := grade3Mastered(rel)
		g.in.Summary.Modules = []ModuleChoices{{Module: releaseGrade3, Choices: 13}}
		return g
	}},
	{"queue-off-the-turn", "7", func(rel map[string]release) golden {
		g := grade3Mastered(rel)
		g.in.Summary.Modules = []ModuleChoices{{Module: releaseGrade3, Choices: 14}}
		return g
	}},
	{"all-mastered", "8", grade3AllMastered},
	{"repetition-off-the-turn", "5", func(rel map[string]release) golden {
		g := grade3AllMastered(rel)
		g.state.Modules[0].Skills[1] = masteredState(gDivision, 1, 1, goldenTime-20*day)
		return g
	}},
}

// grade3Mastered is school maths with times-tables mastered and due for review — practised four
// days ago, with a stability of a day —, fractions-of-a-whole mastered and not due, and the
// learner at work on division-with-remainder; every skill handed out, every lesson passed.
func grade3Mastered(rel map[string]release) golden {
	m := rel[releaseGrade3].forAge("under_13")
	var runs []Run
	for _, l := range m.Lessons {
		runs = append(runs, passed(l, goldenTime-10*day))
	}
	return golden{modules: []Module{m}, grade: learnerGrade3,
		state: learnermodel.State{Modules: []learnermodel.ModuleState{module(releaseGrade3, 0.5, grade3Band,
			masteredState(gTimes, 1, 1, goldenTime-4*day), masteredState(gFractions, 1, 400, goldenTime-2*day),
			practised(gDivision, 0, 3, goldenTime-day), practised(gComparing, 0, 2, goldenTime-day))}},
		in: Input{Runs: runs, Summary: Summary{Skills: slices.Concat(
			issued(goldenTime-4*day, gTimes), issued(goldenTime-2*day, gFractions), issued(goldenTime-day, gDivision),
			issued(goldenTime-day+60, gComparing))}}}
}

// Short essays, by the author's order (section 3.3): a step for each lesson in turn, then the
// queue; and a child, who gets no writing-about-a-day.
var essayScenarios = func() []scenario {
	order := []string{"first-sentences", "reading-model-essays", "paragraphs", "endings", "improving-an-essay", "writing-about-a-day"}
	var out []scenario
	for k := range len(order) + 1 {
		for _, age := range []string{"adult", "under_13"} {
			branch := "6a"
			if k == len(order) || age == "under_13" && k == len(order)-1 {
				branch = "7"
			}
			if age == "under_13" && k != len(order)-1 {
				continue
			}
			name := fmt.Sprintf("%s-step-%d", age, k+1)
			out = append(out, scenario{name, branch, func(rel map[string]release) golden {
				m := rel[releaseEssays].forAge(age)
				g := golden{modules: []Module{m}}
				for _, n := range order[:k] {
					if slices.ContainsFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(releaseEssays, n) }) {
						l := lessonNamed(m, n)
						// The lesson is handed out an hour before its run ends, on a day of its own.
						ended := goldenTime - int64(10-len(g.in.Runs))*day
						g.in.Runs = append(g.in.Runs, passed(l, ended))
						g.in.Summary.Skills = slices.DeleteFunc(g.in.Summary.Skills, func(s SkillIssue) bool { return slices.Contains(l.Skills, s.Skill) })
						g.in.Summary.Skills = append(g.in.Summary.Skills, issued(ended-3600, l.Skills...)...)
					}
				}
				return g
			}})
		}
	}
	return out
}()

// Both maths at once (section 3.2), and the AI closed (section 3.3).
var moduleScenarios = []scenario{
	{"both-maths-first-turn", "6", func(rel map[string]release) golden {
		return golden{modules: []Module{rel[releaseOlympiad].forAge("under_13"), rel[releaseGrade3].forAge("under_13")}, grade: learnerGrade3}
	}},
	{"both-maths-consolidation-in-its-module", "3", func(rel map[string]release) golden {
		g := olympiadQueue(rel)
		gr := grade3Mastered(rel)
		g.modules = append(g.modules, gr.modules...)
		g.state.Modules = append(g.state.Modules, gr.state.Modules...)
		g.in.Runs = append(g.in.Runs, gr.in.Runs...)
		g.in.Summary.Skills = append(g.in.Summary.Skills, gr.in.Summary.Skills...)
		g.in.Last, g.in.HasLast = failedOn(tid(releaseOlympiad, "gaps-34-d3-1"), oGaps), true
		return g
	}},
	{"ai-closed", "nothing", func(rel map[string]release) golden {
		o, gr := rel[releaseOlympiad].forAge("under_13"), rel[releaseGrade3].forAge("under_13")
		o.Pool, o.Lessons, gr.Pool, gr.Lessons = nil, nil, nil, nil
		return golden{modules: []Module{o, gr}, grade: learnerGrade3}
	}},
}

// vector writes the line of one vector: its name and branch, the choice, and for a task the bits
// of the P the engine computed of it on each of its primary skills with a verdict.
func vector(t *testing.T, s scenario, in Input) string {
	t.Helper()
	c, err := Choose(in)
	if err != nil {
		t.Fatalf("%s: %v", s.name, err)
	}
	kind := map[ChoiceKind]string{ChoiceLesson: "lesson", ChoiceTask: "task", ChoiceNothing: "nothing"}[c.Kind]
	line := fmt.Sprintf("%s branch=%s %s", s.name, s.branch, kind)
	if c.Kind != ChoiceNothing {
		line += " " + c.ID
	}
	line += " " + reasonNames[c.Reason]
	if c.Kind == ChoiceTask {
		pr, err := prepare(in)
		if err != nil {
			t.Fatal(err)
		}
		for i, m := range in.Modules {
			for j, cd := range m.Pool {
				if cd.Task.ID != c.ID {
					continue
				}
				pt := newPoolTask(cd, pr.p[i][j])
				for _, skill := range slices.Compact(sortedCopy(pt.primaries)) {
					if v, ok := candidateOf(pt, skill, nil); ok {
						line += fmt.Sprintf(" p:%s=%016x", strings.TrimPrefix(skill, m.Graph.Module+"/skill/"), math.Float64bits(v.p))
					}
				}
			}
		}
	}
	return line
}

// TestGoldenChoose holds the choice of every scenario to its vector in testdata/choose.
func TestGoldenChoose(t *testing.T) {
	rel := readRelease(t)
	for _, set := range goldenSets {
		t.Run(set.file, func(t *testing.T) {
			var b strings.Builder
			for _, s := range set.scenarios {
				b.WriteString(vector(t, s, s.input(rel).input(t)) + "\n")
			}
			got := b.String()
			path := filepath.Join("testdata", "choose", set.file+".txt")
			if os.Getenv(updateEnv) == "1" {
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want := string(readFile(t, path))
			if got != want {
				gl, wl := strings.Split(got, "\n"), strings.Split(want, "\n")
				for i := range min(len(gl), len(wl)) {
					if gl[i] != wl[i] {
						t.Fatalf("%s:%d: the choice is\n%s\nnot\n%s", path, i+1, gl[i], wl[i])
					}
				}
				t.Fatalf("%s: %d lines, the file has %d", path, len(gl), len(wl))
			}
		})
	}
}

// TestGoldenCodes holds the code of each scenario to the code of its branch in table 3.5, and
// requires a scenario of every branch: so every branch gives the code of the table, and the rule
// gives each of the nine codes, REASON_CODE_ROTATION at its three places.
func TestGoldenCodes(t *testing.T) {
	rel := readRelease(t)
	seen := map[string]bool{}
	codes := map[ReasonCode]bool{}
	for _, set := range goldenSets {
		for _, s := range set.scenarios {
			k := slices.IndexFunc(branchCodes, func(b branchCode) bool { return b.branch == s.branch })
			if k < 0 {
				t.Fatalf("%s: branch %q is not in table 3.5", s.name, s.branch)
			}
			c, err := Choose(s.input(rel).input(t))
			if err != nil {
				t.Fatal(err)
			}
			if want := branchCodes[k].code; c.Reason != want || (c.Kind == ChoiceNothing) != (want == ReasonUnspecified) {
				t.Errorf("%s/%s: branch %s gives %+v, want the code %s", set.file, s.name, s.branch, c, reasonNames[want])
			}
			seen[s.branch], codes[c.Reason] = true, true
		}
	}
	for _, b := range branchCodes {
		if !seen[b.branch] {
			t.Errorf("no scenario of branch %s", b.branch)
		}
	}
	for code := ReasonUnspecified; code <= ReasonRotation; code++ {
		if !codes[code] {
			t.Errorf("no scenario gives %s", reasonNames[code])
		}
	}
}

// olympiadQueue is olympiad maths past its trial series, at about the start of grade 3, with
// the lessons of combinatorics-enumeration and counting-gaps at grades-3-4 passed and both skills
// handed out, combinatorics-enumeration first.
func olympiadQueue(rel map[string]release) golden {
	m := rel[releaseOlympiad].forAge("under_13")
	return golden{modules: []Module{m}, grade: learnerGrade3,
		state: learnermodel.State{Modules: []learnermodel.ModuleState{module(releaseOlympiad, 2, band34,
			practised(oEnum, 0, 4, goldenTime-3*day), practised(oGaps, 0, 4, goldenTime-2*day))}},
		in: Input{Runs: []Run{passed(lessonNamed(m, "combinatorics-enumeration-grades-3-4"), goldenTime-3*day),
			passed(lessonNamed(m, "counting-gaps-grades-3-4"), goldenTime-2*day)},
			Summary: Summary{Skills: slices.Concat(issued(goldenTime-3*day, oEnum), issued(goldenTime-2*day, oGaps))}}}
}

// olympiadLadder is olympiad maths at a level so low that no skill is within reach.
func olympiadLadder(rel map[string]release) golden {
	return golden{modules: []Module{rel[releaseOlympiad].forAge("under_13")}, grade: learnerGrade3,
		state: learnermodel.State{Modules: []learnermodel.ModuleState{module(releaseOlympiad, -4, band34)}}}
}

// grade3AllMastered is school maths with all four skills mastered and none due for review.
func grade3AllMastered(rel map[string]release) golden {
	g := grade3Mastered(rel)
	g.state.Modules[0].Level = 2
	g.state.Modules[0].Skills = []learnermodel.SkillState{
		masteredState(gTimes, 1, 400, goldenTime-4*day), masteredState(gDivision, 1, 400, goldenTime-3*day),
		masteredState(gFractions, 1, 400, goldenTime-2*day), masteredState(gComparing, 1, 400, goldenTime-day),
	}
	g.in.Summary.Modules = []ModuleChoices{{Module: releaseGrade3, Choices: 14}}
	return g
}
