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

// ObjectType builds structural object descriptors.
//
// ObjectType describes fixed-field records. It is not a map, not a Go struct
// reflection wrapper, and not a runtime object implementation. Field order is
// preserved exactly as declared so diagnostics, generated output, and future
// schema exports can remain deterministic.
type ObjectType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact object shape under construction.
	payload objectPayload
}

// Object returns an object descriptor builder containing fields in order.
//
// The constructor accepts package-sealed FieldExpr values so external packages
// cannot inject arbitrary field implementations. Nil field expressions become
// zero field descriptors; ValidateType reports those with descriptor paths
// instead of forcing constructors to panic.
//
// Typical reusable declaration:
//
//	specType := Object(
//		Field("image").String().
//			Required().
//			MinLen(1),
//	).UnknownFields(UnknownReject)
func Object(fields ...FieldExpr) ObjectType {
	payload := objectPayload{unknown: UnknownReject, fields: make([]FieldDescriptor, 0, len(fields))}

	for _, expr := range fields {
		payload.fields = append(payload.fields, fieldFromExpr(expr))
	}

	return ObjectType{header: newHeader(TypeObject), payload: payload}
}

// Nullable returns an object descriptor that admits null values.
func (t ObjectType) Nullable() ObjectType {
	t.header = t.header.withNullable()

	return t
}

// UnknownFields records the structural policy for undeclared object fields.
//
// The policy is only descriptor metadata here. This package does not reject,
// prune, or preserve concrete unknown fields; future value-processing layers
// interpret the policy when they operate on actual API objects.
func (t ObjectType) UnknownFields(policy UnknownFieldPolicy) ObjectType {
	t.payload.unknown = policy

	return t
}

// Type returns a detached Type descriptor.
func (t ObjectType) Type() Type {
	out := typeFromHeader(t.header)
	out.object = cloneObjectPayload(t.payload)

	return out
}

// typeExpr marks ObjectType as a sealed TypeExpr implementation.
func (t ObjectType) typeExpr() {}
