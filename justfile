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
# another version than GOPLS_VERSION in the Dockerfile; then ci-release.
# Check formatting, go vet, the version of gopls and the release's file.
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
    # gopls comes into the devcontainer image from its lock file, .devcontainer/gopls/, which the
    # go kind of spec's repo-template gives (R16; R102 spec), and the image's build stops when the
    # gopls it built does not name GOPLS_VERSION. CI never builds that stage, so the two are held
    # to each other here: a pull request that raises one alone, or a move to a new tag of spec
    # that drops the step from the Dockerfile, fails here.
    lock="$(sed -n 's/^[[:space:]]*\(require[[:space:]]\{1,\}\)\{0,1\}golang\.org\/x\/tools\/gopls[[:space:]]\{1,\}v\([^[:space:]]\{1,\}\).*$/\2/p' \
        .devcontainer/gopls/go.mod)"
    arg="$(sed -n 's/^ARG GOPLS_VERSION=\(.*\)$/\1/p' .devcontainer/Dockerfile)"
    if [ -z "$lock" ] || [ "$arg" != "$lock" ]; then
        echo "ci-lint: .devcontainer/gopls/go.mod names gopls ${lock:-at no version} and" \
            ".devcontainer/Dockerfile has GOPLS_VERSION ${arg:-unset}; raise them together (R16)" >&2
        exit 1
    fi
    # The release's file is the code, the same bytes on every build (R15).
    {{ just_executable() }} ci-release

# The coverage goes into coverage/go.out, one Go profile, which the job coverage of .github's shared
# go.yml sends to Codecov (R48 spec; R12). Go counts statements and leaves the tests themselves out;
# a package that no test runs is in it too, its statements at zero. Every test counts for every
# package of the module, -coverpkg=./...: a test of golden vectors runs the code of other packages,
# and its lines count where they are.
# The checks of the same choice on every build (R10) are the shared go.yml's, with same-bits:
# "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules". "Go / arm64" runs this recipe on
# arm64 hardware, then the targets. The recipes test-targets, fma-diff and rules below run the same
# in the container, and same-bits runs all three; CI calls none of them.
# Run the tests, with the race detector, and write their coverage.
ci-test:
    mkdir -p coverage
    go test -race -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage/go.out ./...

# The checks of the same choice on every build (R10; docs/architecture.md, section 3.7) in the
# container, as the shared go.yml of .github runs them with same-bits, by its own scripts: .github
# at the commit ci.yml pins its go.yml to, so that the container runs what CI runs and nothing here
# is kept in step by hand (the fourth amendment of 2026-11-08 to R10). CI calls none of these
# recipes. ci.yml must call go.yml once, by a commit's SHA. The checkout of that commit is kept in
# the user's cache, by the SHA, with the tool of the rules built in it, and reused: the first run of
# a commit reaches the network. The path of the tool's module holds .github, which go run does not
# take.
# Give the folder of the checkout of .github that ci.yml pins, fetching it once.
[private]
github:
    #!/usr/bin/env bash
    set -euo pipefail
    shas="$(sed -n 's|^[[:space:]]*uses: altnovatrail/\.github/\.github/workflows/go\.yml@\([0-9a-f]\{40\}\).*$|\1|p' .github/workflows/ci.yml)"
    if [ "$(printf '%s' "$shas" | grep -c .)" != 1 ]; then
        echo "same-bits: .github/workflows/ci.yml must call the go.yml of .github once, by a commit's SHA; it calls it ${shas:+by }$(printf '%s' "$shas" | tr '\n' ' ')" >&2
        exit 1
    fi
    dir="${XDG_CACHE_HOME:-$HOME/.cache}/engine-github/$shas"
    if [ ! -f "$dir/.complete" ]; then
        rm -rf "$dir.partial"
        mkdir -p "$(dirname "$dir")"
        git clone --quiet --filter=blob:none --no-checkout https://github.com/altnovatrail/.github.git "$dir.partial"
        git -C "$dir.partial" checkout --quiet "$shas"
        go -C "$dir.partial/tools/coderules" build -o "$dir.partial/coderules" .
        touch "$dir.partial/.complete"
        rm -rf "$dir"
        mv "$dir.partial" "$dir"
    fi
    echo "$dir"

# The script a job of the pinned go.yml runs in the toolchain stage: the heredoc SCRIPT of the job,
# without the indentation of the YAML. A job with no such script stops the recipe with the reason.
[private]
script job:
    #!/usr/bin/env bash
    set -euo pipefail
    file="$({{ just_executable() }} github)/.github/workflows/go.yml"
    awk -v job="{{ job }}" '
        $0 == "  " job ":" { injob = 1; next }
        injob && /^  [a-z0-9-]+:$/ { injob = 0 }
        injob && !inscript && /<<'"'"'SCRIPT'"'"'$/ { inscript = 1; next }
        inscript && /^[[:space:]]*SCRIPT$/ { done = 1; exit }
        inscript {
            if (indent == "") { match($0, /^ */); indent = RLENGTH }
            print substr($0, indent + 1)
        }
        END { if (!done) exit 1 }
    ' "$file" || { echo "same-bits: the job {{ job }} of $file has no script SCRIPT" >&2; exit 1; }

