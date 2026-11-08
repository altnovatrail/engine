package engine

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Several modules, the skill chosen on the map and the anchor skills, on the example of
// docs/architecture.md, section 3.2, and at their edges. The reading of the learner model and the
// P of the candidates are set by hand, as in skills_internal_test.go.

// many is an input over several modules, with the P of their candidates beside it.
type many struct {
	t  *testing.T
	in Input
	p  [][]learnermodel.TaskReadout
}

func newMany(t *testing.T, modules ...Module) *many {
	w := &many{t: t, in: Input{Time: exampleTime, Modules: modules}}
	for _, m := range modules {
		w.in.Readout.Modules = append(w.in.Readout.Modules, learnermodel.ModuleReadout{Module: m.Graph.Module})
		ps := make([]learnermodel.TaskReadout, len(m.Pool))
		for j := range ps {
			ps[j] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1, Probability: 0.8, HasVerdict: true}}}
		}
		w.p = append(w.p, ps)
	}
	return w
}

func (w *many) probabilities(ps []px) {
	for i, m := range w.in.Modules {
		for j, c := range m.Pool {
			name := c.Task.ID[strings.LastIndex(c.Task.ID, "/")+1:]
			for _, x := range ps {
				if x.name == name || strings.HasSuffix(x.name, "*") && strings.HasPrefix(name, strings.TrimSuffix(x.name, "*")) {
					w.p[i][j] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1, Probability: x.p, HasVerdict: true}}}
				}
			}
		}
	}
}

// moduleOf gives the place of the module whose graph has a skill.
func (w *many) moduleOf(skill string) int {
	w.t.Helper()
	i := slices.IndexFunc(w.in.Modules, func(m Module) bool {
		return slices.ContainsFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == skill })
	})
	if i < 0 {
		w.t.Fatalf("no module has skill %s", skill)
	}
	return i
}

func (w *many) readout(module string) *learnermodel.ModuleReadout {
	i := slices.IndexFunc(w.in.Readout.Modules, func(r learnermodel.ModuleReadout) bool { return r.Module == module })
	return &w.in.Readout.Modules[i]
}

func (w *many) skill(r learnermodel.SkillReadout) {
	mr := w.readout(w.in.Modules[w.moduleOf(r.Skill)].Graph.Module)
	mr.Skills = slices.DeleteFunc(mr.Skills, func(s learnermodel.SkillReadout) bool { return s.Skill == r.Skill })
	mr.Skills = append(mr.Skills, r)
	slices.SortFunc(mr.Skills, func(a, b learnermodel.SkillReadout) int { return strings.Compare(a.Skill, b.Skill) })
}

func (w *many) mastered(skill string) {
	w.skill(learnermodel.SkillReadout{Skill: skill, EvidenceCount: 10, Mastered: true, Review: notDue, MasteryBand: learnermodel.BandRef{Place: 1}})
}

func (w *many) anchor(id string, status learnermodel.AnchorStatus, reason learnermodel.AnchorReason) {
	w.in.Readout.Anchors = slices.DeleteFunc(w.in.Readout.Anchors, func(a learnermodel.AnchorReadout) bool { return a.Anchor == id })
	w.in.Readout.Anchors = append(w.in.Readout.Anchors, learnermodel.AnchorReadout{Anchor: id, Status: status, Reason: reason})
}

func (w *many) issue(skill string, at int64) {
	w.in.Summary.Skills = slices.DeleteFunc(w.in.Summary.Skills, func(s SkillIssue) bool { return s.Skill == skill })
	w.in.Summary.Skills = append(w.in.Summary.Skills, SkillIssue{Skill: skill, Issued: true, LastIssued: at})
}

func (w *many) open(skill string) bool { return ready(w.in, w.p).isOpen(skill) }

// want checks the choice: a task by its full ID and a code, or nothing to give when task is empty.
func (w *many) want(task string, code ReasonCode) {
	w.t.Helper()
	want := Choice{Kind: ChoiceNothing}
	if task != "" {
		want = Choice{Kind: ChoiceTask, ID: task, Reason: code}
	}
	if got := ready(w.in, w.p).choose(); got != want {
		w.t.Errorf("chose %+v, want %+v", got, want)
	}
}

func tid(module, name string) string { return module + "/task/" + name }

const fractionsAnchor = "core.fractions-of-a-whole"

