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

	"arcoris.dev/admission"
)

var (
	// ErrInvalidReasonDescriptor identifies locally invalid reason descriptor
	// metadata.
	ErrInvalidReasonDescriptor = errors.New("admissioncatalog: invalid reason descriptor")

	// ErrDuplicateReasonDeclaration identifies a repeated reason declaration.
	ErrDuplicateReasonDeclaration = errors.New("admissioncatalog: duplicate reason declaration")
)

// InvalidReasonDescriptorError reports invalid reason descriptor metadata.
type InvalidReasonDescriptorError struct {
	// Descriptor is the rejected value.
	Descriptor ReasonDescriptor

	// Path identifies the input location when known.
	Path string
}

// Error returns a concise diagnostic for the rejected reason descriptor.
func (e InvalidReasonDescriptorError) Error() string {
	return formatPathError("invalid reason descriptor", e.Path, e.Descriptor.Reason.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e InvalidReasonDescriptorError) Unwrap() error {
	return ErrInvalidReasonDescriptor
}

// DuplicateReasonDeclarationError reports a repeated reason declaration.
type DuplicateReasonDeclarationError struct {
	// Reason is the duplicated reason.
	Reason admission.Reason

	// Path identifies the duplicate input location when known.
	Path string
}

// Error returns a concise diagnostic for the repeated reason declaration.
func (e DuplicateReasonDeclarationError) Error() string {
	return formatPathError("duplicate reason declaration", e.Path, e.Reason.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e DuplicateReasonDeclarationError) Unwrap() error {
	return ErrDuplicateReasonDeclaration
}
