// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


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
	// NormalizeOptions currently resolves it to StrategyBalanced because
	// contiguous ranges have the lowest scheduling overhead for uniform
	// in-memory reductions.
	StrategyAuto Strategy = iota

	// StrategySequential forces one [0:n) range in the current goroutine. It
	// produces one partial and performs no merge.
	StrategySequential

	// StrategyBalanced splits work into a bounded number of balanced contiguous
	// ranges. It is the default for uniform per-element work and gives the
	// Reduce family deterministic range-order merge input.
	StrategyBalanced

	// StrategyFixedChunks splits work into fixed-size chunks assigned to workers
	// deterministically. It is useful for grain-size tuning and locality-sensitive
	// reductions that do not need load-balancing through an atomic cursor.
	StrategyFixedChunks

	// StrategyDynamicChunks lets workers claim fixed-size chunks from an atomic
	// cursor. It is intended for variable-cost chunks and can produce different
	// floating-point grouping than balanced range execution.
	StrategyDynamicChunks
)
