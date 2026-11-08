package engine

import (
	"slices"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The repetition, on the example of docs/architecture.md, section 3.4 — school maths for grade 3 —
// and at its edges. The reading of the learner model, the review dates and the P of the candidates
// are set by hand, as in modules_internal_test.go.

var (
	divisionSkill  = sid(grade3Math, "division-with-remainder")
	comparingSkill = sid(grade3Math, "comparing-fractions")
)

// repWorld is school maths for grade 3 past its trial series, with no lessons, every skill handed
// out — times-tables first, then fractions-of-a-whole, division-with-remainder and
// comparing-fractions —, and times-tables and fractions-of-a-whole mastered with the review dates
// given.
func repWorld(t *testing.T, times, fractions learnermodel.ReviewDate) *many {
	t.Helper()
	m := gradeLessons(t)
	m.Lessons = nil
	w := newMany(t, m)
	w.readout(grade3Math).TrialOver = true
	for k, s := range []string{timesTablesSkill, fractionsSkill, divisionSkill, comparingSkill} {
		w.issue(s, int64(100*(k+1)))
	}
	w.masteredBy(timesTablesSkill, times)
	w.masteredBy(fractionsSkill, fractions)
	return w
}

func (w *many) masteredBy(skill string, review learnermodel.ReviewDate) {
	w.skill(learnermodel.SkillReadout{Skill: skill, EvidenceCount: 10, Mastered: true, MasteryBand: learnermodel.BandRef{Place: 1, ID: "grade-3"}, Review: review})
}

// choices sets n, the number of choices of a module before this one.
func (w *many) choices(module string, n int64) {
	w.in.Summary.Modules = slices.DeleteFunc(w.in.Summary.Modules, func(m ModuleChoices) bool { return m.Module == module })
	w.in.Summary.Modules = append(w.in.Summary.Modules, ModuleChoices{Module: module, Choices: n})
}

func due(at int64) learnermodel.ReviewDate { return learnermodel.ReviewDate{At: at} }

// Steps 1–5 of the example. The review date of times-tables has come: on n = 13, the turn of the
// repetition, branch 5 gives it, a task in the corridor; on n = 14 the queue goes on, or
// consolidation after a failed repetition. With both due, the earlier date first, and the other on
// the next turn, n = 16. With every skill mastered, the repetitions come one after another off the
// turn, by date, before branch 8; with no date come, branch 8.
func TestRepetitionGrade3(t *testing.T) {
	w := repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	permuted(t, w)
	w.choices(grade3Math, 14)
	w.wantChoice(taskChoice(tid(grade3Math, "division-with-remainder-d1"), ReasonQueue))
	w.in.Last, w.in.HasLast = Last{Task: tid(grade3Math, "times-tables-d3"), Results: []SkillResult{{Skill: timesTablesSkill, Answered: 1}}}, true
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonConsolidation))
	w.in.HasLast = false

	w.masteredBy(fractionsSkill, due(exampleTime-20))
	w.choices(grade3Math, 13)
	w.wantChoice(taskChoice(tid(grade3Math, "fractions-of-a-whole-d1"), ReasonRepetition)) // the earlier date
	w.masteredBy(fractionsSkill, due(exampleTime+86_400))                                  // reviewed: a date further on
	w.choices(grade3Math, 16)
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	for _, n := range []int64{14, 15} {
		w.choices(grade3Math, n)
		w.wantChoice(taskChoice(tid(grade3Math, "division-with-remainder-d1"), ReasonQueue)) // the new goes on
	}

	w.masteredBy(divisionSkill, due(exampleTime-5))
	w.masteredBy(comparingSkill, due(exampleTime-30))
	w.choices(grade3Math, 14)
	w.wantChoice(taskChoice(tid(grade3Math, "comparing-fractions-d1"), ReasonRepetition)) // off the turn, by date
	for _, s := range []string{timesTablesSkill, divisionSkill, comparingSkill} {
		w.masteredBy(s, due(exampleTime+86_400))
	}
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRotation)) // branch 8, handed out longest ago
}

