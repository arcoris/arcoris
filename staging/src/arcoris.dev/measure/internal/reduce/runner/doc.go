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
// ReduceIndexedInto are the map-then-merge entry points. AccumulateInto and
// AccumulateIndexedInto are the worker-local accumulation entry points.
// Strategy-specific paths stay private so dispatch rules remain visible in one
// place per API family.
//
// The Reduce family calls mappers with a fresh destination for one range or
// chunk. Mappers may assign dst, and the runner merges produced partials later.
// Chunk-heavy fixed and dynamic strategies may therefore call mergeFn once per
// chunk when they fold chunk partials into worker-local storage. The Accumulate
// family calls accumulators directly on worker-local partials. Accumulators may
// receive the same dst many times and must add to existing state, avoiding a
// chunkPartial plus mergeFn call for every chunk.
//
// Execution model:
//
//   - StrategySequential plans one [0:n) range, executes in the current
//     goroutine, produces one partial, and performs no merge.
//   - StrategyBalanced plans balanced contiguous ranges. Reduce uses
//     range-local partials merged in range order; Accumulate assigns balanced
//     ranges to worker-local partials. The current balanced planner produces at
//     most one planned range per worker, so no over-partitioned queued range
//     filler is active in normal dispatch.
//   - StrategyFixedChunks uses deterministic contiguous chunk blocks per worker.
//     Reduce folds complete chunk partials into worker partials; Accumulate
//     updates worker partials directly. No atomic cursor is used.
//   - StrategyDynamicChunks uses worker-local execution with atomic chunk
//     claiming. Idle workers are compacted away before merge.
//
// Worker-local paths never assume that the zero value of a generic partial is a
// merge identity. A worker publishes a partial only after claiming at least one
// chunk, and compactUsedPartials removes inactive slots before folding.
// Accumulate-family worker partials still start from the zero value of T, so
// accumulators that need internal maps, slices, or buffers must lazily
// initialize them on first use.

// ChunkSize is an execution-grain hint for chunk strategies, not a hard maximum
// callback range length. Sequential fallback may still call a mapper or
// accumulator once with the full [0:n) range when the input is too small to
// justify parallelism.

// Pairwise merge uses scratch partial storage as working memory. When partials
// contain references, intermediate values can remain retained in scratch slots
// until the caller uses Scratch.Clear, Runner.Clear, Scratch.Release, or
// Runner.Release.
//
// This kernel deliberately avoids context cancellation, error aggregation, panic
// recovery, and worker pools on the hot path. Callers that need those policies
// should wrap the mapper before entering runner. The module targets Go 1.25, so
// loops rely on the language's per-iteration loop variable semantics instead of
// legacy closure-capture copies.
package runner
