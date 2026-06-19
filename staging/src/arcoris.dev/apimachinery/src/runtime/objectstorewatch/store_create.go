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

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Create delegates to the backend and publishes the committed create change.
//
// No event is published if the backend rejects the operation. If the backend
// reports success but returns a state that cannot form a valid created change,
// all live streams are terminated with continuity loss because the wrapper can
// no longer prove the stream is complete and well-formed.
func (s *Store) Create(ctx context.Context, key objectstore.Key, state objectstore.State) (objectstore.State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	created, err := s.backend.Create(ctx, key, state)
	if err != nil {
		return objectstore.State{}, err
	}

	change, err := createdChange(key, created)
	if err != nil {
		return created, s.loseCommittedContinuityLocked(created.Revision, err)
	}
	if err := s.publishLocked(change); err != nil {
		return created, err
	}

	return created, nil
}
