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

package merge

import "arcoris.dev/measure/internal/reduce"

// PairwiseInPlace merges partials in pairwise rounds.
//
// PairwiseInPlace returns false for an empty slice. For non-empty input it
// reuses partials as working storage, reducing adjacent pairs until one value
// remains. Callers must not depend on partials preserving its original contents
// after this call.
func PairwiseInPlace[T any](partials []T, mergeFn reduce.Merger[T]) (T, bool) {
	var zero T
	if len(partials) == 0 {
		return zero, false
	}
	n := len(partials)
	for n > 1 {
		write := 0
		for read := 0; read+1 < n; read += 2 {
			dst := partials[read]
			mergeFn(&dst, partials[read+1])
			partials[write] = dst
			write++
		}
		if n%2 == 1 {
			partials[write] = partials[n-1]
			write++
		}
		n = write
	}
	return partials[0], true
}
