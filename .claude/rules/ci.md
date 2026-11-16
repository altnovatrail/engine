---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - ".claude/settings.json"
  - "justfile"
---

# CI, the release and the toolchain stage

- **The checks** are this repository's own recipes, in the `justfile`, and run in the toolchain stage of `.devcontainer/Dockerfile`, which is spec's (R12). CI calls NovaTrail's shared workflows of `.github`, pinned by SHA with the version beside it, and Dependabot moves both: "Secrets", "DCO", "Documents" and "Go" — with `same-bits`, "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" (R10; R12). Its own jobs are "Contracts", which checks out the `spec` submodule for `just ci-contracts`, and "Gate".
- **The image of the stage** comes from the cache of GitHub Actions, by the hash of the Dockerfile, and is built only on a miss: in the shared jobs, and in "Contracts" and "Gate" by the same three steps as the shared `go.yml` of the pinned release, so a move to a release of `.github` that changes them changes these too (R28). `actions/cache` in `ci.yml` keeps the version of the pinned `go.yml`: a Dependabot pull request that raises it in `ci.yml` alone waits for a release of `.github` that raises it too. Dependabot raises the references to `.github` in one pull request, the group `shared-workflows`.
- **The gate.** "Gate" runs `just gate` in the toolchain stage of `sim`'s devcontainer, by the tag and the SHA of `sim` that `.github/workflows/ci.yml` names, on `ubuntu-24.04-arm`, where its numbers are the same to the bit as on amd64 (R17 and its second 2026-11-16 amendment). The tag and the SHA are inputs of a checkout, so Dependabot does not move them: a task raises them, by the order of releases of repositories 4 spec.
- **The release** comes out by itself through the shared `release.yml` of `.github`, on the line `0.1`, for a commit of `main` that changes a path of its `paths`; `just release-files <version>` builds its one file, and `just ci-release` holds it to the code (R15). Nobody tags by hand.
- **The container is spec's.** `.devcontainer/` and `.claude/settings.json` are what `just -f spec/justfile new-repo-draft engine go 3.8 <tag>` builds for the kind `go` at the pinned tag: a move to a new tag compares them byte for byte. gopls comes from its lock file, `.devcontainer/gopls/`, by the kind `go` (R16, which replaced R13), and the plugin `gopls-lsp`, which `CLAUDE_CODE_PLUGIN_DIRS` of `.devcontainer/devcontainer.json` loads, runs it (R143 spec; the 2026-11-16 amendment to R24). No decision records a difference of the container from spec's so far (CLAUDE.md, Environment).
