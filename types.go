package engine

import learnermodel "github.com/altnovatrail/learner-model"

// The types of the engine's input and output (docs/architecture.md, section 3.6; R09). The graph,
// the description of a task and the learner model's reading are learner-model's, taken as they
// are; the rest is the engine's own, simple as they are. Times are whole seconds of Unix time,
// UTC, as learner-model's. A field that may have no value by its meaning is a value with a flag
// beside it, with no pointer.

// Input is all that one choice reads. The order of its lists does not change the choice: the
// rule orders by itself, by the order of the graph, of the lessons and of the IDs.
type Input struct {
	Time int64 // the time of the choice
	// The learner model's reading at Time, by learnermodel.Read: the modules of the engine and
	// the modules of the evidence log.
	Readout learnermodel.Readout
	Modules []Module // the modules of the engine: connected, and suited to the learner
	// The full ID of the skill the learner chose on the skill map; empty in "My subjects".
	MapSkill   string
	Runs       []Run        // the lesson runs of the profile in the modules of the engine
	RubricDone []RubricDone // the attempts of tasks with a rubric closed by "Done"
	Last       Last         // the last closed attempt of the profile, when HasLast
	HasLast    bool
	Summary    Summary // what was handed out before
}

// Module is a module of the engine: the graph of its current version, and the pool and the
// lessons api gives the learner of it.
type Module struct {
	Graph   learnermodel.Graph
	Pool    []Candidate
	Lessons []Lesson // the lessons that suit the learner whole
}

// Candidate is a task of the pool: its description, and the traps of its items.
type Candidate struct {
	Task  learnermodel.Task
	Traps []ItemTraps
}

// ItemTraps is the traps of an item, by its place: the trap_values of its checks, by full ID.
type ItemTraps struct {
	Place int32
	Traps []string
}

// Lesson is a lesson of a module: its full ID, its skills, and the descriptions of the tasks of
// its segments, in the order of the segments.
type Lesson struct {
	ID     string
	Skills []string
	Tasks  []learnermodel.Task
}

// Run is a lesson run of the profile: its lesson and the lesson's skills — translated, on a
// package before — how it ended, and when. The runs of one lesson repeat by right: a lesson is
// given again by a new run.
type Run struct {
	Lesson  string
	Skills  []string
	End     RunEnd // RunOpen while it is open
	EndedAt int64  // zero while it is open
}

// RubricDone is an attempt of a task with a rubric that "Done" closed, and when. The attempts of
// one task repeat by right: a task with a rubric is given again by a new attempt.
type RubricDone struct {
	Task string
	At   int64
}

// Last is the last closed attempt of the profile.
type Last struct {
	Task     string        // the full ID of its task, whose beginning is its module
	Left     bool          // whether it was closed as left: its outcome is "left"
	Lesson   string        // the full ID of the lesson, when InLesson
	InLesson bool          // whether the attempt is a segment of a run of that lesson
	Results  []SkillResult // the result on each of its primary skills
}

// SkillResult is the result of an attempt on a skill, in facts: of the items of the attempt the
// skill is primary in, how many got a first verdict, how many of those are right, and whether the
// solution was shown. The engine tells a failure from them, by the same rule as the learner model.
type SkillResult struct {
	Skill         string
	Answered      int32
	FirstRight    int32
	SolutionShown bool
}

// Summary is what api gives of what was handed out before, from the records of choices and
// attempts it keeps.
type Summary struct {
	Skills  []SkillIssue
	Modules []ModuleChoices
	Tasks   []TaskSolution
}

// SkillIssue is when a skill was last handed out, by a task on it or by a lesson that teaches it,
// when Issued; and, when HasLast, the reason code of the choice of its last closed attempt and the
// attempt's result on it.
type SkillIssue struct {
	Skill      string
	LastIssued int64
	Issued     bool
	LastCode   ReasonCode
	LastResult SkillResult
	HasLast    bool
}

// ModuleChoices is the number of choices of a module, a lesson counted once.
type ModuleChoices struct {
	Module  string
	Choices int64
}

// TaskSolution is when the solution of a task last went to the learner's device.
type TaskSolution struct {
	Task           string
	SolutionSentAt int64
}

// Choice is the engine's choice: a lesson or a task by its full ID — of a template, the
// template's — or nothing to give, with no ID; and the reason code of the branch that made it,
// ReasonUnspecified beside nothing to give.
type Choice struct {
	Kind   ChoiceKind
	ID     string
	Reason ReasonCode
}

// ReasonCode is why the engine gave what it chose: one value of the closed list ReasonCode of
// cause.proto in spec's contracts, with its number, which api translates the code by. The adult
// sees it; the learner is not told it. A test holds the names and numbers to the file.
type ReasonCode int32

const (
	ReasonUnspecified   ReasonCode = 0 // a choice of nothing
	ReasonConsolidation ReasonCode = 1 // after a wrong answer, another task on the skill
	ReasonNewSkill      ReasonCode = 2 // a skill not given before, whose prerequisites are met
	ReasonQueue         ReasonCode = 3 // the skill within reach not given for longest in its module
	ReasonRepetition    ReasonCode = 4 // the review date of a mastered skill has come
	ReasonProbe         ReasonCode = 5 // a skill put aside comes back as a check
	ReasonSkillChosen   ReasonCode = 6 // the skill chosen on the skill map
	ReasonAuthorOrder   ReasonCode = 7 // a skill with a rubric or without assessment, by the graph
	ReasonOpenRun       ReasonCode = 8 // the lesson of a run that is open
	ReasonRotation      ReasonCode = 9 // a skill given before comes round again, outside the queue
)

// ChoiceKind is what a choice gives: ChoiceKind of cause.proto, names and numbers.
type ChoiceKind int32

const (
	ChoiceUnspecified ChoiceKind = 0 // never given
	ChoiceLesson      ChoiceKind = 1
	ChoiceTask        ChoiceKind = 2 // a task outside a lesson
	ChoiceNothing     ChoiceKind = 3
)

// RunEnd is how a lesson run ended: RunEndReason of cause.proto, names and numbers. Its zero is a
// run that is open.
type RunEnd int32

const (
	RunOpen        RunEnd = 0
	RunPassed      RunEnd = 1
	RunUnchecked   RunEnd = 2
	RunNotSuitable RunEnd = 3
	RunAbandoned   RunEnd = 4
)
