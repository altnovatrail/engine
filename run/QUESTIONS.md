# QUESTIONS.md — questions to the author of `engine`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В50 ' run/QUESTIONS.md`.
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В51.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В50 | When to close Д10, the comment of `TaskSolution` in `types.go`. Its «When to return», the next change of the Go code, holds it back so that a comment alone does not release the same compiled code (R15); no change of the Go code is planned, and T68 releases a patch anyway: `.devcontainer/Dockerfile` and `.github/workflows/release.yml` are among the `paths` | The night run of 2026-11-16 (T171 spec, R156 spec): Д10 is the only row of table 4 closable inside the repository, and T68 edits `.devcontainer/`, which the run may not | (a) with T68: a step of its «Сделать», and Д10's «When to return» says so — no release of its own; (b) now, as a task outside an epic — one more patch with the same compiled code; (c) as it is, with the next change of the Go code | (a): the comment ships in a release that comes out anyway, and T68 is the next task |
