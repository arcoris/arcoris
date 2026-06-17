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

package objectstorewatch

import "arcoris.dev/apimachinery/api/objectstore"

// Boundary is the collection/revision pair from which a watch can continue.
//
// Boundary revision may be zero. Zero is the initial list-to-watch boundary and
// means a subsequent watch should deliver matching committed changes with
// revision greater than zero. Boundary validates shape only; it does not prove
// history availability in a concrete source.
type Boundary struct {
	// collection is the structural collection to continue watching.
	collection objectstore.ListRequest
	// revision is the source-local boundary revision for continuation.
	revision objectstore.Revision
}

// NewBoundary validates collection and stores revision as a watch boundary.
func NewBoundary(collection objectstore.ListRequest, revision objectstore.Revision) (Boundary, error) {
	boundary := Boundary{collection: collection, revision: revision}
	if err := boundary.Validate(); err != nil {
		return Boundary{}, err
	}

	return boundary, nil
}

// Collection returns the structural collection associated with the boundary.
func (b Boundary) Collection() objectstore.ListRequest {
	return b.collection
}

// Revision returns the source-local revision boundary.
func (b Boundary) Revision() objectstore.Revision {
	return b.revision
}

// Clone returns an independent value copy of b.
func (b Boundary) Clone() Boundary {
	return b
}

// IsZero reports whether b contains no collection and no revision.
func (b Boundary) IsZero() bool {
	return b.collection == (objectstore.ListRequest{}) && b.revision.IsZero()
}

// IsValid reports whether b passes boundary validation.
func (b Boundary) IsValid() bool {
	return b.Validate() == nil
}
