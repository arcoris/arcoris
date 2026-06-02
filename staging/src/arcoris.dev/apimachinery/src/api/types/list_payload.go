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

// listPayload stores element structure, length limits, and list semantics.
//
// List payloads carry structural metadata only. They do not compare concrete
// elements, enforce set uniqueness, perform keyed merges, or implement patch
// and apply semantics.
type listPayload struct {
	// elem points at the structural descriptor of each list element.
	//
	// The pointer is an internal representation detail. Public APIs still expose
	// Type by value, but the pointer breaks the otherwise recursive Type ->
	// listPayload -> Type storage cycle that Go cannot represent directly.
	elem *Type
	// minLen is the inclusive minimum list length.
	//
	// The limit wrapper distinguishes an explicit zero from an unset rule.
	minLen limit[int]
	// maxLen is the inclusive maximum list length.
	//
	// The limit wrapper distinguishes an explicit zero from an unset rule.
	maxLen limit[int]
	// semantics records future merge/apply intent.
	//
	// This package validates that the semantic value is known, but does not
	// execute merge/apply behavior.
	semantics ListSemantics
	// mapKeys stores object field names used by ListMap semantics.
	//
	// ValidateType checks that these names are valid, present on the element
	// object, and required.
	mapKeys []FieldName
}

// cloneListPayload detaches the list element and map keys.
func cloneListPayload(p listPayload) listPayload {
	if p.elem != nil {
		elem := cloneType(*p.elem)
		p.elem = &elem
	}
	p.mapKeys = cloneSlice(p.mapKeys)

	return p
}

// emptyListPayload reports whether p has no configured TypeList state.
func emptyListPayload(p listPayload) bool {
	return p.elem == nil &&
		!p.minLen.set &&
		!p.maxLen.set &&
		p.semantics == ListAtomic &&
		len(p.mapKeys) == 0
}
