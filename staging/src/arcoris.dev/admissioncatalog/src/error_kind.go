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
	// ErrInvalidComponentKindDescriptor identifies locally invalid component kind
	// descriptor metadata.
	ErrInvalidComponentKindDescriptor = errors.New("admissioncatalog: invalid component kind descriptor")

	// ErrDuplicateComponentKindDeclaration identifies a repeated component kind
	// declaration.
	ErrDuplicateComponentKindDeclaration = errors.New("admissioncatalog: duplicate component kind declaration")
)

// InvalidComponentKindDescriptorError reports invalid component kind metadata.
type InvalidComponentKindDescriptorError struct {
	// Descriptor is the rejected value.
	Descriptor ComponentKindDescriptor

	// Path identifies the input location when known.
	Path string
}

// Error returns a concise diagnostic for the rejected kind descriptor.
func (e InvalidComponentKindDescriptorError) Error() string {
	return formatPathError("invalid component kind descriptor", e.Path, e.Descriptor.Kind.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e InvalidComponentKindDescriptorError) Unwrap() error {
	return ErrInvalidComponentKindDescriptor
}

// DuplicateComponentKindDeclarationError reports a repeated component kind
// declaration.
type DuplicateComponentKindDeclarationError struct {
	// Kind is the duplicated component kind.
	Kind admission.ComponentKind

	// Path identifies the duplicate input location when known.
	Path string
}

// Error returns a concise diagnostic for the repeated kind declaration.
func (e DuplicateComponentKindDeclarationError) Error() string {
	return formatPathError("duplicate component kind declaration", e.Path, e.Kind.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e DuplicateComponentKindDeclarationError) Unwrap() error {
	return ErrDuplicateComponentKindDeclaration
}
