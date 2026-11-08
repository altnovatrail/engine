package engine

import (
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The branches of the rule in one module, on the examples of docs/architecture.md, section 3.1,
// and at their edges, over the worlds of skills_internal_test.go.

// trial sets the two flags of the module's trial series.
func (w *world) trial(over, lastIn bool) {
	w.in.Readout.Modules[0].TrialOver, w.in.Readout.Modules[0].LastInTrial = over, lastIn
}

// choices sets n, the module's choices before this one.
func (w *world) choices(n int64) {
	w.in.Summary.Modules = []ModuleChoices{{Module: w.in.Modules[0].Graph.Module, Choices: n}}
}

// issue sets when a skill was last handed out.
func (w *world) issue(skill string, at int64) {
	for i := range w.in.Summary.Skills {
		if w.in.Summary.Skills[i].Skill == skill {
			w.in.Summary.Skills[i].Issued, w.in.Summary.Skills[i].LastIssued = true, at
			return
		}
	}
	w.in.Summary.Skills = append(w.in.Summary.Skills, SkillIssue{Skill: skill, Issued: true, LastIssued: at})
}

// last sets the last closed attempt of the profile: a task of the module by name, and its results.
func (w *world) last(task string, left bool, results ...SkillResult) {
	w.in.Last = Last{Task: w.in.Modules[0].Graph.Module + "/task/" + task, Left: left, Results: results}
	w.in.HasLast = true
}

func (w *world) choose(second bool) (Choice, bool) {
	return prepared{in: w.in, p: [][]learnermodel.TaskReadout{w.p}}.chooseIn(0, second)
}

// want checks the choice of the first pass: a task by name and a code, or nothing when task is
// empty.
func (w *world) want(task string, code ReasonCode) {
	w.t.Helper()
	got, ok := w.choose(false)
	if task == "" {
		if ok {
			w.t.Errorf("chose %+v, want nothing", got)
		}
		return
	}
	want := Choice{Kind: ChoiceTask, ID: w.in.Modules[0].Graph.Module + "/task/" + task, Reason: code}
	if !ok || got != want {
		w.t.Errorf("chose %+v (%v), want %+v", got, ok, want)
	}
}

func wrong(skill string) SkillResult   { return SkillResult{Skill: skill, Answered: 1} }
func rightOn(skill string) SkillResult { return SkillResult{Skill: skill, Answered: 1, FirstRight: 1} }

// The olympiad maths of section 3.1, step by step.
func olympiadWorld(t *testing.T) *world {
	w := newWorld(t, olympiad(t))
	w.probabilities([]px{
		{"enum-12-*", 0.95}, {"enum-34-d1-1", 0.78}, {"enum-34-d3-1", 0.6}, {"enum-34-d5-1", 0.4}, {"enum-56-*", 0.2},
		{"gaps-12-*", 0.9}, {"gaps-34-d1-1", 0.76}, {"gaps-34-d3-1", 0.55}, {"gaps-34-d5-1", 0.3}, {"gaps-56-*", 0.1},
		{"pigeon-*", 0.8},
	})
	return w
}

var (
	enumSkill   = sid(olympiadMath, "combinatorics-enumeration")
	gapsSkill   = sid(olympiadMath, "counting-gaps")
	pigeonSkill = sid(olympiadMath, "pigeonhole-basic")
)

// Steps 1–3: in the trial series new skills come first, then the queue; no mistake is worked on.
func TestOlympiadSeries(t *testing.T) {
	w := olympiadWorld(t)
	w.want("enum-34-d1-1", ReasonNewSkill)
	w.issue(enumSkill, exampleTime-100)
	w.want("gaps-34-d1-1", ReasonNewSkill)
	w.issue(gapsSkill, exampleTime-50)
	w.want("enum-34-d1-1", ReasonQueue)

	w.last("enum-34-d3-1", false, wrong(enumSkill))
	w.issue(enumSkill, exampleTime-10)
	w.want("gaps-34-d1-1", ReasonQueue) // in the series
	w.trial(true, true)
	w.want("gaps-34-d1-1", ReasonQueue) // on the series' last answer
}

// Step 4: past the series, a mistake keeps the next choice on its skill.
func TestOlympiadConsolidation(t *testing.T) {
	w := olympiadWorld(t)
	w.issue(enumSkill, exampleTime-100)
	w.issue(gapsSkill, exampleTime-10)
	w.trial(true, false)
	w.last("gaps-34-d3-1", false, wrong(gapsSkill))
	w.want("gaps-34-d1-1", ReasonConsolidation)

	t.Run("right", func(t *testing.T) {
		w.last("gaps-34-d3-1", false, rightOn(gapsSkill))
		w.want("enum-34-d1-1", ReasonQueue)
	})
	t.Run("left", func(t *testing.T) {
		w.last("gaps-34-d3-1", true, wrong(gapsSkill))
		w.want("enum-34-d1-1", ReasonQueue)
	})
	t.Run("of another module", func(t *testing.T) {
		w.last("gaps-34-d3-1", false, wrong(gapsSkill))
		w.in.Last.Task = "novatrail/grade-3-math/task/times-tables-d1"
		w.want("enum-34-d1-1", ReasonQueue)
	})
	t.Run("of the first skill in the order of the graph", func(t *testing.T) {
		w.last("gaps-34-d3-1", false, wrong(gapsSkill), wrong(enumSkill))
		w.want("enum-34-d1-1", ReasonConsolidation)
	})
	t.Run("out of reach", func(t *testing.T) {
		w.probabilities([]px{{"gaps-*", 0.5}})
		w.last("gaps-34-d3-1", false, wrong(gapsSkill))
		w.want("enum-34-d1-1", ReasonQueue)
	})
}

// Step 5: counting-gaps on pause comes back as a probe on the turn of the probe, the easiest
// task; consolidation takes that turn first.
func TestOlympiadProbe(t *testing.T) {
	w := olympiadWorld(t)
	w.probabilities([]px{{"gaps-12-d2-1", 0.66}, {"gaps-12-d3-1", 0.5}, {"gaps-12-d5-1", 0.3}, {"gaps-34-*", 0.2}})
	w.skill(learnermodel.SkillReadout{Skill: gapsSkill, EvidenceCount: 9})
	w.issue(enumSkill, exampleTime-100)
	w.issue(gapsSkill, exampleTime-10)
	w.trial(true, false)
	for _, tc := range []struct {
		n    int64
		task string
		code ReasonCode
	}{
		{5, "enum-34-d1-1", ReasonQueue},
		{10, "enum-34-d1-1", ReasonQueue},
		{11, "gaps-12-d2-1", ReasonProbe},
		{12, "enum-34-d1-1", ReasonQueue},
		{17, "gaps-12-d2-1", ReasonProbe},
	} {
		w.choices(tc.n)
		w.want(tc.task, tc.code)
	}
	w.choices(11)
	w.trial(false, false)
	w.want("enum-34-d1-1", ReasonQueue) // in the series
	w.trial(true, false)
	w.last("enum-34-d3-1", false, wrong(enumSkill))
	w.want("enum-34-d1-1", ReasonConsolidation)

	// A skill on pause with no candidates — every task of combinatorics-enumeration on the pause
	// before a seen solution — gives no probe, though it was handed out longer ago.
	w.in.HasLast = false
	w.skill(learnermodel.SkillReadout{Skill: enumSkill, EvidenceCount: 4})
	for _, c := range w.in.Modules[0].Pool {
		if strings.Contains(c.Task.ID, "/enum-") {
			w.in.Summary.Tasks = append(w.in.Summary.Tasks, TaskSolution{Task: c.Task.ID, SolutionSentAt: exampleTime - 1})
		}
	}
	w.want("gaps-12-d2-1", ReasonProbe)
	w.in.Summary.Tasks = nil

	// On the turn of the probe with no skill on pause, the branches below choose.
	w.in.HasLast = false
	w.probabilities([]px{{"gaps-12-d2-1", 0.8}})
	w.want("enum-34-d1-1", ReasonQueue)
}

// Step 6: "mastered" of combinatorics-enumeration opens pigeonhole-basic, a new skill; once its
// recommended task moves above grades-3-4, the mastered skill comes back in rotation.
func TestOlympiadMasteredSkills(t *testing.T) {
	w := olympiadWorld(t)
	w.skill(learnermodel.SkillReadout{Skill: enumSkill, EvidenceCount: 12, Mastered: true, MasteryBand: learnermodel.BandRef{Place: 2}})
	w.issue(enumSkill, exampleTime-100)
	w.issue(gapsSkill, exampleTime-50)
	w.want("pigeon-12-d1-1", ReasonNewSkill)

	w.issue(pigeonSkill, exampleTime-10)
	w.want("gaps-34-d1-1", ReasonQueue) // enum is mastered for the queue
	w.probabilities([]px{{"enum-34-*", 0.93}, {"enum-56-d1-1", 0.79}})
	w.want("enum-56-d1-1", ReasonRotation)
}

func grade3World(t *testing.T) *world {
	w := newWorld(t, grade3(t))
	w.probabilities([]px{{"times-tables-d1", 0.9}, {"times-tables-d3", 0.78}, {"times-tables-d5", 0.6}})
	return w
}

var timesTables = sid(grade3Math, "times-tables")

// School maths for grade 3, step 1: times-tables is the new skill; when even its easiest task is
// too hard, no skill is within reach and the bottom of the ladder gives it, the easiest task.
func TestGrade3Branches(t *testing.T) {
	w := grade3World(t)
	w.want("times-tables-d3", ReasonNewSkill)
	w.probabilities([]px{{"times-tables-d1", 0.6}, {"times-tables-d3", 0.4}, {"times-tables-d5", 0.2}})
	w.want("times-tables-d1", ReasonNewSkill)
	w.issue(timesTables, exampleTime-10)
	w.want("times-tables-d1", ReasonRotation)
}

// School maths for grade 3, step 4: all four mastered for the queue — branch 8, the one not
// handed out for longest.
func TestGrade3AllMastered(t *testing.T) {
	w := grade3World(t)
	w.probabilities([]px{{"*", 0.78}})
	for k, s := range []string{"comparing-fractions", "division-with-remainder", "fractions-of-a-whole", "times-tables"} {
		w.skill(learnermodel.SkillReadout{Skill: sid(grade3Math, s), EvidenceCount: 10, Mastered: true, MasteryBand: learnermodel.BandRef{Place: 1}})
		w.issue(sid(grade3Math, s), exampleTime-int64(100*(k+1)))
	}
	w.want("times-tables-d1", ReasonRotation)
	w.issue(timesTables, exampleTime)
	w.want("fractions-of-a-whole-d1", ReasonRotation)
}

// The bottom of the ladder is the lowest band taught among the open skills with candidates, not
// the first skill in the order of the graph.
func TestBottomOfTheLadder(t *testing.T) {
	module := "example/ladder"
	bands := []learnermodel.Band{{Place: 1, ID: "low"}, {Place: 2, ID: "high"}}
	low, high := learnermodel.BandRef{Place: 1, ID: "low"}, learnermodel.BandRef{Place: 2, ID: "high"}
	a, b := sid(module, "a"), sid(module, "b")
	m := exampleModule(t, module, bands, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
	}, []learnermodel.Task{
		exampleTask(module, "a-high", high, 1, a),
		exampleTask(module, "b-low", low, 3, b),
		exampleTask(module, "b-high", high, 1, b),
	})
	w := newWorld(t, m)
	w.probabilities([]px{{"a-high", 0.6}, {"b-low", 0.3}, {"b-high", 0.5}})
	w.want("b-high", ReasonNewSkill)
	w.issue(b, exampleTime-10)
	w.want("b-high", ReasonRotation)
	// With both within reach and mastered for the queue, branch 8 takes the one handed out longest
	// ago, whatever its band.
	w.probabilities([]px{{"a-high", 0.8}, {"b-high", 0.8}})
	for _, s := range []string{a, b} {
		w.skill(learnermodel.SkillReadout{Skill: s, EvidenceCount: 5, Mastered: true, MasteryBand: high})
	}
	w.issue(a, exampleTime-100)
	w.want("a-high", ReasonRotation)
	w.probabilities([]px{{"a-high", 0.6}, {"b-high", 0.5}})
	w.in.Summary.Skills = []SkillIssue{{Skill: b, Issued: true, LastIssued: exampleTime - 10}}

	// Without the task at the low band, both skills are at the bottom, and a is not handed out yet.
	w.in.Modules[0].Graph.Skills[1].BandsTaught = []int32{2}
	w.want("a-high", ReasonNewSkill)
	// At the bottom, branch 7 takes a skill not mastered for the queue first, though a mastered one
	// was handed out longer ago; with all mastered, the one handed out longest ago.
	w.issue(a, exampleTime-100)
	w.want("a-high", ReasonRotation)
	w.skill(learnermodel.SkillReadout{Skill: b, EvidenceCount: 5})
	w.want("b-high", ReasonRotation)
}

