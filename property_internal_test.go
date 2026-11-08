package engine

import (
	"math/rand/v2"
	"slices"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The properties of the choice that the stand sim would check and that are properties of the code
// (docs/architecture.md, sections 3.1–3.5, "What the stand will check"), on random inputs over the
// graphs of the release, drawn by PCG with a pinned seed:
//   - the same input gives the same choice, and so does any order of its lists;
//   - the choice has a code when, and only when, it gives a lesson or a task;
//   - it gives a task of a pool or a lesson of a module of the engine;
//   - it gives nothing only when, in the second pass too, no open skill has a candidate — but one
//     that awaits its lesson of the author's order —, the author's order has nothing to give, and
//     on the map the skill chosen there has no task.
//
// Shares, rhythm and "no less often than" are the stand's.

const propertyInputs = 300

func TestChoiceProperties(t *testing.T) {
	rel := readRelease(t)
	rng := rand.New(rand.NewPCG(21, 2026))
	kinds := map[ChoiceKind]int{}
	for k := range propertyInputs {
		in := randomInput(t, rng, rel)
		c, err := Choose(in)
		if err != nil {
			t.Fatalf("input %d: %v", k, err)
		}
		if again, _ := Choose(in); again != c {
			t.Fatalf("input %d: chose %+v, then %+v", k, c, again)
		}
		if got, _ := Choose(shuffledInput(rng, in)); got != c {
			t.Fatalf("input %d: chose %+v, and %+v in another order", k, c, got)
		}
		kinds[c.Kind]++
		switch c.Kind {
		case ChoiceTask, ChoiceLesson:
			if c.Reason == ReasonUnspecified {
				t.Fatalf("input %d: %+v gives with no code", k, c)
			}
			if !givable(in, c) {
				t.Fatalf("input %d: %+v is no task of a pool or lesson of a module of the engine", k, c)
			}
		case ChoiceNothing:
			if c.Reason != ReasonUnspecified || c.ID != "" {
				t.Fatalf("input %d: nothing to give as %+v", k, c)
			}
			if why := somethingToGive(t, in); why != "" {
				t.Fatalf("input %d: nothing to give, though %s", k, why)
			}
		default:
			t.Fatalf("input %d: chose %+v", k, c)
		}
	}
	for _, kind := range []ChoiceKind{ChoiceLesson, ChoiceTask, ChoiceNothing} {
		if kinds[kind] == 0 {
			t.Errorf("no input of %d gives a choice of kind %d: the properties of that kind went unchecked", propertyInputs, kind)
		}
	}
}

func givable(in Input, c Choice) bool {
	return slices.ContainsFunc(in.Modules, func(m Module) bool {
		if c.Kind == ChoiceTask {
			return slices.ContainsFunc(m.Pool, func(cd Candidate) bool { return cd.Task.ID == c.ID })
		}
		return slices.ContainsFunc(m.Lessons, func(l Lesson) bool { return l.ID == c.ID })
	})
}

// somethingToGive tells, in words, what the rule could have given in the second pass, or "".
func somethingToGive(t *testing.T, in Input) string {
	t.Helper()
	pr, err := prepare(in)
	if err != nil {
		t.Fatal(err)
	}
	if in.MapSkill != "" {
		for i, m := range in.Modules {
			if !slices.ContainsFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == in.MapSkill }) {
				continue
			}
			for _, t := range pr.poolOf(i) {
				if _, ok := candidateOf(t, in.MapSkill, nil); ok || slices.Contains(pr.mods[i].rubricOf(t.c.Task), in.MapSkill) {
					return "task " + t.c.Task.ID + " is on the skill chosen on the map"
				}
			}
			// A lesson on it with no open run: of a skill with a verdict, one whose skills all have a
			// verdict and that is not passed — with no task on the skill, any band suits —; of one with
			// a rubric or none, any.
			verdict := pr.mods[i].verdict(in.MapSkill)
			for _, l := range m.Lessons {
				if slices.Contains(l.Skills, in.MapSkill) && !pr.isRunning(l.ID) &&
					(!verdict || !pr.mods[i].authorOrder(l.Skills) && !pr.isPassed(l.ID)) {
					return "lesson " + l.ID + " is on the skill chosen on the map"
				}
			}
		}
		return ""
	}
	for i := range in.Modules {
		for _, v := range pr.skills(i, true) {
			if v.open && len(v.candidates) > 0 && !v.awaitsLesson {
				return "skill " + v.id + " is open with a candidate"
			}
		}
		if c, ok := pr.authorOrder(i, true); ok {
			return "the author's order gives " + c.ID
		}
	}
	return ""
}