// The turn of the repetition is n > 5 with n mod 3 = 1: not n = 4, nor the turn of the probe,
// n = 11, which never meets it. It needs no end of the trial series (the amendment of 2026-11-08 to
// R06). A date marked already comes before the earliest date; equal dates go by the order of the
// graph; the date is come at the time of the choice, not a second before it.
func TestRepetitionTurn(t *testing.T) {
	w := repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	for _, tc := range []struct {
		n          int64
		repetition bool
	}{{4, false}, {5, false}, {6, false}, {7, true}, {10, true}, {11, false}, {16, true}} {
		w.choices(grade3Math, tc.n)
		want := taskChoice(tid(grade3Math, "division-with-remainder-d1"), ReasonQueue)
		if tc.repetition {
			want = taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition)
		}
		if got := ready(w.in, w.p).choose(); got != want {
			t.Errorf("n = %d: chose %+v, want %+v", tc.n, got, want)
		}
	}
	w.choices(grade3Math, 13)
	w.readout(grade3Math).TrialOver = false
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))

	w.masteredBy(fractionsSkill, learnermodel.ReviewDate{At: exampleTime + 1000, Already: true})
	w.wantChoice(taskChoice(tid(grade3Math, "fractions-of-a-whole-d1"), ReasonRepetition))
	w.masteredBy(fractionsSkill, due(exampleTime-10))
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition)) // equal dates: tier 1 first
	w.masteredBy(timesTablesSkill, due(exampleTime))
	w.masteredBy(fractionsSkill, due(exampleTime+1))
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	w.masteredBy(timesTablesSkill, due(exampleTime+1))
	w.wantChoice(taskChoice(tid(grade3Math, "division-with-remainder-d1"), ReasonQueue))
}

// Only a skill mastered for the queue with candidates is due: a date the reading gives a skill not
// mastered counts nothing, a skill whose tasks are all on pause waits, and a skill whose
// recommended task went above its mastery band goes into the queue by branch 7.
func TestRepetitionDue(t *testing.T) {
	w := repWorld(t, due(exampleTime+1000), due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.skill(learnermodel.SkillReadout{Skill: divisionSkill, EvidenceCount: 3, Review: due(exampleTime - 10)})
	w.wantChoice(taskChoice(tid(grade3Math, "division-with-remainder-d1"), ReasonQueue)) // its date counts nothing
	w.issue(comparingSkill, 50)
	w.wantChoice(taskChoice(tid(grade3Math, "comparing-fractions-d1"), ReasonQueue))

	w.masteredBy(timesTablesSkill, due(exampleTime-10))
	w.solutions(TaskSolution{tid(grade3Math, "times-tables-d1"), exampleTime - 10}, TaskSolution{tid(grade3Math, "times-tables-d3"), exampleTime - 10},
		TaskSolution{tid(grade3Math, "times-tables-d5"), exampleTime - 10})
	w.wantChoice(taskChoice(tid(grade3Math, "comparing-fractions-d1"), ReasonQueue))

	o := olympiad(t)
	enum := sid(olympiadMath, "combinatorics-enumeration")
	w = newMany(t, o)
	w.readout(olympiadMath).TrialOver = true
	w.probabilities([]px{{"enum-12-*", 0.99}, {"enum-34-d1-1", 0.78}, {"enum-34-d3-1", 0.6}, {"enum-34-d5-1", 0.4}, {"enum-56-*", 0.3}})
	for _, s := range []string{enum, sid(olympiadMath, "counting-gaps"), sid(olympiadMath, "pigeonhole-basic")} {
		w.issue(s, 100)
	}
	w.skill(learnermodel.SkillReadout{Skill: enum, EvidenceCount: 10, Mastered: true, MasteryBand: learnermodel.BandRef{Place: 1, ID: "grades-1-2"}, Review: due(exampleTime - 10)})
	w.choices(olympiadMath, 13)
	w.wantChoice(taskChoice(tid(olympiadMath, "enum-34-d1-1"), ReasonRotation)) // above its mastery band: branch 7
}

// A skill due for review out of reach gets the task nearest the aim, the easiest, and a right
// answer brings it back within reach. The repetition gives a task, never a lesson, even with a
// lesson of its band not passed. On the map the code is the map's.
func TestRepetitionTask(t *testing.T) {
	w := repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.probabilities([]px{{"times-tables-d1", 0.5}, {"times-tables-d3", 0.4}, {"times-tables-d5", 0.3}})
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	view := func() skillView {
		vs := ready(w.in, w.p).skills(0, false)
		return vs[slices.IndexFunc(vs, func(v skillView) bool { return v.id == timesTablesSkill })]
	}
	if view().reachable {
		t.Fatalf("%s within reach with every P below 0.70", timesTablesSkill)
	}
	is := &w.in.Summary.Skills[slices.IndexFunc(w.in.Summary.Skills, func(s SkillIssue) bool { return s.Skill == timesTablesSkill })]
	is.LastCode, is.LastResult, is.HasLast = ReasonRepetition, SkillResult{Answered: 1, FirstRight: 1}, true
	if !view().reachable {
		t.Errorf("%s out of reach after a right repetition, want it back", timesTablesSkill)
	}

	w = repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.in.Modules[0].Lessons = []Lesson{lessonOf(w.in.Modules[0], "times-tables", "times-tables-", timesTablesSkill)}
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	w.in.MapSkill = fractionsSkill
	w.in.Modules[0].Lessons = nil
	w.wantChoice(taskChoice(tid(grade3Math, "fractions-of-a-whole-d1"), ReasonSkillChosen))
}

// Each module has its turn of the repetition by its own n.
func TestRepetitionModules(t *testing.T) {
	w := repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	other := otherModule(t)
	w2 := newMany(t, other, w.in.Modules[0])
	w2.in.Readout.Modules[1] = w.in.Readout.Modules[0]
	w2.in.Summary = w.in.Summary
	w2.choices(grade3Math, 13)
	w2.choices("example/other", 14)
	w2.issue(sid("example/other", "x"), 50) // the other module goes first
	w2.wantChoice(taskChoice(tid("example/other", "x-1"), ReasonQueue))
	w2.issue(sid("example/other", "x"), 500) // now school maths, on its turn of the repetition
	w2.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	permuted(t, w2)
}

// A skill due for review within reach gets the task in the corridor by its P with retention, not
// the easiest; a skill not handed out that a lesson of the author's order not passed teaches waits
// for that lesson, even with its date come.
func TestRepetitionCorridorAndAwaits(t *testing.T) {
	w := repWorld(t, due(exampleTime-10), due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.probabilities([]px{{"times-tables-d1", 0.95}, {"times-tables-d3", 0.78}, {"times-tables-d5", 0.5}})
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d3"), ReasonRepetition))

	module := "example/review"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	v, r := sid(module, "v"), sid(module, "r")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: v, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: r, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "v-1", band, 1, v)})
	m.Lessons = []Lesson{{ID: lid(module, "vr"), Skills: []string{v, r}}}
	w = newMany(t, m)
	w.skill(learnermodel.SkillReadout{Skill: v, EvidenceCount: 10, Mastered: true, MasteryBand: band, Review: due(exampleTime - 10)})
	w.choices(module, 13)
	w.wantChoice(lessonChoice(lid(module, "vr"), ReasonAuthorOrder))
}

