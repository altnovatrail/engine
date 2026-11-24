# engine — the documents

What the [README](../README.md) of the repository leaves out, which says only what the component is and where it stands in NovaTrail: the index of the documents, for people and for the session of Claude Code that changes `engine`, the releases, the gate and how to contribute.

## Documents

- [architecture.md](architecture.md) — the component: what it does and does not do, its data and links, what the card left to decide — all of it decided — and how it is built.
- [decisions.md](decisions.md) — the decision log: what was decided, why, and what was rejected.
- [sources.md](sources.md) — the records of the facts no run checks that the documents rest on, each with its address, the date read and a quote.
- [measurements.md](measurements.md) — the results of measurements: what was measured, how, on what, and the numbers.
- `plan/` — the archives of the done epics, in Russian: a task's body and its gist sentence.
- [RUN.md](../RUN.md) — the work plan, in epics; `run/` — its other tables: [the epics for later](../run/EPICS.md), [the open questions](../run/QUESTIONS.md), [the technical debt](../run/DEBT.md) and [the tasks outside an epic](../run/TASKS.md).

## Releases

Releases come out by themselves: each commit of `main` whose CI has passed and which changes the code, `go.mod`, `go.sum` or what builds the release's file gets the next patch of its line, `RELEASE_LINE` of [release.yml](../.github/workflows/release.yml), with its tag, by .github's shared release workflow. Nobody tags by hand, and a release is immutable.

- **Go code** takes the library as a Go module, by the tag: `go get github.com/altnovatrail/engine@<tag>`. `go.sum` and the Go checksum database hold the version from then on.
- **The release's one file** is the code a consumer builds: `engine-<tag>.tar.gz`, with `LICENSE`, `go.mod`, `go.sum` and the Go files of the module but its tests and `testdata/`, signed without a key. Verify it where it enters NovaTrail with spec's recipe, `just -f spec/justfile release-verify engine <tag> engine-<tag>.tar.gz`, with its bundle `<file>.sigstore.json` beside it; `gh release download <tag> -R altnovatrail/engine` fetches both. `just release-files <tag>`, in the toolchain stage of the tag's `.devcontainer/Dockerfile`, builds the same bytes from the tag's tree.

## The gate

Every change after the baseline passes the gate of [sim](https://github.com/altnovatrail/sim), the learner simulator: CI runs sim, by the tag and the SHA of its commit that `.github/workflows/ci.yml` names, on the commit's code in place of engine, in four jobs "Gate / part k" and the check "Gate" after them, and a pull request whose numbers leave their bands, or whose choices break an invariant, does not merge. `just gate` runs the same gate on the working tree in one run, with sim fetched once into the user's cache, and `just gate-part <k> <n> <folder>` and `just gate-check <folder>` run it as the jobs do.

## Contributing

The working rules, for people and coding agents alike, are in [CLAUDE.md](../CLAUDE.md) and [AGENTS.md](../AGENTS.md): the work happens in the devcontainer, one task of [RUN.md](../RUN.md) at a time. The platform's shared vocabulary, formats, contracts and decisions are in [spec](https://github.com/altnovatrail/spec), checked out here as a git submodule: clone this repository with `git clone --recurse-submodules`. The code of conduct, the contributing guide and the security policy are the organization's, in its repository `.github`, and GitHub shows them here.
