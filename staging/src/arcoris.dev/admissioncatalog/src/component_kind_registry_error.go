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

import (
	"errors"
	"fmt"

	"arcoris.dev/admission"
)

// ErrInvalidComponentKindDescriptor identifies invalid kind descriptor values.
//
// Use errors.Is to check for this class of registry failure. Use errors.As with
// InvalidComponentKindDescriptorError when the rejected descriptor is needed.
var ErrInvalidComponentKindDescriptor = errors.New("admissioncatalog: invalid component kind descriptor")

// ErrComponentKindAlreadyRegistered identifies duplicate kind registration.
//
// Use errors.Is to check for duplicate registration and errors.As with
// DuplicateComponentKindError to inspect the conflicting kind.
var ErrComponentKindAlreadyRegistered = errors.New("admissioncatalog: component kind already registered")

// InvalidComponentKindDescriptorError reports the descriptor rejected by a kind
// registry.
type InvalidComponentKindDescriptorError struct {
	// Descriptor is the invalid descriptor passed by the caller. It is stored by
	// value so inspecting the error cannot mutate caller-owned input or registry
	// state.
	Descriptor ComponentKindDescriptor
}

// Error returns a stable diagnostic for invalid kind descriptor registration.
func (e InvalidComponentKindDescriptorError) Error() string {
	return fmt.Sprintf("%v: kind %q", ErrInvalidComponentKindDescriptor, e.Descriptor.Kind)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e InvalidComponentKindDescriptorError) Unwrap() error {
	return ErrInvalidComponentKindDescriptor
}

// DuplicateComponentKindError reports the kind that was already registered.
type DuplicateComponentKindError struct {
	// Kind is the duplicate kind rejected by the registry. The value is
	// syntactically valid, but catalog-level uniqueness failed.
	Kind admission.ComponentKind
}

// Error returns a stable diagnostic for duplicate kind registration.
func (e DuplicateComponentKindError) Error() string {
	return fmt.Sprintf("%v: %q", ErrComponentKindAlreadyRegistered, e.Kind)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e DuplicateComponentKindError) Unwrap() error {
	return ErrComponentKindAlreadyRegistered
}
