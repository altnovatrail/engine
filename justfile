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
# Run the tests, with the race detector, and write their coverage.
ci-test:
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...

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
