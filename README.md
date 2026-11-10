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

A component of [NovaTrail](https://github.com/altnovatrail/spec/tree/v0.14.0), a universal open platform for adaptive learning, interactive lessons and AI tutoring. It belongs to stage 1, Format and core, with spec, modules, module-kit, learner-model and sim; it comes together with learner-model and sim, since the simulator needs both the learner model and the rule.

## Status

Detailing is done: the first epic, archived in [docs/plan/E1.md](docs/plan/E1.md), decided everything the component's card left open, and spec has recorded it. The second epic, E2, archived in [docs/plan/E2.md](docs/plan/E2.md), built the first version of the rule, which starts from the one by which MathTrail's standalone chose the next task: lessons and tasks with their reason codes over several modules, with golden vectors of the choice on the graphs of the first slice, and the same bits on every build. It came out as `v0.1.0`, the baseline before sim's first tag, and its gate now holds every change after it. The current epic of [RUN.md](RUN.md), E3, works on the first candidates of the rule, which sim judges: why the base gives tasks outside the corridor of success, and variants that close the gap.

## The gate

Every change after the baseline passes the gate of [sim](https://github.com/altnovatrail/sim), the learner simulator: the job "Gate" of CI runs sim, by the tag and the SHA of its commit that `.github/workflows/ci.yml` names, on the commit's code in place of engine, and a pull request whose numbers leave their bands, or whose choices break an invariant, does not merge ([R17](docs/decisions.md)). `just gate` runs the same gate on the working tree, with sim fetched once into the user's cache.

## Releases

Releases come out by themselves: each commit of `main` whose CI has passed and which changes the code, `go.mod`, `go.sum` or what builds the release's file gets the next patch of the line `0.1` — the first, `v0.1.0`, came out on 2026-11-08 — with its tag, by .github's shared release workflow ([R15](docs/decisions.md)). Nobody tags by hand, and a release is immutable.

- **Go code** takes the library as a Go module, by the tag: `go get github.com/altnovatrail/engine@v0.1.0`. `go.sum` and the Go checksum database hold the version from then on.
- **The release's one file** is the code a consumer builds: `engine-<tag>.tar.gz`, with `LICENSE`, `go.mod`, `go.sum` and the Go files of the module but its tests and `testdata/`, signed without a key. Verify it where it enters NovaTrail with spec's recipe, `just -f spec/justfile release-verify engine <tag> engine-<tag>.tar.gz`, with its bundle `<file>.sigstore.json` beside it. `just release-files <tag>`, in the toolchain stage of the tag's `.devcontainer/Dockerfile`, builds the same bytes from the tag's tree.

## Documents

The documents are in Russian.

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links, what the card left to decide — all of it decided — and how it is built.
- [docs/decisions.md](docs/decisions.md) — the decision log: what was decided, why, and what was rejected.
- [RUN.md](RUN.md) — the work plan, in epics.

The working rules, for people and coding agents alike, are in [CLAUDE.md](CLAUDE.md) and [AGENTS.md](AGENTS.md). The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec/tree/v0.14.0), checked out here as a git submodule at `v0.14.0`: clone this repository with `git clone --recurse-submodules`.

## License

[Apache-2.0](LICENSE).
