package engine

import (
	"slices"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The benchmark of the choice on a large pool (docs/measurements.md, section 2.1): olympiad maths
// 0.1.1, 772 tasks, for a learner of grade 3, in four states of the learner model, each read by its
// real Read at goldenTime. Read is api's call, not the engine's: it runs once, before the loop,
// which times Choose alone. The states are synthetic, built from the graph by its order, so that a
// run before a change of the engine and a run after it time the same inputs.

// benchCases are the inputs of the benchmark, with the choice each gives, its ID in the module:
// TestBenchInputs holds them, so that a change that keeps the choice keeps them too.
var benchCases = []struct {
	name   string
	id     string
	input  func(m Module) golden
	kind   ChoiceKind
	reason ReasonCode
}{
	{"new", "lesson/arithmetic-tricks-grades-1-2", benchNew, ChoiceLesson, ReasonNewSkill},
	{"middle", "task/tri-34-d1-2", benchMiddle, ChoiceTask, ReasonQueue},
	{"middle-consolidation", "task/tri-34-d1-2", benchConsolidation, ChoiceTask, ReasonConsolidation},
	{"mastered", "task/tri-12-d1-2", benchMastered, ChoiceTask, ReasonRepetition},
}

// benchCorrections are the corrections of the skills practised, by their place in the graph.
var benchCorrections = []float64{0, 1, -1, 2, 0.5, -0.5, 1.5}

// benchStabilities are the stabilities in days of the skills mastered, by their place: a review
// date falls due for some of them at goldenTime and not for others.
var benchStabilities = []float64{3, 60, 400, 10, 200}

var bandLarge34 = learnermodel.BandRef{Place: 2, ID: "grades-3-4"}

// benchNew is a learner who has just connected the module: no evidence, and nothing handed out.
func benchNew(m Module) golden {
	return golden{modules: []Module{m}, grade: learnerGrade3}
}

// benchMiddle is a learner half way, past the trial series at about the start of grade 3: the
// first ten skills of the graph practised, one a day, their lessons at grades-3-4 passed the day
// before, and the solutions of their tasks at grades-3-4 at even places of the pool sent to the
// device: about half of each skill's, as the pool is ordered by ID and a skill's IDs stand together.
func benchMiddle(m Module) golden {
	skills := graphOrder(m.Graph.Skills)[:10]
	g := golden{modules: []Module{m}, grade: learnerGrade3}
	var states []learnermodel.SkillState
	for k, s := range skills {
		last := goldenTime - int64(len(skills)-k)*day
		states = append(states, practised(s.ID, benchCorrections[k%len(benchCorrections)], 6, last))
		g.in.Summary.Skills = append(g.in.Summary.Skills, issued(last, s.ID)...)
		g.in.Runs = append(g.in.Runs, benchPassed(m, s.ID, "-grades-3-4", last-day)...)
	}
	g.state = learnermodel.State{Modules: []learnermodel.ModuleState{module(m.Graph.Module, 2, bandLarge34, states...)}}
	g.in.Summary.Tasks = benchSent(m, skills, func(k int, t learnermodel.Task) bool { return k%2 == 0 && t.Band.ID == bandLarge34.ID })
	g.in.Summary.Modules = []ModuleChoices{{Module: m.Graph.Module, Choices: 60}}
	return g
}

// benchConsolidation is the learner half way whose last attempt failed on a task of the first
// skill of the graph.
func benchConsolidation(m Module) golden {
	g := benchMiddle(m)
	first := graphOrder(m.Graph.Skills)[0].ID
	i := slices.IndexFunc(m.Pool, func(c Candidate) bool {
		return c.Task.Band.ID == bandLarge34.ID && slices.Contains(c.Task.Primary, first)
	})
	if i >= 0 { // with no such task there is no last attempt, and TestBenchInputs names the case
		g.in.Last, g.in.HasLast = failedOn(m.Pool[i].Task.ID, first), true
	}
	return g
}

// benchMastered is a learner far on: fifteen skills of the graph mastered at grades-3-4, with
// review dates some of which have come, the other four practised, the lessons of every skill
// passed, and the solutions of the tasks at places of the pool not a multiple of three sent to the
// device: about two in three of each skill's.
func benchMastered(m Module) golden {
	skills := graphOrder(m.Graph.Skills)
	g := golden{modules: []Module{m}, grade: learnerGrade3}
	var states []learnermodel.SkillState
	for k, s := range skills {
		last := goldenTime - int64(len(skills)-k)*day
		if k < 15 {
			states = append(states, masteredState(s.ID, bandLarge34.Place, benchStabilities[k%len(benchStabilities)], last))
		} else {
			states = append(states, practised(s.ID, benchCorrections[k%len(benchCorrections)], 8, last))
		}
		g.in.Summary.Skills = append(g.in.Summary.Skills, issued(last, s.ID)...)
		g.in.Runs = append(g.in.Runs, benchPassed(m, s.ID, "", last-day)...)
	}
	g.state = learnermodel.State{Modules: []learnermodel.ModuleState{module(m.Graph.Module, 3, bandLarge34, states...)}}
	g.in.Summary.Tasks = benchSent(m, skills, func(k int, _ learnermodel.Task) bool { return k%3 != 0 })
	g.in.Summary.Modules = []ModuleChoices{{Module: m.Graph.Module, Choices: 301}} // the turn of the repetition, and not of the probe
	return g
}

// benchPassed gives a passed run of each lesson of the module on a skill whose ID ends as given.
func benchPassed(m Module, skill, suffix string, at int64) []Run {
	var out []Run
	for _, l := range m.Lessons {
		if slices.Contains(l.Skills, skill) && strings.HasSuffix(l.ID, suffix) {
			out = append(out, passed(l, at))
		}
	}
	return out
}

// benchSent gives the solutions sent to the device of the tasks of the pool primary on one of the
// skills that keep keeps, by the task's place in the pool, each some days before goldenTime.
func benchSent(m Module, skills []learnermodel.Skill, keep func(int, learnermodel.Task) bool) []TaskSolution {
	var out []TaskSolution
	for k, c := range m.Pool {
		on := slices.ContainsFunc(skills, func(s learnermodel.Skill) bool { return slices.Contains(c.Task.Primary, s.ID) })
		if on && keep(k, c.Task) {
			out = append(out, TaskSolution{Task: c.Task.ID, SolutionSentAt: goldenTime - int64(k%9+2)*day})
		}
	}
	return out
}

// benchInputs gives the inputs of the benchmark, read once.
func benchInputs(tb testing.TB) []Input {
	tb.Helper()
	m := readLarge(tb).forAge("under_13")
	var out []Input
	for _, c := range benchCases {
		out = append(out, c.input(m).input(tb))
	}
	return out
}

// The inputs of the benchmark give the choice they name: so the benchmark times the branches it
// says, before a change of the engine and after it.
func TestBenchInputs(t *testing.T) {
	for k, in := range benchInputs(t) {
		c := benchCases[k]
		got, err := Choose(in)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if want := releaseOlympiad + "/" + c.id; got.Kind != c.kind || got.ID != want || got.Reason != c.reason {
			t.Errorf("%s: %v %s with the code %v, want %v %s with %v", c.name, got.Kind, got.ID, got.Reason, c.kind, want, c.reason)
		}
	}
}

// BenchmarkChoose times one choice on the large pool, in each state:
//
//	go test -run '^$' -bench BenchmarkChoose -benchmem
func BenchmarkChoose(b *testing.B) {
	inputs := benchInputs(b)
	for k, c := range benchCases {
		in := inputs[k]
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := Choose(in); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
