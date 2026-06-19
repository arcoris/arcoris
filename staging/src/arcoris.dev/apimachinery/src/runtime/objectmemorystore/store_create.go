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

// Create commits state for a missing or tombstoned key.
//
// The caller-supplied state must be an input state with zero revision; the
// store assigns the committed revision. Recreating a tombstoned key publishes a
// new live record in the existing per-object slot so stale delete/update races
// cannot succeed through physical map removal.
func (s *Store) Create(ctx context.Context, key objectstore.Key, state objectstore.State) (objectstore.State, error) {
	prepared, err := s.prepareCreate(ctx, key, state)
	if err != nil {
		return objectstore.State{}, err
	}

	st := s.shardFor(key).getOrCreate(key)
	for {
		current := st.load()
		if current != nil && !current.deleted {
			return objectstore.State{}, objectstoreError(
				objectstore.ErrorReasonAlreadyExists,
				key,
				0,
				current.state.Revision,
				objectstore.ErrAlreadyExists,
			)
		}

		next := liveRecord(prepared, s.nextRevision())
		if st.compareAndSwap(current, next) {
			return next.visibleState(), nil
		}
	}
}