# The tests on the builds where the bits are most at risk — GOAMD64 v3, v1 with
# GODEBUG=cpu.all=off, 386 and js/wasm in the container's Node, and wasip1 built and vetted —, by
# the script of the job targets.
# Run the tests on the builds the same bits are held on, as "Go / targets" does.
test-targets:
    #!/usr/bin/env bash
    set -euo pipefail
    script="$({{ just_executable() }} script targets)"
    bash -c "$script"

# The assembler of each package with fusion and without it on six targets, with the job's own
# control, by the script of the job fma-diff.
# Compare the assembler with fusion and without it, as "Go / fma-diff" does.
fma-diff:
    #!/usr/bin/env bash
    set -euo pipefail
    script="$({{ just_executable() }} script fma-diff)"
    bash -c "$script"

# The rules of code, by the tool tools/coderules of the pinned .github and this repository's
# lists, .github/coderules.json, as the job rules runs it.
# Check the rules of code, as "Go / rules" does.
rules:
    "$({{ just_executable() }} github)/coderules" -lists .github/coderules.json .

# What a task's report names as checked on every target.
# Run test-targets, fma-diff and rules.
same-bits: test-targets fma-diff rules

# The gate of sim (R17; R103 spec): sim's own recipe gate, in its bands mode, on this working tree
# in place of engine — GATE_ENGINE —, and on learner-model at the commit sim's gate.json names. sim
# is taken by the tag and the SHA of its commit that the job gate of ci.yml names, SIM_TAG and
# SIM_COMMIT, once each, so that here and in CI it is the same sim with the same bands: from
# SIM_DIR when it is set, as the job sets it to its checkout of sim, or else fetched once into the
# user's cache, by the SHA. Either way the folder must be at that commit with nothing changed in it,
# and the tag must still be on it. Its status is sim's: 0 green, 1 red or refused, 2 a call it does
# not know. CI runs it in the toolchain stage of sim's devcontainer at that tag; here it runs in
# this repository's, which has the same Go today but is not held to it. sim names engine by its
# HEAD, so a working tree with changes is run but named so on the way.
# Run the gate of sim on this working tree, as "Gate" does.
gate:
    #!/usr/bin/env bash
    set -euo pipefail
    # A value of the job's env, bare or in quotes, as YAML takes either.
    pinned() {
        local values
        values="$(sed -n "s/^[[:space:]]*$1:[[:space:]]*[\"']\{0,1\}\([^[:space:]#\"']\{1,\}\)[\"']\{0,1\}[[:space:]]*\(#.*\)\{0,1\}$/\1/p" .github/workflows/ci.yml)"
        if [ "$(printf '%s' "$values" | grep -c .)" != 1 ]; then
            echo "gate: .github/workflows/ci.yml must name $1 once; it names it ${values:+as }$(printf '%s' "$values" | tr '\n' ' ')" >&2
            exit 1
        fi
        printf '%s' "$values"
    }
    tag="$(pinned SIM_TAG)"
    commit="$(pinned SIM_COMMIT)"
    if ! printf '%s' "$commit" | grep -qx '[0-9a-f]\{40\}'; then
        echo "gate: SIM_COMMIT in .github/workflows/ci.yml is $commit, not the SHA of a commit" >&2
        exit 1
    fi
    repo=https://github.com/altnovatrail/sim
    dir="${SIM_DIR:-}"
    if [ -z "$dir" ]; then
        dir="${XDG_CACHE_HOME:-$HOME/.cache}/engine-sim/$commit"
        # The mark of a whole checkout lies in .git, out of sim's working tree. Each run fills a
        # folder of its own, so that two runs at once do not undo each other: the first to finish
        # puts its checkout in place, and the other drops its own.
        if [ ! -f "$dir/.git/engine-complete" ]; then
            mkdir -p "$(dirname "$dir")"
            partial="$(mktemp -d "$dir.partial.XXXXXX")"
            git -C "$partial" init --quiet
            if ! git -C "$partial" fetch --quiet --depth 1 "$repo" "$commit"; then
                rm -rf "$partial"
                echo "gate: could not fetch $commit, which .github/workflows/ci.yml names, from $repo: git says why above" >&2
                exit 1
            fi
            git -C "$partial" -c advice.detachedHead=false checkout --quiet FETCH_HEAD
            touch "$partial/.git/engine-complete"
            if [ -f "$dir/.git/engine-complete" ]; then
                rm -rf "$partial"
            else
                rm -rf "$dir"
                mv "$partial" "$dir"
            fi
        fi
    fi
    head="$(git -C "$dir" rev-parse HEAD)"
    if [ "$head" != "$commit" ]; then
        echo "gate: sim in $dir is at $head, and .github/workflows/ci.yml names $commit" >&2
        exit 1
    fi
    if [ -n "$(git -C "$dir" status --porcelain)" ]; then
        echo "gate: sim in $dir has changes of its own, so its bands may not be those of $commit:" >&2
        git -C "$dir" status --short >&2
        exit 1
    fi
    # The commit of the tag: the peeled ref of an annotated tag, or else the tag's own.
    tagged="$(git ls-remote "$repo" "refs/tags/$tag" "refs/tags/$tag^{}" \
        | awk -v ref="refs/tags/$tag" '$2 == ref "^{}" { peeled = $1 } $2 == ref { own = $1 } END { print (peeled != "" ? peeled : own) }')"
    if [ "$tagged" != "$commit" ]; then
        echo "gate: sim $tag is on ${tagged:-no commit}, and .github/workflows/ci.yml names $commit" >&2
        exit 1
    fi
    if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
        echo "gate: engine has changes not committed; sim runs them but names them $(git rev-parse --short HEAD)" >&2
    fi
    echo "gate: sim $tag at $commit, engine from $PWD" >&2
    GATE_ENGINE="$PWD" {{ just_executable() }} -f "$dir/justfile" gate

