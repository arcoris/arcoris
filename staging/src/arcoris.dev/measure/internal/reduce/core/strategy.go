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

package core

// Strategy selects how an input interval is split and scheduled.
//
// Strategy values are part of the core reduction contract. Planner and runner
// packages interpret the same values so callers can configure planning and
// execution with one Options value.
type Strategy uint8

const (
	// StrategyAuto lets implementations select the package default.
	//
	// NormalizeOptions currently resolves it to StrategyStatic because
	// contiguous ranges have the lowest scheduling overhead for uniform
	// in-memory reductions.
	StrategyAuto Strategy = iota

	// StrategySequential forces one range covering the full input and avoids
	// worker goroutines.
	StrategySequential

	// StrategyStatic splits work into a bounded number of balanced contiguous
	// ranges. It is the default for uniform per-element work.
	StrategyStatic

	// StrategyFixed splits work into fixed-size chunks. It is useful for
	// grain-size tuning and for work where fixed chunks improve locality.
	StrategyFixed

	// StrategyDynamic lets workers claim fixed-size chunks from a shared cursor.
	// It is intended for variable-cost chunks and can produce different
	// floating-point grouping than static range execution.
	StrategyDynamic
)

// MergeMode selects how completed partial results are combined.
//
// Merge modes operate on already-computed partials. They do not affect planning
// or mapper scheduling.
type MergeMode uint8

const (
	// MergeLinear merges partial results from left to right in slice order.
	MergeLinear MergeMode = iota

	// MergePairwise merges partial results in pairwise rounds, reusing the
	// partial slice as working storage.
	MergePairwise
)
