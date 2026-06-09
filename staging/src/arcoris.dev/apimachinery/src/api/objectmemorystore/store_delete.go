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

package objectmemorystore

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Delete tombstones live state when expected matches the current revision.
//
// Delete returns the live State that was removed. The tombstone commit revision
// is kept internally so a future event layer can distinguish the delete commit,
// but the current objectstore.Store contract does not expose that revision.
func (s *Store) Delete(ctx context.Context, key objectstore.Key, expected objectstore.Revision) (objectstore.State, error) {
	if err := s.prepareDelete(ctx, key, expected); err != nil {
		return objectstore.State{}, err
	}

	st := s.shardFor(key).get(key)
	if st == nil {
		return objectstore.State{}, objectstoreError(
			objectstore.ErrorReasonNotFound,
			key,
			expected,
			0,
			objectstore.ErrNotFound,
		)
	}

	for {
		current := st.load()
		if current == nil || current.deleted {
			return objectstore.State{}, objectstoreError(
				objectstore.ErrorReasonNotFound,
				key,
				expected,
				0,
				objectstore.ErrNotFound,
			)
		}
		if current.state.Revision != expected {
			return objectstore.State{}, objectstoreError(
				objectstore.ErrorReasonStaleRevision,
				key,
				expected,
				current.state.Revision,
				objectstore.ErrStaleRevision,
			)
		}

		deleted := current.visibleState()
		next := tombstoneRecord(deleted, s.nextRevision())
		if st.compareAndSwap(current, next) {
			return deleted, nil
		}

		if latest := st.load(); latest != nil && latest.state.Revision != expected {
			return objectstore.State{}, objectstoreError(
				objectstore.ErrorReasonConflict,
				key,
				expected,
				latest.state.Revision,
				objectstore.ErrConflict,
			)
		}
	}
}
