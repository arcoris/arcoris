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
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// Get delegates to the backend under the observable store lock.
//
// Locking reads is conservative but keeps the first implementation simple:
// callers observe backend state through the same serialization point used for
// writes, collection reads, and watch registration.
func (s *Store) Get(ctx context.Context, key objectstore.Key) (objectstore.State, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.backend.Get(ctx, key)
}

// List delegates raw structural objectstore listing under the observable lock.
//
// List is the plain objectstore.Store operation. It does not create an
// api/objectstorewatch CollectionRead; callers that need a watch boundary
// should use ListCollection instead.
func (s *Store) List(ctx context.Context, request objectstore.ListRequest) (objectstore.ListResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.backend.List(ctx, request)
}

// ListCollection returns a validated boundary-safe collection read.
//
// The backend List call and CollectionRead validation run under Store.mu. That
// lock coverage serializes the collection boundary with writes and later watch
// registration, which is the core continuity guarantee this wrapper can provide
// over a generic objectstore.Store.
func (s *Store) ListCollection(
	ctx context.Context,
	collection objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := objectstore.ValidateListRequest(collection); err != nil {
		return storewatchapi.CollectionRead{}, err
	}
	result, err := s.backend.List(ctx, collection)
	if err != nil {
		return storewatchapi.CollectionRead{}, err
	}

	return storewatchapi.NewCollectionRead(collection, result)
}