// Of two dates marked already, the skill practised longest ago comes first, whatever the order of
// the graph (the amendment of 2026-11-08 to R06): practice does not clear the mark, and the first in
// the graph would take every repetition. A skill closed, whatever its reading says, is not due.
func TestRepetitionAlready(t *testing.T) {
	already := learnermodel.ReviewDate{Already: true}
	w := repWorld(t, already, already)
	w.choices(grade3Math, 13)
	set := func(skill string, last int64) {
		w.skill(learnermodel.SkillReadout{Skill: skill, EvidenceCount: 10, LastPractice: last, Mastered: true,
			MasteryBand: learnermodel.BandRef{Place: 1, ID: "grade-3"}, Review: already})
	}
	set(timesTablesSkill, exampleTime-100)
	set(fractionsSkill, exampleTime-200)
	w.wantChoice(taskChoice(tid(grade3Math, "fractions-of-a-whole-d1"), ReasonRepetition))
	set(fractionsSkill, exampleTime-50)
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition))
	set(fractionsSkill, exampleTime-100)
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonRepetition)) // as long ago: the graph

	// comparing-fractions, mastered and due by its reading, with no evidence of its own and its
	// prerequisite no longer mastered, is closed.
	w = repWorld(t, notDue, due(exampleTime+1000))
	w.choices(grade3Math, 13)
	w.skill(learnermodel.SkillReadout{Skill: fractionsSkill, EvidenceCount: 10})
	w.skill(learnermodel.SkillReadout{Skill: comparingSkill, Mastered: true, MasteryBand: learnermodel.BandRef{Place: 1, ID: "grade-3"}, Review: due(exampleTime - 10)})
	if w.open(comparingSkill) {
		t.Fatalf("%s open, want it closed", comparingSkill)
	}
	w.wantChoice(taskChoice(tid(grade3Math, "fractions-of-a-whole-d1"), ReasonQueue)) // the queue, not a repetition of comparing-fractions
}
