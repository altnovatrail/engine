package engine

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Lessons, the author's order and nothing to give, on the examples of docs/architecture.md,
// section 3.3 — the three modules of the first slice and the AI closed — and at their edges. The
// reading of the learner model and the P of the candidates are set by hand, as in
// modules_internal_test.go.

func lid(module, name string) string { return module + "/lesson/" + name }

const englishEssays = "novatrail/english-essays"

// wantChoice checks the choice, a lesson or a task or nothing.
func (w *many) wantChoice(want Choice) {
	w.t.Helper()
	if got := ready(w.in, w.p).choose(); got != want {
		w.t.Errorf("chose %+v, want %+v", got, want)
	}
}

func lessonChoice(id string, code ReasonCode) Choice {
	return Choice{Kind: ChoiceLesson, ID: id, Reason: code}
}

func taskChoice(id string, code ReasonCode) Choice {
	return Choice{Kind: ChoiceTask, ID: id, Reason: code}
}

// lessonOf is a lesson of a module on skills, with the tasks of its pool whose name begins with
// prefix, in the order of the pool.
func lessonOf(m Module, name, prefix string, skills ...string) Lesson {
	l := Lesson{ID: lid(m.Graph.Module, name), Skills: skills}
	for _, c := range m.Pool {
		if prefix != "" && strings.HasPrefix(c.Task.ID, tid(m.Graph.Module, prefix)) {
			l.Tasks = append(l.Tasks, c.Task)
		}
	}
	return l
}

// passRun ends a lesson as passed at a time, and hands out its skills then, as api records it.
func (w *many) passRun(l Lesson, at int64) {
	w.in.Runs = slices.DeleteFunc(w.in.Runs, func(r Run) bool { return r.Lesson == l.ID && r.End == RunOpen })
	w.in.Runs = append(w.in.Runs, Run{Lesson: l.ID, Skills: l.Skills, End: RunPassed, EndedAt: at})
	for _, s := range l.Skills {
		w.issue(s, at)
	}
}

// solutions sends the solutions of tasks to the device at times, as a lesson's tasks are.
func (w *many) solutions(sent ...TaskSolution) {
	for _, s := range sent {
		w.in.Summary.Tasks = slices.DeleteFunc(w.in.Summary.Tasks, func(t TaskSolution) bool { return t.Task == s.Task })
		w.in.Summary.Tasks = append(w.in.Summary.Tasks, s)
	}
}

// gradeLessons builds school maths for grade 3 of section 3.1 with a lesson on each skill — an
// explanation and its three pages, of difficulty 1, 3 and 5.
func gradeLessons(t *testing.T) Module {
	t.Helper()
	band := learnermodel.BandRef{Place: 1, ID: "grade-3"}
	times, division, fractions, comparing := sid(grade3Math, "times-tables"), sid(grade3Math, "division-with-remainder"),
		sid(grade3Math, "fractions-of-a-whole"), sid(grade3Math, "comparing-fractions")
	var tasks []learnermodel.Task
	for _, s := range []string{times, division, fractions, comparing} {
		for _, d := range []int32{1, 3, 5} {
			tasks = append(tasks, exampleTask(grade3Math, fmt.Sprintf("%s-d%d", s[strings.LastIndex(s, "/")+1:], d), band, d, s))
		}
	}
	m := exampleModule(t, grade3Math, []learnermodel.Band{{Place: 1, ID: "grade-3", Grades: []int32{3}}}, []learnermodel.Skill{
		{ID: times, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: division, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{times}, Tier: 2},
		{ID: fractions, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{times}, Tier: 2},
		{ID: comparing, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{fractions}, Tier: 3},
	}, tasks)
	for _, s := range []string{times, division, fractions, comparing} {
		name := s[strings.LastIndex(s, "/")+1:]
		m.Lessons = append(m.Lessons, lessonOf(m, name, name+"-", s))
	}
	return m
}

// otherModule is a module of one skill and one task, for a module that is not the example's.
func otherModule(t *testing.T) Module {
	t.Helper()
	module := "example/other"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	return exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: sid(module, "x"), Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "x-1", band, 1, sid(module, "x"))})
}