# The release's one file (R15): engine-<version>.tar.gz, a folder of that name with LICENSE,
# go.mod, go.sum and the Go code of the module but its tests and testdata/ — what a consumer builds,
# the paths that give a release. The tar is written the same by every build: names in order, no
# times or owners, mode 0644, and gzip with no name or time. .github's shared release runs it in a
# tree git archive made, with no git to ask, so the files come from the tree. Go code takes the
# module by the tag; this file binds the tag to the signed release (R77 spec).
# Build the release's file into dist/.
release-files version:
    #!/usr/bin/env bash
    set -euo pipefail
    version={{ quote(version) }}
    if ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
        echo "release-files: the version is v<major>.<minor>.<patch>, as v0.1.0, not ${version}" >&2
        exit 1
    fi
    name="engine-${version}"
    rm -rf dist
    mkdir dist
    (printf '%s\n' LICENSE go.mod go.sum
        find . -name testdata -prune -o -path './.*' -prune -o -path ./spec -prune -o \
            -type f -name '*.go' ! -name '*_test.go' -print | sed 's|^\./||') |
        LC_ALL=C sort |
        tar --create --verbatim-files-from --files-from=- --format=posix \
            --pax-option=exthdr.name=%d/PaxHeaders/%f,delete=atime,delete=ctime \
            --transform="s,^,${name}/," --mtime=@0 --owner=0 --group=0 --numeric-owner --mode=0644 |
        gzip -9n > "dist/${name}.tar.gz"

# Builds the release's file twice, a second apart, and fails unless the two are the same bytes;
# unless the archive holds LICENSE, go.mod, go.sum and every Go file but tests and testdata/ git
# tracks or would, as they are in the tree, and nothing else; and unless the unpacked archive builds.
# So the file a release signs is the code of its commit, and anyone can build it again and
# compare. It refuses to run beside a dist/ of one's own, which release-files would replace.
# ci-lint runs it, so "Go / lint" holds it (R15).
# Check that the release's file is the code, and the same bytes on every build.
ci-release:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -e dist ]; then
        echo "ci-release: dist/ is there already, and release-files would replace it: move it away first" >&2
        exit 1
    fi
    version=v0.0.0-ci
    name="engine-${version}"
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp" dist' EXIT
    {{ just_executable() }} release-files "$version"
    mv dist "$tmp/first"
    sleep 1
    {{ just_executable() }} release-files "$version"
    if [ "$(ls dist)" != "${name}.tar.gz" ] || ! cmp -s "$tmp/first/${name}.tar.gz" "dist/${name}.tar.gz"; then
        echo "ci-release: two builds of the release's files differ, or dist/ holds another file:" >&2
        ls -l "$tmp/first" dist >&2
        exit 1
    fi
    mkdir "$tmp/unpacked"
    tar -xzf "dist/${name}.tar.gz" -C "$tmp/unpacked"
    got="$(cd "$tmp/unpacked/$name" && find . -type f | sed 's|^\./||' | LC_ALL=C sort)"
    # The files git tracks or would, as the gofmt step of ci-lint takes them, but those deleted from
    # the tree: a new file before git add, or one removed before git rm, is the tree's as the release
    # sees it.
    want="$(git ls-files --cached --others --exclude-standard -- LICENSE go.mod go.sum '*.go' \
        ':(exclude)*_test.go' ':(exclude,glob)**/testdata/**' |
        while IFS= read -r file; do if [ -e "$file" ]; then printf '%s\n' "$file"; fi; done | LC_ALL=C sort -u)"
    if [ "$got" != "$want" ]; then
        echo "ci-release: the archive holds" >&2
        echo "$got" >&2
        echo "and git tracks" >&2
        echo "$want" >&2
        exit 1
    fi
    while IFS= read -r file; do
        if ! cmp -s "$file" "$tmp/unpacked/$name/$file"; then
            echo "ci-release: ${file} in the archive is not the tree's" >&2
            exit 1
        fi
    done <<< "$want"
    go -C "$tmp/unpacked/$name" build ./...
    echo "ci-release: ${name}.tar.gz is the code, the same bytes on two builds" >&2

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
