# QUESTIONS.md — questions to the author of `engine`

Table 3 of [RUN.md](../RUN.md): the author's open questions from tasks' «Открыто» (R153 spec, R154 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| В50 ' run/QUESTIONS.md`.
- A row comes from a task's «Открыто» (RUN.md rule 9). Record the author's answer where it belongs — the task, a document or [DEBT.md](DEBT.md) — and delete the row.
- A row whose question opens with the mark «В реестр spec» is for `spec`: its session finds it on `main` (R81 spec; RULES.md, Across repositories). It leaves once `spec` has recorded it: the task that rests on it names its number.

## Table 3. Questions to the author

Numbers are not reused: next — В52.

| ID | Question | From | Options | Recommendation |
|---|---|---|---|---|
| В51 | «В реестр spec»: in a copy, the template's `.devcontainer/Dockerfile` of `spec` v0.18.4 speaks of `spec`: the stage `toolchain` of a copy of the kind `go`, to which the kind adds no lines, reads "spec's tests need nothing the core lacks, so here the stage adds nothing", and the step of gopls "It is built with the Go of the toolchain stage", where Go now comes from `toolchain-core`. A session of the copy that weighs where a tool of its tests goes reads of `spec`'s tests, not its own; `engine` keeps the bytes, since a move to a tag compares them byte for byte | T68; R29 | (a) `new-repo` gives a copy its own sentence for a stage `toolchain` to which its kind adds nothing — "this repository's tests need nothing the core lacks" —, and the step of gopls names "the Go of toolchain-core, which toolchain takes"; (b) one sentence in the template that holds for `spec` and every copy alike; (c) leave it: R47 .github and each repository's rule of CI say it | (a): a copy's session reads the Dockerfile to put a tool in a stage, and the sentence should be about the copy; the cost is one line of `new-repo` and its test |
