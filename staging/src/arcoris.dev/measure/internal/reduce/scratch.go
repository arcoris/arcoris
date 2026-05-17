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

package reduce

// Scratch owns reusable planning and partial-result buffers for one partial
// result type.
//
// Scratch is for sequential reuse by one caller. Do not share one Scratch
// between simultaneous reductions because runners mutate its slices and may
// treat reused partial slots as dirty.
type Scratch[T any] struct {
	// Ranges stores the most recent planned ranges. Planner packages may reuse
	// and overwrite this slice on each call.
	Ranges []Range

	// Partials stores worker-local or range-local partial results. Runner
	// packages may reuse and overwrite this slice on each call.
	Partials []T
}

// Reset clears logical scratch contents while retaining backing storage for
// future reductions.
func (s *Scratch[T]) Reset() {
	s.Ranges = s.Ranges[:0]
	s.Partials = s.Partials[:0]
}

// EnsurePartials returns a zeroed partial-result slice of length n.
//
// Use this when the caller needs old partial values removed before mapping or
// merging. Existing capacity is reused when possible.
func (s *Scratch[T]) EnsurePartials(n int) []T {
	if n <= 0 {
		s.Partials = s.Partials[:0]
		return s.Partials
	}
	if cap(s.Partials) < n {
		s.Partials = make([]T, n)
		return s.Partials
	}
	s.Partials = s.Partials[:n]
	var zero T
	for i := range s.Partials {
		s.Partials[i] = zero
	}
	return s.Partials
}

// EnsurePartialsDirty returns a partial-result slice of length n without
// clearing reused slots.
//
// Runners use this when they will overwrite every used slot before merging. It
// saves clearing cost for large pointer-free partial types and dense counters.
// Callers must not read returned slots before writing them.
func (s *Scratch[T]) EnsurePartialsDirty(n int) []T {
	if n <= 0 {
		s.Partials = s.Partials[:0]
		return s.Partials
	}
	if cap(s.Partials) < n {
		s.Partials = make([]T, n)
		return s.Partials
	}
	s.Partials = s.Partials[:n]
	return s.Partials
}
