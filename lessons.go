package engine

import (
	"cmp"
	"math"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// Lessons, the author's order and nothing to give (docs/architecture.md, section 3.3; R05): the
// order of the lessons, the lesson on a skill, branch 1 — the lesson of an open run —, branch 6a —
// the author's order —, and the skill chosen on the map with a rubric or none.

// moduleSkills is what the lessons read of the graph of a module: its skills by ID, sorted, with
// the assessability and the tier of each, and the skills of the module each needs through its
// prerequisites, sorted.
type moduleSkills struct {
	ids    []string
	assess []learnermodel.Assessability
	tiers  []int32
	needs  [][]string
}

func newModuleSkills(g learnermodel.Graph) moduleSkills {
	skills := sortedBy(g.Skills, func(s learnermodel.Skill) string { return s.ID })
	ms := moduleSkills{
		ids: make([]string, len(skills)), assess: make([]learnermodel.Assessability, len(skills)),
		tiers: make([]int32, len(skills)), needs: make([][]string, len(skills)),
	}
	for k, s := range skills {
		ms.ids[k], ms.assess[k], ms.tiers[k] = s.ID, s.Assessability, s.Tier
	}
	for k := range skills {
		seen := make([]bool, len(skills))
		stack := []int{k}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, p := range skills[v].Prerequisites {
				if j, ok := ms.at(p); ok && !seen[j] {
					seen[j] = true
					stack = append(stack, j)
				}
			}
		}
		for j, s := range seen {
			if s {
				ms.needs[k] = append(ms.needs[k], ms.ids[j])
			}
		}
	}
	return ms
}

func (ms moduleSkills) at(id string) (int, bool) { return slices.BinarySearch(ms.ids, id) }

// verdict tells whether a skill of the module has the assessability verdict.
func (ms moduleSkills) verdict(id string) bool {
	k, ok := ms.at(id)
	return ok && ms.assess[k] == learnermodel.AssessabilityVerdict
}

// authorOrder tells whether a lesson of these skills goes by the author's order: not all of its
// skills have a verdict, a mixed one too.
func (ms moduleSkills) authorOrder(skills []string) bool {
	return slices.ContainsFunc(skills, func(s string) bool { return !ms.verdict(s) })
}

// after tells whether a lesson of the skills a comes after one of the skills b, as the format's
// text, section 21, has it: a skill of a needs, through its prerequisites among the module's
// skills, a skill of b that is not one of a's own. An anchor skill orders nothing.
func (ms moduleSkills) after(a, b []string) bool {
	for _, x := range a {
		k, ok := ms.at(x)
		if !ok {
			continue
		}
		for _, y := range b {
			if _, needs := slices.BinarySearch(ms.needs[k], y); needs && !slices.Contains(a, y) {
				return true
			}
		}
	}
	return false
}

// tier gives the tier of a lesson of these skills: the least of theirs (the amendment of
// 2026-11-08 to R05). The format computes a tier from the graph (the format's text, section 12).
func (ms moduleSkills) tier(skills []string) int32 {
	t := int32(math.MaxInt32)
	for _, s := range skills {
		if k, ok := ms.at(s); ok {
			t = min(t, ms.tiers[k])
		}
	}
	return t
}

// lessonAt is a lesson — or the lesson of a run — among those whose first in the order of the
// lessons is sought: its module's place in Input.Modules, its full ID and its skills.
type lessonAt struct {
	module int
	id     string
	skills []string
}

