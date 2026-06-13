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

package objectstore

import "context"

// Store is the authoritative committed object state contract.
//
// Store implementations commit already-computed live object state and ownership
// documents. They do not apply requests, run admission, validate against
// resource descriptors, or stamp object metadata. Implementations should be
// safe for concurrent use unless documented otherwise.
type Store interface {
	// Get reads the latest live state for key.
	//
	// A missing or deleted object returns ok=false and a nil error. Invalid
	// keys and context cancellation are reported as errors.
	Get(ctx context.Context, key Key) (state State, ok bool, err error)

	// Create commits state for a missing or deleted key and assigns a revision.
	//
	// The input state's Revision must be zero so callers cannot forge committed
	// revisions. The returned state contains the committed revision.
	Create(ctx context.Context, key Key, state State) (State, error)

	// Update replaces the live state for key if expected matches the current
	// committed revision.
	//
	// The input state's Revision must be zero. The returned state contains the
	// newly assigned committed revision.
	Update(ctx context.Context, key Key, expected Revision, state State) (State, error)

	// Delete tombstones the live state for key when expected matches the
	// current committed revision.
	//
	// DeleteResult.Deleted is the live state that was deleted and keeps its
	// previous live revision. DeleteResult.Revision is the newly assigned
	// tombstone commit revision.
	Delete(ctx context.Context, key Key, expected Revision) (DeleteResult, error)

	// List reads live committed states for one resource collection and scope.
	//
	// List returns only live states. Missing, deleted, and tombstoned objects
	// are omitted. It is a storage collection read: it does not validate
	// resource descriptors, authorize callers, run admission, apply selectors,
	// paginate results, watch changes, or interpret API-server request policy.
	List(ctx context.Context, request ListRequest) (ListResult, error)
}
