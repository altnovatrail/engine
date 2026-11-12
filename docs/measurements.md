# engine — measurements

> The results of the engine's measurements: what was measured, how, on what, and what numbers came out. The document that describes how it is built — [architecture.md](architecture.md) — refers here, and itself holds only what was decided from the measurements (RULES.md, How we work; R143 spec). A measurement is cited by section, like any document: `measurements 2.1`, and from another repository `measurements 2.1 engine`.
>
> An external fact a measurement rests on is a record of [sources.md](sources.md). Sections 1 and 2 were moved here word for word on 2026-11-12 from [architecture.md](architecture.md), subsections 2.9 and 2.10, by task T55 (R24); the analysis and the decisions stayed there. The measurement of the tokens spent by an epic's sessions, by the rules of epics of [RUN.md](../RUN.md), will go here as a section of its own when it is first made.

## 1. The corridor of the base

The run of T29, before the first candidate of E3: which branches and codes give items outside the corridor 0.70–0.85. The analysis by it — the hypotheses, the levers and the judgement — is architecture 2.9.

**How it was counted.** `sim` v0.1.2 (e61f90a) — the pin of T30 of `sim`: `modules` v0.1.4 with `olympiad-math` 0.1.1, `learner-model` v0.1.0. The populations `G0`, `G2-half` and `G2`, 1,000 working learners a cell (seed 20261107), the rule is the base. A copy of `engine` from `main` 167b96d (the code is `v0.1.3`) writes into each choice its branch, and for a task on a skill also its P by the model, by which the rule chose it, and the skill's candidates: in the corridor, above it, and in the corridor with a solution that has not gone to the device. A copy of `sim` sorts the items by these fields by the kind, the code and the branch of the choice. P* is `PGiven`, as in `E1`; the model's P is the item's, from the same `Readout` that the engine got, at the same moment of issue, and for the tasks of a lesson it is read anew before each task, at the same moment. The copies lay in a temporary folder; they do not go into the repository. The cross-check: in all six cells the sums by code, down to the item, match `E1` of plain `sim` v0.1.2 on `engine` v0.1.0 — so `v0.1.3` too chooses here the same as the base.

**The tables.** The shares are in percent: the choices and the items of the cell, P* and the model's P of the items of the row, the last five columns of the choices of a task on a skill. "In the corridor by the rule's P" is the P of the issued task on the skill by the model. "No unseen ones" is an issued task with a seen solution: the skill had no unseen candidates. Branch 5 is on the turn of repetition and outside it together, branch 9 all three of its steps. The tables hold the rows from 1% of the cell's items; the rest together are 0.6–2.0% on the olympiad maths and up to 0.7% on the school maths.

Olympiad maths:

