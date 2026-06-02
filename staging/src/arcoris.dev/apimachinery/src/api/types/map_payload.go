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

// mapPayload stores dynamic-map descriptor rules.
//
// Map payloads describe dynamic key/value dictionaries. Fixed schema fields
// belong to objectPayload instead, which preserves declaration order and field
// presence.
type mapPayload struct {
	// key records the structural key family.
	//
	// Only string keys are supported in this design pass.
	key MapKeyType
	// value points at the structural descriptor of every map value.
	//
	// The pointer is private and preserves public value semantics while avoiding
	// a recursive Type -> mapPayload -> Type storage cycle.
	value *Type
	// minLen is the inclusive minimum map size.
	//
	// The limit wrapper distinguishes an explicit zero from an unset rule.
	minLen limit[int]
	// maxLen is the inclusive maximum map size.
	//
	// The limit wrapper distinguishes an explicit zero from an unset rule.
	maxLen limit[int]
}

// cloneMapPayload detaches the map value descriptor.
func cloneMapPayload(p mapPayload) mapPayload {
	if p.value != nil {
		value := cloneType(*p.value)
		p.value = &value
	}

	return p
}

// emptyMapPayload reports whether p has no configured TypeMap state.
func emptyMapPayload(p mapPayload) bool {
	return p.key == MapKeyString && p.value == nil && !p.minLen.set && !p.maxLen.set
}
