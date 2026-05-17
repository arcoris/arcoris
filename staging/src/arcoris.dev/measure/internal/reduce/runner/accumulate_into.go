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

// AccumulateInto executes a reduction with a non-indexed worker-local
// accumulator.
//
// Unlike ReduceInto, the accumulator may receive the same dst many times and
// must add to existing state. This avoids one chunkPartial and one mergeFn call
// per chunk for algorithms that naturally update reusable partial storage.
func AccumulateInto[T any](n int, opts core.Options, scratch *core.Scratch[T], accumulate core.Accumulator[T], mergeFn core.Merger[T]) (T, bool) {
	return AccumulateIndexedInto(n, opts, scratch, func(_ int, r core.Range, dst *T) {
		accumulate(r, dst)
	}, mergeFn)
}