| Cell | Kind, code, branch | Choices | Items | P*: in the corridor / below 0.50 / above 0.95 | Model's P: in the corridor / below 0.50 / above 0.95 | Mean P* / model's P | In the corridor by the rule's P | The skill has one in the corridor | Only above the corridor | Has an unseen one in the corridor | No unseen ones |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `G0` | task, `CONSOLIDATION`, 3 | 40.2 | 32.4 | 19.6 / 47.8 / 3.9 | 38.3 / 46.5 / 1.3 | 0.48 / 0.49 | 38.3 | 94.1 | 5.9 | 36.8 | 21.3 |
| `G0` | task, `QUEUE`, 7 | 25.6 | 20.6 | 36.6 / 22.6 / 1.3 | 49.2 / 19.9 / 0.3 | 0.65 / 0.66 | 49.2 | 97.3 | 1.4 | 46.5 | 13.7 |
| `G0` | lesson, `NEW_SKILL`, 6 | 7.6 | 18.4 | 13.6 / 37.3 / 19.0 | 20.6 / 37.5 / 16.6 | 0.61 / 0.62 | — | — | — | — | — |
| `G0` | task, `ROTATION`, 7 | 11.6 | 9.3 | 49.1 / 5.6 / 4.0 | 83.2 / 2.6 / 3.4 | 0.76 / 0.77 | 83.2 | 99.6 | 0.4 | 83.1 | 1.2 |
| `G0` | lesson, `QUEUE`, 7 | 2.7 | 6.5 | 30.5 / 17.3 / 19.7 | 31.9 / 18.6 / 19.1 | 0.74 / 0.74 | — | — | — | — | — |
| `G0` | task, `REPETITION`, 5 | 4.7 | 3.8 | 66.6 / 4.7 / 2.0 | 78.2 / 3.5 / 1.5 | 0.75 / 0.76 | 78.2 | 95.9 | 0.1 | 78.1 | 0.4 |
| `G0` | lesson, `ROTATION`, 7 | 1.5 | 3.7 | 21.8 / 40.8 / 2.8 | 28.1 / 44.8 / 1.8 | 0.56 / 0.54 | — | — | — | — | — |
| `G0` | task, `PROBE`, 4 | 3.9 | 3.1 | 7.0 / 57.3 / 0.0 | 0.0 / 58.2 / 0.0 | 0.42 / 0.41 | 0.0 | 0.0 | 0.0 | 0.0 | 11.5 |
| `G0` | task, `ROTATION`, 9 | 1.8 | 1.5 | 4.1 / 84.0 / 0.0 | 0.0 / 90.5 / 0.0 | 0.26 / 0.23 | 0.0 | 0.0 | 0.0 | 0.0 | 8.1 |
| `G2-half` | task, `CONSOLIDATION`, 3 | 35.6 | 28.3 | 20.2 / 41.3 / 6.3 | 40.1 / 43.7 / 1.6 | 0.54 / 0.51 | 40.1 | 94.9 | 5.1 | 38.6 | 19.1 |
| `G2-half` | task, `QUEUE`, 7 | 28.4 | 22.6 | 38.4 / 16.6 / 2.5 | 48.5 / 19.7 / 0.4 | 0.70 / 0.65 | 48.5 | 97.3 | 1.6 | 45.8 | 13.8 |
| `G2-half` | lesson, `NEW_SKILL`, 6 | 8.0 | 19.2 | 13.8 / 35.5 / 20.7 | 20.3 / 37.4 / 16.3 | 0.63 / 0.61 | — | — | — | — | — |
| `G2-half` | task, `ROTATION`, 7 | 13.5 | 10.8 | 45.3 / 3.1 / 5.7 | 81.6 / 2.7 / 3.5 | 0.81 / 0.77 | 81.6 | 99.1 | 0.9 | 81.4 | 2.0 |
| `G2-half` | lesson, `QUEUE`, 7 | 3.0 | 7.1 | 29.0 / 16.3 / 20.7 | 31.8 / 20.8 / 17.3 | 0.75 / 0.72 | — | — | — | — | — |
| `G2-half` | lesson, `ROTATION`, 7 | 1.7 | 4.1 | 22.6 / 35.9 / 3.7 | 28.8 / 47.6 / 1.2 | 0.60 / 0.52 | — | — | — | — | — |
| `G2-half` | task, `REPETITION`, 5 | 4.2 | 3.4 | 67.0 / 4.2 / 2.7 | 79.3 / 4.0 / 1.5 | 0.76 / 0.75 | 79.3 | 95.9 | 0.0 | 79.2 | 0.4 |
| `G2-half` | task, `PROBE`, 4 | 3.2 | 2.6 | 11.1 / 49.8 / 0.0 | 0.0 / 58.0 / 0.0 | 0.46 / 0.41 | 0.0 | 0.0 | 0.0 | 0.0 | 10.1 |
| `G2-half` | task, `ROTATION`, 9 | 1.4 | 1.1 | 5.3 / 78.3 / 0.0 | 0.0 / 89.3 / 0.0 | 0.31 / 0.25 | 0.0 | 0.0 | 0.0 | 0.0 | 6.1 |
| `G2` | task, `QUEUE`, 7 | 31.2 | 24.3 | 38.0 / 10.6 / 4.7 | 50.5 / 17.3 / 0.4 | 0.75 / 0.66 | 50.5 | 97.3 | 1.8 | 47.8 | 13.2 |
| `G2` | task, `CONSOLIDATION`, 3 | 30.8 | 24.0 | 20.1 / 34.5 / 10.3 | 44.1 / 39.1 / 1.9 | 0.60 / 0.53 | 44.1 | 94.8 | 5.2 | 42.9 | 17.4 |
| `G2` | lesson, `NEW_SKILL`, 6 | 8.5 | 19.9 | 13.9 / 33.6 / 22.6 | 20.1 / 38.6 / 15.5 | 0.64 / 0.61 | — | — | — | — | — |
| `G2` | task, `ROTATION`, 7 | 15.2 | 11.8 | 33.4 / 2.0 / 10.7 | 80.8 / 2.9 / 3.5 | 0.85 / 0.77 | 80.8 | 99.0 | 1.0 | 80.5 | 2.4 |
| `G2` | lesson, `QUEUE`, 7 | 3.4 | 7.9 | 26.9 / 16.8 / 21.9 | 30.8 / 24.7 / 15.5 | 0.75 / 0.70 | — | — | — | — | — |
| `G2` | lesson, `ROTATION`, 7 | 2.1 | 5.0 | 21.7 / 30.7 / 7.2 | 28.6 / 49.7 / 1.3 | 0.64 / 0.51 | — | — | — | — | — |
| `G2` | task, `REPETITION`, 5 | 3.9 | 3.1 | 67.8 / 3.7 / 2.9 | 82.1 / 4.2 / 1.2 | 0.78 / 0.75 | 82.1 | 96.5 | 0.0 | 81.9 | 0.7 |
| `G2` | task, `PROBE`, 4 | 2.8 | 2.1 | 17.7 / 35.4 / 0.1 | 0.0 / 49.9 / 0.0 | 0.53 / 0.45 | 0.0 | 0.0 | 0.0 | 0.0 | 8.7 |

