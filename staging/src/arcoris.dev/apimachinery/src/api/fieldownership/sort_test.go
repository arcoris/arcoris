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

package fieldownership

import (
	"slices"
	"testing"
)

func TestCompareOwners(t *testing.T) {
	requireEqual(t, compareOwners("a", "b") < 0, true)
	requireEqual(t, compareOwners("b", "a") > 0, true)
	requireEqual(t, compareOwners("a", "a"), 0)
}

func TestCompactSortedOwners(t *testing.T) {
	got := compactSortedOwners([]Owner{"a", "a", "b", "b", "c"})

	requireOwners(t, got, "a", "b", "c")
}

func TestCompactSortedOwnersEmptyIsNil(t *testing.T) {
	requireEqual(t, compactSortedOwners(nil) == nil, true)
}

func TestSortedOwnersUseCompareOwners(t *testing.T) {
	owners := []Owner{"z", "a", "m"}

	slices.SortFunc(owners, compareOwners)

	requireOwners(t, owners, "a", "m", "z")
}
