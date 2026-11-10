package engine

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The states of a skill and the task on it, on the examples of docs/architecture.md, section
// 3.1, and at their edges. The reading of the learner model and the P of the candidates are set
// by hand, as the examples set them: there they are for the example, and the model computes them.

const (
	olympiadMath = "novatrail/olympiad-math"
	grade3Math   = "novatrail/grade-3-math"
	exampleTime  = int64(10_000_000)
)

func sid(module, name string) string { return module + "/skill/" + name }

// notDue is a review date that never comes, for a mastered skill whose repetition is not what a
// test is about: the reading gives every mastered skill a date.
var notDue = learnermodel.ReviewDate{At: math.MaxInt64}

// exampleTask is a task of one item with an answer to choose, on the task's primary skills.
func exampleTask(module, name string, band learnermodel.BandRef, difficulty int32, primary ...string) learnermodel.Task {
	return learnermodel.Task{
		ID: module + "/task/" + name, Band: band, Difficulty: difficulty, Primary: primary,
		Items: []learnermodel.TaskItem{{
			Place: 1, Slots: []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}},
		}},
	}
}

func exampleModule(t *testing.T, module string, bands []learnermodel.Band, skills []learnermodel.Skill, tasks []learnermodel.Task) Module {
	t.Helper()
	g, err := learnermodel.NewGraph(module, bands, skills, nil, tasks)
	if err != nil {
		t.Fatal(err)
	}
	m := Module{Graph: g}
	for _, task := range tasks {
		m.Pool = append(m.Pool, Candidate{Task: task})
	}
	return m
}

// olympiad is the graph of the first batch of olympiad maths in section 3.1: two skills at tier
// 1 and pigeonhole-basic after combinatorics-enumeration, three bands, and on each skill and band
// tasks of difficulty 1, 3 and 5 — of counting-gaps at grades-1-2, 2, 3 and 5.
func olympiad(t *testing.T) Module {
	t.Helper()
	bands := []learnermodel.Band{
		{Place: 1, ID: "grades-1-2", Grades: []int32{1, 2}},
		{Place: 2, ID: "grades-3-4", Grades: []int32{3, 4}},
		{Place: 3, ID: "grades-5-6", Grades: []int32{5, 6}},
	}
	enum, gaps, pigeon := sid(olympiadMath, "combinatorics-enumeration"), sid(olympiadMath, "counting-gaps"), sid(olympiadMath, "pigeonhole-basic")
	skills := []learnermodel.Skill{
		{ID: pigeon, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{enum}, Tier: 2},
		{ID: gaps, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: enum, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
	}
	var tasks []learnermodel.Task
	for _, s := range []struct{ code, skill string }{{"enum", enum}, {"gaps", gaps}, {"pigeon", pigeon}} {
		for _, b := range bands {
			band := learnermodel.BandRef{Place: b.Place, ID: b.ID}
			difficulties := []int32{1, 3, 5}
			if s.code == "gaps" && b.Place == 1 {
				difficulties = []int32{2, 3, 5}
			}
			for _, d := range difficulties {
				name := fmt.Sprintf("%s-%d%d-d%d-1", s.code, 2*b.Place-1, 2*b.Place, d)
				tasks = append(tasks, exampleTask(olympiadMath, name, band, d, s.skill))
			}
		}
	}
	return exampleModule(t, olympiadMath, bands, skills, tasks)
}

// grade3 is the graph of school maths for grade 3 in section 3.1: times-tables, then
// division-with-remainder and fractions-of-a-whole, then comparing-fractions; one band, and
// tasks of difficulty 1, 3 and 5 on each skill.
func grade3(t *testing.T) Module {
	t.Helper()
	bands := []learnermodel.Band{{Place: 1, ID: "grade-3", Grades: []int32{3}}}
	tt, div, frac, comp := sid(grade3Math, "times-tables"), sid(grade3Math, "division-with-remainder"),
		sid(grade3Math, "fractions-of-a-whole"), sid(grade3Math, "comparing-fractions")
	skills := []learnermodel.Skill{
		{ID: comp, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{frac}, Tier: 3},
		{ID: div, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{tt}, Tier: 2},
		{ID: frac, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{tt}, Tier: 2},
		{ID: tt, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
	}
	var tasks []learnermodel.Task
	for _, s := range []string{"times-tables", "division-with-remainder", "fractions-of-a-whole", "comparing-fractions"} {
		for _, d := range []int32{1, 3, 5} {
			tasks = append(tasks, exampleTask(grade3Math, fmt.Sprintf("%s-d%d", s, d), learnermodel.BandRef{Place: 1, ID: "grade-3"}, d, sid(grade3Math, s)))
		}
	}
	return exampleModule(t, grade3Math, bands, skills, tasks)
}

// world is an input over one module, with the P of its candidates beside it, by place in the pool.
type world struct {
	t  *testing.T
	in Input
	p  []learnermodel.TaskReadout
}

func newWorld(t *testing.T, m Module) *world {
	w := &world{t: t, in: Input{
		Time:    exampleTime,
		Readout: learnermodel.Readout{Modules: []learnermodel.ModuleReadout{{Module: m.Graph.Module}}},
		Modules: []Module{m},
	}, p: make([]learnermodel.TaskReadout, len(m.Pool))}
	for j := range m.Pool {
		w.setP(j, 0.5)
	}
	return w
}

func (w *world) setP(j int, p float64) {
	w.p[j] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1, Probability: p, HasVerdict: true}}}
}

