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

// Reduce maps each planned range to a complete partial result and merges
// partials.
//
// Reduce is the simplest entry point when the mapper naturally returns a
// complete partial value. It allocates temporary scratch as needed; repeated
// callers should prefer ReduceInto with caller-owned Scratch or Runner.
func Reduce[T any](n int, opts core.Options, mapRange core.Mapper[T], mergeFn core.Merger[T]) (T, bool) {
	return ReduceInto(n, opts, nil, func(r core.Range, dst *T) { *dst = mapRange(r) }, mergeFn)
}
