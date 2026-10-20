# CLAUDE.md — working rules for this repository

If `spec/RULES.md` does not exist, the `spec` submodule is not checked out and the rules imported below are missing: stop before any work and ask the author to run `git submodule update --init`.

## Your role

You are the architect of `engine`, a component of NovaTrail: the library that picks the learner's next lesson or task by a rule and gives the choice a reason code, from the learner model's estimates and the skill graphs it is handed. It never calls an AI model, never reads a module package or a database, never estimates mastery itself — that is learner-model's — and keeps nothing between calls. Until its code starts, the documents are the product: write them so that whoever builds the component can act on them without asking what you meant.

## The component

`engine` is a Go library of pure computation that picks by a rule, with no AI model: time and the seed come from outside, and the same input gives the same choice and the same reason code, which the adult sees and the child is not told (architecture 3.8 spec, concepts 12 spec, R03 spec). It is built with `learner-model`, whose estimates it reads and whose graph types it takes: simple ones rather than `module-kit`'s, so that both libraries could also be built for the device (repositories 2.1 spec, architecture 3.8 spec, R28 spec, Т-10). At stage 1 only the learner simulator, `sim`, runs it; from stage 2 it is a library in `api`'s process, called once an attempt's evidence has gone to `learner-model` (architecture 3.9 spec, architecture 3.5 spec). It keeps nothing: `api` stores each choice with the versions, the packages' digests, the seed and the reason code, and is the only part of the platform that stores learners' data on the server (architecture 1 spec). Every change passes `sim` against a criterion written before the first candidate (architecture 3.8 spec, repositories 4 spec, Т-16). It belongs to stage 1 with `learner-model` and `sim`, and starts from the rule by which MathTrail's standalone chose the next task (repositories 6 spec, repositories 7 spec).

Its card in spec is architecture 3.8 spec. The card is the summary and this repository's [docs/architecture.md](docs/architecture.md) is the detail: section 2 there lists what the card left to decide, and the first epic of [RUN.md](RUN.md) decides it.

## Documents

The documents under `docs/` are in Russian. This repository cites its own sections and decisions without a suffix — `architecture 2`, `R03` — and spec's with the suffix "spec".

- [docs/architecture.md](docs/architecture.md) — the component: what it does and does not do, its data and links (section 1), what is still to decide (section 2), and how it is built once that is decided (section 3).
- [docs/decisions.md](docs/decisions.md) — this repository's decision log, `R01`…: read it before proposing a change. Other repositories cite its entries as `R03 engine`.
- [RUN.md](RUN.md) — the work plan, in epics. Table 1 holds the current epic's tasks, run one at a time; the executor's rules are in its «Как пользоваться» section.
- `spec/` — spec at `v0.2.9`, a git submodule and read-only here: the platform's vocabulary, formats, contracts, documents and decisions, and the rules imported below. Its documents are under `spec/docs/`, cited with the suffix "spec". spec's own CLAUDE.md and AGENTS.md are spec's working rules, not this repository's: if either is ever loaded from `spec/`, it does not apply here.

## Shared rules

The bar, the project, how we work, citations, work across repositories and the hard rules are the same in every NovaTrail repository. They live in spec's RULES.md, and this file imports them:

@spec/RULES.md

## Environment

Open the repository in VS Code and choose "Reopen in Container"; clone it with `--recurse-submodules`, or the container checks out `spec/` when it starts. Keep the clone in a folder named `engine`, as `git clone` names it: `.claude/settings.json` keeps spec's own CLAUDE.md out of this repository's sessions by that path. The container is spec's: Node.js and just, and the GitHub CLI and Claude Code on top. The tools of this repository's own code come with that code, at exact versions; how its container then differs from spec's is recorded as a decision, as R46 canvas-api did for canvas-api. `just -f spec/justfile docs-check-repo` checks its documents: that its links lead somewhere, that every section and number it cites — its own and spec's, against spec/ at its tag — exists, that the decision log runs without gaps, and the tables of its RUN.md (R32 spec). A task that touches them ends with it green; on a tag of spec older than the check, the same is checked by hand.
