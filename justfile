# engine's own recipes. The ci- ones are what .github's shared go.yml runs, one job each, in the
# toolchain stage of .devcontainer/Dockerfile (R16 .github; R12 in docs/decisions.md): their names
# are the same in every Go repository of NovaTrail, and what they run is this repository's own.

# The tools only the checks need, at exact versions. Each is a Go program run straight from its
# module, as learner-model's justfile runs them, so nothing is installed for it and go.mod gains no
# require: the Go module proxy serves that version and the checksum database vouches for it.
# Dependabot does not see them, so they are raised by hand.
GOVULNCHECK := "golang.org/x/vuln/cmd/govulncheck@v1.9.0"
GO_LICENSES := "github.com/google/go-licenses/v2@v2.0.1"

# List the recipes.
default:
    @just --list

# Fails on code gofmt would rewrite, on what go vet reports, and on a gopls lock file that names
# another version than GOPLS_VERSION in the Dockerfile.
# Check formatting, go vet and the version of gopls.
ci-lint:
    #!/usr/bin/env bash
    set -euo pipefail
    # The Go files git tracks or would, outside every testdata/, the root's too, which glob magic
    # takes and a bare **/ would not: spec/ is a submodule, another repository's, and git lists
    # none of its files. A file deleted from the working tree but still in the index is skipped,
    # so that gofmt does not fail on it with no word of formatting.
    unformatted="$(git ls-files -z --cached --others --exclude-standard -- '*.go' ':(exclude,glob)**/testdata/**' \
        | while IFS= read -r -d '' file; do if [ -e "$file" ]; then printf '%s\0' "$file"; fi; done \
        | xargs -0 -r gofmt -l)"
    if [ -n "$unformatted" ]; then
        echo "ci-lint: gofmt would rewrite these files; run gofmt -w on them:" >&2
        echo "$unformatted" >&2
        exit 1
    fi
    go vet ./...
    # gopls comes into the devcontainer image from its lock file, .devcontainer/gopls/, and the
    # image's build stops when the gopls it built does not name GOPLS_VERSION (R13). CI never
    # builds that stage, so the two are held to each other here: a pull request that raises one
    # alone, or a move to a new tag of spec that drops the step from the Dockerfile, fails here.
    lock="$(sed -n 's/^[[:space:]]*\(require[[:space:]]\{1,\}\)\{0,1\}golang\.org\/x\/tools\/gopls[[:space:]]\{1,\}v\([^[:space:]]\{1,\}\).*$/\2/p' \
        .devcontainer/gopls/go.mod)"
    arg="$(sed -n 's/^ARG GOPLS_VERSION=\(.*\)$/\1/p' .devcontainer/Dockerfile)"
    if [ -z "$lock" ] || [ "$arg" != "$lock" ]; then
        echo "ci-lint: .devcontainer/gopls/go.mod names gopls ${lock:-at no version} and" \
            ".devcontainer/Dockerfile has GOPLS_VERSION ${arg:-unset}; raise them together (R13)" >&2
        exit 1
    fi

# The coverage goes into coverage/go.out, one Go profile, which the job coverage of .github's shared
# go.yml sends to Codecov (R48 spec; R12). Go counts statements and leaves the tests themselves out;
# a package that no test runs is in it too, its statements at zero. Every test counts for every
# package of the module, -coverpkg=./...: a test of golden vectors runs the code of other packages,
# and its lines count where they are.
# Then come the checks of the same choice on every build (R10): the tests on the builds where the
# bits are most at risk, and the comparison of the assembler with and without fusion.
# Run the tests, with the race detector, and write their coverage; then the same-bits checks.
ci-test: && test-targets fma-diff
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...

