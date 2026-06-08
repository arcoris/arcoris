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

package memory

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Update replaces live state when expected matches the current revision.
func (s *Store) Update(ctx context.Context, key objectstore.Key, expected objectstore.Revision, state objectstore.State) (objectstore.State, error) {
	if err := requireStore(s); err != nil {
		return objectstore.State{}, err
	}
	if err := checkContext(ctx); err != nil {
		return objectstore.State{}, err
	}
	if err := validateKey(key); err != nil {
		return objectstore.State{}, err
	}
	if err := validateExpectedRevision(key, expected); err != nil {
		return objectstore.State{}, err
	}
	prepared, err := prepareInputState(state)
	if err != nil {
		return objectstore.State{}, err
	}

	slot := s.shardFor(key).get(key)
	if slot == nil {
		return objectstore.State{}, objectstoreError(
			objectstore.ReasonNotFound,
			key,
			expected,
			0,
			objectstore.ErrNotFound,
		)
	}

	for {
		current := slot.load()
		if current == nil || current.deleted {
			return objectstore.State{}, objectstoreError(
				objectstore.ReasonNotFound,
				key,
				expected,
				0,
				objectstore.ErrNotFound,
			)
		}
		if current.state.Revision != expected {
			return objectstore.State{}, objectstoreError(
				objectstore.ReasonStaleRevision,
				key,
				expected,
				current.state.Revision,
				objectstore.ErrStaleRevision,
			)
		}

		next := liveRecord(prepared, s.nextRevision())
		if slot.compareAndSwap(current, next) {
			return next.visibleState(), nil
		}

		if latest := slot.load(); latest != nil && latest.state.Revision != expected {
			return objectstore.State{}, objectstoreError(
				objectstore.ReasonConflict,
				key,
				expected,
				latest.state.Revision,
				objectstore.ErrConflict,
			)
		}
	}
}