// Step 1 of the example: school maths for grade 3. A new times-tables gets its lesson. After the
// lesson its three pages are on pause, and consolidation of a mistake in the last finds no
// candidate: with another module the turn goes there; with this module alone the second pass gives
// consolidation by the page whose solution went longest ago, not the last.
func TestLessonsGrade3(t *testing.T) {
	m := gradeLessons(t)
	times := sid(grade3Math, "times-tables")
	w := newMany(t, m)
	w.wantChoice(lessonChoice(lid(grade3Math, "times-tables"), ReasonNewSkill))

	w.readout(grade3Math).TrialOver = true
	w.skill(learnermodel.SkillReadout{Skill: times, EvidenceCount: 6})
	w.passRun(m.Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(grade3Math, "times-tables") })], exampleTime-50)
	w.solutions(
		TaskSolution{tid(grade3Math, "times-tables-d1"), exampleTime - 300}, TaskSolution{tid(grade3Math, "times-tables-d3"), exampleTime - 200},
		TaskSolution{tid(grade3Math, "times-tables-d5"), exampleTime - 100},
	)
	w.in.Last = Last{Task: tid(grade3Math, "times-tables-d5"), InLesson: true, Lesson: lid(grade3Math, "times-tables"),
		Results: []SkillResult{{Skill: times, Answered: 1}}}
	w.in.HasLast = true
	if c, ok := ready(w.in, w.p).pass(false); ok {
		t.Errorf("first pass after the lesson: %+v, want nothing", c)
	}
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonConsolidation))
	// Consolidation gives no lesson, even with one on the skill not passed.
	w.in.Runs = append(w.in.Runs, Run{Lesson: lid(grade3Math, "times-tables-again"), Skills: []string{times}, End: RunAbandoned, EndedAt: 1})
	w.in.Modules[0].Lessons = append(slices.Clone(w.in.Modules[0].Lessons), Lesson{ID: lid(grade3Math, "times-tables-again"), Skills: []string{times}})
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonConsolidation))
	w.in.Modules[0].Lessons = m.Lessons

	other := otherModule(t)
	w2 := newMany(t, other, m)
	w2.in.Readout.Modules[1] = w.in.Readout.Modules[0]
	w2.in.Runs, w2.in.Summary, w2.in.Last, w2.in.HasLast = w.in.Runs, w.in.Summary, w.in.Last, true
	w2.wantChoice(taskChoice(tid("example/other", "x-1"), ReasonNewSkill))
}

// olympiadLessons builds olympiad maths of section 3.1 with a lesson on each skill and band, of the
// three tasks of the skill at the band.
func olympiadLessons(t *testing.T) Module {
	t.Helper()
	m := olympiad(t)
	for _, s := range []struct{ code, skill string }{{"enum", "combinatorics-enumeration"}, {"gaps", "counting-gaps"}, {"pigeon", "pigeonhole-basic"}} {
		for _, b := range []struct{ code, band string }{{"12", "grades-1-2"}, {"34", "grades-3-4"}, {"56", "grades-5-6"}} {
			m.Lessons = append(m.Lessons, lessonOf(m, s.skill+"-"+b.band, s.code+"-"+b.code+"-", sid(olympiadMath, s.skill)))
		}
	}
	return m
}

// Step 2: olympiad maths for a learner of grade 3. A new combinatorics-enumeration gets the lesson
// of the band of its recommended task, grades-3-4; once that task goes up to grades-5-6, branch 7
// gives the lesson of grades-5-6 — of the queue, or round again for a skill mastered at grades-3-4.
// The lesson of grades-1-2 never suits such a learner.
func TestLessonsOlympiad(t *testing.T) {
	m := olympiadLessons(t)
	enum, gaps := sid(olympiadMath, "combinatorics-enumeration"), sid(olympiadMath, "counting-gaps")
	w := newMany(t, m)
	w.probabilities([]px{{"enum-12-*", 0.95}, {"enum-34-d1-1", 0.78}, {"enum-34-d3-1", 0.6}, {"enum-34-d5-1", 0.4}, {"enum-56-*", 0.3}})
	w.wantChoice(lessonChoice(lid(olympiadMath, "combinatorics-enumeration-grades-3-4"), ReasonNewSkill))

	lesson := func(name string) Lesson {
		return m.Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(olympiadMath, name) })]
	}
	w.passRun(lesson("combinatorics-enumeration-grades-3-4"), 100)
	w.issue(gaps, 200)
	w.wantChoice(taskChoice(tid(olympiadMath, "enum-34-d1-1"), ReasonQueue)) // the lesson of grades-1-2 does not suit

	w.probabilities([]px{{"enum-12-*", 0.99}, {"enum-34-*", 0.97}, {"enum-56-d1-1", 0.78}})
	w.wantChoice(lessonChoice(lid(olympiadMath, "combinatorics-enumeration-grades-5-6"), ReasonQueue))
	w.skill(learnermodel.SkillReadout{Skill: enum, EvidenceCount: 10, Mastered: true, Review: notDue, MasteryBand: learnermodel.BandRef{Place: 2, ID: "grades-3-4"}})
	w.wantChoice(lessonChoice(lid(olympiadMath, "pigeonhole-basic-grades-1-2"), ReasonNewSkill)) // mastered, enum opens pigeonhole-basic
	w.issue(sid(olympiadMath, "pigeonhole-basic"), 300)
	w.wantChoice(lessonChoice(lid(olympiadMath, "combinatorics-enumeration-grades-5-6"), ReasonRotation))
	permuted(t, w)
}

