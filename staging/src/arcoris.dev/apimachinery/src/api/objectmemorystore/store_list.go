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

// List returns a detached current live collection read.
//
// List scans the sharded key index without a global list lock. It copies
// matching key/slot pairs under shard locks, releases those locks, then
// atomically loads immutable records. The result is deterministic and detached,
// but it is not a historical MVCC snapshot under concurrent writes.
func (s *Store) List(ctx context.Context, request objectstore.ListRequest) (objectstore.ListResult, error) {
	if err := s.prepareList(ctx, request); err != nil {
		return objectstore.ListResult{}, err
	}

	candidates := s.listCandidates(request)
	items := make([]objectstore.ListItem, 0, len(candidates))
	for _, candidate := range candidates {
		current := candidate.slot.load()
		if current == nil || current.deleted {
			continue
		}
		items = append(items, objectstore.ListItem{
			Key:   candidate.key,
			State: current.visibleState(),
		})
	}

	sortListItems(items)

	return objectstore.ListResult{
		Items:    items,
		Revision: objectstore.Revision(s.revision.Load()),
	}, nil
}

// listCandidates copies matching slots while holding only shard map locks.
func (s *Store) listCandidates(request objectstore.ListRequest) []listCandidate {
	candidates := []listCandidate{}
	for i := range s.shards {
		shard := &s.shards[i]
		shard.mu.Lock()
		for key, st := range shard.slots {
			if listMatches(request, key) {
				candidates = append(candidates, listCandidate{key: key, slot: st})
			}
		}
		shard.mu.Unlock()
	}

	return candidates
}

// listCandidate is a stable key/slot pair copied from a shard index.
type listCandidate struct {
	// key is the validated storage identity.
	key objectstore.Key

	// slot is the per-object publication slot for key.
	slot *slot
}
