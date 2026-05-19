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

// nilComponentRegistryPanic is the stable panic string for nil
// *ComponentRegistry method receivers.
//
// A nil component registry is a wiring/configuration bug. Missing components
// are represented by Lookup returning false; a nil receiver is intentionally a
// panic so the failure is caught close to initialization.
const nilComponentRegistryPanic = "admission.ComponentRegistry: nil registry"

// ComponentRegistry is an owner-created catalog of stable component
// descriptors.
//
// The registry validates each component descriptor against a KindRegistry so a
// component can reference only a known catalog kind. It is not a runtime
// instance registry, not a global registry, and not a service-discovery
// mechanism. All methods are safe for concurrent use, and reads return copies.
type ComponentRegistry struct {
	// mu protects byID. Register takes the write lock; all read methods take the
	// read lock. KindRegistry owns its own lock, so component and kind catalog
	// reads remain independently safe.
	mu sync.RWMutex

	// kinds is the owner-provided kind catalog used for registry-level
	// validation. It is intentionally a reference to the owner's registry rather
	// than a copied snapshot, allowing owners to decide whether kind catalogs are
	// immutable after construction or extended before component registration.
	kinds *KindRegistry

	// byID stores descriptors by stable component ID.
	//
	// Descriptors are values rather than pointers, so read methods can return
	// copies and preserve internal registry ownership.
	byID map[ComponentID]ComponentDescriptor
}

// NewComponentRegistry creates a component catalog backed by kinds.
//
// Construction fails when kinds is nil, when any descriptor is syntactically
// invalid, when a descriptor references an unknown kind, or when a duplicate
// component ID is provided.
func NewComponentRegistry(
	kinds *KindRegistry,
	descriptors ...ComponentDescriptor,
) (*ComponentRegistry, error) {
	if kinds == nil {
		return nil, ErrNilKindRegistry
	}

	registry := &ComponentRegistry{
		kinds: kinds,
		byID:  make(map[ComponentID]ComponentDescriptor, len(descriptors)),
	}
	for _, descriptor := range descriptors {
		if err := registry.Register(descriptor); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

// MustComponentRegistry creates a component registry or panics on invalid
// catalog input.
//
// It is intended for static catalog assembly and tests where invalid descriptor
// literals are programming errors.
func MustComponentRegistry(
	kinds *KindRegistry,
	descriptors ...ComponentDescriptor,
) *ComponentRegistry {
	registry, err := NewComponentRegistry(kinds, descriptors...)
	if err != nil {
		panic(err)
	}
	return registry
}

// Register adds descriptor to r after syntax and catalog validation.
//
// ComponentDescriptor.IsValid performs only local syntax checks. Register adds
// catalog-level checks: the kind must be present in r's KindRegistry and the ID
// must not already be registered.
func (r *ComponentRegistry) Register(descriptor ComponentDescriptor) error {
	r.requireNonNil()
	if !descriptor.IsValid() {
		return InvalidComponentDescriptorError{Descriptor: descriptor}
	}
	// This guard primarily protects manually allocated zero-value registries.
	// NewComponentRegistry rejects nil kind catalogs before returning.
	if r.kinds == nil {
		return ErrNilKindRegistry
	}
	if !r.kinds.Contains(descriptor.Kind) {
		return UnknownComponentKindError{Kind: descriptor.Kind}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// A registry normally initializes byID in NewComponentRegistry. The lazy path
	// keeps a manually allocated zero-value ComponentRegistry from panicking on
	// map assignment after it has passed the explicit kind-registry check.
	if r.byID == nil {
		r.byID = make(map[ComponentID]ComponentDescriptor, 1)
	}
	if _, exists := r.byID[descriptor.ID]; exists {
		return DuplicateComponentError{ID: descriptor.ID}
	}
	r.byID[descriptor.ID] = descriptor

	return nil
}

// Lookup returns the descriptor registered for id.
//
// Invalid IDs and missing valid IDs both return the zero descriptor and false.
// Invalid lookup keys are treated as absence because lookup is a read operation,
// not a validation report.
func (r *ComponentRegistry) Lookup(id ComponentID) (ComponentDescriptor, bool) {
	r.requireNonNil()
	if !id.IsValid() {
		return ComponentDescriptor{}, false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// A nil map behaves like an empty registry. Lookup still returns a copy-safe
	// zero descriptor and false for the absence case.
	descriptor, ok := r.byID[id]
	return descriptor, ok
}

// Contains reports whether id is registered.
func (r *ComponentRegistry) Contains(id ComponentID) bool {
	_, ok := r.Lookup(id)
	return ok
}

// List returns registered descriptors sorted by ComponentID string.
//
// The returned slice is a fresh copy. Sorting makes catalog output stable for
// tests, docs generation, and operator-facing config validation.
func (r *ComponentRegistry) List() []ComponentDescriptor {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Build and sort a fresh slice so callers can freely mutate the returned
	// descriptors without affecting registry state or future deterministic list
	// ordering.
	descriptors := make([]ComponentDescriptor, 0, len(r.byID))
	for _, descriptor := range r.byID {
		descriptors = append(descriptors, descriptor)
	}
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].ID.String() < descriptors[j].ID.String()
	})

	return descriptors
}

// Len reports the number of registered components.
func (r *ComponentRegistry) Len() int {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	// len(nil) is zero, so Len is safe for a manually allocated zero-value
	// registry as long as the receiver itself is non-nil.
	return len(r.byID)
}

// requireNonNil panics with the stable nil receiver message used by all
// ComponentRegistry methods.
func (r *ComponentRegistry) requireNonNil() {
	if r == nil {
		panic(nilComponentRegistryPanic)
	}
}
