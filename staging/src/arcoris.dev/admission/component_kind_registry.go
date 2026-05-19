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

package admission

import (
	"sort"
	"sync"
)

const nilKindRegistryPanic = "admission.KindRegistry: nil registry"

// KindRegistry is an owner-created catalog of known ComponentKind descriptors.
//
// The registry is intentionally not global. Callers create a registry for the
// catalog they own, pass it to validation or documentation code, and decide when
// to share or discard it. All methods are safe for concurrent use. Returned
// values are copies, so callers cannot mutate internal registry storage through
// Lookup or List.
type KindRegistry struct {
	mu     sync.RWMutex
	byKind map[ComponentKind]ComponentKindDescriptor
}

// NewKindRegistry creates a registry populated with descriptors.
//
// Descriptors are registered in argument order. Invalid descriptors and
// duplicate kinds fail construction, and no partially initialized registry is
// returned on error.
func NewKindRegistry(descriptors ...ComponentKindDescriptor) (*KindRegistry, error) {
	registry := &KindRegistry{
		byKind: make(map[ComponentKind]ComponentKindDescriptor, len(descriptors)),
	}
	for _, descriptor := range descriptors {
		if err := registry.Register(descriptor); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

// MustKindRegistry creates a registry or panics when descriptors are invalid.
//
// It is intended for package-level catalog assembly and tests where invalid
// descriptor literals are programming errors.
func MustKindRegistry(descriptors ...ComponentKindDescriptor) *KindRegistry {
	registry, err := NewKindRegistry(descriptors...)
	if err != nil {
		panic(err)
	}
	return registry
}

// Register adds descriptor to r.
//
// Register validates the descriptor and rejects duplicate kinds. It preserves
// the open-world nature of ComponentKind: any syntactically valid kind may be
// registered by the registry owner.
func (r *KindRegistry) Register(descriptor ComponentKindDescriptor) error {
	r.requireNonNil()
	if !descriptor.IsValid() {
		return InvalidComponentKindDescriptorError{Descriptor: descriptor}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.byKind == nil {
		r.byKind = make(map[ComponentKind]ComponentKindDescriptor, 1)
	}
	if _, exists := r.byKind[descriptor.Kind]; exists {
		return DuplicateComponentKindError{Kind: descriptor.Kind}
	}
	r.byKind[descriptor.Kind] = descriptor

	return nil
}

// Lookup returns the descriptor registered for kind.
//
// Invalid kinds and missing valid kinds both return the zero descriptor and
// false. Invalid lookup keys are treated as absence because lookup is a read
// operation, not a validation report.
func (r *KindRegistry) Lookup(kind ComponentKind) (ComponentKindDescriptor, bool) {
	r.requireNonNil()
	if !kind.IsValid() {
		return ComponentKindDescriptor{}, false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptor, ok := r.byKind[kind]
	return descriptor, ok
}

// Contains reports whether kind is registered.
func (r *KindRegistry) Contains(kind ComponentKind) bool {
	_, ok := r.Lookup(kind)
	return ok
}

// List returns registered descriptors sorted by ComponentKind string.
//
// The returned slice is a fresh copy. Mutating it cannot affect later registry
// reads. Sorting makes docs, tests, and config validation deterministic.
func (r *KindRegistry) List() []ComponentKindDescriptor {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	descriptors := make([]ComponentKindDescriptor, 0, len(r.byKind))
	for _, descriptor := range r.byKind {
		descriptors = append(descriptors, descriptor)
	}
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Kind.String() < descriptors[j].Kind.String()
	})

	return descriptors
}

// Len reports the number of registered kinds.
func (r *KindRegistry) Len() int {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.byKind)
}

// requireNonNil panics with the stable nil receiver message used by all
// KindRegistry methods.
func (r *KindRegistry) requireNonNil() {
	if r == nil {
		panic(nilKindRegistryPanic)
	}
}
