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

package planner

import "arcoris.dev/measure/internal/reduce/core"

// Static returns a stable contiguous range plan for uniform-cost reductions.
//
// Static uses MinItemsPerWorker to avoid parallel plans that are too fine to
// amortize worker startup and merge costs. It emits ranges in increasing index
// order, caps the range count by Workers, and keeps range sizes as balanced as
// integer division allows. For small inputs it returns one range so runners can
// use the sequential fast path.
func Static(n int, opts core.Options, dst []core.Range) []core.Range {
	dst = dst[:0]
	if n <= 0 {
		return dst
	}

	opts = core.NormalizeOptions(opts)
	workers := opts.Workers
	if workers <= 1 || n/2 < opts.MinItemsPerWorker {
		return append(dst, core.Range{Start: 0, End: n})
	}
	if workers > n {
		workers = n
	}

	maxWorkersByGrain := n / opts.MinItemsPerWorker
	if maxWorkersByGrain < 1 {
		maxWorkersByGrain = 1
	}
	if workers > maxWorkersByGrain {
		workers = maxWorkersByGrain
	}
	if workers <= 1 {
		return append(dst, core.Range{Start: 0, End: n})
	}
	if cap(dst) < workers {
		dst = make([]core.Range, 0, workers)
	}

	base := n / workers
	rem := n % workers

	start := 0
	for worker := 0; worker < workers; worker++ {
		size := base
		if worker < rem {
			size++
		}
		end := start + size
		dst = append(dst, core.Range{Start: start, End: end})
		start = end
	}
	return dst
}