// firstLesson gives the first of the lessons in the order of the lessons. Of each module's, the
// first is one that comes after no other of them, of the least tier, then by ID byte by byte; of
// the modules' firsts, by full ID: the lessons of two modules are not ordered, and each module's
// tiers are its own. The runs of one module on a package before may each come after another by
// the graph of the current one; with none that comes after no other, its first is by tier and ID
// among them all.
func (pr prepared) firstLesson(ls []lessonAt) (lessonAt, bool) {
	comesAfter := func(l lessonAt) bool {
		return slices.ContainsFunc(ls, func(o lessonAt) bool { return o.module == l.module && pr.mods[l.module].after(l.skills, o.skills) })
	}
	var firsts []lessonAt
	for _, l := range ls {
		ms := pr.mods[l.module]
		if comesAfter(l) && slices.ContainsFunc(ls, func(o lessonAt) bool { return o.module == l.module && !comesAfter(o) }) {
			continue
		}
		j := slices.IndexFunc(firsts, func(f lessonAt) bool { return f.module == l.module })
		switch {
		case j < 0:
			firsts = append(firsts, l)
		case cmp.Or(cmp.Compare(ms.tier(l.skills), ms.tier(firsts[j].skills)), strings.Compare(l.id, firsts[j].id)) < 0:
			firsts[j] = l
		}
	}
	if len(firsts) == 0 {
		return lessonAt{}, false
	}
	return slices.MinFunc(firsts, func(a, b lessonAt) int { return strings.Compare(a.id, b.id) }), true
}

// isPassed tells whether a lesson has a run that ended "passed", on any package; isRunning,
// whether it has an open run.
func (pr prepared) isPassed(lesson string) bool {
	_, found := slices.BinarySearch(pr.passed, lesson)
	return found
}

func (pr prepared) isRunning(lesson string) bool {
	_, found := slices.BinarySearch(pr.running, lesson)
	return found
}

// isDone tells whether an attempt of a task with a rubric was closed by "Done".
func (pr prepared) isDone(task string) bool {
	_, ok := pr.lastDone(task)
	return ok
}

// lastDone gives when the latest attempt of a task with a rubric was closed by "Done".
func (pr prepared) lastDone(task string) (int64, bool) {
	d, ok := lookup(pr.done, func(d RubricDone) string { return d.Task }, task)
	return d.At, ok
}

// isIssued tells whether a skill was handed out before, by the summary.
func (pr prepared) isIssued(skill string) bool {
	_, found := slices.BinarySearch(pr.issued, skill)
	return found
}

// lastPassed gives when the latest run of a lesson that ended "passed" ended.
func (pr prepared) lastPassed(lesson string) int64 {
	var at int64
	found := false
	for _, r := range pr.in.Runs {
		if r.Lesson == lesson && r.End == RunPassed && (!found || r.EndedAt > at) {
			at, found = r.EndedAt, true
		}
	}
	return at
}

// suits tells whether a lesson suits a band: the band is one of its tasks', or it has no task.
func suits(l Lesson, band int32) bool {
	return len(l.Tasks) == 0 || slices.ContainsFunc(l.Tasks, func(t learnermodel.Task) bool { return t.Band.Place == band })
}

// lessonOnSkill gives the lesson on a skill with a verdict of the module at place i: a lesson not
// passed and with no open run, whose skills all have a verdict and are open — on the map, whatever
// their state, as the map's tasks are —, the skill among them, and which suits the band of the
// skill's recommended task; with no recommended task, any. Off the map, a lesson that names another
// skill which awaits its lesson of the author's order is none: that skill comes by that lesson. Of
// several, the first in the order of the lessons.
func (pr prepared) lessonOnSkill(i int, skill string, rec candidate, hasRec, onMap bool) (string, bool) {
	ms := pr.mods[i]
	var ls []lessonAt
	for _, l := range pr.in.Modules[i].Lessons {
		if !slices.Contains(l.Skills, skill) || pr.isPassed(l.ID) || pr.isRunning(l.ID) || ms.authorOrder(l.Skills) ||
			!onMap && slices.ContainsFunc(l.Skills, func(s string) bool { return !pr.isOpen(s) || s != skill && !pr.isIssued(s) && pr.awaitsLesson(i, s) }) ||
			hasRec && !suits(l, rec.band) {
			continue
		}
		ls = append(ls, lessonAt{i, l.ID, l.Skills})
	}
	l, ok := pr.firstLesson(ls)
	return l.id, ok
}

// awaitsLesson tells whether a lesson of the author's order not passed teaches a skill of the
// module at place i: a skill with a verdict not handed out comes by that lesson, and branches 6–9
// do not take it (the amendment of 2026-11-08 to R05).
func (pr prepared) awaitsLesson(i int, skill string) bool {
	ms := pr.mods[i]
	return slices.ContainsFunc(pr.in.Modules[i].Lessons, func(l Lesson) bool {
		return slices.Contains(l.Skills, skill) && ms.authorOrder(l.Skills) && !pr.isPassed(l.ID)
	})
}

