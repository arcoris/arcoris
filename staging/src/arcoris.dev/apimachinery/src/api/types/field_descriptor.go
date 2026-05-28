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

// FieldDescriptor is a finalized object field descriptor.
//
// FieldDescriptor values are produced by field builders and stored by object
// descriptors in declaration order. The fields are immutable-by-convention;
// accessors return detached Type copies so callers cannot rewrite object
// descriptors through a field returned by a view or resolved definition.
type FieldDescriptor struct {
	// name is the lowerCamelCase API field name.
	name FieldName
	// presence records whether the field key is required or optional.
	presence FieldPresence
	// typ is the structural descriptor of the field value.
	typ Type
	// description is optional human-facing descriptor text.
	description string
}

// Name returns the field name.
func (f FieldDescriptor) Name() FieldName {
	return f.name
}

// Presence returns the field presence requirement.
func (f FieldDescriptor) Presence() FieldPresence {
	return f.presence
}

// Type returns a detached copy of the field value descriptor.
func (f FieldDescriptor) Type() Type {
	return cloneType(f.typ)
}

// IsRequired reports whether the field key must be present.
func (f FieldDescriptor) IsRequired() bool {
	return f.presence == PresenceRequired
}

// IsOptional reports whether the field key may be omitted.
func (f FieldDescriptor) IsOptional() bool {
	return f.presence == PresenceOptional
}

// Description returns optional human-facing descriptor text.
func (f FieldDescriptor) Description() string {
	return f.description
}

// IsZero reports whether f is the empty, unfinalized descriptor.
//
// A zero field can be produced only by zero values or nil FieldExpr entries in
// package-local construction paths. Object validation rejects zero fields
// because real API fields must have a valid name, required/optional presence,
// and a valid normalized Type.
func (f FieldDescriptor) IsZero() bool {
	return f.name == "" &&
		f.presence == PresenceUnspecified &&
		f.typ.IsZero() &&
		f.description == ""
}
