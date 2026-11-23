# DEBT.md — technical debt of `engine`

Table 4 of [RUN.md](../RUN.md): what was put off on purpose, its risk, when to return (R153 spec).

- Never read whole: take the rows the task and its sections name — `grep -n '^| Д5 ' run/DEBT.md` — and those `grep -n` finds by the paths and parts it changes (RUN.md rule 2).
- A row is what a task leaves open (RUN.md rule 9) and holds only what is open: drop a closed part, delete the row once closed.

## Table 4. Technical debt

Numbers are not reused: next — Д14.

| ID | Debt | From | Risk while open | When to return |
|---|---|---|---|---|
| Д5 | No framing of a task for the learner's interests: format v0 has no framing variants (R51 spec), and the rule reads no interests | T01; R03; R09 | Tasks are the same for every learner in one state | Format 0.4 (R111 spec) brought none. Once the format gives them, by a minor version, with a field of the input (R09) |
| Д6 | Traps the AI suggested are not read by the task choice: they do not go to the model (concepts 9.4 spec), and frequent traps are the program's findings only | T01; R03 | A mistake only the AI sees in the learner's work is not worked on by the choice | E3's judgement did not touch traps. Once the stand or the adult's progress show the program's traps are too few |
| Д7 | No glue with the device app, no offline | T09, T05; R04, R10 | None — stage 5; but glue added without an amendment to R10 fails the test of the rules | Stage 5, by Т-10 |
| Д13 | The jobs of the gate, "Gate / part k" and "Gate", keep a copy of the three image steps of the composite action `toolchain` of `.github`: the action hashes and builds `.devcontainer/Dockerfile` of the checkout's root alone, and the gate builds `sim`'s, `sim/.devcontainer/Dockerfile`, of a checkout in two trees (R31; sources 7.7). An input of the action for the Dockerfile's path would end the copy | T72; R31, R28 | A release of `.github` that changes the action's image steps reaches "Contracts" with Dependabot's SHA, and the gate only by hand, by the rule of `.claude/rules/ci.md` | A release of `.github` with such an input — a row «В реестр spec» if a second caller needs it —, or a move that changes the action's image steps |
