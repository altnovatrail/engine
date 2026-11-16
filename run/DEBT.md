# DEBT.md — technical debt of `engine`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д5 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row comes from a task's «Открыто» (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed.

## Table 4. Technical debt

Numbers are not reused: next — Д11.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д5 | No framing of a task for the learner's interests: format v0 has no framing variants (R51 spec), and the rule reads no interests | T01; R03; R09 | Tasks are the same for every learner in one state | Format 0.4 (R111 spec) brought none. Once the format gives them, by a minor version, with a field of the input (R09) |
| Д6 | Traps the AI suggested are not read by the task choice: they do not go to the model (concepts 9.4 spec), and frequent traps are the program's findings only | T01; R03 | A mistake only the AI sees in the learner's work is not worked on by the choice | E3's judgement did not touch traps. Once the stand or the adult's progress show the program's traps are too few |
| Д7 | No tasks on paper, no glue with the device app, no offline | T09, T05; R04, R10 | None — stage 5; but glue added without an amendment to R10 fails the test of the rules | Stage 5, by О-11 and Т-10 |
| Д10 | The comment of `TaskSolution` in `types.go` speaks of any task, not of a task with a key: the solution of a task with a rubric, open after "Done", does not go to the summary of issues (the second 2026-11-15 amendment to R82 spec; the 2026-11-16 amendment to R03). The code does not change, and a fix of one comment in `paths` would release the same bits (R15) | T62; R26 | Whoever builds the input themselves, as `sim` does, may by the comment put a task with a rubric into `Summary.Tasks`, and the engine holds it back for a day | With T68, step 4 of its «Сделать»: its patch comes out anyway, as `.devcontainer/Dockerfile` and `.github/workflows/release.yml` are among the `paths` (R15) — the author's answer of 2026-11-16 to В50 |
