# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T93. `sources` по правилу записей, часть 3: тема `sim` — вне эпика

- **Цель:** тема `sim` `docs/sources.md` — его код, гейт, судьи и релизы — по правилу записей `spec/RULES.md`; после неё у `docs/sources.md` записи — только у того, что прогон не проверит или у чего нет версии.
- **Почему:** строка 7.66 COORDINATION.md `spec`; T91 разбита по правилу 11; после T93 строка закрыта.
- **Сделать:** записи 6.1–6.15 — файлы и теги `sim` — удалить, утверждения — ссылка на файл `sim` на теге или коммите, цитаты сверить на клоне во временной папке; что держит прогон гейта — `measurements`, где его там нет; повторные чтения — свернуть до текущего; номера не переиспользовать, `docs/plan/` не трогать; смысл утверждений не менять. Около 40 ссылок — журнал, `architecture`, `measurements`, `run/DEBT.md`.
- **Как проверить:** `just -f spec/justfile docs-check-repo` и проверки частей, чьи файлы задеты, — зелёные; размер `docs/sources.md` до и после — в отчёт; `/code-review high`, один круг.
- **Разделы:** `spec/RULES.md`, «The bar» и «Citations»; R171 `spec`; [docs/sources.md](../docs/sources.md), тема `sim`.
