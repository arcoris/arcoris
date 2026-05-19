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
	"errors"
	"fmt"
)

// ErrNilKindRegistry identifies ComponentRegistry construction without a kind
// catalog.
//
// Component registries need a KindRegistry to distinguish syntactically valid
// but unknown kinds from registered catalog kinds. Use errors.Is to check for
// this configuration failure.
var ErrNilKindRegistry = errors.New("admission: nil kind registry")

// ErrInvalidComponentDescriptor identifies invalid component descriptor values.
//
// Invalid descriptor errors wrap this sentinel and include the rejected
// descriptor in InvalidComponentDescriptorError for callers that need detail.
var ErrInvalidComponentDescriptor = errors.New("admission: invalid component descriptor")

// ErrUnknownComponentKind identifies component descriptors that reference a kind
// absent from the registry's KindRegistry.
//
// Unknown kind errors wrap this sentinel and expose the missing kind through
// UnknownComponentKindError.
var ErrUnknownComponentKind = errors.New("admission: unknown component kind")

// ErrComponentAlreadyRegistered identifies duplicate component registration.
//
// Duplicate component errors wrap this sentinel and expose the conflicting
// component ID through DuplicateComponentError.
var ErrComponentAlreadyRegistered = errors.New("admission: component already registered")

// InvalidComponentDescriptorError reports the descriptor rejected by a component
// registry.
type InvalidComponentDescriptorError struct {
	// Descriptor is the invalid descriptor passed by the caller. It is stored by
	// value so inspecting the error cannot mutate caller-owned input or registry
	// state.
	Descriptor ComponentDescriptor
}

// Error returns a stable diagnostic for invalid component descriptors.
func (e InvalidComponentDescriptorError) Error() string {
	return fmt.Sprintf("%v: id %q", ErrInvalidComponentDescriptor, e.Descriptor.ID)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e InvalidComponentDescriptorError) Unwrap() error {
	return ErrInvalidComponentDescriptor
}

// UnknownComponentKindError reports a descriptor kind missing from the registry
// kind catalog.
type UnknownComponentKindError struct {
	// Kind is the valid but unregistered kind referenced by a component. The
	// syntax is acceptable; the owner-created kind catalog simply does not know
	// it.
	Kind ComponentKind
}

// Error returns a stable diagnostic for unknown component kinds.
func (e UnknownComponentKindError) Error() string {
	return fmt.Sprintf("%v: %q", ErrUnknownComponentKind, e.Kind)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e UnknownComponentKindError) Unwrap() error {
	return ErrUnknownComponentKind
}

// DuplicateComponentError reports the component ID that was already registered.
type DuplicateComponentError struct {
	// ID is the duplicate component ID rejected by the registry. The ID is
	// syntactically valid, but catalog-level uniqueness failed.
	ID ComponentID
}

// Error returns a stable diagnostic for duplicate component registration.
func (e DuplicateComponentError) Error() string {
	return fmt.Sprintf("%v: %q", ErrComponentAlreadyRegistered, e.ID)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e DuplicateComponentError) Unwrap() error {
	return ErrComponentAlreadyRegistered
}