# The same tests, without the race detector and coverage, which Go does not have on 386 and wasm
# (docs/architecture.md, section 3.7): on amd64 with GOAMD64 v3, where the compiler fuses an
# addition; with v1 and GODEBUG=cpu.all=off, with no processor extension for the runtime and the
# standard library to pick their path by; on 386; and in js/wasm, in the container's Node. A host
# other than amd64 — the arm64 job of ci.yml — runs no amd64 build and none for 386, and runs its
# own build with cpu.all=off. wasip1 is built and vetted alone: there is no WASI runtime in the
# container.
# Run the tests on the builds the same bits are held on.
test-targets:
    #!/usr/bin/env bash
    set -euo pipefail
    run() { echo "== $*" >&2; env "$@"; }
    if [ "$(go env GOHOSTARCH)" = amd64 ]; then
        run GOAMD64=v3 go test -count=1 ./...
        run GOAMD64=v1 GODEBUG=cpu.all=off go test -count=1 ./...
        run GOARCH=386 go test -count=1 ./...
    else
        run GODEBUG=cpu.all=off go test -count=1 ./...
    fi
    run GOOS=js GOARCH=wasm go test -count=1 -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./...
    run GOOS=wasip1 GOARCH=wasm go build ./...
    run GOOS=wasip1 GOARCH=wasm go vet ./...

# Each package of the module is compiled for every target whose compiler fuses a multiplication
# with an addition or a subtraction, twice: as it is, and with -d=fmahash=n, which turns every rule
# of fusion off. The assembler must be the same, and the second build must report no place where a
# rule would have fired, "fmahash triggered". The control, testdata/fma, must differ on every target,
# or the check has stopped seeing fusion (docs/architecture.md, section 3.7). The packages are built
# one at a time, so that their output comes in one order. The output does not depend on the host,
# so the arm64 job repeats it, as the same ci-test.
# Compare the assembler with fusion and without it, on every target that fuses.
fma-diff:
    #!/usr/bin/env bash
    set -euo pipefail
    # A build that fails says why, rather than stopping the script with nothing printed.
    asm() {
        local out
        if ! out="$(go build -gcflags="-S $1" "$2" 2>&1)"; then
            echo "fma-diff: $GOARCH: $2 does not build:" >&2
            echo "$out" >&2
            return 1
        fi
        echo "$out"
    }
    # The packages with code: one of tests alone has nothing to build.
    pkgs="$(go list -f '{{{{if .GoFiles}}{{{{.ImportPath}}{{{{end}}' ./...)"
    status=0
    for target in arm64 ppc64le s390x riscv64 loong64 amd64/v3; do
        export GOOS=linux GOARCH="${target%%/*}" GOAMD64=v1
        if [ "$target" = amd64/v3 ]; then GOAMD64=v3; fi
        ok=1
        for pkg in $pkgs; do
            fused="$(asm "" "$pkg")"
            plain="$(asm -d=fmahash=n "$pkg")"
            if [ "$fused" != "$plain" ] || grep -q 'fmahash triggered' <<<"$plain"; then
                echo "fma-diff: $target: the compiler fuses in $pkg; put float64() round the * or / (R10):" >&2
                diff <(echo "$fused") <(echo "$plain") | head -n 40 >&2 || true
                ok=0
            fi
        done
        if [ "$(asm "" ./testdata/fma)" = "$(asm -d=fmahash=n ./testdata/fma)" ]; then
            echo "fma-diff: $target: the control in testdata/fma was not caught: the check sees no fusion" >&2
            ok=0
        fi
        if [ "$ok" = 1 ]; then echo "fma-diff: $target: ok" >&2; else status=1; fi
    done
    exit "$status"

# Fail on a known vulnerability in code this module reaches.
ci-vuln:
    go run {{ GOVULNCHECK }} ./...

# Fails if a package this module is built from carries a license outside ALLOWED_LICENSES, SPDX
# identifiers separated by commas, which .github's go.yml sets.
# Check the licenses of what the module is built from against ALLOWED_LICENSES.
ci-licenses:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "${ALLOWED_LICENSES:-}" ]; then
        echo "ci-licenses: ALLOWED_LICENSES is not set. CI takes it from .github's go.yml; to run this here, set it to the list there." >&2
        exit 1
    fi
    go run {{ GO_LICENSES }} check ./... --allowed_licenses="$ALLOWED_LICENSES"

# The test that holds ReasonCode, ChoiceKind and RunEnd to cause.proto in the spec submodule
# (docs/architecture.md, section 3.6): the job Contracts of ci.yml runs it with the submodule
# checked out, where a missing file fails rather than skips, as it does in the other runs.
# Check the engine's codes against the contracts of spec, in its submodule.
ci-contracts:
    ENGINE_CONTRACTS=1 go test -count=1 -run '^TestContracts$' ./...
