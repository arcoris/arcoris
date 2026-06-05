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

import (
	"fmt"
	"sort"
)

// NewVector returns a canonical immutable vector from entries.
func NewVector(entries ...Entry) (Vector, error) {
	if len(entries) == 0 {
		return Vector{}, nil
	}

	copied := append([]Entry(nil), entries...)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Resource < copied[j].Resource
	})

	for i, entry := range copied {
		entryPath := fmt.Sprintf("entries[%d]", i)

		if !entry.Resource.IsValid() {
			return Vector{}, errorAt(
				entryPath+".resource",
				ErrInvalidResource,
				ErrorReasonInvalidResource,
				fmt.Sprintf("resource %q must be dot-separated lower_snake_case", entry.Resource),
			)
		}

		if entry.Amount.IsZero() {
			return Vector{}, errorAt(
				entryPath+".amount",
				ErrZeroAmount,
				ErrorReasonZeroAmount,
				"amount must be positive",
			)
		}

		if i > 0 && copied[i-1].Resource == entry.Resource {
			return Vector{}, errorAt(
				entryPath+".resource",
				ErrDuplicateResource,
				ErrorReasonDuplicateResource,
				fmt.Sprintf("resource %q appears more than once", entry.Resource),
			)
		}
	}

	return Vector{entries: copied}, nil
}

// MustVector returns NewVector(entries...) or panics when entries are invalid.
func MustVector(entries ...Entry) Vector {
	vector, err := NewVector(entries...)
	if err != nil {
		panic(err)
	}
	return vector
}

// vectorFromSorted returns a vector from already canonicalized entries.
func vectorFromSorted(entries []Entry) Vector {
	return Vector{entries: append([]Entry(nil), entries...)}
}
