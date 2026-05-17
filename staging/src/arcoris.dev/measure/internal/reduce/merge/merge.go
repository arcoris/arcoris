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

import "arcoris.dev/measure/internal/reduce/core"

// Merge dispatches to the merge algorithm selected by mode.
//
// The input slice order is the merge order contract. Unknown modes fall back to
// linear merging so the zero value and any future invalid values remain safe.
func Merge[T any](partials []T, mode core.MergeMode, mergeFn core.Merger[T]) (T, bool) {
	if mode == core.MergePairwise {
		return PairwiseInPlace(partials, mergeFn)
	}
	return Linear(partials, mergeFn)
}