var (
	timesTablesSkill = sid(grade3Math, "times-tables")
	fractionsSkill   = sid(grade3Math, "fractions-of-a-whole")
	tricksSkill      = sid(olympiadMath, "arithmetic-tricks")
	partsSkill       = sid(olympiadMath, "fractions-parts")
	ratioSkill       = sid(olympiadMath, "ratio-sharing")
)

// bothModules builds the two maths of the example of section 3.2: school maths for grade 3, where
// fractions-of-a-whole teaches core.fractions-of-a-whole, and olympiad maths, where
// fractions-parts needs arithmetic-tricks and that anchor skill, and ratio-sharing needs
// fractions-parts. Each skill has tasks of difficulty 1, 3 and 5 at one band.
func bothModules(t *testing.T) (grade, olympiad Module) {
	t.Helper()
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	bands := []learnermodel.Band{{Place: 1, ID: "b"}}
	tasks := func(module string, skills ...string) []learnermodel.Task {
		var out []learnermodel.Task
		for _, s := range skills {
			name := s[strings.LastIndex(s, "/")+1:]
			for _, d := range []int32{1, 3, 5} {
				out = append(out, exampleTask(module, name+"-d"+string(rune('0'+d)), band, d, s))
			}
		}
		return out
	}
	grade = exampleModule(t, grade3Math, bands, []learnermodel.Skill{
		{ID: timesTablesSkill, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: fractionsSkill, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{timesTablesSkill}, Teaches: []string{fractionsAnchor}, Tier: 2},
	}, tasks(grade3Math, timesTablesSkill, fractionsSkill))
	olympiad = exampleModule(t, olympiadMath, bands, []learnermodel.Skill{
		{ID: tricksSkill, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: partsSkill, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{tricksSkill, fractionsAnchor}, Tier: 2},
		{ID: ratioSkill, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{partsSkill}, Tier: 3},
	}, tasks(olympiadMath, tricksSkill, partsSkill, ratioSkill))
	return grade, olympiad
}

func bothWorld(t *testing.T) *many {
	grade, olympiad := bothModules(t)
	w := newMany(t, olympiad, grade)
	w.anchor(fractionsAnchor, learnermodel.AnchorHasNot, learnermodel.AnchorReasonNoEvidence)
	return w
}

// Step 1: the modules take turns — school maths first, never handed out and less by ID —, and
// consolidation stays in the module of the attempt, past its trial series.
func TestModulesTakeTurns(t *testing.T) {
	w := bothWorld(t)
	w.want(tid(grade3Math, "times-tables-d1"), ReasonNewSkill)
	w.issue(timesTablesSkill, exampleTime-100)
	w.want(tid(olympiadMath, "arithmetic-tricks-d1"), ReasonNewSkill)
	w.issue(tricksSkill, exampleTime-50)
	w.want(tid(grade3Math, "times-tables-d1"), ReasonQueue)

	w.in.Last, w.in.HasLast = Last{Task: tid(olympiadMath, "arithmetic-tricks-d3"), Results: []SkillResult{wrong(tricksSkill)}}, true
	w.want(tid(grade3Math, "times-tables-d1"), ReasonQueue) // in the olympiad's trial series
	w.readout(olympiadMath).TrialOver = true
	w.want(tid(olympiadMath, "arithmetic-tricks-d1"), ReasonConsolidation)
	w.in.Last.Left = true
	w.want(tid(grade3Math, "times-tables-d1"), ReasonQueue)
	w.in.Last.Left = false

	// Of two modules handed out at one time, the less by ID; the order of the input's lists
	// changes nothing.
	w.in.HasLast = false
	w.issue(tricksSkill, exampleTime-100)
	want := tid(grade3Math, "times-tables-d1")
	w.want(want, ReasonQueue)
	rng := rand.New(rand.NewPCG(7, 8))
	for range 10 {
		v := &many{t: t, in: w.in, p: slices.Clone(w.p)}
		v.in.Modules = slices.Clone(w.in.Modules)
		rng.Shuffle(len(v.in.Modules), func(i, j int) {
			v.in.Modules[i], v.in.Modules[j] = v.in.Modules[j], v.in.Modules[i]
			v.p[i], v.p[j] = v.p[j], v.p[i]
		})
		v.in.Readout.Modules = shuffled(rng, w.in.Readout.Modules)
		v.in.Summary.Skills = shuffled(rng, w.in.Summary.Skills)
		v.want(want, ReasonQueue)
	}

	// A module whose branches give nothing hands its turn on: every task of times-tables is on
	// pause, and fractions-of-a-whole is closed.
	for _, d := range []string{"d1", "d3", "d5"} {
		w.in.Summary.Tasks = append(w.in.Summary.Tasks, TaskSolution{Task: tid(grade3Math, "times-tables-"+d), SolutionSentAt: exampleTime - 1})
	}
	w.want(tid(olympiadMath, "arithmetic-tricks-d1"), ReasonQueue)
}

