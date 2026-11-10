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
	// Whether the skill was handed out, by a task on it or a lesson that teaches it, and when last
	// (the summary).
	issued     bool
	lastIssued int64
	// The lowest band the skill is taught at, by place, when the graph has one: the bottom of the
	// ladder of branch 9.
	lowestBand    int32
	hasLowestBand bool
	// Not handed out, and taught by a lesson of the author's order not passed: it comes by that
	// lesson, and branches 6–9 do not take it (the amendment of 2026-11-08 to R05).
	awaitsLesson bool
	// Due for review: open, mastered for the queue, with candidates in the pass, and the time of the
	// choice not before its review date, or the date marked already (docs/architecture.md, section
	// 3.4; R06). A skill that awaits its lesson is not.
	due bool
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
	mr := moduleReadout(in, m.Graph.Module)
	reads := sortedBy(mr.Skills, func(r learnermodel.SkillReadout) string { return r.Skill })
	readOf := func(id string) learnermodel.SkillReadout {
		r, _ := lookup(reads, func(r learnermodel.SkillReadout) string { return r.Skill }, id)
		return r
	}

	pool := pr.poolOf(i)
	issues := sortedBy(in.Summary.Skills, func(s SkillIssue) string { return s.Skill })
	var out []skillView
	for _, s := range graphOrder(m.Graph.Skills) {
		if s.Assessability != learnermodel.AssessabilityVerdict {
			continue
		}
		v := skillView{id: s.ID, read: readOf(s.ID), open: pr.isOpen(s.ID)}
		v.evidence = v.read.EvidenceCount > 0
		frequent := frequentTraps(v.read.Traps)
		cands, all := pr.candidates(pool, s.ID, frequent, second, false)
		v.candidates = cands
		if len(all) > 0 {
			v.recommended, v.hasRecommended = slices.MinFunc(all, byCorridor), true
		}
		withinReach := slices.ContainsFunc(v.candidates, func(c candidate) bool { return c.p >= corridorLow })
		issue, _ := lookup(issues, func(s SkillIssue) string { return s.Skill }, s.ID)
		v.issued, v.lastIssued = issue.Issued, issue.LastIssued
		if len(s.BandsTaught) > 0 {
			v.lowestBand, v.hasLowestBand = slices.Min(s.BandsTaught), true
		}
		v.reachable = v.open && len(v.candidates) > 0 && (withinReach || backFromPause(issue))
		v.masteredForQueue = v.read.Mastered && (!v.hasRecommended || v.recommended.band <= v.read.MasteryBand.Place)
		v.awaitsLesson = !v.issued && pr.awaitsLesson(i, s.ID)
		v.due = v.open && v.masteredForQueue && len(v.candidates) > 0 && !v.awaitsLesson && (v.read.Review.Already || in.Time >= v.read.Review.At)
		out = append(out, v)
	}
	return out
}

// candidates gives the candidates of a skill of the module at place i in the pass, in the order
// of the task on the skill, and all of its tasks without the pause before a seen solution, for the
// recommended task. A task with another primary skill closed is no candidate, but on the map: the
// skill chosen there takes all its tasks (R04). In the second pass a task on pause is a candidate
// again, after the others, but the task of the last closed attempt is not (the amendment of
// 2026-11-07 to R03) — on the map it is, after all the others (the amendment of 2026-11-07 to
// R04).
func (pr prepared) candidates(pool []poolTask, skill string, frequent []string, second, onMap bool) (cands, all []candidate) {
	for _, t := range pool {
		cd, ok := candidateOf(t, skill, frequent)
		if !ok || !onMap && !othersOpen(t.primaries, skill, pr.isOpen) {
			continue
		}
		all = append(all, cd)
		sp, ok := pr.pause(t.c.Task, second, onMap)
		if !ok {
			continue
		}
		cd.seen, cd.seenAt, cd.recentSolution, cd.last = sp.seen, sp.seenAt, sp.recentSolution, sp.last
		cands = append(cands, cd)
	}
	slices.SortFunc(cands, byCorridor)
	return cands, all
}

// pause reads the pause before a task with a seen solution in the pass: whether the task is a
// candidate in it, and, as a candidate with its ID alone, what bySolution and byCorridor order it
// by — seen, seenAt, recentSolution and last. A task whose solution went to the device less than
// 24 hours ago is none in the first pass; in the second it is, after the rest, but the task of the
// last closed attempt is not — on the map it is, after all the others. The task of the last closed
// attempt whose solution went to the device is on pause however long ago that was: past the 24
// hours it would come first by the corridor, and R82 of spec keeps it from being the next task.
func (pr prepared) pause(t learnermodel.Task, second, onMap bool) (candidate, bool) {
	in := pr.in
	cd := candidate{id: t.ID, band: t.Band.Place}
	if !t.Template {
		if sent, ok := lookup(pr.solutions, func(t TaskSolution) string { return t.Task }, t.ID); ok {
			cd.seen, cd.seenAt = true, sent.SolutionSentAt
		}
	}
	last := in.HasLast && in.Last.Task == cd.id
	if cd.seen && (pausedSince(in.Time, cd.seenAt) || last) {
		if !second || last && !onMap {
			return cd, false
		}
		cd.recentSolution, cd.last = true, last
	}
	return cd, true
}

// poolOf reads the tasks of the pool of the module at place i.
func (pr prepared) poolOf(i int) []poolTask {
	pool := make([]poolTask, len(pr.in.Modules[i].Pool))
	for j, c := range pr.in.Modules[i].Pool {
		pool[j] = newPoolTask(c, pr.p[i][j])
	}
	return pool
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

// moduleReadout gives the learner model's reading of a module, which checkReadout has made sure
// of for a module of the engine.
func moduleReadout(in Input, module string) learnermodel.ModuleReadout {
	if j := slices.IndexFunc(in.Readout.Modules, func(r learnermodel.ModuleReadout) bool { return r.Module == module }); j >= 0 {
		return in.Readout.Modules[j]
	}
	return learnermodel.ModuleReadout{}
}