// px is the P of a task by its name, or of every task whose name begins with a prefix ending in
// "*".
type px struct {
	name string
	p    float64
}

func (w *world) probabilities(ps []px) {
	for j, c := range w.in.Modules[0].Pool {
		name := c.Task.ID[strings.LastIndex(c.Task.ID, "/")+1:]
		for _, x := range ps {
			if x.name == name || strings.HasSuffix(x.name, "*") && strings.HasPrefix(name, strings.TrimSuffix(x.name, "*")) {
				w.setP(j, x.p)
			}
		}
	}
}

// skill sets what the reading gives of a skill.
func (w *world) skill(r learnermodel.SkillReadout) {
	mr := &w.in.Readout.Modules[0]
	mr.Skills = slices.DeleteFunc(mr.Skills, func(s learnermodel.SkillReadout) bool { return s.Skill == r.Skill })
	mr.Skills = append(mr.Skills, r)
	slices.SortFunc(mr.Skills, func(a, b learnermodel.SkillReadout) int { return strings.Compare(a.Skill, b.Skill) })
}

func (w *world) views(second bool) []skillView {
	return ready(w.in, [][]learnermodel.TaskReadout{w.p}).skills(0, second)
}

func (w *world) view(id string, second bool) skillView {
	w.t.Helper()
	vs := w.views(second)
	i := slices.IndexFunc(vs, func(v skillView) bool { return v.id == id })
	if i < 0 {
		w.t.Fatalf("no skill %s among %v", id, ids(vs))
	}
	return vs[i]
}

func ids(vs []skillView) []string {
	var out []string
	for _, v := range vs {
		out = append(out, v.id)
	}
	return out
}

func candidateIDs(cs []candidate) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.id[strings.LastIndex(c.id, "/")+1:])
	}
	return out
}

func taskName(t *testing.T, s skillView, easiest bool) string {
	t.Helper()
	c, ok := s.task(easiest)
	if !ok {
		t.Fatalf("skill %s has no task", s.id)
	}
	return c.id[strings.LastIndex(c.id, "/")+1:]
}

// Olympiad maths, step 1: no evidence, the trial series is on. Both skills of tier 1 are open and
// within reach, pigeonhole-basic is closed, and the task on combinatorics-enumeration is the one
// nearest to 0.775.
func TestOlympiadStart(t *testing.T) {
	w := newWorld(t, olympiad(t))
	w.probabilities([]px{
		{"enum-12-*", 0.95}, {"enum-34-d1-1", 0.78}, {"enum-34-d3-1", 0.6}, {"enum-34-d5-1", 0.4}, {"enum-56-*", 0.2},
		{"gaps-12-*", 0.9}, {"gaps-34-*", 0.5}, {"gaps-56-*", 0.2}, {"pigeon-*", 0.8},
	})
	vs := w.views(false)
	if want := []string{sid(olympiadMath, "combinatorics-enumeration"), sid(olympiadMath, "counting-gaps"), sid(olympiadMath, "pigeonhole-basic")}; !slices.Equal(ids(vs), want) {
		t.Fatalf("the order of the graph is %v, want %v", ids(vs), want)
	}
	enum, gaps, pigeon := vs[0], vs[1], vs[2]
	if !enum.open || !enum.reachable || !gaps.open || !gaps.reachable {
		t.Errorf("tier 1: %+v, %+v; want both open and within reach", enum, gaps)
	}
	if pigeon.open || pigeon.reachable {
		t.Errorf("pigeonhole-basic: %+v; want it closed and out of reach", pigeon)
	}
	if got := taskName(t, enum, false); got != "enum-34-d1-1" {
		t.Errorf("the task on combinatorics-enumeration is %s, want enum-34-d1-1", got)
	}
}

