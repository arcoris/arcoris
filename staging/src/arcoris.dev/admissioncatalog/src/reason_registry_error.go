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

// ErrInvalidReasonDescriptor identifies invalid reason descriptor values.
//
// Invalid descriptor errors wrap this sentinel and include the rejected
// descriptor in InvalidReasonDescriptorError for callers that need detail.
var ErrInvalidReasonDescriptor = errors.New("admissioncatalog: invalid reason descriptor")

// ErrReasonAlreadyRegistered identifies duplicate reason registration.
//
// Duplicate reason errors wrap this sentinel and expose the conflicting reason
// through DuplicateReasonError.
var ErrReasonAlreadyRegistered = errors.New("admissioncatalog: reason already registered")

// InvalidReasonDescriptorError reports the descriptor rejected by a reason
// registry.
type InvalidReasonDescriptorError struct {
	// Descriptor is the invalid descriptor passed by the caller. It is stored by
	// value so inspecting the error cannot mutate caller-owned input or registry
	// state.
	Descriptor ReasonDescriptor
}

// Error returns a stable diagnostic for invalid reason descriptor registration.
func (e InvalidReasonDescriptorError) Error() string {
	return fmt.Sprintf("%v: reason %q", ErrInvalidReasonDescriptor, e.Descriptor.Reason)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e InvalidReasonDescriptorError) Unwrap() error {
	return ErrInvalidReasonDescriptor
}

// DuplicateReasonError reports the reason that was already registered.
type DuplicateReasonError struct {
	// Reason is the duplicate reason rejected by the registry. The value is
	// syntactically valid, but catalog-level uniqueness failed.
	Reason admission.Reason
}

// Error returns a stable diagnostic for duplicate reason registration.
func (e DuplicateReasonError) Error() string {
	return fmt.Sprintf("%v: %q", ErrReasonAlreadyRegistered, e.Reason)
}

// Unwrap exposes the sentinel error for errors.Is.
func (e DuplicateReasonError) Unwrap() error {
	return ErrReasonAlreadyRegistered
}
