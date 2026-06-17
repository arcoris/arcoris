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
