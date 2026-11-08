package engine

import (
	"cmp"
	"math"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The states of a skill in one module, for the skills with a verdict (docs/architecture.md,
// section 3.1; R03): open, within reach, on pause, not reached yet and mastered for the queue,
// and the candidates the task on the skill is chosen among.

// The constants of the rule, as in standalone. Only a run of sim changes them (R23 of spec).
const (
	corridorLow  = 0.70  // the corridor of success, from; and the least P of a skill within reach
	corridorHigh = 0.85  // the corridor of success, to
	corridorAim  = 0.775 // the aim inside the corridor

	solutionPause int64 = 86_400 // a task whose solution went to the device is no candidate for 24 hours
	frequentTrap  int64 = 2      // a trap that this many pieces of evidence of a skill found is frequent
)

// skillView is a skill with a verdict of a module of the engine, as the rule reads it in one pass.
type skillView struct {
	id       string
	read     learnermodel.SkillReadout // zero when the skill has no state
	evidence bool                      // the skill has evidence of its own, where it is primary
	open     bool
	// The candidates of the skill in the pass, in the order of the task on the skill.
	candidates []candidate
	reachable  bool
	// The recommended task: the task on the skill without the pause before a task with a seen
	// solution (docs/architecture.md, section 3.3).
	recommended    candidate
	hasRecommended bool
	// Mastered by the model while the recommended task is not above the band it is mastered at.
	masteredForQueue bool
}

// paused is a skill open, with evidence, and out of reach.
func (s skillView) paused() bool { return s.open && s.evidence && !s.reachable }

// notYet is a skill open, without evidence, and out of reach: it comes when the learner grows.
func (s skillView) notYet() bool { return s.open && !s.evidence && !s.reachable }

// skills reads the skills with a verdict of the module at place i of the input, in the order of
// the graph. In the second pass a task with a solution that went to the device less than 24
// hours ago is a candidate again, after the others, but the task of the last closed attempt is
// not (the amendment of 2026-11-07 to R03).
func (pr prepared) skills(i int, second bool) []skillView {
	in := pr.in
	m := in.Modules[i]
	var mr learnermodel.ModuleReadout
	if j := slices.IndexFunc(in.Readout.Modules, func(r learnermodel.ModuleReadout) bool { return r.Module == m.Graph.Module }); j >= 0 {
		mr = in.Readout.Modules[j]
	}
	reads := sortedBy(mr.Skills, func(r learnermodel.SkillReadout) string { return r.Skill })
	graph := sortedBy(m.Graph.Skills, func(s learnermodel.Skill) string { return s.ID })
	readOf := func(id string) learnermodel.SkillReadout {
		r, _ := lookup(reads, func(r learnermodel.SkillReadout) string { return r.Skill }, id)
		return r
	}

	// A skill is open with evidence of its own, whatever came of its prerequisites (the first
	// amendment of 2026-11-07 to R03), or when each of its prerequisites is counted.
	var open []string
	for _, s := range graph {
		if readOf(s.ID).EvidenceCount > 0 || allCounted(s.Prerequisites, graph, readOf) {
			open = append(open, s.ID)
		}
	}
	isOpen := func(id string) bool {
		_, found := slices.BinarySearch(open, id)
		return found
	}

	pool := make([]poolTask, len(m.Pool))
	for j, c := range m.Pool {
		pool[j] = newPoolTask(c, pr.p[i][j])
	}
	solutions := sortedBy(in.Summary.Tasks, func(t TaskSolution) string { return t.Task })
	issues := sortedBy(in.Summary.Skills, func(s SkillIssue) string { return s.Skill })
	var out []skillView
	for _, s := range graphOrder(m.Graph.Skills) {
		if s.Assessability != learnermodel.AssessabilityVerdict {
			continue
		}
		v := skillView{id: s.ID, read: readOf(s.ID), open: isOpen(s.ID)}
		v.evidence = v.read.EvidenceCount > 0
		frequent := frequentTraps(v.read.Traps)
		var all []candidate
		for _, t := range pool {
			cd, ok := candidateOf(t, s.ID, frequent)
			if !ok || !othersOpen(t.primaries, s.ID, isOpen) {
				continue
			}
			all = append(all, cd)
			if !t.c.Task.Template {
				if sent, ok := lookup(solutions, func(t TaskSolution) string { return t.Task }, cd.id); ok {
					cd.seen, cd.seenAt = true, sent.SolutionSentAt
				}
			}
			if cd.seen && pausedSince(in.Time, cd.seenAt) && (!second || in.HasLast && in.Last.Task == cd.id) {
				continue
			}
			v.candidates = append(v.candidates, cd)
		}
		slices.SortFunc(v.candidates, byCorridor)
		if len(all) > 0 {
			v.recommended, v.hasRecommended = slices.MinFunc(all, byCorridor), true
		}
		withinReach := slices.ContainsFunc(v.candidates, func(c candidate) bool { return c.p >= corridorLow })
		issue, _ := lookup(issues, func(s SkillIssue) string { return s.Skill }, s.ID)
		v.reachable = v.open && len(v.candidates) > 0 && (withinReach || backFromPause(issue))
		v.masteredForQueue = v.read.Mastered && (!v.hasRecommended || v.recommended.band <= v.read.MasteryBand.Place)
		out = append(out, v)
	}
	return out
}

// allCounted tells whether each prerequisite is counted: a skill of the module with a verdict,
// mastered by the model at any band (R03). A prerequisite with a rubric or none and an anchor
// skill hold the skill until the rules for them come (RUN.md, T18 and T19).
func allCounted(prerequisites []string, graph []learnermodel.Skill, readOf func(string) learnermodel.SkillReadout) bool {
	for _, p := range prerequisites {
		s, ok := lookup(graph, func(s learnermodel.Skill) string { return s.ID }, p)
		if !ok || s.Assessability != learnermodel.AssessabilityVerdict || !readOf(p).Mastered {
			return false
		}
	}
	return true
}

// graphOrder orders the skills of a graph by tier, then by ID byte by byte: the order of the
// skills in a graph carries no meaning (the format's text, section 12).
func graphOrder(skills []learnermodel.Skill) []learnermodel.Skill {
	s := slices.Clone(skills)
	slices.SortFunc(s, func(a, b learnermodel.Skill) int {
		return cmp.Or(cmp.Compare(a.Tier, b.Tier), strings.Compare(a.ID, b.ID))
	})
	return s
}

// primaryOf gives the primary skills of an item: its own, or else the task's.
func primaryOf(t learnermodel.Task, it learnermodel.TaskItem) []string {
	if it.OwnSkills {
		return it.Primary
	}
	return t.Primary
}

// primaries gives the primary skills of a task, as learner-model counts them: the items' own,
// and the task's when an item takes them or the task has no items.
func primaries(t learnermodel.Task) []string {
	var out []string
	takesTask := len(t.Items) == 0
	for _, it := range t.Items {
		if it.OwnSkills {
			out = append(out, it.Primary...)
		} else {
			takesTask = true
		}
	}
	if takesTask {
		out = append(out, t.Primary...)
	}
	return out
}

// othersOpen tells whether every primary skill of a task other than skill is open: otherwise
// the task would give evidence to a closed skill past its prerequisites (the amendment of
// 2026-11-07 to R03 found by the review of T03).
func othersOpen(primaries []string, skill string, isOpen func(string) bool) bool {
	for _, p := range primaries {
		if p != skill && !isOpen(p) {
			return false
		}
	}
	return true
}

// pausedSince tells whether a solution that went to the device at sent did so less than 24 hours
// before now, without overflow: a time of sending after now is less than 24 hours ago too.
func pausedSince(now, sent int64) bool {
	if now < math.MinInt64+solutionPause {
		return true
	}
	return sent > now-solutionPause
}

// failed tells whether an attempt is a failure on a skill, as the model counts one: fewer than
// half of its first verdicts right, or a solution shown without a right first verdict (R06 of
// learner-model).
func failed(r SkillResult) bool {
	return 2*int64(r.FirstRight) < int64(r.Answered) || r.SolutionShown && r.FirstRight == 0
}

// right tells whether an attempt is right on a skill: one that is not a failure and has a first
// verdict on it. An attempt with none is neither (the amendment of 2026-11-08 to R03).
func right(r SkillResult) bool { return r.Answered > 0 && !failed(r) }

// backFromPause tells whether the last closed attempt of a skill is a probe or a repetition, and
// right: it puts the skill within reach until the next closed attempt (R06).
func backFromPause(s SkillIssue) bool {
	return s.HasLast && (s.LastCode == ReasonProbe || s.LastCode == ReasonRepetition) && right(s.LastResult)
}

// sortedBy gives a copy of s sorted by ID byte by byte, stable.
func sortedBy[T any](s []T, id func(T) string) []T {
	c := slices.Clone(s)
	slices.SortStableFunc(c, func(a, b T) int { return strings.Compare(id(a), id(b)) })
	return c
}

// lookup finds the value with an ID in s, sorted by sortedBy.
func lookup[T any](s []T, id func(T) string, key string) (T, bool) {
	i, found := slices.BinarySearchFunc(s, key, func(v T, k string) int { return strings.Compare(id(v), k) })
	if !found {
		var zero T
		return zero, false
	}
	return s[i], true
}
