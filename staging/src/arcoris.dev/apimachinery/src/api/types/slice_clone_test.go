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

func TestCloneSlice(t *testing.T) {
	requireEqual(t, cloneSlice[int8](nil) == nil, true)
	requireEqual(t, cloneSlice[string]([]string{}) == nil, true)

	ints := []int8{1}
	clonedInts := cloneSlice(ints)
	clonedInts[0] = 2
	requireEqual(t, ints[0], int8(1))

	strings := []string{"a"}
	clonedStrings := cloneSlice(strings)
	clonedStrings[0] = "b"
	requireEqual(t, strings[0], "a")

	names := []FieldName{"name"}
	clonedNames := cloneSlice(names)
	clonedNames[0] = "changed"
	requireEqual(t, names[0], FieldName("name"))
}
