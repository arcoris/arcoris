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

import "arcoris.dev/admission"

// Builder is an owner-controlled mutable catalog assembly helper.
//
// The zero value is usable. Builder is not safe for concurrent mutation; owners
// should build catalogs during setup and then share the immutable Catalog. A nil
// *Builder receiver panics because it is a wiring bug, not ordinary invalid
// descriptor input.
type Builder struct {
	// reasons stores pending reason declarations by reason.
	reasons descriptorStore[admission.Reason, ReasonDescriptor]

	// kinds stores pending component kind declarations by kind.
	kinds descriptorStore[admission.ComponentKind, ComponentKindDescriptor]

	// components stores pending component declarations by component ID.
	components descriptorStore[admission.ComponentID, ComponentDescriptor]
}

// newInitializedBuilder returns a builder whose stores are ready for
// path-aware internal declaration helpers.
func newInitializedBuilder() *Builder {
	builder := &Builder{}
	builder.init()
	return builder
}

// init makes the zero-value builder's stores usable without dropping existing
// declarations.
func (b *Builder) init() {
	initReasonStore(&b.reasons)
	initKindStore(&b.kinds)
	initComponentStore(&b.components)
}

// requireNonNil enforces the package-wide nil receiver policy for Builder.
func (b *Builder) requireNonNil() {
	if b == nil {
		panic(nilBuilderPanic)
	}
}
