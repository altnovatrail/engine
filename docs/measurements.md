# engine — measurements

> The results of the engine's measurements: what was measured, how, on what, and what numbers came out. The document that describes how it is built — [architecture.md](architecture.md) — refers here, and itself holds only what was decided from the measurements. A claim cites a measurement by its number, `measurements 2.1`, and from another repository `measurements 2.1 engine`.
>
> A fact no run checks that a measurement rests on is a record of [sources.md](sources.md); a program's documentation or code at a version or commit is a link beside it. A section's heading names what was measured, not a task. The measurement of the tokens spent by an epic's sessions, by the plan's rules for epics in [RUN.md](../RUN.md), is measurements 4.

## 1. The corridor of the base

The run before the first candidate of E3: which branches and codes give items outside the corridor 0.70–0.85. The analysis by it — the hypotheses, the levers and the judgement — is in [architecture.md](architecture.md).

**How it was counted.** `sim` v0.1.2 (e61f90a) — the pin of `sim` then: `modules` v0.1.4 with `olympiad-math` 0.1.1, `learner-model` v0.1.0. The populations `G0`, `G2-half` and `G2`, 1,000 working learners a cell (seed 20261107), the rule is the base. A copy of `engine` from `main` 167b96d (the code is `v0.1.3`) writes into each choice its branch, and for a task on a skill also its P by the model, by which the rule chose it, and the skill's candidates: in the corridor, above it, and in the corridor with a solution that has not gone to the device. A copy of `sim` sorts the items by these fields by the kind, the code and the branch of the choice. P* is `PGiven`, as in `E1`; the model's P is the item's, from the same `Readout` that the engine got, at the same moment of issue, and for the tasks of a lesson it is read anew before each task, at the same moment. The copies lay in a temporary folder; they do not go into the repository. The cross-check: in all six cells the sums by code, down to the item, match `E1` of plain `sim` v0.1.2 on `engine` v0.1.0 — so `v0.1.3` too chooses here the same as the base.

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

### 2.1. The profile of the choice before its speed-up

Before the first fix of E4: where the time of `engine.Choose` goes on the pool of `olympiad-math` 0.1.1, 772 tasks, and whose it is — the engine's or `learner-model`'s. What came of it — the levers of the speed-up — is in [architecture.md](architecture.md).

**How it was counted.** The benchmark `BenchmarkChoose` (`bench_internal_test.go`) — the package `olympiad-math` 0.1.1 from `testdata/modules`, a learner of grade 3, `under_13`, and four states of the model, built in the order of the graph, with the real `Read` of `learner-model` v0.1.0: "new" — no evidence and no issues; "in the middle" — the first ten skills of the graph practised, their lessons of the band `grades-3-4` passed, the solutions of half of their tasks of that band gone to the device, 60 choices of the module; "in the middle, consolidation" — the same and an error in the last attempt; "mastered" — 15 skills mastered, the rest practised, the lessons passed, the solutions of two tasks of three gone, the turn of repetition. `Read` is outside the measurement: `api` calls it, not the engine. Which choice each state gives is held by `TestBenchInputs`. The machine is an AMD Ryzen AI 9 365, linux/amd64, go1.27.3; the code of `engine` is `main` c6504bb, which the measurement does not change. The numbers are by `go test -run '^$' -bench BenchmarkChoose -benchmem -count 6 -cpu 1`, the median of six; the profile is the same benchmark with `-benchtime 3s -cpuprofile -memprofile`, the four states together, 14.0 s in `Choose`; the shares are of the time of `Choose`, by `go tool pprof -top -cum`. The time depends on the load on the machine: the same run an hour later, under someone else's load on the host, gave 1.4–2.4 ms; the memory and the number of allocations were the same to the byte. So the speed-up compares the time before and after in one session, through `-count`, and the memory and the allocations exactly.

| State | Choice, code | Time per choice | Memory per choice | Allocations per choice |
|---|---|---|---|---|
| new | lesson `arithmetic-tricks-grades-1-2`, `NEW_SKILL` | 1.34 ms | 1.32 MB | 16,699 |
| in the middle | task `tri-34-d1-2`, `QUEUE` | 1.43 ms | 1.28 MB | 15,963 |
| in the middle, consolidation | task `tri-34-d1-2`, `CONSOLIDATION` | 1.34 ms | 1.26 MB | 15,191 |
| mastered | task `tri-12-d1-2`, `REPETITION` | 1.40 ms | 1.31 MB | 15,202 |

**The cross-check with `sim`.** The about 30 ms a call that the epic began with, by `sim`'s profile (row В10 of its [RUN.md](https://github.com/altnovatrail/sim/blob/0364b40/RUN.md) at 0364b40), were not confirmed here — the choice is 20 times faster. A copy of `sim` v0.1.2 (e61f90a) in a temporary folder, with a `go.work` on the working tree of `engine`: 8 learners each of `G0` and `G2` on the olympiad maths 0.1.1, on the base. The inputs of all 2,628 choices, taken from the records of the steps, `Choose` goes through anew in 1.60–1.61 ms on average, p90 2.0–2.1 ms, the largest 2.9 ms; with `engine` v0.1.0, which `sim` takes without `go.work`, the same 1.58–1.63 ms. A learner runs 0.33–0.36 s, a step 2.0–2.1 ms, and `Choose` in the profile of the run is 4.40 s of the 5.51 s of the steps, 80%; the shares inside are the same as in the benchmark. Where the 30 ms come from cannot be found out here: so the synthetic states of the benchmark do represent the choice of `sim`, and the 30 ms do not. The copies do not go into the repository.

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

### 2.2. After the speed-up (2026-11-10)

The first four levers were taken, with the small things that go with them: the index "skill → tasks" and the pool once per choice, `poolTask` without copies, candidates without sorting and without a slice of all the tasks, the tasks by ID once per choice, and the checks without copies. The same benchmark, the binaries before and after in turn in one session, the median of ten: 0.69–0.78 ms instead of 1.33–1.39, about 1.8×; memory per choice 0.67–0.73 MB instead of 1.26–1.32, allocations 9.5–11.0 thousand instead of 15.2–16.7. The engine's own share is about 0.31 ms instead of 0.93, three times less; now `Readout.Task` is 61% of `Choose`, building the pool and the index 15%, the rule 16%, of it the candidates 12.5%, the checks of the input about 2.5%. On the inputs of `sim`'s learners in the same session, 0.91–1.10 ms, and before the fix 1.76–1.87 ms; the time there is under someone else's load on the host and only for the cross-check. The choice is the same: the golden vectors, the hash of 7,677 choices of `sim`'s learners and all 69 numbers of `just gate`, to the bit. With the choice 1.8 times faster, a run of `sim`, where the choice is 80% of a step, is about 1.55× faster; what this gives the gate and the decisive run, `sim` will decide by a measurement of its own.

