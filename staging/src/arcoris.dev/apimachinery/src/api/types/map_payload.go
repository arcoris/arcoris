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
// presence. Map keys remain concrete string tokens, but the key descriptor can
// constrain those tokens or reference a reusable string-like definition.
type mapPayload struct {
	// key points at the structural descriptor for concrete string map keys.
	//
	// The pointer is private and preserves public value semantics while avoiding
	// a recursive Descriptor -> mapPayload -> Descriptor storage cycle.
	key *Descriptor
	// value points at the structural descriptor of every map value.
	//
	// The pointer is private and preserves public value semantics while avoiding
	// a recursive Descriptor -> mapPayload -> Descriptor storage cycle.
	value *Descriptor
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
	if p.key != nil {
		key := cloneDescriptor(*p.key)
		p.key = &key
	}

	if p.value != nil {
		value := cloneDescriptor(*p.value)
		p.value = &value
	}

	return p
}

// emptyMapPayload reports whether p has no configured DescriptorMap state.
func emptyMapPayload(p mapPayload) bool {
	return p.key == nil && p.value == nil && !p.minLen.set && !p.maxLen.set
}