// essays builds the module of short essays of the first slice (altnovatrail/modules, main
// 6d061e9): reading model essays, none, and choosing a first sentence, with a verdict, in tier 1;
// writing one, paragraphs, endings and improving an essay, with a rubric, after them; six lessons,
// and every task in a lesson. adult says whether the learner is an adult: a child gets neither the
// lesson writing-about-a-day nor its task.
func essays(t *testing.T, adult bool) Module {
	t.Helper()
	band := learnermodel.BandRef{Place: 1, ID: "first-steps"}
	s := func(name string) string { return sid(englishEssays, name) }
	tasks := []learnermodel.Task{
		exampleTask(englishEssays, "best-first-sentence-1", band, 1, s("choosing-a-first-sentence")),
		exampleTask(englishEssays, "best-first-sentence-2", band, 2, s("choosing-a-first-sentence")),
		exampleTask(englishEssays, "best-first-sentence-3", band, 3, s("choosing-a-first-sentence")),
		exampleTask(englishEssays, "talking-umbrella-1", band, 2, s("writing-a-first-sentence")),
		exampleTask(englishEssays, "weather-shop-1", band, 2, s("paragraphs")),
		exampleTask(englishEssays, "undiscovered-island-1", band, 3, s("endings")),
		exampleTask(englishEssays, "dragon-draft-1", band, 3, s("improving-an-essay")),
	}
	if adult {
		tasks = append(tasks, exampleTask(englishEssays, "a-day-you-remember-1", band, 3, s("paragraphs"), s("endings")))
	}
	m := exampleModule(t, englishEssays, []learnermodel.Band{{Place: 1, ID: "first-steps"}}, []learnermodel.Skill{
		{ID: s("reading-model-essays"), Assessability: learnermodel.AssessabilityNone, Tier: 1},
		{ID: s("choosing-a-first-sentence"), Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: s("writing-a-first-sentence"), Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{s("choosing-a-first-sentence")}, Tier: 2},
		{ID: s("paragraphs"), Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{s("writing-a-first-sentence")}, Tier: 3},
		{ID: s("endings"), Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{s("paragraphs")}, Tier: 4},
		{ID: s("improving-an-essay"), Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{s("paragraphs")}, Tier: 4},
	}, tasks)
	m.Lessons = []Lesson{
		lessonOf(m, "first-sentences", "", s("choosing-a-first-sentence"), s("writing-a-first-sentence")),
		lessonOf(m, "reading-model-essays", "", s("reading-model-essays")),
		lessonOf(m, "paragraphs", "weather-shop", s("paragraphs")),
		lessonOf(m, "endings", "undiscovered-island", s("endings")),
		lessonOf(m, "improving-an-essay", "dragon-draft", s("improving-an-essay")),
	}
	m.Lessons[0].Tasks = slices.Concat(lessonOf(m, "", "best-first-sentence", "").Tasks, lessonOf(m, "", "talking-umbrella", "").Tasks)
	if adult {
		m.Lessons = append(m.Lessons, lessonOf(m, "writing-about-a-day", "a-day-you-remember", s("endings")))
	}
	return m
}

// Step 3: short essays, by the author's order, all by branch 6a. first-sentences and
// reading-model-essays the graph leaves unordered, both of tier 1: first-sentences first, by ID;
// then reading-model-essays, of tier 1, before paragraphs, of tier 3 (the amendment of 2026-11-08
// to R05); of the three lessons after paragraphs, by ID. choosing-a-first-sentence, with a
// verdict, comes by first-sentences, not by branch 6, and once handed out goes into the queue, but
// branch 6a comes before branch 7. A child gets no writing-about-a-day.
func TestLessonsEssays(t *testing.T) {
	for _, adult := range []bool{true, false} {
		m := essays(t, adult)
		w := newMany(t, m)
		order := []string{"first-sentences", "reading-model-essays", "paragraphs", "endings", "improving-an-essay"}
		if adult {
			order = append(order, "writing-about-a-day")
		}
		for k, name := range order {
			w.wantChoice(lessonChoice(lid(englishEssays, name), ReasonAuthorOrder))
			if t.Failed() {
				t.Fatalf("adult %v: at lesson %d", adult, k)
			}
			permuted(t, w)
			w.passRun(m.Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(englishEssays, name) })], int64(100*(k+1)))
		}
		w.wantChoice(taskChoice(tid(englishEssays, "best-first-sentence-1"), ReasonQueue))
	}
}

// Step 4: the AI closed. The pool holds only tasks answered by a touch, and the lessons only such
// tasks; with none in any module, the second pass changes nothing, and the choice is nothing to
// give.
func TestLessonsAIClosed(t *testing.T) {
	g, o := gradeLessons(t), olympiadLessons(t)
	for _, m := range []*Module{&g, &o} {
		m.Pool, m.Lessons = nil, nil
	}
	w := newMany(t, g, o)
	w.wantChoice(Choice{Kind: ChoiceNothing})
}

