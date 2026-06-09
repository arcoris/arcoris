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

package types

// Definition binds a reusable structural Descriptor to a stable TypeName.
//
// Definitions are resolved through Resolver implementations supplied by the
// caller. Package types does not own concrete mutable catalog storage and
// performs no global or init-time registration.
type Definition struct {
	// name is the stable resolver key for the reusable descriptor.
	name TypeName
	// descriptor is the structural descriptor associated with name.
	descriptor Descriptor
	// description is optional human-facing descriptor text.
	description string
}

// Define creates a named descriptor definition from descriptor.
//
// The name parameter accepts string-like values so call sites can use string
// literals in descriptor declarations while the finalized definition still
// stores the closed TypeName value.
//
// Typical declaration:
//
//	nameDef := Define(
//		"meta.arcoris.dev.Name",
//		String().
//			MinBytes(1).
//			MaxBytes(253),
//	)
func Define[N ~string](name N, descriptor DescriptorExpr) Definition {
	return Definition{name: TypeName(name), descriptor: descriptorFromExpr(descriptor)}
}

// Name returns the definition name.
func (d Definition) Name() TypeName {
	return d.name
}

// Descriptor returns a detached copy of the definition descriptor.
func (d Definition) Descriptor() Descriptor {
	return cloneDescriptor(d.descriptor)
}

// Description returns optional human-facing descriptor text.
func (d Definition) Description() string {
	return d.description
}

// WithDescription returns a copy of d with human-facing descriptor text.
func (d Definition) WithDescription(text string) Definition {
	d.description = text

	return d
}
