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

// ReduceInto executes a reduction using a mapper that does not need worker-slot
// information.
//
// ReduceIndexedInto owns all strategy dispatch. Keeping ReduceInto as a thin
// adapter makes the exported entry points consistent and prevents one strategy
// path from gaining different fallback behavior than the indexed path.
func ReduceInto[T any](n int, opts core.Options, scratch *core.Scratch[T], mapRange core.IntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
	return ReduceIndexedInto(n, opts, scratch, func(_ int, r core.Range, dst *T) {
		mapRange(r, dst)
	}, mergeFn)
}
