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

package admissioncatalog

import "sort"

// descriptorStore owns the repeated immutable-map mechanics for descriptor
// categories.
//
// The type is deliberately private. Public APIs stay domain-specific, while the
// implementation shares copy, lookup, deterministic listing, and duplicate
// detection behavior.
type descriptorStore[K comparable, D any] struct {
	// byKey stores descriptor values, never pointers, so catalog reads can
	// return detached values without cloning descriptor internals.
	byKey map[K]D

	// key extracts the stable catalog key for duplicate detection and lookup.
	key func(D) K

	// less gives each descriptor family deterministic list ordering.
	less func(a, b D) bool
}

// newDescriptorStore returns an initialized store for one descriptor family.
func newDescriptorStore[K comparable, D any](
	key func(D) K,
	less func(a, b D) bool,
) descriptorStore[K, D] {
	return descriptorStore[K, D]{
		byKey: make(map[K]D),
		key:   key,
		less:  less,
	}
}

// init attaches descriptor-family behavior and allocates storage for a zero
// descriptor store.
func (s *descriptorStore[K, D]) init(
	key func(D) K,
	less func(a, b D) bool,
) {
	if s.key == nil {
		s.key = key
	}
	if s.less == nil {
		s.less = less
	}
	if s.byKey == nil {
		s.byKey = make(map[K]D)
	}
}

// len reports the number of descriptors in the store.
func (s descriptorStore[K, D]) len() int {
	return len(s.byKey)
}

// get returns the descriptor stored for key.
func (s descriptorStore[K, D]) get(key K) (D, bool) {
	descriptor, ok := s.byKey[key]
	return descriptor, ok
}

// has reports whether key exists in the store.
func (s descriptorStore[K, D]) has(key K) bool {
	_, ok := s.get(key)
	return ok
}

// list returns a sorted detached descriptor slice.
func (s descriptorStore[K, D]) list() []D {
	descriptors := make([]D, 0, len(s.byKey))
	for _, descriptor := range s.byKey {
		descriptors = append(descriptors, descriptor)
	}
	if s.less != nil {
		sort.Slice(descriptors, func(i, j int) bool {
			return s.less(descriptors[i], descriptors[j])
		})
	}
	return descriptors
}

// clone returns a detached store with the same descriptor-family behavior and
// descriptor values.
func (s descriptorStore[K, D]) clone() descriptorStore[K, D] {
	clone := descriptorStore[K, D]{
		byKey: make(map[K]D, len(s.byKey)),
		key:   s.key,
		less:  s.less,
	}
	for key, descriptor := range s.byKey {
		clone.byKey[key] = descriptor
	}
	return clone
}

// declare records descriptor when its key is not already present.
func (s *descriptorStore[K, D]) declare(descriptor D) bool {
	key := s.key(descriptor)
	if _, exists := s.byKey[key]; exists {
		return false
	}
	s.byKey[key] = descriptor
	return true
}