// Steps 2–5: the anchor skill the learner has not holds fractions-parts — with no evidence, or
// not mastered —, "has" opens it, and losing it again closes it only for a skill with no evidence
// of its own; with only olympiad maths, the anchor is unknown and holds nothing.
func TestAnchorHolds(t *testing.T) {
	w := bothWorld(t)
	w.mastered(tricksSkill)
	if w.open(partsSkill) {
		t.Errorf("fractions-parts is open with the anchor not had, no evidence")
	}
	w.anchor(fractionsAnchor, learnermodel.AnchorHasNot, learnermodel.AnchorReasonNotMastered)
	if w.open(partsSkill) {
		t.Errorf("fractions-parts is open with the anchor not had, not mastered")
	}
	w.anchor(fractionsAnchor, learnermodel.AnchorHas, learnermodel.AnchorReasonNone)
	if !w.open(partsSkill) {
		t.Errorf("fractions-parts is closed with the anchor had")
	}
	w.issue(timesTablesSkill, exampleTime-10)
	w.issue(tricksSkill, exampleTime-100)
	w.want(tid(olympiadMath, "fractions-parts-d1"), ReasonNewSkill)

	w.anchor(fractionsAnchor, learnermodel.AnchorHasNot, learnermodel.AnchorReasonNotMastered)
	w.skill(learnermodel.SkillReadout{Skill: partsSkill, EvidenceCount: 3})
	if !w.open(partsSkill) || w.open(ratioSkill) {
		t.Errorf("fractions-parts open %v, ratio-sharing open %v; want fractions-parts open with its evidence, ratio-sharing closed", w.open(partsSkill), w.open(ratioSkill))
	}
	w.skill(learnermodel.SkillReadout{Skill: ratioSkill, EvidenceCount: 1})
	if !w.open(ratioSkill) {
		t.Errorf("ratio-sharing with evidence of its own is closed")
	}

	_, olympiad := bothModules(t)
	only := newMany(t, olympiad)
	only.anchor(fractionsAnchor, learnermodel.AnchorUnknown, learnermodel.AnchorReasonNoTeacher)
	if only.open(partsSkill) {
		t.Errorf("fractions-parts is open before arithmetic-tricks is mastered")
	}
	only.mastered(tricksSkill)
	if !only.open(partsSkill) {
		t.Errorf("fractions-parts is closed with the anchor unknown")
	}
}

// An anchor the learner has not holds while a skill with a verdict that teaches it has a task in
// the pool — on pause too, and whether the teacher is open or not — or a lesson; with neither, it
// holds nothing. A task of the teacher's whose other primary skill is closed counts for nothing.
func TestAnchorNothingToTeach(t *testing.T) {
	w := bothWorld(t)
	w.mastered(tricksSkill)
	if w.open(fractionsSkill) || w.open(partsSkill) {
		t.Fatalf("the teacher is closed and holds: fractions-parts open %v", w.open(partsSkill))
	}
	g := w.moduleOf(fractionsSkill)
	for _, d := range []string{"d1", "d3", "d5"} {
		w.in.Summary.Tasks = append(w.in.Summary.Tasks, TaskSolution{Task: tid(grade3Math, "fractions-of-a-whole-"+d), SolutionSentAt: exampleTime - 1})
	}
	if w.open(partsSkill) {
		t.Errorf("fractions-parts is open while its teacher's tasks are only on pause")
	}
	pool, p := w.in.Modules[g].Pool, w.p[g]
	w.in.Modules[g].Pool, w.p[g] = pool[:3], p[:3] // times-tables alone
	if !w.open(partsSkill) {
		t.Errorf("fractions-parts is closed with nothing to teach its anchor")
	}
	w.in.Modules[g].Lessons = []Lesson{{ID: grade3Math + "/lesson/fractions", Skills: []string{fractionsSkill}}}
	if w.open(partsSkill) {
		t.Errorf("fractions-parts is open while a lesson teaches its anchor")
	}
	w.in.Modules[g].Lessons = nil

	// The teacher's one task shares it with times-tables-2, closed until times-tables is mastered.
	other := sid(grade3Math, "times-tables-2")
	shared := exampleTask(grade3Math, "shared", learnermodel.BandRef{Place: 1, ID: "b"}, 3, fractionsSkill, other)
	w.in.Modules[g].Graph.Skills = append(w.in.Modules[g].Graph.Skills,
		learnermodel.Skill{ID: other, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{timesTablesSkill}, Tier: 2})
	slices.SortFunc(w.in.Modules[g].Graph.Skills, func(a, b learnermodel.Skill) int { return strings.Compare(a.ID, b.ID) })
	w.in.Modules[g].Pool = append(pool[:3:3], Candidate{Task: shared})
	w.p[g] = append(p[:3:3], learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1, Probability: 0.8, HasVerdict: true}}})
	if !w.open(partsSkill) {
		t.Errorf("fractions-parts is closed, though its teacher's one task has another primary skill closed")
	}
	w.mastered(timesTablesSkill)
	if w.open(partsSkill) {
		t.Errorf("fractions-parts is open, though its teacher's one task is a candidate")
	}
}

