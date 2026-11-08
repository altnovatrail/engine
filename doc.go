// Package engine is NovaTrail's engine: it picks a learner's next lesson or task by a rule, and
// gives the choice a reason code, which the adult sees and the learner is not told.
//
// It is a library with no input or output of its own. The api service — at stage 1, the learner
// simulator, sim — hands it the learner model's estimates, the skill graphs of the learner's
// modules, the lessons and the pool of tasks, the learner's lesson runs, and a summary of what has
// been handed out before, and passes the time in as a number; api keeps every choice the library
// makes. The library estimates no mastery —
// that is the learner model's — gives no verdict, calls no AI model, reads no module package or
// database, and keeps nothing between calls. The rule takes no randomness, so it takes no seed.
//
// The same input gives the same choice and the same reason code, bit for bit, on any build by Go's
// gc compiler of one version of this module and of the learner model — on any architecture it
// knows, WebAssembly included, with any Go no older than go.mod names. The code imports nothing
// but Go's standard library and the learner model, and only the packages its rules allow.
//
// What the library picks and why is in docs/architecture.md, section 3, and in docs/decisions.md
// of this repository.
package engine
