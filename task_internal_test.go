package engine

import (
	"math"
	"slices"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The corridor is from 0.70 to 0.85, both inclusive.
func TestInCorridor(t *testing.T) {
	below := math.Float64frombits(math.Float64bits(corridorLow) - 1)
	above := math.Float64frombits(math.Float64bits(corridorHigh) + 1)
	for _, tc := range []struct {
		p  float64
		in bool
	}{{corridorLow, true}, {corridorHigh, true}, {corridorAim, true}, {below, false}, {above, false}, {0, false}, {1, false}} {
		if got := inCorridor(tc.p); got != tc.in {
			t.Errorf("inCorridor(%v) = %v, want %v", tc.p, got, tc.in)
		}
	}
}

// The order of the task on a skill, case by case: each list is in the order byCorridor gives, and
// so, from any order, is sorting it.
func TestByCorridor(t *testing.T) {
	// 0.775 ± 0.0625 are exact, and as near the aim as each other.
	near, nearLow := corridorAim+0.0625, corridorAim-0.0625
	for _, tc := range []struct {
		name string
		want []candidate
	}{
		{"in the corridor first", []candidate{{id: "z", p: 0.71}, {id: "b", p: 0.69}, {id: "a", p: 0.775 + 0.2}}},
		{"nearest to the aim", []candidate{{id: "z", p: 0.78}, {id: "a", p: 0.84}}},
		{"of two as near, the easier", []candidate{{id: "z", p: near}, {id: "a", p: nearLow}}},
		{"out of the corridor, nearest, then the easier", []candidate{{id: "z", p: 0.9}, {id: "y", p: 0.65}, {id: "a", p: 0.6}}},
		{"then by ID", []candidate{{id: "a", p: 0.8}, {id: "b", p: 0.8}}},
		{"more frequent traps first", []candidate{{id: "z", p: 0.84, traps: []int{1, 2}}, {id: "a", p: 0.775, traps: []int{0}}, {id: "b", p: 0.775}}},
		{"of as many, the higher trap", []candidate{{id: "z", p: 0.84, traps: []int{0, 3}}, {id: "a", p: 0.775, traps: []int{1, 2}}}},
		{"traps count only in the corridor", []candidate{{id: "z", p: 0.69}, {id: "a", p: 0.6, traps: []int{0}}}},
		{"in the corridor before an unseen solution outside it", []candidate{{id: "z", p: 0.775, seen: true, seenAt: 1}, {id: "a", p: 0.1}}},
		{"in the corridor, a seen solution counts not", []candidate{{id: "z", p: 0.775, seen: true, seenAt: 1}, {id: "a", p: 0.8}}},
		{"outside the corridor, a seen solution counts not", []candidate{{id: "z", p: 0.6, seen: true, seenAt: 1}, {id: "a", p: 0.1}}},
		{"of seen solutions, not the longest ago but the nearest", []candidate{{id: "z", p: 0.775, seen: true, seenAt: 2}, {id: "a", p: 0.8, seen: true, seenAt: 1}}},
		{"on pause after the rest", []candidate{{id: "z", p: 0.1}, {id: "a", p: 0.775, seen: true, seenAt: 9, recentSolution: true}}},
		{"the last after them all", []candidate{{id: "z", p: 0.1, seen: true, seenAt: 9, recentSolution: true}, {id: "a", p: 0.775, seen: true, seenAt: 1, recentSolution: true, last: true}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Clone(tc.want)
			slices.Reverse(got)
			slices.SortFunc(got, byCorridor)
			if !slices.EqualFunc(got, tc.want, func(a, b candidate) bool { return a.id == b.id }) {
				t.Errorf("order %v, want %v", got, tc.want)
			}
		})
	}
}

// The probe and the tasks with a rubric keep the order of a seen solution: those never seen first,
// then the one seen longest ago, and the task of the last closed attempt after them all.
func TestBySolution(t *testing.T) {
	want := []candidate{{id: "z"}, {id: "y", seen: true, seenAt: 1}, {id: "x", seen: true, seenAt: 2}, {id: "a", seen: true, last: true}}
	got := slices.Clone(want)
	slices.Reverse(got)
	slices.SortStableFunc(got, bySolution)
	if !slices.EqualFunc(got, want, func(a, b candidate) bool { return a.id == b.id }) {
		t.Errorf("order %v, want %v", got, want)
	}
}

// The task of a probe is the easiest, after those with a seen solution.
func TestByEase(t *testing.T) {
	want := []candidate{{id: "z", p: 0.9}, {id: "a", p: 0.8}, {id: "b", p: 0.8}, {id: "c", p: 0.99, seen: true, seenAt: 5}}
	got := slices.Clone(want)
	slices.Reverse(got)
	slices.SortFunc(got, byEase)
	if !slices.EqualFunc(got, want, func(a, b candidate) bool { return a.id == b.id }) {
		t.Errorf("order %v, want %v", got, want)
	}
	s := skillView{candidates: got}
	slices.SortFunc(s.candidates, byCorridor)
	if c, _ := s.task(true); c.id != "z" {
		t.Errorf("the easiest is %s, want z", c.id)
	}
	if c, _ := s.task(false); c.id != "a" {
		t.Errorf("the task is %s, want a, nearest to the aim", c.id)
	}
	if _, ok := (skillView{}).task(false); ok {
		t.Error("a skill with no candidates has a task")
	}
}

// A trap is frequent when at least two pieces of evidence found it, ordered by their number, then
// by the time of the last, then by ID.
func TestFrequentTraps(t *testing.T) {
	got := frequentTraps([]learnermodel.TrapCount{
		{ID: "once", Count: 1, Last: 100},
		{ID: "b", Count: 2, Last: 7},
		{ID: "a", Count: 2, Last: 7},
		{ID: "late", Count: 2, Last: 9},
		{ID: "most", Count: 3, Last: 1},
	})
	if want := []string{"most", "late", "a", "b"}; !slices.Equal(got, want) {
		t.Errorf("frequentTraps = %v, want %v", got, want)
	}
}

// The aim of the corridor is 0.775 exactly: of 0.772 and 0.781, the first is nearer it, though an
// aim of 0.78 would take the second. Found by a negative probe of T21: no vector told them apart.
func TestCorridorAim(t *testing.T) {
	a, b := candidate{id: "a", p: 0.772}, candidate{id: "b", p: 0.781}
	if byCorridor(a, b) >= 0 || byCorridor(b, a) <= 0 {
		t.Errorf("0.781 comes before 0.772 for the aim %v", corridorAim)
	}
}
