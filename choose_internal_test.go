package engine

import (
	"errors"
	"math"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// A P outside [0, 1] is an error, NaN and the infinities too, and an item with no verdict has no
// P to check (R10, rule 3).
func TestCheckProbabilities(t *testing.T) {
	nan := math.Float64frombits(0x7ff8000000000001) // a quiet NaN, by its bits: math.NaN is off the list of rule 2
	for _, tc := range []struct {
		name    string
		p       float64
		verdict bool
		ok      bool
	}{
		{"zero", 0, true, true},
		{"negative zero", math.Copysign(0, -1), true, true},
		{"one", 1, true, true},
		{"inside", 0.775, true, true},
		{"just above one", math.Float64frombits(math.Float64bits(1) + 1), true, false},
		{"just below zero", -math.Float64frombits(1), true, false},
		{"NaN", nan, true, false},
		{"plus infinity", math.Inf(1), true, false},
		{"minus infinity", math.Inf(-1), true, false},
		{"NaN with no verdict", nan, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tr := learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{
				{Place: 1, Probability: 0.5, HasVerdict: true},
				{Place: 2, Probability: tc.p, HasVerdict: tc.verdict},
			}}
			err := checkProbabilities("example/m/task/t", tr)
			if tc.ok && err != nil {
				t.Errorf("checkProbabilities: %v, want none", err)
			}
			if !tc.ok && (!errors.Is(err, ErrInvalid) || err.Error() != `engine: invalid input: item 2 of task "example/m/task/t" has a P outside [0, 1]`) {
				t.Errorf("checkProbabilities: %v, want ErrInvalid naming item 2", err)
			}
		})
	}
}

// The P of each candidate is kept by the input's places, for the rule.
func TestPrepareKeepsProbabilities(t *testing.T) {
	task := learnermodel.Task{
		ID: "example/m/task/t", Band: learnermodel.BandRef{Place: 1, ID: "b"}, Difficulty: 2,
		Primary: []string{"example/m/skill/s"},
		Items: []learnermodel.TaskItem{
			{Place: 2, Slots: []learnermodel.Slot{{Check: learnermodel.CheckKindNumber}}},
			{Place: 1},
		},
	}
	g, err := learnermodel.NewGraph("example/m", []learnermodel.Band{{Place: 1, ID: "b"}},
		[]learnermodel.Skill{{ID: "example/m/skill/s", Assessability: learnermodel.AssessabilityVerdict}}, nil,
		[]learnermodel.Task{task})
	if err != nil {
		t.Fatal(err)
	}
	r, err := learnermodel.Read(learnermodel.State{},
		learnermodel.Graphs{"example/m": {Graph: g, Current: true}}, 0, learnermodel.Grade{}, []string{"example/m"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := prepare(Input{Readout: r, Modules: []Module{{Graph: g, Pool: []Candidate{{Task: task}}}}})
	if err != nil {
		t.Fatal(err)
	}
	items := p.pools[0][0].p
	if len(items) != 2 || items[0].Place != 1 || items[0].HasVerdict || items[1].Place != 2 || !items[1].HasVerdict {
		t.Fatalf("P of the candidate: %+v, want item 1 with no verdict and item 2 with one", items)
	}
	if want, _ := r.Task(task); items[1].Probability != want.Items[1].Probability {
		t.Errorf("P of item 2: %v, want Readout.Task's %v", items[1].Probability, want.Items[1].Probability)
	}
}

// An ID is shown byte by byte: printable ASCII as it is, a quote and a backslash escaped, and any
// other byte in hex.
func TestShown(t *testing.T) {
	for in, want := range map[string]string{
		"":                `""`,
		"example/m":       `"example/m"`,
		`a"b\c`:           `"a\"b\\c"`,
		"tab\there":       `"tab\x09here"`,
		"é":               `"\xc3\xa9"`,
		"\x7f\x00":        `"\x7f\x00"`,
		"space ~ tilde 1": `"space ~ tilde 1"`,
	} {
		if got := shown(in); got != want {
			t.Errorf("shown(%q) = %s, want %s", in, got, want)
		}
	}
}

func TestErrorWithCause(t *testing.T) {
	cause := errors.New("learner-model said no")
	err := &detailError{detail: "task \"x\" of the pool", cause: cause}
	if !errors.Is(err, ErrInvalid) || !errors.Is(err, cause) {
		t.Errorf("%v does not wrap both its kind and its cause", err)
	}
	if want := `engine: invalid input: task "x" of the pool: learner-model said no`; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

// The traps of a task are checked pair by pair when it has few, and by sorted places when it has
// many: both ways name the same least place, whether the traps of an item come twice or stand where
// the task has no item.
func TestCheckTrapsFewAndMany(t *testing.T) {
	const id = "example/m/task/t"
	for _, n := range []int32{3, trapsByPair + 4} {
		task := learnermodel.Task{ID: id}
		var traps []ItemTraps
		for p := n; p >= 1; p-- { // places from the last, so that the least is found, not the first
			task.Items = append(task.Items, learnermodel.TaskItem{Place: p})
			traps = append(traps, ItemTraps{Place: p})
		}
		for _, tc := range []struct {
			name  string
			extra []ItemTraps
			want  string
		}{
			{"none", nil, ""},
			{"twice", []ItemTraps{{Place: 3}, {Place: 2}}, `item 2 of task "` + id + `" has its traps twice`},
			{"no item", []ItemTraps{{Place: n + 5}, {Place: n + 2}}, `traps of item ` + itoa(int64(n+2)) + ` of task "` + id + `", which has no item there`},
			{"both", []ItemTraps{{Place: n + 2}, {Place: 1}}, `item 1 of task "` + id + `" has its traps twice`},
		} {
			c := Candidate{Task: task, Traps: append(append([]ItemTraps(nil), traps...), tc.extra...)}
			in := Input{Modules: []Module{{Pool: []Candidate{c}}}}
			err := checkTraps(in, tasksByID(in))
			if got := ""; err != nil {
				got = err.Error()
				if !errors.Is(err, ErrInvalid) || got != "engine: invalid input: "+tc.want {
					t.Errorf("%d traps, %s: %q, want %q", len(c.Traps), tc.name, got, tc.want)
				}
			} else if tc.want != "" {
				t.Errorf("%d traps, %s: no error, want %q", len(c.Traps), tc.name, tc.want)
			}
		}
	}
}
