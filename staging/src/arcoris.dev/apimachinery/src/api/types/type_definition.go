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

// TypeDefinition binds a reusable structural Type to a stable TypeName.
//
// Definitions are resolved through Resolver implementations supplied by the
// caller. Package types does not own concrete mutable catalog storage and
// performs no global or init-time registration.
type TypeDefinition struct {
	// name is the stable resolver key for the reusable descriptor.
	name TypeName
	// typ is the structural descriptor associated with name.
	typ Type
	// description is optional human-facing descriptor text.
	description string
}

// Define creates a named type definition from typ.
//
// The name parameter accepts string-like values so call sites can use string
// literals in descriptor declarations while the finalized definition still
// stores the closed TypeName value.
//
// Typical declaration:
//
//	nameDef := Define(
//		"arcoris.meta.Name",
//		String().
//			MinLen(1).
//			MaxLen(253),
//	)
func Define[N ~string](name N, typ TypeExpr) TypeDefinition {
	return TypeDefinition{name: TypeName(name), typ: typeFromExpr(typ)}
}

// Name returns the definition name.
func (d TypeDefinition) Name() TypeName {
	return d.name
}

// Type returns a detached copy of the definition type.
func (d TypeDefinition) Type() Type {
	return cloneType(d.typ)
}

// Description returns optional human-facing descriptor text.
func (d TypeDefinition) Description() string {
	return d.description
}

// WithDescription returns a copy of d with human-facing descriptor text.
func (d TypeDefinition) WithDescription(text string) TypeDefinition {
	d.description = text
	return d
}

// cloneDefinition detaches the type payload stored in d.
func cloneDefinition(d TypeDefinition) TypeDefinition {
	d.typ = cloneType(d.typ)
	return d
}
