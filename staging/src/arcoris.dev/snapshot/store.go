/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

import (
	"sync"
	"time"

	"arcoris.dev/chrono/clock"
)

// Store is a concurrency-safe holder for one always-present mutable value.
//
// Store owns its internal value and uses CloneFunc to isolate writes and reads.
// It is the safe baseline for state that may contain slices, maps, pointers, or
// other mutable data. Snapshot and Stamped return cloned values; mutating a value
// returned from Store must not affect Store's internal state when the CloneFunc is
// correct.
//
// Store starts at revision 1 because NewStore commits the initial value. Use a
// value-level container such as maybe.Maybe[T] when the logical state can be
// absent.
//
// Store is safe for concurrent use. Store must not be copied after first use.
type Store[T any] struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects value, revision, and updated.
	mu sync.RWMutex

	// clone copies values across the Store ownership boundary.
	clone CloneFunc[T]

	// clock provides local commit timestamps for Stamped values.
	clock clock.PassiveClock

	// value is the currently committed internal value owned by the Store.
	value T

	// revision is the source-local revision of value.
	revision Revision

	// updated is the local time at which value was committed.
	updated time.Time
}

// NewStore creates a Store containing initial.
//
// NewStore clones initial before storing it and commits that initial value at
// revision 1. The clone function must be non-nil and must provide the ownership
// isolation required by T.
func NewStore[T any](initial T, clone CloneFunc[T], opts ...Option) *Store[T] {
	clone = requireClone(clone)
	cfg := newConfig(opts...)

	return &Store[T]{
		clone:    clone,
		clock:    cfg.clock,
		value:    clone(initial),
		revision: ZeroRevision.Next(),
		updated:  cfg.clock.Now(),
	}
}

// Snapshot returns the Store's current lightweight snapshot.
//
// The returned value is cloned from the Store's internal value. Mutating the
// returned value must not affect Store when the Store's CloneFunc is correct.
func (s *Store[T]) Snapshot() Snapshot[T] {
	stamped := s.Stamped()
	return stamped.Snapshot()
}

// Stamped returns the Store's current stamped snapshot.
//
// Stamped returns a cloned value and the local time at which the current value
// was committed.
func (s *Store[T]) Stamped() Stamped[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return Stamped[T]{
		Revision: s.revision,
		Updated:  s.updated,
		Value:    s.clone(s.value),
	}
}

// Revision returns the Store's current source-local revision.
func (s *Store[T]) Revision() Revision {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.revision
}

// Replace replaces the Store value and returns the resulting lightweight
// snapshot.
//
// Replace clones next before storing it and advances the revision exactly once.
func (s *Store[T]) Replace(next T) Snapshot[T] {
	return s.ReplaceStamped(next).Snapshot()
}

// ReplaceStamped replaces the Store value and returns the resulting stamped
// snapshot.
//
// ReplaceStamped clones next before storing it, advances the revision exactly
// once, and records the local commit time using the configured PassiveClock. If
// revision overflow is detected, ReplaceStamped panics before exposing next.
func (s *Store[T]) ReplaceStamped(next T) Stamped[T] {
	stored := s.clone(next)

	s.mu.Lock()
	defer s.mu.Unlock()

	revision := s.revision.Next()
	updated := s.clock.Now()

	s.value = stored
	s.revision = revision
	s.updated = updated

	return Stamped[T]{
		Revision: revision,
		Updated:  updated,
		Value:    s.clone(s.value),
	}
}

// Update applies update to a cloned copy of the current value and returns the
// resulting lightweight snapshot.
//
// Update is an atomic read-modify-write operation with respect to other Store
// operations. The update function runs while Store's write lock is held. It must
// not call back into the same Store, block indefinitely, or perform unrelated
// long-running work.
//
// The value passed to update is a clone of the current internal value. Store
// clones the value returned by update before committing it, so retaining or
// mutating values outside Store does not affect Store when the CloneFunc is
// correct.
func (s *Store[T]) Update(update func(T) T) Snapshot[T] {
	return s.UpdateStamped(update).Snapshot()
}

// UpdateStamped applies update to a cloned copy of the current value and returns
// the resulting stamped snapshot.
//
// UpdateStamped has the same ownership and locking semantics as Update and also
// records the local commit time using the configured PassiveClock. If update
// panics or revision overflow is detected, the current Store value is left
// unchanged.
func (s *Store[T]) UpdateStamped(update func(T) T) Stamped[T] {
	if update == nil {
		panic("snapshot: nil update function")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	working := s.clone(s.value)
	next := update(working)
	stored := s.clone(next)
	revision := s.revision.Next()
	updated := s.clock.Now()

	s.value = stored
	s.revision = revision
	s.updated = updated

	return Stamped[T]{
		Revision: revision,
		Updated:  updated,
		Value:    s.clone(s.value),
	}
}
