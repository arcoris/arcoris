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

// DeclareReason adds descriptor to the builder.
//
// The method validates local descriptor shape and rejects duplicate reason
// declarations. It does not require the reason to appear in any global
// vocabulary.
func (b *Builder) DeclareReason(descriptor ReasonDescriptor) error {
	b.requireNonNil()
	b.init()
	return b.declareReason(descriptor, "")
}

// DeclareKind adds descriptor to the builder.
//
// The method validates local descriptor shape and rejects duplicate component
// kind declarations.
func (b *Builder) DeclareKind(descriptor ComponentKindDescriptor) error {
	b.requireNonNil()
	b.init()
	return b.declareKind(descriptor, "")
}

// DeclareComponent adds descriptor to the builder.
//
// The descriptor's kind must already be declared in the builder. Batch
// construction through Build may list kinds and components in any order because
// Build processes all kinds before all components.
func (b *Builder) DeclareComponent(descriptor ComponentDescriptor) error {
	b.requireNonNil()
	b.init()
	return b.declareComponent(descriptor, "")
}

// declareReason validates and records descriptor with an optional diagnostic
// path. The caller must initialize b before calling it.
func (b *Builder) declareReason(descriptor ReasonDescriptor, path string) error {
	if !descriptor.IsValid() {
		return InvalidReasonDescriptorError{Descriptor: descriptor, Path: path}
	}
	if !b.reasons.declare(descriptor) {
		return DuplicateReasonDeclarationError{Reason: descriptor.Reason, Path: path}
	}
	return nil
}

// declareKind validates and records descriptor with an optional diagnostic path.
// The caller must initialize b before calling it.
func (b *Builder) declareKind(descriptor ComponentKindDescriptor, path string) error {
	if !descriptor.IsValid() {
		return InvalidComponentKindDescriptorError{Descriptor: descriptor, Path: path}
	}
	if !b.kinds.declare(descriptor) {
		return DuplicateComponentKindDeclarationError{Kind: descriptor.Kind, Path: path}
	}
	return nil
}

// declareComponent validates descriptor, checks kind membership, and records it
// with an optional diagnostic path. The caller must initialize b before calling
// it.
func (b *Builder) declareComponent(descriptor ComponentDescriptor, path string) error {
	if !descriptor.IsValid() {
		return InvalidComponentDescriptorError{Descriptor: descriptor, Path: path}
	}
	if !b.kinds.has(descriptor.Kind) {
		return UnknownComponentKindError{
			ComponentID: descriptor.ID,
			Kind:        descriptor.Kind,
			Path:        path,
		}
	}
	if !b.components.declare(descriptor) {
		return DuplicateComponentDeclarationError{ID: descriptor.ID, Path: path}
	}
	return nil
}
