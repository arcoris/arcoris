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