### 2.3. After `learner-model` v0.1.6 (2026-11-11)

`learner-model` `v0.1.6` sped up `Readout.Task` ([`readout.go` of `learner-model` v0.1.6](https://github.com/altnovatrail/learner-model/blob/v0.1.6/readout.go)), and `go.mod` requires it. The same benchmark, the binaries before the raise and after in turn in one session, 10 measurements each, go1.27.4; the host was loaded, and the time is from the calmest of three sessions. The best measurement is 0.38–0.46 ms instead of 0.72–0.80, the medians 0.41–0.50 ms instead of 0.75–1.13, about 1.8×; memory per choice 0.27–0.33 MB instead of 0.67–0.73, allocations 1.75–2.56 thousand instead of 9.5–11.0. `Readout.Task` is now 36% of `Choose`, not 61%. The choice is the same: the golden vectors and `TestBenchInputs` did not change.

## 3. The rough look of the variant of E7

The variant of consolidation in a lesson against the base on `sim` before its judgement, as the variants of E3 were looked at: a rough look, which decides nothing. The rule of the variant is in [architecture.md](architecture.md); the step and the criterion of `sim` are the entries on the step within a lesson and on the judgement of its [log](https://github.com/altnovatrail/sim/blob/v0.1.15/docs/decisions.md) at `v0.1.15`.

**How it was counted.** 2026-11-14, in the container of `engine`. A clone of `sim` `v0.1.15` (52a3658) in a temporary folder, and its recipe of a rough look: `just rough engine runs/base-300.jsonl` — the base, `engine` `v0.1.6` by `go.mod` of `sim`, — and `CANDIDATE_ENGINE=<folder> just rough engine runs/e7-300.jsonl`, where the folder is a copy of the working tree of `engine` with the variant and a local commit, 492f939, also in the temporary folder: the recipe takes a full SHA on GitHub or a clean folder, and nothing was committed to the repository. The build of the candidate went through the workspace of the gate, with `learner-model` at ab688b2, and held to `gate.json`. 300 working learners a cell, the whole matrix, each run about 6 min 15 s by the clock and 107–109 CPU minutes. SHA-256 of the runs: the base `bcacb08e…125d`, the candidate `d85202c9…b5e`. Then `sim judge` with the judgement file of one candidate: `{"name": "E7", "group": "corridor", "constants": 0, "new_field": false, "distance": 1, "equal_on": ["novatrail/grade-3-math", "novatrail/english-essays"]}`, as `sim`'s [log](https://github.com/altnovatrail/sim/blob/v0.1.15/docs/decisions.md) at `v0.1.15` asks of it. The look ran on the code before the review of the variant. Its fixes touch the map, which `sim` does not have, and read the skills of a run from `Run.Skills`; `just gate` after them gives the same 69 numbers to the bit.

**The outcome.** The report opens "A rough look, which decides nothing": E7 does not hold every constraint and is not better than the base. Its group, the corridor, is +0.00002, with the interval [−0.0019, +0.0020]. It fails "no worse" on three cells of olympiad maths: `E2.stuck` on `G2-half` — worse by at most 0.0333 with the tolerance 0.0316 — and `E4.given` on `G5` — 0.0119 with 0.0103 — and on `S-scale3.5` — 0.0100 with 0.0100. The constraint `neutral` holds: on `grade-3-math` and `english-essays` every learner's sums are the base's to the bit, and the report names no "fails neutral". `E3.end` on `F-model`, `F-exp` and `F-growth5` of `grade-3-math` left the score within the author's tolerance of its path.

**The note `E6`, on the 34 cells of olympiad maths and of the pair with it** — the rows of `G0`, `G2-half` and `G2`, the three cells of the corridor, and the range over all the cells:

| Number | `G0` base → E7 | `G2-half` base → E7 | `G2` base → E7 | Over the 34 cells: E7 − base, or E7 where the base has none |
|---|---|---|---|---|
| `E6.consolidations`, consolidations in a lesson a pass | 0 → 0.968 | 0 → 0.909 | 0 → 0.847 | 0.744 to 1.007 |
| `E6.lesson.corridor`, the items of the lessons' task segments in 0.70–0.85 | 0.183 → 0.172 | 0.182 → 0.170 | 0.179 → 0.170 | −0.018 to +0.003 |
| `E6.lesson.hard`, the same below 0.50 | 0.341 → 0.358 | 0.312 → 0.325 | 0.286 → 0.298 | +0.004 to +0.019 |
| `E6.consolidation.corridor`, the items of the consolidations in 0.70–0.85 | — → 0.375 | — → 0.363 | — → 0.299 | 0.277 to 0.429 |
| `E6.consolidation.hard`, the same below 0.50 | — → 0.103 | — → 0.069 | — → 0.067 | 0.046 to 0.167 |
| `E6.next`, the mean P* of the next segment after a failed one | 0.286 → 0.266 | 0.301 → 0.278 | 0.318 → 0.316 | −0.028 to −0.001 |

A consolidation's items are in the corridor about twice as often as a segment's, and below 0.50 three to five times less often, yet the corridor of the cell does not move: a consolidation is about one task a pass, and on the segments of the lessons the share in the corridor falls and that below 0.50 grows. The interval of `E6.lesson.hard` is above zero on all 34 cells, that of `E6.lesson.corridor` below zero on 31, and that of `E6.next` below zero on 28: the next segment after a failure is given harder than in the base.

**"Gate" on the variant.** `just gate` on the working tree with the variant, on `sim` `v0.1.15`, 250 learners a cell: red by two bands of zero width alone — `E2.stuck` on `G2-half` of olympiad maths, 0.016, and on `A-linked`, 0.008, with the band 0 ± 0, — the other 67 numbers within their bands, of them the numbers of `grade-3-math` alone and of `english-essays` to the bit; nobody broke an invariant, and no run ended with an error. The same gate on the code of `main` e0fdb3f was green, all 69 numbers to the bit.

## 4. The tokens of an epic's sessions

`just -f spec/justfile tokens-report --since … --row` prints one row of the table below, for the window of an epic: from the merge of the previous epic's branch, on the epic's first day, to the end of the session of its outcome. The row "before" is a baseline, not a part of a sum: the window of E7 overlaps it. It runs at the archive of each epic, by the plan's rules, and once a month from `spec` (measurements 1.9 spec); the transcripts are the folder `-workspaces-engine` of this machine, and the cost is at the API's prices.

| Project | Window, UTC | Sessions | Calls, main / subagents | Context of a call, k, main / subagents | Compactions | Cost at API prices |
|---|---|---|---|---|---|---|
| `-workspaces-engine` — before, the row of measurements 1.9 spec | 2026-10-13T00:00 … 2026-11-12T14:27 | 16 | 3,016 / 1,232 | 328.6 / 79.3 | 3 | $400.70 |
| `-workspaces-engine` — epic E7 and the tasks outside it in its window | 2026-11-11T22:51 … 2026-11-14T18:39 | 13 | 1,072 / 681 | 169.2 / 103.1 | 1 | $109.14 |
| `-workspaces-engine` — epic E5, but its dropped task | 2026-11-14T18:48 … 2026-11-15T02:03 | 6 | 514 / 140 | 195.7 / 71.6 | 1 | $42.86 |

## 5. The corridor on the new base

The run before the candidates of E5: the rough look of the base and of the variants A and C of the round `engine-1` on the pin of `sim` after `model-2`, the base on the model before `theta`, and the items, the lag of the model and the coverage by the kind of the step, the reason code of its choice and whether the task's solution was seen. A rough look decides nothing, by `sim`'s criterion. The analysis by it — the hypotheses and the levers — is in [architecture.md](architecture.md); the facts of `sim` and `learner-model` it rests on are [`sim` v0.1.17](https://github.com/altnovatrail/sim/tree/v0.1.17) with the [report](https://github.com/altnovatrail/sim/blob/v0.1.17/judgements/engine-1/report.md) of the judgement `engine-1` of `sim` and the entry on it in its [log](https://github.com/altnovatrail/sim/blob/v0.1.17/docs/decisions.md) at `v0.1.17`, and [`learner-model` v0.1.8](https://github.com/altnovatrail/learner-model/tree/v0.1.8).

**How it was counted.** 2026-11-14, in the container of `engine`: AMD Ryzen AI 9 365, 20 processors, go1.27.4. A clone of `sim` `v0.1.17` (6f5151f) in a temporary folder, and its recipe of a rough look, 300 working learners a cell, the whole matrix:
- `just rough engine runs/base-300.jsonl` — the base, `engine` `v0.1.6` and `learner-model` `v0.1.8` by `go.mod` of `sim`;
- `CANDIDATE_ENGINE=7ef8288c677a8879bd50b89e6e6320584807d7d5 just rough engine runs/a-300.jsonl`, and the same with 97abff0d1078025a646ff574c73ba10a3298337a into `c-300.jsonl` — the heads of #46 and #48, the code `engine-1` judged. It is the code of `v0.1.4` with the fix of `I4` that the base got too, and the speed-up of the choice does not change the choice, so a port onto the base's code would choose the same. Both built in the workspace of the gate with `learner-model` at ab688b2 and held to `gate.json`;
- `CANDIDATE_LEARNER_MODEL=7d796e290599eeb15481f0098c9356bf6c31afe9 just rough engine runs/pretheta-300.jsonl` — the base on `learner-model` `v0.1.6`, the base of `sim` before `theta`.

By the clock 7 min 44 s, 19 min 58 s, 15 min 17 s and 12 min 05 s: the heads lack the speed-up of the choice, and the last ran under the load of other work. SHA-256: the base `bcacb08e…125d` — to the byte the base of the rough look of E7 on `v0.1.15` (measurements 3) —, A `84c8ddd4…0a6d`, C `de37ae53…f801`, before `theta` `4f58129c…11ad`. Then `sim compare` of the base with each, and `sim judge` with the judgement file of `engine-1` for A and C: `{"name": "A", "group": "corridor", "constants": 0, "new_field": false, "distance": 1}` and the same for C with `distance` 3.

The split is a second clone of the same commit with a file of its own, `internal/measure/split.go`, called beside the note `E6` in `internal/bench/bench.go` under the engine alone, and the variable `POPS` for the populations in the recipe `cells`. For each attempt with a verdict it adds, by the group "the kind of the step — a task, a lesson, a lesson resumed after a consolidation —, the reason code of its choice, and `Obs.Repeat`", the items and their shares of P* at issue as `E1` reads them; the sums of P* and of the model's P; the lag of `m3` — θ + δ − θ* − δ* of the skills with evidence after the attempt, from the 101st answer; and, in the same window, the moves of θ + δ and of θ* + δ* through the attempt of the skills with evidence before and after it. For `E4` it adds each skill within reach and the code that first gave it. Its keys start with `X.`, beside those of `sim`. The populations `G2`, `G2-half`, `G2-fading`, `G3`, `F-growth5`, `F-none` and `S-scale1.5` — where `engine-1`, E7 or the rough look below failed `M3.lag`, `M3.settle` or `E4.given` —, 1,000 working learners a cell, for the same four rules: `POPS=… just rough engine runs/<name>-1000.jsonl 1000`. SHA-256: the base `eed984ce…785c`, A `073867f8…133f`, C `8cf4b1cf…a11a`, before `theta` `71dc28eb…4097`; of the file `split.go` `b7cc462e…65b8`. The cross-check: the sums of `sim`'s own of the first 300 learners of each cell, 3,900 learners a rule, are the same to the bit as in the rough look of the same rule, and in each of the 13 cells the groups sum to `E1.items`, `E1.corridor` and the sums of `M3.lag`. A small program in the temporary folder reads the files. The clones, the file and the program do not go into the repository.

**The rough look.** The cells of the populations where `engine-1`, E7 or this look failed `E1.corridor`, `M3.lag`, `M3.settle` or `E4.given`, and `G0` beside them, by `sim compare`, 300 learners a cell; the failures on `W-repeat` and `F-model` are in the list of `sim judge` below. `E1.corridor` here is the share of items in the corridor (more is better), `M3.lag` the signed lag (the criterion reads its absolute value), `E4.given` the share of the skills within reach given at least once. Each cell: before `theta` → the base / A / C.

| Cell | `E1.corridor` | `M3.lag` | `E4.given` |
|---|---|---|---|
| olympiad `G0` | 0.259 → 0.259 / 0.381 / 0.389 | −0.028 → −0.028 / −0.020 / −0.010 | 0.868 → 0.868 / 0.868 / 0.868 |
| olympiad `G2` | 0.254 → 0.254 / 0.280 / 0.286 | −1.026 → −1.026 / −1.030 / −1.021 | 0.798 → 0.798 / 0.796 / 0.797 |
| olympiad `G2-half` | 0.265 → 0.265 / 0.339 / 0.344 | −0.529 → −0.529 / −0.516 / −0.524 | 0.827 → 0.827 / 0.826 / 0.830 |
| olympiad `G2-fading` | 0.267 → 0.267 / 0.354 / 0.356 | −0.431 → −0.431 / −0.418 / −0.418 | 0.833 → 0.833 / 0.834 / 0.834 |
| olympiad `G3` | 0.256 → 0.256 / 0.325 / 0.329 | −0.616 → −0.616 / −0.617 / −0.619 | 0.831 → 0.831 / 0.832 / 0.831 |
| olympiad `S-scale1.5` | 0.267 → 0.267 / 0.365 / 0.374 | 0.929 → 0.929 / 0.843 / 0.813 | 0.889 → 0.889 / 0.882 / 0.883 |
| olympiad `F-none` | 0.136 → 0.137 / 0.368 / 0.387 | 0.050 → 0.053 / 0.022 / 0.017 | 0.899 → 0.897 / 0.890 / 0.894 |
| olympiad `F-growth5` | 0.167 → 0.167 / 0.424 / 0.437 | −0.152 → −0.165 / −0.189 / −0.195 | 0.879 → 0.879 / 0.874 / 0.874 |
| school `G0` | 0.190 → 0.190 / 0.261 / 0.261 | 0.112 → 0.112 / 0.111 / 0.112 | 0.659 → 0.659 / 0.659 / 0.659 |
| school `G2` | 0.131 → 0.131 / 0.124 / 0.124 | −2.448 → −2.448 / −2.476 / −2.474 | 0.291 → 0.291 / 0.291 / 0.291 |
| school `G2-half` | 0.193 → 0.193 / 0.183 / 0.183 | −1.067 → −1.067 / −1.064 / −1.063 | 0.327 → 0.327 / 0.325 / 0.325 |
| school `G3` | 0.178 → 0.178 / 0.206 / 0.206 | −0.495 → −0.495 / −0.475 / −0.475 | 0.424 → 0.424 / 0.424 / 0.424 |

`theta` changes the base on no cell of `G0`, `G2`, `G2-half`, `G2-fading`, `G3` and `S-scale1.5`: `sim compare` gives zero with the interval [0, 0] on every number there. It moves only the populations with breaks, `F-`: on `F-growth5` of olympiad maths `M3.lag` from −0.152 to −0.165.

`sim judge`: "A rough look, which decides nothing: no candidate holds every constraint and is better than the base". The group `corridor` is 0.146 [0.125; 0.167] for A and 0.149 [0.129; 0.170] for C; the score is not read on a rough look. The checks "no worse" whose worst end is past the tolerance:
- both: school `E1.corridor` on `G2` (A 0.0106 at 0.0100, C 0.0109) and on `G2-half` (0.0148 at 0.0123), school `M3.lag` on `G2` (A 0.0436 at 0.0334, C 0.0416), school `M3.settle` on `G3` (10.67 at 8.18), olympiad `E4.given` on `S-scale1.5` (A 0.0123, C 0.0113, at 0.0100) and olympiad `M4.false.all` on `W-repeat` (A 0.0203, C 0.0234, at 0.0200);
- A alone: olympiad `E4.given` on `F-none` (0.0114 at 0.0100) and school `E3.retention` on `F-model` (0.0380 at 0.0349).

No olympiad `M3.lag` fails, though `engine-1` failed it on `G2` for A and on `G2-fading` for C (the [report](https://github.com/altnovatrail/sim/blob/v0.1.17/judgements/engine-1/report.md) of the judgement `engine-1` of `sim` and the entry on it in its [log](https://github.com/altnovatrail/sim/blob/v0.1.17/docs/decisions.md) at `v0.1.17`).

**The split by the seen solution**, 1,000 learners a cell: the shares of all the items of the cell, the signed lag, and the items of attempts whose task had its solution seen — their share, their mean P* and the model's mean P at issue, and their share in the corridor.

| Cell | Rule | In the corridor / below 0.50 / above 0.95 | `M3.lag` | Seen: share of items | Seen: P* / model's P | Seen: in the corridor |
|---|---|---|---|---|---|---|
| olympiad `G2` | base / A / C | 0.258 / 0.255 / 0.109 · 0.283 / 0.137 / 0.115 · 0.286 / 0.128 / 0.115 | −0.993 / −1.000 / −0.994 | 4.3 / 21.7 / 22.3% | 0.645 / 0.539 · 0.872 / 0.783 · 0.874 / 0.785 | 0.186 / 0.261 / 0.260 |
| olympiad `G2-half` | base / A / C | 0.268 / 0.318 / 0.078 · 0.347 / 0.157 / 0.075 · 0.351 / 0.145 / 0.075 | −0.513 / −0.504 / −0.507 | 4.2 / 24.4 / 25.5% | 0.562 / 0.494 · 0.831 / 0.779 · 0.833 / 0.782 | 0.179 / 0.403 / 0.398 |
| olympiad `G2-fading` | base / A / C | 0.271 / 0.320 / 0.078 · 0.358 / 0.154 / 0.076 · 0.361 / 0.140 / 0.076 | −0.420 / −0.407 / −0.409 | 4.1 / 24.5 / 25.7% | 0.553 / 0.498 · 0.823 / 0.779 · 0.825 / 0.783 | 0.177 / 0.424 / 0.420 |
| olympiad `S-scale1.5` | base / A / C | 0.260 / 0.347 / 0.050 · 0.366 / 0.181 / 0.048 · 0.374 / 0.166 / 0.047 | +0.915 / +0.817 / +0.805 | 5.8 / 28.3 / 29.7% | 0.463 / 0.440 · 0.775 / 0.775 · 0.780 / 0.779 | 0.148 / 0.477 / 0.476 |
| olympiad `F-none` | base / A / C | 0.134 / 0.507 / 0.099 · 0.369 / 0.153 / 0.066 · 0.386 / 0.127 / 0.064 | +0.051 / +0.020 / +0.018 | 16.5 / 56.9 / 58.9% | 0.511 / 0.496 · 0.763 / 0.760 · 0.774 / 0.771 | 0.132 / 0.462 / 0.475 |
| olympiad `F-growth5` | base / A / C | 0.166 / 0.512 / 0.051 · 0.424 / 0.160 / 0.027 · 0.438 / 0.137 / 0.026 | −0.152 / −0.186 / −0.192 | 15.7 / 55.5 / 57.1% | 0.480 / 0.478 · 0.756 / 0.760 · 0.767 / 0.770 | 0.147 / 0.524 / 0.530 |
| school `G2` | base / A / C | 0.130 / 0.123 / 0.414 · 0.122 / 0.115 / 0.446 · 0.122 / 0.115 / 0.446 | −2.449 / −2.477 / −2.476 | 86.1 / 86.3 / 86.3% | 0.858 / 0.688 · 0.870 / 0.708 · 0.870 / 0.708 | 0.123 / 0.115 / 0.115 |
| school `G3` | base / A / C | 0.180 / 0.376 / 0.107 · 0.206 / 0.330 / 0.093 · 0.207 / 0.330 / 0.093 | −0.500 / −0.488 / −0.488 | 86.7 / 86.9 / 86.9% | 0.594 / 0.542 · 0.626 / 0.573 · 0.626 / 0.573 | 0.183 / 0.213 / 0.213 |
| school `F-none` | base / A / C | 0.167 / 0.451 / 0.074 · 0.232 / 0.393 / 0.063 · 0.232 / 0.393 / 0.063 | +0.112 / +0.113 / +0.112 | 99.5 / 99.5 / 99.5% | 0.529 / 0.522 · 0.559 / 0.552 · 0.559 / 0.552 | 0.167 / 0.232 / 0.232 |

On school maths every task outside a lesson the engine gives has its solution seen, on every cell and under each rule; the rest of the items are those of the lessons.

**Olympiad `G2` by group**, base → A, 1,000 learners: the share of the cell's items, the share of them below 0.50, the mean P* and the model's mean P, and the moves of θ + δ and of θ* + δ* through an attempt, from the 101st answer.

| Group | Items | Below 0.50 | P* / model's P | θ + δ / θ* + δ* through an attempt |
|---|---|---|---|---|
| task, `CONSOLIDATION`, unseen | 24.5 → 12.7% | 0.428 → 0.123 | 0.547 / 0.481 → 0.757 / 0.712 | +0.0076 / +0.0106 → +0.0064 / +0.0104 |
| task, `CONSOLIDATION`, seen | 1.6 → 4.1% | 0.548 → 0.009 | 0.477 / 0.380 → 0.875 / 0.775 | +0.0069 / +0.0106 → +0.0072 / +0.0105 |
| task, `QUEUE`, unseen | 20.5 → 15.6% | 0.106 → 0.040 | 0.757 / 0.668 → 0.808 / 0.734 | +0.0075 / +0.0104 → +0.0066 / +0.0104 |
| task, `QUEUE`, seen | 1.4 → 13.6% | 0.230 → 0.002 | 0.681 / 0.548 → 0.872 / 0.785 | +0.0086 / +0.0106 → +0.0064 / +0.0105 |
| task, `ROTATION`, unseen | 13.5 → 12.7% | 0.051 → 0.034 | 0.825 / 0.742 → 0.832 / 0.748 | +0.0050 / +0.0103 → +0.0051 / +0.0103 |
| task, `ROTATION`, seen | 0 → 2.6% | — → 0.001 | — → 0.885 / 0.797 | — → +0.0054 / +0.0103 |
| lesson, `NEW_SKILL` | 20.3 → 20.3% | 0.345 → 0.345 | 0.636 / 0.597 → 0.636 / 0.597 | +0.0113 / +0.0099 → +0.0122 / +0.0099 |

Left out are the probe, repetition and the lessons of the queue and of rotation: about 18% of the items under either rule, and A moves each of them by less than one point. On school `G2` the queue of seen tasks is 64% of the items under both rules; its P* / model's P is 0.872 / 0.708 under the base and 0.903 / 0.756 under A, its share above 0.95 0.538 and 0.603, and through an attempt θ + δ moves by 0.081 and 0.074 while θ* + δ* moves by 0.235 and 0.241: a task of school maths is a page of several items.

**`E4` by the code**, 1,000 learners a cell: the skills within reach that one rule gave at least once and the other never, by the code that first gave them under the rule that did. Under A and C every such skill on olympiad maths was first given by the lesson of a new skill, `NEW_SKILL`; on school maths there are none.

| Cell | Within reach under A | A loses / gains against the base | C loses / gains | `theta` loses / gains against the base before it |
|---|---|---|---|---|
| olympiad `S-scale1.5` | 13,256 | 105 / 27 | 110 / 39 | 0 / 0 |
| olympiad `F-none` | 14,231 | 101 / 31 | 99 / 49 | 41 / 34 |
| olympiad `F-growth5` | 14,194 | 120 / 54 | 134 / 58 | 59 / 50 |
| olympiad `G2` | 17,080 | 146 / 113 | 163 / 113 | 0 / 0 |
| olympiad `G2-half` | 15,786 | 78 / 107 | 100 / 126 | 0 / 0 |
| olympiad `G2-fading` | 15,541 | 62 / 113 | 75 / 129 | 0 / 0 |
| olympiad `G3` | 15,766 | 86 / 95 | 107 / 132 | 0 / 0 |

## 6. The rough look of the variants of E5

The variants D, E and F against the base on the pin of `sim` after `model-2`, as A and C were looked at in measurements 5: a rough look, which decides nothing, by `sim`'s criterion. The rules of the variants and what came of the look are in [architecture.md](architecture.md); the facts of `sim` are [`sim` v0.1.17](https://github.com/altnovatrail/sim/tree/v0.1.17).

**How it was counted.** 2026-11-14, in the container of `engine`: AMD Ryzen AI 9 365, 20 processors, go1.27.4. A clone of `sim` `v0.1.17` (6f5151f; `gh release verify` accepts it) in a temporary folder, and its recipe of a rough look, 300 working learners a cell, the whole matrix: `just rough engine runs/base-300.jsonl` for the base, and `CANDIDATE_ENGINE=<folder> just rough engine runs/<x>-300.jsonl` for each variant, where the folder is a clean clone of `engine` at the variant's commit — D 29e205a, E e9832fc, F 389798b —, which the recipe takes in place of a SHA on GitHub (the README of [`sim` v0.1.15](https://github.com/altnovatrail/sim/blob/v0.1.15/README.md), on judging a candidate). The base built by `go.mod` of `sim` — `engine` `v0.1.6` and `learner-model` `v0.1.8` —, as in measurements 5; each variant in the workspace of the gate, with `learner-model` at ab688b2, and held to `gate.json`. By the clock 9 min 44 s, 6 min 49 s, 12 min 53 s and 7 min 18 s, the base and E beside other work. SHA-256: the base `bcacb08e…125d` — to the byte the base of measurements 5 —, D `04a19d34…23da`, E `fa487a7b…4a92`, F `40408d8d…69f7`. Then `sim compare` of the base with each, and `sim judge` with one judgement file of the three: `{"name": "E", "group": "corridor", "constants": 0, "new_field": false, "distance": 1, "equal_on": ["novatrail/grade-3-math", "novatrail/english-essays"]}`, D the same with `"distance": 2`, F with `"constants": 1, "distance": 3`. The clones and the runs do not go into the repository.

**The rough look.** The cells of olympiad maths of measurements 5 and those where this look failed "no worse", by `sim compare`, 300 learners a cell. On the graphs of school maths alone and of the essays every number of every variant is the base's, to the bit; the pair of olympiad and school maths, `A-linked`, moves with olympiad maths, and the synthetic sheet of `S-predicate`, which no graph of the choice has, moves too. `E1.corridor` is the share of items in the corridor (more is better), `M3.lag` the signed lag, `E4.given` the share of the skills within reach given at least once. Each cell: the base → D / E / F.

| Cell | `E1.corridor` | `M3.lag` | `E4.given` |
|---|---|---|---|
| olympiad `G0` | 0.259 → 0.379 / 0.322 / 0.347 | −0.028 → −0.020 / −0.012 / −0.023 | 0.868 → 0.868 / 0.869 / 0.871 |
| olympiad `G2` | 0.254 → 0.280 / 0.276 / 0.293 | −1.026 → −1.032 / −1.013 / −1.024 | 0.798 → 0.795 / 0.799 / 0.799 |
| olympiad `G2-half` | 0.265 → 0.338 / 0.302 / 0.323 | −0.529 → −0.516 / −0.512 / −0.518 | 0.827 → 0.826 / 0.829 / 0.828 |
| olympiad `G2-fading` | 0.267 → 0.353 / 0.310 / 0.332 | −0.431 → −0.417 / −0.406 / −0.413 | 0.833 → 0.834 / 0.837 / 0.836 |
| olympiad `G3` | 0.256 → 0.324 / 0.293 / 0.315 | −0.616 → −0.618 / −0.611 / −0.608 | 0.831 → 0.832 / 0.834 / 0.834 |
| olympiad `G3-drop` | 0.220 → 0.324 / 0.282 / 0.302 | 0.502 → 0.463 / 0.475 / 0.467 | 0.863 → 0.865 / 0.864 / 0.864 |
| olympiad `S-scale1.5` | 0.267 → 0.363 / 0.309 / 0.341 | 0.929 → 0.844 / 0.878 / 0.854 | 0.889 → 0.882 / 0.889 / 0.890 |
| olympiad `F-none` | 0.137 → 0.351 / 0.221 / 0.282 | 0.053 → 0.024 / 0.059 / 0.040 | 0.897 → 0.891 / 0.899 / 0.893 |
| olympiad `F-growth5` | 0.167 → 0.402 / 0.258 / 0.341 | −0.165 → −0.195 / −0.159 / −0.188 | 0.879 → 0.871 / 0.877 / 0.874 |
| olympiad `F-C1` | 0.150 → 0.381 / 0.242 / 0.321 | −0.103 → −0.101 / −0.111 / −0.090 | 0.900 → 0.895 / 0.898 / 0.895 |

`M3.settle` on olympiad `G3-drop` — 43.8 → 53.4 / 55.8 / 56.3, and on `G3` 51.0 → 54.0 / 56.1 / 54.6, with intervals across zero; on olympiad `W-repeat` `M4.false.all` — 0.032 → 0.047 / 0.041 / 0.037 — and `M3.lag` — 0.023 → 0.165 / 0.092 / 0.088, wider for all three with the interval clear of zero.

`sim judge`: "A rough look, which decides nothing: no candidate holds every constraint and is better than the base". `neutral` holds for all three. The group `corridor` is 0.054 [0.050; 0.059] for D, 0.029 [0.026; 0.031] for E and 0.043 [0.040; 0.047] for F; the score is not read on a rough look. The checks "no worse" whose worst end is past the tolerance:
- all three: olympiad `M3.settle` on `G3-drop` — D 17.00 at 16.70, E 19.28 at 16.50, F 20.25 at 17.77;
- D also: olympiad `E4.given` on `S-scale1.5` (0.0123 at 0.0100), `F-none` (0.0105) and `F-growth5` (0.0119), and `M4.false.all` on `W-repeat` (0.0212 at 0.0200);
- F also: olympiad `E4.given` on `F-C1` (0.0102 at 0.0100).

## 7. The constraint `fair` on the rough look of the variants of E5

The variants D, E and F against the base on `sim` `v0.1.20`, the tag that the judgement `engine-2` pins, with the constraint of fairness `fair`: a rough look, which decides nothing, by `sim`'s criterion. What came of it is in [architecture.md](architecture.md); the facts of `sim` are [`sim` v0.1.20](https://github.com/altnovatrail/sim/tree/v0.1.20), its [log](https://github.com/altnovatrail/sim/blob/v0.1.20/docs/decisions.md) and [RUN.md](https://github.com/altnovatrail/sim/blob/v0.1.20/RUN.md).

**How it was counted.** 2026-11-15, in the container of `engine`: AMD Ryzen AI 9 365, 20 processors, go1.27.4. A clone of `sim` `v0.1.20` (63288d7; `gh release verify` accepts it) in a temporary folder, and its recipe of a rough look, 300 working learners a cell, the whole matrix: `just rough engine runs/base-300.jsonl` for the base, and `CANDIDATE_ENGINE=<SHA> just rough engine runs/<x>-300.jsonl` for each variant by its full SHA — D 29e205a, E e9832fc, F 389798b —, each held to `gate.json`, with `learner-model` at ab688b2. By the clock 7 min 38 s, 8 min 22 s, 7 min 58 s and 8 min 5 s. SHA-256: the base `0aeb0e02…a9bb`, D `154250f4…4c8a`, E `47d99b0e…a84b`, F `92a925e9…629b`. Each file has 20,401 lines: the header, 19,500 lines of the 65 cells of measurements 6 and 900 of the cells `D-grade1`, `D-grade3` and `D-grade5`. Without the header and the lines of `D-`, each file is to the byte the file of the same rule in measurements 6, so every number of measurements 6 stands. Then `sim judge` with the judgement file of measurements 6, and, for the numbers of the checks that pass, which its report does not print, a program in a scratch copy of `sim` that calls `criterion.Judge` and `measure.ReadPair` of the same tag and prints every check of `fair` and the numbers of `M11` of both sides. `sim evaluate` reads only a decisive run of 2,000 and refuses these files. The clones, the program and the runs do not go into the repository.

**The report** of `sim judge` is the report of measurements 6 line for line, with one line more, under D: "fails fair D-grade1/novatrail/olympiad-math M11.band: worse by at most 0.02523, the tolerance 0.0219". The outcome is the same — "A rough look, which decides nothing: no candidate holds every constraint and is better than the base" —, and so are the groups `corridor`, 0.054 for D, 0.029 for E and 0.043 for F.

**`M11`.** On olympiad maths, each grade against the reference, grade 3, 300 learners; `M11.answers` is the answers until the estimate holds within 0.5 of the reference's (fewer is better), `M11.same` the share for whom that takes at most 50 (more is better), `M11.band` the largest difference of the shares of the tasks at a band after the 50th answer (less is better). Each cell: the base → D / E / F, with the 95 % interval of the base.

| Cell | `M11.answers` | `M11.same` | `M11.band` |
|---|---|---|---|
| `D-grade1` | 30.42 [25.3; 35.41] → 29.12 / 30.00 / 29.89 | 0.850 [0.81; 0.89] → 0.853 / 0.850 / 0.847 | 0.0055 [0.0018; 0.0151] → 0.0229 / 0.0154 / 0.0196 |
| `D-grade5` | 74.04 [64.56; 83.65] → 72.97 / 73.94 / 73.04 | 0.610 [0.553; 0.663] → 0.607 / 0.610 / 0.613 | 0.0165 [0.0063; 0.0297] → 0.0110 / 0.0118 / 0.0151 |

**The checks of `fair`** — the worse end of the paired interval, the candidate against the base, and the tolerance, the larger of the author's — 1 answer, 0.01, 0.01 — and 4.3 standard errors of the difference (the entry on `fair` in `sim judge` of `sim`'s [log](https://github.com/altnovatrail/sim/blob/v0.1.20/docs/decisions.md) at `v0.1.20`):

| Check | D | E | F |
|---|---|---|---|
| `D-grade1` `M11.answers` | 0.127 at 3.198 | 0.797 at 2.667 | 0.827 at 3.084 |
| `D-grade1` `M11.same` | 0 at 0.0141 | 0.0100 at 0.0200 | 0.0100 at 0.0138 |
| `D-grade1` `M11.band` | **0.0252 at 0.0219** | 0.0165 at 0.0205 | 0.0208 at 0.0212 |
| `D-grade5` `M11.answers` | 0.563 at 3.698 | 1.363 at 3.213 | 0.477 at 3.272 |
| `D-grade5` `M11.same` | 0.0133 at 0.0238 | 0.0100 at 0.0194 | 0 at 0.0139 |
| `D-grade5` `M11.band` | 0.0070 at 0.0256 | 0.0073 at 0.0258 | 0.0108 at 0.0274 |

E and F hold `fair`, F on `D-grade1` `M11.band` by 0.0004 within the tolerance; D fails it there by 0.0033. On `D-grade1` all three widen `M11.band`, the base's 0.0055 to 0.015–0.023, the paired interval clear of zero for D, [0.0052; 0.0252] worse, and for F, [0.0018; 0.0208], and across it for E, [−0.0025; 0.0165]; on `D-grade5` all three narrow it. `M11.answers` and `M11.same` move within their intervals.

## 8. "Gate" on arm64

The time of the job "Gate" of `ci.yml` on `ubuntu-24.04` and on `ubuntu-24.04-arm`; `sim` `v0.1.20` at 63288d7 on both, 250 learners a cell.

**How it was measured.** 2026-11-16, from the GitHub API, read only: `actions/runs/<id>/jobs`, the start and the end of the job and of its steps. amd64 — the last 6 successful runs of CI on `main`, 2026-11-14 to 2026-11-16 (43647847694, 43664203248, 43666671717, 43799483044, 43805279819, 43810269337), the image of `sim`'s stage from the cache. arm64 — the runs 43812293508 and 43813275834 of pull request #82: in the first the image was built on a miss of the new entry of arm64 and saved, in the second it came from that entry.

| Runner | The job | The step `gate` | The image |
|---|---|---|---|
| `ubuntu-24.04`, 6 runs | 379–654 s, median 441 s | 348–622 s, median 398 s | 20–35 s, from the cache |
| `ubuntu-24.04-arm`, 2 runs | 387 s, 360 s | 342 s, 327 s | 34 s built and 4 s saved; 23 s from the cache |

**The bits.** The table of the gate on arm64, in both runs, is that of the run 43810269337 on amd64 line for line: 69 numbers, each "in band" and "to the bit" against `gate.json`. Not measured: the spread of arm64 over more runs.

## 9. The gate as jobs

The time of the gate of `ci.yml` run as four jobs "Gate / part k" and the check "Gate" after them, against the one job "Gate" before; `ubuntu-24.04-arm` on both, 250 learners a cell. Before: `sim` `v0.1.20` at 63288d7; after: `sim` `v0.1.24` at 33d70f2, whose gate is the same (the entry on the gate as jobs of `sim`'s [log](https://github.com/altnovatrail/sim/blob/v0.1.24/docs/decisions.md) at `v0.1.24`).

**How it was measured.** 2026-11-16, from the GitHub API, read only: `actions/runs/<id>/jobs`, the start and the end of the jobs and of their steps; the gate's path is from the start of its first job to the end of "Gate". Before — the run 43814890320 of CI on `main`, 5db1df5, the merge of "Gate" on arm64, the image of `sim`'s stage from the cache, and the two runs of pull request #82 of measurements 8. After — two runs of pull request #83: 43818734321, on a miss of the cache's new entry — `sim`'s Dockerfile at `v0.1.24` has a new base image, so each part built the image, part 0 saved it, and "Gate" loaded it from the cache —, and 43820905022, after the fixes of the review, every job loading the image from the cache.

| The gate | Its path | The run of the learners | The image |
|---|---|---|---|
| one job, `main` | 375 s | 335 s, the step `gate` | from the cache |
| one job, #82, 2 runs | 387 s, 360 s | 342 s, 327 s | built; from the cache |
| four parts and the check, #83, first run | 196 s: the parts 115–144 s, a wait of 5 s, the check 47 s | 89–94 s a part, the step `gate-part`; 13 s the check, the step `gate-check` | 20–37 s a part, built; 23 s the check, from the cache |
| four parts and the check, #83, second run | 186 s: the parts 122–133 s, a wait of 4 s, the check 48 s | 88–94 s a part; 13 s the check | 22–26 s a part and 24 s the check, all from the cache |

The run of CI, from its start to the end of its last check: 392 s on `main` before, 201 s and 192 s on #83 after; "Gate" was the last check of both. On #83 every other check ended by 06:45:21, "Go / coverage" the last, 109 s after the start of the run, so the gate is still the longest path.

**The bits.** The table of "Gate" in both runs of #83 is that of `just gate` on `sim` `v0.1.20` in the container of `engine` on the same code, 2026-11-16, line for line: 69 numbers, each "in band" and "to the bit" against `gate.json`. Loading the image from the cache takes about as long as building it on this runner, 22–26 s against 20–37 s. Not measured: the spread over more runs.

## 10. CI on `.github` v0.1.30

The time of the jobs of CI and CodeQL that the move to the shared workflows of `.github` v0.1.30 changed, before it and on its pull request.

**How it was measured.** 2026-11-16, from the GitHub API, read only: `actions/runs/<id>/jobs`, the start and the end of the jobs and of their steps, and the wait of a job, from its creation to its start. Before — the run 43822630871 of CI and 43822630818 of CodeQL on `main`, 4d4654d, the merge of the gate as jobs, on `.github` v0.1.24. After — the first run of pull request #84, 43856171341 of CI and 43856171330 of CodeQL: the caches of Go had no entry yet, since only a push to `main` saves one, so every job ran with empty ones; the image of amd64 came from `main`'s entry of `novatrail-toolchain.tar`, which `documents.yml` had saved already, and the gate's arm64 image of `sim` had no entry under the new file and was built in each job. The jobs of the gate took no caches of Go in that run: they take those of `sim`'s `gate.yml` from the fixes of the review, the second run, which had no entry of them either.

| Job | Before, s | After, s | Of it after |
|---|---|---|---|
| "Go / targets" | 85, the five targets one after another, 56 | 5, after the matrix | "Go / target (…)", 43–55 each, the tests 8–13 s each |
| "Go / fma-diff" | 58, the step 32 | 51, the step 25 | |
| "Go / lint", "test", "vuln", "licenses" | 36, 49, 47, 55 | 46, 53, 56, 54 | the image 27–34 s, from the cache; the recipes 12, 19, 17, 15 s with empty caches of Go |
| "Go / toolchain" | — | 11 | the entry looked up, a hit |
| "Contracts" | 55: the image 34, `ci-contracts` 10 | 38: the image 13, `ci-contracts` 10 | the image from the cache; the caches of Go empty |
| "Gate / part k" | 119–137, `gate-part` 90–95 | 121–126, `gate-part` 91–93 | the image built, 22–24 s, as long as it loaded before |
| "Gate" | 57, `gate-check` 14 | 43, `gate-check` 13 | the image built, 22 s |
| "CodeQL / analyze (go)" | 75, the step of Go 20 | 52, the step of Go 2 | go.dev's archive, no image |

The run of CI, from its start to the end of its last job: 259 s before, 345 s after. Its jobs waited for a runner 33 s on average and 76 s at most before, 81 s and 151 s after, with 24 jobs against 17: the matrix of targets adds five jobs and `toolchain` one, and other repositories of the organization ran at the same hours. "Gate" was the last job of both runs, 268 s from the start of its first part after, 253 s before.

The second run, 43858005943, after the fixes of the review, again with no entry of the caches of Go: "Contracts" 51 s — the image 30 s, `ci-contracts` 7 s; "Gate / part k" 118–124 s, `gate-part` 89–92 s, the image built 21–24 s; "Gate" 38 s, `gate-check` 11 s. Its jobs waited for runners longer still, and the run ended 12 minutes after its start.

**Not measured.** A run with the caches of Go filled: the first push to `main` after the merge saves them, and a pull request after it reads them. The share of the queue that the six new jobs cause, apart from the other repositories' runs.

## 11. CI and the release on `spec` v0.18.4 and `.github` v0.1.31

The time of the jobs of CI and of the release that the stage `toolchain-core` and the move to `.github` v0.1.31 could change, before it and on its pull request.

**How it was measured.** 2026-11-16, from the GitHub API, read only: `actions/runs/<id>/jobs`, the start and the end of the jobs and of their steps. Before — the run 43884931449 of CI of pull request #85, b2c32be, on `.github` v0.1.30 and `spec` v0.18.3, with the image and the caches of Go from `main`'s entries; the release — the runs 43861429097, which released `v0.1.12`, and 43886299085, which released nothing, of `main`. After — the first run of pull request #86, 43910838447 of CI and its run of CodeQL: the Dockerfile's hash is new, so the image of `engine` and the caches of Go had no entry, and every job built the stages `toolchain-core` and `toolchain` and its Go from nothing; the gate's image of `sim` came from its entry. The step "Toolchain stages" of `documents.yml` and of the job `toolchain` of `go.yml` named `toolchain`: "The checks without a browser run in the stage toolchain."

| Job | Before, s | After, s | Of it |
|---|---|---|---|
| "Go / lint", "test", "vuln", "licenses" | 43, 33, 31, 48 | 46, 68, 63, 70 | the image 12–14 s → 14–34 s, from the artifact of "Go / toolchain"; the recipes 1–3 s → 10–20 s, with empty caches of Go |
| "Go / arm64", "Go / fma-diff", "Go / rules" | 40, 44, 31 | 64, 99, 51 | `ci-test` of "Go / arm64" 4 → 15 s, with empty caches of Go |
| "Go / target (…)" | 30–43 | 43–68 | the image 13–21 s → 13–30 s |
| "Documents / docs-check-repo" | 45 | 46 | the image 21 → 25 s, built; the check 13 s, the same |
| "Contracts" | 43 | 39 | the image 22 → 21 s, built; `ci-contracts` under 1 → 11 s, with empty caches of Go |
| "Gate / part k" | 117–124 | 121–125 | the image 24–27 s, `gate-part` 83–86 s: `sim`'s, from its entry |
| "Gate" | 42 | 47 | the image 23 s, `gate-check` 5 s |
| "Go / toolchain" | 6 | 53 | the image built on a miss, 26 s; the core named no entry of its own |
| "CodeQL / analyze (go)" | 63 | 49 | the step of Go 3 → 2 s, by the `ARG`s of `toolchain-core` |
| "Release / release" | 50, 59 | — | "Work out the version", which builds the stage, 21 s and 50 s |

The run of CI, from the start of its first job to the end of its last: 173 s before, 573 s after. Its 24 jobs waited for a runner 18 s on average and 54 s at most before, 154 s and 315 s after, at an hour when other repositories of the organization moved to the same tags.

**What it shows.** The stage `toolchain-core` changes no job of CI in `engine`: its `toolchain` holds the same files, the checks take it, as `.github` has it for a caller of the kind `go` ([its decisions at v0.1.31](https://github.com/altnovatrail/.github/blob/v0.1.31/docs/decisions.md)), and the jobs that grew grew by the miss of the new Dockerfile — the image built, the caches of Go empty — and by the queue.

**Not measured.** The release after the move — "Work out the version" and "Build the files" in `toolchain-core` — runs only on `main`, after the merge: the stage holds the same layers as `toolchain`, so its build should take what it took. A run of the pull request with the entries of the new Dockerfile filled: the first push to `main` saves them.

## 12. The fusion of floating-point operations by Go's compiler

Which expressions of the kind the choice rule computes Go's compiler fuses into one operation, on which architectures.

**How it was measured.** 2026-11-07, go1.27.3 in the container, `go build -gcflags=-S` of small functions for amd64 with GOAMD64 v3, arm64, ppc64le, s390x, riscv64 and loong64, and the instructions of the assembler read.

**What came out.**
- `(p1+p2)/2 - 0.775` — `FNMSUBD` on arm64, `FMSUB` on ppc64le and s390x, `FMSUBD` on riscv64 and loong64; `float64((p1+p2)/2) - 0.775` — a separate multiplication and subtraction on all five; on amd64 with GOAMD64 v3 — fused in neither case.
- `s*0.5 - 0.775` — `FNMSUBD` on arm64, `FMSUB` on ppc64le and `FMSUBD` on riscv64.
- `s/float64(n) - 0.775` — `FDIVD` and `FSUBD` on arm64 and riscv64, `FDIV` and `FSUB` on ppc64le: a division by a variable does not fuse.
