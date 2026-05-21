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

import "runtime"

const (
	// DefaultMinItemsPerWorker is the default minimum range size for parallel
	// execution.
	//
	// The value is intentionally conservative because measurement reducers often
	// run cheap slice loops where goroutine startup and final merging can dominate
	// small inputs.
	DefaultMinItemsPerWorker = 64 * 1024

	// DefaultChunkSize is the default grain size for fixed-chunk and
	// dynamic-chunk plans.
	//
	// Callers with unusual locality, cache behavior, or load-balance needs should
	// benchmark their algorithm and override this value with Options.ChunkSize.
	DefaultChunkSize = 64 * 1024
)

// Options configures planning, execution, and merge behavior for one reduction.
//
// Zero Options are valid. NormalizeOptions fills safe defaults before planners
// and runners consume the value. The defaults favor balanced contiguous ranges
// because many measurement reductions are memory-bound and cheap per element;
// callers should benchmark before choosing fine chunk sizes or more workers.
type Options struct {
	// Workers is the maximum number of worker goroutines used by parallel
	// runners. Values less than or equal to zero resolve to runtime.GOMAXPROCS(0)
	// during normalization. Chunk strategies cap this further by chunk count.
	Workers int

	// MinItemsPerWorker prevents parallel execution when the input is too small
	// to amortize worker startup, synchronization, and merge costs. It is a
	// fallback policy, not a planning correctness requirement.
	MinItemsPerWorker int

	// ChunkSize controls StrategyFixedChunks and StrategyDynamicChunks grain
	// size. Smaller chunks can improve load balance or locality experiments but
	// can also increase mapper and merge overhead. ChunkSize is not a hard
	// maximum callback range length: sequential fallback may still pass one full
	// [0:n) range even when a chunk strategy is selected. Callbacks must handle
	// any valid range they receive. ChunkSize is ignored by StrategyBalanced and
	// StrategySequential.
	ChunkSize int

	// Strategy selects the range planning and execution strategy. StrategyAuto
	// normalizes to StrategyBalanced.
	Strategy Strategy

	// MergeMode selects how completed partials are folded. It affects only final
	// partial-result folding, not range planning or mapper scheduling. The zero
	// value is MergeLinear.
	MergeMode MergeMode
}

// NormalizeOptions returns options with invalid or zero values replaced by core
// defaults.
//
// NormalizeOptions does not mutate its input. It also resolves StrategyAuto so
// downstream packages can switch on a concrete strategy without repeating
// defaulting logic.
func NormalizeOptions(opts Options) Options {
	if opts.Workers <= 0 {
		opts.Workers = runtime.GOMAXPROCS(0)
	}
	if opts.Workers < 1 {
		opts.Workers = 1
	}
	if opts.MinItemsPerWorker <= 0 {
		opts.MinItemsPerWorker = DefaultMinItemsPerWorker
	}
	if opts.ChunkSize <= 0 {
		opts.ChunkSize = DefaultChunkSize
	}
	if opts.Strategy == StrategyAuto {
		opts.Strategy = StrategyBalanced
	}
	return opts
}
