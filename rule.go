package engine

import (
	"cmp"
	"slices"
	"strings"
)

// The branches of the rule in one module, for the skills with a verdict (docs/architecture.md,
// section 3.1; R03): consolidation, the probe, a new skill, the queue, all mastered and no skill
// within reach. The first branch that fits makes the choice, and its code is the choice's (R07).
// The open run, the skill chosen on the map, the repetition, the author's order and the lessons
// come with T18–T20 (RUN.md).

// The turn of the probe: the module's n-th choice, counted before this one, with n > 5 and
// n mod 6 = 5, as in standalone. Only a run of sim changes them.
const (
	probeEvery int64 = 6
	probeAfter int64 = 5
)

// chooseIn runs the branches of the rule in the module at place i of the input, in one pass, and
// tells whether one of them gave a task.
func (pr prepared) chooseIn(i int, second bool) (Choice, bool) {
	in := pr.in
	module := in.Modules[i].Graph.Module
	mr := moduleReadout(in, module)
	var n int64
	if j := slices.IndexFunc(in.Summary.Modules, func(m ModuleChoices) bool { return m.Module == module }); j >= 0 {
		n = in.Summary.Modules[j].Choices
	}
	skills := pr.skills(i, second)

	// given is the task on a skill a branch chose: every branch chooses a skill with candidates.
	given := func(s skillView, easiest bool, code ReasonCode) (Choice, bool) {
		c, _ := s.task(easiest)
		return Choice{Kind: ChoiceTask, ID: c.id, Reason: code}, true
	}

	// Branch 3: the last closed attempt is the module's, a failure on a primary skill, past the
	// trial series, and not left; of its skills within reach, the first in the order of the graph.
	// In the trial series, and on its last answer, a mistake is not worked on.
	if in.HasLast && !in.Last.Left && strings.HasPrefix(in.Last.Task, module+"/") && mr.TrialOver && !mr.LastInTrial {
		for _, s := range skills {
			if s.reachable && slices.ContainsFunc(in.Last.Results, func(r SkillResult) bool { return r.Skill == s.id && failed(r) }) {
				return given(s, false, ReasonConsolidation)
			}
		}
	}

	// Branch 4: past the trial series, on the turn of the probe, a skill on pause with candidates,
	// the easiest of them.
	if mr.TrialOver && n > probeAfter && n%probeEvery == probeAfter {
		if s, ok := longestAgo(skills, func(s skillView) bool { return s.paused() && len(s.candidates) > 0 }); ok {
			return given(s, true, ReasonProbe)
		}
	}

	// Branch 6: a skill within reach not handed out before, the first in the order of the graph.
	for _, s := range skills {
		if s.reachable && !s.issued {
			return given(s, false, ReasonNewSkill)
		}
	}

	// Branch 7: a skill within reach, handed out before and not mastered for the queue, the one
	// not handed out for longest. A skill the model has mastered is back in rotation.
	if s, ok := longestAgo(skills, func(s skillView) bool { return s.reachable && s.issued && !s.masteredForQueue }); ok {
		if s.read.Mastered {
			return given(s, false, ReasonRotation)
		}
		return given(s, false, ReasonQueue)
	}

	// Branch 8: every skill within reach is mastered for the queue — the choice of branch 7 among
	// them all.
	if s, ok := longestAgo(skills, func(s skillView) bool { return s.reachable }); ok {
		return given(s, false, ReasonRotation)
	}

	// Branch 9: no skill within reach — branches 6 and 7 among the open skills with candidates at
	// the bottom of the ladder, whose lowest band taught is the lowest of theirs; and when every
	// skill there given before is mastered for the queue, the choice of branch 7 among them all, as
	// branch 8 does, so that an open skill with candidates is never left with nothing to give.
	var ladder []skillView
	for _, s := range skills {
		if s.open && len(s.candidates) > 0 && s.hasLowestBand {
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
			return given(s, false, ReasonNewSkill)
		}
	}
	if s, ok := longestAgo(ladder, func(s skillView) bool { return !s.masteredForQueue }); ok {
		return given(s, false, ReasonRotation)
	}
	s, _ := longestAgo(ladder, func(skillView) bool { return true })
	return given(s, false, ReasonRotation)
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
