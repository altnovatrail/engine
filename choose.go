package engine

import (
	"cmp"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Choose picks what the learner is given next, and the reason code of the choice
// (docs/architecture.md, section 3). It refuses, with ErrInvalid, an input whose rules it relies
// on are broken: a module, a task of the pool or a lesson twice by ID; a skill of a graph, a task
// of the pool or a lesson that is not its module's by ID; a skill of a lesson that is not in the
// graph of its module; the traps of an item twice, or at a place where its task has no item; a
// run's end outside its list, an open run with a time of its end, and two open runs of one lesson;
// a skill, a module or a task twice in the summary, and a reason code of the summary outside its
// list; a skill twice in the results of the last attempt, and a result with first right verdicts
// fewer than none or more than its first verdicts; a module of the engine with no reading of the
// learner model, and an anchor skill a graph names with none or with a status off its list; and a
// candidate task the reading refuses, or one with the P of an item outside [0, 1], NaN and the
// infinities included. The checks run in one order and each names the least ID that breaks it, so
// that the order of the input's lists does not change the error either.
//
// Choose runs the rule for the skills with a verdict: on the skill chosen on the map, or in "My
// subjects" — consolidation in the module of the last closed attempt, then the modules in turn —,
// then a second pass with the tasks on pause, and gives a task with the code of its branch or
// nothing to give. The open run, the repetition, the author's order and the lessons are still to
// come (RUN.md, T19 and T20): an input with an open run gives nothing to give, since branch 1,
// which comes first for it, is not there yet.
func Choose(in Input) (Choice, error) {
	pr, err := prepare(in)
	if err != nil {
		return Choice{}, err
	}
	return pr.choose(), nil
}

// prepared is the input once its checks have passed, with what the rule reads beside it: the P
// of each candidate, by the module's place in Input.Modules and the candidate's in its pool.
type prepared struct {
	in   Input
	p    [][]learnermodel.TaskReadout
	open []string // the open skills of the modules of the engine, by full ID, sorted
	// The tasks whose solution went to the device, by ID, sorted.
	solutions []TaskSolution
}

// ready puts what the rule reads beside the input: the P of the candidates, and the open skills.
func ready(in Input, p [][]learnermodel.TaskReadout) prepared {
	pr := prepared{in: in, p: p, solutions: sortedBy(in.Summary.Tasks, func(t TaskSolution) string { return t.Task })}
	pr.open = pr.openSkills()
	return pr
}

// isOpen tells whether a skill of the modules of the engine is open.
func (pr prepared) isOpen(id string) bool {
	_, found := slices.BinarySearch(pr.open, id)
	return found
}

func prepare(in Input) (prepared, error) {
	for _, check := range []func(Input) error{
		checkIDs, checkOwnModule, checkLessonSkills, checkTraps, checkRuns, checkSummary, checkResults, checkReadout, checkAnchors,
	} {
		if err := check(in); err != nil {
			return prepared{}, err
		}
	}
	p, err := probabilities(in)
	if err != nil {
		return prepared{}, err
	}
	return ready(in, p), nil
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

// checkOwnModule refuses a skill of a graph, a task of the pool or a lesson whose full ID does not
// begin with its module's: the learner model reads a candidate by the module its ID names, and the
// candidate would get the P of another module's graph; and the rule reads the skills of all the
// modules together — the skill chosen on the map, the anchor skills —, where a skill of two
// modules would be read by the order of the modules.
func checkOwnModule(in Input) error {
	graphs := acrossModules(in, func(m Module) int { return len(m.Graph.Skills) }, func(m Module, j int) string { return m.Graph.Skills[j].ID })
	for _, at := range graphs {
		m := in.Modules[at.module]
		if id := m.Graph.Skills[at.item].ID; !strings.HasPrefix(id, m.Graph.Module+"/") {
			return invalid("skill " + shown(id) + " of the graph of module " + shown(m.Graph.Module) + " is not that module's")
		}
	}
	for _, at := range tasksByID(in) {
		m := in.Modules[at.module]
		if id := m.Pool[at.item].Task.ID; !strings.HasPrefix(id, m.Graph.Module+"/") {
			return invalid("task " + shown(id) + " of the pool of module " + shown(m.Graph.Module) + " is not that module's")
		}
	}
	for _, at := range lessonsByID(in) {
		m := in.Modules[at.module]
		if id := m.Lessons[at.item].ID; !strings.HasPrefix(id, m.Graph.Module+"/") {
			return invalid("lesson " + shown(id) + " of module " + shown(m.Graph.Module) + " is not that module's")
		}
	}
	return nil
}

// checkLessonSkills refuses a skill of a lesson that the graph of the lesson's module has not.
func checkLessonSkills(in Input) error {
	skills := make([][]string, len(in.Modules))
	for i, m := range in.Modules {
		for _, s := range m.Graph.Skills {
			skills[i] = append(skills[i], s.ID)
		}
		slices.Sort(skills[i])
	}
	for _, at := range lessonsByID(in) {
		m := in.Modules[at.module]
		l := m.Lessons[at.item]
		for _, s := range sortedCopy(l.Skills) {
			if _, found := slices.BinarySearch(skills[at.module], s); !found {
				return invalid("skill " + shown(s) + " of lesson " + shown(l.ID) + " is not in the graph of module " + shown(m.Graph.Module))
			}
		}
	}
	return nil
}

// checkTraps refuses the traps of an item twice for a task of the pool, and traps at a place
// where the task has no item.
func checkTraps(in Input) error {
	for _, at := range tasksByID(in) {
		c := in.Modules[at.module].Pool[at.item]
		var places, items []int32
		for _, it := range c.Traps {
			places = append(places, it.Place)
		}
		for _, it := range c.Task.Items {
			items = append(items, it.Place)
		}
		if p, ok := firstDuplicate(places); ok {
			return invalid("item " + itoa(int64(p)) + " of task " + shown(c.Task.ID) + " has its traps twice")
		}
		slices.Sort(items)
		for _, p := range sortedCopy(places) {
			if _, found := slices.BinarySearch(items, p); !found {
				return invalid("traps of item " + itoa(int64(p)) + " of task " + shown(c.Task.ID) + ", which has no item there")
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

// checkResults refuses a skill twice in the results of the last attempt, and a result, there or
// in the summary, whose first right verdicts are fewer than none or more than its first verdicts.
// A result is checked where it has a meaning: of the last attempt when HasLast, and of a skill of
// the summary when its HasLast, as the reason code of the summary is.
func checkResults(in Input) error {
	if in.HasLast {
		var skills []string
		for _, r := range in.Last.Results {
			skills = append(skills, r.Skill)
		}
		if id, ok := firstDuplicate(skills); ok {
			return invalid("skill " + shown(id) + " comes twice in the results of the last attempt")
		}
		for _, i := range byID(len(in.Last.Results), func(i int) string { return in.Last.Results[i].Skill }) {
			if err := checkResult("the last attempt", in.Last.Results[i]); err != nil {
				return err
			}
		}
	}
	for _, i := range byID(len(in.Summary.Skills), func(i int) string { return in.Summary.Skills[i].Skill }) {
		if !in.Summary.Skills[i].HasLast {
			continue
		}
		r := in.Summary.Skills[i].LastResult
		r.Skill = in.Summary.Skills[i].Skill
		if err := checkResult("the summary", r); err != nil {
			return err
		}
	}
	return nil
}

func checkResult(where string, r SkillResult) error {
	if r.FirstRight < 0 || r.FirstRight > r.Answered {
		return invalid("result on skill " + shown(r.Skill) + " of " + where + " has " + itoa(int64(r.FirstRight)) +
			" first right verdicts of " + itoa(int64(r.Answered)))
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

// checkAnchors refuses an anchor skill that a graph of a module of the engine names, as a
// prerequisite or as a skill it teaches, with no reading of the learner model, or with a status
// outside its list: the rule would have to guess whether the anchor holds a skill (the amendment
// of 2026-11-08 to R04).
func checkAnchors(in Input) error {
	var named []string
	for _, m := range in.Modules {
		for _, s := range m.Graph.Skills {
			for _, id := range slices.Concat(s.Prerequisites, s.Teaches) {
				if strings.HasPrefix(id, anchorPrefix) {
					named = append(named, id)
				}
			}
		}
	}
	anchors := sortedBy(in.Readout.Anchors, func(a learnermodel.AnchorReadout) string { return a.Anchor })
	for _, id := range slices.Compact(sortedCopy(named)) {
		a, ok := lookup(anchors, func(a learnermodel.AnchorReadout) string { return a.Anchor }, id)
		if !ok {
			return invalid("anchor skill " + shown(id) + " has no reading of the learner model")
		}
		if a.Status < learnermodel.AnchorHas || a.Status > learnermodel.AnchorUnknown {
			return invalid("anchor skill " + shown(id) + " has the status " + itoa(int64(a.Status)) + ", which the learner model does not name")
		}
	}
	return nil
}

// probabilities reads the P of every candidate of the pools once, before the rule, by
// Readout.Task (R09): a P outside [0, 1] is an error whichever branch of the rule would come to
// the candidate. The result is by the input's places; the reading goes by ID.
func probabilities(in Input) ([][]learnermodel.TaskReadout, error) {
	p := make([][]learnermodel.TaskReadout, len(in.Modules))
	for i, m := range in.Modules {
		p[i] = make([]learnermodel.TaskReadout, len(m.Pool))
	}
	for _, at := range tasksByID(in) {
		task := in.Modules[at.module].Pool[at.item].Task
		tr, err := in.Readout.Task(task)
		if err != nil {
			return nil, &detailError{detail: "task " + shown(task.ID) + " of the pool", cause: err}
		}
		if err := checkProbabilities(task.ID, tr); err != nil {
			return nil, err
		}
		p[at.module][at.item] = tr
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

// firstDuplicate gives the least value that comes twice: an ID byte by byte, a number by value.
func firstDuplicate[T cmp.Ordered](s []T) (T, bool) {
	c := sortedCopy(s)
	for i := 1; i < len(c); i++ {
		if c[i] == c[i-1] {
			return c[i], true
		}
	}
	var zero T
	return zero, false
}

// at is a task of a pool or a lesson by its places: of its module in Input.Modules and of it in
// the module's list.
type at struct{ module, item int }

// tasksByID gives the tasks of the pools of all the modules in the order of their full IDs, byte
// by byte, and lessonsByID the lessons: one module's ID may begin another's, so the order of the
// modules is not the order of their tasks.
func tasksByID(in Input) []at {
	return acrossModules(in, func(m Module) int { return len(m.Pool) }, func(m Module, j int) string { return m.Pool[j].Task.ID })
}

func lessonsByID(in Input) []at {
	return acrossModules(in, func(m Module) int { return len(m.Lessons) }, func(m Module, j int) string { return m.Lessons[j].ID })
}

func acrossModules(in Input, n func(Module) int, id func(Module, int) string) []at {
	var out []at
	for i, m := range in.Modules {
		for j := range n(m) {
			out = append(out, at{i, j})
		}
	}
	slices.SortStableFunc(out, func(a, b at) int {
		return strings.Compare(id(in.Modules[a.module], a.item), id(in.Modules[b.module], b.item))
	})
	return out
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

func sortedCopy[T cmp.Ordered](s []T) []T {
	c := slices.Clone(s)
	slices.Sort(c)
	return c
}
