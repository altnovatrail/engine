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

### 6.1. The choice at `sim` and the translation of a package

Read 2026-11-10: RUN.md [`sim`](https://github.com/altnovatrail/sim/blob/0364b40e961dbef87b07572baee03d408e8e090f/RUN.md), `main` 0364b40, the row В10 — «один вызов `engine.Choose` идёт около 30 мс на машине разработки», «В профиле `sim`, 2026-11-09: `engine.Choose` — около 60 % процессора, `prepared.candidates` и `probabilities` — его большая часть, сборка мусора — около 23 %», — and the outcome recorded there — «ученик олимпиадного графа — около 0,65 с вместо около 0,007 с». The translation of a package at `sim` is [`internal/pack/translate.go`](https://github.com/altnovatrail/sim/blob/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2/internal/pack/translate.go) `v0.1.2` e61f90a: `Template: st.Template != nil`, and `KeyProvenance` and `SkillTagsProvenance` — `Readout.Task` of `learner-model` v0.1.0 does not read them. `learner-model` v0.1.0, `readout.go` and `fold.go`, from the Go module cache: `Readout.Task` on each call checks the description of the task, parses its ID by `strings.SplitN`, and copies and sorts the items.

### 6.2. The gate of `sim`

[sim `v0.1.0`](https://github.com/altnovatrail/sim/tree/09e5a84f8a26fedcf5c8996bf81bf29c616917a2), commit 09e5a84, an immutable release, read 2026-11-09: the README, the sections "Gate" and "Taking a release", — "`just gate` builds the bench on the commits of the libraries that gate.json names, through a `go.work` it writes in a temporary folder, runs the gate and holds its numbers to their bands […] A caller puts its own code in place of a library with `GATE_LEARNER_MODEL=<folder>` or `GATE_ENGINE=<folder>`"; the recipe `gate` in the `justfile`; `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `spec` v0.9.9, `go 1.27.3`, 65 bands; the job `gate` in `.github/workflows/ci.yml` — the model for the job "Gate"; `docs/architecture.md` of `sim` — what the caller takes and why Dependabot does not raise the tag.

### 6.3. The gate of `sim` on the new pin

[sim `v0.1.2`](https://github.com/altnovatrail/sim/tree/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2), commit e61f90a, an immutable release of 2026-11-09T23:44:12Z, read 2026-11-10: `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `modules` v0.1.4 with `olympiad-math` 0.1.1, `spec` v0.9.9, `go 1.27.3`, 69 bands with the reason "a new pin, modules v0.1.4 with olympiad-math 0.1.1, A-linked on the two published maths, 100 learners a cell, beside the decisive run that fixed the paths anew […]"; `internal/gate/gate.go` — `const Learners = 100`; RUN.md of `sim`, on the new pin, — "Gate / amd64" 11 min in the CI of `sim`. The stage `toolchain` of `.devcontainer/Dockerfile` is the same as at `v0.1.0`.

### 6.4. The gate of `sim` after the choice was sped up

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, an immutable release of 2026-11-11T02:54:16Z, read 2026-11-11: `gate.json` — `engine` `v0.1.6` 57e0ab8, `learner-model` `v0.1.4` acfc4f2, `modules` v0.2.0, `spec` v0.12.0, `go 1.27.3`; `internal/gate/gate.go` — `const Learners = 150`; `docs/decisions.md` of `sim`, the entry on the gate, — «база гейта — `engine` v0.1.6, и в гейте 150 учеников на ячейку, а не 100», «на `engine` v0.1.6 гейт на 100 учениках читает все 69 чисел те же до бита, что на v0.1.4», «в пересчёте по самому медленному замеру на 150 учеников: у `sim` — около 8 мин, у `engine` — около 11,5, у `learner-model` — около 11,2»; its entry on the decisive run — on 1,500 learners. The stage `toolchain` of `.devcontainer/Dockerfile` is the same as at `v0.1.2`: the files differ only in the comments of the gopls stage.

### 6.5. The corridor of the base

[altnovatrail/sim](https://github.com/altnovatrail/sim/tree/v0.1.2), the tag `v0.1.2`, commit e61f90a, read 2026-11-10: `internal/measure/engine.go` — `E1` counts the items with a verdict by `PGiven`, P* at issue, with the corridor 0.70–0.85 and the thresholds 0.50 and 0.95; `internal/run/pool.go` — "Ceiling is the largest share of items with a verdict whose P* is in the corridor, of one task of the pool the rule could give"; `docs/decisions.md`, the entry on the paths, the new pin («новое закрепление») — the table of paths: `E1.corridor` — the share outside the corridor, the less the better — for the base and for the ceiling: olympiad — `G0` 0.7347 and 0.0163, `G2-half` 0.7210 and 0.0110, `G2` 0.7313 and 0.0125; school — `G0` 0.8280 and 0.7510, `G2-half` 0.8113 and 0.6866, `G2` 0.8694 and 0.7001. A solution goes to the device «после верного ответа или показа» — by the decision log of `spec`.

### 6.6. The truth of `sim` about an anchor

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, read 2026-11-11: `internal/run/pool.go`, `holds`, — "an anchor is held when a skill of the learner's modules that teaches it is truly mastered with its R* then, or when none teaches it"; `internal/measure/model.go`, `m9`, — "the share of readings "holds" an anchor skill when the truth does not hold it, and of "does not hold" when it does", and the count goes only for `AnchorHas` and `AnchorHasNot`; `internal/learner/assumptions.go` — "TeacherBoost makes the skills that teach an anchor skill stronger by that much: the learner knows the anchor outside the platform (A-unpractised)".

### 6.7. A lesson at `sim`

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, read 2026-11-11: `internal/run/run.go` — the choice of a lesson is played by `PassLesson` (`internal/learner/attempt.go`) in one step, and `engine.Choose` is called once per step; `internal/measure/engine.go`, `e1`, — the items of each `Obs` of a step, the tasks of a lesson too. On `main` d957fc7 there are tags up to `v0.1.8`, and `internal/run`, `internal/learner` and `internal/measure` are the same as at `v0.1.6` (`git diff --stat v0.1.6 d957fc7`).

### 6.8. The simulator for the criterion of consolidation in a lesson

[sim `v0.1.9`](https://github.com/altnovatrail/sim/tree/e1d5e8e33d90f2c3a88634a192988f23b6e56aaf), commit e1d5e8e, `main`, read 2026-11-12 from a clone in a temporary folder:

- `docs/architecture.md`, its part on the populations («Популяции»): `G2` — «растёт: θ* — на 0,01 за ответ, δ* основных навыков ответа — ещё на 0,02, а δ* остальных навыков возвращается на 1 % пути к начальной»; its part on the virtual learner («Виртуальный ученик») — «Базовая популяция не учится, как G0»; its part on the measures («Меры»), `I2`, — «„закрепление“ — прошлая попытка на этом навыке была с неверным ответом, и последнее свидетельство модуля уже не в пробной серии»; its part on the criteria («Критерии») — the criterion of the engine and the group "Corridor" («Коридор»);
- `RUN.md`, E5, «Общее для задач»: «одно закрепление на круг: от тега […] до выбора суда `model-2` `go.mod`, `gate.json`, пакеты и пути не меняются»; table 2, E6 — the judgement `engine-2` «на базе после суда `model-2` … с подтверждением на нетронутом `held-out-engine-1`»;
- `internal/run/run.go`, `Step`: the choice of a lesson is played by `PassLesson` in one step, then each skill of the lesson gets an issue at the time of the choice (`is.Skill, is.Issued, is.LastIssued = s, true, at`), and the module's counter grows by one (`r.choices[m.ID]++`) — the same as after a task;
- `internal/published/published.go`: the release `modules` `v0.2.5`, `olympiad-math` and `grade-3-math` 0.3.0, `english-essays` 0.2.0. `jq` over the `server.json` of their archives (`v0.2.1/sha256-e4de172e….tar.gz`, `v0.2.1/sha256-dd2c6eb2….tar.gz`, `v0.2.0/sha256-85681ea4….tar.gz`): `olympiad-math` 0.3.0 has 41 lessons, each with one skill and one band and three task segments, the difficulties 1, 3, 5 in 37 lessons, and 1, 2, 4, 1, 2, 5, 1, 4, 5 and 2, 3, 5 in the rest; 775 tasks, 123 in lessons, and outside lessons per skill and band of a lesson from 2 to 27. `grade-3-math` 0.3.0 has 4 lessons of three task segments each, of the difficulties 1, 3, 5, and all 12 tasks are in lessons; `english-essays` 0.2.0 has 6 lessons, and all 8 tasks are in lessons.

### 6.9. The step and the criterion of consolidation in a lesson at `sim`

[sim `v0.1.15`](https://github.com/altnovatrail/sim/tree/52a365813ebe807f85b4c3e8d34b1df39f506e4d), commit 52a3658, an immutable release of 2026-11-14T01:52:36Z, read 2026-11-14 from a clone in a temporary folder:

- `docs/decisions.md` of `sim`, the entry on the step within a lesson — "Within a lesson whose skills all have a verdict, the run calls the engine at the end of each task segment, as `api` will call it: this attempt is the last one, the run is open, the time is the end of the segment. The lesson of the same run with `REASON_CODE_OPEN_RUN` is the next segment and no new choice; a task outside the lesson is a consolidation — a choice of its module and an issue of its skill, given at once —, after which the engine is called again and must give the lesson of the same run; any other answer is a run that cannot go on. […] On the base the step gives the same bits." The code is `internal/run/run.go`: `passLesson`, `consolidate` — a task "neither the lesson nor a task outside it" stops the run —, `further` and `Verdicts`;
- the same file, the entry on the judgement — "The judgement of the variant of consolidation in a lesson […] goes by the engine's criterion […] as it is, and the variant's target is the group "Corridor". The note of consolidation in a lesson, the measure `E6`, stands beside it in the report and decides nothing. On school maths and on the essays the candidate gives the base's runs: the constraint `neutral`. The judgement is `judgements/engine-lesson-1/` […] Its held-out learners are its own: `held-out-engine-lesson-1`, seed 20261113. A rough look is run by the recipe `just rough`." The judgement file names `equal_on` — "The judgement file of E7 names `novatrail/grade-3-math` and `novatrail/english-essays`";
- the README, "Judging a candidate" — a candidate is `CANDIDATE_ENGINE=<full SHA>`, or "A folder works in place of the SHA, if it has nothing uncommitted"; "A rough look goes the same way, on fewer learners a cell […] `just rough engine runs/base-300.jsonl` for the base and `CANDIDATE_ENGINE=<full SHA> just rough engine runs/engine-a-300.jsonl` for the candidate — 300 learners a cell unless a third argument names from 1 to 1,999";
- `gate.json` — `engine` `v0.1.6` 57e0ab8, `learner-model` `v0.1.8` ab688b2, `modules` v0.2.5, `spec` v0.16.3, `go 1.27.3`, 69 bands with the reason "[…] the base moves to the release of theta, learner-model v0.1.8 at ab688b2, with the code of a063835 and the same numbers to the bit"; `internal/gate/gate.go` — `const Learners = 250`; the entry on the gate — "the gate has 250 learners a cell, not 150", and at the callers without the Go cache "about 11 min 45 s – 13 min 17 s";
- `.devcontainer/Dockerfile` — the stage `toolchain` differs from that of `v0.1.6` only by Go: `GO_VERSION=1.27.4` and its sha256, where `v0.1.6` has 1.27.3;
- `server.json` of `grade-3-math` 0.3.0 (`internal/published/v0.2.1/sha256-dd2c6eb2….tar.gz`), by `jq`: four lessons, one a skill — `times-tables`, `division-with-remainder`, `fractions-of-a-whole`, `comparing-fractions` —, each with the three tasks of its skill, of the difficulties 1, 3 and 5.

### 6.10. The judgement of consolidation in a lesson at `sim`

[sim `main`](https://github.com/altnovatrail/sim/tree/b171f6f617cfd3b6c8e6342d1b0e98cc78bc104e), commit b171f6f, the merge of pull request #65 of 2026-11-14T15:51:25Z, and the tag `v0.1.17`, read 2026-11-14 from a clone in a temporary folder:

- `docs/decisions.md` of `sim`, the entry on the judgement `engine-lesson-1` — "The judgement of the variant of consolidation in a lesson — […] the commit 0e4680f of the branch `e7-consolidation-in-lesson` — ended in an exit: the candidate fails one constraint of the engine's criterion and is not better than the base. The base stays as it is. There was no confirmation, and `held-out-engine-lesson-1` is still held out. `gate.json`, `go.mod`, the paths and the bands do not change." The judgement ran "on `sim` at the tag `v0.1.17`, one pin for the whole judgement", by decisive runs of "2,000 working learners a cell";
- the same entry, the table of the report — "`E7` | −0.0018 [−0.0047; 0.0010] | −0.0002 [−0.0010; 0.0006] | 1" (score, group `corridor`, failed "no worse") — and "It fails one check, `M3.lag` of `G2-half` on olympiad maths: worse by at most 0.01812, at a tolerance of 0.01785. `I1`–`I4`, `E2` and `neutral` hold: on school maths and on the essays every learner's sums are the base's to the bit";
- the same entry, on the note `E6` — "the items of a consolidation are in the corridor 0.27–0.43 of the time, about twice as often as those of the lesson's segments, but on almost every cell the segments themselves go up to 1.4 points less often into the corridor and more often into "hard", and `E6.next`, P* of the next segment after a failure, is lower by up to 0.023. What a consolidation gains in the corridor the lesson loses"; and "The lag of the model on `G2` is the common failure of A and C in `engine-1` too";
- the same entry, the consequences — "[…] of `engine` records the exit by the link to `judgements/engine-lesson-1/` on `main`; […] of `engine` does not change `spec`'s contract or cards, which change only if a judgement chooses the variant. The outcome goes to `spec` as row В19 «в реестр spec» of table 3 of RUN.md"; "A next judgement of the variant, if `engine` makes one, gets its folder, its own set and its row […] before its first run";
- [`judgements/engine-lesson-1/report.md`](https://github.com/altnovatrail/sim/blob/b171f6f617cfd3b6c8e6342d1b0e98cc78bc104e/judgements/engine-lesson-1/report.md) — "holds every constraint: false; better than the base: false"; "fails no worse G2-half/novatrail/olympiad-math M3.lag: worse by at most 0.01812, the tolerance 0.01785"; `E6.next` on `G2-half/novatrail/olympiad-math` — "0.3033 | 0.281 | -0.02229 [-0.02672, -0.01789]";
- the tag `v0.1.17` — commit 6f5151f5eb241b5b44b11b71296a0e3cd2c201a1, an immutable release of 2026-11-14T05:36:14Z: `gh release verify v0.1.17 -R altnovatrail/sim` — "✓ Release v0.1.17 verified!". Its `gate.json` differs from that of `v0.1.15` only by `"github.com/altnovatrail/spec": "v0.17.3"` in place of `"v0.16.3"`, and is the same on b171f6f; `.devcontainer/Dockerfile` is the same as at `v0.1.15`; `internal/gate/gate.go` — `const Learners = 250`.

### 6.11. The order of the judgements at `sim` after E7

[RUN.md of `sim`](https://github.com/altnovatrail/sim/blob/40092d7f1d09580c92a46f21ea42a1dbcf2d3d61/RUN.md), `main` at commit 40092d7 of 2026-11-14T16:34:35Z, read 2026-11-14 through `gh api`:

- table 1 — «Таблица 1. Текущий эпик — E9. Справедливость модели по классу»;
- table 2, «Эпики на потом», whose order of rows is the order of the queue («Таблица 2 — эпики на потом, в порядке очереди»): first E8 — «Третий круг модели | Суд `model-3` по кандидатам E7 `learner-model` … на базе после суда E7: он кончился выходом, и база та же, `learner-model` v0.1.8 и `engine` v0.1.6», — then E6 — «Второй круг движка | Суд `engine-2` по кандидатам E5 `engine` («Второй круг коридора» в его таблице 2) на базе после судов `model-2`, E7 и E8 — с выбранными кандидатами или без них — и на размерах […], с подтверждением на нетронутом `held-out-engine-1`», «Выбранный вышел релизом, и база стоит на нём, — или выход записан».

### 6.12. The rough look and the measures of `sim` for the corridor on the new base

[sim `v0.1.17`](https://github.com/altnovatrail/sim/tree/6f5151f5eb241b5b44b11b71296a0e3cd2c201a1), commit 6f5151f, read 2026-11-14 from a clone in a temporary folder (`gh release verify v0.1.17 -R altnovatrail/sim` — "✓ Release v0.1.17 verified!"), with the report of the round `engine-1` and the code of `learner-model` `v0.1.8` from the Go module cache:

- `justfile`, the recipe `rough` — "the run of just decisive, the same build and the same candidates, on fewer learners a cell — from 1 to 1,999, 300 unless named, as engine's rough looks ran — and on the working learners alone"; the recipe `workspace` — a library "from <prefix>_<NAME> — GATE_ENGINE, CANDIDATE_LEARNER_MODEL — when the caller sets it: a commit by its full SHA, fetched from GitHub even when it is on no branch, as the head of a pull request is";
- `internal/measure/model.go`, `m3` — "the lag — θ + δ − θ* − δ* of the skills with evidence from the 101st answer on, on average"; `internal/criterion/measures.go` — `"M3.lag": {absolute, 0.01, false}`; `docs/architecture.md` of `sim`, its part on the criteria — "a signed measure by its absolute value";
- `internal/measure/engine.go`, `e4` — "of the skills within reach at some choice, the share given at least once";
- `internal/run/observe.go`, `Obs` — `Repeat bool // the task, not a template, had its solution seen before: a repeat (M8)`; `ItemObs.Help` — "the answer had help: the technique on request, or a task whose solution was seen"; `internal/run/run.go`, `Step` — `Choice engine.Choice`, `Obs []*Obs`, and `Pool` — "what the truth says of the pool at the choice";
- `internal/learner/attempt.go`, `learn` — "the module's level grows by LearnLevel and the offsets of the item's primary skills by LearnSkill more"; `internal/population/population.go` — `a.LearnLevel, a.LearnSkill, a.LearnReturn = 0.01, 0.02, 0.01` for `G2`; `docs/architecture.md` of `sim`, the population `W-repeat` — "on a repeat of a task whose solution has already gone to the device, within 30 days, P* = 0.95";
- `learner-model` `v0.1.8`, `step.go`, `helped` — "the solution of this task sent to the device before, or a sketch. […] An answer with help moves the step, the counts and τ as it is, but masters nothing";
- the base before `theta`: the history of `gate.json` on `v0.1.17` (`gh api repos/altnovatrail/sim/commits?path=gate.json`) — 77828ca `Move the base to learner-model v0.1.6 and read the paths anew (…)`, cc18998 `Judge model-2: theta chosen and confirmed, the base on its commit (…)`, 5c864bb `Move the base to the release of theta, learner-model v0.1.8 (…)`; `docs/decisions.md` of `sim`, the entry on the gate — "the gate file names `learner-model` v0.1.6 at 7d796e2"; `learner-model` `v0.1.6` is commit 7d796e290599eeb15481f0098c9356bf6c31afe9 (`gh api repos/altnovatrail/learner-model/commits/v0.1.6`); the entry on `theta` — "after a pause of a module of at least 14 days its weight of answers n falls to 20 before the step";
- [`judgements/engine-1/report.md`](https://github.com/altnovatrail/sim/blob/9cd59ef7a51c0fa2c383686b8f9a2c4e1bee9ecd/judgements/engine-1/report.md), read through `gh api` — A: "fails no worse G2/novatrail/olympiad-math M3.lag: worse by at most 0.01855, the tolerance 0.01651", "fails no worse G2/novatrail/grade-3-math M3.lag: worse by at most 0.03632, the tolerance 0.01854", "fails no worse F-growth5/novatrail/olympiad-math E4.given: worse by at most 0.01015, the tolerance 0.01"; C: "fails no worse G2/novatrail/grade-3-math E1.corridor: worse by at most 0.01004, the tolerance 0.01", "fails no worse G2/novatrail/grade-3-math M3.lag: worse by at most 0.03553, the tolerance 0.01836", "fails no worse G2-fading/novatrail/olympiad-math M3.lag: worse by at most 0.02152, the tolerance 0.02019". The heads of the pull requests are 7ef8288c677a8879bd50b89e6e6320584807d7d5 (#46, A) and 97abff0d1078025a646ff574c73ba10a3298337a (#48, C) (`gh pr view 46` and `48 --json headRefOid`). The entry of the log of `sim` on the judgement `engine-1` — the score and the group `corridor`: "| A | 0.124 [0.115; 0.133] | 0.133 [0.122; 0.144] | 3 |" and "| C | 0.127 [0.117; 0.138] | 0.137 [0.126; 0.148] | 3 |".

### 6.13. The queue of the judgements and the constraint `fair` at `sim`, before `engine-2`

Read 2026-11-14 through `gh api`:

- [RUN.md of `sim`](https://github.com/altnovatrail/sim/blob/7151c42bd028e30b71f177a1e8a5995cbdb1a126/RUN.md), `main` at 7151c42 of 2026-11-14T21:08:38Z: table 1 — «Таблица 1. Текущий эпик — E9. Справедливость модели по классу», three tasks `[x]`, the run of the fairness of the base and the end of the epic `[ ]`; table 2, in the order of the queue, — E8 «Третий круг модели», whose basis says «разведка — […], на 2026-11-14 не начата», then E6 «Второй круг движка», whose basis says «Закрепление — тег […] или позже: ограничение справедливости держит и критерий движка»; in the common part of E9 — «E8 и E6 ждут тега […] Если кандидаты `learner-model` будут готовы раньше, порядок решаю я»; row В20 — «Ограничение справедливости `fair` держит и критерий движка […] Популяции `D-` не двигают строк 65 ячеек, так что числа разведки «Второго круга коридора» `engine` на прежнем теге `sim` остаются в силе, но `fair` она не видит; на теге […] `just rough` гоняет и ячейки `D-`. Тег […] — `v0.1.19`, вышел 2026-11-14 на мерже 7553035», with the recommendation «(а) `engine` берёт тег […] в «Gate» и грубые взгляды, как он выйдет, до выбора своих кандидатов»;
- `docs/decisions.md` of `sim` at the same commit: the entry on the order of the first judgements — "Of the first two judgements — of the model and of the engine — the first is the one whose candidates are ready earlier, and the second goes on the base with what the first chose", and among its rejected options "judge as a pair — the engine's candidates would wait for the model's candidates"; the entry on fairness — "The constraint `fair` — no worse than the base — holds the criteria of the model and of the engine and a pair"; the entry on `fair` in `sim judge` — "`sim judge` gives [it] in the model's criterion, the engine's and a pair — once in a pair — at the choice and in the confirmation", "A grade whose cell neither run has fails it";
- the release [sim `v0.1.19`](https://github.com/altnovatrail/sim/releases/tag/v0.1.19), published 2026-11-14T20:59:33Z on 75530353c563e4d21440a5564f02eb7d5eb2b3af; `gh release verify v0.1.19 -R altnovatrail/sim` — "✓ Release v0.1.19 verified!";
- [RUN.md of `learner-model`](https://github.com/altnovatrail/learner-model/blob/d2d7ccd475c22c7dbdf41efb8565f3faf525a045/RUN.md), `main` at d2d7ccd of 2026-11-14T00:07:50Z: «Таблица 1. Текущий эпик — E7. Третий круг кандидатов модели», the reconnaissance of the third round, «Разведка третьего круга», and its candidates, «Кандидаты — черновые pull request», both `[ ]`.

Read again 2026-11-15, from a clone of `sim` in a temporary folder, `main` at d9e7bcd2680bc886ed283d25ec97fe87aafdf680 of 2026-11-15T00:25:16Z:

- RUN.md: the run of the fairness of the base is `[x]`, the end of E9 `[ ]`; table 2 in the same order, E8 then E6, whose basis now says «Закрепление — тег […], `v0.1.20`, или позже»; row В20 is extended — «[…] `engine-2` закрепляет тег […] или более поздний — на нём `paths.json` называет файл базы с 68 ячейками; биты и `fair` те же […]. Тег […] — `v0.1.20`, вышел 2026-11-15 на мерже 63288d7»;
- `docs/decisions.md` of `sim`: the entry on `fair` in `sim judge` — "the judgements `model-3` and `engine-2` pin the tag of […], not `v0.1.19`. The merge of […] changes `paths.json`, which names the base's file of 68 cells […]. The code, the bits and `fair` of the two tags are the same"; the entry on the fairness of the base — "The base — `learner-model` v0.1.8 and `engine` v0.1.6 — does not hold the criterion of fairness by the grade of […]", "The base's file of 68 cells is the base of `model-3` and `engine-2`, and `paths.json` names it"; what follows from it is row В22 of RUN.md of `sim`;
- the release [sim `v0.1.20`](https://github.com/altnovatrail/sim/releases/tag/v0.1.20), published 2026-11-15T00:13:14Z on 63288d7631a6341333d1aa63975e67765de1b637; `gh release verify v0.1.20 -R altnovatrail/sim` — "✓ Release v0.1.20 verified!";
- `git diff --stat v0.1.17 v0.1.20 -- gate.json .devcontainer/ go.mod go.sum justfile` is empty; from `v0.1.19` to `v0.1.20` `internal/criterion/paths.json` changes only `"base"`, `2e53e18c…` → `9e89cb7c…`;
- RUN.md of `learner-model`: `main` is still d2d7ccd, the reconnaissance of the third round `[ ]`.

Read 2026-11-15, from a clone of `sim` `v0.1.20` in a temporary folder, at 63288d7 (`gh release verify v0.1.20 -R altnovatrail/sim` — "✓ Release v0.1.20 verified!"; `git ls-remote` — `63288d7631a6341333d1aa63975e67765de1b637 refs/tags/v0.1.20`):

- `docs/decisions.md` of `sim`: the entry on `fair` in `sim judge` — "on it `just decisive engine` gives the file that `paths.json` names, as step 1 of “Judging a candidate” in README says, while on `v0.1.19` it gives a file that `paths.json` there does not name"; the entry on fairness — "The populations `D-grade1`, `D-grade3` and `D-grade5` on olympiad maths differ by the grade alone, with one truth and one unit for the learner n of all three"; the entry on `fair` in `sim judge` — "`fair` gives, for `D-grade1` and `D-grade5` on olympiad maths, three checks each: `M11.same` by its shortfall, `M11.answers` and `M11.band` as they are", "within the larger of the author's tolerance — 0.01, one answer, 0.01 — and 4.3 standard errors of that difference";
- `gate.json`: `.bands.numbers` holds 69 numbers, and `.libraries` names `engine` `v0.1.6` at 57e0ab8 and `learner-model` `v0.1.8` at ab688b2;
- `sim evaluate` on a rough look of 300 learners refuses it: "cell D-grade1/novatrail/olympiad-math has 300 learners, and an evaluation reads a decisive run of 2000"; the report of `sim judge` prints only the checks that fail (`internal/criterion/report.go`).

### 6.14. The rounds by plan stop at `sim`, after E9

Read 2026-11-15 through `gh api`, `main` of `sim` at 59885beb618ac89ffd22a4d659b04045bfa669eb of 2026-11-15T01:00:57Z, the merge of the end of E9:

- `docs/decisions.md` of `sim`, the entry on the rounds by plan — "`sim` holds no rounds by plan. It stays the gate in the CI of `learner-model` and `engine`, and a judgement is laid out as an epic only when a library brings a candidate with a reason of its own — a fault that the gate or a measurement found, a need of the product, a row of `spec` — and names it. The epics E8, `model-3`, and E6, `engine-2`, leave table 2 of RUN.md"; among its findings — "the candidates of a round seek small gains on virtual learners whose speeds and laws of forgetting are assumptions, not checked against live learners (Д3)"; among its reasons — "A round by plan tunes the libraries to the simulator's assumptions; a candidate with a reason of its own answers something that matters, and its judgement is worth its runs. D, E and F may come this way, with the reason named by `engine`"; among the rejected — "E6 next, the proposal of […] — a round by plan, on candidates whose rough look already fails a constraint";
- RUN.md of `sim`, table 1 — «Текущего эпика нет: E9 сделан 2026-11-15», and «[…] Кругов по плану нет: `sim` остаётся гейтом в CI `learner-model` и `engine`, а суд разворачивается, только когда библиотека приносит кандидата со своей причиной. E8 и E6 ушли из таблицы 2»; table 2 holds one row, E10 «Суд круга справедливости модели»;
- RUN.md of `sim`, row В24, «в реестр spec» — «Так могут прийти и кандидаты D, E и F E5 `engine` […]: […] ждут теперь не очереди `sim`, а такой причины», with the recommendation «(а) записать в […] `spec`: суды по плану сняты, `engine` и `learner-model` решают у себя, нести ли кандидатов с причиной или закрыть круги»;
- the latest release of `sim` is still `v0.1.20` (`gh api repos/altnovatrail/sim/tags`), and `gh api repos/altnovatrail/sim/compare/v0.1.20...main` changes only `DEBT.md`, `README.md`, `RUN.md`, `docs/architecture.md`, `docs/decisions.md`, `docs/measurements.md`, `docs/plan/E9.md` and `docs/sources.md`: `gate.json`, the code and the container are those of the tag.

### 6.15. The gate as jobs at `sim` `v0.1.24`

Read 2026-11-16, from a clone of `sim` in a temporary folder at the tag `v0.1.24`, 33d70f288682d95e12f6072c4e9ff09253417497 of 2026-11-16T00:24:27-06:00:

- the release [sim `v0.1.24`](https://github.com/altnovatrail/sim/releases/tag/v0.1.24), published 2026-11-16T06:33:22Z, immutable, the latest (`gh release list -R altnovatrail/sim`), its one change "[…] the gate of CI as jobs, TestEngineOnRuns by processes […]"; `gh release verify v0.1.24 -R altnovatrail/sim` — "✓ Release v0.1.24 verified!", the file `sim-v0.1.24.tar.gz` sha256 `31e41aa7ae999b2c32d96a089c1f15d2bfd75f3dd2b6f6f8ad935d4a850cc9bf`;
- `docs/decisions.md` of `sim`, the entry on the gate as jobs — "CI runs the gate as jobs: the workflow `.github/workflows/gate.yml`, which `ci.yml` calls once an architecture and a caller can call by a tag of `sim`, runs four jobs of `just gate-part <k> 4` and after them one of `just gate-check`"; "Job k of n runs, as processes, the parts k + n·j of n·p, j below p, with p its processors, so the parts of all the jobs are each part of n·p once; `gate-check` builds the workspace again, holds the build to `gate.json` and merges the parts by `bench.Merge`, which refuses a part missing, repeated or of another run"; "on 2026-11-16 `just gate-check exact` over four `just gate-part` in the container gave the table of `just gate exact` to the byte"; "a fixed number of processes a job — the runners of a gate are of one kind, and a job with other processors fails the check loudly, by the number of parts" (among the rejected);
- `justfile` of `sim`: `gate-part k n dir` — "Run part k of n of the gate into a folder, as one job of n", "It holds the numbers to nothing; gate-check does"; `gate-check mode dir` — "Check the parts of the gate in a folder: mode bands or exact"; both resolve a relative folder against the caller's folder, and `gate` is the same run in one job, `gate-bench … "$mode" 0 1 -`;
- README of `sim`, "Gate" — "A caller calls it by the SHA of a tag of sim", with `name: Gate`, `uses: altnovatrail/sim/.github/workflows/gate.yml@<SHA of the tag>`, `name: amd64`, and "the check is "Gate / amd64""; `gate.yml` names its check `${{ inputs.name }}` and its parts `${{ inputs.name }} / part ${{ matrix.part }}`, each `timeout-minutes: 15`, on `inputs.runner`;
- `git diff v0.1.20 v0.1.24 -- gate.json` changes only `"github.com/altnovatrail/spec": "v0.17.3"` → `"v0.18.3"`, so the 69 bands and the commits of the libraries are those of `v0.1.20`; `.devcontainer/Dockerfile` changes the comments of the stages and the base image of `toolchain`, `mcr.microsoft.com/devcontainers/base:3.0.12-noble` → `3.0.13-noble@sha256:48be4ade…`, so its hash, and the key of the cache's entry, are new.

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
