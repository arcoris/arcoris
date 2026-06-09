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

// Get reads the latest live state for key.
//
// Missing keys and tombstones are ordinary cache-style misses and return
// ok=false with nil error. The returned state is detached from the immutable
// internal record, so caller mutation cannot alter subsequent reads.
func (s *Store) Get(ctx context.Context, key objectstore.Key) (objectstore.State, bool, error) {
	if err := s.prepareGet(ctx, key); err != nil {
		return objectstore.State{}, false, err
	}

	st := s.shardFor(key).get(key)
	if st == nil {
		return objectstore.State{}, false, nil
	}

	current := st.load()
	if current == nil || current.deleted {
		return objectstore.State{}, false, nil
	}

	return current.visibleState(), true, nil
}
