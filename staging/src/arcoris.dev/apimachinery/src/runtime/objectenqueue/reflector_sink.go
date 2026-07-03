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
	"context"
	"sync"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

// ReflectorSink maps reflected collection state and changes to queue work.
//
// ReflectorSink is replace-aware. It keeps a small local known set so a relist
// can enqueue previously admitted objects that disappear or stop matching the
// configured predicate without a corresponding watch change. The known set is
// only a loss-prevention boundary for enqueueing; it is not a read model.
type ReflectorSink struct {
	// queue receives every mapped reconciliation item.
	queue Enqueuer

	// predicate selects list items and projects committed changes.
	predicate objectquery.Predicate

	// listed maps live list items observed during Replace.
	listed ListItemMapper

	// changed maps committed changes observed during ApplyChange.
	changed Mapper

	// mu protects known and order for accidental concurrent public calls.
	mu sync.Mutex

	// known stores admitted live items by authoritative object key.
	known map[objectstore.Key]objectstore.ListItem

	// order stores admitted keys in their latest observed list order.
	order []objectstore.Key
}

// NewReflectorSink constructs a stateful reflected-state enqueue sink.
//
// The constructor validates only local wiring. It initializes empty known state
// and does not touch queue state or start background work.
func NewReflectorSink(config ReflectorSinkConfig) (*ReflectorSink, error) {
	if isNilInterface(config.Queue) {
		return nil, ErrNilQueue
	}
	if isNilInterface(config.Listed) {
		return nil, ErrNilListItemMapper
	}
	if isNilInterface(config.Changed) {
		return nil, ErrNilMapper
	}

	return &ReflectorSink{
		queue:     config.Queue,
		predicate: config.Predicate,
		listed:    config.Listed,
		changed:   config.Changed,
		known:     make(map[objectstore.Key]objectstore.ListItem),
		order:     make([]objectstore.Key, 0),
	}, nil
}

// Replace enqueues work for a complete reflected collection read.
//
// Current matching list items are emitted first in collection-read order.
// Previously known keys that are missing from the new admitted set are emitted
// afterward in previous known order. Internal known state changes only after
// every enqueue operation succeeds.
func (s *ReflectorSink) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	if err := s.validate(); err != nil {
		return err
	}
	if err := s.validateKnown(); err != nil {
		return err
	}
	if err := read.Validate(); err != nil {
		return err
	}

	nextKnown, nextOrder := s.projectRead(read.Items())

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.validateKnownLocked(); err != nil {
		return err
	}

	emitter := handlerEmitter{ctx: ctx, queue: s.queue}
	if err := s.emitListedLocked(nextOrder, nextKnown, &emitter); err != nil {
		return err
	}
	if err := s.emitMissingLocked(nextKnown, &emitter); err != nil {
		return err
	}

	s.known = nextKnown
	s.order = nextOrder

	return nil
}

// ApplyChange enqueues work for one committed reflected change.
//
// Predicate.ProjectChange owns change membership semantics. Ignored changes do
// not mutate known state. Entered, Updated, and Left changes enqueue through the
// committed-change mapper; known state changes only after enqueue succeeds.
func (s *ReflectorSink) ApplyChange(ctx context.Context, change objectstore.Change) error {
	if err := s.validate(); err != nil {
		return err
	}
	if err := s.validateKnown(); err != nil {
		return err
	}

	projection, err := s.predicate.ProjectChange(change)
	if err != nil {
		return err
	}
	if projection.Kind == objectquery.ChangeProjectionIgnored {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.validateKnownLocked(); err != nil {
		return err
	}

	emitter := handlerEmitter{ctx: ctx, queue: s.queue}
	mapErr := s.changed.Map(change, emitter.emit)
	if err := emitter.result(mapErr); err != nil {
		return err
	}

	s.applyProjectedChangeLocked(projection.Kind, change)

	return nil
}

// validate checks immutable receiver wiring before any operation work begins.
func (s *ReflectorSink) validate() error {
	if s == nil ||
		isNilInterface(s.queue) ||
		isNilInterface(s.listed) ||
		isNilInterface(s.changed) {
		return ErrInvalidReflectorSink
	}

	return nil
}

// validateKnown checks mutable receiver state under s.mu.
func (s *ReflectorSink) validateKnown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.validateKnownLocked()
}

// validateKnownLocked checks mutable receiver state while s.mu is held.
func (s *ReflectorSink) validateKnownLocked() error {
	if s.known == nil {
		return ErrInvalidReflectorSink
	}

	return nil
}
