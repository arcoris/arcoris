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

// compactUsedPartials compacts active partial slots in worker-index order.
//
// The function never allocates and does not preserve inactive slots. The two
// slices must have identical length because they describe the same worker-slot
// storage; a mismatch indicates an internal runner bug and panics.
func compactUsedPartials[T any](partials []T, used []bool) []T {
	if len(partials) != len(used) {
		panic("reduce runner: partial and used slices have different lengths")
	}
	write := 0
	for read := range partials {
		if !used[read] {
			continue
		}
		if write != read {
			partials[write] = partials[read]
		}
		write++
	}
	return partials[:write]
}
