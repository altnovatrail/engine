# engine — источники

> Внешние факты, на которых стоят документы `docs/` репозитория: правило и код standalone, исследования и документация других программ, документы, код и файлы релизов других репозиториев NovaTrail, документация Go, GitHub и Codecov. Здесь — их записи: адрес, версия или дата страницы, дата прочтения и цитата на языке источника, — а [architecture.md](architecture.md) держит утверждение и ссылку на запись по номеру, как на раздел: `sources 1.1` (RULES.md, The bar; R143 spec).
>
> Запись — подраздел «N.M» темы: один источник или одно чтение нескольких источников для одного решения. Номер записи не меняется: новая встаёт в конец своей темы, новая тема — в конец файла. Источник, прочитанный снова, — пункт с датой в той же записи, а прежний текст не правится. Образец — sources spec. Итоги замеров — не здесь, а в [measurements.md](measurements.md).
>
> Записи 1.1–8.5 перенесены 2026-11-12 дословно из раздела 4 [architecture.md](architecture.md), «Источники», задачей T55 (R24). Первая строка записи называет подразделы architecture и решения, для которых она прочитана, если пункт их называл. Раздел 4 остался указателем: на него ссылаются журнал и архивы, а их не переписывают.

## 1. MathTrail: standalone

### 1.1. Правило выбора задачи standalone

Для architecture 1, 2.1–2.3.

[mathtrail/mathtrail-standalone, `SPEC.md`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/SPEC.md), ветка `main`, коммит c6cbeb2, прочитано 2026-10-19:

