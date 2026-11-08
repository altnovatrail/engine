# engine

[![CI](https://github.com/altnovatrail/engine/actions/workflows/ci.yml/badge.svg)](https://github.com/altnovatrail/engine/actions/workflows/ci.yml)
[![CodeQL](https://github.com/altnovatrail/engine/actions/workflows/codeql.yml/badge.svg)](https://github.com/altnovatrail/engine/actions/workflows/codeql.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/altnovatrail/engine)](go.mod)
[![Coverage](https://codecov.io/gh/altnovatrail/engine/branch/main/graph/badge.svg)](https://codecov.io/gh/altnovatrail/engine)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/altnovatrail/engine/badge)](https://scorecard.dev/viewer/?uri=github.com/altnovatrail/engine)
[![License](https://img.shields.io/github/license/altnovatrail/engine)](LICENSE)

[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_engine&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=altnovatrail_engine)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_engine&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_engine)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_engine&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_engine)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_engine&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=altnovatrail_engine)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=altnovatrail_engine&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=altnovatrail_engine)

The NovaTrail engine: it picks what a learner does next — a lesson or a task, on a new skill or on one due for practice again — and gives each choice a reason code that the adult sees. It is a Go library of pure computation that works by a rule, with no AI model: it is handed the skill graphs, the lessons and the task pool and the learner model's estimates, time comes from outside, it takes no seed because the rule needs no randomness, and the same input always gives the same choice, on any build. From stage 2 it runs inside the platform's service, `api`; at stage 1 the learner simulator, `sim`, runs it.

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.9.8), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and sim; it comes together with learner-model and sim, since the simulator needs both the learner model and the rule.

## Status

Detailing is done: the first epic, archived in [docs/plan/E1.md](docs/plan/E1.md), decided everything the component's card left open, and spec has recorded it. The current epic of [RUN.md](RUN.md), E2, builds the first version of the rule: so far there is the Go module, with no rule in it yet, and the CI its code will pass. The rule starts from the one by which MathTrail's standalone chose the next task.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what the card left to decide — all of it decided — and how it is built.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.9.8), checked out here as a git submodule at `v0.9.8`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
