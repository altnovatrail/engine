# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T86. Ссылки по R171 spec: архитектура, разделы 3 и 4 — вне эпика

- **Цель:** разделы 3 и 4 [docs/architecture.md](../docs/architecture.md) — по правилу ссылок R171 spec.
- **Почему:** седьмая часть T80 (архив [docs/plan/E5.md](../docs/plan/E5.md), T80); R171 spec.
- **Сделать:** как у T85, в своих разделах.
- **Как проверить:** как у T85, по своим разделам.
- **Разделы:** как у T85.
- **Зависит от:** T80.

#### T87. Ссылки по R171 spec: источники, замеры и ключ проверки — вне эпика

- **Цель:** [docs/sources.md](../docs/sources.md) и [docs/measurements.md](../docs/measurements.md) — по правилу ссылок R171 spec, а `docs-check-repo` держит правило ключом `"citations": true` в `docs-budgets.json`, как в `spec/repo-template/files/docs-budgets.json`.
- **Почему:** последняя часть T80 (архив [docs/plan/E5.md](../docs/plan/E5.md), T80): ключ включает она; R171 spec.
- **Сделать:**
  1. Записи источников и замеров держат дату и не держат номера задачи, в заголовке тоже; ссылки — по шагу 4 T80; ссылку на поправку — на запись.
  2. `docs-budgets.json` — `"citations": true`; что ключ ещё находит в документах и комментариях кода — по правилу.
- **Как проверить:** `just -f spec/justfile docs-check-repo` с ключом — зелёный; `/code-review high`, один круг; в отчёте — ссылок было до T80 и стало, по видам, командой `measurements 1.18 spec`.
- **Разделы:** R171 spec, R32 spec; `spec/repo-template/files/docs-budgets.json`.
- **Зависит от:** T81–T86.

Найдено в T84: в `justfile` сообщение об ошибке рецепта несёт номер решения — `(R16)` — в строке кода, не в комментарии, и ключ его не находит; по строгому правилу — причина словами.
