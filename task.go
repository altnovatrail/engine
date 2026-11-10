package engine

import (
	"cmp"
	"math"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The task on a skill (docs/architecture.md, section 3.1; R03): once the rule has chosen a
// skill, which of its candidates it gives.

// candidate is a task of the pool as a candidate of a skill.
type candidate struct {
	id   string
	band int32   // the place of the task's band
	p    float64 // the P of the task on the skill
	// The places of the frequent traps of the skill the task catches, ascending: a place is the
	// trap's in the order of frequentTraps.
	traps  []int
	seen   bool  // the solution went to the device: never for a template
	seenAt int64 // when it last did
	// The solution went less than 24 hours ago: a candidate only in the second pass, after the rest.
	recentSolution bool
	last           bool // the task of the last closed attempt, back on the map in the second pass
}

// poolTask is a task of the pool with what the rule reads of it whatever the skill: its items and
// their P, by place, and its primary skills.
type poolTask struct {
	c         Candidate
	items     []learnermodel.TaskItem
	p         []learnermodel.ItemProbability
	primaries []string
}

func newPoolTask(c Candidate, tr learnermodel.TaskReadout) poolTask {
	t := poolTask{c: c, items: slices.Clone(c.Task.Items), p: slices.Clone(tr.Items), primaries: primaries(c.Task)}
	slices.SortStableFunc(t.items, func(a, b learnermodel.TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	slices.SortStableFunc(t.p, func(a, b learnermodel.ItemProbability) int { return cmp.Compare(a.Place, b.Place) })
	return t
}

// onSkill tells whether a skill is primary in the item at a place.
func (t poolTask) onSkill(place int32, skill string) bool {
	i, found := slices.BinarySearchFunc(t.items, place, func(it learnermodel.TaskItem, p int32) int { return cmp.Compare(it.Place, p) })
	return found && slices.Contains(primaryOf(t.c.Task, t.items[i]), skill)
}

// candidateOf reads a task of the pool as a candidate of a skill, from the P of its items: the
// task on the skill has an item with a verdict where the skill is primary, or it is no candidate
// of the skill (the amendment of 2026-11-08 to R03).
func candidateOf(t poolTask, skill string, frequent []string) (candidate, bool) {
	// The P of the task on the skill is the mean P of its items where the skill is primary,
	// summed by place: addition of floating-point numbers is not associative (R10).
	var sum float64
	var n int
	for _, ip := range t.p {
		if ip.HasVerdict && t.onSkill(ip.Place, skill) {
			sum += ip.Probability
			n++
		}
	}
	if n == 0 {
		return candidate{}, false
	}

	var caught []string
	for _, it := range t.c.Traps {
		if t.onSkill(it.Place, skill) {
			caught = append(caught, it.Traps...)
		}
	}
	var traps []int
	for place, id := range frequent {
		if slices.Contains(caught, id) {
			traps = append(traps, place)
		}
	}
	return candidate{id: t.c.Task.ID, band: t.c.Task.Band.Place, p: float64(sum / float64(n)), traps: traps}, true
}

// frequentTraps gives the frequent traps of a skill — those at least two pieces of its evidence
// found — by their number, then by the time of the last, then by ID.
func frequentTraps(traps []learnermodel.TrapCount) []string {
	var f []learnermodel.TrapCount
	for _, t := range traps {
		if t.Count >= frequentTrap {
			f = append(f, t)
		}
	}
	slices.SortFunc(f, func(a, b learnermodel.TrapCount) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), cmp.Compare(b.Last, a.Last), strings.Compare(a.ID, b.ID))
	})
	ids := make([]string, 0, len(f))
	for _, t := range f {
		ids = append(ids, t.ID)
	}
	return ids
}

// bySolution puts the task of the last closed attempt, back on the map, after them all, and orders
// the rest by bySeen.
func bySolution(a, b candidate) int {
	return cmp.Or(falseFirst(a.last, b.last), bySeen(a, b))
}

// bySeen puts the candidates whose solution never went to the device first, and of the rest the
// one whose solution went longest ago (R82 of spec).
func bySeen(a, b candidate) int {
	if c := falseFirst(a.seen, b.seen); c != 0 || !a.seen {
		return c
	}
	return cmp.Compare(a.seenAt, b.seenAt)
}

// falseFirst orders false before true.
func falseFirst(a, b bool) int {
	switch {
	case a == b:
		return 0
	case !a:
		return -1
	}
	return 1
}

func inCorridor(p float64) bool { return p >= corridorLow && p <= corridorHigh }

// byCorridor orders the candidates of a skill for the task on it, the first being the task: the
// task of the last closed attempt, back on the map, after them all, and the tasks on pause, back in
// the second pass, after the rest; then those with a P in the corridor first, so that a task in the
// corridor whose solution was seen comes before one outside it whose solution was not; within
// each, bySeen; then, in the corridor, the one that catches more frequent traps, then the one
// that catches the higher; then the one nearest by P to the aim, of two as near the easier, with
// the greater P; then by ID byte by byte.
func byCorridor(a, b candidate) int {
	if c := cmp.Or(falseFirst(a.last, b.last), falseFirst(a.recentSolution, b.recentSolution)); c != 0 {
		return c
	}
	inA, inB := inCorridor(a.p), inCorridor(b.p)
	if inA != inB {
		if inA {
			return -1
		}
		return 1
	}
	if c := bySeen(a, b); c != 0 {
		return c
	}
	if inA {
		if c := cmp.Or(cmp.Compare(len(b.traps), len(a.traps)), slices.Compare(a.traps, b.traps)); c != 0 {
			return c
		}
	}
	return cmp.Or(
		cmp.Compare(math.Abs(a.p-corridorAim), math.Abs(b.p-corridorAim)),
		cmp.Compare(b.p, a.p),
		strings.Compare(a.id, b.id),
	)
}

// byEase orders the candidates of a skill for a probe, the first being the easiest: after
// bySolution, the greatest P first, then by ID byte by byte.
func byEase(a, b candidate) int {
	return cmp.Or(bySolution(a, b), cmp.Compare(b.p, a.p), strings.Compare(a.id, b.id))
}

// task gives the task on a skill among its candidates in the pass: the first by byCorridor, or
// by byEase for a probe.
func (s skillView) task(easiest bool) (candidate, bool) {
	if len(s.candidates) == 0 {
		return candidate{}, false
	}
	if easiest {
		return slices.MinFunc(s.candidates, byEase), true
	}
	return s.candidates[0], true
}
