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


package runner

import "arcoris.dev/measure/internal/reduce/core"

// activeWorkers caps requested worker count to available jobs while preserving a
// valid positive worker count for non-empty job sets.
func activeWorkers(workers, jobs int) int {
	if jobs <= 0 {
		return 0
	}
	if workers <= 0 || workers > jobs {
		return jobs
	}
	return workers
}

// enoughItemsForParallel reports whether an input is large enough to justify
// starting more than one reduction worker.
func enoughItemsForParallel(n int, minItemsPerWorker int) bool {
	if n <= 0 {
		return false
	}
	if minItemsPerWorker <= 0 {
		return true
	}
	// Use division instead of minItemsPerWorker*2 to avoid overflow for very
	// large thresholds.
	return n/2 >= minItemsPerWorker
}

// shouldReduceSequentially reports whether an automatic fallback should bypass
// worker startup and execute as one full-input range.
//
// Explicit StrategySequential is handled by ReduceIndexedInto before this
// policy check so strategy dispatch remains easy to audit.
func shouldReduceSequentially(n int, opts core.Options) bool {
	return opts.Workers <= 1 || !enoughItemsForParallel(n, opts.MinItemsPerWorker)
}