// Step 7: a cycle between modules — a in A needs core.y and teaches core.x, b in B needs core.x and
// teaches core.y. Neither anchor holds a or b, which lie in one component, but core.x still holds
// c, which needs it and is in no cycle.
func TestAnchorCycle(t *testing.T) {
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	bands := []learnermodel.Band{{Place: 1, ID: "b"}}
	a, b, c := sid("example/a", "a"), sid("example/b", "b"), sid("example/b", "c")
	ma := exampleModule(t, "example/a", bands, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{"core.y"}, Teaches: []string{"core.x"}, Tier: 1},
	}, []learnermodel.Task{exampleTask("example/a", "a-1", band, 1, a)})
	mb := exampleModule(t, "example/b", bands, []learnermodel.Skill{
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{"core.x"}, Teaches: []string{"core.y"}, Tier: 1},
		{ID: c, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{"core.x"}, Tier: 1},
	}, []learnermodel.Task{exampleTask("example/b", "b-1", band, 1, b), exampleTask("example/b", "c-1", band, 1, c)})
	w := newMany(t, mb, ma)
	w.anchor("core.x", learnermodel.AnchorHasNot, learnermodel.AnchorReasonNoEvidence)
	w.anchor("core.y", learnermodel.AnchorHasNot, learnermodel.AnchorReasonNoEvidence)
	if !w.open(a) || !w.open(b) || w.open(c) {
		t.Errorf("open: a %v, b %v, c %v; want a and b, not c", w.open(a), w.open(b), w.open(c))
	}
	w.want(tid("example/a", "a-1"), ReasonNewSkill)
}

// Step 6: the skill chosen on the map comes closed, past a mistake and in the trial series, with
// its own code; left, the next choice is "My subjects" with no consolidation.
func TestMapSkill(t *testing.T) {
	w := olympiadWorld(t)
	w.in.MapSkill = pigeonSkill
	pick := func() Choice { return ready(w.in, [][]learnermodel.TaskReadout{w.p}).choose() }
	want := Choice{Kind: ChoiceTask, ID: olympiadMath + "/task/pigeon-12-d1-1", Reason: ReasonSkillChosen}
	if got := pick(); got != want {
		t.Errorf("closed: %+v, want %+v", got, want)
	}
	w.trial(true, false)
	w.last("pigeon-12-d1-1", false, wrong(pigeonSkill))
	w.in.Summary.Tasks = []TaskSolution{{Task: olympiadMath + "/task/pigeon-12-d1-1", SolutionSentAt: exampleTime - 1}}
	want.ID = olympiadMath + "/task/pigeon-12-d3-1"
	if got := pick(); got != want {
		t.Errorf("after a mistake: %+v, want %+v", got, want)
	}

	w.in.MapSkill = ""
	w.in.Last.Left = true
	w.issue(enumSkill, exampleTime-100)
	w.issue(gapsSkill, exampleTime-50)
	w.issue(pigeonSkill, exampleTime-1)
	w.skill(learnermodel.SkillReadout{Skill: pigeonSkill, EvidenceCount: 1})
	if got := pick(); got.Reason == ReasonConsolidation || got.ID != olympiadMath+"/task/enum-34-d1-1" {
		t.Errorf("left: %+v, want the queue in My subjects", got)
	}

	for _, skill := range []string{sid("example/other", "a"), sid(olympiadMath, "no-such-skill")} {
		w.in.MapSkill = skill
		if got := pick(); got.Kind != ChoiceNothing {
			t.Errorf("map skill %s outside the modules of the engine: %+v, want nothing to give", skill, got)
		}
	}
}

