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

import "testing"

func TestHasDuplicates(t *testing.T) {
	requireEqual(t, hasDuplicates([]int8{}), false)
	requireEqual(t, hasDuplicates([]int8{1}), false)
	requireEqual(t, hasDuplicates([]int8{1, 2}), false)
	requireEqual(t, hasDuplicates([]int8{1, 1}), true)
	requireEqual(t, hasDuplicates([]string{"a", "b"}), false)
	requireEqual(t, hasDuplicates([]string{"a", "a"}), true)
}

func TestEnumBounds(t *testing.T) {
	requireEqual(t, enumBelowMin([]int8{0}, limit[int8]{}), false)
	requireEqual(t, enumAboveMax([]uint64{10}, limit[uint64]{}), false)
	requireEqual(t, enumBelowMin([]int8{-1, 0}, limit[int8]{value: 0, set: true}), true)
	requireEqual(t, enumAboveMax([]uint64{1, 2}, limit[uint64]{value: 1, set: true}), true)
	requireEqual(t, enumBelowMin([]int8{0, 1}, limit[int8]{value: 0, set: true}), false)
	requireEqual(t, enumAboveMax([]uint64{0, 1}, limit[uint64]{value: 1, set: true}), false)
}
