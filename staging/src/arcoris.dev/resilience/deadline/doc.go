/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

// Package deadline provides local execution-budget helpers for resilience
// components.
//
// The package inspects context deadlines, computes remaining time budgets,
// decides whether work may start, clamps requested durations to the remaining
// budget, reserves tail budget for caller-owned cleanup, and derives child
// contexts whose deadlines do not exceed their parent deadline.
//
// Reserve returns three values:
//   - duration: a finite derived duration, meaningful only when bounded and ok
//     are both true;
//   - bounded: whether duration comes from a parent context deadline;
//   - ok: whether the observed deadline and runtime context state still permit
//     continuing.
//
// Reserve does not choose fallback timeouts for unbounded contexts. When
// bounded is false, callers must apply their own timeout policy if they still
// need a finite child budget.
//
// Deadline is a lower-level resilience primitive. Retry loops, queue waits,
// bulkhead acquisition, circuit-breaker probes, cooldown paths, and admission
// checks can use this package to avoid starting work that cannot complete inside
// the caller's remaining execution budget.
//
// The package does not execute operations, retry failed work, sleep, create
// timers, randomize delays, classify errors, propagate deadlines across HTTP or
// gRPC boundaries, export metrics, log, trace, or make health, lifecycle,
// routing, admission, scheduling, or overload-control decisions.
//
// All time calculations are explicit: callers pass the observation time to each
// operation. This keeps the package deterministic in tests and prevents it from
// owning a clock, timer, ticker, goroutine, or runtime loop.
//
// Nil contexts and negative durations are programming errors. The package
// panics at the API boundary for those inputs instead of returning values that
// fail later inside retry, waiting, or admission code.
package deadline