// Olympiad maths, step 4: after the mistake in gaps-34-d3-1 the corridor has moved to easier
// tasks, and the task on counting-gaps is gaps-34-d1-1.
func TestOlympiadMistake(t *testing.T) {
	w := newWorld(t, olympiad(t))
	gaps := sid(olympiadMath, "counting-gaps")
	w.skill(learnermodel.SkillReadout{Skill: gaps, EvidenceCount: 7, LastPractice: exampleTime - 60})
	w.probabilities([]px{{"gaps-12-*", 0.92}, {"gaps-34-d1-1", 0.76}, {"gaps-34-d3-1", 0.55}, {"gaps-34-d5-1", 0.3}, {"gaps-56-*", 0.1}})
	v := w.view(gaps, false)
	if !v.reachable || v.paused() {
		t.Errorf("counting-gaps: %+v; want it within reach", v)
	}
	if got := taskName(t, v, false); got != "gaps-34-d1-1" {
		t.Errorf("the task on counting-gaps is %s, want gaps-34-d1-1", got)
	}
}

// Olympiad maths, step 5: after two more mistakes even the easiest task of counting-gaps has a P
// below 0.70, and the skill is on pause; the task of a probe is the easiest, gaps-12-d2-1, and a
// right probe puts the skill within reach again.
func TestOlympiadPause(t *testing.T) {
	w := newWorld(t, olympiad(t))
	gaps := sid(olympiadMath, "counting-gaps")
	w.skill(learnermodel.SkillReadout{Skill: gaps, EvidenceCount: 9, LastPractice: exampleTime - 60})
	w.probabilities([]px{{"gaps-12-d2-1", 0.66}, {"gaps-12-d3-1", 0.5}, {"gaps-12-d5-1", 0.3}, {"gaps-34-*", 0.2}, {"gaps-56-*", 0.1}})
	v := w.view(gaps, false)
	if v.reachable || !v.paused() || v.notYet() {
		t.Fatalf("counting-gaps: %+v; want it on pause", v)
	}
	if got := taskName(t, v, true); got != "gaps-12-d2-1" {
		t.Errorf("the easiest task on counting-gaps is %s, want gaps-12-d2-1", got)
	}
	// The task on the skill, out of the corridor, is the one nearest to 0.775: the same here.
	if got := taskName(t, v, false); got != "gaps-12-d2-1" {
		t.Errorf("the task on counting-gaps is %s, want gaps-12-d2-1", got)
	}

	for _, tc := range []struct {
		name   string
		issue  SkillIssue
		within bool
	}{
		{"right probe", SkillIssue{LastCode: ReasonProbe, LastResult: SkillResult{Answered: 1, FirstRight: 1}, HasLast: true}, true},
		{"right repetition", SkillIssue{LastCode: ReasonRepetition, LastResult: SkillResult{Answered: 2, FirstRight: 1}, HasLast: true}, true},
		{"failed probe", SkillIssue{LastCode: ReasonProbe, LastResult: SkillResult{Answered: 3, FirstRight: 1}, HasLast: true}, false},
		{"probe with no answer", SkillIssue{LastCode: ReasonProbe, HasLast: true}, false},
		{"right task of the queue", SkillIssue{LastCode: ReasonQueue, LastResult: SkillResult{Answered: 1, FirstRight: 1}, HasLast: true}, false},
		{"right probe, not the last", SkillIssue{LastCode: ReasonProbe, LastResult: SkillResult{Answered: 1, FirstRight: 1}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.issue.Skill = gaps
			w.in.Summary.Skills = []SkillIssue{tc.issue}
			if v := w.view(gaps, false); v.reachable != tc.within || v.paused() == tc.within {
				t.Errorf("counting-gaps: %+v; want within reach %v", v, tc.within)
			}
		})
	}
}

