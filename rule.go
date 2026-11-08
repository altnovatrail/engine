package engine

import (
	"cmp"
	"slices"
	"strings"

	learnermodel "github.com/altnovatrail/learner-model"
)

// The rule (docs/architecture.md, sections 3.1–3.3; R03, R04, R05): branch 1, the lesson of an
// open run; then "My subjects" — consolidation in the module of the last closed attempt, then the
// modules in turn, each running branches 4–9 with 6a, the author's order — or the skill chosen on
// the map, branch 2. Branches 6–9 give the lesson on the skill when there is one. The first branch
// that fits makes the choice, and its code is the choice's (R07). The repetition comes with T20
// (RUN.md).

// The turn of the probe: the module's n-th choice, counted before this one, with n > 5 and
// n mod 6 = 5, as in standalone. Only a run of sim changes them.
const (
	probeEvery int64 = 6
	probeAfter int64 = 5
)

// choose runs the rule: branch 1, which the pause does not touch; then a first pass, a second pass
// with the tasks on pause, then nothing to give.
func (pr prepared) choose() Choice {
	if c, ok := pr.openRun(); ok {
		return c
	}
	for _, second := range []bool{false, true} {
		if c, ok := pr.pass(second); ok {
			return c
		}
	}
	return Choice{Kind: ChoiceNothing}
}

// pass runs one pass of the rule after branch 1: on the skill chosen on the map, or in "My
// subjects".
func (pr prepared) pass(second bool) (Choice, bool) {
	if pr.in.MapSkill != "" {
		return pr.onMap(second)
	}
	views := make([][]skillView, len(pr.in.Modules))
	skillsOf := func(i int) []skillView {
		if views[i] == nil {
			views[i] = pr.skills(i, second)
		}
		return views[i]
	}
	if i, ok := pr.lastModule(); ok {
		if c, ok := pr.consolidation(i, skillsOf(i)); ok {
			return c, true
		}
	}
	for _, i := range pr.turns() {
		if c, ok := pr.turn(i, skillsOf(i), second); ok {
			return c, true
		}
	}
	return Choice{}, false
}

// given is the task on a skill a branch chose: every branch chooses a skill with candidates.
func given(s skillView, easiest bool, code ReasonCode) (Choice, bool) {
	c, _ := s.task(easiest)
	return Choice{Kind: ChoiceTask, ID: c.id, Reason: code}, true
}

// give is the choice of branches 6–9 on a skill of the module at place i: the lesson on the skill
// when there is one, else the task on it, with the branch's code (docs/architecture.md, section
// 3.3). Consolidation and the probe give only tasks.
func (pr prepared) give(i int, s skillView, code ReasonCode) (Choice, bool) {
	if id, ok := pr.lessonOnSkill(i, s.id, s.recommended, s.hasRecommended, false); ok {
		return Choice{Kind: ChoiceLesson, ID: id, Reason: code}, true
	}
	return given(s, false, code)
}

// lastModule gives the place of the module of the engine whose task the last closed attempt was.
func (pr prepared) lastModule() (int, bool) {
	if !pr.in.HasLast {
		return 0, false
	}
	i := slices.IndexFunc(pr.in.Modules, func(m Module) bool { return strings.HasPrefix(pr.in.Last.Task, m.Graph.Module+"/") })
	return i, i >= 0
}

// consolidation is branch 3 in the module of the last closed attempt: a failure on a primary
// skill, past the trial series, and not left; of its skills within reach, the first in the order
// of the graph. In the trial series, and on its last answer, a mistake is not worked on. It runs
// once, before the turn of the modules, and not again on a module's turn.
func (pr prepared) consolidation(i int, skills []skillView) (Choice, bool) {
	in := pr.in
	mr := moduleReadout(in, in.Modules[i].Graph.Module)
	if in.Last.Left || !mr.TrialOver || mr.LastInTrial {
		return Choice{}, false
	}
	for _, s := range skills {
		if s.reachable && slices.ContainsFunc(in.Last.Results, func(r SkillResult) bool { return r.Skill == s.id && failed(r) }) {
			return given(s, false, ReasonConsolidation)
		}
	}
	return Choice{}, false
}

// turns orders the modules of the engine for their turns: the one whose last handing out was
// longest ago first, a module never handed out before them all; of as long ago, by full ID byte by
// byte. A module was last handed out when the latest of its skills was, by the summary (R04): the
// rows whose full ID begins with the module's, so that a skill a newer version of the module
// dropped still counts.
func (pr prepared) turns() []int {
	type turn struct {
		place  int
		issued bool
		at     int64
	}
	var out []turn
	for i, m := range pr.in.Modules {
		t := turn{place: i}
		for _, is := range pr.in.Summary.Skills {
			if is.Issued && strings.HasPrefix(is.Skill, m.Graph.Module+"/") && (!t.issued || is.LastIssued > t.at) {
				t.issued, t.at = true, is.LastIssued
			}
		}
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b turn) int {
		if a.issued != b.issued {
			if !a.issued {
				return -1
			}
			return 1
		}
		return cmp.Or(cmp.Compare(a.at, b.at), strings.Compare(pr.in.Modules[a.place].Graph.Module, pr.in.Modules[b.place].Graph.Module))
	})
	places := make([]int, len(out))
	for k, t := range out {
		places[k] = t.place
	}
	return places
}

