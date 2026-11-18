# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T77. Комментарии кода по R160 spec — вне эпика

- **Цель:** комментарии кода говорят только то, чего не говорит код (R160 spec, «Comments»).
- **Почему:** R160 spec; T73 разделила строку 7.30 COORDINATION.md spec по правилу 11 RUN.md и поправила номера решений, которые `docs-check` не находил в журнале.
- **Сделать:**
  - Комментарии Go, `.github/`, `justfile` и `.devcontainer/` (кроме того, что совпадает с заготовкой `spec` байт в байт, и копий шагов `.github` и `sim` в `.github/workflows/ci.yml` — их держат `.claude/rules/ci.md`, R28 и R31: там — только суффикс номера решения другого репозитория): почему, ограничение, ловушка, контракт экспортируемого API; без пересказа кода, истории, номеров задач и разделов документов; номер решения — якорем в конце, `(R25)`.
  - Код не меняется: только комментарии; `go vet` и линтер — зелёные.
- **Как проверить:** `just -f spec/justfile docs-check-repo`; проверки частей по `.claude/rules/go.md` и `.claude/rules/ci.md`; `/code-review high`.
- **Разделы:** R160 spec; RULES.md spec — The bar; R32 spec.
