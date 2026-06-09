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

// ListDescriptor builds list descriptors.
//
// ListDescriptor describes a homogeneous sequence.
//
// The descriptor records element descriptor, length constraints, and future
// merge/apply intent, but it does not implement scheduling, queues, patch/apply,
// field ownership, or concrete value validation.
type ListDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact list shape under construction.
	payload listPayload
}

// ListOf returns a list descriptor builder for elem.
//
// A nil DescriptorExpr is accepted as an invalid zero element descriptor so
// ValidateResolved can report the shape error at list.elem. The default list
// semantics are atomic, which is the most conservative structural merge intent.
//
// Typical reusable declaration:
//
//	conditionListType := ListOf(
//		Ref("arcoris.meta.Condition"),
//	).Map("type")
func ListOf(elem DescriptorExpr) ListDescriptor {
	elemType := descriptorFromExpr(elem)

	return ListDescriptor{
		header: newHeader(DescriptorList),
		payload: listPayload{
			elem:      &elemType,
			semantics: ListAtomic,
		},
	}
}

// Nullable returns a list descriptor that admits null values.
func (desc ListDescriptor) Nullable() ListDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// MinItems sets the inclusive minimum list item count.
func (desc ListDescriptor) MinItems(n int) ListDescriptor {
	desc.payload.minLen = limit[int]{n, true}

	return desc
}

// MaxItems sets the inclusive maximum list item count.
func (desc ListDescriptor) MaxItems(n int) ListDescriptor {
	desc.payload.maxLen = limit[int]{n, true}

	return desc
}

// Atomic records atomic list semantics.
//
// Atomic semantics are the conservative default. They mean future field-set,
// ownership, merge, and apply layers should treat the complete list as a single
// replaceable semantic field. Validation may still inspect individual items and
// report item-level diagnostics by index.
func (desc ListDescriptor) Atomic() ListDescriptor {
	desc.payload.semantics = ListAtomic
	desc.payload.mapKeys = nil

	return desc
}

// Ordered records index-addressable list semantics.
//
// Ordered semantics mean physical item indexes are part of the API contract.
// Future field-set, diff, and apply layers may therefore treat indexes as
// semantic addresses. Use ordered lists only for truly positional values, such
// as ordered command arguments or ordered pipeline stages.
//
// Atomic remains the default because it treats the complete list as one field.
func (desc ListDescriptor) Ordered() ListDescriptor {
	desc.payload.semantics = ListOrdered
	desc.payload.mapKeys = nil

	return desc
}

// Set records set-like list semantics.
//
// Set semantics record that future merge/apply layers may treat list elements
// as identity-less set members. Until a stable value-based set identity model
// exists, field-set extraction treats the complete list as one field and must
// not use physical indexes for ownership purposes.
//
// This package does not compare elements or enforce set uniqueness for concrete
// values.
func (desc ListDescriptor) Set() ListDescriptor {
	desc.payload.semantics = ListSet
	desc.payload.mapKeys = nil

	return desc
}

// Map records ListMap semantics keyed by object field names.
//
// ValidateResolved later checks that map keys are valid field names, that the list
// element is an object or resolvable object reference, that each key field is
// required, and that each key field resolves to a non-nullable stable scalar
// identity type suitable for selector-based validation, field-set extraction,
// diff, and apply layers.
//
// Field-set extraction can then address items by selector rather than unstable
// physical index. The builder only records the declared key order.
func (desc ListDescriptor) Map(keys ...string) ListDescriptor {
	desc.payload.semantics = ListMap
	desc.payload.mapKeys = make([]FieldName, len(keys))

	for i, key := range keys {
		desc.payload.mapKeys[i] = FieldName(key)
	}

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc ListDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.list = cloneListPayload(desc.payload)

	return out
}

// descriptorExpr marks ListDescriptor as a sealed DescriptorExpr implementation.
func (desc ListDescriptor) descriptorExpr() {}