// openRun is branch 1, the lesson of an open run of a module of the engine, with the code
// REASON_CODE_OPEN_RUN: of several, the first in the order of the lessons, read by the skills of
// the run, which may be on a package before. It does not give the run whose segment's attempt the
// learner has just left, asking for something else; on the map it gives only the run of a lesson
// that teaches the skill chosen there.
func (pr prepared) openRun() (Choice, bool) {
	in := pr.in
	var ls []lessonAt
	for _, r := range in.Runs {
		i := slices.IndexFunc(in.Modules, func(m Module) bool { return strings.HasPrefix(r.Lesson, m.Graph.Module+"/") })
		if r.End != RunOpen || i < 0 || in.HasLast && in.Last.Left && in.Last.InLesson && in.Last.Lesson == r.Lesson ||
			in.MapSkill != "" && !slices.Contains(r.Skills, in.MapSkill) {
			continue
		}
		ls = append(ls, lessonAt{i, r.Lesson, r.Skills})
	}
	if l, ok := pr.firstLesson(ls); ok {
		return Choice{Kind: ChoiceLesson, ID: l.id, Reason: ReasonOpenRun}, true
	}
	return Choice{}, false
}

// inLessons tells whether a lesson the engine gets names a task.
func inLessons(lessons []Lesson, task string) bool {
	return slices.ContainsFunc(lessons, func(l Lesson) bool {
		return slices.ContainsFunc(l.Tasks, func(t learnermodel.Task) bool { return t.ID == task })
	})
}

// rubricOf gives the primary skills with a rubric of a task of the module, sorted: a task with a
// rubric has one.
func (ms moduleSkills) rubricOf(t learnermodel.Task) []string {
	var out []string
	for _, s := range primaries(t) {
		if k, ok := ms.at(s); ok && ms.assess[k] == learnermodel.AssessabilityRubric {
			out = append(out, s)
		}
	}
	return slices.Compact(sortedCopy(out))
}

// authorOrder is branch 6a in the module at place i, with the code REASON_CODE_AUTHOR_ORDER: the
// first lesson of the author's order not passed and with no open run, all of whose lessons of the
// author's order it comes after are passed; then a task with a rubric outside the lessons, with no
// attempt closed by "Done", as a lesson of its skills with a rubric would — after the lessons of
// the author's order of its own skills and of those it would come after are passed (the amendment
// of 2026-11-08 to R05) —, in the order of its skills in the graph, then by ID. A lesson whose
// skills all have a verdict holds neither: it is given on the new and by band.
func (pr prepared) authorOrder(i int, second bool) (Choice, bool) {
	ms := pr.mods[i]
	m := pr.in.Modules[i]
	waits := func(skills []string, own bool) bool {
		return slices.ContainsFunc(m.Lessons, func(o Lesson) bool {
			return ms.authorOrder(o.Skills) && !pr.isPassed(o.ID) &&
				(ms.after(skills, o.Skills) || own && slices.ContainsFunc(o.Skills, func(s string) bool { return slices.Contains(skills, s) }))
		})
	}
	var ls []lessonAt
	for _, l := range m.Lessons {
		if ms.authorOrder(l.Skills) && !pr.isPassed(l.ID) && !pr.isRunning(l.ID) && !waits(l.Skills, false) {
			ls = append(ls, lessonAt{i, l.ID, l.Skills})
		}
	}
	if l, ok := pr.firstLesson(ls); ok {
		return Choice{Kind: ChoiceLesson, ID: l.id, Reason: ReasonAuthorOrder}, true
	}

	order := graphOrder(m.Graph.Skills)
	rank := func(skills []string) int {
		return slices.IndexFunc(order, func(s learnermodel.Skill) bool { return slices.Contains(skills, s.ID) })
	}
	type rubricTask struct {
		cd   candidate
		rank int
	}
	var ts []rubricTask
	for _, c := range m.Pool {
		rubric := ms.rubricOf(c.Task)
		if len(rubric) == 0 || inLessons(m.Lessons, c.Task.ID) || pr.isDone(c.Task.ID) || waits(rubric, true) || !pr.verdictsOpen(ms, c.Task) {
			continue
		}
		if cd, ok := pr.pause(c.Task, second, false); ok {
			ts = append(ts, rubricTask{cd, rank(rubric)})
		}
	}
	if len(ts) == 0 {
		return Choice{}, false
	}
	t := slices.MinFunc(ts, func(a, b rubricTask) int {
		return cmp.Or(bySolution(a.cd, b.cd), cmp.Compare(a.rank, b.rank), strings.Compare(a.cd.id, b.cd.id))
	})
	return Choice{Kind: ChoiceTask, ID: t.cd.id, Reason: ReasonAuthorOrder}, true
}

