# TASKS.md — tasks outside an epic of `engine`

Open tasks outside an epic, each as its body, in Russian; each also stands as a row of table 1 of [RUN.md](../RUN.md), with its status and what it depends on (R154 spec, R155 spec; R27).

- Never read whole: a session reads only its own task — `grep -n '^#### T' run/TASKS.md` lists them, `awk '/^#### T64\./{p=1;print;next} /^#### /{p=0} p' run/TASKS.md` gives one body.
- A new task outside an epic goes here as a body, its number the one free when it is written down (RUN.md, Citations), and as a row of table 1. When done, the body leaves for the archive (RUN.md rule 9).

## Tasks outside an epic

#### T72. Переход на тег `spec` v0.18.27 — вне эпика

- **Цель:** `engine` стоит на последнем теге `spec`, v0.18.27, — его правила, план и документы по этому тегу.
- **Почему:** R156 spec — параллельная работа волны 4 (слово автора 2026-11-18: пока `app` идёт эпиком E5, остальные репозитории не простаивают); переход — задача и решение журнала (RULES.md spec, Across repositories). Задача вне эпика или в текущем эпике — по правилам этого RUN.md.
- **Сделать:**
  1. `gh release verify v0.18.27 -R altnovatrail/spec`; подмодуль `spec/` — на v0.18.27; где репозиторий берёт модуль Go или пакет npm `spec` — та же версия, пакет npm — через `release-verify-npm` (R77 spec); тег — везде, где этот репозиторий его называет.
  2. Что теги после v0.18.14 до v0.18.27 несут этому репозиторию: `git -C spec log --oneline v0.18.14..v0.18.27`; строки таблицы 4 `spec` с именем репозитория в обратных кавычках — `grep -n '`engine`' spec/run/DEBT.md`; карточка этого репозитория в architecture 3 spec; RULES.md spec. Что из них берёт эта задача, а что — задача своя, решает сессия; отложенное — строкой таблицы 4.
  3. Заготовка `spec/repo-template/` (repositories 2.2 spec) на клоне `spec` v0.18.27 во временной папке: где она изменилась после v0.18.14 — README, `docs/README.md`, `.claude/`, CLAUDE.md, правила RUN.md —, взять, с отличиями, которые держат решения этого журнала.
  4. Общие workflow и действие `toolchain` `.github` — с v0.1.31 на релиз v0.1.33 (по SHA его тега, как сейчас): v0.1.33 восстанавливает кеши Go по префиксу без хеша `go.sum` (строка 2.7 COORDINATION.md spec, T75 `.github`), а v0.1.32 — составное действие `toolchain` для своих заданий (T69 `.github`). Проверки этого pull request — на новом SHA, зелёные.
  5. Решение журнала о переходе — что тег принёс и что взято —, поправки к записям, которые он меняет; строки таблиц 3 и 4, которые тег закрывает. Больше одного ревью — разбить по правилу 11 и назвать следующие задачи в отчёте.
- **Как проверить:** `just -f spec/justfile docs-check-repo` на новом теге — зелёный; проверки частей, которые задача меняет, — правило 6 этого RUN.md; `/code-review high`, один круг.
- **Разделы:** R156 spec, R76 spec; RULES.md spec; repositories 2.2 spec; карточка этого репозитория в architecture 3 spec.