// Branch 1: an open run of a module of the engine, in the order of the lessons read by its skills —
// one that comes after another last, then by tier, then by ID —, of the modules' by full ID; not
// the run the learner has just left, in either mode; on the map only a run that teaches the skill
// chosen there; a run of no module of the engine is not read. A lesson with an open run no other
// branch gives.
func TestOpenRun(t *testing.T) {
	m := gradeLessons(t)
	times, fractions, comparing := sid(grade3Math, "times-tables"), sid(grade3Math, "fractions-of-a-whole"), sid(grade3Math, "comparing-fractions")
	w := newMany(t, m, otherModule(t))
	open := func(lesson string, skills ...string) Run { return Run{Lesson: lid(grade3Math, lesson), Skills: skills} }
	w.in.Runs = []Run{open("a-comparing", comparing), open("z-fractions", fractions), {Lesson: "example/gone/lesson/a"}}
	w.wantChoice(lessonChoice(lid(grade3Math, "z-fractions"), ReasonOpenRun)) // a-comparing comes after it
	w.in.Runs = append(w.in.Runs, open("a-div-comp", sid(grade3Math, "division-with-remainder"), comparing))
	w.wantChoice(lessonChoice(lid(grade3Math, "z-fractions"), ReasonOpenRun)) // a-div-comp, of tier 2 too and less by ID, comes after it
	w.in.Runs = append(w.in.Runs, open("y-times", times))
	w.wantChoice(lessonChoice(lid(grade3Math, "y-times"), ReasonOpenRun)) // tier 1, before z-fractions of tier 2
	w.in.Runs = append(w.in.Runs, Run{Lesson: lid("example/other", "x"), Skills: []string{sid("example/other", "x")}})
	w.wantChoice(lessonChoice(lid("example/other", "x"), ReasonOpenRun)) // the modules' firsts by full ID
	permuted(t, w)

	w.in.Last, w.in.HasLast = Last{Task: tid("example/other", "x-1"), Left: true, InLesson: true, Lesson: lid("example/other", "x")}, true
	w.wantChoice(lessonChoice(lid(grade3Math, "y-times"), ReasonOpenRun))
	w.in.MapSkill = sid("example/other", "x")
	w.wantChoice(taskChoice(tid("example/other", "x-1"), ReasonSkillChosen)) // left, and the others teach no x
	w.in.MapSkill = fractions
	w.wantChoice(lessonChoice(lid(grade3Math, "z-fractions"), ReasonOpenRun))

	// A lesson with an open run the learner has just left no other branch gives.
	w = newMany(t, m)
	w.in.Runs = []Run{open("times-tables", times)}
	w.in.Last, w.in.HasLast = Last{Task: tid(grade3Math, "times-tables-d1"), Left: true, InLesson: true, Lesson: lid(grade3Math, "times-tables")}, true
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonNewSkill))
}

// A lesson is passed only by a run that ended "passed": one ended unchecked, not suitable or
// abandoned is given again by a new run. A lesson with no task suits any band.
func TestLessonPassed(t *testing.T) {
	m := gradeLessons(t)
	l := m.Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(grade3Math, "times-tables") })]
	for _, end := range []RunEnd{RunUnchecked, RunNotSuitable, RunAbandoned} {
		w := newMany(t, m)
		w.in.Runs = []Run{{Lesson: l.ID, Skills: l.Skills, End: end, EndedAt: 10}}
		w.wantChoice(lessonChoice(l.ID, ReasonNewSkill))
	}
	w := newMany(t, m)
	w.passRun(l, 10)
	w.in.Summary.Skills = nil
	w.wantChoice(taskChoice(tid(grade3Math, "times-tables-d1"), ReasonNewSkill))

	for i := range m.Lessons {
		m.Lessons[i].Tasks = nil
	}
	w = newMany(t, m)
	w.wantChoice(lessonChoice(l.ID, ReasonNewSkill))
}

// authorWorld is a module where lesson m teaches r, with a rubric, and lesson l teaches v, with a
// verdict, and w, with a rubric, that needs r: l comes after m.
func authorWorld(t *testing.T) (*many, Module) {
	t.Helper()
	module := "example/author"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	r, v, wr := sid(module, "r"), sid(module, "v"), sid(module, "w")
	mod := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: r, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
		{ID: v, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: wr, Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{r}, Tier: 2},
	}, []learnermodel.Task{exampleTask(module, "r-1", band, 1, r), exampleTask(module, "v-1", band, 1, v), exampleTask(module, "w-1", band, 1, wr)})
	mod.Lessons = []Lesson{lessonOf(mod, "m", "r-", r), lessonOf(mod, "l", "v-", v, wr)}
	mod.Lessons[1].Tasks = append(mod.Lessons[1].Tasks, lessonOf(mod, "", "w-", "").Tasks...)
	return newMany(t, mod), mod
}

// A skill with a verdict not handed out, taught by a lesson of the author's order not passed,
// comes by that lesson: branches 6–9 do not take it (the amendment of 2026-11-08 to R05) — even
// when the learner has just left the run of the lesson before it, and branch 6a gives nothing.
func TestAwaitsLesson(t *testing.T) {
	w, mod := authorWorld(t)
	w.wantChoice(lessonChoice(lid(mod.Graph.Module, "m"), ReasonAuthorOrder))
	w.in.Runs = []Run{{Lesson: lid(mod.Graph.Module, "m"), Skills: []string{sid(mod.Graph.Module, "r")}}}
	w.in.Last, w.in.HasLast = Last{Task: tid(mod.Graph.Module, "r-1"), Left: true, InLesson: true, Lesson: lid(mod.Graph.Module, "m")}, true
	w.wantChoice(Choice{Kind: ChoiceNothing}) // v within reach would be branch 8's
	w.probabilities([]px{{"v-1", 0.5}})
	w.wantChoice(Choice{Kind: ChoiceNothing}) // and out of reach, branch 9's

	w.passRun(mod.Lessons[0], 100)
	w.in.HasLast = false
	w.wantChoice(lessonChoice(lid(mod.Graph.Module, "l"), ReasonAuthorOrder))
}

