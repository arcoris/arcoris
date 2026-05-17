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

// Package runner executes in-memory range reductions.
//
// Runner functions are responsible for strategy normalization, worker startup,
// partial-result storage, and final merge. Reduce, ReduceInto, and
// ReduceIndexedInto are the only exported entry points; strategy-specific paths
// stay private so all dispatch rules remain visible in one place.
//
// Execution model:
//
//   - StrategySequential plans one [0:n) range, executes in the current
//     goroutine, produces one partial, and performs no merge.
//   - StrategyStatic plans balanced contiguous ranges and uses range-local
//     partials: one partial slot per planned range, merged in range order for
//     deterministic non-commutative folding.
//   - StrategyFixed uses fixed chunk sizes but worker-local execution by
//     default: one partial per active worker, merged in active worker order.
//     This avoids one partial per chunk when fine chunk sizes create far more
//     chunks than workers.
//   - StrategyDynamic also uses worker-local execution. Workers claim chunks
//     from an atomic cursor, fill one active partial each, and idle workers are
//     compacted away before merge.
//
// Worker-local paths never assume that the zero value of a generic partial is a
// merge identity. A worker publishes a partial only after claiming at least one
// chunk, and compactUsedPartials removes inactive slots before folding.
//
// This kernel deliberately avoids context cancellation, error aggregation, panic
// recovery, and worker pools on the hot path. Callers that need those policies
// should wrap the mapper before entering runner. The module targets Go 1.25, so
// loops rely on the language's per-iteration loop variable semantics instead of
// legacy closure-capture copies.
package runner
