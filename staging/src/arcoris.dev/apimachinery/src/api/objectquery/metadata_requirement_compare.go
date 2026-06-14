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

package objectquery

import (
	"slices"
	"strings"
)

// compareMetadataRequirement orders requirements by key, then operator, then
// canonical values. It is the single ordering used by compiled selectors.
func compareMetadataRequirement(left, right metadataRequirement) int {
	if cmp := strings.Compare(left.key, right.key); cmp != 0 {
		return cmp
	}
	if left.op < right.op {
		return -1
	}
	if left.op > right.op {
		return 1
	}

	return compareStringSlices(left.values, right.values)
}

// sameMetadataRequirement reports exact canonical requirement equality.
func sameMetadataRequirement(left, right metadataRequirement) bool {
	return left.key == right.key &&
		left.op == right.op &&
		slices.Equal(left.values, right.values)
}

// compareStringSlices orders string slices lexicographically by content and
// then by length.
func compareStringSlices(left, right []string) int {
	for i := 0; i < len(left) && i < len(right); i++ {
		if cmp := strings.Compare(left[i], right[i]); cmp != 0 {
			return cmp
		}
	}
	switch {
	case len(left) < len(right):
		return -1
	case len(left) > len(right):
		return 1
	default:
		return 0
	}
}