// A task with a rubric outside the lessons comes in branch 6a as a lesson of its skills with a
// rubric would: after the lessons of the author's order of its own skills and of those it would
// come after are passed (the amendment of 2026-11-08 to R05), after the lessons, in the order of
// its skills in the graph, then by ID, while no attempt of it is closed by "Done".
func TestRubricTasksOutsideLessons(t *testing.T) {
	module := "example/rubric"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	r1, r2 := sid(module, "r1"), sid(module, "r2")
	mod := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: r1, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
		{ID: r2, Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{r1}, Tier: 2},
	}, []learnermodel.Task{exampleTask(module, "a-r2", band, 1, r2), exampleTask(module, "z-r1", band, 1, r1), exampleTask(module, "r1-in", band, 1, r1)})
	mod.Lessons = []Lesson{lessonOf(mod, "r1", "r1-in", r1)}
	w := newMany(t, mod)
	w.wantChoice(lessonChoice(lid(module, "r1"), ReasonAuthorOrder))
	w.in.Runs = []Run{{Lesson: lid(module, "r1"), Skills: []string{r1}}}
	w.in.Last, w.in.HasLast = Last{Task: tid(module, "r1-in"), Left: true, InLesson: true, Lesson: lid(module, "r1")}, true
	w.wantChoice(Choice{Kind: ChoiceNothing}) // both tasks wait for the lesson of r1

	w.passRun(mod.Lessons[0], 100)
	w.wantChoice(taskChoice(tid(module, "z-r1"), ReasonAuthorOrder)) // r1 before r2 in the graph, though a-r2 is less by ID
	permuted(t, w)
	w.in.RubricDone = []RubricDone{{Task: tid(module, "z-r1"), At: 200}}
	w.wantChoice(taskChoice(tid(module, "a-r2"), ReasonAuthorOrder))
	w.in.RubricDone = append(w.in.RubricDone, RubricDone{Task: tid(module, "a-r2"), At: 300})
	w.wantChoice(Choice{Kind: ChoiceNothing})
}

// A prerequisite with a rubric or none is counted by a lesson that teaches it passed — on a
// package before too —, or by a task where it is primary closed by "Done", which is read by its
// description in a pool or a lesson; one with nothing to teach it holds nothing, nor does a
// prerequisite with a verdict with neither a lesson nor a task the rule could give it.
func TestPrerequisiteCounted(t *testing.T) {
	module := "example/counted"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	essay, reading, after, tail, a, b, c := sid(module, "essay"), sid(module, "reading"), sid(module, "after"), sid(module, "tail"),
		sid(module, "a"), sid(module, "b"), sid(module, "c")
	build := func() Module {
		m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
			{ID: essay, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
			{ID: reading, Assessability: learnermodel.AssessabilityNone, Tier: 1},
			{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
			{ID: c, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{reading}, Tier: 2},
			{ID: after, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{essay}, Tier: 2},
			{ID: tail, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{reading}, Tier: 2},
			{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}, Tier: 2},
		}, []learnermodel.Task{
			exampleTask(module, "essay-1", band, 1, essay), exampleTask(module, "after-1", band, 1, after),
			exampleTask(module, "a-1", band, 1, a), exampleTask(module, "ac", band, 1, a, c), exampleTask(module, "b-1", band, 1, b),
		})
		m.Lessons = []Lesson{lessonOf(m, "essay", "essay-", essay), lessonOf(m, "reading", "", reading)}
		return m
	}
	m := build()
	w := newMany(t, m)
	for _, id := range []string{after, c, b} {
		if w.open(id) {
			t.Errorf("%s is open, want it held", id)
		}
	}
	w.in.Runs = []Run{{Lesson: lid(module, "essay-before"), Skills: []string{essay}, End: RunPassed, EndedAt: 10}}
	if !w.open(after) {
		t.Errorf("%s is held after a lesson on %s passed on a package before, want it open", after, essay)
	}
	w.in.Runs = nil
	w.in.RubricDone = []RubricDone{{Task: tid(module, "unknown"), At: 10}}
	if w.open(after) {
		t.Errorf("%s is open by a task closed by Done that no pool or lesson describes, want it held", after)
	}
	w.in.RubricDone = []RubricDone{{Task: tid(module, "essay-1"), At: 10}}
	if !w.open(after) {
		t.Errorf("%s is held after a task on %s closed by Done, want it open", after, essay)
	}
	w.in.Modules[0].Pool = slices.DeleteFunc(slices.Clone(w.in.Modules[0].Pool), func(c Candidate) bool { return c.Task.ID == tid(module, "essay-1") })
	w.p[0] = w.p[0][:len(w.in.Modules[0].Pool)]
	if !w.open(after) {
		t.Errorf("%s is held after a task on %s closed by Done, described by its lesson, want it open", after, essay)
	}

	// Nothing to teach: no lesson and no task on the prerequisite.
	m = build()
	m.Lessons = nil
	m.Pool = slices.DeleteFunc(m.Pool, func(c Candidate) bool {
		return c.Task.ID == tid(module, "essay-1") || c.Task.ID == tid(module, "a-1")
	})
	w = newMany(t, m)
	for _, id := range []string{after, c, tail} {
		if !w.open(id) {
			t.Errorf("%s is held by a prerequisite with nothing to teach it, want it open", id)
		}
	}
	if !w.open(b) {
		t.Errorf("%s is held by a, whose one task ac has c closed in the open without nothing to teach, want it open", b)
	}
	m.Lessons = []Lesson{{ID: lid(module, "a"), Skills: []string{a}}}
	w = newMany(t, m)
	if w.open(b) {
		t.Errorf("%s is open with a lesson on a, want it held", b)
	}
}

