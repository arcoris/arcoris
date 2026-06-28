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

package snapshot

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
func (s *Store[T]) Update(update func(T) T) Snapshot[LocalRevision, T] {
	return s.UpdateStamped(update).Snapshot()
}

// UpdateStamped applies update to a cloned copy of the current value and returns
// the resulting stamped snapshot.
//
// UpdateStamped has the same ownership and locking semantics as Update and also
// records the local commit time using the configured PassiveClock. If update
// panics, cloning before commit panics, or revision overflow is detected, the
// current Store value is left unchanged.
func (s *Store[T]) UpdateStamped(update func(T) T) Stamped[LocalRevision, T] {
	if update == nil {
		panic("snapshot: nil update function")
	}

	stamped, err := s.UpdateStampedErr(func(current T) (T, error) {
		return update(current), nil
	})
	if err != nil {
		panic("snapshot: unexpected update error")
	}

	return stamped
}

// UpdateErr applies update to a cloned copy of the current value and returns
// the resulting lightweight snapshot.
//
// UpdateErr has the same locking and ownership semantics as Update. If update
// returns an error, Store leaves the previous committed value unchanged and
// returns a zero snapshot with that error.
func (s *Store[T]) UpdateErr(update func(T) (T, error)) (Snapshot[LocalRevision, T], error) {
	stamped, err := s.UpdateStampedErr(update)
	if err != nil {
		return Snapshot[LocalRevision, T]{}, err
	}

	return stamped.Snapshot(), nil
}

// UpdateStampedErr applies update to a cloned copy of the current value and
// returns the resulting stamped snapshot.
//
// UpdateStampedErr is the fallible form of UpdateStamped. If update returns an
// error, Store does not advance its revision, does not change its value, and
// returns a zero stamped snapshot with that error.
func (s *Store[T]) UpdateStampedErr(update func(T) (T, error)) (Stamped[LocalRevision, T], error) {
	if update == nil {
		panic("snapshot: nil update function")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	working := s.clone(s.value)
	next, err := update(working)
	if err != nil {
		return Stamped[LocalRevision, T]{}, err
	}

	stored := s.clone(next)
	returned := s.clone(stored)
	rev := s.revision.Next()
	updated := s.clock.Now()

	s.value = stored
	s.revision = rev
	s.updated = updated

	return Stamped[LocalRevision, T]{
		Revision: rev,
		Updated:  updated,
		Value:    returned,
	}, nil
}
