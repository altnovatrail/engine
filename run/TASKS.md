# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T64. Переход CI на релиз `.github` `v0.1.24` — вне эпика

- **Цель:** общие проверки `.github` — на его последнем релизе, не старше `v0.1.24`: с кешем образа `toolchain` (Д146 spec) и с тестами wasip1 в «Go / targets» и «Go / arm64» (Д213 spec), а в `.github/dependabot.yml` — группа `shared-workflows` (Д88 spec).
- **Почему:** третья часть перехода T62 по правилу 11 (R26), как T68 `sim`. `engine` вызывает общие workflow `.github` `v0.1.21`: каждое задание собирает образ из сети, и сбой nodejs.org, GitHub или снимка Ubuntu роняет проверки, а сборка для wasip1 не проверяется тестами (Д146, Д213 spec). Тега `spec` задача не ждёт.
- **Сделать:**
  1. Последний релиз `.github`: `gh release list -R altnovatrail/.github`, `gh release verify <тег> -R altnovatrail/.github`; что общие workflow меняют от `v0.1.21` — входы, кеш `toolchain`, wasip1 — запись [docs/sources.md](../docs/sources.md), раздел 7.
  2. `.github/workflows/` — каждое `uses: altnovatrail/.github/…@<SHA> # v0.1.21` — на SHA коммита последнего тега с его именем в комментарии; входы, которые новая версия требует или меняет; у задания «Gate» `ci.yml` — тот же кеш образа, что у `go.yml` нового тега.
  3. `.github/dependabot.yml` — группа `shared-workflows`, `patterns: ["altnovatrail/.github/*"]`, под `github-actions` (Д88 spec).
  4. Решение журнала, поправки к R10 и R15 о версии `.github`; где документы называют `.github` `v0.1.21` — [docs/architecture.md](../docs/architecture.md), раздел 3, README и [.claude/rules/ci.md](../.claude/rules/ci.md).
- **Не делать:** строку 7.20 `COORDINATION.md` `spec` — гейт частями и «Gate» на arm: она ждёт релиза `.github` строки 7.13 и тега `sim` строки 7.19.
- **Как проверить:** проверки CI pull request — зелёные, «Go / targets» — с wasip1, «Gate» — с кешем образа; правило [.claude/rules/ci.md](../.claude/rules/ci.md) и `just test-targets`; `just -f spec/justfile docs-check-repo`; `/code-review high`, один круг.
- **Разделы:** Д146 spec, Д213 spec, Д88 spec; R10, R15, R26; sources 7.2, 7.3.
