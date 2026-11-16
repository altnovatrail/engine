# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T63. План по R153–R155 spec: таблицы в `run/`, план на английском — вне эпика

- **Цель:** план `engine` — в раскладке и на языке заготовки `spec` v0.18.3: таблицы 2, 3 и 4 — `run/EPICS.md`, `run/QUESTIONS.md`, `run/DEBT.md`; правила RUN.md и таблицы — по-английски, коротко, для сессии, заголовки — «Table N. …»; тела задач и архивы эпиков — по-русски (R154 spec; Д250 spec).
- **Почему:** вторая часть перехода T62 по правилу 11 (R26): R153–R155 spec приходят с тегом v0.18.3, а `docs-check` и `just inbox` `spec` держат старую раскладку только до последнего переехавшего репозитория (Д250 spec). Так же разбил свой переход `sim`, T66 и T67 `sim`.
- **Сделать:**
  1. По заготовке `just new-repo-draft engine go 3.8 v0.18.3` в копии `spec` во временной папке: `run/EPICS.md` — таблица 2; `run/QUESTIONS.md` — таблица 3 со следующим свободным номером, В50; `run/DEBT.md` — таблица 4 вместо `DEBT.md` (`git mv`), строки Д5–Д7 — по-английски; `docs-budgets.json` — пределы новых файлов, если нужны.
  2. RUN.md — правила «Как пользоваться» и «Как устроен план» — по-английски, коротко, по заготовке, номера правил прежние, с отличиями, которые держат решения этого журнала: правило 9 и задача вне эпика между эпиками (R25; поправка 2026-11-13 к R19 spec) против `run/TASKS.md` R155 spec — взвесить и записать; фраза итога уходит в архив вместе с телом (R153 spec). Тела задач — по-русски; что из таблицы 1 ещё остаётся по-русски, — по заготовке и RULES.md, Hard rules — Language. История начал таблиц — закрытые строки таблиц 3 и 4, список сделанных эпиков — уходит из начал дословно, туда, где её найдут ссылки (R154 spec).
  3. CLAUDE.md — раздел Documents по заготовке: RUN.md и `run/`, и Environment — таблицы RUN.md и `run/`; README.md — план по-английски; ссылки на `DEBT.md` в `docs/` и `.claude/rules/` — на `run/DEBT.md`; архивы `docs/plan/` иначе не переписываются.
  4. Решение журнала — раскладка и язык плана, что взято из заготовки и что нет, поправки к записям, которые оно меняет (R19, R24, R25).
- **Мой шаг:** до задачи — пересобрать контейнер после мержа T62: без этого у сессий `engine` нет `gopls-lsp` — `.claude/settings.json` его больше не включает, а переменной `CLAUDE_CODE_PLUGIN_DIRS` у старого контейнера нет (R26).
- **Как проверить:** `just -f spec/justfile docs-check-repo` — зелёный, в раскладке `run/`; `/code-review high`, один круг.
- **Разделы:** R153 spec, R154 spec, R155 spec; Д250 spec; R24, R25, R26; RULES.md spec, Hard rules — Language.

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
