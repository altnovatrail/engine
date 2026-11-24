# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T91. Переход на тег `spec` v0.19.20 и `sources` по правилу записей, часть 1 — вне эпика

- **Цель:** `docs/sources.md` по правилу записей `spec/RULES.md` нового тега: запись — только у того, что прогон не проверит или у чего нет версии; остальное — ссылка на месте с версией или коммитом, тест или `measurements`; теги `spec` и файлы репозиториев NovaTrail — git.
- **Почему:** строка 7.66 COORDINATION.md `spec`, по тексту T291 `spec`; переход — на последний тег, v0.19.20. Задача вне эпика — по правилам этого RUN.md.
- **Сделать:** переход на тег по `spec/tools/orchestrator/tasks/spec-tag.md` — CLAUDE.md, RUN.md и вводная `docs/sources.md` из заготовки шаблона; свернуть повторные чтения до текущего; записи о тегах `spec` и файлах других репозиториев — удалить, утверждения — ссылка на тег или репозиторий@коммит; поведение инструментов и сервисов под тестом, пробным PR или записью `measurements` — ссылка на месте с версией или коммитом, запись удалить; номера не переиспользовать, архивы `docs/plan/` не трогать; ссылки на записи `spec` (`… spec`) сверить с записью нового тега: удалённая или суженная, что утверждения уже не держит (так 5.4 — теперь только белая книга Edge), — так же. Смысл утверждений не менять. Больше одного ревью — по правилу 11 остальное в `run/TASKS.md`, строка ждёт всех.
- **Как проверить:** `just -f spec/justfile docs-check-repo` и проверки частей, чьи файлы задеты, — зелёные; размер `docs/sources.md` до и после — в отчёт.
- **Разделы:** `spec/RULES.md`, «The bar» и «Citations»; R171 `spec`; строка «Порядка».
- **Выяснилось при планировании (2026-11-24):** `docs/sources.md` — 82 805 байт, около 127 ссылок вне `docs/plan/`; почти все записи — файлы репозиториев NovaTrail и MathTrail и теги `spec`. Одно ревью этого не вмещает: по правилу 11 эта задача — переход на тег, вводные `docs/`, ссылки на записи `spec` и темы `.github` и «Go, GitHub and Codecov» с тегами `spec`; темы standalone, исследований, `modules`, `api` и `learner-model` — T92, тема `sim` — T93. Строка 7.66 ждёт всех трёх.

#### T92. `sources` по правилу записей, часть 2: темы standalone, исследований, `modules`, `api` и `learner-model` — вне эпика

- **Цель:** темы `docs/sources.md` standalone, исследований и других программ, `modules`, `api` и `learner-model` — по правилу записей `spec/RULES.md`, как T91 сделала с темами `.github` и Go.
- **Почему:** строка 7.66 COORDINATION.md `spec`; T91 разбита по правилу 11.
- **Сделать:** файлы MathTrail (standalone) и репозиториев NovaTrail (`modules`, `api`, `learner-model`), их релизы и `testdata/modules` — записи удалить, утверждения — ссылка на файл в репозитории@коммит или на тег, цитаты сверить на клоне во временной папке (MathTrail — тоже во временной папке, по RULES.md); исследования и документация других программ (2.1, 2.2) — статья вне репозиториев или страница без версии остаётся записью с одним текущим чтением, документация на версии или коммите — ссылка на месте; повторные чтения — свернуть до текущего; номера не переиспользовать, `docs/plan/` не трогать; смысл утверждений не менять. Около 60 ссылок — журнал, `architecture`, `measurements`.
- **Как проверить:** `just -f spec/justfile docs-check-repo` и проверки частей, чьи файлы задеты, — зелёные; размер `docs/sources.md` до и после — в отчёт; `/code-review high`, один круг.
- **Разделы:** `spec/RULES.md`, «The bar» и «Citations»; R171 `spec`; [docs/sources.md](../docs/sources.md), его темы от standalone до `learner-model`.

#### T93. `sources` по правилу записей, часть 3: тема `sim` — вне эпика

- **Цель:** тема `sim` `docs/sources.md` — его код, гейт, судьи и релизы — по правилу записей `spec/RULES.md`; после неё у `docs/sources.md` записи — только у того, что прогон не проверит или у чего нет версии.
- **Почему:** строка 7.66 COORDINATION.md `spec`; T91 разбита по правилу 11; после T93 строка закрыта.
- **Сделать:** записи 6.1–6.15 — файлы и теги `sim` — удалить, утверждения — ссылка на файл `sim` на теге или коммите, цитаты сверить на клоне во временной папке; что держит прогон гейта — `measurements`, где его там нет; повторные чтения — свернуть до текущего; номера не переиспользовать, `docs/plan/` не трогать; смысл утверждений не менять. Около 40 ссылок — журнал, `architecture`, `measurements`, `run/DEBT.md`.
- **Как проверить:** `just -f spec/justfile docs-check-repo` и проверки частей, чьи файлы задеты, — зелёные; размер `docs/sources.md` до и после — в отчёт; `/code-review high`, один круг.
- **Разделы:** `spec/RULES.md`, «The bar» и «Citations»; R171 `spec`; [docs/sources.md](../docs/sources.md), тема `sim`.
