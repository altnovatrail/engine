# engine

The NovaTrail engine: it picks what a learner does next — a lesson, a task or an item, on a new skill or on one due for practice again — and gives each choice a reason code that the adult sees. It is a Go library of pure computation that works by a rule, with no AI model: it is handed the skill graphs, the task pool and the learner model's estimates, time and the seed come from outside, and the same input always gives the same choice. From stage 2 it runs inside the platform's service, `api`; at stage 1 the learner simulator, `sim`, runs it.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.1.0), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and sim; it comes together with learner-model and sim, since the simulator needs both the learner model and the rule.

## Status

Detailing. The documents come first: the first epic of [RUN.md](RUN.md) decides what the component's card left open, and the code follows in the order RUN.md sets. There is no code here yet: the rule starts from the one by which MathTrail's standalone chose the next task, and the epic decides what it keeps and what it changes.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, and what is still to decide.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.1.0), checked out here as a git submodule at `v0.1.0`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
