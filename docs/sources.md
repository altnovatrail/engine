# engine — sources

> The facts of the outside world the repository's documents under `docs/` rest on and no run of ours checks: a law, a regulator, a standard's normative text, a provider's terms and policies, a price, a license, a page that changes in place with no version, a text outside the repositories. Here are their records — the address, the date it was read and a quote in the source's language — and [architecture.md](architecture.md) and the decision log hold the claim and a citation of its one record by its number, `sources 1.1`, not a chain. A record's heading names its topic, not a task.
>
> No record is needed for a library's documentation or code at a version or commit — a link beside the claim —, for a behaviour a test or a measurement checks — they are its evidence —, or for spec's tags and the files of NovaTrail's and MathTrail's repositories — git, at a commit.
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. A record's number never changes and is never given again: a new one goes at the end of its topic, a new topic at the end of the file, and a record removed leaves a gap. A source read again with the same text changes the record's date; with another text, the claims and decisions on it are rewritten in the same change and the record holds the new reading alone — where a decision waits on the author, an open question, and both readings stand until it is closed. The model is `spec/docs/sources.md`. The results of measurements are not here but in [measurements.md](measurements.md).

## 1. MathTrail: standalone

No record for standalone's files: they are MathTrail's git, at the commit a claim links beside it.

## 2. Research and other programs

### 2.1. A corridor near 85% correct

Wilson R. C., Shenhav A., Straccia M., Cohen J. D. The Eighty Five Percent Rule for optimal learning. *Nature Communications* 10, 4646 (2019), [doi:10.1038/s41467-019-12552-4](https://doi.org/10.1038/s41467-019-12552-4), via [PMC6831579](https://pmc.ncbi.nlm.nih.gov/articles/PMC6831579/), read 2026-11-07: "the optimal error rate for training is around 15.87% or, conversely, that the optimal training accuracy is about 85%". This is the upper bound of standalone's corridor. The lower one, about 0.75, [standalone's SPEC](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md) at c6cbeb2 takes from Math Garden ("Math Garden keeps success near 0.75"); the primary source — the article of Math Garden — could not be read: the publisher answers 403, and it is not stated here as fact.

### 2.2. Repetitions and the new in Anki

[The Anki manual, Deck Options](https://docs.ankiweb.net/deck-options.html), read 2026-11-07. The review limit limits the new too, with a backlog the advice is not to introduce the new, and the order of the new and the repetitions is a setting: "By default, the review limit also applies to new cards, and no new cards will be shown when the review limit has been reached"; "When this limit is reached, Anki will not show any more review cards for the day, even if there are more waiting"; the order of new cards and repetitions — "Whether new cards should be mixed in with review cards, shown before them, or shown after them"; "If you have a backlog of overdue review cards, it is recommended that you stop introducing new cards until you catch up with that backlog".

## 3. `modules`

No record for `modules`' files and releases, nor for its packages in `testdata/modules`: they are git, at the commit or tag a claim links beside it, and the release's index and packages are held by the digest that `release_internal_test.go` pins.

## 4. `api`

No record for `api`'s files: they are git, at the commit a claim links beside it.

## 5. `learner-model`

No record for `learner-model`'s files and releases: they are git, at the commit or tag a claim links beside it, and its Go module is held by `go.sum`.

## 6. `sim`

No record for `sim`'s files, releases and judgements: they are git, at the commit or tag a claim links beside it, and what a run of its gate or of a rough look gave is in [measurements.md](measurements.md).

## 7. `.github`

No record for the shared workflows: they are `.github`'s files at the tag the workflows of `.github/workflows/` pin, git. What stays is a repository's setting, an answer of GitHub's API with no version.

### 7.3. The ruleset of tags of `learner-model`

By GitHub's API, `repos/altnovatrail/learner-model/rulesets`, read 2026-11-08: the ruleset `v-tags` — id 26770474, `refs/tags/v*`, `update`, `deletion`, `non_fast_forward`, without bypass.

## 8. GitHub, Codecov, the base image and the licenses of Go's tools

No record for Go's specification and code, the Go Wiki or spec's tags: a link at a version, a commit or a tag beside the claim, or the measurement.

### 8.3. The arm64 runners

[GitHub Docs, GitHub-hosted runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), read 2026-11-07: "Use of the standard GitHub-hosted runners is free and unlimited on public repositories"; in the table of the standard runners for public repositories — "Linux | 4 | 16 GB | 14 GB | arm64 | `ubuntu-24.04-arm`, `ubuntu-22.04-arm`, `ubuntu-26.04-arm`". `altnovatrail/engine` is public (`gh repo view`, 2026-11-07).

### 8.4. The licenses of the tools of the checks and of gopls

By GitHub's API, `repos/<repository>/license`, read 2026-11-08: govulncheck and gopls — BSD-3-Clause (`golang/vuln` and `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`).

### 8.5. Codecov

Read 2026-11-08: the Codecov app is installed in the organization for all repositories — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `engine` is already in Codecov, but not activated — `https://api.codecov.io/api/v2/github/altnovatrail/repos/engine/`, `"active":false,"activated":false`.

### 8.6. The digest of the base image `devcontainers/base` 3.0.13

Read 2026-11-16: the registry, `https://mcr.microsoft.com/v2/devcontainers/base/manifests/3.0.13-noble` with the media types of an image index — `docker-content-digest: sha256:48be4adec1c1ecbe030064700aef85382f5205786a8f82937fc44a7473717eae`.