// mapAuthorWorld is a module where r, with a rubric, has lessons a and b and a task in a lesson
// and one outside them, and r2, with a rubric, needs r and has lesson c.
func mapAuthorWorld(t *testing.T) (*many, Module) {
	t.Helper()
	module := "example/map-author"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	r, r2 := sid(module, "r"), sid(module, "r2")
	mod := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: r, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
		{ID: r2, Assessability: learnermodel.AssessabilityRubric, Prerequisites: []string{r}, Tier: 2},
	}, []learnermodel.Task{exampleTask(module, "r-in", band, 1, r), exampleTask(module, "r-out", band, 1, r)})
	mod.Lessons = []Lesson{lessonOf(mod, "a", "r-in", r), lessonOf(mod, "b", "", r), lessonOf(mod, "c", "", r2)}
	return newMany(t, mod), mod
}

// The skill chosen on the map with a rubric or none: a lesson not passed that teaches it — not
// waiting for those before it —, then a task with a rubric on it outside the lessons with no
// "Done", then the lesson whose latest run "passed" ended longest ago, then the task whose latest
// "Done" closed longest ago, one never done first.
func TestMapAuthor(t *testing.T) {
	w, mod := mapAuthorWorld(t)
	module := mod.Graph.Module
	w.in.MapSkill = sid(module, "r2")
	w.wantChoice(lessonChoice(lid(module, "c"), ReasonSkillChosen))
	w.in.MapSkill = sid(module, "r")
	w.wantChoice(lessonChoice(lid(module, "a"), ReasonSkillChosen))
	w.passRun(mod.Lessons[0], 200)
	w.wantChoice(lessonChoice(lid(module, "b"), ReasonSkillChosen))
	w.passRun(mod.Lessons[1], 100)
	w.wantChoice(taskChoice(tid(module, "r-out"), ReasonSkillChosen))
	w.in.RubricDone = []RubricDone{{Task: tid(module, "r-out"), At: 250}}
	w.wantChoice(lessonChoice(lid(module, "b"), ReasonSkillChosen)) // passed longest ago
	permuted(t, w)
	w.in.Runs = append(w.in.Runs, Run{Lesson: lid(module, "b"), Skills: []string{sid(module, "r")}})
	w.wantChoice(lessonChoice(lid(module, "b"), ReasonOpenRun))

	w.in.Runs = nil
	w.in.Modules[0].Lessons = nil
	w.wantChoice(taskChoice(tid(module, "r-in"), ReasonSkillChosen)) // never done
	w.in.RubricDone = append(w.in.RubricDone, RubricDone{Task: tid(module, "r-in"), At: 300})
	w.wantChoice(taskChoice(tid(module, "r-out"), ReasonSkillChosen))
	w.in.RubricDone = append(w.in.RubricDone, RubricDone{Task: tid(module, "r-out"), At: 400}) // its latest "Done" counts
	w.wantChoice(taskChoice(tid(module, "r-in"), ReasonSkillChosen))
	w.in.Modules[0].Pool, w.p[0] = nil, nil
	w.wantChoice(Choice{Kind: ChoiceNothing})
}

