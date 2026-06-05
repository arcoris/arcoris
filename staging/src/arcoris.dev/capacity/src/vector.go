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

package capacity

import "sort"

// Vector is a canonical immutable resource vector.
//
// The zero value is a valid empty vector. Non-empty vectors contain valid
// resources, positive amounts, no duplicate resources, and entries sorted by
// resource string. Public methods return copies so callers cannot mutate vector
// internals.
type Vector struct {
	// entries is sorted by Resource and never contains zero amounts.
	entries []Entry
}

// IsZero reports whether v has no resource entries.
func (v Vector) IsZero() bool {
	return len(v.entries) == 0
}

// IsValid reports whether v is canonical.
func (v Vector) IsValid() bool {
	for i, entry := range v.entries {
		if !entry.IsValid() {
			return false
		}

		if i > 0 && v.entries[i-1].Resource >= entry.Resource {
			return false
		}
	}

	return true
}

// Len reports the number of resources in v.
func (v Vector) Len() int {
	return len(v.entries)
}

// Amount returns the amount for resource, or zero when absent or invalid.
func (v Vector) Amount(resource Resource) Amount {
	if !resource.IsValid() {
		return 0
	}

	i, ok := v.find(resource)
	if !ok {
		return 0
	}

	return v.entries[i].Amount
}

// Has reports whether v contains resource.
func (v Vector) Has(resource Resource) bool {
	if !resource.IsValid() {
		return false
	}

	_, ok := v.find(resource)

	return ok
}

// Entries returns a detached copy of v's canonical entries.
func (v Vector) Entries() []Entry {
	return append([]Entry(nil), v.entries...)
}

// Equal reports whether v and other contain exactly the same canonical entries.
func (v Vector) Equal(other Vector) bool {
	if len(v.entries) != len(other.entries) {
		return false
	}

	for i := range v.entries {
		if v.entries[i] != other.entries[i] {
			return false
		}
	}

	return true
}

// find returns the canonical entry index for resource.
func (v Vector) find(resource Resource) (int, bool) {
	i := sort.Search(len(v.entries), func(i int) bool {
		return v.entries[i].Resource >= resource
	})

	return i, i < len(v.entries) && v.entries[i].Resource == resource
}
