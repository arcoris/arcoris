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

// Update replaces live state when expected matches the current revision.
//
// Update is an optimistic compare-and-swap transition. If another writer
// commits first, the loser reports stale revision, conflict, or not found
// according to the record shape it observes after the race.
func (s *Store) Update(ctx context.Context, key objectstore.Key, expected objectstore.Revision, state objectstore.State) (objectstore.State, error) {
	prepared, err := s.prepareUpdate(ctx, key, expected, state)
	if err != nil {
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

		next := liveRecord(prepared, s.nextRevision())
		if st.compareAndSwap(current, next) {
			return next.visibleState(), nil
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
