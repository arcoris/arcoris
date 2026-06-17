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

// Snapshot is a validated collection read plus its list-to-watch boundary.
//
// Snapshot preserves the exact objectstore.ListRequest used for validation and
// stores a detached clone of the list result. It is immutable by API
// convention: methods that expose list data return detached copies.
type Snapshot struct {
	// collection is the structural collection that produced result.
	collection objectstore.ListRequest
	// result is the detached validated list output for collection.
	result objectstore.ListResult
}

// NewSnapshot validates collection/result and stores a detached snapshot.
//
// Validation runs before cloning so malformed large results are rejected
// without first deep-copying their item payloads. Successful construction still
// detaches the stored result from caller-owned data.
func NewSnapshot(collection objectstore.ListRequest, result objectstore.ListResult) (Snapshot, error) {
	if err := objectstore.ValidateListRequest(collection); err != nil {
		return Snapshot{}, errorFor("snapshot.collection", ErrorReasonInvalidSnapshot, ErrInvalidSnapshot, err)
	}
	if err := objectstore.ValidateListResult(collection, result); err != nil {
		return Snapshot{}, errorFor("snapshot.result", ErrorReasonInvalidSnapshot, ErrInvalidSnapshot, err)
	}

	return Snapshot{collection: collection, result: result.Clone()}, nil
}

// Collection returns the structural collection that produced this snapshot.
func (s Snapshot) Collection() objectstore.ListRequest {
	return s.collection
}

// Result returns a detached copy of the validated list result.
func (s Snapshot) Result() objectstore.ListResult {
	return s.result.Clone()
}

// Items returns detached list items in snapshot order.
func (s Snapshot) Items() []objectstore.ListItem {
	return s.Result().Items
}

// Len returns the number of live items captured by the snapshot.
func (s Snapshot) Len() int {
	return s.result.Len()
}

// Revision returns the collection boundary revision captured by the snapshot.
func (s Snapshot) Revision() objectstore.Revision {
	return s.result.Revision
}

// Boundary returns the collection/revision pair represented by this snapshot.
//
// Invalid snapshots return the zero Boundary. Callers that operate on arbitrary
// Snapshot values should call Validate before using the returned Boundary.
func (s Snapshot) Boundary() Boundary {
	if !s.IsValid() {
		return Boundary{}
	}

	return Boundary{collection: s.collection, revision: s.result.Revision}
}

// Clone returns a detached copy of s.
func (s Snapshot) Clone() Snapshot {
	return Snapshot{collection: s.collection, result: s.result.Clone()}
}

// IsZero reports whether s contains no collection and no list result.
func (s Snapshot) IsZero() bool {
	return s.collection == (objectstore.ListRequest{}) && s.result.IsZero()
}

// IsValid reports whether s passes snapshot validation.
func (s Snapshot) IsValid() bool {
	return s.Validate() == nil
}