// turn runs branches 4–9 in the module at place i, on its turn.
func (pr prepared) turn(i int, skills []skillView, second bool) (Choice, bool) {
	in := pr.in
	module := in.Modules[i].Graph.Module
	mr := moduleReadout(in, module)
	var n int64
	if j := slices.IndexFunc(in.Summary.Modules, func(m ModuleChoices) bool { return m.Module == module }); j >= 0 {
		n = in.Summary.Modules[j].Choices
	}

	// Branch 4: past the trial series, on the turn of the probe, a skill on pause with candidates,
	// the easiest of them.
	if mr.TrialOver && n > probeAfter && n%probeEvery == probeAfter {
		if s, ok := longestAgo(skills, func(s skillView) bool { return s.paused() && len(s.candidates) > 0 }); ok {
			return given(s, true, ReasonProbe)
		}
	}

	// Branch 6: a skill within reach not handed out before, the first in the order of the graph.
	// One a lesson of the author's order not passed teaches comes by that lesson, and branches 6–9
	// do not take it.
	for _, s := range skills {
		if s.reachable && !s.issued && !s.awaitsLesson {
			return pr.give(i, s, ReasonNewSkill)
		}
	}

	// Branch 6a: the author's order.
	if c, ok := pr.authorOrder(i, second); ok {
		return c, true
	}

	// Branch 7: a skill within reach, handed out before and not mastered for the queue, the one
	// not handed out for longest. A skill the model has mastered is back in rotation.
	if s, ok := longestAgo(skills, func(s skillView) bool { return s.reachable && s.issued && !s.masteredForQueue }); ok {
		if s.read.Mastered {
			return pr.give(i, s, ReasonRotation)
		}
		return pr.give(i, s, ReasonQueue)
	}

	// Branch 8: every skill within reach is mastered for the queue — the choice of branch 7 among
	// them all.
	if s, ok := longestAgo(skills, func(s skillView) bool { return s.reachable && !s.awaitsLesson }); ok {
		return pr.give(i, s, ReasonRotation)
	}

	// Branch 9: no skill within reach — branches 6 and 7 among the open skills with candidates at
	// the bottom of the ladder, whose lowest band taught is the lowest of theirs; and when every
	// skill there given before is mastered for the queue, the choice of branch 7 among them all, as
	// branch 8 does, so that an open skill with candidates is never left with nothing to give.
	var ladder []skillView
	for _, s := range skills {
		if s.open && len(s.candidates) > 0 && s.hasLowestBand && !s.awaitsLesson {
			ladder = append(ladder, s)
		}
	}
	if len(ladder) == 0 {
		return Choice{}, false
	}
	bottom := slices.MinFunc(ladder, func(a, b skillView) int { return cmp.Compare(a.lowestBand, b.lowestBand) }).lowestBand
	ladder = slices.DeleteFunc(ladder, func(s skillView) bool { return s.lowestBand != bottom })
	for _, s := range ladder {
		if !s.issued {
			return pr.give(i, s, ReasonNewSkill)
		}
	}
	if s, ok := longestAgo(ladder, func(s skillView) bool { return !s.masteredForQueue }); ok {
		return pr.give(i, s, ReasonRotation)
	}
	s, _ := longestAgo(ladder, func(skillView) bool { return true })
	return pr.give(i, s, ReasonRotation)
}

// longestAgo gives, of the skills that fit, the one not handed out for longest — one never
// handed out first —, and of as long ago the first in the order of the graph: the skills come in
// that order.
func longestAgo(skills []skillView, fits func(skillView) bool) (skillView, bool) {
	var best skillView
	found := false
	for _, s := range skills {
		if !fits(s) {
			continue
		}
		if !found || byIssue(s, best) < 0 {
			best, found = s, true
		}
	}
	return best, found
}

// byIssue orders skills by when they were last handed out: one never handed out first.
func byIssue(a, b skillView) int {
	if a.issued != b.issued {
		if !a.issued {
			return -1
		}
		return 1
	}
	return cmp.Compare(a.lastIssued, b.lastIssued)
}

// onMap is branch 2, the skill chosen on the map (R04), with the code REASON_CODE_SKILL_CHOSEN,
// whatever its state — closed, on pause, in the trial series. A skill with a verdict gets the
// lesson on it when there is one, else a task on it among all its tasks; one with a rubric or none
// goes by mapAuthor (docs/architecture.md, section 3.3). A skill outside the modules of the engine
// gives nothing.
func (pr prepared) onMap(second bool) (Choice, bool) {
	skill := pr.in.MapSkill
	for i, m := range pr.in.Modules {
		j := slices.IndexFunc(m.Graph.Skills, func(s learnermodel.Skill) bool { return s.ID == skill })
		if j < 0 {
			continue
		}
		if m.Graph.Skills[j].Assessability != learnermodel.AssessabilityVerdict {
			return pr.mapAuthor(i, skill, second)
		}
		read, _ := lookup(sortedBy(moduleReadout(pr.in, m.Graph.Module).Skills, func(r learnermodel.SkillReadout) string { return r.Skill }),
			func(r learnermodel.SkillReadout) string { return r.Skill }, skill)
		cands, all := pr.candidates(pr.poolOf(i), skill, frequentTraps(read.Traps), second, true)
		var rec candidate
		if len(all) > 0 {
			rec = slices.MinFunc(all, byCorridor)
		}
		if id, ok := pr.lessonOnSkill(i, skill, rec, len(all) > 0, true); ok {
			return Choice{Kind: ChoiceLesson, ID: id, Reason: ReasonSkillChosen}, true
		}
		if len(cands) == 0 {
			return Choice{}, false
		}
		return Choice{Kind: ChoiceTask, ID: cands[0].id, Reason: ReasonSkillChosen}, true
	}
	return Choice{}, false
}
