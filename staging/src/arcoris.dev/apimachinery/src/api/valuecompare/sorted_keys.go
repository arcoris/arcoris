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

package valuecompare

import (
	"slices"

	"arcoris.dev/apimachinery/api/value"
)

// unionSortedKeys returns canonical traversal order for dynamic object members.
func unionSortedKeys(left, right map[string]value.Value) []string {
	return unionSortedMapKeys(left, right)
}

// unionSortedMapKeys returns sorted keys present in either lookup.
//
// Fast paths avoid the temporary "seen" map when one side is empty. Callers use
// this for deterministic diff traversal; equality-only paths avoid it when they
// can compare by direct lookup.
func unionSortedMapKeys[T any](left, right map[string]T) []string {
	switch {
	case len(left) == 0:
		return sortedMapKeys(right)
	case len(right) == 0:
		return sortedMapKeys(left)
	}

	seen := make(map[string]struct{}, len(left)+len(right))
	for key := range left {
		seen[key] = struct{}{}
	}
	for key := range right {
		seen[key] = struct{}{}
	}

	return sortedKeysFromSet(seen)
}

// sortedMapKeys returns sorted keys from one lookup.
func sortedMapKeys[T any](items map[string]T) []string {
	if len(items) == 0 {
		return nil
	}

	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	return keys
}

// sortedKeysFromSet returns sorted keys from a string set.
func sortedKeysFromSet(seen map[string]struct{}) []string {
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	return keys
}