// The skill chosen on the map takes its tasks whose other primary skill is closed, and in the
// second pass even the task of the last closed attempt, after all the others; one with a rubric
// gets its task.
func TestMapSkillTasks(t *testing.T) {
	module := "example/map"
	band := learnermodel.BandRef{Place: 1, ID: "b"}
	a, b, essay := sid(module, "a"), sid(module, "b"), sid(module, "essay")
	m := exampleModule(t, module, []learnermodel.Band{{Place: 1, ID: "b"}}, []learnermodel.Skill{
		{ID: a, Assessability: learnermodel.AssessabilityVerdict, Tier: 1},
		{ID: b, Assessability: learnermodel.AssessabilityVerdict, Prerequisites: []string{a}, Tier: 2},
		{ID: essay, Assessability: learnermodel.AssessabilityRubric, Tier: 1},
	}, []learnermodel.Task{exampleTask(module, "ab", band, 1, a, b), exampleTask(module, "a-1", band, 3, a), exampleTask(module, "essay-1", band, 1, essay)})
	w := newMany(t, m)
	w.probabilities([]px{{"ab", 0.78}, {"a-1", 0.9}})
	w.want(tid(module, "a-1"), ReasonNewSkill) // ab is no candidate of a: b is closed
	w.in.MapSkill = a
	w.want(tid(module, "ab"), ReasonSkillChosen)

	w.in.Summary.Tasks = []TaskSolution{{Task: tid(module, "ab"), SolutionSentAt: exampleTime - 10}, {Task: tid(module, "a-1"), SolutionSentAt: exampleTime - 20}}
	w.in.Last, w.in.HasLast = Last{Task: tid(module, "a-1")}, true
	if c, ok := ready(w.in, w.p).pass(false); ok {
		t.Errorf("first pass: %+v, want nothing", c)
	}
	w.want(tid(module, "ab"), ReasonSkillChosen) // a-1, the last, comes after it though sent earlier
	w.in.Modules[0].Pool, w.p[0] = w.in.Modules[0].Pool[1:], w.p[0][1:]
	w.want(tid(module, "a-1"), ReasonSkillChosen)

	w.in.MapSkill = essay
	w.want(tid(module, "essay-1"), ReasonSkillChosen)
}

func TestComponents(t *testing.T) {
	// 0 → 1 → 2 → 0 is a cycle; 3 → 0 and 4 alone are components of their own.
	got := components([][]int{{1}, {2}, {0}, {0}, nil})
	if got[0] != got[1] || got[1] != got[2] || got[3] == got[0] || got[4] == got[0] || got[3] == got[4] {
		t.Errorf("components = %v", got)
	}
}

// A teacher teaches only by a task the rule could give it — with an item with a verdict where it
// is primary — and only if the model may master it, with a band it is taught at.
func TestAnchorTeacherMayMaster(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(w *many, g int)
	}{
		{"items with no verdict", func(w *many, g int) {
			for j := 3; j < 6; j++ {
				w.p[g][j] = learnermodel.TaskReadout{Items: []learnermodel.ItemProbability{{Place: 1}}}
			}
		}},
		{"no band taught", func(w *many, g int) {
			skills := slices.Clone(w.in.Modules[g].Graph.Skills)
			k := slices.IndexFunc(skills, func(s learnermodel.Skill) bool { return s.ID == fractionsSkill })
			skills[k].BandsTaught = nil
			w.in.Modules[g].Graph.Skills = skills
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := bothWorld(t)
			w.mastered(tricksSkill)
			w.mastered(timesTablesSkill)
			if w.open(partsSkill) {
				t.Fatalf("fractions-parts is open while its teacher has tasks")
			}
			tc.change(w, w.moduleOf(fractionsSkill))
			if !w.open(partsSkill) {
				t.Errorf("fractions-parts is closed, though its teacher can never be mastered")
			}
		})
	}
}

// A skill a newer version of the module dropped still counts in when the module was last handed
// out: its row in the summary names the module.
func TestModuleTimeDroppedSkill(t *testing.T) {
	w := bothWorld(t)
	w.want(tid(grade3Math, "times-tables-d1"), ReasonNewSkill)
	w.issue(sid(grade3Math, "retired"), exampleTime-1)
	w.want(tid(olympiadMath, "arithmetic-tricks-d1"), ReasonNewSkill)
}