The cell as a whole, P* in the corridor / below 0.50 / above 0.95: `G0` — 26.7 / 33.7 / 6.9%, `G2-half` — 27.7 / 27.7 / 8.7%, `G2` — 26.7 / 22.2 / 11.7%. Consolidation gives 46% of all the items below 0.50 on `G0`, 42% on `G2-half` and 37% on `G2`.

School maths — a skill in a choice has two candidates on average:

| Cell | Kind, code, branch | Choices | Items | P*: in the corridor / below 0.50 / above 0.95 | Model's P: in the corridor / below 0.50 / above 0.95 | Mean P* / model's P | In the corridor by the rule's P | The skill has one in the corridor | Only above the corridor | Has an unseen one in the corridor | No unseen ones |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `G0` | task, `QUEUE`, 7 | 31.2 | 26.8 | 22.9 / 52.3 / 5.7 | 27.3 / 52.8 / 5.8 | 0.49 / 0.50 | 22.4 | 36.3 | 62.0 | 0.0 | 100.0 |
| `G0` | task, `ROTATION`, 9 | 26.7 | 26.6 | 6.3 / 61.2 / 0.0 | 0.0 / 63.5 / 0.0 | 0.40 / 0.39 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G0` | task, `CONSOLIDATION`, 3 | 29.5 | 25.0 | 28.9 / 24.6 / 12.0 | 28.4 / 24.6 / 13.1 | 0.66 / 0.68 | 27.4 | 60.2 | 39.8 | 0.0 | 100.0 |
| `G0` | lesson, `NEW_SKILL`, 6 | 4.9 | 13.4 | 17.2 / 43.8 / 9.1 | 0.4 / 49.1 / 2.5 | 0.54 / 0.54 | — | — | — | — | — |
| `G0` | task, `PROBE`, 4 | 6.9 | 7.5 | 7.0 / 59.7 / 0.0 | 0.0 / 64.9 / 0.0 | 0.45 / 0.42 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G2-half` | task, `QUEUE`, 7 | 54.3 | 47.5 | 12.0 / 26.5 / 24.8 | 24.2 / 39.5 / 14.5 | 0.71 / 0.60 | 20.0 | 34.6 | 64.1 | 0.0 | 100.0 |
| `G2-half` | task, `ROTATION`, 9 | 17.1 | 17.5 | 38.8 / 23.2 / 0.0 | 0.0 / 50.6 / 0.0 | 0.62 / 0.46 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G2-half` | task, `CONSOLIDATION`, 3 | 18.1 | 15.6 | 9.7 / 8.4 / 46.3 | 25.4 / 8.5 / 19.1 | 0.86 / 0.82 | 25.1 | 53.9 | 46.1 | 0.0 | 100.0 |
| `G2-half` | lesson, `NEW_SKILL`, 6 | 5.1 | 13.7 | 17.2 / 40.4 / 9.9 | 0.7 / 47.2 / 3.4 | 0.57 / 0.56 | — | — | — | — | — |
| `G2-half` | task, `PROBE`, 4 | 4.3 | 4.6 | 47.5 / 6.8 / 0.0 | 0.0 / 41.0 / 0.0 | 0.75 / 0.51 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G2-half` | task, `ROTATION`, 8 | 1.2 | 1.0 | 10.8 / 0.3 / 45.4 | 9.8 / 0.3 / 33.8 | 0.91 / 0.88 | 11.0 | 31.1 | 68.9 | 0.0 | 100.0 |
| `G2` | task, `QUEUE`, 7 | 71.7 | 63.8 | 8.9 / 8.5 / 53.6 | 21.4 / 25.0 / 20.8 | 0.87 / 0.71 | 18.4 | 32.0 | 67.4 | 0.0 | 100.0 |
| `G2` | lesson, `NEW_SKILL`, 6 | 5.1 | 13.9 | 17.0 / 37.6 / 9.8 | 0.6 / 46.4 / 3.8 | 0.59 / 0.56 | — | — | — | — | — |
| `G2` | task, `ROTATION`, 9 | 11.4 | 11.7 | 35.8 / 11.4 / 2.9 | 0.0 / 45.5 / 0.0 | 0.73 / 0.49 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G2` | task, `CONSOLIDATION`, 3 | 8.4 | 7.3 | 4.7 / 4.3 / 60.9 | 27.3 / 4.3 / 19.1 | 0.91 / 0.85 | 27.1 | 52.5 | 47.5 | 0.0 | 100.0 |
| `G2` | task, `PROBE`, 4 | 1.8 | 2.0 | 13.2 / 0.0 / 18.6 | 0.0 / 27.9 / 0.0 | 0.90 / 0.55 | 0.0 | 0.0 | 0.0 | 0.0 | 100.0 |
| `G2` | task, `ROTATION`, 8 | 1.5 | 1.3 | 10.9 / 0.0 / 56.4 | 10.7 / 0.3 / 31.9 | 0.93 / 0.88 | 12.4 | 30.6 | 69.4 | 0.0 | 100.0 |

The cell as a whole: `G0` — 18.0 / 46.8 / 6.0%, `G2-half` — 18.7 / 23.8 / 20.8%, `G2` — 13.0 / 12.3 / 41.4%.

## 2. The speed of the choice

### 2.1. The profile of T37

Before the first fix of E4: where the time of `engine.Choose` goes on the pool of `olympiad-math` 0.1.1, 772 tasks, and whose it is — the engine's or `learner-model`'s. What came of it — the levers of T38 and the revision of the epic's goal — is architecture 2.10.

**How it was counted.** The benchmark `BenchmarkChoose` (`bench_internal_test.go`) — the package `olympiad-math` 0.1.1 from `testdata/modules` (the 2026-11-10 amendment to R10), a learner of grade 3, `under_13`, and four states of the model, built in the order of the graph, with the real `Read` of `learner-model` v0.1.0: "new" — no evidence and no issues; "in the middle" — the first ten skills of the graph practised, their lessons of the band `grades-3-4` passed, the solutions of half of their tasks of that band gone to the device, 60 choices of the module; "in the middle, consolidation" — the same and an error in the last attempt; "mastered" — 15 skills mastered, the rest practised, the lessons passed, the solutions of two tasks of three gone, the turn of repetition. `Read` is outside the measurement: `api` calls it, not the engine. Which choice each state gives is held by `TestBenchInputs`. The machine is an AMD Ryzen AI 9 365, linux/amd64, go1.27.3; the code of `engine` is `main` c6504bb, which T37 does not change. The numbers are by `go test -run '^$' -bench BenchmarkChoose -benchmem -count 6 -cpu 1`, the median of six; the profile is the same benchmark with `-benchtime 3s -cpuprofile -memprofile`, the four states together, 14.0 s in `Choose`; the shares are of the time of `Choose`, by `go tool pprof -top -cum`. The time depends on the load on the machine: the same run an hour later, under someone else's load on the host, gave 1.4–2.4 ms; the memory and the number of allocations were the same to the byte. So T38 compares the time before and after in one session, through `-count`, and the memory and the allocations exactly.

| State | Choice, code | Time per choice | Memory per choice | Allocations per choice |
|---|---|---|---|---|
| new | lesson `arithmetic-tricks-grades-1-2`, `NEW_SKILL` | 1.34 ms | 1.32 MB | 16,699 |
| in the middle | task `tri-34-d1-2`, `QUEUE` | 1.43 ms | 1.28 MB | 15,963 |
| in the middle, consolidation | task `tri-34-d1-2`, `CONSOLIDATION` | 1.34 ms | 1.26 MB | 15,191 |
| mastered | task `tri-12-d1-2`, `REPETITION` | 1.40 ms | 1.31 MB | 15,202 |

**The cross-check with `sim`.** The about 30 ms a call that the epic began with (Д9; Д167 spec; the row В10 of RUN.md `sim`) were not confirmed here — the choice is 20 times faster. A copy of `sim` v0.1.2 (e61f90a) in a temporary folder, with a `go.work` on the working tree of `engine`: 8 learners each of `G0` and `G2` on the olympiad maths 0.1.1, on the base. The inputs of all 2,628 choices, taken from the records of the steps, `Choose` goes through anew in 1.60–1.61 ms on average, p90 2.0–2.1 ms, the largest 2.9 ms; with `engine` v0.1.0, which `sim` takes without `go.work`, the same 1.58–1.63 ms. A learner runs 0.33–0.36 s, a step 2.0–2.1 ms, and `Choose` in the profile of the run is 4.40 s of the 5.51 s of the steps, 80%; the shares inside are the same as in the benchmark. Where the 30 ms come from cannot be found out here: so the synthetic states of the benchmark do represent the choice of `sim`, and the number in Д9 does not. The copies do not go into the repository.

**Where the time goes.** The shares are of `Choose`; the memory, of the bytes it allocates.

| Part | Functions | Time | Memory |
|---|---|---|---|
| The checks of the input | `checkIDs` … `checkAnchors`: most of all `checkTraps` (4.8%), `checkOwnModule` (2.5%), `checkIDs` (1.4%) | 9.7% | 10% |
| The P of tasks — `learner-model` | `probabilities` → `Readout.Task` on each task of the pool, 772 calls per choice | 32.0% | 32.8% |
| The P of tasks — the engine | `probabilities` outside `Readout.Task`: the sorting `tasksByID` and the check of P | 1.6% | 3.6% |
| Preparation | `ready` — the sets of the input and the open skills | 2.4% | 2.5% |
| The rule | `prepared.choose` — the branches and the skills | 54.4% | 49.4% |
| — of it, the view of the skills | `skills`, on each pass | 53.4% | 49.1% |
| — — the candidates | `candidates`: the whole pool for each of the 19 skills, `candidateOf` — 14,668 calls per pass (26.2%), the sorting `byCorridor` (6.8%), the pause (2.9%), appends to slices (4.8%) | 43.9% | 25.8% |
| — — the pool | `poolOf` — `poolTask` anew on each call: copies and sorts of the items and P | 6.9% | 21.3% |
| Garbage collection and allocation | `mallocgc` — through all the parts above, of it the assist to the collector 7.1% | 23.1% | — |

Inside `Readout.Task` — 0.56 µs per call, 0.43 ms per choice: the P of the items, `chanceOf`, is 45% of its time; parsing the task's ID, `moduleOf`, through `strings.SplitN` 16%; `checkTaskContent` and the copy with the sorting of the items 12% each. All of this is the code of `learner-model`, and the engine does not change it.

### 2.2. After T38 (2026-11-10)

The first four levers were taken, with the small things that go with them (R20): the index "skill → tasks" and the pool once per choice, `poolTask` without copies, candidates without sorting and without a slice of all the tasks, the tasks by ID once per choice, and the checks without copies. The same benchmark, the binaries before and after in turn in one session, the median of ten: 0.69–0.78 ms instead of 1.33–1.39, about 1.8×; memory per choice 0.67–0.73 MB instead of 1.26–1.32, allocations 9.5–11.0 thousand instead of 15.2–16.7. The engine's own share is about 0.31 ms instead of 0.93, three times less; now `Readout.Task` is 61% of `Choose`, building the pool and the index 15%, the rule 16%, of it the candidates 12.5%, the checks of the input about 2.5%. On the inputs of `sim`'s learners in the same session, 0.91–1.10 ms, and before the fix 1.76–1.87 ms; the time there is under someone else's load on the host and only for the cross-check. The choice is the same: the golden vectors, the hash of 7,677 choices of `sim`'s learners and all 69 numbers of `just gate`, to the bit. With the choice 1.8 times faster, a run of `sim`, where the choice is 80% of a step, is about 1.55× faster; what this gives the gate and the decisive run, `sim` will decide by a measurement of its own (the row В33).

### 2.3. After T48 (2026-11-11)

`learner-model` `v0.1.6` sped up `Readout.Task` (R29 learner-model), and `go.mod` requires it (R23). The same benchmark, the binaries before the raise and after in turn in one session, 10 measurements each, go1.27.4; the host was loaded, and the time is from the calmest of three sessions. The best measurement is 0.38–0.46 ms instead of 0.72–0.80, the medians 0.41–0.50 ms instead of 0.75–1.13, about 1.8×; memory per choice 0.27–0.33 MB instead of 0.67–0.73, allocations 1.75–2.56 thousand instead of 9.5–11.0. `Readout.Task` is now 36% of `Choose`, not 61%. The choice is the same: the golden vectors and `TestBenchInputs` did not change.
