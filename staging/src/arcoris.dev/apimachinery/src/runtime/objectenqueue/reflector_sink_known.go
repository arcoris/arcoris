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

package objectenqueue

import (
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// projectRead builds the next admitted known set from collection read items.
//
// Predicate.Match is the replace-path semantic source. The input items already
// came from CollectionRead.Items, so they are detached from the collection read
// and safe to retain if the replace operation succeeds.
func (s *ReflectorSink) projectRead(items []objectstore.ListItem) (
	map[objectstore.Key]objectstore.ListItem,
	[]objectstore.Key,
) {
	nextKnown := make(map[objectstore.Key]objectstore.ListItem, len(items))
	nextOrder := make([]objectstore.Key, 0, len(items))

	for _, item := range items {
		if !s.predicate.Match(item) {
			continue
		}
		nextKnown[item.Key] = item
		nextOrder = append(nextOrder, item.Key)
	}

	return nextKnown, nextOrder
}

// emitListedLocked emits all current admitted list items in collection order.
func (s *ReflectorSink) emitListedLocked(
	order []objectstore.Key,
	known map[objectstore.Key]objectstore.ListItem,
	emitter *handlerEmitter,
) error {
	for _, key := range order {
		item := known[key]
		mapErr := s.listed.MapListItem(item.Clone(), emitter.emit)
		if err := emitter.result(mapErr); err != nil {
			return err
		}
	}

	return nil
}

// emitMissingLocked emits previously admitted keys absent from the next set.
//
// Previous known order is preserved so replace repair work is deterministic and
// does not depend on map iteration order.
func (s *ReflectorSink) emitMissingLocked(
	nextKnown map[objectstore.Key]objectstore.ListItem,
	emitter *handlerEmitter,
) error {
	for _, key := range s.order {
		if _, stillKnown := nextKnown[key]; stillKnown {
			continue
		}

		item, wasKnown := s.known[key]
		if !wasKnown {
			continue
		}

		mapErr := s.listed.MapListItem(item.Clone(), emitter.emit)
		if err := emitter.result(mapErr); err != nil {
			return err
		}
	}

	return nil
}

// applyProjectedChangeLocked applies a successful change projection to known.
func (s *ReflectorSink) applyProjectedChangeLocked(kind objectquery.ChangeProjectionKind, change objectstore.Change) {
	switch kind {
	case objectquery.ChangeProjectionEntered, objectquery.ChangeProjectionUpdated:
		if _, exists := s.known[change.Key]; !exists {
			s.order = append(s.order, change.Key)
		}
		s.known[change.Key] = objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	case objectquery.ChangeProjectionLeft:
		delete(s.known, change.Key)
		s.removeKnownOrderLocked(change.Key)
	}
}

// removeKnownOrderLocked removes key from order while preserving relative order.
func (s *ReflectorSink) removeKnownOrderLocked(key objectstore.Key) {
	for i, existing := range s.order {
		if !existing.Equal(key) {
			continue
		}

		copy(s.order[i:], s.order[i+1:])
		s.order = s.order[:len(s.order)-1]
		return
	}
}
