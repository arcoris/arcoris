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

var (
	// ErrInvalidComponentDescriptor identifies locally invalid component
	// descriptor metadata.
	ErrInvalidComponentDescriptor = errors.New("admissioncatalog: invalid component descriptor")

	// ErrDuplicateComponentDeclaration identifies a repeated component
	// declaration.
	ErrDuplicateComponentDeclaration = errors.New("admissioncatalog: duplicate component declaration")

	// ErrUnknownComponentKind identifies a component descriptor that references
	// an undeclared component kind.
	ErrUnknownComponentKind = errors.New("admissioncatalog: unknown component kind")
)

// InvalidComponentDescriptorError reports invalid component descriptor
// metadata.
type InvalidComponentDescriptorError struct {
	// Descriptor is the rejected value.
	Descriptor ComponentDescriptor

	// Path identifies the input location when known.
	Path string
}

// Error returns a concise diagnostic for the rejected component descriptor.
func (e InvalidComponentDescriptorError) Error() string {
	return formatPathError("invalid component descriptor", e.Path, e.Descriptor.ID.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e InvalidComponentDescriptorError) Unwrap() error {
	return ErrInvalidComponentDescriptor
}

// UnknownComponentKindError reports a component descriptor whose kind has not
// been declared in the assembled catalog.
type UnknownComponentKindError struct {
	// ComponentID identifies the component whose kind is missing.
	ComponentID admission.ComponentID

	// Kind is the undeclared component kind.
	Kind admission.ComponentKind

	// Path identifies the input location when known.
	Path string
}

// Error returns a concise diagnostic for the undeclared component kind.
func (e UnknownComponentKindError) Error() string {
	detail := e.Kind.String()
	if e.ComponentID.IsValid() {
		detail = fmt.Sprintf("%s for component %s", e.Kind, e.ComponentID)
	}
	return formatPathError("unknown component kind", e.Path, detail)
}

// Unwrap returns the sentinel error for errors.Is.
func (e UnknownComponentKindError) Unwrap() error {
	return ErrUnknownComponentKind
}

// DuplicateComponentDeclarationError reports a repeated component declaration.
type DuplicateComponentDeclarationError struct {
	// ID is the duplicated component ID.
	ID admission.ComponentID

	// Path identifies the duplicate input location when known.
	Path string
}

// Error returns a concise diagnostic for the repeated component declaration.
func (e DuplicateComponentDeclarationError) Error() string {
	return formatPathError("duplicate component declaration", e.Path, e.ID.String())
}

// Unwrap returns the sentinel error for errors.Is.
func (e DuplicateComponentDeclarationError) Unwrap() error {
	return ErrDuplicateComponentDeclaration
}
