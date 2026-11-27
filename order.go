package engine

import (
	"slices"

	learnermodel "github.com/altnovatrail/learner-model"
)

// OrderedLesson is a lesson of a module in the author's order: its full ID, and the place of its
// group from 0. The lessons of one group have no order between them by the skill graph, so a page
// may show them in any order; the order of the list is still the one the rule takes them in.
type OrderedLesson struct {
	ID    string
	Group int32
}

// LessonOrder gives the lessons of one module in the author's order, by the same rule a choice
// takes them by: first a lesson that comes after none of the others by the prerequisites of the
// graph, of the least tier, then by ID byte by byte; then the first of the rest, and so on. A
// lesson starts a new group when it and a lesson of the group before it come one after the other
// by the graph; otherwise it joins that group. Only the lessons' IDs and skills are read, and the
// lessons need not suit the learner: the caller gives those it lists. A lesson twice by its ID is
// ErrInvalid.
func LessonOrder(g learnermodel.Graph, lessons []Lesson) ([]OrderedLesson, error) {
	ids := make([]string, len(lessons))
	for k, l := range lessons {
		ids[k] = l.ID
	}
	if id, ok := firstDuplicate(ids); ok {
		return nil, invalid("lesson " + shown(id) + " comes twice")
	}
	pr := prepared{mods: []moduleSkills{newModuleSkills(g)}}
	ms := pr.mods[0]
	left := make([]lessonAt, len(lessons))
	for k, l := range lessons {
		left[k] = lessonAt{id: l.ID, skills: l.Skills}
	}
	out := make([]OrderedLesson, 0, len(lessons))
	var group []lessonAt
	place := int32(0)
	for len(left) > 0 {
		first, _ := pr.firstLesson(left)
		ordered := func(o lessonAt) bool { return ms.after(first.skills, o.skills) || ms.after(o.skills, first.skills) }
		if slices.ContainsFunc(group, ordered) {
			place++
			group = group[:0]
		}
		group = append(group, first)
		out = append(out, OrderedLesson{ID: first.id, Group: place})
		left = slices.DeleteFunc(left, func(l lessonAt) bool { return l.id == first.id })
	}
	return out, nil
}
