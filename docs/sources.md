# engine — sources

> The external facts that the documents of the repository's `docs/` rest on: the rule and the code of standalone, research and the documentation of other programs, the documents, code and release files of other NovaTrail repositories, the documentation of Go, GitHub and Codecov. Here lie their records — the address, the version or the date of the page, the date read and the quote in the language of the source, — and [architecture.md](architecture.md) holds the claim and a reference to the record by its number, as to a section: `sources 1.1` (RULES.md, The bar; R143 spec).
>
> A record is a subsection "N.M" of a topic: one source, or one reading of several sources for one decision. The number of a record does not change: a new one goes at the end of its topic, a new topic at the end of the file. A source read again is an item with its date in the same record, and the earlier text is not edited. The model is sources spec. The results of measurements are not here but in [measurements.md](measurements.md).
>
> Records 1.1–8.5 were moved here word for word on 2026-11-12 from [architecture.md](architecture.md), section 4, "Sources", by task T55 (R24). The first line of a record names the subsections of architecture and the decisions it was read for, if the item named them. Section 4 stayed as a pointer: the decision log and the archives refer to it, and they are not rewritten.

## 1. MathTrail: standalone

### 1.1. The rule by which standalone chooses a task

For architecture 1, 2.1–2.3.

[mathtrail/mathtrail-standalone, `SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), branch `main`, commit c6cbeb2, read 2026-10-19:

- SPEC 3.1: "It is deterministic by design (the prototype's D40): the same profile always produces the same brief"; "Nor does the rule read the links between topics (1.2)". The fields of the brief: `pedagogical_goal` — "`reinforce` or `new_topic`", `rationale` — "One sentence naming what the choice was based on, and the model's reason when the model overrode the choice";
- SPEC 3.2, step 5: "If every topic within reach is mastered, the same choice runs over all of them";
- SPEC 2.5, read 2026-10-20: "The topic counts as mastered only while its recommended point (2.3) stays at that level or below: once the child's level in the topic moves the recommended point up a level, the topic is back in the rotation"; "Mastery is lost after two consecutive wrong answers in that topic, its `wrong_streak`: both fields are cleared, at whatever level mastery was held, and the topic returns to the rotation"; "A mastered topic is skipped when the rule looks for something new, until every topic within reach is mastered — after which all of them are back in the rotation";
- SPEC 2.3: "A topic is **within reach** when its easiest point — difficulty 1 of the lowest level it is taught at — has P ≥ 0.70"; a topic on pause "comes back once in six tasks as a probe at that point";
- SPEC 3.2, "The idea": a task's idea "is drawn among the rest by a hash of `student_id`, the topic, the level, the round and how many ideas the record of the list holds";
- SPEC 3.4: "The model may ask for a different topic, level or difficulty, and must give a reason"; "`open_request.tutor_mode` becomes `llm` instead of `rule`". SPEC 3.6: when the lessons keep to a topic that a person chose, "`open_request.tutor_mode` is **`person`**".

### 1.2. The code of standalone's rule

The same commit c6cbeb2 as in sources 1.1, read 2026-10-19:

- the choice is the package [`internal/domain/tutor`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/tutor): "Package tutor decides what the child is asked next"; "Nothing here is random and nothing here is a model". On the choice of an idea in `idea.go`: "looks drawn at random — another order for every child and every round — and is not: it is read off a hash of the child's id". This package has nothing to do with NovaTrail's `tutor`: only its name is the same. Beside it lie its tests, among them `golden_test.go`, `property_test.go` and `simulation_test.go`;
- the corridor, reachability and "mastered" are the package [`internal/domain/rating`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), the files `corridor.go`, `ladder.go` and `mastery.go`. In `corridor.go`: `CorridorLow  = 0.70`, `CorridorHigh = 0.85`, `CorridorMiddle = (CorridorLow + CorridorHigh) / 2`, and above the last one — "CorridorMiddle is what the recommendation aims at, and the line an answer has to be at or below to count toward mastering a topic" (read 2026-10-20).

### 1.3. The choice rule of standalone, in addition for T01

The same commit c6cbeb2, read 2026-11-07:

- SPEC 2.3: a topic "on pause" — "comes back once in six tasks as a probe at that point, however small its chance there, whose right answer lifts it into reach"; SPEC 3.2, step 4: "those never handed out come first, in catalog order; then the one whose `topics[t].last_issued` is the oldest", and a topic in which the child is "behind" counts twice (R300 standalone);
- SPEC 3.2: "A child below every topic of the catalog […] is given instead the topics whose easiest point is the easiest the catalog has"; step 1: "In the trial series […] a failure is not worked over";
- [`internal/domain/tutor/pause.go`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/internal/domain/tutor/pause.go): `const ProbeEvery = 6`, the turn of a probe is `tasks > rating.TrialAnswers && tasks%ProbeEvery == ProbeEvery-1`.

### 1.4. A topic chosen by a person, in standalone

For architecture 3.2.

SPEC 3.6 standalone, the same commit c6cbeb2, read 2026-11-07: "The lessons may be kept to one topic of the catalog"; the choice writes `student.lesson_topic`, "and an empty one gives the choice back to the rule"; "Once the trial series is over, and while the catalog has the topic, the rule sets it in place of its own"; "During the trial series the choice is kept and waits".

## 2. Research and other programs

### 2.1. A corridor near 85% correct

Wilson R. C., Shenhav A., Straccia M., Cohen J. D. The Eighty Five Percent Rule for optimal learning. *Nature Communications* 10, 4646 (2019), [doi:10.1038/s41467-019-12552-4](https://doi.org/10.1038/s41467-019-12552-4), via [PMC6831579](https://pmc.ncbi.nlm.nih.gov/articles/PMC6831579/), read 2026-11-07: "the optimal error rate for training is around 15.87% or, conversely, that the optimal training accuracy is about 85%". This is the upper bound of standalone's corridor. The lower one, about 0.75, SPEC 2.3 standalone takes from Math Garden ("Math Garden keeps success near 0.75"); the primary source — the article of Math Garden — could not be read: the publisher answers 403, and it is not stated here as fact.

### 2.2. Repetitions and the new in Anki

For architecture 3.4.

[The Anki manual, Deck Options](https://docs.ankiweb.net/deck-options.html), read 2026-11-07. The review limit limits the new too, with a backlog the advice is not to introduce the new, and the order of the new and the repetitions is a setting: "By default, the review limit also applies to new cards, and no new cards will be shown when the review limit has been reached"; "When this limit is reached, Anki will not show any more review cards for the day, even if there are more waiting"; the order of new cards and repetitions — "Whether new cards should be mixed in with review cards, shown before them, or shown after them"; "If you have a backlog of overdue review cards, it is recommended that you stop introducing new cards until you catch up with that backlog".

## 3. `modules`

### 3.1. The graphs of the first slice

[altnovatrail/modules](https://github.com/altnovatrail/modules/tree/6d061e9/modules), `main` 6d061e9, read 2026-11-07: `olympiad-math/skill-graph.yaml`, `grade-3-math/skill-graph.yaml`, the module manifests with bands and the files of `tasks/` — the example of architecture 3.1.

### 3.2. The anchor skills of the first slice's modules

For architecture 3.2.

`altnovatrail/modules`, `main` 6d061e9, read 2026-11-07: the module manifest `olympiad-math/manifest.yaml` — "Its skills need no anchor skill yet"; in [docs/architecture.md](https://github.com/altnovatrail/modules/blob/6d061e9/docs/architecture.md) of `modules`, the table of the port, — «у `number-divisibility` есть ещё навык-опора `core.division-with-remainder`, а у `fractions-parts` — `core.fractions-of-a-whole`» and «Остальные 14 навыков вносит порт»; `grade-3-math/skill-graph.yaml` — `fractions-of-a-whole` with `teaches: [core.fractions-of-a-whole]`. The port's graph is the [format example](../spec/format/examples/skill-graph/valid/olympiad-math.yaml) `olympiad-math.yaml`, 17 skills.

### 3.3. The lessons of the first slice

For architecture 3.3.

`altnovatrail/modules`, `main` 6d061e9, read 2026-11-07, the folders `lessons/` and `tasks/` of the three modules:

- olympiad maths — nine lessons, a lesson per skill and band, `counting-gaps-grades-3-4` and so on; the lesson `counting-gaps-grades-3-4` is explanations, a worked example and the tasks `gaps-34-d1-1`, `gaps-34-d3-1`, `gaps-34-d5-1`: all three tasks of the skill on the band;
- school maths of grade 3 — four lessons, a lesson per skill; the lesson `times-tables` — "an explanation and the three pages of the skill, in the order of their difficulty";
- essays — six lessons and eight tasks, all the tasks in lessons. `reading-model-essays` — "explanations alone, with no task"; `first-sentences` — "on two skills […]: choosing a first sentence, which gets a verdict, and writing one, which gets feedback by a rubric"; `writing-about-a-day` — "for adults alone", `age_categories: [adult]`, the skill `endings`; `endings` — "the engine's rule chooses between it and the lesson of improving an essay".

### 3.4. The packages of the release `modules` v0.1.0 in `testdata/modules`

For architecture 3.7; the third 2026-11-08 amendment to R10.

2026-11-08: the release [v0.1.0](https://github.com/altnovatrail/modules/releases/tag/v0.1.0) is the only one of `altnovatrail/modules` (`gh release list -R altnovatrail/modules`), its tag on the commit `6d061e940808b1ee9f64fe9a30bd209d408a756e` (`gh api repos/altnovatrail/modules/git/ref/tags/v0.1.0`), the same one by which the examples of 3.1–3.3 were written. `just -f spec/justfile release-verify modules v0.1.0 <file> <file>.sigstore.json` accepted each of the four files — "modules v0.1.0, signed by the shared release.yml of .github v0.1.15":

- `module-index.json`, sha256 `ed4d8bb325fdcb3ea54da439e08c58d927f773ba0f097338f7767ad76be8c90e`, — `"module_format_version":"0.3"` and, for the three modules of version `0.1.0`, the package digest;
- `sha256-1d8a962d….tar.gz` — `novatrail/olympiad-math`, `sha256-2e2f7a2f….tar.gz` — `novatrail/grade-3-math`, `sha256-8c777871….tar.gz` — `novatrail/english-essays`; in each, `inventory.json`, `server.json`, `device.json` and the files of the device;
- the package digest is the sha256 of its `inventory.json`, not of the archive: "`package_digest`: the digest of the package, which is the digest of its inventory" ([format text](../spec/format/README.md), section 26); `inventory.json` names the sha256 of `server.json` in `server_part.digest` — for the three packages they matched the sha256 of the files.

### 3.5. The index of `modules` v0.1.4 and the package `olympiad-math` 0.1.1 in `testdata/modules`

For architecture 2.10 and 3.7; the 2026-11-10 amendment to R10.

2026-11-10: the release [v0.1.4](https://github.com/altnovatrail/modules/releases/tag/v0.1.4), its tag on the commit `c23ed6267bc5854e7e728b23cb6541db4d15f38a`, and the release [v0.1.1](https://github.com/altnovatrail/modules/releases/tag/v0.1.1), its tag on the commit `5f6d7a9668ae172e36ca349e3f89e96faa6597ab` (`gh api repos/altnovatrail/modules/git/ref/tags/<tag>`). `just -f spec/justfile release-verify modules v0.1.4 module-index.json` — "module-index.json: modules v0.1.4, signed by the shared release.yml of .github v0.1.15"; `just -f spec/justfile release-verify modules v0.1.1 sha256-123cb0d3….tar.gz` — "sha256-123cb0d3….tar.gz: modules v0.1.1, signed by the shared release.yml of .github v0.1.15"; the bundles are `<file>.sigstore.json` of their releases:

- `module-index.json`, sha256 `a244c0ba2832ba6360f45bc74d293572b84b491fc823c7b01bdf6004ff9edb2a`, — `"module_format_version":"0.3"`; the entries `0.1.0` of the three modules are the same as in the index v0.1.0, byte for byte, and also `novatrail/olympiad-math` `0.1.1` — `"package_digest":"sha256:123cb0d3dec7a094dce8e6c62f3b5f0f8e4563b23d148793fe3b9ce6c4839cb1"`, `"release":"v0.1.1"`;
- `sha256-123cb0d3….tar.gz`, 2,477,140 bytes: in `server.json`, 19 skills, all with `verdict`, three bands, 772 tasks, of them 19 templates — tasks with the field `template` — and 41 lessons; the check kinds are `choice`, `number` and `fraction`; the tasks and the lessons have no age categories of their own, and the module has `under_13` and `adult`; the chain held, as for the packages 0.1.0: the sha256 of `inventory.json` is the package digest in the index, and the sha256 of `server.json` is its `server_part.digest`.

### 3.6. The lessons and tasks of the packages of `testdata/modules`

For architecture 2.12.

2026-11-11, `jq` over the `server.json` of each package in a temporary folder:

- `olympiad-math` 0.1.1, `sha256-123cb0d3….tar.gz`: 41 lessons; each has one skill, the tasks of one band and three task segments; the difficulty of the segments in order is 1, 3, 5 in 37 lessons, and 1, 2, 4, 1, 2, 5, 1, 4, 5 and 2, 3, 5 in one each; the lesson `arithmetic-tricks-grades-1-2` is the tasks `tri-12-d1-1`, `tri-12-d3-1`, `tri-12-d5-1`. There are 123 distinct tasks in lessons out of 772; outside lessons 649, of them 19 templates, and per skill and band of a lesson from 2 to 26;
- `olympiad-math` 0.1.0 — 27 tasks and 9 lessons of three task segments each, `grade-3-math` 0.1.0 — 12 tasks and 4 lessons of three each, `english-essays` 0.1.0 — 8 tasks and 6 lessons: all the tasks are in lessons.

## 4. `api`

### 4.1. What the engine gets from `api`

For architecture 3.2.

[api, `docs/architecture.md`](https://github.com/altnovatrail/api/blob/06dad44/docs/architecture.md) and [`docs/decisions.md`](https://github.com/altnovatrail/api/blob/06dad44/docs/decisions.md), `main` 06dad44, read 2026-11-07: «задачу даёт только ученику её возрастных категорий, `age_categories` (R51 spec), а урок — только если ученику подходят и урок, и каждая его задача»; on a lesson that does not suit — «движок таких уроков и не получает»; the 2026-11-02 amendment to R07 api — «Движок, как и прежде, получает графы текущих версий подключённых модулей, которые ученику подходят».

### 4.2. A lesson at `api`

For architecture 2.12.

[api, `docs/architecture.md`](https://github.com/altnovatrail/api/blob/2024d52015965ad018c54f07ca2d891da5e7b42c/docs/architecture.md), section 3.7, and [`docs/decisions.md`](https://github.com/altnovatrail/api/blob/2024d52015965ad018c54f07ca2d891da5e7b42c/docs/decisions.md), R13, R48 and R50 api, `main` 2024d52, read 2026-11-11: «Попытки закрепления в уроке первый срез не знает (R107 spec); если её введёт новая версия движка, её место — «после звена k»»; «Введёт их новая версия движка — `api` будет звать его на конце звена урока, чьи навыки все с вердиктом, а ответ — следующее звено, закрепление или «досрочно»»; on "early" («досрочно») — «До тех пор `api` урок досрочно не кончает, а среди причин конца в базе и в коде её нет». R48 api — `OpenLesson`, `GetLessonRun` and `CompleteLesson` in `internal/learning`; R50 api — the engine is called on "something else", on `next` outside a lesson and on `RunEnd`.

## 5. `learner-model`

### 5.1. Knowing an anchor skill

For architecture 3.2.

[learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/98ba29f/docs/decisions.md), `main` 98ba29f, — R07 learner-model, read 2026-11-07: «Держать ли урок при «нет свидетельств», решает `engine`»; rejected is "unknown" while the skills that teach it have no answers — «урок сначала открыт, а с первыми ответами закрывается».

### 5.2. What `Read` gives the engine

[learner-model, `docs/architecture.md`](https://github.com/altnovatrail/learner-model/blob/98ba29f/docs/architecture.md), `main` 98ba29f, read 2026-11-07, the table of operations: for a skill with practice — «число свидетельств, время последней практики, θ + δ, v и удержание на t₀, ловушки с числом и временем последней»; for a mastered one — «полоса освоения, осторожная вероятность с удержанием и срок повторения»; `Readout.Task` — «P каждого пункта на t₀». There was no history of issues and no mark of the end of the trial series there at 98ba29f.

### 5.3. The types of `learner-model` in code

For architecture 2.4.

[learner-model, `types.go`](https://github.com/altnovatrail/learner-model/blob/6a6c843/types.go) and [`readout.go`](https://github.com/altnovatrail/learner-model/blob/6a6c843/readout.go), `main` 6a6c843, read 2026-11-07: `Graph`, `Skill`, `Band`, `RetiredSkill`, `Task`, `TaskItem`, `Slot`, `Readout`, `ModuleReadout`, `SkillReadout`, `AnchorReadout` and `TaskReadout`; `Task` — "describes a task, the same for evidence and for a candidate task of the engine", and `readout.go` — "These are the types alone: Read and Readout.Task compute them from T20 on". `ModuleReadout` has `HasEvidence`, `StartBand` and `Skills`, and the code has no marks of the trial series yet. In [docs/architecture.md](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/architecture.md) of `learner-model` — in the section «Типы и кодировка» — they are already in the table of operations: `Read` gives «у модуля — кончилась ли пробная серия и было ли последнее свидетельство ещё в ней», and «Уроки, возрастные категории и другое, что нужно только выбору, движок держит своими типами вокруг описания задачи».

### 5.4. The build and the same bits at `learner-model`

For architecture 3.7.

[learner-model, `docs/architecture.md`](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/architecture.md), the section «Сборка и одни биты», and [`docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/decisions.md) — R09 learner-model with its amendments, `main` 6a6c843, read 2026-11-07: eight rules of code, the lists of packages, the recipes `test-targets` and `fma-diff`, the test of the rules `rules_test.go` with the control `testdata/coderules`, the job "Go arm64" and, on Т-10, «`engine` для повтора выбора на любой сборке нужна та же дисциплина — его T05».