// Olympiad maths, step 6: "mastered" of combinatorics-enumeration opens pigeonhole-basic, and the
// skill is mastered for the queue while its recommended task is not above grades-3-4.
func TestOlympiadMastered(t *testing.T) {
	enum, pigeon := sid(olympiadMath, "combinatorics-enumeration"), sid(olympiadMath, "pigeonhole-basic")
	for _, tc := range []struct {
		name     string
		ps       []px
		forQueue bool
	}{
		{"recommended at grades-3-4", []px{{"enum-12-*", 0.97}, {"enum-34-d3-1", 0.78}, {"enum-34-d1-1", 0.9}, {"enum-56-*", 0.4}}, true},
		{"recommended at grades-5-6", []px{{"enum-12-*", 0.98}, {"enum-34-*", 0.93}, {"enum-56-d1-1", 0.79}, {"enum-56-d3-1", 0.6}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := newWorld(t, olympiad(t))
			w.skill(learnermodel.SkillReadout{
				Skill: enum, EvidenceCount: 12, Mastered: true, Review: notDue, MasteryBand: learnermodel.BandRef{Place: 2, ID: "grades-3-4"},
			})
			w.probabilities(tc.ps)
			if v := w.view(enum, false); v.masteredForQueue != tc.forQueue {
				t.Errorf("combinatorics-enumeration: recommended %+v, mastered for the queue %v; want %v", v.recommended, v.masteredForQueue, tc.forQueue)
			}
			if v := w.view(pigeon, false); !v.open {
				t.Errorf("pigeonhole-basic: %+v; want it open", v)
			}
		})
	}
}

// School maths for grade 3, step 1: only times-tables is open; when even times-tables-d1 is too
// hard, it is out of reach and, with no evidence, not reached yet.
func TestGrade3Start(t *testing.T) {
	w := newWorld(t, grade3(t))
	w.probabilities([]px{{"times-tables-d1", 0.6}, {"times-tables-d3", 0.4}, {"times-tables-d5", 0.2}})
	vs := w.views(false)
	want := []string{sid(grade3Math, "times-tables"), sid(grade3Math, "division-with-remainder"), sid(grade3Math, "fractions-of-a-whole"), sid(grade3Math, "comparing-fractions")}
	if !slices.Equal(ids(vs), want) {
		t.Fatalf("the order of the graph is %v, want %v", ids(vs), want)
	}
	for _, v := range vs[1:] {
		if v.open {
			t.Errorf("%s is open, want it closed", v.id)
		}
	}
	if tt := vs[0]; !tt.open || tt.reachable || !tt.notYet() || tt.paused() {
		t.Errorf("times-tables: %+v; want it open and not reached yet", tt)
	}
	if got := taskName(t, vs[0], false); got != "times-tables-d1" {
		t.Errorf("the task on times-tables is %s, want times-tables-d1, the nearest", got)
	}
}

// School maths for grade 3, step 2: the solution of times-tables-d3 was shown, and for 24 hours
// it is no candidate; then it comes first again, in the corridor, before times-tables-d1 and
// times-tables-d5 outside it; in the second pass it comes back on pause, after them.
func TestGrade3Solution(t *testing.T) {
	w := newWorld(t, grade3(t))
	tt := sid(grade3Math, "times-tables")
	w.probabilities([]px{{"times-tables-d1", 0.9}, {"times-tables-d3", 0.78}, {"times-tables-d5", 0.6}})
	for _, tc := range []struct {
		name   string
		sent   int64
		second bool
		want   []string
	}{
		{"an hour ago", exampleTime - 3600, false, []string{"times-tables-d1", "times-tables-d5"}},
		{"a second short of 24 hours", exampleTime - solutionPause + 1, false, []string{"times-tables-d1", "times-tables-d5"}},
		{"24 hours ago", exampleTime - solutionPause, false, []string{"times-tables-d3", "times-tables-d1", "times-tables-d5"}},
		{"an hour ago, second pass", exampleTime - 3600, true, []string{"times-tables-d1", "times-tables-d5", "times-tables-d3"}},
		{"in the future", exampleTime + 10, false, []string{"times-tables-d1", "times-tables-d5"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w.in.Summary.Tasks = []TaskSolution{{Task: grade3Math + "/task/times-tables-d3", SolutionSentAt: tc.sent}}
			v := w.view(tt, tc.second)
			if got := candidateIDs(v.candidates); !slices.Equal(got, tc.want) {
				t.Errorf("candidates %v, want %v", got, tc.want)
			}
			// The recommended task does not see the pause.
			if got := v.recommended.id; got != grade3Math+"/task/times-tables-d3" {
				t.Errorf("recommended %s, want times-tables-d3", got)
			}
		})
	}
}

