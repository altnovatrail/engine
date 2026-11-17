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

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring.

## The gate

Every change after the first version passes the gate of [sim](https://github.com/altnovatrail/sim), the learner simulator, in CI and locally with `just gate`; how it runs — [docs/README.md](docs/README.md#the-gate).

## Releases

Each commit of `main` that changes the code and passes CI comes out as the next patch: a Go module by its tag and a signed archive of its code. Details — [docs/README.md](docs/README.md#releases).

```sh
go get github.com/altnovatrail/engine@<tag>
gh release download <tag> -R altnovatrail/engine
just -f spec/justfile release-verify engine <tag> engine-<tag>.tar.gz
```

## Documents

- [docs/](docs/README.md) — the component, its decisions, sources and measurements, and how to contribute.
- [RUN.md](RUN.md) — the work plan; its other tables are in `run/`.
