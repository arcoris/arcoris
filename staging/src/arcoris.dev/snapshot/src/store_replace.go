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
// cloning before commit panics or revision overflow is detected, the current
// Store value is left unchanged.
func (s *Store[T]) ReplaceStamped(next T) Stamped[T] {
	stored := s.clone(next)
	returned := s.clone(stored)

	s.mu.Lock()
	defer s.mu.Unlock()

	rev := s.revision.Next()
	updated := s.clock.Now()

	s.value = stored
	s.revision = rev
	s.updated = updated

	return Stamped[T]{
		Revision: rev,
		Updated:  updated,
		Value:    returned,
	}
}
