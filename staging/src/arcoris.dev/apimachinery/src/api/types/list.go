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

// ListType builds list descriptors.
//
// ListType describes a homogeneous sequence. The descriptor records element
// type, length constraints, and future merge/apply intent, but it does not
// implement scheduling, queues, patch/apply, field ownership, or concrete value
// validation.
type ListType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact list shape under construction.
	payload listPayload
}

// ListOf returns a list descriptor builder for elem.
//
// A nil TypeExpr is accepted as an invalid zero element descriptor so
// ValidateType can report the shape error at list.elem. The default list
// semantics are atomic, which is the most conservative structural merge intent.
//
// Typical reusable declaration:
//
//	conditionListType := ListOf(
//		Ref("arcoris.meta.Condition"),
//	).Map("type")
func ListOf(elem TypeExpr) ListType {
	elemType := typeFromExpr(elem)

	return ListType{
		header: newHeader(TypeList),
		payload: listPayload{
			elem:      &elemType,
			semantics: ListAtomic,
		},
	}
}

// Nullable returns a list descriptor that admits null values.
func (t ListType) Nullable() ListType {
	t.header = t.header.withNullable()

	return t
}

// MinLen sets the inclusive minimum list length.
func (t ListType) MinLen(n int) ListType {
	t.payload.minLen = limit[int]{n, true}

	return t
}

// MaxLen sets the inclusive maximum list length.
func (t ListType) MaxLen(n int) ListType {
	t.payload.maxLen = limit[int]{n, true}

	return t
}

// Atomic records atomic list semantics.
//
// Atomic semantics mean future merge/apply layers should treat the complete
// list as a single replaceable value. This package records that intent only.
func (t ListType) Atomic() ListType {
	t.payload.semantics = ListAtomic
	t.payload.mapKeys = nil

	return t
}

// Ordered records index-addressable list semantics.
//
// Ordered semantics mean future field-set, diff, and apply layers may treat
// physical item indexes as semantic addresses. Use ordered lists only when item
// position is part of the API contract. Atomic remains the conservative default
// because it treats the complete list as one field.
func (t ListType) Ordered() ListType {
	t.payload.semantics = ListOrdered
	t.payload.mapKeys = nil

	return t
}

// Set records set-like list semantics.
//
// Set semantics record that future merge/apply layers may treat list elements
// as identity-less set members. Until a stable value-based identity model
// exists, field-set extraction treats the complete list as one field. This
// package does not compare elements or enforce set uniqueness for concrete
// values.
func (t ListType) Set() ListType {
	t.payload.semantics = ListSet
	t.payload.mapKeys = nil

	return t
}

// Map records map-like list semantics keyed by object field names.
//
// ValidateType later checks that map keys are valid field names, that the list
// element is an object or resolvable object reference, that each key field is
// required, and that each key field resolves to a non-nullable stable scalar
// identity type suitable for future selector-based validation, diff, and apply
// layers. Field-set extraction can then address items by selector rather than
// unstable physical index. The builder only records the declared key order.
func (t ListType) Map(keys ...string) ListType {
	t.payload.semantics = ListMap
	t.payload.mapKeys = make([]FieldName, len(keys))

	for i, key := range keys {
		t.payload.mapKeys[i] = FieldName(key)
	}

	return t
}

// Type returns a detached Type descriptor.
func (t ListType) Type() Type {
	out := typeFromHeader(t.header)
	out.list = cloneListPayload(t.payload)

	return out
}

// typeExpr marks ListType as a sealed TypeExpr implementation.
func (t ListType) typeExpr() {}
