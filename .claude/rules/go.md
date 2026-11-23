---
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
  - "testdata/**"
  - ".github/coderules.json"
---

# The Go code of engine

The library is one Go module, `github.com/altnovatrail/engine`, at the root: the rule's code and its tests. `testdata/` holds the golden vectors of the choice, `testdata/choose/*.txt`, and the files of `modules` releases in `testdata/modules` that `release-verify` accepted (R10). A change of code comes with its tests, and Codecov waits for `codecov/patch`, at least 90% of the changed lines covered (R48 spec).

- `just ci-lint`, `just ci-test`, `just ci-vuln` and `just ci-licenses` — `gofmt` and `go vet`, the version of gopls against its lock file and the release's file; the tests with `-race` and their coverage in `coverage/go.out`; govulncheck; go-licenses with `ALLOWED_LICENSES`. A task that touches Go code ends with them green.
- `just same-bits` — `test-targets`, `fma-diff` and `rules`: the same bits on every build, WebAssembly too, by the rules of code of CLAUDE.md, "The component", and the lists of `.github/coderules.json`. A task that touches the rule's arithmetic or its imports ends with it green.
- `just ci-contracts` — `ReasonCode`, `ChoiceKind` and `RunEnd` against `cause.proto` in the `spec` submodule, at its tag. A move to a new tag of spec and a change of a code end with it green.
- `just gate` — the gate of `sim` on the working tree, in one run, as the jobs "Gate / part k" and the check "Gate" run it in parts by `just gate-part` and `just gate-check`, with the same numbers to the bit (R17): every change after the baseline `v0.1.0` passes it. A change that keeps the choice keeps the golden vectors and all the gate's numbers to the bit; the vectors are written anew only where the choice changes on purpose, and such a change is judged by `sim` before it is released (R103 spec).
- `learner-model` enters `go.mod` only by a task that has `just -f spec/justfile release-verify` accept its release's file (R09; R77 spec): Dependabot does not raise it.
