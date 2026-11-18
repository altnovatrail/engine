# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T74. Без истории: CLAUDE.md, README, `docs/README.md`, раздел 1 архитектуры — вне эпика

- **Цель:** CLAUDE.md, README.md, `docs/README.md` и раздел 1 [docs/architecture.md](../docs/architecture.md) держит текущее состояние, а не историю (R160 spec, «One home», «Decisions»).
- **Почему:** R160 spec; T73 разделила по правилу 11 RUN.md строку 7.30 COORDINATION.md spec: переход на тег, `docs-budgets.json` и план — T73, остальное — T74–T77.
- **Сделать:**
  - Убрать хроники версий, статусов, повышений и итогов эпиков (пример — абзац «The component» CLAUDE.md о разделах 2.8–2.13, абзац «Status» `docs/architecture.md` `infra`): документ называет файл, где задано значение, — подмодуль, `go.mod`, тег релиза, — а не значение.
  - Раздел 1 архитектуры — текущее устройство; решение — номером, без пересказа.
  - Ссылка на поправку записи `spec`, которую тег переписал на месте (R143, R145, R154, R155 spec и другие, что E16 `spec` перепишет), — на запись.
  - Записи журнала, которые задача меняет, — переписать на месте целиком, с поправками внутри (R160 spec); других записей не трогать.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec; RULES.md spec — The bar; R32 spec.

#### T75. Без истории: раздел 2 архитектуры — вне эпика

- **Цель:** раздел 2 [docs/architecture.md](../docs/architecture.md) держит текущее состояние, а не историю (R160 spec, «One home», «Decisions»).
- **Почему:** R160 spec; T73 разделила по правилу 11 RUN.md строку 7.30 COORDINATION.md spec: переход на тег, `docs-budgets.json` и план — T73, остальное — T74–T77. После предыдущей задачи — тот же файл `docs/architecture.md`, и её правки ссылок на записи этот раздел видит.
- **Сделать:**
  - Подразделы 2.1–2.13 — что решено и где это в коде и журнале, номером решения, без пересказа записей и хроник кругов, судов и эпиков; что история — в архивах `docs/plan/` и git.
  - Ссылка на поправку записи `spec`, которую тег переписал на месте (R143, R145, R154, R155 spec и другие, что E16 `spec` перепишет), — на запись.
  - Записи журнала, которые задача меняет, — переписать на месте целиком, с поправками внутри (R160 spec); других записей не трогать.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec; RULES.md spec — The bar; R32 spec.

#### T76. Без истории: раздел 3 архитектуры — вне эпика

- **Цель:** раздел 3 [docs/architecture.md](../docs/architecture.md) держит текущее состояние, а не историю (R160 spec, «One home», «Decisions»).
- **Почему:** R160 spec; T73 разделила по правилу 11 RUN.md строку 7.30 COORDINATION.md spec: переход на тег, `docs-budgets.json` и план — T73, остальное — T74–T77. После предыдущей задачи — тот же файл `docs/architecture.md`, и её правки ссылок на записи этот раздел видит.
- **Сделать:**
  - Раздел 3 — библиотека, как она в коде: без версий, задач и поправок; значение, заданное в файле (`go.mod`, `justfile`, `.github/workflows/`), — ссылкой на файл.
  - Ссылка на поправку записи `spec`, которую тег переписал на месте (R143, R145, R154, R155 spec и другие, что E16 `spec` перепишет), — на запись.
  - Записи журнала, которые задача меняет, — переписать на месте целиком, с поправками внутри (R160 spec); других записей не трогать.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; `/code-review high`.
- **Разделы:** R160 spec; RULES.md spec — The bar; R32 spec.

#### T77. Комментарии кода по R160 spec — вне эпика

- **Цель:** комментарии кода говорят только то, чего не говорит код (R160 spec, «Comments»).
- **Почему:** R160 spec; T73 разделила строку 7.30 COORDINATION.md spec по правилу 11 RUN.md и поправила номера решений, которые `docs-check` не находил в журнале.
- **Сделать:**
  - Комментарии Go, `.github/`, `justfile` и `.devcontainer/` (кроме того, что совпадает с заготовкой `spec` байт в байт, и копий шагов `.github` и `sim` в `.github/workflows/ci.yml` — их держат `.claude/rules/ci.md`, R28 и R31: там — только суффикс номера решения другого репозитория): почему, ограничение, ловушка, контракт экспортируемого API; без пересказа кода, истории, номеров задач и разделов документов; номер решения — якорем в конце, `(R25)`.
  - Код не меняется: только комментарии; `go vet` и линтер — зелёные.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; проверки частей по `.claude/rules/go.md` и `.claude/rules/ci.md`; `/code-review high`.
- **Разделы:** R160 spec; RULES.md spec — The bar; R32 spec.
