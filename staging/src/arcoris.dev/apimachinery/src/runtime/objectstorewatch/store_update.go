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
	"errors"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Update delegates to the backend and publishes the committed update change.
//
// The wrapper reads the committed before state under Store.mu before delegating
// Update so the resulting objectstore.Change contains both sides of the
// transition. If the backend succeeds without a readable before state, the
// wrapper treats that as continuity loss rather than inventing a partial event.
func (s *Store) Update(
	ctx context.Context,
	key objectstore.Key,
	expected objectstore.Revision,
	state objectstore.State,
) (objectstore.State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	before, ok, err := s.backend.Get(ctx, key)
	if err != nil {
		return objectstore.State{}, err
	}

	after, err := s.backend.Update(ctx, key, expected, state)
	if err != nil {
		return objectstore.State{}, err
	}
	if !ok {
		cause := errors.New("backend update succeeded without readable before state")
		return after, s.loseCommittedContinuityLocked(after.Revision, cause)
	}

	change, err := updatedChange(key, before, after)
	if err != nil {
		return after, s.loseCommittedContinuityLocked(after.Revision, err)
	}

	if err := s.publishLocked(change); err != nil {
		return after, err
	}

	return after, nil
}
