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

// prepareGet validates the common boundary for a read operation.
func (s *Store) prepareGet(ctx context.Context, key objectstore.Key) error {
	return s.prepareKeyed(ctx, key)
}

// prepareCreate validates Create inputs and returns detached state for storage.
func (s *Store) prepareCreate(ctx context.Context, key objectstore.Key, state objectstore.State) (objectstore.State, error) {
	if err := s.prepareKeyed(ctx, key); err != nil {
		return objectstore.State{}, err
	}

	return prepareInputState(state)
}

// prepareUpdate validates Update inputs and returns detached state for storage.
func (s *Store) prepareUpdate(
	ctx context.Context,
	key objectstore.Key,
	expected objectstore.Revision,
	state objectstore.State,
) (objectstore.State, error) {
	if err := s.prepareRevisioned(ctx, key, expected); err != nil {
		return objectstore.State{}, err
	}

	return prepareInputState(state)
}

// prepareDelete validates Delete inputs before the tombstone CAS loop starts.
func (s *Store) prepareDelete(ctx context.Context, key objectstore.Key, expected objectstore.Revision) error {
	return s.prepareRevisioned(ctx, key, expected)
}

// prepareRevisioned validates the common boundary for expected-revision writes.
func (s *Store) prepareRevisioned(ctx context.Context, key objectstore.Key, expected objectstore.Revision) error {
	if err := s.prepareKeyed(ctx, key); err != nil {
		return err
	}

	return validateExpectedRevision(key, expected)
}

// prepareKeyed keeps receiver, context, and key validation order consistent.
func (s *Store) prepareKeyed(ctx context.Context, key objectstore.Key) error {
	if err := requireStore(s); err != nil {
		return err
	}
	if err := checkContext(ctx); err != nil {
		return err
	}

	return validateKey(key)
}
