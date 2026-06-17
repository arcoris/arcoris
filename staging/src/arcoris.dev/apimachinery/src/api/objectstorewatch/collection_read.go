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

// CollectionRead is a validated collection read plus its list-to-watch boundary.
//
// CollectionRead preserves the exact objectstore.ListRequest used for validation and
// stores a detached clone of the list result. It is immutable by API
// convention: methods that expose list data return detached copies.
type CollectionRead struct {
	// collection is the structural collection that produced result.
	collection objectstore.ListRequest
	// result is the detached validated list output for collection.
	result objectstore.ListResult
}

// NewCollectionRead validates collection/result and stores a detached collection read.
//
// Validation runs before cloning so malformed large results are rejected
// without first deep-copying their item payloads. Successful construction still
// detaches the stored result from caller-owned data.
func NewCollectionRead(collection objectstore.ListRequest, result objectstore.ListResult) (CollectionRead, error) {
	if err := objectstore.ValidateListRequest(collection); err != nil {
		return CollectionRead{}, errorFor("collection_read.collection", ErrorReasonInvalidCollectionRead, ErrInvalidCollectionRead, err)
	}
	if err := objectstore.ValidateListResult(collection, result); err != nil {
		return CollectionRead{}, errorFor("collection_read.result", ErrorReasonInvalidCollectionRead, ErrInvalidCollectionRead, err)
	}

	return CollectionRead{collection: collection, result: result.Clone()}, nil
}

// Collection returns the structural collection that produced this collection read.
func (s CollectionRead) Collection() objectstore.ListRequest {
	return s.collection
}

// Result returns a detached copy of the validated list result.
func (s CollectionRead) Result() objectstore.ListResult {
	return s.result.Clone()
}

// Items returns detached list items in collection read order.
func (s CollectionRead) Items() []objectstore.ListItem {
	return s.Result().Items
}

// Len returns the number of live items captured by the collection read.
func (s CollectionRead) Len() int {
	return s.result.Len()
}

// Revision returns the collection boundary revision captured by the collection read.
func (s CollectionRead) Revision() objectstore.Revision {
	return s.result.Revision
}

// Boundary returns the collection/revision pair represented by this collection read.
//
// Invalid collection reads return the zero Boundary. Callers that operate on arbitrary
// CollectionRead values should call Validate before using the returned Boundary.
func (s CollectionRead) Boundary() Boundary {
	if !s.IsValid() {
		return Boundary{}
	}

	return Boundary{collection: s.collection, revision: s.result.Revision}
}

// Clone returns a detached copy of s.
func (s CollectionRead) Clone() CollectionRead {
	return CollectionRead{collection: s.collection, result: s.result.Clone()}
}

// IsZero reports whether s contains no collection and no list result.
func (s CollectionRead) IsZero() bool {
	return s.collection == (objectstore.ListRequest{}) && s.result.IsZero()
}

// IsValid reports whether s passes collection read validation.
func (s CollectionRead) IsValid() bool {
	return s.Validate() == nil
}
