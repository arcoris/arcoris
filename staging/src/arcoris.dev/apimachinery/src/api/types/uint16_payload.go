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

// uint16Payload stores DescriptorUint16 constraints in their exact native width.
type uint16Payload struct {
	// min stores the inclusive lower bound for DescriptorUint16.
	min limit[uint16]
	// max stores the inclusive upper bound for DescriptorUint16.
	max limit[uint16]
	// enum stores accepted uint16 literals in declaration order.
	enum []uint16
}

// cloneUint16Payload detaches uint16 enum values.
func cloneUint16Payload(p uint16Payload) uint16Payload {
	p.enum = cloneSlice(p.enum)

	return p
}

// emptyUint16Payload reports whether p has no configured DescriptorUint16 state.
func emptyUint16Payload(p uint16Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