- SPEC 3.1: «It is deterministic by design (the prototype's D40): the same profile always produces the same brief»; «Nor does the rule read the links between topics (1.2)». Поля постановки: `pedagogical_goal` — «`reinforce` or `new_topic`», `rationale` — «One sentence naming what the choice was based on, and the model's reason when the model overrode the choice»;
- SPEC 3.2, шаг 5: «If every topic within reach is mastered, the same choice runs over all of them»;
- SPEC 2.5, прочитано 2026-10-20: «The topic counts as mastered only while its recommended point (2.3) stays at that level or below: once the child's level in the topic moves the recommended point up a level, the topic is back in the rotation»; «Mastery is lost after two consecutive wrong answers in that topic, its `wrong_streak`: both fields are cleared, at whatever level mastery was held, and the topic returns to the rotation»; «A mastered topic is skipped when the rule looks for something new, until every topic within reach is mastered — after which all of them are back in the rotation»;
- SPEC 2.3: «A topic is **within reach** when its easiest point — difficulty 1 of the lowest level it is taught at — has P ≥ 0.70»; тема на паузе «comes back once in six tasks as a probe at that point»;
- SPEC 3.2, «The idea»: идея задачи «is drawn among the rest by a hash of `student_id`, the topic, the level, the round and how many ideas the record of the list holds»;
- SPEC 3.4: «The model may ask for a different topic, level or difficulty, and must give a reason»; «`open_request.tutor_mode` becomes `llm` instead of `rule`». SPEC 3.6: когда уроки держатся темы, которую выбрал человек, «`open_request.tutor_mode` is **`person`**».

### 1.2. Код правила standalone

Тот же коммит c6cbeb2, что у sources 1.1, прочитано 2026-10-19:

- выбор — пакет [`internal/domain/tutor`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/tutor): «Package tutor decides what the child is asked next»; «Nothing here is random and nothing here is a model». О выборе идеи в `idea.go`: «looks drawn at random — another order for every child and every round — and is not: it is read off a hash of the child's id». К `tutor` NovaTrail этот пакет отношения не имеет: у него только имя то же. Рядом лежат его тесты, в том числе `golden_test.go`, `property_test.go` и `simulation_test.go`;
- коридор, досягаемость и «освоено» — пакет [`internal/domain/rating`](https://github.com/mathtrail/mathtrail-standalone/tree/c6cbeb2/internal/domain/rating), файлы `corridor.go`, `ladder.go` и `mastery.go`. В `corridor.go`: `CorridorLow  = 0.70`, `CorridorHigh = 0.85`, `CorridorMiddle = (CorridorLow + CorridorHigh) / 2`, а над последней — «CorridorMiddle is what the recommendation aims at, and the line an answer has to be at or below to count toward mastering a topic» (прочитано 2026-10-20).

### 1.3. Правило выбора standalone, дополнительно для T01

Тот же коммит c6cbeb2, прочитано 2026-11-07:

- SPEC 2.3: тема «on pause» — «comes back once in six tasks as a probe at that point, however small its chance there, whose right answer lifts it into reach»; SPEC 3.2, шаг 4: «those never handed out come first, in catalog order; then the one whose `topics[t].last_issued` is the oldest», и тема, в которой ребёнок «behind», считается дважды (R300 standalone);
- SPEC 3.2: «A child below every topic of the catalog […] is given instead the topics whose easiest point is the easiest the catalog has»; шаг 1: «In the trial series […] a failure is not worked over»;
- [`internal/domain/tutor/pause.go`](https://github.com/mathtrail/mathtrail-standalone/blob/c6cbeb2/internal/domain/tutor/pause.go): `const ProbeEvery = 6`, ход пробы — `tasks > rating.TrialAnswers && tasks%ProbeEvery == ProbeEvery-1`.

### 1.4. Тема, выбранная человеком, в standalone

Для architecture 3.2.

SPEC 3.6 standalone, тот же коммит c6cbeb2, прочитано 2026-11-07: «The lessons may be kept to one topic of the catalog»; выбор пишет `student.lesson_topic`, «and an empty one gives the choice back to the rule»; «Once the trial series is over, and while the catalog has the topic, the rule sets it in place of its own»; «During the trial series the choice is kept and waits».

## 2. Исследования и другие программы

### 2.1. Коридор около 85 % верных

Wilson R. C., Shenhav A., Straccia M., Cohen J. D. The Eighty Five Percent Rule for optimal learning. *Nature Communications* 10, 4646 (2019), [doi:10.1038/s41467-019-12552-4](https://doi.org/10.1038/s41467-019-12552-4), по [PMC6831579](https://pmc.ncbi.nlm.nih.gov/articles/PMC6831579/), прочитано 2026-11-07: «the optimal error rate for training is around 15.87% or, conversely, that the optimal training accuracy is about 85%». Это верхняя граница коридора standalone. Нижнюю, около 0,75, SPEC 2.3 standalone берёт у Math Garden («Math Garden keeps success near 0.75»); первоисточник — статью Math Garden — прочитать не удалось: издатель отвечает 403, и здесь это не утверждается как факт.

### 2.2. Повторения и новое в Anki

Для architecture 3.4.

[руководство Anki, Deck Options](https://docs.ankiweb.net/deck-options.html), прочитано 2026-11-07. Лимит повторений ограничивает и новое, при долге новое советуют не вводить, а порядок нового и повторений — настройка: «By default, the review limit also applies to new cards, and no new cards will be shown when the review limit has been reached»; «When this limit is reached, Anki will not show any more review cards for the day, even if there are more waiting»; порядок новых карточек и повторений — «Whether new cards should be mixed in with review cards, shown before them, or shown after them»; «If you have a backlog of overdue review cards, it is recommended that you stop introducing new cards until you catch up with that backlog».

## 3. `modules`

### 3.1. Графы первого среза

[altnovatrail/modules](https://github.com/altnovatrail/modules/tree/6d061e9/modules), `main` 6d061e9, прочитано 2026-11-07: `olympiad-math/skill-graph.yaml`, `grade-3-math/skill-graph.yaml`, паспорта с полосами и файлы `tasks/` — пример architecture 3.1.

### 3.2. Навыки-опоры у модулей первого среза

Для architecture 3.2.

`altnovatrail/modules`, `main` 6d061e9, прочитано 2026-11-07: паспорт `olympiad-math/manifest.yaml` — «Its skills need no anchor skill yet»; в [docs/architecture.md](https://github.com/altnovatrail/modules/blob/6d061e9/docs/architecture.md) `modules`, таблица порта, — «у `number-divisibility` есть ещё навык-опора `core.division-with-remainder`, а у `fractions-parts` — `core.fractions-of-a-whole`» и «Остальные 14 навыков вносит порт»; `grade-3-math/skill-graph.yaml` — `fractions-of-a-whole` с `teaches: [core.fractions-of-a-whole]`. Граф порта — [пример формата](../spec/format/examples/skill-graph/valid/olympiad-math.yaml) `olympiad-math.yaml`, 17 навыков.

### 3.3. Уроки первого среза

Для architecture 3.3.

`altnovatrail/modules`, `main` 6d061e9, прочитано 2026-11-07, папки `lessons/` и `tasks/` трёх модулей:

- олимпиадная математика — девять уроков, по уроку на навык и полосу, `counting-gaps-grades-3-4` и так далее; урок `counting-gaps-grades-3-4` — объяснения, пример с разбором и задачи `gaps-34-d1-1`, `gaps-34-d3-1`, `gaps-34-d5-1`: все три задачи навыка на полосе;
- школьная математика 3 класса — четыре урока, по уроку на навык; урок `times-tables` — «an explanation and the three pages of the skill, in the order of their difficulty»;
- сочинения — шесть уроков и восемь задач, все задачи — в уроках. `reading-model-essays` — «explanations alone, with no task»; `first-sentences` — «on two skills […]: choosing a first sentence, which gets a verdict, and writing one, which gets feedback by a rubric»; `writing-about-a-day` — «for adults alone», `age_categories: [adult]`, навык `endings`; `endings` — «the engine's rule chooses between it and the lesson of improving an essay».

### 3.4. Пакеты релиза `modules` v0.1.0 в `testdata/modules`

Для architecture 3.7; третья поправка 2026-11-08 к R10.

2026-11-08: релиз [v0.1.0](https://github.com/altnovatrail/modules/releases/tag/v0.1.0) — единственный у `altnovatrail/modules` (`gh release list -R altnovatrail/modules`), тег на коммите `6d061e940808b1ee9f64fe9a30bd209d408a756e` (`gh api repos/altnovatrail/modules/git/ref/tags/v0.1.0`), том же, по которому написаны примеры 3.1–3.3. `just -f spec/justfile release-verify modules v0.1.0 <файл> <файл>.sigstore.json` принял каждый из четырёх файлов — «modules v0.1.0, signed by the shared release.yml of .github v0.1.15»:

- `module-index.json`, sha256 `ed4d8bb325fdcb3ea54da439e08c58d927f773ba0f097338f7767ad76be8c90e`, — `"module_format_version":"0.3"` и у трёх модулей версии `0.1.0` дайджест пакета;
- `sha256-1d8a962d….tar.gz` — `novatrail/olympiad-math`, `sha256-2e2f7a2f….tar.gz` — `novatrail/grade-3-math`, `sha256-8c777871….tar.gz` — `novatrail/english-essays`; в каждом `inventory.json`, `server.json`, `device.json` и файлы устройства;
- дайджест пакета — sha256 его `inventory.json`, а не архива: «`package_digest`: the digest of the package, which is the digest of its inventory» ([текст формата](../spec/format/README.md), раздел 26); `inventory.json` называет sha256 `server.json` в `server_part.digest` — у трёх пакетов они совпали с sha256 файлов.

### 3.5. Индекс `modules` v0.1.4 и пакет `olympiad-math` 0.1.1 в `testdata/modules`

Для architecture 2.10 и 3.7; поправка 2026-11-10 к R10.

2026-11-10: релиз [v0.1.4](https://github.com/altnovatrail/modules/releases/tag/v0.1.4), тег на коммите `c23ed6267bc5854e7e728b23cb6541db4d15f38a`, и релиз [v0.1.1](https://github.com/altnovatrail/modules/releases/tag/v0.1.1), тег на коммите `5f6d7a9668ae172e36ca349e3f89e96faa6597ab` (`gh api repos/altnovatrail/modules/git/ref/tags/<тег>`). `just -f spec/justfile release-verify modules v0.1.4 module-index.json` — «module-index.json: modules v0.1.4, signed by the shared release.yml of .github v0.1.15»; `just -f spec/justfile release-verify modules v0.1.1 sha256-123cb0d3….tar.gz` — «sha256-123cb0d3….tar.gz: modules v0.1.1, signed by the shared release.yml of .github v0.1.15»; бандлы — `<файл>.sigstore.json` своих релизов:

- `module-index.json`, sha256 `a244c0ba2832ba6360f45bc74d293572b84b491fc823c7b01bdf6004ff9edb2a`, — `"module_format_version":"0.3"`; записи `0.1.0` трёх модулей — те же, что у индекса v0.1.0, побайтно, и ещё `novatrail/olympiad-math` `0.1.1` — `"package_digest":"sha256:123cb0d3dec7a094dce8e6c62f3b5f0f8e4563b23d148793fe3b9ce6c4839cb1"`, `"release":"v0.1.1"`;
- `sha256-123cb0d3….tar.gz`, 2 477 140 байт: в `server.json` — 19 навыков, все с `verdict`, три полосы, 772 задачи, из них 19 шаблонов — задач с полем `template`, — и 41 урок; виды проверок — `choice`, `number` и `fraction`; своих возрастных категорий у задач и уроков нет, а у модуля — `under_13` и `adult`; цепочка сошлась, как у пакетов 0.1.0: sha256 `inventory.json` — дайджест пакета в индексе, а sha256 `server.json` — его `server_part.digest`.

### 3.6. Уроки и задачи пакетов `testdata/modules`

Для architecture 2.12.

2026-11-11, `jq` по `server.json` каждого пакета во временной папке:

- `olympiad-math` 0.1.1, `sha256-123cb0d3….tar.gz`: 41 урок; у каждого — один навык, задачи одной полосы и три звена-задачи; сложность звеньев по порядку — 1, 3, 5 у 37 уроков, 1, 2, 4, 1, 2, 5, 1, 4, 5 и 2, 3, 5 — по одному; урок `arithmetic-tricks-grades-1-2` — задачи `tri-12-d1-1`, `tri-12-d3-1`, `tri-12-d5-1`. Задач в уроках — 123 разных из 772; вне уроков — 649, из них 19 шаблонов, и на навык и полосу урока — от 2 до 26;
- `olympiad-math` 0.1.0 — 27 задач и 9 уроков по три звена-задачи, `grade-3-math` 0.1.0 — 12 задач и 4 урока по три, `english-essays` 0.1.0 — 8 задач и 6 уроков: все задачи — в уроках.

## 4. `api`

### 4.1. Что движок получает от `api`

Для architecture 3.2.

[api, `docs/architecture.md`](https://github.com/altnovatrail/api/blob/06dad44/docs/architecture.md) и [`docs/decisions.md`](https://github.com/altnovatrail/api/blob/06dad44/docs/decisions.md), `main` 06dad44, прочитано 2026-11-07: «задачу даёт только ученику её возрастных категорий, `age_categories` (R51 spec), а урок — только если ученику подходят и урок, и каждая его задача»; об уроке, который не подходит, — «движок таких уроков и не получает»; поправка 2026-11-02 к R07 api — «Движок, как и прежде, получает графы текущих версий подключённых модулей, которые ученику подходят».

### 4.2. Урок у `api`

Для architecture 2.12.

[api, `docs/architecture.md`](https://github.com/altnovatrail/api/blob/2024d52015965ad018c54f07ca2d891da5e7b42c/docs/architecture.md), раздел 3.7, и [`docs/decisions.md`](https://github.com/altnovatrail/api/blob/2024d52015965ad018c54f07ca2d891da5e7b42c/docs/decisions.md), R13, R48 и R50 api, `main` 2024d52, прочитано 2026-11-11: «Попытки закрепления в уроке первый срез не знает (R107 spec); если её введёт новая версия движка, её место — «после звена k»»; «Введёт их новая версия движка — `api` будет звать его на конце звена урока, чьи навыки все с вердиктом, а ответ — следующее звено, закрепление или «досрочно»»; о «досрочно» — «До тех пор `api` урок досрочно не кончает, а среди причин конца в базе и в коде её нет». R48 api — `OpenLesson`, `GetLessonRun` и `CompleteLesson` в `internal/learning`; R50 api — движок зовётся на «другое», `next` вне урока и `RunEnd`.

## 5. `learner-model`

### 5.1. Владение навыком-опорой

Для architecture 3.2.

[learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/98ba29f/docs/decisions.md), `main` 98ba29f, — R07 learner-model, прочитано 2026-11-07: «Держать ли урок при «нет свидетельств», решает `engine`»; отвергнуто «неизвестно», пока у навыков, которые учат, нет ответов, — «урок сначала открыт, а с первыми ответами закрывается».

### 5.2. Что `Read` отдаёт движку

[learner-model, `docs/architecture.md`](https://github.com/altnovatrail/learner-model/blob/98ba29f/docs/architecture.md), `main` 98ba29f, прочитано 2026-11-07, таблица операций: у навыка с практикой — «число свидетельств, время последней практики, θ + δ, v и удержание на t₀, ловушки с числом и временем последней»; у освоенного — «полоса освоения, осторожная вероятность с удержанием и срок повторения»; `Readout.Task` — «P каждого пункта на t₀». Истории выдач и признака конца пробной серии там на 98ba29f не было.

### 5.3. Типы `learner-model` в коде

Для architecture 2.4.

[learner-model, `types.go`](https://github.com/altnovatrail/learner-model/blob/6a6c843/types.go) и [`readout.go`](https://github.com/altnovatrail/learner-model/blob/6a6c843/readout.go), `main` 6a6c843, прочитано 2026-11-07: `Graph`, `Skill`, `Band`, `RetiredSkill`, `Task`, `TaskItem`, `Slot`, `Readout`, `ModuleReadout`, `SkillReadout`, `AnchorReadout` и `TaskReadout`; `Task` — «describes a task, the same for evidence and for a candidate task of the engine», а `readout.go` — «These are the types alone: Read and Readout.Task compute them from T20 on». В `ModuleReadout` — `HasEvidence`, `StartBand` и `Skills`, а признаков пробной серии в коде ещё нет. В [docs/architecture.md](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/architecture.md) `learner-model` — в разделе «Типы и кодировка» — они уже в таблице операций: `Read` отдаёт «у модуля — кончилась ли пробная серия и было ли последнее свидетельство ещё в ней», а «Уроки, возрастные категории и другое, что нужно только выбору, движок держит своими типами вокруг описания задачи».

### 5.4. Сборка и одни биты у `learner-model`

Для architecture 3.7.

[learner-model, `docs/architecture.md`](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/architecture.md), раздел «Сборка и одни биты», и [`docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/6a6c843/docs/decisions.md) — R09 learner-model с поправками, `main` 6a6c843, прочитано 2026-11-07: восемь правил кода, списки пакетов, рецепты `test-targets` и `fma-diff`, тест правил `rules_test.go` с контролем `testdata/coderules`, job «Go arm64» и о Т-10 — «`engine` для повтора выбора на любой сборке нужна та же дисциплина — его T05».

### 5.5. Чтение `learner-model` в коде

Для architecture 3.6.

[learner-model, `readout.go`](https://github.com/altnovatrail/learner-model/blob/fc444ba/readout.go), `main` fc444ba, T20, прочитано 2026-11-07: `func Read(s State, g Graphs, at int64, grade Grade, engineModules []string) (Readout, error)` и `func (r Readout) Task(t Task) (TaskReadout, error)`; у `ModuleReadout` — «TrialOver tells whether the module's trial series is over, and LastInTrial whether the last» свидетельство было ещё в ней, поля `TrialOver` и `LastInTrial`. Первый подписанный релиз — [`v0.1.0`](https://github.com/altnovatrail/learner-model/releases/tag/v0.1.0), тег на коммите d2620a1, опубликован 2026-11-07T20:12:08Z: `gh release verify v0.1.0 -R altnovatrail/learner-model` — «Release v0.1.0 verified!» (2026-11-08). В код `engine` он входит, только когда его примет `just -f spec/justfile release-verify` (R77 spec).

### 5.6. Каркас и gopls у `learner-model`

Для architecture 3.7; R12, R13.

[learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/50f1b51/docs/decisions.md), R13 и R14 learner-model с поправками, и `justfile`, `.github/`, `.devcontainer/` того же коммита, `main` 50f1b51, прочитано 2026-11-08: образец рецептов `ci-*`, job «Go», `codecov.yml` и шага gopls. Из поправок к R13 learner-model: на pull request без строк кода Codecov не ставит `codecov/patch` — «Codecov видит у него 0 файлов и 0 строк […] и статуса не ставит, так что мерж ждал его без конца», — а первый статус пришёл с первым кодом под тестами.

### 5.7. `learner-model` `v0.1.0` в коде

Для architecture 3.6; поправка 2026-11-08, вторая, к R09.

2026-11-08: `gh release download v0.1.0 -R altnovatrail/learner-model` и `just -f spec/justfile release-verify learner-model v0.1.0 learner-model-proto-v0.1.0.tar.gz` — «learner-model-proto-v0.1.0.tar.gz: learner-model v0.1.0, signed by the shared release.yml of .github v0.1.15»; тег `v0.1.0` — на коммите `d2620a18d865c813ec2556c93fca92a1445e3a74` (`gh api repos/altnovatrail/learner-model/git/ref/tags/v0.1.0`); `go mod download -json github.com/altnovatrail/learner-model@v0.1.0` через `proxy.golang.org` — `"Origin":{"Hash":"d2620a18d865c813ec2556c93fca92a1445e3a74","Ref":"refs/tags/v0.1.0"}`, `"Sum":"h1:KK3u66KFWfWV8NmI+ueCXZQVXyUrg8zfDtiyo8+Me9Q="` — те же суммы, что `go get` записал в `go.sum`. Что модуль Go — не файл релиза и его держат `go.sum` и sum.golang.org, — R16 learner-model ([learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/50f1b51/docs/decisions.md), `main` 50f1b51, прочитано 2026-11-08). API — [`readout.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.0/readout.go) того же тега: у `Readout` есть невидимые снаружи поля, и строит его только `Read`.

### 5.8. `learner-model` `v0.1.6` в коде

Для architecture 2.10, 3.6 и 3.7; R23.

2026-11-11: `gh release list -R altnovatrail/learner-model` — релизы `v0.1.0`…`v0.1.6`, последний `v0.1.6`, опубликован 2026-11-11T17:32:38Z, `isImmutable` — `true`; `gh release download v0.1.6` и `just -f spec/justfile release-verify learner-model v0.1.6 learner-model-proto-v0.1.6.tar.gz` — «learner-model-proto-v0.1.6.tar.gz: learner-model v0.1.6, signed by the shared release.yml of .github v0.1.21», sha256 файла `2ec63cc68d2396cc96cee34c09281290bfc3e209cb5037f100c85fc55b789d5f`; `gh release verify v0.1.6 -R altnovatrail/learner-model` — «Release v0.1.6 verified!»; тег `v0.1.6` — на коммите `7d796e290599eeb15481f0098c9356bf6c31afe9`; `go mod download -json github.com/altnovatrail/learner-model@v0.1.6` через `proxy.golang.org` — Origin `refs/tags/v0.1.6` на том же коммите, `"Sum":"h1:xoifdkdFlXvTSgjFu939p2Ol2oB7s6G6pCO3h2GxfxU="` — та же сумма, что `go get` записал в `go.sum`; строка `go` её `go.mod` — `1.27.3`, `require` нет. Что изменилось с `v0.1.0` — `git diff v0.1.0 v0.1.6` того же репозитория: [`types.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/types.go) — `CheckKindTime`; [`evidence.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/evidence.go) — `EvidenceKind = 2` и `LeadingQuestionOnRequest`; [`readout.go`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/readout.go) — у `Read`: «a module of the engine with no state and a current version gives the starting point by the grade, BandByGrade, or none when that version has no bands», у `Readout.Task`: «It copies nothing it need not»; R28 и R29 learner-model и поправка 2026-11-09 к R07 learner-model — [learner-model, `docs/decisions.md`](https://github.com/altnovatrail/learner-model/blob/v0.1.6/docs/decisions.md) того же тега.

### 5.9. Навык-опора без учителя у модели и в пакете

Для architecture 2.11.

Прочитано 2026-11-11:

- пакет `olympiad-math` 0.1.1 в `testdata/modules`, `sha256-123cb0d3….tar.gz`, файл `server.json`: у `fractions-parts` — `"prerequisites":["arithmetic-tricks","core.fractions-of-a-whole"]`, у `number-divisibility` — `["time-calendar","parity-alternation","core.division-with-remainder"]`, у `percent-basic` и `ratio-sharing` — `["fractions-parts"]`; у 16 задач — `"skills":{"primary":["fractions-parts"],"supporting":["core.fractions-of-a-whole"]}`, у 16 — `"skills":{"primary":["number-divisibility"],"supporting":["core.division-with-remainder"]}`, и других задач у этих двух навыков нет;
- `learner-model` v0.1.0 из кэша модулей Go: `step.go`, `itemSkills` — свидетельство двигает основные навыки пункта или задачи, а вспомогательные не читает; `readout.go` — `AnchorUnknown` с причинами `AnchorReasonNoTeacher` и `AnchorReasonNoVerdictTeacher`; `docs/decisions.md`, R07 learner-model, — «навык-опора как вспомогательный навык задачи ничего не получает (R06), и модель его в свидетельстве не читает»;
- `spec` v0.15.0: R63 spec — «Ученику модуль предлагает `api` — только из тех, что держит сам и что ученику подходят»; [текст формата](../spec/format/README.md), раздел 26, — «The platform, for a module that a learner lacks […]: the modules that teach an anchor skill are those whose entry names it in `teaches`»; `progress.proto` — у `AnchorSkillState` поля `anchor_skill`, `holding`, `reason` и `teaching_skills`, «The full IDs of the skills that teach it, in the modules whose graphs the learner model has»; `module.proto` — у `CatalogModule` только `module_id` и `package_digest`; R100 spec, «Отвергнуто», — среди ветвей standalone «короткая проверка навыка-опоры», и о них — «движок может их не ввести».

## 6. `sim`

### 6.1. Выбор у `sim` и перевод пакета

Для architecture 2.10.

Прочитано 2026-11-10: RUN.md [`sim`](https://github.com/altnovatrail/sim/blob/0364b40e961dbef87b07572baee03d408e8e090f/RUN.md), `main` 0364b40, строка В10 — «один вызов `engine.Choose` идёт около 30 мс на машине разработки», «В профиле `sim`, 2026-11-09: `engine.Choose` — около 60 % процессора, `prepared.candidates` и `probabilities` — его большая часть, сборка мусора — около 23 %», — и итог T30 — «ученик олимпиадного графа — около 0,65 с вместо около 0,007 с». Перевод пакета у `sim` — [`internal/pack/translate.go`](https://github.com/altnovatrail/sim/blob/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2/internal/pack/translate.go) `v0.1.2` e61f90a: `Template: st.Template != nil`, а `KeyProvenance` и `SkillTagsProvenance` — их `Readout.Task` `learner-model` v0.1.0 не читает. `learner-model` v0.1.0, `readout.go` и `fold.go`, из кэша модулей Go: `Readout.Task` на каждый вызов проверяет описание задачи, разбирает её ID `strings.SplitN` и копирует и сортирует пункты.

### 6.2. Гейт `sim`

Для architecture 3.7; R17.

[sim `v0.1.0`](https://github.com/altnovatrail/sim/tree/09e5a84f8a26fedcf5c8996bf81bf29c616917a2), коммит 09e5a84, неизменяемый релиз, прочитано 2026-11-09: README, разделы «Gate» и «Taking a release», — «`just gate` builds the bench on the commits of the libraries that gate.json names, through a `go.work` it writes in a temporary folder, runs the gate and holds its numbers to their bands […] A caller puts its own code in place of a library with `GATE_LEARNER_MODEL=<folder>` or `GATE_ENGINE=<folder>`»; рецепт `gate` в `justfile`; `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `spec` v0.9.9, `go 1.27.3`, 65 полос; job `gate` в `.github/workflows/ci.yml` — образец job «Gate»; architecture 3.9 sim — что берёт вызывающий и почему Dependabot тег не поднимает.

### 6.3. Гейт `sim` на новом закреплении

Для architecture 3.7; поправка 2026-11-10 к R17.

[sim `v0.1.2`](https://github.com/altnovatrail/sim/tree/e61f90aa1bd894d1ce8f35547e0757d38f28c0b2), коммит e61f90a, неизменяемый релиз 2026-11-09T23:44:12Z, прочитано 2026-11-10: `gate.json` — `engine` `v0.1.0` c6c5bb9, `learner-model` `v0.1.0` d2620a1, `modules` v0.1.4 с `olympiad-math` 0.1.1, `spec` v0.9.9, `go 1.27.3`, 69 полос с причиной «a new pin, modules v0.1.4 with olympiad-math 0.1.1, A-linked on the two published maths, 100 learners a cell, beside the decisive run that fixed the paths anew (T30)»; `internal/gate/gate.go` — `const Learners = 100`; RUN.md `sim`, T30, — «Gate / amd64» 11 мин в CI `sim`. Стадия `toolchain` `.devcontainer/Dockerfile` — та же, что у `v0.1.0`.

### 6.4. Гейт `sim` после ускорения выбора

Для architecture 3.7; поправка 2026-11-11 к R17.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), коммит a97be52, неизменяемый релиз 2026-11-11T02:54:16Z, прочитано 2026-11-11: `gate.json` — `engine` `v0.1.6` 57e0ab8, `learner-model` `v0.1.4` acfc4f2, `modules` v0.2.0, `spec` v0.12.0, `go 1.27.3`; `internal/gate/gate.go` — `const Learners = 150`; `docs/decisions.md` `sim`, поправка 2026-11-11 к R31 sim, — «база гейта — `engine` v0.1.6, и в гейте 150 учеников на ячейку, а не 100», «на `engine` v0.1.6 гейт на 100 учениках читает все 69 чисел те же до бита, что на v0.1.4», «в пересчёте по самому медленному замеру на 150 учеников: у `sim` — около 8 мин, у `engine` — около 11,5, у `learner-model` — около 11,2»; поправка 2026-11-11 к R33 sim — решающий прогон на 1 500 учениках. Стадия `toolchain` `.devcontainer/Dockerfile` — та же, что у `v0.1.2`: файлы расходятся только в комментариях стадии gopls.

### 6.5. Коридор базы

Для architecture 2.9.

[altnovatrail/sim](https://github.com/altnovatrail/sim/tree/v0.1.2), тег `v0.1.2`, коммит e61f90a, прочитано 2026-11-10: `internal/measure/engine.go` — `E1` считает пункты с вердиктом по `PGiven`, P* при выдаче, с коридором 0,70–0,85 и порогами 0,50 и 0,95; `internal/run/pool.go` — «Ceiling is the largest share of items with a verdict whose P* is in the corridor, of one task of the pool the rule could give»; `docs/decisions.md`, поправка 2026-11-09 к R04 sim, «новое закрепление» — таблица путей: `E1.corridor` — доля вне коридора, меньше — лучше, — у базы и у потолка: олимпиадная — `G0` 0,7347 и 0,0163, `G2-half` 0,7210 и 0,0110, `G2` 0,7313 и 0,0125; школьная — `G0` 0,8280 и 0,7510, `G2-half` 0,8113 и 0,6866, `G2` 0,8694 и 0,7001. Решение уходит на устройство «после верного ответа или показа» — R82 spec.

### 6.6. Правда `sim` об опоре

Для architecture 2.11.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), коммит a97be52, прочитано 2026-11-11: `internal/run/pool.go`, `holds`, — «an anchor is held when a skill of the learner's modules that teaches it is truly mastered with its R* then, or when none teaches it»; `internal/measure/model.go`, `m9`, — «the share of readings "holds" an anchor skill when the truth does not hold it, and of "does not hold" when it does», и счёт идёт только у `AnchorHas` и `AnchorHasNot`; `internal/learner/assumptions.go` — «TeacherBoost makes the skills that teach an anchor skill stronger by that much: the learner knows the anchor outside the platform (A-unpractised)».

### 6.7. Урок у `sim`

Для architecture 2.12.

[sim `v0.1.6`](https://github.com/altnovatrail/sim/tree/a97be529249042595e544da6192555fa783ff9ec), коммит a97be52, прочитано 2026-11-11: `internal/run/run.go` — выбор урока играет `PassLesson` (`internal/learner/attempt.go`) одним шагом, а `engine.Choose` зовётся раз на шаг; `internal/measure/engine.go`, `e1`, — пункты каждого `Obs` шага, задачи урока тоже. На `main` d957fc7 — теги до `v0.1.8`, а `internal/run`, `internal/learner` и `internal/measure` те же, что у `v0.1.6` (`git diff --stat v0.1.6 d957fc7`).

### 6.8. Стенд для критерия закрепления в уроке

Для architecture 2.12, «Критерий суда».

[sim `v0.1.9`](https://github.com/altnovatrail/sim/tree/e1d5e8e33d90f2c3a88634a192988f23b6e56aaf), коммит e1d5e8e, `main`, прочитано 2026-11-12 по клону во временной папке:

- `docs/architecture.md`, раздел 3.2: `G2` — «растёт: θ* — на 0,01 за ответ, δ* основных навыков ответа — ещё на 0,02, а δ* остальных навыков возвращается на 1 % пути к начальной»; раздел 3.1 — «Базовая популяция не учится, как G0»; раздел 3.5, `I2`, — «„закрепление“ — прошлая попытка на этом навыке была с неверным ответом, и последнее свидетельство модуля уже не в пробной серии»; раздел 3.6 — критерий движка и группа «Коридор»;
- `RUN.md`, E5, «Общее для задач»: «одно закрепление на круг: от тега T46 до выбора суда `model-2` `go.mod`, `gate.json`, пакеты и пути не меняются»; таблица 2, E6 — суд `engine-2` «на базе после суда `model-2` … с подтверждением на нетронутом `held-out-engine-1`»;
- `internal/run/run.go`, `Step`: выбор урока играет `PassLesson` одним шагом, затем каждый навык урока получает выдачу на время выбора (`is.Skill, is.Issued, is.LastIssued = s, true, at`), а счётчик модуля растёт на один (`r.choices[m.ID]++`) — так же, как после задачи;
- `internal/published/published.go`: релиз `modules` `v0.2.5`, `olympiad-math` и `grade-3-math` 0.3.0, `english-essays` 0.2.0. `jq` по `server.json` их архивов (`v0.2.1/sha256-e4de172e….tar.gz`, `v0.2.1/sha256-dd2c6eb2….tar.gz`, `v0.2.0/sha256-85681ea4….tar.gz`): у `olympiad-math` 0.3.0 — 41 урок, у каждого один навык и одна полоса и три звена-задачи, сложности 1, 3, 5 у 37 уроков, 1, 2, 4, 1, 2, 5, 1, 4, 5 и 2, 3, 5 у остальных; 775 задач, в уроках 123, вне уроков на навык и полосу урока — от 2 до 27. У `grade-3-math` 0.3.0 — 4 урока по три звена-задачи сложностей 1, 3, 5, и все 12 задач в уроках; у `english-essays` 0.2.0 — 6 уроков, и все 8 задач в уроках.

## 7. `.github`

### 7.1. Общие workflow `.github` v0.1.15

Для architecture 3.7; R12.

Релиз [v0.1.15](https://github.com/altnovatrail/.github/releases/tag/v0.1.15), коммит 89d3cee, опубликован 2026-11-02, прочитано 2026-11-08; `gh release verify v0.1.15 -R altnovatrail/.github` — «Release v0.1.15 verified!»; более нового релиза нет (`gh release list -R altnovatrail/.github`):

- `go.yml` — «ci-test writes the coverage of the tests as one Go coverage profile, coverage/go.out, and the job coverage sends it to Codecov by OIDC, with no secret»; `ci-licenses` получает `ALLOWED_LICENSES` — «the ASF's Category A as SPDX identifiers»; вызывающий job даёт ему `contents: read` и `id-token: write`;
- `codeql.yml` — «Besides actions, javascript-typescript and go are taken», а Go читается из сборки Go стадии `toolchain`;
- `secrets.yml`, `dco.yml`, `documents.yml` и `scorecard.yml` у v0.1.15 те же, что у v0.1.0 — коммит ff2a359, на котором стояли ссылки здесь: `git diff --stat v0.1.0 v0.1.15` по этим четырём файлам пуст (2026-11-08).

### 7.2. Общие workflow `.github` v0.1.21 и `same-bits`

Для architecture 3.7; четвёртая поправка 2026-11-08 к R10.

2026-11-08: релиз [v0.1.21](https://github.com/altnovatrail/.github/releases/tag/v0.1.21), коммит `706f4067a66c91b0442033f82a7924eeb8b60497`, опубликован 2026-11-08T13:54:51Z; `gh release verify v0.1.21 -R altnovatrail/.github` — «Release v0.1.21 verified!»:

- `go.yml` — вход `same-bits`: «With same-bits, four more jobs hold the caller's bits the same on every build (…), as learner-model and engine held them each in a recipe of its own»; «"Go / targets" runs the tests on the builds where the bits are most at risk, and builds and vets wasip1; "Go / fma-diff" compares the assembler of each package with fusion and without it, with a control of its own here; "Go / arm64" runs ci-test and the targets of its host on an arm64 runner; "Go / rules" holds the code to the rules of code by the tool tools/coderules of .github, of this workflow's own version, and the caller's lists, .github/coderules.json»;
- `tools/coderules` — модуль `github.com/altnovatrail/.github/tools/coderules`; списки — `code-imports`, `test-imports`, `test-math`, `requires` и `nested-modules`, без поля сверх них; правила — те же, что у порта `rules_test.go`, и `nested`. На коде `engine` со списками `.github/coderules.json` — без находок;
- у `secrets.yml`, `dco.yml`, `codeql.yml` и `scorecard.yml` между v0.1.15 и v0.1.21 — только версии действий (`git diff v0.1.15 v0.1.21`); `documents.yml` строит стадию скриптом `spec/tools/toolchain/toolchain.sh`, который `spec` даёт с v0.2.17.

### 7.3. Общий релиз `.github` v0.1.21

Для architecture 3.6 и 3.7; R15.

Прочитано 2026-11-08: `release.yml` того же коммита 706f4067 — «Nobody tags by hand: the version is the next patch of the line the caller's release workflow sets, RELEASE_LINE»; «A commit gets a release only when a path the caller names in paths has changed since the release before it; the first release comes out always»; рецепт `release-files` вызывающего «gets the version, has no git to ask, and writes the files into dist/»; вызывающий «keeps its releases immutable and its tags v* under a ruleset». Образец — R16 learner-model и её `release.yml` и `justfile` (`main`, прочитано 2026-11-08): файл релиза — архив схемы, ruleset `v-tags` — id 26770474, `refs/tags/v*`, `update`, `deletion`, `non_fast_forward`, без обхода. У `engine` 2026-11-08 — `immutable-releases` `enabled: false`, ruleset только `main`, тегов нет (`gh api`).

## 8. Go, GitHub и Codecov

### 8.1. Слияние операций с плавающей точкой

Для architecture 3.7.

Спецификация Go, «Language version go1.27 (May 26, 2026)», файл `doc/go_spec.html` go1.27.3 в контейнере, прочитано 2026-11-07: «An implementation may combine multiple floating-point operations into a single fused operation, possibly across statements, and produce a result that differs from the value obtained by executing and rounding the instructions individually. An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion that would discard that rounding»; о картах — «The iteration order over maps is not specified». Опыт того же дня, go1.27.3, `go build -gcflags=-S`: `(p1+p2)/2 - 0.775` — `FNMSUBD` на arm64, `FMSUB` на ppc64le и s390x, `FMSUBD` на riscv64 и loong64; `float64((p1+p2)/2) - 0.775` — отдельные умножение и вычитание на всех пяти; на amd64 при GOAMD64 v3 — ни там, ни там. `s*0.5 - 0.775` — `FNMSUBD` на arm64, `FMSUB` на ppc64le и `FMSUBD` на riscv64, а `s/float64(n) - 0.775` — `FDIVD` и `FSUBD` на arm64 и riscv64, `FDIV` и `FSUB` на ppc64le: деление на переменное не сливается.

### 8.2. Порты wasm

Для architecture 3.7.

`go tool dist list -json` go1.27.3 в контейнере, 2026-11-07: у `js`/`wasm` и у `wasip1`/`wasm` — `"FirstClass": false`. [Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly), прочитано 2026-11-07: «At present, Go generates large Wasm files, with the smallest possible size being around ~2MB»; «Using `gz` compression reduces the ~2MB (minimum file size) example WASM file down to around 500kB»; «To build for WASI, use the `wasip1` port»; о TinyGo — «TinyGo supports a subset of the Go language targeted for embedded devices».

### 8.3. Раннеры arm64

Для architecture 3.7.

[GitHub Docs, GitHub-hosted runners](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), прочитано 2026-11-07: «Use of the standard GitHub-hosted runners is free and unlimited on public repositories»; в таблице стандартных раннеров для публичных репозиториев — «Linux | 4 | 16 GB | 14 GB | arm64 | `ubuntu-24.04-arm`, `ubuntu-22.04-arm`, `ubuntu-26.04-arm`». `altnovatrail/engine` — публичный (`gh repo view`, 2026-11-07).

### 8.4. Инструменты проверок и gopls

Для architecture 3.7; R12, R13.

Прокси модулей Go, `https://proxy.golang.org/<модуль>/@latest`, прочитано 2026-11-08: `golang.org/x/vuln` — v1.9.0 от 2026-10-11, `github.com/google/go-licenses/v2` — v2.0.1 от 2025-09-08, `golang.org/x/tools/gopls` — v0.24.0 от 2026-10-14. Лицензии — по API GitHub, `repos/<репозиторий>/license`, 2026-11-08: govulncheck и gopls — BSD-3-Clause (`golang/vuln` и `golang/tools`), go-licenses — Apache-2.0 (`google/go-licenses`). Заготовка `spec` v0.9.8 gopls виду `go` не даёт: в `repo-template/kinds.json` у него только расширение `golang.go@0.56.1`.

### 8.5. Codecov

Для architecture 3.7; R12.

Прочитано 2026-11-08: приложение Codecov стоит у организации на все репозитории — `gh api orgs/altnovatrail/installations`, `"app_slug":"codecov","repository_selection":"all"`; `engine` в Codecov уже есть, но не включён — `https://api.codecov.io/api/v2/github/altnovatrail/repos/engine/`, `"active":false,"activated":false`.
