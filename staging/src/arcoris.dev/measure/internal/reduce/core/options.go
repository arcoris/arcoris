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

	// DefaultChunkSize is the default grain size for fixed and dynamic plans.
	//
	// Callers with unusual locality, cache behavior, or load-balance needs should
	// benchmark their algorithm and override this value with Options.ChunkSize.
	DefaultChunkSize = 64 * 1024
)

// Options configures planning, execution, and merge behavior for one reduction.
//
// Zero Options are valid. NormalizeOptions fills safe defaults before planners
// and runners consume the value.
type Options struct {
	// Workers is the maximum number of worker goroutines used by parallel
	// runners. Values less than or equal to zero resolve to runtime.GOMAXPROCS(0)
	// during normalization.
	Workers int

	// MinItemsPerWorker prevents parallel execution when the input is too small
	// to amortize worker startup, synchronization, and merge costs.
	MinItemsPerWorker int

	// ChunkSize controls StrategyFixed and StrategyDynamic grain size. It is
	// ignored by StrategyStatic and StrategySequential.
	ChunkSize int

	// Strategy selects the range planning and execution strategy.
	Strategy Strategy

	// MergeMode selects the post-worker merge algorithm. The zero value is
	// MergeLinear.
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
		opts.Strategy = StrategyStatic
	}
	return opts
}