// School maths for grade 3, step 3: "mastered" of times-tables opens both skills of tier 2, and
// comparing-fractions waits for fractions-of-a-whole.
func TestGrade3Mastered(t *testing.T) {
	w := newWorld(t, grade3(t))
	w.skill(learnermodel.SkillReadout{Skill: sid(grade3Math, "times-tables"), EvidenceCount: 10, Mastered: true, Review: notDue, MasteryBand: learnermodel.BandRef{Place: 1, ID: "grade-3"}})
	var open []string
	for _, v := range w.views(false) {
		if v.open {
			open = append(open, v.id)
		}
	}
	if want := []string{sid(grade3Math, "times-tables"), sid(grade3Math, "division-with-remainder"), sid(grade3Math, "fractions-of-a-whole")}; !slices.Equal(open, want) {
		t.Errorf("open %v, want %v", open, want)
	}
}

// A skill with evidence of its own is open whatever came of its prerequisites; a prerequisite
// with a rubric holds it while it has a lesson not passed, and an anchor skill with no reading,
// which Choose refuses, holds it too.
func TestOpen(t *testing.T) {
	module := "example/open"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	a, b, essay := sid(module, "a"), sid(module, "b"), sid(module, "essay")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: essay, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}, Tier: 2},
		{ID: sid(module, "c"), Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{essay}, Tier: 2},
		{ID: sid(module, "d"), Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{"core.counting"}, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "a-1", band, 1, a), exampleTask(module, "b-1", band, 1, b)})
	m.Lessons = []Lesson{{ID: module + "/lesson/essay", Skills: []string{essay}}}
	w := newWorld(t, m)
	if v := w.view(b, false); v.open {
		t.Fatalf("b is open with a not mastered, want it closed")
	}
	w.skill(learnermodel.SkillReadout{Skill: b, EvidenceCount: 1})
	if v := w.view(b, false); !v.open {
		t.Errorf("b with evidence of its own is closed, want it open")
	}
	// Mastered by the model: a rubric skill never is, but a reading that says so counts nothing.
	w.skill(learnermodel.SkillReadout{Skill: essay, Mastered: true, Review: notDue})
	for _, id := range []string{sid(module, "c"), sid(module, "d")} {
		if v := w.view(id, false); v.open {
			t.Errorf("%s is open, want it held", id)
		}
	}
	if slices.Contains(ids(w.views(false)), essay) {
		t.Errorf("the skills with a verdict include %s, a skill with a rubric", essay)
	}
}

