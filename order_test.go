package engine_test

import (
	"errors"
	"slices"
	"testing"

	learnermodel "github.com/altnovatrail/learner-model"

	"github.com/altnovatrail/engine"
)

const orderModule = "novatrail/olympiad-math"

func orderSkill(name string, tier int32, prerequisites ...string) learnermodel.Skill {
	pre := make([]string, len(prerequisites))
	for k, p := range prerequisites {
		pre[k] = skillID(orderModule, p)
	}
	return learnermodel.Skill{ID: skillID(orderModule, name), Assessability: learnermodel.AssessabilityVerdict, Tier: tier, Prerequisites: pre}
}

func orderLesson(name string, skills ...string) engine.Lesson {
	ids := make([]string, len(skills))
	for k, s := range skills {
		ids[k] = skillID(orderModule, s)
	}
	return engine.Lesson{ID: orderModule + "/lesson/" + name, Skills: ids}
}

// The graph: counting (tier 1) and shapes (tier 1) need nothing; gaps (tier 2) needs counting; and
// patterns (tier 3) needs gaps.
var orderGraph = learnermodel.Graph{Module: orderModule, Skills: []learnermodel.Skill{
	orderSkill("counting", 1),
	orderSkill("gaps", 2, "counting"),
	orderSkill("patterns", 3, "gaps"),
	orderSkill("shapes", 1),
}}

func TestLessonOrderFollowsTheGraphThenTierThenID(t *testing.T) {
	lessons := []engine.Lesson{
		orderLesson("e-patterns", "patterns"),
		orderLesson("d-shapes", "shapes"),
		orderLesson("c-gaps", "gaps"),
		orderLesson("b-counting", "counting"),
		orderLesson("a-counting-again", "counting"),
	}
	got, err := engine.LessonOrder(orderGraph, lessons)
	if err != nil {
		t.Fatal(err)
	}
	// The two lessons of counting and the one of shapes need nothing and share tier 1: by ID, with
	// no order between them. Gaps, of tier 2, comes after counting, so it starts a group; patterns
	// comes after gaps.
	want := []engine.OrderedLesson{
		{ID: orderModule + "/lesson/a-counting-again", Group: 0},
		{ID: orderModule + "/lesson/b-counting", Group: 0},
		{ID: orderModule + "/lesson/d-shapes", Group: 0},
		{ID: orderModule + "/lesson/c-gaps", Group: 1},
		{ID: orderModule + "/lesson/e-patterns", Group: 2},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// A lesson of a lower tier that needs nothing waits behind a lesson the graph lets in earlier,
// and joins its group: the order is the rule's, and the group says only that no prerequisite
// stands between them.
func TestLessonOrderGroupsLessonsWithNoOrderBetweenThem(t *testing.T) {
	g := learnermodel.Graph{Module: orderModule, Skills: []learnermodel.Skill{
		orderSkill("counting", 1),
		orderSkill("gaps", 1, "counting"),
		orderSkill("shapes", 2),
	}}
	got, err := engine.LessonOrder(g, []engine.Lesson{orderLesson("shapes", "shapes"), orderLesson("gaps", "gaps"), orderLesson("counting", "counting")})
	if err != nil {
		t.Fatal(err)
	}
	want := []engine.OrderedLesson{
		{ID: orderModule + "/lesson/counting", Group: 0},
		{ID: orderModule + "/lesson/gaps", Group: 1},
		{ID: orderModule + "/lesson/shapes", Group: 1},
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// The order of the input does not change the order.
func TestLessonOrderReadsNoInputOrder(t *testing.T) {
	a := []engine.Lesson{orderLesson("gaps", "gaps"), orderLesson("counting", "counting"), orderLesson("patterns", "patterns")}
	b := []engine.Lesson{a[2], a[0], a[1]}
	ga, errA := engine.LessonOrder(orderGraph, a)
	gb, errB := engine.LessonOrder(orderGraph, b)
	if errA != nil || errB != nil || !slices.Equal(ga, gb) {
		t.Errorf("%+v %v and %+v %v differ", ga, errA, gb, errB)
	}
}

// Lessons that need each other through a cycle still come out, each once, by tier and ID.
func TestLessonOrderOfACycle(t *testing.T) {
	g := learnermodel.Graph{Module: orderModule, Skills: []learnermodel.Skill{orderSkill("a", 1, "b"), orderSkill("b", 1, "a")}}
	got, err := engine.LessonOrder(g, []engine.Lesson{orderLesson("b", "b"), orderLesson("a", "a")})
	if err != nil {
		t.Fatal(err)
	}
	want := []engine.OrderedLesson{{ID: orderModule + "/lesson/a", Group: 0}, {ID: orderModule + "/lesson/b", Group: 1}}
	if !slices.Equal(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestLessonOrderOfNoLessons(t *testing.T) {
	got, err := engine.LessonOrder(orderGraph, nil)
	if err != nil || len(got) != 0 {
		t.Errorf("got %+v, %v", got, err)
	}
}

func TestLessonOrderRefusesALessonTwice(t *testing.T) {
	_, err := engine.LessonOrder(orderGraph, []engine.Lesson{orderLesson("gaps", "gaps"), orderLesson("gaps", "counting")})
	if !errors.Is(err, engine.ErrInvalid) || err.Error() != `engine: invalid input: lesson "novatrail/olympiad-math/lesson/gaps" comes twice` {
		t.Errorf("got %v", err)
	}
}
