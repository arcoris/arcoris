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

// Field creates a name-first object field builder.
//
// Field is used only inside Object declarations. It captures the intended API
// field name and then lets the next typed method choose the value family, for
// example Field("replicas").Int64() or Field("labels").MapOf(String()). The
// typed field wrapper then owns presence, nullability, description, and
// type-specific constraints.
//
// Field does not validate the name immediately. Builders stay simple,
// immutable-by-value, and allocation-light; ValidateType owns descriptor-shape
// diagnostics and reports invalid field names with full descriptor paths.
func Field(name string) FieldBuilder {
	return FieldBuilder{name: FieldName(name)}
}

// FieldBuilder is the field-first entrypoint used inside Object declarations.
//
// It is intentionally small: the builder only carries a FieldName until a
// family-specific method creates a concrete field wrapper. FieldBuilder itself
// does not implement FieldExpr, because an object field is incomplete until it
// has a value type and Required or Optional presence.
type FieldBuilder struct {
	// name is retained until a typed field wrapper finalizes the descriptor.
	name FieldName
}

// state converts the entrypoint into shared field builder state.
//
// Family-specific field wrappers embed this state so every wrapper follows the
// same descriptor rules for name, presence, and description while keeping
// type-specific constraints in its own builder value.
func (b FieldBuilder) state() fieldState {
	return fieldState{name: b.name}
}