// The second pass takes the tasks on pause back, after the others, but not the task of the last
// closed attempt; with nothing even then, the rule gives nothing.
func TestSecondPassChoice(t *testing.T) {
	w := grade3World(t)
	w.in.Summary.Tasks = []TaskSolution{
		{Task: grade3Math + "/task/times-tables-d1", SolutionSentAt: exampleTime - 300},
		{Task: grade3Math + "/task/times-tables-d3", SolutionSentAt: exampleTime - 200},
		{Task: grade3Math + "/task/times-tables-d5", SolutionSentAt: exampleTime - 100},
	}
	w.last("times-tables-d5", false, wrong(timesTables))
	w.want("", ReasonUnspecified)
	got, ok := w.choose(true)
	if want := (Choice{Kind: ChoiceTask, ID: grade3Math + "/task/times-tables-d1", Reason: ReasonNewSkill}); !ok || got != want {
		t.Errorf("second pass: %+v, want %+v", got, want)
	}

	w.in.Modules[0].Pool = w.in.Modules[0].Pool[2:3] // times-tables-d5 alone
	w.p = w.p[2:3]
	if got, ok := w.choose(true); ok {
		t.Errorf("second pass with the last task alone: %+v, want nothing", got)
	}
	if !strings.HasSuffix(w.in.Modules[0].Pool[0].Task.ID, "times-tables-d5") {
		t.Fatalf("the pool is %s", w.in.Modules[0].Pool[0].Task.ID)
	}
}

// longestAgo puts a skill never handed out first, then the one handed out longest ago, then the
// order of the skills.
func TestLongestAgo(t *testing.T) {
	all := func(skillView) bool { return true }
	for _, tc := range []struct {
		skills []skillView
		want   string
	}{
		{[]skillView{{id: "a", issued: true, lastIssued: 5}, {id: "b"}, {id: "c"}}, "b"},
		{[]skillView{{id: "b"}, {id: "a", issued: true, lastIssued: 1}}, "b"},
		{[]skillView{{id: "a", issued: true, lastIssued: 5}, {id: "b", issued: true, lastIssued: 3}, {id: "c", issued: true, lastIssued: 3}}, "b"},
	} {
		if got, ok := longestAgo(tc.skills, all); !ok || got.id != tc.want {
			t.Errorf("longestAgo = %s, want %s", got.id, tc.want)
		}
	}
	if _, ok := longestAgo(nil, all); ok {
		t.Error("longestAgo of no skills found one")
	}
}