// A candidate of a skill: a task where it is primary, of the task or of an item, with an item with
// a verdict where it is primary, and whose every other primary skill is open; its P is the mean P
// of those items.
func TestCandidates(t *testing.T) {
	module := "example/candidates"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	a, b, c := sid(module, "a"), sid(module, "b"), sid(module, "c")
	slot := []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}}
	both := exampleTask(module, "both", band, 1, a, b) // b is closed: no candidate of a
	ownItem := learnermodel.Task{ID: module + "/task/own", Band: band, Difficulty: 1, Primary: []string{c}, Items: []learnermodel.TaskItem{
		{Place: 1, OwnSkills: true, Primary: []string{a}, Slots: slot},
		{Place: 2, Slots: slot},
		{Place: 3, OwnSkills: true, Primary: []string{a}, Slots: slot},
		{Place: 4, OwnSkills: true, Primary: []string{a}},
	}}
	noVerdict := exampleTask(module, "no-verdict", band, 1, a)
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}, Tier: 2},
		{ID: c, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
	}, []learnermodel.Task{both, ownItem, noVerdict})
	m.Pool[1].Traps = []ItemTraps{{Place: 2, Traps: []string{"t/c-only"}}, {Place: 3, Traps: []string{"t/a"}}}
	w := newWorld(t, m)
	w.p[1] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{
		{Place: 1, Probability: 0.5, HasVerdict: true},
		{Place: 2, Probability: 0.1, HasVerdict: true},
		{Place: 3, Probability: 0.75, HasVerdict: true},
		{Place: 4},
	}}
	w.p[2] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1}}}
	w.skill(learnermodel.SkillReadout{Skill: a, Traps: []learnermodel.TrapCount{{ID: "t/a", Count: 2}, {ID: "t/c-only", Count: 5}}})

	v := w.view(a, false)
	if got := candidateIDs(v.candidates); !slices.Equal(got, []string{"own"}) {
		t.Fatalf("candidates of a %v, want [own]", got)
	}
	if got := v.candidates[0]; got.p != 0.625 || !slices.Equal(got.traps, []int{1}) {
		t.Errorf("own on a: P %v, traps %v; want 0.625 of items 1 and 3, and only the trap of item 3", got.p, got.traps)
	}
	// The task's primary skill counts only for the items that take it.
	if got := candidateIDs(w.view(c, false).candidates); !slices.Equal(got, []string{"own"}) {
		t.Errorf("candidates of c %v, want [own], by item 2", got)
	}
	if got := w.view(c, false).candidates[0].p; got != 0.1 {
		t.Errorf("own on c: P %v, want 0.1", got)
	}
	w.skill(learnermodel.SkillReadout{Skill: a, Mastered: true, Review: notDue})
	if got := candidateIDs(w.view(a, false).candidates); !slices.Equal(got, []string{"own", "both"}) {
		t.Errorf("candidates of a with b open %v, want [own both]", got)
	}
}

// The edge of reach is 0.70, inclusive, and a skill with no candidates is out of reach.
func TestReach(t *testing.T) {
	tt := sid(grade3Math, "times-tables")
	for _, tc := range []struct {
		p      float64
		within bool
	}{{0.70, true}, {math.Float64frombits(math.Float64bits(0.70) - 1), false}, {1, true}, {0, false}} {
		w := newWorld(t, grade3(t))
		w.probabilities([]px{{"times-tables-*", tc.p}})
		if v := w.view(tt, false); v.reachable != tc.within {
			t.Errorf("P %v: within reach %v, want %v", tc.p, v.reachable, tc.within)
		}
	}
	w := newWorld(t, grade3(t))
	w.probabilities([]px{{"times-tables-*", 0.9}})
	w.in.Modules[0].Pool = w.in.Modules[0].Pool[3:]
	w.p = w.p[3:]
	w.in.Summary.Skills = []SkillIssue{{Skill: tt, LastCode: ReasonProbe, LastResult: SkillResult{Answered: 1, FirstRight: 1}, HasLast: true}}
	if v := w.view(tt, false); v.reachable || v.hasRecommended || !v.notYet() {
		t.Errorf("times-tables with no task: %+v; want it out of reach, even after a right probe", v)
	}
}

// The second pass takes back the tasks on pause, after the others, the one in the corridor first,
// but not the task of the last closed attempt; a template is never on pause.
func TestSecondPass(t *testing.T) {
	w := newWorld(t, grade3(t))
	tt := sid(grade3Math, "times-tables")
	w.probabilities([]px{{"times-tables-d1", 0.9}, {"times-tables-d3", 0.78}, {"times-tables-d5", 0.6}})
	w.in.Summary.Tasks = []TaskSolution{
		{Task: grade3Math + "/task/times-tables-d1", SolutionSentAt: exampleTime - 100},
		{Task: grade3Math + "/task/times-tables-d3", SolutionSentAt: exampleTime - 50},
		{Task: grade3Math + "/task/times-tables-d5", SolutionSentAt: exampleTime - 10},
	}
	w.in.Last, w.in.HasLast = Last{Task: grade3Math + "/task/times-tables-d5"}, true
	if got := candidateIDs(w.view(tt, false).candidates); len(got) != 0 {
		t.Errorf("first pass: candidates %v, want none", got)
	}
	if got := candidateIDs(w.view(tt, true).candidates); !slices.Equal(got, []string{"times-tables-d3", "times-tables-d1"}) {
		t.Errorf("second pass: candidates %v, want [times-tables-d3 times-tables-d1]", got)
	}
	w.in.Modules[0].Pool[1].Task.Template = true // times-tables-d3
	if got := candidateIDs(w.view(tt, false).candidates); !slices.Equal(got, []string{"times-tables-d3"}) {
		t.Errorf("a template: candidates %v, want [times-tables-d3]", got)
	}
}