### 5.5. The reading of `learner-model` in code

For architecture 3.6.

[learner-model, `readout.go`](https://github.com/altnovatrail/learner-model/blob/fc444ba/readout.go), `main` fc444ba, T20, read 2026-11-07: `func Read(s State, g Graphs, at int64, grade Grade, engineModules []string) (Readout, error)` and `func (r Readout) Task(t Task) (TaskReadout, error)`; `ModuleReadout` has "TrialOver tells whether the module's trial series is over, and LastInTrial whether the last" evidence was still in it, the fields `TrialOver` and `LastInTrial`. The first signed release is [`v0.1.0`](https://github.com/altnovatrail/learner-model/releases/tag/v0.1.0), its tag on the commit d2620a1, published 2026-11-07T20:12:08Z: `gh release verify v0.1.0 -R altnovatrail/learner-model` — "Release v0.1.0 verified!" (2026-11-08). It enters the code of `engine` only once `just -f spec/justfile release-verify` accepts it (R77 spec).

### 5.6. The frame and gopls at `learner-model`

For architecture 3.7; R12, R13.

[learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/50f1b51/docs/decisions.md), R13 and R14 learner-model with their amendments, and the `justfile`, `.github/`, `.devcontainer/` of the same commit, `main` 50f1b51, read 2026-11-08: the model for the recipes `ci-*`, the job "Go", `codecov.yml` and the step of gopls. From the amendments to R13 learner-model: on a pull request without lines of code Codecov does not set `codecov/patch` — «Codecov видит у него 0 файлов и 0 строк […] и статуса не ставит, так что мерж ждал его без конца», — and the first status came with the first code under tests.

### 5.7. `learner-model` `v0.1.0` in code

For architecture 3.6; the second 2026-11-08 amendment to R09.

2026-11-08: `gh release download v0.1.0 -R altnovatrail/learner-model` and `just -f spec/justfile release-verify learner-model v0.1.0 learner-model-proto-v0.1.0.tar.gz` — "learner-model-proto-v0.1.0.tar.gz: learner-model v0.1.0, signed by the shared release.yml of .github v0.1.15"; the tag `v0.1.0` is on the commit `d2620a18d865c813ec2556c93fca92a1445e3a74` (`gh api repos/altnovatrail/learner-model/git/ref/tags/v0.1.0`); `go mod download -json github.com/altnovatrail/learner-model@v0.1.0` through `proxy.golang.org` — `"Origin":{"Hash":"d2620a18d865c813ec2556c93fca92a1445e3a74","Ref":"refs/tags/v0.1.0"}`, `"Sum":"h1:KK3u66KFWfWV8NmI+ueCXZQVXyUrg8zfDtiyo8+Me9Q="` — the same sums that `go get` wrote into `go.sum`. That a Go module is not a release file, and that `go.sum` and sum.golang.org hold it, is R16 learner-model ([learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/50f1b51/docs/decisions.md), `main` 50f1b51, read 2026-11-08). The API is [`readout.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.0/readout.go) of the same tag: `Readout` has fields not visible from outside, and only `Read` builds it.

### 5.8. `learner-model` `v0.1.6` in code

For architecture 2.10, 3.6 and 3.7; R23.

2026-11-11: `gh release list -R altnovatrail/learner-model` — the releases `v0.1.0`…`v0.1.6`, the last `v0.1.6`, published 2026-11-11T17:32:38Z, `isImmutable` — `true`; `gh release download v0.1.6` and `just -f spec/justfile release-verify learner-model v0.1.6 learner-model-proto-v0.1.6.tar.gz` — "learner-model-proto-v0.1.6.tar.gz: learner-model v0.1.6, signed by the shared release.yml of .github v0.1.21", the file's sha256 `2ec63cc68d2396cc96cee34c09281290bfc3e209cb5037f100c85fc55b789d5f`; `gh release verify v0.1.6 -R altnovatrail/learner-model` — "Release v0.1.6 verified!"; the tag `v0.1.6` is on the commit `7d796e290599eeb15481f0098c9356bf6c31afe9`; `go mod download -json github.com/altnovatrail/learner-model@v0.1.6` through `proxy.golang.org` — Origin `refs/tags/v0.1.6` on the same commit, `"Sum":"h1:xoifdkdFlXvTSgjFu939p2Ol2oB7s6G6pCO3h2GxfxU="` — the same sum that `go get` wrote into `go.sum`; the `go` line of its `go.mod` is `1.27.3`, and there is no `require`. What changed since `v0.1.0` — `git diff v0.1.0 v0.1.6` of the same repository: [`types.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/types.go) — `CheckKindTime`; [`evidence.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/evidence.go) — `EvidenceKind = 2` and `LeadingQuestionOnRequest`; [`readout.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/readout.go) — at `Read`: "a module of the engine with no state and a current version gives the starting point by the grade, BandByGrade, or none when that version has no bands", at `Readout.Task`: "It copies nothing it need not"; R28 and R29 learner-model and the 2026-11-09 amendment to R07 learner-model — [learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/docs/decisions.md) of the same tag.

### 5.9. An anchor skill without a teacher in the model and in the package

For architecture 2.11.

Read 2026-11-11:

- the package `olympiad-math` 0.1.1 in `testdata/modules`, `sha256-123cb0d3….tar.gz`, the file `server.json`: `fractions-parts` has `"prerequisites":["arithmetic-tricks","core.fractions-of-a-whole"]`, `number-divisibility` — `["time-calendar","parity-alternation","core.division-with-remainder"]`, `percent-basic` and `ratio-sharing` — `["fractions-parts"]`; 16 tasks have `"skills":{"primary":["fractions-parts"],"supporting":["core.fractions-of-a-whole"]}`, 16 — `"skills":{"primary":["number-divisibility"],"supporting":["core.division-with-remainder"]}`, and these two skills have no other tasks;
- `learner-model` v0.1.0 from the Go module cache: `step.go`, `itemSkills` — evidence moves the primary skills of the item or of the task and does not read the supporting ones; `readout.go` — `AnchorUnknown` with the reasons `AnchorReasonNoTeacher` and `AnchorReasonNoVerdictTeacher`; `docs/decisions.md`, R07 learner-model, — «навык-опора как вспомогательный навык задачи ничего не получает (R06), и модель его в свидетельстве не читает»;
- `spec` v0.15.0: R63 spec — «Ученику модуль предлагает `api` — только из тех, что держит сам и что ученику подходят»; [format text](../spec/format/README.md), section 26, — "The platform, for a module that a learner lacks […]: the modules that teach an anchor skill are those whose entry names it in `teaches`"; `progress.proto` — `AnchorSkillState` has the fields `anchor_skill`, `holding`, `reason` and `teaching_skills`, "The full IDs of the skills that teach it, in the modules whose graphs the learner model has"; `module.proto` — `CatalogModule` has only `module_id` and `package_digest`; R100 spec, «Отвергнуто», — among the branches of standalone «короткая проверка навыка-опоры», and about them — «движок может их не ввести».

## 6. `sim`

### 6.1. The choice at `sim` and the translation of a package

For architecture 2.10.

Read 2026-11-10: RUN.md [`sim`](https://github.com/altnovatrail/sim/blob/0364b40e961dbef87b07572baee03d408e8e090f/RUN.md), `main` 0364b40, the row В10 — «один вызов `engine.Choose` идёт около 30 мс на машине разработки», «В профиле `sim`, 2026-11-09: `engine.Choose` — около 60 % процессора, `prepared.candidates` и `probabilities` — его большая часть, сборка мусора — около 23 %», — and the outcome of T30 — «ученик олимпиадного графа — около 0,65 с вместо около 0,007 с». The translation of a package at `sim` is [`internal/pack/translate.go`](https://github.com/altnovatrail/sim/blob/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2/internal/pack/translate.go) `v0.1.2` e61f90a: `Template: st.Template != nil`, and `KeyProvenance` and `SkillTagsProvenance` — `Readout.Task` of `learner-model` v0.1.0 does not read them. `learner-model` v0.1.0, `readout.go` and `fold.go`, from the Go module cache: `Readout.Task` on each call checks the description of the task, parses its ID by `strings.SplitN`, and copies and sorts the items.

### 6.2. The gate of `sim`

For architecture 3.7; R17.

[sim `v0.1.0`](https://github.com/altnovatrail/sim/tree/09e5a84f8a26fedcf5c8996bf81bf29c616917a2), commit 09e5a84, an immutable release, read 2026-11-09: the README, the sections "Gate" and "Taking a release", — "`just gate` builds the bench on the commits of the libraries that gate.json names, through a `go.work` it writes in a temporary folder, runs the gate and holds its numbers to their bands […] A caller puts its own code in place of a library with `GATE_LEARNER_MODEL=<folder>` or `GATE_ENGINE=<folder>`"; the recipe `gate` in the `justfile`; `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `spec` v0.9.9, `go 1.27.3`, 65 bands; the job `gate` in `.github/workflows/ci.yml` — the model for the job "Gate"; architecture 3.9 sim — what the caller takes and why Dependabot does not raise the tag.

### 6.3. The gate of `sim` on the new pin

For architecture 3.7; the 2026-11-10 amendment to R17.

[sim `v0.1.2`](https://github.com/altnovatrail/sim/tree/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2), commit e61f90a, an immutable release of 2026-11-09T23:44:12Z, read 2026-11-10: `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `modules` v0.1.4 with `olympiad-math` 0.1.1, `spec` v0.9.9, `go 1.27.3`, 69 bands with the reason "a new pin, modules v0.1.4 with olympiad-math 0.1.1, A-linked on the two published maths, 100 learners a cell, beside the decisive run that fixed the paths anew (T30)"; `internal/gate/gate.go` — `const Learners = 100`; RUN.md `sim`, T30, — "Gate / amd64" 11 min in the CI of `sim`. The stage `toolchain` of `.devcontainer/Dockerfile` is the same as at `v0.1.0`.

### 6.4. The gate of `sim` after the choice was sped up

For architecture 3.7; the 2026-11-11 amendment to R17.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, an immutable release of 2026-11-11T02:54:16Z, read 2026-11-11: `gate.json` — `engine` `v0.1.6` 57e0ab8, `learner-model` `v0.1.4` acfc4f2, `modules` v0.2.0, `spec` v0.12.0, `go 1.27.3`; `internal/gate/gate.go` — `const Learners = 150`; `docs/decisions.md` of `sim`, the 2026-11-11 amendment to R31 sim, — «база гейта — `engine` v0.1.6, и в гейте 150 учеников на ячейку, а не 100», «на `engine` v0.1.6 гейт на 100 учениках читает все 69 чисел те же до бита, что на v0.1.4», «в пересчёте по самому медленному замеру на 150 учеников: у `sim` — около 8 мин, у `engine` — около 11,5, у `learner-model` — около 11,2»; the 2026-11-11 amendment to R33 sim — the decisive run on 1,500 learners. The stage `toolchain` of `.devcontainer/Dockerfile` is the same as at `v0.1.2`: the files differ only in the comments of the gopls stage.

### 6.5. The corridor of the base

For architecture 2.9.

[altnovatrail/sim](https://github.com/altnovatrail/sim/tree/v0.1.2), the tag `v0.1.2`, commit e61f90a, read 2026-11-10: `internal/measure/engine.go` — `E1` counts the items with a verdict by `PGiven`, P* at issue, with the corridor 0.70–0.85 and the thresholds 0.50 and 0.95; `internal/run/pool.go` — "Ceiling is the largest share of items with a verdict whose P* is in the corridor, of one task of the pool the rule could give"; `docs/decisions.md`, the 2026-11-09 amendment to R04 sim, the new pin («новое закрепление») — the table of paths: `E1.corridor` — the share outside the corridor, the less the better — for the base and for the ceiling: olympiad — `G0` 0.7347 and 0.0163, `G2-half` 0.7210 and 0.0110, `G2` 0.7313 and 0.0125; school — `G0` 0.8280 and 0.7510, `G2-half` 0.8113 and 0.6866, `G2` 0.8694 and 0.7001. A solution goes to the device «после верного ответа или показа» — R82 spec.

### 6.6. The truth of `sim` about an anchor

For architecture 2.11.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, read 2026-11-11: `internal/run/pool.go`, `holds`, — "an anchor is held when a skill of the learner's modules that teaches it is truly mastered with its R* then, or when none teaches it"; `internal/measure/model.go`, `m9`, — "the share of readings "holds" an anchor skill when the truth does not hold it, and of "does not hold" when it does", and the count goes only for `AnchorHas` and `AnchorHasNot`; `internal/learner/assumptions.go` — "TeacherBoost makes the skills that teach an anchor skill stronger by that much: the learner knows the anchor outside the platform (A-unpractised)".

### 6.7. A lesson at `sim`

For architecture 2.12.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), commit a97be52, read 2026-11-11: `internal/run/run.go` — the choice of a lesson is played by `PassLesson` (`internal/learner/attempt.go`) in one step, and `engine.Choose` is called once per step; `internal/measure/engine.go`, `e1`, — the items of each `Obs` of a step, the tasks of a lesson too. On `main` d957fc7 there are tags up to `v0.1.8`, and `internal/run`, `internal/learner` and `internal/measure` are the same as at `v0.1.6` (`git diff --stat v0.1.6 d957fc7`).

### 6.8. The simulator for the criterion of consolidation in a lesson

For architecture 2.12, "The criterion of the judgement".

[sim `v0.1.9`](https://github.com/altnovatrail/sim/tree/e1d5e8e33d90f2c3a88634a192988f23b6e56aaf), commit e1d5e8e, `main`, read 2026-11-12 from a clone in a temporary folder:

- `docs/architecture.md`, section 3.2: `G2` — «растёт: θ* — на 0,01 за ответ, δ* основных навыков ответа — ещё на 0,02, а δ* остальных навыков возвращается на 1 % пути к начальной»; section 3.1 — «Базовая популяция не учится, как G0»; section 3.5, `I2`, — «„закрепление“ — прошлая попытка на этом навыке была с неверным ответом, и последнее свидетельство модуля уже не в пробной серии»; section 3.6 — the criterion of the engine and the group "Corridor" («Коридор»);
- `RUN.md`, E5, «Общее для задач»: «одно закрепление на круг: от тега T46 до выбора суда `model-2` `go.mod`, `gate.json`, пакеты и пути не меняются»; table 2, E6 — the judgement `engine-2` «на базе после суда `model-2` … с подтверждением на нетронутом `held-out-engine-1`»;
- `internal/run/run.go`, `Step`: the choice of a lesson is played by `PassLesson` in one step, then each skill of the lesson gets an issue at the time of the choice (`is.Skill, is.Issued, is.LastIssued = s, true, at`), and the module's counter grows by one (`r.choices[m.ID]++`) — the same as after a task;
- `internal/published/published.go`: the release `modules` `v0.2.5`, `olympiad-math` and `grade-3-math` 0.3.0, `english-essays` 0.2.0. `jq` over the `server.json` of their archives (`v0.2.1/sha256-e4de172e….tar.gz`, `v0.2.1/sha256-dd2c6eb2….tar.gz`, `v0.2.0/sha256-85681ea4….tar.gz`): `olympiad-math` 0.3.0 has 41 lessons, each with one skill and one band and three task segments, the difficulties 1, 3, 5 in 37 lessons, and 1, 2, 4, 1, 2, 5, 1, 4, 5 and 2, 3, 5 in the rest; 775 tasks, 123 in lessons, and outside lessons per skill and band of a lesson from 2 to 27. `grade-3-math` 0.3.0 has 4 lessons of three task segments each, of the difficulties 1, 3, 5, and all 12 tasks are in lessons; `english-essays` 0.2.0 has 6 lessons, and all 8 tasks are in lessons.

## 7. `.github`

### 7.1. The shared workflows of `.github` v0.1.15

For architecture 3.7; R12.

The release [v0.1.15](https://github.com/altnovatrail/.github/releases/tag/v0.1.15), commit 89d3cee, published 2026-11-02, read 2026-11-08; `gh release verify v0.1.15 -R altnovatrail/.github` — "Release v0.1.15 verified!"; there is no newer release (`gh release list -R altnovatrail/.github`):

- `go.yml` — "ci-test writes the coverage of the tests as one Go coverage profile, coverage/go.out, and the job coverage sends it to Codecov by OIDC, with no secret"; `ci-licenses` gets `ALLOWED_LICENSES` — "the ASF's Category A as SPDX identifiers"; the calling job gives it `contents: read` and `id-token: write`;
- `codeql.yml` — "Besides actions, javascript-typescript and go are taken", and Go is read from the Go build of the stage `toolchain`;
- `secrets.yml`, `dco.yml`, `documents.yml` and `scorecard.yml` at v0.1.15 are the same as at v0.1.0 — the commit ff2a359 that the links here stood on: `git diff --stat v0.1.0 v0.1.15` over these four files is empty (2026-11-08).

### 7.2. The shared workflows of `.github` v0.1.21 and `same-bits`

For architecture 3.7; the fourth 2026-11-08 amendment to R10.

2026-11-08: the release [v0.1.21](https://github.com/altnovatrail/.github/releases/tag/v0.1.21), commit `706f4067a66c91b0442033f82a7924eeb8b60497`, published 2026-11-08T13:54:51Z; `gh release verify v0.1.21 -R altnovatrail/.github` — "Release v0.1.21 verified!":

- `go.yml` — the input `same-bits`: "With same-bits, four more jobs hold the caller's bits the same on every build (…), as learner-model and engine held them each in a recipe of its own"; ""Go / targets" runs the tests on the builds where the bits are most at risk, and builds and vets wasip1; "Go / fma-diff" compares the assembler of each package with fusion and without it, with a control of its own here; "Go / arm64" runs ci-test and the targets of its host on an arm64 runner; "Go / rules" holds the code to the rules of code by the tool tools/coderules of .github, of this workflow's own version, and the caller's lists, .github/coderules.json";
- `tools/coderules` is the module `github.com/altnovatrail/.github/tools/coderules`; the lists are `code-imports`, `test-imports`, `test-math`, `requires` and `nested-modules`, with no field beyond them; the rules are the same as those of the port `rules_test.go`, and `nested`. On the code of `engine` with the lists of `.github/coderules.json` — no findings;
- between v0.1.15 and v0.1.21, `secrets.yml`, `dco.yml`, `codeql.yml` and `scorecard.yml` differ only in the versions of actions (`git diff v0.1.15 v0.1.21`); `documents.yml` builds the stage by the script `spec/tools/toolchain/toolchain.sh`, which `spec` gives since v0.2.17.

### 7.3. The shared release of `.github` v0.1.21

For architecture 3.6 and 3.7; R15.

Read 2026-11-08: `release.yml` of the same commit 706f4067 — "Nobody tags by hand: the version is the next patch of the line the caller's release workflow sets, RELEASE_LINE"; "A commit gets a release only when a path the caller names in paths has changed since the release before it; the first release comes out always"; the caller's recipe `release-files` "gets the version, has no git to ask, and writes the files into dist/"; the caller "keeps its releases immutable and its tags v* under a ruleset". The model is R16 learner-model and its `release.yml` and `justfile` (`main`, read 2026-11-08): the release file is the archive of the schema, the ruleset `v-tags` — id 26770474, `refs/tags/v*`, `update`, `deletion`, `non_fast_forward`, without bypass. At `engine` on 2026-11-08 — `immutable-releases` `enabled: false`, a ruleset only for `main`, no tags (`gh api`).

## 8. Go, GitHub and Codecov

### 8.1. The fusion of floating-point operations

For architecture 3.7.

The Go specification, "Language version go1.27 (May 26, 2026)", the file `doc/go_spec.html` of go1.27.3 in the container, read 2026-11-07: "An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually. An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding"; on maps — "The iteration order over maps is not specified". An experiment of the same day, go1.27.3, `go build -gcflags=-S`: `(p1+p2)/2 - 0.775` — `FNMSUBD` on arm64, `FMSUB` on ppc64le and s390x, `FMSUBD` on riscv64 and loong64; `float64((p1+p2)/2) - 0.775` — a separate multiplication and subtraction on all five; on amd64 with GOAMD64 v3 — in neither case. `s*0.5 - 0.775` — `FNMSUBD` on arm64, `FMSUB` on ppc64le and `FMSUBD` on riscv64, and `s/float64(n) - 0.775` — `FDIVD` and `FSUBD` on arm64 and riscv64, `FDIV` and `FSUB` on ppc64le: a division by a variable does not fuse.

### 8.2. The wasm ports

For architecture 3.7.

`go tool dist list -json` of go1.27.3 in the container, 2026-11-07: `js`/`wasm` and `wasip1`/`wasm` have `"FirstClass": false`. [Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly), read 2026-11-07: "At present, Go generates large Wasm files, with the smallest possible size being around ~2MB"; "Using `gz` compression reduces the ~2MB (minimum file size) example WASM file down to around 500kB"; "To build for WASI, use the `wasip1` port"; on TinyGo — "TinyGo supports a subset of the Go language targeted for embedded devices".

### 8.3. The arm64 runners

For architecture 3.7.

[GitHub Docs, GitHub-hosted runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), read 2026-11-07: "Use of the standard GitHub-hosted runners is free and unlimited on public repositories"; in the table of the standard runners for public repositories — "Linux | 4 | 16 GB | 14 GB | arm64 | `ubuntu-24.04-arm`, `ubuntu-22.04-arm`, `ubuntu-26.04-arm`". `altnovatrail/engine` is public (`gh repo view`, 2026-11-07).

### 8.4. The tools of the checks and gopls

For architecture 3.7; R12, R13.

The Go module proxy, `https://proxy.golang.org/<module>/@latest`, read 2026-11-08: `golang.org/x/vuln` — v1.9.0 of 2026-10-11, `github.com/google/go-licenses/v2` — v2.0.1 of 2025-09-08, `golang.org/x/tools/gopls` — v0.24.0 of 2026-10-14. The licenses — by GitHub's API, `repos/<repository>/license`, 2026-11-08: govulncheck and gopls — BSD-3-Clause (`golang/vuln` and `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`). The template of `spec` v0.9.8 gives no gopls to the kind `go`: in `repo-template/kinds.json` it has only the extension `golang.go@0.56.1`.

### 8.5. Codecov

For architecture 3.7; R12.

Read 2026-11-08: the Codecov app is installed in the organization for all repositories — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `engine` is already in Codecov, but not activated — `https://api.codecov.io/api/v2/github/altnovatrail/repos/engine/`, `"active":false,"activated":false`.
