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