// The skill chosen on the map with a verdict gets the lesson on it, closed and whatever the state
// of the lesson's other skills; with no task in the pool, any lesson on it suits.
func TestMapLesson(t *testing.T) {
	m := olympiadLessons(t)
	pigeon := sid(olympiadMath, "pigeonhole-basic")
	m.Lessons = append(m.Lessons, Lesson{ID: lid(olympiadMath, "a-pigeon-enum"), Skills: []string{pigeon, sid(olympiadMath, "combinatorics-enumeration")}})
	w := newMany(t, m)
	w.in.MapSkill = pigeon
	w.probabilities([]px{{"pigeon-12-*", 0.95}, {"pigeon-34-d1-1", 0.78}, {"pigeon-34-d3-1", 0.6}, {"pigeon-34-d5-1", 0.4}, {"pigeon-56-*", 0.3}})
	w.wantChoice(lessonChoice(lid(olympiadMath, "a-pigeon-enum"), ReasonSkillChosen))
	w.in.Modules[0].Lessons = m.Lessons[:len(m.Lessons)-1]
	w.wantChoice(lessonChoice(lid(olympiadMath, "pigeonhole-basic-grades-3-4"), ReasonSkillChosen))
	w.passRun(w.in.Modules[0].Lessons[slices.IndexFunc(m.Lessons, func(l Lesson) bool { return l.ID == lid(olympiadMath, "pigeonhole-basic-grades-3-4") })], 100)
	w.wantChoice(taskChoice(tid(olympiadMath, "pigeon-34-d1-1"), ReasonSkillChosen))

	w.in.Modules[0].Pool, w.p[0] = nil, nil
	w.wantChoice(lessonChoice(lid(olympiadMath, "pigeonhole-basic-grades-1-2"), ReasonSkillChosen))
}

// The order of the lessons: one after another by the graph, transitively through the module's
// prerequisites, but not after a lesson of one of its own skills; an anchor skill orders nothing.
func TestLessonAfter(t *testing.T) {
	module := "example/after"
	a, b, c := sid(module, "a"), sid(module, "b"), sid(module, "c")
	g, err := learnermodel.NewGraph(module, nil, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}, Tier: 2},
		{ID: c, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{b, "core.x"}, Tier: 3},
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	ms := newModuleSkills(g)
	for _, tc := range []struct {
		x, y []string
		want bool
	}{
		{[]string{c}, []string{a}, true},
		{[]string{c}, []string{b}, true},
		{[]string{a}, []string{c}, false},
		{[]string{b, a}, []string{a}, false},
		{[]string{c, b}, []string{a, b}, true},
		{[]string{c}, []string{"core.x"}, false},
		{[]string{"gone"}, []string{a}, false},
	} {
		if got := ms.after(tc.x, tc.y); got != tc.want {
			t.Errorf("after(%v, %v) = %v, want %v", tc.x, tc.y, got, tc.want)
		}
	}
	if got := ms.tier([]string{c, b}); got != 2 {
		t.Errorf("tier of c and b = %d, want 2", got)
	}
}

// permuted checks that the order of the input's lists, those of the lessons and runs among them,
// does not change the choice.
func permuted(t *testing.T, w *many) {
	t.Helper()
	want := ready(w.in, w.p).choose()
	rng := rand.New(rand.NewPCG(19, 3))
	for range 10 {
		in := w.in
		p := slices.Clone(w.p)
		in.Modules = slices.Clone(in.Modules)
		order := rng.Perm(len(in.Modules))
		mods, ps := make([]Module, len(order)), make([][]learnermodel.TaskReadout, len(order))
		for k, j := range order {
			m := in.Modules[j]
			m.Lessons = shuffled(rng, m.Lessons)
			for i := range m.Lessons {
				m.Lessons[i].Skills = shuffled(rng, m.Lessons[i].Skills)
				m.Lessons[i].Tasks = shuffled(rng, m.Lessons[i].Tasks)
			}
			m.Graph.Skills = shuffled(rng, m.Graph.Skills)
			perm := rng.Perm(len(m.Pool))
			pool, pp := make([]Candidate, len(perm)), make([]learnermodel.TaskReadout, len(perm))
			for a, b := range perm {
				pool[a], pp[a] = m.Pool[b], p[j][b]
			}
			m.Pool = pool
			mods[k], ps[k] = m, pp
		}
		in.Modules = mods
		in.Runs = shuffled(rng, in.Runs)
		for i := range in.Runs {
			in.Runs[i].Skills = shuffled(rng, in.Runs[i].Skills)
		}
		in.RubricDone = shuffled(rng, in.RubricDone)
		in.Summary.Skills = shuffled(rng, in.Summary.Skills)
		in.Summary.Tasks = shuffled(rng, in.Summary.Tasks)
		if got := ready(in, ps).choose(); got != want {
			t.Fatalf("chose %+v in another order, want %+v", got, want)
		}
	}
}

