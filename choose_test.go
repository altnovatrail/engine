package engine_test

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/engine"
	learnermodel "github.com/altnovatrail/learner-model"
)

// Two synthetic modules, each with two skills, a task on each skill, a lesson on the first and a
// lesson on both.
var fixtureModules = []string{"example/olympiad", "example/workbook"}

func skillID(module, name string) string { return module + "/skill/" + name }

func fixtureTask(module, name, skill string) learnermodel.Task {
	t := learnermodel.Task{
		ID: module + "/task/" + name, Band: learnermodel.BandRef{Place: 1, ID: "grade-3"}, Difficulty: 3,
		Primary: []string{skill},
	}
	for i := range 2 {
		t.Items = append(t.Items, learnermodel.TaskItem{
			Place: int32(i + 1),
			Slots: []learnermodel.Slot{{Check: learnermodel.CheckKindChoice, Options: 3, HasOptions: true}},
		})
	}
	return t
}

func fixtureModule(t *testing.T, module string) engine.Module {
	t.Helper()
	a, b := skillID(module, "a"), skillID(module, "b")
	tasks := []learnermodel.Task{fixtureTask(module, "a-1", a), fixtureTask(module, "b-1", b)}
	g, err := learnermodel.NewGraph(module,
		[]learnermodel.Band{{Place: 1, ID: "grade-3", Grades: []int32{3}}},
		[]learnermodel.Skill{
			{ID: a, Assessability: learnermodel.AssessabilityVerdict},
			{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}},
		}, nil, tasks)
	if err != nil {
		t.Fatal(err)
	}
	return engine.Module{
		Graph: g,
		Pool: []engine.Candidate{
			{Task: tasks[0], Traps: []engine.ItemTraps{{Place: 1, Traps: []string{module + "/trap/off-by-one"}}}},
			{Task: tasks[1]},
		},
		Lessons: []engine.Lesson{
			{ID: module + "/lesson/b", Skills: []string{b, a}, Tasks: tasks},
			{ID: module + "/lesson/a", Skills: []string{a}, Tasks: tasks[:1]},
		},
	}
}

