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

// Sequential returns a single range covering [0:n).
//
// For n <= 0 it returns dst truncated to zero length. The function never
// allocates when dst has enough capacity for one range.
func Sequential(n int, dst []core.Range) []core.Range {
	dst = dst[:0]
	if n <= 0 {
		return dst
	}
	return append(dst, core.Range{Start: 0, End: n})
}
