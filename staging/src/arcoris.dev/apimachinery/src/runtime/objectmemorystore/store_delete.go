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
func (s *Store) Delete(
	ctx context.Context,
	key objectstore.Key,
	expected objectstore.Revision,
) (objectstore.DeleteResult, error) {
	if err := s.prepareDelete(ctx, key, expected); err != nil {
		return objectstore.DeleteResult{}, err
	}

	st := s.shardFor(key).get(key)
	if st == nil {
		return objectstore.DeleteResult{}, objectstoreError(
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
			return objectstore.DeleteResult{}, objectstoreError(
				objectstore.ErrorReasonNotFound,
				key,
				expected,
				0,
				objectstore.ErrNotFound,
			)
		}
		if current.state.Revision != expected {
			return objectstore.DeleteResult{}, objectstoreError(
				objectstore.ErrorReasonStaleRevision,
				key,
				expected,
				current.state.Revision,
				objectstore.ErrStaleRevision,
			)
		}

		deleted := current.visibleState()
		deleteRevision := s.nextRevision()
		next := tombstoneRecord(deleted, deleteRevision)
		if st.compareAndSwap(current, next) {
			return objectstore.DeleteResult{
				Deleted:  deleted,
				Revision: deleteRevision,
			}, nil
		}

		if latest := st.load(); latest != nil && latest.state.Revision != expected {
			return objectstore.DeleteResult{}, objectstoreError(
				objectstore.ErrorReasonConflict,
				key,
				expected,
				latest.state.Revision,
				objectstore.ErrConflict,
			)
		}
	}
}