// readout reads the empty state of the learner model over the graphs of the modules, with the
// modules named as the engine's.
func readout(t *testing.T, modules []engine.Module, engineModules []string) learnermodel.Readout {
	t.Helper()
	graphs := learnermodel.Graphs{}
	for _, m := range modules {
		graphs[m.Graph.Module] = learnermodel.ModuleGraph{Graph: m.Graph, Current: true}
	}
	r, err := learnermodel.Read(learnermodel.State{}, graphs, 1_000_000, learnermodel.Grade{Value: 3, Known: true}, engineModules)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// fixture is a valid input over both modules, with runs, attempts with a rubric, a last attempt
// and a summary.
func fixture(t *testing.T) engine.Input {
	t.Helper()
	var modules []engine.Module
	for _, id := range fixtureModules {
		modules = append(modules, fixtureModule(t, id))
	}
	olympiad := fixtureModules[0]
	return engine.Input{
		Time:    1_000_000,
		Readout: readout(t, modules, fixtureModules),
		Modules: modules,
		Runs: []engine.Run{
			{Lesson: olympiad + "/lesson/a", Skills: []string{skillID(olympiad, "a")}, End: engine.RunPassed, EndedAt: 900_000},
			{Lesson: olympiad + "/lesson/a", Skills: []string{skillID(olympiad, "a")}},
		},
		RubricDone: []engine.RubricDone{{Task: olympiad + "/task/essay", At: 800_000}},
		Last: engine.Last{
			Task:    olympiad + "/task/a-1",
			Results: []engine.SkillResult{{Skill: skillID(olympiad, "a"), Answered: 2, FirstRight: 1}},
		},
		HasLast: true,
		Summary: engine.Summary{
			Skills: []engine.SkillIssue{
				{Skill: skillID(olympiad, "a"), LastIssued: 900_000, Issued: true, LastCode: engine.ReasonNewSkill, HasLast: true},
				{Skill: skillID(olympiad, "b")},
			},
			Modules: []engine.ModuleChoices{{Module: olympiad, Choices: 3}},
			Tasks:   []engine.TaskSolution{{Task: olympiad + "/task/a-1", SolutionSentAt: 900_000}},
		},
	}
}

// A valid input gives nothing to give until the rule comes, with no code, and the runs of one
// lesson repeat when only one of them is open.
func TestChooseValid(t *testing.T) {
	got, err := engine.Choose(fixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if want := (engine.Choice{Kind: engine.ChoiceNothing}); got != want {
		t.Errorf("Choose = %+v, want %+v", got, want)
	}
}

func TestChooseInvalid(t *testing.T) {
	olympiad, workbook := fixtureModules[0], fixtureModules[1]
	for _, tc := range []struct {
		name   string
		change func(t *testing.T, in *engine.Input)
		want   string // a part of the message: what and where
	}{
		{"module twice", func(t *testing.T, in *engine.Input) {
			in.Modules = append(in.Modules, fixtureModule(t, workbook))
		}, `module "example/workbook" comes twice`},
		{"task of the pool twice", func(t *testing.T, in *engine.Input) {
			in.Modules[1].Pool = append(in.Modules[1].Pool, in.Modules[0].Pool[1])
		}, `task of the pool "example/olympiad/task/b-1" comes twice`},
		{"lesson twice", func(t *testing.T, in *engine.Input) {
			in.Modules[0].Lessons = append(in.Modules[0].Lessons, in.Modules[0].Lessons[1])
		}, `lesson "example/olympiad/lesson/a" comes twice`},
		{"task of the pool of another module", func(t *testing.T, in *engine.Input) {
			// The task is the workbook's, and the learner model would read it by the workbook's graph.
			in.Modules[0].Pool = append(in.Modules[0].Pool, engine.Candidate{Task: fixtureTask(workbook, "x", skillID(workbook, "a"))})
		}, `task "example/workbook/task/x" of the pool of module "example/olympiad" is not that module's`},
		{"lesson of another module", func(t *testing.T, in *engine.Input) {
			in.Modules[1].Lessons[0].ID = olympiad + "/lesson/c"
		}, `lesson "example/olympiad/lesson/c" of module "example/workbook" is not that module's`},
		{"skill of a lesson outside its graph", func(t *testing.T, in *engine.Input) {
			// The skill is the other module's: the graph that counts is the lesson's module's.
			in.Modules[0].Lessons[1].Skills = append(in.Modules[0].Lessons[1].Skills, skillID(workbook, "a"))
		}, `skill "example/workbook/skill/a" of lesson "example/olympiad/lesson/a" is not in the graph of module "example/olympiad"`},
		{"traps of an item twice", func(t *testing.T, in *engine.Input) {
			in.Modules[0].Pool[0].Traps = append(in.Modules[0].Pool[0].Traps, engine.ItemTraps{Place: 1})
		}, `item 1 of task "example/olympiad/task/a-1" has its traps twice`},
		{"traps where the task has no item", func(t *testing.T, in *engine.Input) {
			in.Modules[0].Pool[1].Traps = []engine.ItemTraps{{Place: 3}}
		}, `traps of item 3 of task "example/olympiad/task/b-1", which has no item there`},
		{"end of a run off the list", func(t *testing.T, in *engine.Input) {
			in.Runs[0].End = engine.RunAbandoned + 1
		}, `run of lesson "example/olympiad/lesson/a" ended by 5, which RunEnd does not name`},
		{"open run with a time of its end", func(t *testing.T, in *engine.Input) {
			in.Runs[1].EndedAt = 950_000
		}, `open run of lesson "example/olympiad/lesson/a" has a time of its end, 950000`},
		{"two open runs of one lesson", func(t *testing.T, in *engine.Input) {
			in.Runs[0].End, in.Runs[0].EndedAt = engine.RunOpen, 0
		}, `lesson "example/olympiad/lesson/a" has two open runs`},
		{"skill twice in the summary", func(t *testing.T, in *engine.Input) {
			in.Summary.Skills = append(in.Summary.Skills, engine.SkillIssue{Skill: skillID(olympiad, "a")})
		}, `skill "example/olympiad/skill/a" comes twice in the summary`},
		{"module twice in the summary", func(t *testing.T, in *engine.Input) {
			in.Summary.Modules = append(in.Summary.Modules, in.Summary.Modules[0])
		}, `module "example/olympiad" comes twice in the summary`},
		{"task twice in the summary", func(t *testing.T, in *engine.Input) {
			in.Summary.Tasks = append(in.Summary.Tasks, engine.TaskSolution{Task: olympiad + "/task/a-1"})
		}, `task "example/olympiad/task/a-1" comes twice in the summary`},
		{"reason code of the summary off the list", func(t *testing.T, in *engine.Input) {
			in.Summary.Skills[0].LastCode = engine.ReasonRotation + 1
		}, `skill "example/olympiad/skill/a" of the summary has the reason code 10, which ReasonCode does not name`},
		{"skill twice in the results of the last attempt", func(t *testing.T, in *engine.Input) {
			in.Last.Results = append(in.Last.Results, in.Last.Results[0])
		}, `skill "example/olympiad/skill/a" comes twice in the results of the last attempt`},
		{"more first right verdicts than first verdicts", func(t *testing.T, in *engine.Input) {
			in.Last.Results[0].FirstRight = 3
		}, `result on skill "example/olympiad/skill/a" of the last attempt has 3 first right verdicts of 2`},
		{"fewer first right verdicts than none", func(t *testing.T, in *engine.Input) {
			in.Summary.Skills[1].LastResult.FirstRight, in.Summary.Skills[1].HasLast = -1, true
		}, `result on skill "example/olympiad/skill/b" of the summary has -1 first right verdicts of 0`},
		{"results without a meaning", func(t *testing.T, in *engine.Input) {
			// Without HasLast a result means nothing and is not checked: the next check fails.
			in.HasLast = false
			in.Last.Results = append(in.Last.Results, in.Last.Results[0])
			in.Summary.Skills[1].LastResult.FirstRight = -1
			in.Readout = readout(t, in.Modules, []string{olympiad})
		}, `module "example/workbook" of the engine has no reading of the learner model`},
		{"module of the engine with no reading", func(t *testing.T, in *engine.Input) {
			in.Readout = readout(t, in.Modules, []string{olympiad})
		}, `module "example/workbook" of the engine has no reading of the learner model`},
		{"candidate the reading refuses", func(t *testing.T, in *engine.Input) {
			bad := fixtureTask(workbook, "a-2", skillID(workbook, "a"))
			bad.Difficulty = 9
			in.Modules[1].Pool = append(in.Modules[1].Pool, engine.Candidate{Task: bad})
		}, `task "example/workbook/task/a-2" of the pool: learnermodel: invalid input`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := fixture(t)
			tc.change(t, &in)
			_, err := engine.Choose(in)
			if !errors.Is(err, engine.ErrInvalid) {
				t.Fatalf("Choose: %v, want engine.ErrInvalid", err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Choose: %q, want it to say %q", err, tc.want)
			}
		})
	}
}

// The error of learner-model a candidate brought is wrapped: errors.Is finds both kinds.
func TestChooseWrapsReadoutError(t *testing.T) {
	in := fixture(t)
	bad := fixtureTask(fixtureModules[1], "a-2", skillID(fixtureModules[1], "a"))
	bad.Difficulty = 0
	in.Modules[1].Pool = append(in.Modules[1].Pool, engine.Candidate{Task: bad})
	_, err := engine.Choose(in)
	if !errors.Is(err, engine.ErrInvalid) || !errors.Is(err, learnermodel.ErrInvalid) {
		t.Errorf("Choose: %v, want both engine.ErrInvalid and learnermodel.ErrInvalid", err)
	}
}

// The order of the input's lists changes neither the choice nor the error: an input that breaks
// one check twice names the least ID, and the least value, whatever the order. Each broken input
// passes the checks before its own.
func TestChooseOrder(t *testing.T) {
	olympiad, workbook := fixtureModules[0], fixtureModules[1]
	for _, tc := range []struct {
		name   string
		change func(*engine.Input)
		want   string
	}{
		{"valid", func(*engine.Input) {}, ""},
		{"lessons twice", func(in *engine.Input) {
			for i := range in.Modules {
				in.Modules[i].Lessons = append(in.Modules[i].Lessons, in.Modules[i].Lessons[1])
			}
		}, `lesson "example/olympiad/lesson/a" comes twice`},
		{"tasks of other modules", func(in *engine.Input) {
			in.Modules[0].Pool = append(in.Modules[0].Pool, engine.Candidate{Task: fixtureTask(workbook, "y", skillID(workbook, "a"))})
			in.Modules[1].Pool = append(in.Modules[1].Pool, engine.Candidate{Task: fixtureTask(olympiad, "x", skillID(olympiad, "a"))})
		}, `task "example/olympiad/task/x" of the pool of module "example/workbook"`},
		{"traps off their items", func(in *engine.Input) {
			in.Modules[1].Pool[0].Traps = append(in.Modules[1].Pool[0].Traps, engine.ItemTraps{Place: 1})
			in.Modules[0].Pool[1].Traps = []engine.ItemTraps{{Place: 5}, {Place: 4}}
		}, `traps of item 4 of task "example/olympiad/task/b-1"`},
		{"least task across modules whose IDs begin alike", func(in *engine.Input) {
			// "example/olympiad-2/…" is less than "example/olympiad/…": '-' comes before '/'.
			second := "example/olympiad-2"
			in.Modules = append(in.Modules, fixtureModule(t, second))
			in.Readout = readout(t, in.Modules, append(slices.Clone(fixtureModules), second))
			in.Modules[0].Pool[1].Traps = []engine.ItemTraps{{Place: 7}}
			in.Modules[2].Pool[1].Traps = []engine.ItemTraps{{Place: 8}}
		}, `traps of item 8 of task "example/olympiad-2/task/b-1"`},
		{"ends off the list", func(in *engine.Input) {
			in.Runs = append(in.Runs,
				engine.Run{Lesson: workbook + "/lesson/a", End: 7},
				engine.Run{Lesson: workbook + "/lesson/a", End: 6},
				engine.Run{Lesson: workbook + "/lesson/b", End: 5})
		}, `run of lesson "example/workbook/lesson/a" ended by 6`},
		{"open runs with times", func(in *engine.Input) {
			in.Runs = append(in.Runs,
				engine.Run{Lesson: workbook + "/lesson/a", EndedAt: 7},
				engine.Run{Lesson: workbook + "/lesson/a", EndedAt: 6})
		}, `open run of lesson "example/workbook/lesson/a" has a time of its end, 6`},
		{"codes off the list", func(in *engine.Input) {
			in.Summary.Skills[1].LastCode, in.Summary.Skills[1].HasLast = 12, true
			in.Summary.Skills[0].LastCode = 11
		}, `skill "example/olympiad/skill/a" of the summary has the reason code 11`},
		{"results off", func(in *engine.Input) {
			in.Last.Results = append(in.Last.Results,
				engine.SkillResult{Skill: skillID(workbook, "a"), Answered: 1, FirstRight: 5},
				engine.SkillResult{Skill: skillID(olympiad, "b"), Answered: 1, FirstRight: 4})
		}, `result on skill "example/olympiad/skill/b" of the last attempt has 4`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := func() engine.Input {
				in := fixture(t)
				tc.change(&in)
				return in
			}
			wantChoice, wantErr := engine.Choose(input())
			if !strings.Contains(errorText(wantErr), tc.want) || (tc.want == "") != (wantErr == nil) {
				t.Fatalf("Choose: %v, want it to say %q", wantErr, tc.want)
			}
			rng := rand.New(rand.NewPCG(1, 2))
			for range 20 {
				in := input()
				shuffle(rng, in.Modules)
				for i := range in.Modules {
					shuffle(rng, in.Modules[i].Pool)
					shuffle(rng, in.Modules[i].Lessons)
					shuffle(rng, in.Modules[i].Lessons[0].Skills)
					for k := range in.Modules[i].Pool {
						shuffle(rng, in.Modules[i].Pool[k].Traps)
					}
				}
				shuffle(rng, in.Runs)
				shuffle(rng, in.RubricDone)
				shuffle(rng, in.Summary.Skills)
				shuffle(rng, in.Last.Results)
				got, err := engine.Choose(in)
				if got != wantChoice || errorText(err) != errorText(wantErr) {
					t.Fatalf("Choose = %+v, %v; in the input's order it was %+v, %v", got, err, wantChoice, wantErr)
				}
			}
		})
	}
}

func shuffle[T any](rng *rand.Rand, s []T) {
	rng.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}

func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// Choose reads the input and changes nothing in it: the lists it orders, it orders as copies. The
// fixture's runs and the summary's skills come in an order their sorting would change.
func TestChooseLeavesInput(t *testing.T) {
	in := fixture(t)
	slices.Reverse(in.Summary.Skills)
	runs, skills := slices.Clone(in.Runs), slices.Clone(in.Summary.Skills)
	lessonSkills := slices.Clone(in.Modules[0].Lessons[0].Skills)
	if _, err := engine.Choose(in); err != nil {
		t.Fatal(err)
	}
	for i := range runs {
		if in.Runs[i].Lesson != runs[i].Lesson || in.Runs[i].End != runs[i].End || in.Runs[i].EndedAt != runs[i].EndedAt {
			t.Fatalf("Choose reordered the runs of its input: %+v, was %+v", in.Runs, runs)
		}
	}
	for i := range skills {
		if in.Summary.Skills[i].Skill != skills[i].Skill || in.Summary.Skills[i].LastCode != skills[i].LastCode {
			t.Fatalf("Choose reordered the skills of the summary: %+v, was %+v", in.Summary.Skills, skills)
		}
	}
	if !slices.Equal(in.Modules[0].Lessons[0].Skills, lessonSkills) {
		t.Fatalf("Choose reordered the skills of a lesson: %v, was %v", in.Modules[0].Lessons[0].Skills, lessonSkills)
	}
}

// With one module of the engine, Choose runs the rule, on the learner model's own reading: with no
// evidence, the trial series gives a new skill — the one open, as the graph orders it — whatever
// the order of the input's lists.
func TestChooseOneModule(t *testing.T) {
	olympiad := fixtureModules[0]
	input := func() engine.Input {
		in := fixture(t)
		in.Modules = in.Modules[:1]
		in.Readout = readout(t, in.Modules, fixtureModules[:1])
		in.Runs = in.Runs[:1] // no open run
		in.Summary.Skills = nil
		return in
	}
	want := engine.Choice{Kind: engine.ChoiceTask, ID: olympiad + "/task/a-1", Reason: engine.ReasonNewSkill}
	if got, err := engine.Choose(input()); err != nil || got != want {
		t.Fatalf("Choose = %+v, %v; want %+v", got, err, want)
	}
	rng := rand.New(rand.NewPCG(5, 6))
	for range 20 {
		in := input()
		shuffle(rng, in.Modules[0].Pool)
		shuffle(rng, in.Modules[0].Lessons)
		shuffle(rng, in.Runs)
		in.Modules[0].Graph.Skills = slices.Clone(in.Modules[0].Graph.Skills)
		shuffle(rng, in.Modules[0].Graph.Skills)
		in.Summary.Tasks = append(in.Summary.Tasks, engine.TaskSolution{Task: olympiad + "/task/b-1", SolutionSentAt: 1})
		shuffle(rng, in.Summary.Tasks)
		in.Summary.Modules = append(in.Summary.Modules, engine.ModuleChoices{Module: fixtureModules[1], Choices: 11})
		shuffle(rng, in.Summary.Modules)
		in.Last.Results = append(in.Last.Results, engine.SkillResult{Skill: skillID(olympiad, "b"), Answered: 1})
		shuffle(rng, in.Last.Results)
		if got, err := engine.Choose(in); err != nil || got != want {
			t.Fatalf("Choose = %+v, %v in another order; want %+v", got, err, want)
		}
	}

	// A skill chosen on the map and an open run take the branches that come first, and several
	// modules take turns: until they come (RUN.md, T18 and T19), Choose gives nothing to give.
	for _, tc := range []struct {
		name   string
		change func(*engine.Input)
	}{
		{"skill on the map", func(in *engine.Input) { in.MapSkill = skillID(olympiad, "b") }},
		{"open run", func(in *engine.Input) { in.Runs = append(in.Runs, engine.Run{Lesson: olympiad + "/lesson/a"}) }},
		{"several modules", func(in *engine.Input) {
			full := fixture(t)
			in.Modules, in.Readout = full.Modules, full.Readout
		}},
	} {
		in := input()
		tc.change(&in)
		if got, err := engine.Choose(in); err != nil || got != (engine.Choice{Kind: engine.ChoiceNothing}) {
			t.Errorf("%s: Choose = %+v, %v; want nothing to give", tc.name, got, err)
		}
	}
}

// With every candidate on the pause before a seen solution, the second pass gives one, but not the
// task of the last closed attempt; with nothing else, Choose gives nothing to give.
func TestChooseSecondPass(t *testing.T) {
	olympiad := fixtureModules[0]
	in := fixture(t)
	in.Modules = in.Modules[:1]
	in.Readout = readout(t, in.Modules, fixtureModules[:1])
	in.Runs = in.Runs[:1]
	in.Summary.Skills = nil
	in.Summary.Tasks = []engine.TaskSolution{{Task: olympiad + "/task/a-1", SolutionSentAt: in.Time - 60}}
	in.Last.Task = olympiad + "/task/b-1"
	want := engine.Choice{Kind: engine.ChoiceTask, ID: olympiad + "/task/a-1", Reason: engine.ReasonNewSkill}
	if got, err := engine.Choose(in); err != nil || got != want {
		t.Errorf("Choose = %+v, %v; want %+v", got, err, want)
	}
	in.Last.Task = olympiad + "/task/a-1"
	if got, err := engine.Choose(in); err != nil || got != (engine.Choice{Kind: engine.ChoiceNothing}) {
		t.Errorf("Choose = %+v, %v; want nothing to give", got, err)
	}
}