// On the map, a task with a rubric whose solution went to the device less than 24 hours ago waits
// for the second pass, as a candidate does; a task never closed by "Done" comes before one that was.
func TestMapAuthorPause(t *testing.T) {
	w, mod := mapAuthorWorld(t)
	module := mod.Graph.Module
	w.in.Modules[0].Lessons = nil
	w.in.MapSkill = sid(module, "r")
	w.in.RubricDone = []RubricDone{{Task: tid(module, "r-out"), At: 10}}
	w.solutions(TaskSolution{tid(module, "r-in"), exampleTime - 10})
	if c, ok := ready(w.in, w.p).pass(false); !ok || c != taskChoice(tid(module, "r-out"), ReasonSkillChosen) {
		t.Errorf("first pass: %+v, %v; want r-out, r-in on pause", c, ok)
	}
	w.in.Modules[0].Pool, w.p[0] = w.in.Modules[0].Pool[:1], w.p[0][:1]
	w.wantChoice(taskChoice(tid(module, "r-in"), ReasonSkillChosen))
	for _, tc := range []struct {
		a, b bool
		want int
	}{{false, false, 0}, {true, true, 0}, {false, true, -1}, {true, false, 1}} {
		if got := falseFirst(tc.a, tc.b); got != tc.want {
			t.Errorf("falseFirst(%v, %v) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

// A prerequisite the graph of the module does not have holds the skill: the rule cannot tell
// whether it is counted.
func TestPrerequisiteOutsideGraph(t *testing.T) {
	module := "example/outside"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	a := sid(module, "a")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{sid(module, "gone")}, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "a-1", band, 1, a)})
	if w := newMany(t, m); w.open(a) {
		t.Errorf("%s is open with a prerequisite outside the graph, want it held", a)
	}
}

// Found by the review of T19. Two open runs on a package before may each come after the other by
// the current graph: branch 1 still gives one, by tier and ID.
func TestOpenRunsEachAfterTheOther(t *testing.T) {
	module := "example/old"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	x, y, z, wr := sid(module, "x"), sid(module, "y"), sid(module, "z"), sid(module, "w")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: y, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: wr, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: x, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{y}, Tier: 2},
		{ID: z, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{wr}, Tier: 2},
	}, []learnermodel.Task{exampleTask(module, "y-1", band, 1, y)})
	w := newMany(t, m)
	w.in.Runs = []Run{{Lesson: lid(module, "b"), Skills: []string{x, wr}}, {Lesson: lid(module, "a"), Skills: []string{y, z}}}
	w.wantChoice(lessonChoice(lid(module, "a"), ReasonOpenRun))
}

// Found by the review of T19. Off the map, the lesson on a skill is none when it names another
// skill not handed out that awaits its lesson of the author's order: that skill comes by it.
func TestLessonOnSkillAwaits(t *testing.T) {
	module := "example/awaits"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	x, y, r := sid(module, "x"), sid(module, "y"), sid(module, "r")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: x, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: y, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: r, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "x-1", band, 1, x), exampleTask(module, "y-1", band, 1, y)})
	m.Lessons = []Lesson{{ID: lid(module, "xy"), Skills: []string{x, y}}, {ID: lid(module, "yr"), Skills: []string{y, r}}}
	w := newMany(t, m)
	w.in.Runs = []Run{{Lesson: lid(module, "yr"), Skills: []string{y, r}}}
	w.in.Last, w.in.HasLast = Last{Task: tid(module, "x-1"), Left: true, InLesson: true, Lesson: lid(module, "yr")}, true
	w.wantChoice(taskChoice(tid(module, "x-1"), ReasonNewSkill))
	w.in.MapSkill = x
	w.wantChoice(lessonChoice(lid(module, "xy"), ReasonSkillChosen))
	w.in.MapSkill = ""
	w.issue(y, 10)
	w.wantChoice(lessonChoice(lid(module, "xy"), ReasonNewSkill))
}

// Found by the review of T19. A task that has a skill without assessment primary is no task with a
// rubric on it: the map does not give it, and it teaches that skill nothing. A lesson passed
// teaches nothing more either.
func TestNothingToTeachRubricAndPassed(t *testing.T) {
	module := "example/teach"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	reading, after, v, c, d := sid(module, "reading"), sid(module, "after"), sid(module, "v"), sid(module, "c"), sid(module, "d")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: reading, Assessability: learnermodel.AssessabilityNone, Tier: 1},
		{ID: v, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: after, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{reading}, Tier: 2},
		{ID: c, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{reading}, Tier: 2},
		{ID: d, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{v}, Tier: 2},
	}, []learnermodel.Task{exampleTask(module, "rv", band, 1, reading, v), exampleTask(module, "vc", band, 1, v, c)})
	w := newMany(t, m)
	if !w.open(after) {
		t.Errorf("%s is held by %s, whose one task has no rubric, want it open", after, reading)
	}
	w.in.MapSkill = reading
	w.wantChoice(Choice{Kind: ChoiceNothing})

	// v's one task, with d primary beside it, teaches v nothing while d is closed.
	w.in.MapSkill = ""
	w.in.Modules[0].Pool, w.p[0] = w.in.Modules[0].Pool[:1], w.p[0][:1]
	w.in.Modules[0].Pool[0].Task.Primary = []string{v, d}
	w.in.Modules[0].Lessons = []Lesson{{ID: lid(module, "v"), Skills: []string{v}}}
	if w.open(d) {
		t.Errorf("%s is open with a lesson on %s not passed, want it held", d, v)
	}
	w.in.Runs = []Run{{Lesson: lid(module, "v"), Skills: []string{v}, End: RunPassed, EndedAt: 10}}
	if !w.open(d) {
		t.Errorf("%s is held by %s, whose lesson is passed and whose one task has %s closed, want it open", d, v, d)
	}
}
