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

// Demand is a non-empty immutable resource vector used for reservation checks.
//
// Demand is separate from Vector because empty vectors are useful for limits and
// snapshots, but empty reservation attempts are ambiguous and invalid.
type Demand struct {
	// vector is canonical and non-empty when the demand is valid.
	vector Vector
}

// NewDemand returns a non-empty demand from entries.
func NewDemand(entries ...Entry) (Demand, error) {
	if len(entries) == 0 {
		return Demand{}, errorAt(
			"entries",
			ErrEmptyDemand,
			"demand must contain at least one resource",
		)
	}
	vector, err := NewVector(entries...)
	if err != nil {
		return Demand{}, err
	}
	if vector.IsZero() {
		return Demand{}, errorAt(
			"entries",
			ErrEmptyDemand,
			"demand must contain at least one resource",
		)
	}
	return Demand{vector: vector}, nil
}

// MustDemand returns NewDemand(entries...) or panics when entries are invalid.
func MustDemand(entries ...Entry) Demand {
	demand, err := NewDemand(entries...)
	if err != nil {
		panic(err)
	}
	return demand
}

// IsValid reports whether d contains a non-empty valid vector.
func (d Demand) IsValid() bool {
	return !d.vector.IsZero() && d.vector.IsValid()
}

// Len reports the number of resources in d.
func (d Demand) Len() int {
	return d.vector.Len()
}

// Vector returns d as an immutable vector.
func (d Demand) Vector() Vector {
	return vectorFromSorted(d.vector.entries)
}

// Entries returns a detached copy of d's canonical entries.
func (d Demand) Entries() []Entry {
	return d.vector.Entries()
}

// Amount returns the demanded amount for resource, or zero when absent.
func (d Demand) Amount(resource Resource) Amount {
	return d.vector.Amount(resource)
}

// Has reports whether d contains resource.
func (d Demand) Has(resource Resource) bool {
	return d.vector.Has(resource)
}