// The order of the input's lists changes no state and no candidate: the pool, the graph, the
// reading's skills and their traps, the items of a task and their P, the traps of the items and
// the rows of the summary.
func TestSkillsOrder(t *testing.T) {
	w := newWorld(t, olympiad(t))
	enum, gaps := sid(olympiadMath, "combinatorics-enumeration"), sid(olympiadMath, "counting-gaps")
	w.probabilities([]px{{"enum-12-*", 0.9}, {"enum-34-*", 0.75}, {"gaps-12-d2-1", 0.8}, {"gaps-34-*", 0.72}, {"pigeon-*", 0.7}})
	// A task of three items on counting-gaps, with a trap on each, and P whose sum depends on
	// their order unless it goes by place.
	slot := []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}}
	j := slices.IndexFunc(w.in.Modules[0].Pool, func(c Candidate) bool { return strings.HasSuffix(c.Task.ID, "/gaps-34-d1-1") })
	w.in.Modules[0].Pool[j].Task.Items = []learnermodel.TaskItem{{Place: 1, Slots: slot}, {Place: 2, Slots: slot}, {Place: 3, Slots: slot}}
	w.in.Modules[0].Pool[j].Traps = []ItemTraps{{Place: 1, Traps: []string{"t/1"}}, {Place: 2, Traps: []string{"t/2"}}, {Place: 3, Traps: []string{"t/3"}}}
	w.p[j] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{
		{Place: 1, Probability: 0.1, HasVerdict: true}, {Place: 2, Probability: 0.2, HasVerdict: true}, {Place: 3, Probability: 0.3, HasVerdict: true},
	}}
	w.skill(learnermodel.SkillReadout{Skill: enum, EvidenceCount: 3, Mastered: true, Review: notDue, MasteryBand: learnermodel.BandRef{Place: 2}})
	w.skill(learnermodel.SkillReadout{Skill: gaps, EvidenceCount: 3, Traps: []learnermodel.TrapCount{
		{ID: "t/1", Count: 2, Last: 5}, {ID: "t/2", Count: 3, Last: 1}, {ID: "t/3", Count: 2, Last: 5}, {ID: "t/4", Count: 1},
	}})
	w.in.Summary.Tasks = []TaskSolution{
		{Task: olympiadMath + "/task/enum-34-d1-1", SolutionSentAt: exampleTime - solutionPause},
		{Task: olympiadMath + "/task/enum-34-d3-1", SolutionSentAt: exampleTime - 2*solutionPause},
	}
	w.in.Summary.Skills = []SkillIssue{
		{Skill: gaps, LastCode: ReasonProbe, LastResult: SkillResult{Answered: 1, FirstRight: 1}, HasLast: true},
		{Skill: enum, LastIssued: exampleTime - 10, Issued: true},
	}
	want := viewsText(w.views(false))
	if !strings.Contains(want, "traps:[0 1 2]") {
		t.Fatalf("the views have no task catching three frequent traps:\n%s", want)
	}

	rng := rand.New(rand.NewPCG(3, 4))
	for range 20 {
		in := w.in
		m := in.Modules[0]
		m.Graph.Skills = shuffled(rng, m.Graph.Skills)
		order := rng.Perm(len(w.p))
		m.Pool = make([]Candidate, len(order))
		p := make([]learnermodel.TaskReadout, len(order))
		for k, j := range order {
			c := w.in.Modules[0].Pool[j]
			c.Task.Items = shuffled(rng, c.Task.Items)
			c.Traps = shuffled(rng, c.Traps)
			m.Pool[k] = c
			p[k] = learnermodel.TaskReadout{Items: shuffled(rng, w.p[j].Items)}
		}
		in.Modules = []Module{m}
		mr := in.Readout.Modules[0]
		mr.Skills = shuffled(rng, mr.Skills)
		for k := range mr.Skills {
			mr.Skills[k].Traps = shuffled(rng, mr.Skills[k].Traps)
		}
		in.Readout.Modules = []learnermodel.ModuleReadout{mr}
		in.Summary.Tasks = shuffled(rng, in.Summary.Tasks)
		in.Summary.Skills = shuffled(rng, in.Summary.Skills)
		if got := viewsText((&world{t: t, in: in, p: p}).views(false)); got != want {
			t.Fatalf("in another order:\n%s\nwant\n%s", got, want)
		}
	}
}

