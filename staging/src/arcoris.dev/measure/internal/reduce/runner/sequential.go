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

package runner

import "arcoris.dev/measure/internal/reduce/core"

// reduceSequentially maps the whole input as one non-indexed range.
//
// Sequential execution has exactly one partial and does not call the merger.
// It is used both for explicit StrategySequential and for automatic fallback
// when the input is too small to amortize worker startup.
func reduceSequentially[T any](n int, mapRange core.IntoMapper[T]) (T, bool) {
	return reduceSequentiallyIndexed(n, func(_ int, r core.Range, dst *T) {
		mapRange(r, dst)
	})
}

// reduceSequentiallyIndexed maps the whole input as one range in worker slot
// zero.
//
// Keeping the indexed form separate lets ReduceIndexedInto preserve worker-slot
// semantics even when a parallel strategy falls back to the sequential path.
func reduceSequentiallyIndexed[T any](n int, mapRange core.IndexedIntoMapper[T]) (T, bool) {
	var partial T
	mapRange(0, core.Range{Start: 0, End: n}, &partial)
	return partial, true
}
