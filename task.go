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
	last   bool  // the task of the last closed attempt, back on the map in the second pass
}

// poolTask is a task of the pool with what the rule reads of it whatever the skill: its items and
// their P, by place, and its primary skills.
type poolTask struct {
	c         *Candidate // the input's, which the rule only reads
	items     []learnermodel.TaskItem
	p         []learnermodel.ItemProbability
	primaries []string
}

// newPoolTask reads a task of the pool with its P. The items and their P are read by place: they
// come in that order as a rule, and are copied to be put in order only when they do not — the rule
// reads them and never writes them.
func newPoolTask(c *Candidate, tr learnermodel.TaskReadout) poolTask {
	byPlace := func(a, b learnermodel.TaskItem) int { return cmp.Compare(a.Place, b.Place) }
	pByPlace := func(a, b learnermodel.ItemProbability) int { return cmp.Compare(a.Place, b.Place) }
	t := poolTask{c: c, items: c.Task.Items, p: tr.Items, primaries: primaries(c.Task)}
	if !slices.IsSortedFunc(t.items, byPlace) {
		t.items = slices.Clone(t.items)
		slices.SortStableFunc(t.items, byPlace)
	}
	if !slices.IsSortedFunc(t.p, pByPlace) {
		t.p = slices.Clone(t.p)
		slices.SortStableFunc(t.p, pByPlace)
	}
	return t
}

// skillTasks is a skill primary in tasks of a pool, and the places of those tasks in the pool, in
// the order of the tasks' IDs. Every candidate of the skill is among them: a skill primary in an
// item of a task is among the task's primary skills.
type skillTasks struct {
	skill string
	tasks []int
}

// bySkill gives the skills primary in the tasks of a pool, by ID byte by byte, each with the places
// of its tasks: the rule reads the candidates of a skill there, and not through the whole pool for
// every skill. order is the places of the pool's tasks in the order of their IDs, and the tasks of a
// skill keep it, so that the order of the pool in the input changes nothing.
func bySkill(pool []poolTask, order []int) []skillTasks {
	// The skills, sorted: a pool has many tasks on few skills, so each new one is put in its place.
	var skills []string
	for _, t := range pool {
		for _, s := range t.primaries {
			if g, found := slices.BinarySearch(skills, s); !found {
				skills = slices.Insert(skills, g, s)
			}
		}
	}
	// each calls f with the group of every skill primary in the task at place j, once a skill: a
	// skill primary in two of its items is one of its primaries twice.
	each := func(j int, f func(g int)) {
		ps := pool[j].primaries
		for k, s := range ps {
			if !slices.Contains(ps[:k], s) {
				g, _ := slices.BinarySearch(skills, s)
				f(g)
			}
		}
	}
	counts := make([]int, len(skills))
	total := 0
	for j := range pool {
		each(j, func(g int) { counts[g]++; total++ })
	}
	places := make([]int, total)
	out := make([]skillTasks, len(skills))
	start := 0
	for g, s := range skills {
		out[g] = skillTasks{skill: s, tasks: places[start:start:(start + counts[g])]}
		start += counts[g]
	}
	for _, j := range order {
		each(j, func(g int) { out[g].tasks = append(out[g].tasks, j) })
	}
	return out
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

	cd := candidate{id: t.c.Task.ID, band: t.c.Task.Band.Place, p: float64(sum / float64(n))}
	if len(frequent) == 0 { // no frequent trap to catch
		return cd, true
	}
	var caught []string
	for _, it := range t.c.Traps {
		if t.onSkill(it.Place, skill) {
			caught = append(caught, it.Traps...)
		}
	}
	for place, id := range frequent {
		if slices.Contains(caught, id) {
			cd.traps = append(cd.traps, place)
		}
	}
	return cd, true
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

// bySolution puts the candidates whose solution never went to the device first, and of the rest
// the one whose solution went longest ago (R82 of spec); the task of the last closed attempt, back
// on the map, comes after them all.
func bySolution(a, b candidate) int {
	if c := cmp.Or(falseFirst(a.last, b.last), falseFirst(a.seen, b.seen)); c != 0 || !a.seen {
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

// byCorridor orders the candidates of a skill for the task on it, the first being the task:
// after bySolution, those with a P in the corridor first, and among them the one that catches
// more frequent traps, then the one that catches the higher; then the one nearest by P to the aim,
// of two as near the easier, with the greater P; then by ID byte by byte.
func byCorridor(a, b candidate) int {
	if c := bySolution(a, b); c != 0 {
		return c
	}
	inA, inB := inCorridor(a.p), inCorridor(b.p)
	if inA != inB {
		if inA {
			return -1
		}
		return 1
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
	return slices.MinFunc(s.candidates, byCorridor), true
}
