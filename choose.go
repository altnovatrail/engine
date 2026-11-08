package engine

import (
	"cmp"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Choose picks what the learner is given next, and the reason code of the choice
// (docs/architecture.md, section 3). It refuses, with ErrInvalid, an input whose rules it relies
// on are broken: a module, a task of the pool or a lesson twice by ID; a task of the pool or a
// lesson that is not its module's by ID; a skill of a lesson that is not in the graph of its
// module; a run's end outside its list, an open run with a time of its end, and two open runs of
// one lesson; a skill, a module or a task twice in the summary, and a reason code of the summary
// outside its list; a module of the engine with no reading of the learner model; and a
// candidate task the reading refuses, or one with the P of an item outside [0, 1], NaN and the
// infinities included. The checks run in one order and each names the least ID that breaks it, so
// that the order of the input's lists does not change the error either.
//
// The rule is still to come (RUN.md, T16 and on): an input that passes its checks gives nothing
// to give.
func Choose(in Input) (Choice, error) {
	if _, err := prepare(in); err != nil {
		return Choice{}, err
	}
	return Choice{Kind: ChoiceNothing}, nil
}

// prepared is the input once its checks have passed, with what the rule reads beside it: the P
// of each candidate, by the module's place in Input.Modules and the candidate's in its pool.
type prepared struct {
	in Input
	p  [][]learnermodel.TaskReadout
}

func prepare(in Input) (prepared, error) {
	for _, check := range []func(Input) error{
		checkIDs, checkOwnModule, checkLessonSkills, checkRuns, checkSummary, checkReadout,
	} {
		if err := check(in); err != nil {
			return prepared{}, err
		}
	}
	p, err := probabilities(in)
	if err != nil {
		return prepared{}, err
	}
	return prepared{in: in, p: p}, nil
}

// checkIDs refuses a module, a task of the pool or a lesson twice by its full ID, across the
// modules.
func checkIDs(in Input) error {
	var modules, tasks, lessons []string
	for _, m := range in.Modules {
		modules = append(modules, m.Graph.Module)
		for _, c := range m.Pool {
			tasks = append(tasks, c.Task.ID)
		}
		for _, l := range m.Lessons {
			lessons = append(lessons, l.ID)
		}
	}
	for _, list := range []struct {
		what string
		ids  []string
	}{{"module", modules}, {"task of the pool", tasks}, {"lesson", lessons}} {
		if id, ok := firstDuplicate(list.ids); ok {
			return invalid(list.what + " " + shown(id) + " comes twice")
		}
	}
	return nil
}

// checkOwnModule refuses a task of the pool or a lesson whose full ID does not begin with its
// module's: the learner model reads a candidate by the module its ID names, and the candidate
// would get the P of another module's graph.
func checkOwnModule(in Input) error {
	for _, i := range byID(len(in.Modules), func(i int) string { return in.Modules[i].Graph.Module }) {
		m := in.Modules[i]
		prefix := m.Graph.Module + "/"
		for _, j := range byID(len(m.Pool), func(j int) string { return m.Pool[j].Task.ID }) {
			if id := m.Pool[j].Task.ID; !strings.HasPrefix(id, prefix) {
				return invalid("task " + shown(id) + " of the pool of module " + shown(m.Graph.Module) + " is not that module's")
			}
		}
		for _, j := range byID(len(m.Lessons), func(j int) string { return m.Lessons[j].ID }) {
			if id := m.Lessons[j].ID; !strings.HasPrefix(id, prefix) {
				return invalid("lesson " + shown(id) + " of module " + shown(m.Graph.Module) + " is not that module's")
			}
		}
	}
	return nil
}

// checkLessonSkills refuses a skill of a lesson that the graph of the lesson's module has not.
func checkLessonSkills(in Input) error {
	for _, i := range byID(len(in.Modules), func(i int) string { return in.Modules[i].Graph.Module }) {
		m := in.Modules[i]
		skills := make([]string, 0, len(m.Graph.Skills))
		for _, s := range m.Graph.Skills {
			skills = append(skills, s.ID)
		}
		slices.Sort(skills)
		for _, j := range byID(len(m.Lessons), func(j int) string { return m.Lessons[j].ID }) {
			l := m.Lessons[j]
			for _, s := range sortedCopy(l.Skills) {
				if _, found := slices.BinarySearch(skills, s); !found {
					return invalid("skill " + shown(s) + " of lesson " + shown(l.ID) + " is not in the graph of module " + shown(m.Graph.Module))
				}
			}
		}
	}
	return nil
}

// checkRuns refuses a run whose end is not on the list, an open run with a time of its end, and
// two open runs of one lesson: the runs of one lesson repeat by right, but only one of them is
// open at a time.
func checkRuns(in Input) error {
	runs := slices.Clone(in.Runs)
	slices.SortFunc(runs, func(a, b Run) int {
		return cmp.Or(strings.Compare(a.Lesson, b.Lesson), cmp.Compare(a.End, b.End), cmp.Compare(a.EndedAt, b.EndedAt))
	})
	var open []string
	for _, r := range runs {
		if r.End < RunOpen || r.End > RunAbandoned {
			return invalid("run of lesson " + shown(r.Lesson) + " ended by " + itoa(int64(r.End)) + ", which RunEnd does not name")
		}
		if r.End == RunOpen && r.EndedAt != 0 {
			return invalid("open run of lesson " + shown(r.Lesson) + " has a time of its end, " + itoa(r.EndedAt))
		}
		if r.End == RunOpen {
			open = append(open, r.Lesson)
		}
	}
	if id, ok := firstDuplicate(open); ok {
		return invalid("lesson " + shown(id) + " has two open runs")
	}
	return nil
}

// checkSummary refuses a skill, a module or a task twice in the summary, whose every row is the
// one fact of its ID, and a reason code of the summary that is not on the list.
func checkSummary(in Input) error {
	var skillIDs, modules, tasks []string
	for _, s := range in.Summary.Skills {
		skillIDs = append(skillIDs, s.Skill)
	}
	for _, m := range in.Summary.Modules {
		modules = append(modules, m.Module)
	}
	for _, t := range in.Summary.Tasks {
		tasks = append(tasks, t.Task)
	}
	for _, list := range []struct {
		what string
		ids  []string
	}{{"skill", skillIDs}, {"module", modules}, {"task", tasks}} {
		if id, ok := firstDuplicate(list.ids); ok {
			return invalid(list.what + " " + shown(id) + " comes twice in the summary")
		}
	}
	for _, i := range byID(len(in.Summary.Skills), func(i int) string { return in.Summary.Skills[i].Skill }) {
		s := in.Summary.Skills[i]
		if s.HasLast && (s.LastCode < ReasonUnspecified || s.LastCode > ReasonRotation) {
			return invalid("skill " + shown(s.Skill) + " of the summary has the reason code " + itoa(int64(s.LastCode)) + ", which ReasonCode does not name")
		}
	}
	return nil
}

// checkReadout refuses a module of the engine that the learner model's reading has not read.
func checkReadout(in Input) error {
	read := make([]string, 0, len(in.Readout.Modules))
	for _, m := range in.Readout.Modules {
		read = append(read, m.Module)
	}
	slices.Sort(read)
	for _, i := range byID(len(in.Modules), func(i int) string { return in.Modules[i].Graph.Module }) {
		id := in.Modules[i].Graph.Module
		if _, found := slices.BinarySearch(read, id); !found {
			return invalid("module " + shown(id) + " of the engine has no reading of the learner model")
		}
	}
	return nil
}

// probabilities reads the P of every candidate of the pools once, before the rule, by
// Readout.Task (R09): a P outside [0, 1] is an error whichever branch of the rule would come to
// the candidate. The result is by the input's places; the reading goes by ID.
func probabilities(in Input) ([][]learnermodel.TaskReadout, error) {
	p := make([][]learnermodel.TaskReadout, len(in.Modules))
	for _, i := range byID(len(in.Modules), func(i int) string { return in.Modules[i].Graph.Module }) {
		pool := in.Modules[i].Pool
		p[i] = make([]learnermodel.TaskReadout, len(pool))
		for _, j := range byID(len(pool), func(j int) string { return pool[j].Task.ID }) {
			id := pool[j].Task.ID
			tr, err := in.Readout.Task(pool[j].Task)
			if err != nil {
				return nil, &detailError{detail: "task " + shown(id) + " of the pool", cause: err}
			}
			if err := checkProbabilities(id, tr); err != nil {
				return nil, err
			}
			p[i][j] = tr
		}
	}
	return p, nil
}

// checkProbabilities refuses an item with a verdict whose P is outside [0, 1]: a comparison with
// NaN is always false, and the choice would hang on the order of the comparisons (R10, rule 3).
func checkProbabilities(task string, tr learnermodel.TaskReadout) error {
	for _, it := range tr.Items {
		if it.HasVerdict && !(it.Probability >= 0 && it.Probability <= 1) {
			return invalid("item " + itoa(int64(it.Place)) + " of task " + shown(task) + " has a P outside [0, 1]")
		}
	}
	return nil
}

// firstDuplicate gives the least ID that comes twice, byte by byte.
func firstDuplicate(ids []string) (string, bool) {
	s := sortedCopy(ids)
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			return s[i], true
		}
	}
	return "", false
}

// byID gives the places of n values in the order of their IDs, byte by byte, with a tie in the
// order of the places.
func byID(n int, id func(int) string) []int {
	places := make([]int, n)
	for i := range places {
		places[i] = i
	}
	slices.SortStableFunc(places, func(a, b int) int { return strings.Compare(id(a), id(b)) })
	return places
}

func sortedCopy(s []string) []string {
	c := slices.Clone(s)
	slices.Sort(c)
	return c
}
