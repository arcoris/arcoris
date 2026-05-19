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

// nilReasonRegistryPanic is the stable panic string for nil *ReasonRegistry
// method receivers.
//
// Registry methods panic on nil receivers because nil registry use is a
// programming/configuration error, not an ordinary catalog miss. The exact
// string is tested so downstream initialization failures stay predictable.
const nilReasonRegistryPanic = "admission.ReasonRegistry: nil registry"

// ReasonRegistry is an owner-created catalog of known Reason descriptors.
//
// The registry is intentionally not global. Callers create a registry for the
// reason catalog they own and pass it to docs, config validation, or future
// chain validation code. All methods are safe for concurrent use. Lookup and
// List return copies, so callers cannot mutate internal registry storage.
type ReasonRegistry struct {
	// mu protects byReason. Register takes the write lock; all read methods take
	// the read lock.
	mu sync.RWMutex

	// byReason stores descriptors by stable open-world reason value. Descriptors
	// are values rather than pointers, preserving copy-safe reads.
	byReason map[Reason]ReasonDescriptor
}

// NewReasonRegistry creates a registry populated with descriptors.
//
// Descriptors are registered in argument order. Invalid descriptors and
// duplicate reasons fail construction, and no partially initialized registry is
// returned on error.
func NewReasonRegistry(descriptors ...ReasonDescriptor) (*ReasonRegistry, error) {
	registry := &ReasonRegistry{
		byReason: make(map[Reason]ReasonDescriptor, len(descriptors)),
	}
	for _, descriptor := range descriptors {
		if err := registry.Register(descriptor); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

// MustReasonRegistry creates a registry or panics when descriptors are invalid.
//
// It is intended for static catalog assembly and tests where invalid descriptor
// literals are programming errors.
func MustReasonRegistry(descriptors ...ReasonDescriptor) *ReasonRegistry {
	registry, err := NewReasonRegistry(descriptors...)
	if err != nil {
		panic(err)
	}
	return registry
}

// Register adds descriptor to r.
//
// Register validates descriptor syntax and rejects duplicate reasons. It does
// not close the Reason world: any syntactically valid Reason may be registered
// by the registry owner.
func (r *ReasonRegistry) Register(descriptor ReasonDescriptor) error {
	r.requireNonNil()
	if !descriptor.IsValid() {
		return InvalidReasonDescriptorError{Descriptor: descriptor}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// A registry normally initializes byReason in NewReasonRegistry. The lazy
	// path keeps a manually allocated zero-value ReasonRegistry usable after the
	// nil receiver check without changing the owner-created model.
	if r.byReason == nil {
		r.byReason = make(map[Reason]ReasonDescriptor, 1)
	}
	if _, exists := r.byReason[descriptor.Reason]; exists {
		return DuplicateReasonError{Reason: descriptor.Reason}
	}
	r.byReason[descriptor.Reason] = descriptor

	return nil
}

// Lookup returns the descriptor registered for reason.
//
// Invalid reasons and missing valid reasons both return the zero descriptor and
// false. Invalid lookup keys are treated as absence because lookup is a read
// operation, not a validation report.
func (r *ReasonRegistry) Lookup(reason Reason) (ReasonDescriptor, bool) {
	r.requireNonNil()
	if !reason.IsValid() {
		return ReasonDescriptor{}, false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// A nil map behaves like an empty registry. That is useful for defensive
	// zero-value reads and still returns copy-safe zero values.
	descriptor, ok := r.byReason[reason]
	return descriptor, ok
}

// Contains reports whether reason is registered.
func (r *ReasonRegistry) Contains(reason Reason) bool {
	_, ok := r.Lookup(reason)
	return ok
}

// List returns registered descriptors sorted by Reason string.
//
// The returned slice is a fresh copy. Sorting makes docs, tests, and config
// validation deterministic.
func (r *ReasonRegistry) List() []ReasonDescriptor {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Build a fresh slice before sorting so callers can mutate the returned
	// descriptors without affecting registry state.
	descriptors := make([]ReasonDescriptor, 0, len(r.byReason))
	for _, descriptor := range r.byReason {
		descriptors = append(descriptors, descriptor)
	}
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Reason.String() < descriptors[j].Reason.String()
	})

	return descriptors
}

// Len reports the number of registered reasons.
func (r *ReasonRegistry) Len() int {
	r.requireNonNil()

	r.mu.RLock()
	defer r.mu.RUnlock()

	// len(nil) is zero, so Len is safe even for a manually allocated zero-value
	// registry that has not registered anything yet.
	return len(r.byReason)
}

// requireNonNil panics with the stable nil receiver message used by all
// ReasonRegistry methods.
func (r *ReasonRegistry) requireNonNil() {
	if r == nil {
		panic(nilReasonRegistryPanic)
	}
}
