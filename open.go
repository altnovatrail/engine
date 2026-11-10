package engine

import (
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Which skills of the modules of the engine are open (docs/architecture.md, sections 3.1 and 3.2;
// R03, R04): with evidence of their own, or with each prerequisite counted — a skill of the module
// mastered by the model, or an anchor skill the learner has, does not know, or has not while it
// cannot hold the skill. It does not depend on the pass, and the rule reads it once a call.

// anchorPrefix opens the ID of an anchor skill of the shared taxonomy, `core.<name>`.
const anchorPrefix = "core."

// node is a skill of a module of the engine in the graph G of section 3.2.
type node struct {
	skill  learnermodel.Skill
	module int // its module's place in Input.Modules
	read   learnermodel.SkillReadout
}

// openSkills gives the open skills of the modules of the engine, by full ID, sorted.
//
// An anchor skill the learner has not holds a skill, but not in two cases. One is a cycle: a skill
// that teaches it lies in one strongly connected component of G with the skill. The other is an
// anchor with nothing to teach it: no skill with a verdict that teaches it has a task of the pool,
// on pause or not, whose every other primary skill is open, nor a lesson. That last "open" is the
// one without this case, so that the count is one step and ends (the amendment of 2026-11-08 to
// R04).
//
// A prerequisite of the module is counted when the model has mastered it, and one with a rubric or
// none when a lesson that teaches it is passed or a task where it is primary was closed by "Done"
// (docs/architecture.md, section 3.3). One with nothing to teach it holds nothing either, by the
// same one step: one with a verdict has no lesson and no task as the teacher of an anchor has none,
// and one with a rubric or none has no lesson and no task with a rubric of the pool, where it is
// primary, that the rule could give (the amendments of 2026-11-08 and 2026-11-09 to R05).
func (pr prepared) openSkills() []string {
	in := pr.in
	var nodes []node
	for i, m := range in.Modules {
		reads := sortedBy(moduleReadout(in, m.Graph.Module).Skills, func(r learnermodel.SkillReadout) string { return r.Skill })
		for _, s := range m.Graph.Skills {
			r, _ := lookup(reads, func(r learnermodel.SkillReadout) string { return r.Skill }, s.ID)
			nodes = append(nodes, node{skill: s, module: i, read: r})
		}
	}
	slices.SortStableFunc(nodes, func(a, b node) int { return strings.Compare(a.skill.ID, b.skill.ID) })
	ids := make([]string, len(nodes))
	for k, n := range nodes {
		ids[k] = n.skill.ID
	}
	place := func(id string) (int, bool) { return slices.BinarySearch(ids, id) }
	anchors := sortedBy(in.Readout.Anchors, func(a learnermodel.AnchorReadout) string { return a.Anchor })

	// The anchor skills the graphs name, sorted, and the skills that teach each, by place in nodes.
	var named []string
	for _, n := range nodes {
		for _, id := range slices.Concat(n.skill.Prerequisites, n.skill.Teaches) {
			if strings.HasPrefix(id, anchorPrefix) {
				named = append(named, id)
			}
		}
	}
	named = slices.Compact(sortedCopy(named))
	teaching := make([][]int, len(named))
	for k, n := range nodes {
		for _, id := range n.skill.Teaches {
			if a, ok := slices.BinarySearch(named, id); ok {
				teaching[a] = append(teaching[a], k)
			}
		}
	}
	teachers := func(anchor string) []int {
		a, _ := slices.BinarySearch(named, anchor)
		return teaching[a]
	}
	// The edges of G: from a skill to each prerequisite of its module and to each skill that
	// teaches one of its anchor prerequisites.
	edges := make([][]int, len(nodes))
	for k, n := range nodes {
		for _, p := range n.skill.Prerequisites {
			if strings.HasPrefix(p, anchorPrefix) {
				edges[k] = append(edges[k], teachers(p)...)
			} else if j, ok := place(p); ok {
				edges[k] = append(edges[k], j)
			}
		}
	}
	component := components(edges)

	// counted tells whether a prerequisite of the skill at place k is counted; nothingToTeach, when
	// given, tells whether a prerequisite not counted otherwise has nothing to teach it.
	counted := func(k int, p string, nothingToTeach func(string) bool) bool {
		if !strings.HasPrefix(p, anchorPrefix) {
			j, ok := place(p)
			if !ok {
				return false
			}
			if nodes[j].skill.Assessability == learnermodel.AssessabilityVerdict {
				if nodes[j].read.Mastered {
					return true
				}
			} else if pr.passedOrDone(p) {
				return true
			}
			return nothingToTeach != nil && nothingToTeach(p)
		}
		a, ok := lookup(anchors, func(a learnermodel.AnchorReadout) string { return a.Anchor }, p)
		switch {
		case !ok:
			return false
		case a.Status == learnermodel.AnchorHas || a.Status == learnermodel.AnchorUnknown:
			return true
		case a.Status != learnermodel.AnchorHasNot:
			return false
		}
		if slices.ContainsFunc(teachers(p), func(j int) bool { return component[j] == component[k] }) {
			return true
		}
		return nothingToTeach != nil && nothingToTeach(p)
	}
	openWith := func(nothingToTeach func(string) bool) []string {
		var open []string
		for k, n := range nodes {
			if n.read.EvidenceCount > 0 || !slices.ContainsFunc(n.skill.Prerequisites, func(p string) bool { return !counted(k, p, nothingToTeach) }) {
				open = append(open, n.skill.ID)
			}
		}
		return open
	}

	// Whether each skill has something to teach it, once, by the "open" without the cases of
	// nothing to teach. A skill with a verdict has, when the model may master it — with a band it is
	// taught at, as learner-model reads a teacher — and a lesson not passed names it, or a task the
	// rule could give it, an item with a verdict where it is primary, has every other primary skill
	// open; a skill with a rubric or none, when a lesson not passed names it or a task with a rubric
	// of the pool teaches it — has it primary, outside the lessons of its module or in one not
	// passed, with every primary skill with a verdict open. A lesson passed teaches nothing more: no
	// branch of "My subjects" gives it again.
	first := openWith(nil)
	isOpen := func(id string) bool {
		_, found := slices.BinarySearch(first, id)
		return found
	}
	known, can := make([]bool, len(nodes)), make([]bool, len(nodes))
	canTeach := func(k int) bool {
		if known[k] {
			return can[k]
		}
		known[k] = true
		n := nodes[k]
		verdict := n.skill.Assessability == learnermodel.AssessabilityVerdict
		if verdict && len(n.skill.BandsTaught) == 0 {
			return false
		}
		if slices.ContainsFunc(in.Modules[n.module].Lessons, func(l Lesson) bool { return slices.Contains(l.Skills, n.skill.ID) && !pr.isPassed(l.ID) }) {
			can[k] = true
			return true
		}
		// Only a task the skill is primary in can teach it, of either kind.
		can[k] = slices.ContainsFunc(pr.tasksOn(n.module, n.skill.ID), func(j int) bool {
			t := pr.pools[n.module][j]
			if !verdict {
				// A task with a rubric teaches a skill with a rubric primary in its item with a
				// rubric, and a skill with none primary in it, as its "Done" counts either; and only
				// when the rule could give it: outside the lessons of its module, by branch 6a, or in
				// a lesson of it not passed, by that lesson — one passed is not given again —, and
				// with every primary skill with a verdict open, or the prerequisite and a skill after
				// it, primary in one task, would wait on each other (the amendment of 2026-11-09 to
				// R05).
				ms := pr.mods[n.module]
				rubric := ms.rubricOf(t.c.Task)
				teaches := slices.Contains(rubric, n.skill.ID) ||
					n.skill.Assessability == learnermodel.AssessabilityNone && len(rubric) > 0 && slices.Contains(t.primaries, n.skill.ID)
				lessons := in.Modules[n.module].Lessons
				given := !inLessons(lessons, t.c.Task.ID) || slices.ContainsFunc(lessons, func(l Lesson) bool {
					return !pr.isPassed(l.ID) && slices.ContainsFunc(l.Tasks, func(lt learnermodel.Task) bool { return lt.ID == t.c.Task.ID })
				})
				return teaches && given && !slices.ContainsFunc(t.primaries, func(s string) bool { return ms.verdict(s) && !isOpen(s) })
			}
			_, ok := candidateOf(t, n.skill.ID, nil)
			return ok && othersOpen(t.primaries, n.skill.ID, isOpen)
		})
		return can[k]
	}
	return openWith(func(p string) bool {
		if strings.HasPrefix(p, anchorPrefix) {
			return !slices.ContainsFunc(teachers(p), func(j int) bool {
				return nodes[j].skill.Assessability == learnermodel.AssessabilityVerdict && canTeach(j)
			})
		}
		j, _ := place(p)
		return !canTeach(j)
	})
}

// passedOrDone tells whether a prerequisite with a rubric or none is counted: a run that ended
// "passed" names it, or it is primary in a task of an attempt closed by "Done".
func (pr prepared) passedOrDone(skill string) bool {
	_, passed := slices.BinarySearch(pr.passedSkills, skill)
	_, done := slices.BinarySearch(pr.doneSkills, skill)
	return passed || done
}

// components gives each vertex of a directed graph its strongly connected component, by Tarjan's
// algorithm: vertices in one component get one number. The vertices are visited in order, so the
// numbers do not depend on anything but the graph.
func components(edges [][]int) []int {
	n := len(edges)
	index, low, component := make([]int, n), make([]int, n), make([]int, n)
	for v := range n {
		index[v], component[v] = -1, -1
	}
	var stack []int
	onStack := make([]bool, n)
	next, count := 0, 0
	var visit func(v int)
	visit = func(v int) {
		index[v], low[v] = next, next
		next++
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range edges[v] {
			if index[w] < 0 {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] == index[v] {
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				component[w] = count
				if w == v {
					break
				}
			}
			count++
		}
	}
	for v := range n {
		if index[v] < 0 {
			visit(v)
		}
	}
	return component
}
