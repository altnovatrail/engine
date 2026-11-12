---
paths:
  - ".github/**"
  - ".devcontainer/**"
  - ".claude/settings.json"
  - "justfile"
---

# CI, the release and the toolchain stage

- **The checks** are this repository's own recipes, in the `justfile`, and run in the toolchain stage of `.devcontainer/Dockerfile`, which is spec's (R12). CI calls NovaTrail's shared workflows of `.github`, pinned by SHA with the version beside it, and Dependabot moves both: "Secrets", "DCO", "Documents" and "Go" — with `same-bits`, "Go / targets", "Go / fma-diff", "Go / arm64" and "Go / rules" (R10; R12). Its own jobs are "Contracts", which checks out the `spec` submodule for `just ci-contracts`, and "Gate".
- **The gate.** "Gate" runs `just gate` in the toolchain stage of `sim`'s devcontainer, by the tag and the SHA of `sim` that `.github/workflows/ci.yml` names (R17). They are inputs of a checkout, so Dependabot does not move them: a task raises them, by the order of releases of repositories 4 spec.
- **The release** comes out by itself through the shared `release.yml` of `.github`, on the line `0.1`, for a commit of `main` that changes a path of its `paths`; `just release-files <version>` builds its one file, and `just ci-release` holds it to the code (R15). Nobody tags by hand.
- **The container is spec's.** `.devcontainer/` and `.claude/settings.json` are what `just -f spec/justfile new-repo-draft engine go 3.8 <tag>` builds for the kind `go` at the pinned tag: a move to a new tag compares them byte for byte. gopls comes from its lock file, `.devcontainer/gopls/`, by the kind `go` (R16, which replaced R13), and the plugin `gopls-lsp` of `.claude/settings.json` runs it (R143 spec). No decision records a difference of the container from spec's so far (CLAUDE.md, Environment).
