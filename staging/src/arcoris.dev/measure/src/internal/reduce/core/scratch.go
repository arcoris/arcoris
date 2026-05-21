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


package core

// Scratch owns reusable planning, partial-result, and active-slot buffers for
// one partial result type.
//
// Scratch is for sequential reuse by one caller. Do not share one Scratch
// between simultaneous reductions because planners and runners mutate its slices
// and may treat reused partial slots as dirty. Reset keeps backing arrays and
// can retain references stored in those arrays; pairwise merge paths may also
// leave intermediate reference-bearing partials in retained slots until Clear
// or Release runs. Clear zeroes retained slots while keeping capacity; Release
// drops backing storage.
type Scratch[T any] struct {
	// Ranges stores the most recent planned ranges. Planner packages may reuse
	// and overwrite this slice on each call.
	Ranges []Range

	// Partials stores worker-local or range-local partial results. Runner
	// packages may reuse and overwrite this slice on each call.
	Partials []T

	// Used stores runner-owned active-slot metadata. Worker-local execution uses
	// it to mark which partial slots contain meaningful values, because zero-value
	// partials are not assumed to be merge identities.
	Used []bool
}

// Reset clears logical slice lengths while retaining backing storage.
//
// Reset is the cheapest reuse operation. It does not zero backing arrays, so it
// may retain references in old range or partial slots until those slots are
// overwritten, cleared, or released.
func (s *Scratch[T]) Reset() {
	s.Ranges = s.Ranges[:0]
	s.Partials = s.Partials[:0]
	s.Used = s.Used[:0]
}

// Clear zeroes retained backing slots while keeping backing storage.
//
// Clear scans the full retained capacity of Ranges, Partials, and Used, then
// resets their lengths to zero. It is more expensive than Reset, but it removes
// references, pairwise intermediate values, and stale active flags even after a
// previous Reset shortened the slices.
func (s *Scratch[T]) Clear() {
	var zeroRange Range
	ranges := s.Ranges[:cap(s.Ranges)]
	for i := range ranges {
		ranges[i] = zeroRange
	}
	var zeroPartial T
	partials := s.Partials[:cap(s.Partials)]
	for i := range partials {
		partials[i] = zeroPartial
	}
	used := s.Used[:cap(s.Used)]
	for i := range used {
		used[i] = false
	}
	s.Ranges = ranges[:0]
	s.Partials = partials[:0]
	s.Used = used[:0]
}

// Release drops all scratch backing storage.
//
// Release is the strongest memory-management operation: all retained slices
// become nil, allowing their backing arrays and any references stored in them to
// be garbage collected when no other references exist.
func (s *Scratch[T]) Release() {
	s.Ranges = nil
	s.Partials = nil
	s.Used = nil
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

// EnsureUsed returns an inactive flag slice of length n.
func (s *Scratch[T]) EnsureUsed(n int) []bool {
	if n <= 0 {
		s.Used = s.Used[:0]
		return s.Used
	}
	if cap(s.Used) < n {
		s.Used = make([]bool, n)
		return s.Used
	}
	s.Used = s.Used[:n]
	for i := range s.Used {
		s.Used[i] = false
	}
	return s.Used
}

// EnsureUsedDirty returns an active-slot flag slice of length n without
// clearing reused slots.
//
// Runners use this only when they will set every used flag before reading it.
func (s *Scratch[T]) EnsureUsedDirty(n int) []bool {
	if n <= 0 {
		s.Used = s.Used[:0]
		return s.Used
	}
	if cap(s.Used) < n {
		s.Used = make([]bool, n)
		return s.Used
	}
	s.Used = s.Used[:n]
	return s.Used
}
