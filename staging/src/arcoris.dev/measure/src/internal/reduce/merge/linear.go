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


package merge

import "arcoris.dev/measure/internal/reduce/core"

// Linear merges partials from left to right.
//
// Linear returns false for an empty slice. For non-empty input it copies the
// first partial into the result and folds the remaining partials into that
// result without modifying the input slice.
func Linear[T any](
	partials []T,
	mergeFn core.Merger[T],
) (T, bool) {
	var zero T
	if len(partials) == 0 {
		return zero, false
	}
	result := partials[0]
	for i := 1; i < len(partials); i++ {
		mergeFn(&result, partials[i])
	}
	return result, true
}
