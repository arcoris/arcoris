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

package types

// cloneSlice returns a detached copy of values.
//
// Empty slices normalize to nil because descriptors do not distinguish an
// absent ordered rule from an explicitly empty ordered rule. Exact payload clone
// helpers call this function while preserving their domain-specific names.
func cloneSlice[T any](values []T) []T {
	if len(values) == 0 {
		return nil
	}
	return append([]T(nil), values...)
}
