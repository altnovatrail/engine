# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T68. Переход на тег `spec` строки 7.14 со стадией `toolchain-core` — вне эпика

- **Цель:** `engine` — на теге `spec` строки 7.14: стадия `toolchain-core` в `.devcontainer/Dockerfile` по заготовке, и проверки без Go-тестов — в ней, как дают workflow `.github` строки 7.13 (строки 7.13, 7.14 и 7.20 COORDINATION.md spec).
- **Почему:** последняя из задач скорости CI, записанных T65: переход на тег `spec` меняет `.devcontainer/`, а пересобрать контейнер ночью некому — поэтому он идёт после T66 и T67 (T65, «Сделать», пункты 2 и 3). Переход на тег — задачей своей и решением журнала (RULES.md, Across repositories).
- **Сделать:**
  1. Тег `spec` строки 7.14 — или последний, если он новее: что он приносит `engine` — `git diff v0.18.3 <тег>` в `spec/`; `.devcontainer/` и `.claude/settings.json` — по `just -f spec/justfile new-repo-draft engine go 3.8 <тег>`, байт в байт ([.claude/rules/ci.md](../.claude/rules/ci.md)).
  2. Что из T67 ждало `toolchain-core`, — теперь. T67 перевела CI на `.github` `v0.1.30`: в нём T60–T65 эпика E21 `.github`, но не T66 `.github` — стадия `toolchain-core` в `release.yml` и проверках; она выходит позже, после перехода `.github` на тег `spec` строки 7.14, потому что тот меняет `.devcontainer/` (решение координатора spec, T171 spec). Здесь — релиз `.github` с T66: `gh release list -R altnovatrail/.github`, `gh release verify`; что он даёт вызывающему и какие входы — по RUN.md, `docs/plan/` и журналу решений `.github` на его теге; SHA всех `uses:` `.github/workflows/` — на него, а шаги образа «Contracts» и job гейта — слово в слово с новым `go.yml`, как в T67 (третья поправка 2026-11-16 к R28).
  3. Решение журнала о переходе; поправки, которые оно требует; [docs/architecture.md](../docs/architecture.md), раздел 3; [.claude/rules/ci.md](../.claude/rules/ci.md).
- **Как проверить:** контейнер собирается из нового `.devcontainer/` — сборка стадий `toolchain`, `toolchain-core` и `devcontainer`; проверки частей `go` и `ci` по их правилам; CI pull request зелёный; время изменённых job до и после — в отчёте и в [docs/measurements.md](../docs/measurements.md); `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** строки 7.13, 7.14 и 7.20 COORDINATION.md spec; measurements 3.1 spec; R31 spec, R76 spec, R156 spec; R26, R28.