// viewsText writes the views out to compare, without the traps of the reading: a copy of the
// input in its order, which only frequentTraps reads.
func viewsText(vs []skillView) string {
	for i := range vs {
		vs[i].read.Traps = nil
	}
	return fmt.Sprintf("%+v", vs)
}

// shuffled gives a shuffled copy of s.
func shuffled[T any](rng *rand.Rand, s []T) []T {
	c := slices.Clone(s)
	rng.Shuffle(len(c), func(i, j int) { c[i], c[j] = c[j], c[i] })
	return c
}

func TestFailedAndRight(t *testing.T) {
	for _, tc := range []struct {
		r             SkillResult
		failed, right bool
	}{
		{SkillResult{Answered: 2, FirstRight: 1}, false, true},
		{SkillResult{Answered: 3, FirstRight: 1}, true, false},
		{SkillResult{Answered: 1, FirstRight: 0}, true, false},
		{SkillResult{Answered: 2, FirstRight: 1, SolutionShown: true}, false, true},
		{SkillResult{Answered: 2, FirstRight: 0, SolutionShown: true}, true, false},
		{SkillResult{}, false, false},
		{SkillResult{SolutionShown: true}, true, false},
		{SkillResult{Answered: math.MaxInt32, FirstRight: math.MaxInt32/2 + 1}, false, true},
		{SkillResult{Answered: math.MaxInt32, FirstRight: math.MaxInt32 / 2}, true, false},
	} {
		if got := failed(tc.r); got != tc.failed {
			t.Errorf("failed(%+v) = %v, want %v", tc.r, got, tc.failed)
		}
		if got := right(tc.r); got != tc.right {
			t.Errorf("right(%+v) = %v, want %v", tc.r, got, tc.right)
		}
	}
}

func TestPausedSince(t *testing.T) {
	for _, tc := range []struct {
		now, sent int64
		paused    bool
	}{
		{exampleTime, exampleTime - solutionPause, false},
		{exampleTime, exampleTime - solutionPause + 1, true},
		{exampleTime, exampleTime + 1, true},
		{math.MinInt64, math.MinInt64, true},
		{math.MinInt64 + solutionPause - 1, math.MinInt64, true},
		{math.MinInt64 + solutionPause, math.MinInt64, false},
		{math.MaxInt64, math.MinInt64, false},
		{math.MaxInt64, math.MaxInt64, true},
	} {
		if got := pausedSince(tc.now, tc.sent); got != tc.paused {
			t.Errorf("pausedSince(%d, %d) = %v, want %v", tc.now, tc.sent, got, tc.paused)
		}
	}
}

// The task of the last closed attempt whose solution went to the device is no candidate for the
// next choice, however long ago the solution went (R82 spec): past the pause of 24 hours it would
// otherwise come first by the corridor. Found by invariant I4 of sim's gate on candidate A of T31.
func TestLastSeenNotNext(t *testing.T) {
	w := newWorld(t, grade3(t))
	tt := sid(grade3Math, "times-tables")
	w.probabilities([]px{{"times-tables-d1", 0.9}, {"times-tables-d3", 0.78}, {"times-tables-d5", 0.6}})
	w.in.Summary.Tasks = []TaskSolution{{Task: grade3Math + "/task/times-tables-d3", SolutionSentAt: exampleTime - 2*solutionPause}}
	w.in.Last, w.in.HasLast = Last{Task: grade3Math + "/task/times-tables-d3"}, true
	for _, second := range []bool{false, true} {
		if got := candidateIDs(w.view(tt, second).candidates); !slices.Equal(got, []string{"times-tables-d1", "times-tables-d5"}) {
			t.Errorf("second pass %v: candidates %v, want [times-tables-d1 times-tables-d5]", second, got)
		}
	}
	w.in.HasLast = false
	if got := candidateIDs(w.view(tt, false).candidates); !slices.Equal(got, []string{"times-tables-d3", "times-tables-d1", "times-tables-d5"}) {
		t.Errorf("another last attempt: candidates %v, want times-tables-d3 first, in the corridor", got)
	}
}