// randomInput draws an input over one to three modules of the release, as a learner of grade 3
// under 13 gets them: a state of the learner model, a summary, runs, attempts closed by "Done", a
// last attempt, and now and then a skill chosen on the map.
func randomInput(t *testing.T, rng *rand.Rand, rel map[string]release) Input {
	t.Helper()
	ids := []string{releaseOlympiad, releaseGrade3, releaseEssays}
	rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	ids = ids[:1+rng.IntN(3)]
	g := golden{grade: learnerGrade3}
	days := func(n int) int64 { return goldenTime - int64(rng.IntN(n*86_400)) }
	var tasks []string
	for _, id := range ids {
		m := rel[id].forAge("under_13")
		if rng.IntN(3) == 0 {
			// The AI closed, or a pool api narrowed: no lessons, and some of the tasks, or none.
			m.Lessons = nil
			m.Pool = slices.DeleteFunc(slices.Clone(m.Pool), func(Candidate) bool { return rng.IntN(4) > 0 })
		}
		g.modules = append(g.modules, m)
		for _, c := range m.Pool {
			tasks = append(tasks, c.Task.ID)
		}
		if rng.IntN(5) > 0 {
			ms := learnermodel.ModuleState{Module: id, Level: float64(rng.Float64()*7) - 3, AnswerWeight: float64(rng.IntN(13)),
				StartBand: learnermodel.BandRef{Place: 1, ID: m.Graph.Bands[0].ID}}
			// In the trial series the last evidence is in it too; past it, it may be the one that ended it.
			ms.LastInTrial = ms.AnswerWeight < 5 || rng.IntN(4) == 0
			for _, s := range m.Graph.Skills {
				if s.Assessability != learnermodel.AssessabilityVerdict || rng.IntN(3) == 0 {
					continue
				}
				ss := practised(s.ID, float64(rng.Float64()*6)-3, 1+int64(rng.IntN(10)), days(30))
				ss.Stability = 1 + float64(rng.Float64()*60)
				if rng.IntN(3) == 0 && len(s.BandsTaught) > 0 {
					ss.Mastered, ss.MasteryBand = true, s.BandsTaught[rng.IntN(len(s.BandsTaught))]
				}
				ms.Skills = append(ms.Skills, ss)
			}
			g.state.Modules = append(g.state.Modules, ms)
		}
		for _, s := range m.Graph.Skills {
			if rng.IntN(2) == 0 {
				g.in.Summary.Skills = append(g.in.Summary.Skills, SkillIssue{Skill: s.ID, Issued: true, LastIssued: days(20)})
			}
		}
		if rng.IntN(3) > 0 {
			g.in.Summary.Modules = append(g.in.Summary.Modules, ModuleChoices{Module: id, Choices: int64(rng.IntN(20))})
		}
		for _, l := range m.Lessons {
			switch rng.IntN(4) {
			case 0:
				g.in.Runs = append(g.in.Runs, passed(l, days(20)))
			case 1:
				g.in.Runs = append(g.in.Runs, Run{Lesson: l.ID, Skills: l.Skills, End: RunEnd(2 + rng.IntN(3)), EndedAt: days(20)})
			}
		}
		if len(m.Lessons) > 0 && rng.IntN(12) == 0 {
			l := m.Lessons[rng.IntN(len(m.Lessons))]
			g.in.Runs = append(g.in.Runs, Run{Lesson: l.ID, Skills: l.Skills})
		}
		ms := newModuleSkills(m.Graph)
		for _, c := range m.Pool {
			if len(ms.rubricOf(c.Task)) > 0 && rng.IntN(2) == 0 {
				g.in.RubricDone = append(g.in.RubricDone, RubricDone{Task: c.Task.ID, At: days(20)})
			}
		}
	}
	for _, task := range tasks {
		if rng.IntN(3) == 0 {
			g.in.Summary.Tasks = append(g.in.Summary.Tasks, TaskSolution{Task: task, SolutionSentAt: days(2)})
		}
	}
	if len(tasks) > 0 && rng.IntN(4) > 0 {
		task := tasks[rng.IntN(len(tasks))]
		g.in.Last, g.in.HasLast = Last{Task: task, Left: rng.IntN(5) == 0}, true
		for _, m := range g.modules {
			for _, c := range m.Pool {
				if c.Task.ID == task {
					for _, s := range slices.Compact(sortedCopy(primaries(c.Task))) {
						n := int32(rng.IntN(4))
						g.in.Last.Results = append(g.in.Last.Results, SkillResult{Skill: s, Answered: n, FirstRight: int32(rng.IntN(int(n) + 1)), SolutionShown: rng.IntN(3) == 0})
					}
				}
			}
		}
	}
	if rng.IntN(5) == 0 {
		m := g.modules[rng.IntN(len(g.modules))]
		g.in.MapSkill = m.Graph.Skills[rng.IntN(len(m.Graph.Skills))].ID
	}
	return g.input(t)
}

// shuffledInput gives the input with its lists, and the lists inside them, in another order.
func shuffledInput(rng *rand.Rand, in Input) Input {
	in.Modules = shuffled(rng, in.Modules)
	for i := range in.Modules {
		m := &in.Modules[i]
		m.Pool, m.Lessons = shuffled(rng, m.Pool), shuffled(rng, m.Lessons)
		m.Graph.Skills = shuffled(rng, m.Graph.Skills)
		for j := range m.Lessons {
			m.Lessons[j].Skills = shuffled(rng, m.Lessons[j].Skills)
		}
	}
	in.Runs, in.RubricDone = shuffled(rng, in.Runs), shuffled(rng, in.RubricDone)
	for i := range in.Runs {
		in.Runs[i].Skills = shuffled(rng, in.Runs[i].Skills)
	}
	in.Last.Results = shuffled(rng, in.Last.Results)
	in.Summary.Skills, in.Summary.Modules, in.Summary.Tasks = shuffled(rng, in.Summary.Skills), shuffled(rng, in.Summary.Modules), shuffled(rng, in.Summary.Tasks)
	return in
}