// verdictsOpen tells whether every primary skill with a verdict of a task is open: otherwise its
// evidence would open a closed skill past its prerequisites, as for a candidate.
func (pr prepared) verdictsOpen(ms moduleSkills, t learnermodel.Task) bool {
	return !slices.ContainsFunc(primaries(t), func(s string) bool { return ms.verdict(s) && !pr.isOpen(s) })
}

// mapAuthor is the skill chosen on the map with a rubric or none, in the module at place i, with
// the code REASON_CODE_SKILL_CHOSEN: a lesson not passed that teaches it, without waiting for those
// it comes after, since the learner chose the skill; else a task with a rubric on it outside the
// lessons with no "Done"; else the lesson of it whose latest run "passed" ended longest ago, by a
// new run; else the task with a rubric on it of the pool whose latest "Done" closed longest ago,
// one never done first, by a new attempt (the amendment of 2026-11-07 to R04). A lesson with an
// open run is branch 1's.
func (pr prepared) mapAuthor(i int, skill string, second bool) (Choice, bool) {
	m := pr.in.Modules[i]
	var fresh, passed []lessonAt
	for _, l := range m.Lessons {
		if !slices.Contains(l.Skills, skill) || pr.isRunning(l.ID) {
			continue
		}
		if pr.isPassed(l.ID) {
			passed = append(passed, lessonAt{i, l.ID, l.Skills})
		} else {
			fresh = append(fresh, lessonAt{i, l.ID, l.Skills})
		}
	}
	lesson := func(id string) (Choice, bool) {
		return Choice{Kind: ChoiceLesson, ID: id, Reason: ReasonSkillChosen}, true
	}
	if l, ok := pr.firstLesson(fresh); ok {
		return lesson(l.id)
	}

	type rubricTask struct {
		cd     candidate
		done   bool
		doneAt int64
	}
	var outside, all []rubricTask
	for _, c := range m.Pool {
		if !slices.Contains(pr.mods[i].rubricOf(c.Task), skill) {
			continue
		}
		cd, ok := pr.pause(c.Task, second, true)
		if !ok {
			continue
		}
		at, done := pr.lastDone(c.Task.ID)
		t := rubricTask{cd, done, at}
		all = append(all, t)
		if !done && !inLessons(m.Lessons, c.Task.ID) {
			outside = append(outside, t)
		}
	}
	task := func(ts []rubricTask) Choice {
		t := slices.MinFunc(ts, func(a, b rubricTask) int {
			return cmp.Or(bySolution(a.cd, b.cd), falseFirst(a.done, b.done), cmp.Compare(a.doneAt, b.doneAt), strings.Compare(a.cd.id, b.cd.id))
		})
		return Choice{Kind: ChoiceTask, ID: t.cd.id, Reason: ReasonSkillChosen}
	}
	if len(outside) > 0 {
		return task(outside), true
	}
	if len(passed) > 0 {
		l := slices.MinFunc(passed, func(a, b lessonAt) int {
			return cmp.Or(cmp.Compare(pr.lastPassed(a.id), pr.lastPassed(b.id)), strings.Compare(a.id, b.id))
		})
		return lesson(l.id)
	}
	if len(all) > 0 {
		return task(all), true
	}
	return Choice{}, false
}
