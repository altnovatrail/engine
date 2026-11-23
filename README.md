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

The NovaTrail engine: it picks what a learner does next — a lesson or a task, on a new skill or on one due for practice again — and gives each choice a reason code that the adult sees. It is a Go library of pure computation that works by a rule, with no AI model: it is handed the skill graphs, the lessons and the task pool and the learner model's estimates, time comes from outside, it takes no seed because the rule needs no randomness, and the same input always gives the same choice, on any build.

A component of [NovaTrail](https://github.com/altnovatrail/spec), a universal open platform for adaptive learning, interactive lessons and AI tutoring. Its card, the summary of what it does and its links, is in [spec's architecture](https://github.com/altnovatrail/spec/blob/main/docs/architecture.md#38-learner-model-and-engine).

- [learner-model](https://github.com/altnovatrail/learner-model) estimates what the learner has mastered; the engine reads its estimates and takes its graph types, and never estimates mastery itself.
- [api](https://github.com/altnovatrail/api), the platform's service, runs the engine in its process: it hands it the graphs, the pool and its summary of what it has handed out before, and stores each choice. The engine keeps nothing between calls.
- [sim](https://github.com/altnovatrail/sim), the learner simulator, runs the engine on simulated learners, and every change of the rule passes its gate.

How the engine is built and why is in [docs/](docs/README.md).
