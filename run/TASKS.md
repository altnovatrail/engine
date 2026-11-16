# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic


#### T66. Гейт частями на теге `sim` строки 7.19 — вне эпика

- **Цель:** job «Gate» — частями: матрица из нескольких job, каждый гонит свою часть гейта `sim`, и проверка после них; что гейт доказывает, не меняется (строка 7.20 COORDINATION.md spec).
- **Почему:** гейт — самый долгий job CI `engine`, около 7 минут, и почти всё — прогон учеников `sim` (measurements 3.1 spec). Строка 7.19 даёт `sim` рецепты `gate-part k N` и `gate-check`; без её тега делить нечего. Записана T65 по пункту 2 её «Сделать»: эта часть ждёт тега `sim` строки 7.19 (R103 spec; R17).
- **Сделать:**
  1. Тег `sim` строки 7.19: `gh release list -R altnovatrail/sim`, `gh release verify <тег> -R altnovatrail/sim`; что тег меняет в гейте — рецепты частей, `gate.json`, полосы, контейнер — запись [docs/sources.md](../docs/sources.md), раздел 6.
  2. `.github/workflows/ci.yml` — `SIM_TAG` и `SIM_COMMIT` на этот тег; «Gate» — матрица `just gate-part k N` (или как назовёт `sim`) и job проверки `gate-check` после неё, на `ubuntu-24.04-arm` (вторая поправка 2026-11-16 к R17), образ `toolchain` `sim` — как сейчас, из кеша. Обязательная проверка ruleset `main` — «Gate»: имя job после частей остаётся «Gate», чтобы ruleset не менять; если так нельзя — вопрос автору в [run/QUESTIONS.md](QUESTIONS.md), ruleset — настройка GitHub (правило 4).
  3. `justfile` — свой `gate` и, если нужно, `gate-part`: тег и SHA по-прежнему из `ci.yml`, локально — тот же гейт, что в CI.
  4. Поправка к R17; [docs/architecture.md](../docs/architecture.md), подраздел 3.7; [.claude/rules/ci.md](../.claude/rules/ci.md); README — если меняется вызов.
- **Не делать:** `.devcontainer/` — его меняет переход на тег `spec`, T68.
- **Как проверить:** правило [.claude/rules/ci.md](../.claude/rules/ci.md) и `just gate` в контейнере — зелёный, числа до бита те же, что у прежнего гейта, где тег их не меняет; CI pull request зелёный; время «Gate» до и после — по прогонам `main` и pull request, в отчёте и в [docs/measurements.md](../docs/measurements.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строки 7.19 и 7.20 COORDINATION.md spec; measurements 3.1 spec; R76 spec, R103 spec, R156 spec; R17, R28.

#### T67. Переход CI на релиз `.github` строки 7.13 — вне эпика

- **Цель:** общие проверки `.github` — на его релизе строки 7.13: кеш Go в `go.yml`, один файл образа `toolchain` без сохранения из pull request, «Go / targets» и «Go / fma-diff» — параллельно, сверка релиза — одним вызовом или параллельно (строка 7.20 COORDINATION.md spec).
- **Почему:** CI `engine` компилирует всё с пустым `GOCACHE` в каждом job, а образ из pull request засоряет кеш (measurements 3.1 spec). Записана T65 по пункту 2 её «Сделать»: эта часть ждёт релиза `.github` строки 7.13 (R76 spec; R28).
- **Сделать:**
  1. Релиз `.github` строки 7.13: `gh release list -R altnovatrail/.github`, `gh release verify <тег> -R altnovatrail/.github`; что меняют общие workflow от `v0.1.24` — входы, кеш Go, путь и тег файла образа, сохранение только из `main` — запись [docs/sources.md](../docs/sources.md), раздел 7.
  2. `.github/workflows/` — каждое `uses: altnovatrail/.github/…@<SHA> # v0.1.24` — на SHA этого релиза с его именем; у «Contracts» и «Gate» `ci.yml` — те же шаги образа и кеша, что у `go.yml` нового релиза, слово в слово (R28), и тот же кеш Go, если он берётся.
  3. Поправки к R28, R10 и R15 о версии `.github`; [docs/architecture.md](../docs/architecture.md), подраздел 3.7; [.claude/rules/ci.md](../.claude/rules/ci.md).
- **Решить при выполнении:** если workflow релиза требуют стадии `toolchain-core` в Dockerfile вызывающего, — сейчас только то, что работает без неё, а остальное — в T68.
- **Не делать:** `.devcontainer/` и тег `spec` — это T68.
- **Как проверить:** правило [.claude/rules/ci.md](../.claude/rules/ci.md) и `just same-bits`; CI pull request зелёный; время изменённых job до и после — по прогонам `main` и pull request, в отчёте и в [docs/measurements.md](../docs/measurements.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строки 7.13 и 7.20 COORDINATION.md spec; measurements 3.1 spec; R76 spec, R156 spec; R10, R15, R28.

#### T68. Переход на тег `spec` строки 7.14 со стадией `toolchain-core` — вне эпика

- **Цель:** `engine` — на теге `spec` строки 7.14: стадия `toolchain-core` в `.devcontainer/Dockerfile` по заготовке, и проверки без Go-тестов — в ней, как дают workflow `.github` строки 7.13 (строки 7.13, 7.14 и 7.20 COORDINATION.md spec).
- **Почему:** последняя из задач скорости CI, записанных T65: переход на тег `spec` меняет `.devcontainer/`, а пересобрать контейнер ночью некому — поэтому он идёт после T66 и T67 (T65, «Сделать», пункты 2 и 3). Переход на тег — задачей своей и решением журнала (RULES.md, Across repositories).
- **Сделать:**
  1. Тег `spec` строки 7.14 — или последний, если он новее: что он приносит `engine` — `git diff v0.18.3 <тег>` в `spec/`; `.devcontainer/` и `.claude/settings.json` — по `just -f spec/justfile new-repo-draft engine go 3.8 <тег>`, байт в байт ([.claude/rules/ci.md](../.claude/rules/ci.md)).
  2. Что из T67 ждало `toolchain-core`, — теперь.
  3. Решение журнала о переходе; поправки, которые оно требует; [docs/architecture.md](../docs/architecture.md), раздел 3; [.claude/rules/ci.md](../.claude/rules/ci.md).
- **Как проверить:** контейнер собирается из нового `.devcontainer/` — сборка стадий `toolchain`, `toolchain-core` и `devcontainer`; проверки частей `go` и `ci` по их правилам; CI pull request зелёный; время изменённых job до и после — в отчёте и в [docs/measurements.md](../docs/measurements.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строки 7.13, 7.14 и 7.20 COORDINATION.md spec; measurements 3.1 spec; R31 spec, R76 spec, R156 spec; R26, R28.
