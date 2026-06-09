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

// uint32Payload stores DescriptorUint32 constraints in their exact native width.
type uint32Payload struct {
	// min stores the inclusive lower bound for DescriptorUint32.
	min limit[uint32]
	// max stores the inclusive upper bound for DescriptorUint32.
	max limit[uint32]
	// enum stores accepted uint32 literals in declaration order.
	enum []uint32
}

// cloneUint32Payload detaches uint32 enum values.
func cloneUint32Payload(p uint32Payload) uint32Payload {
	p.enum = cloneSlice(p.enum)

	return p
}

// emptyUint32Payload reports whether p has no configured DescriptorUint32 state.
func emptyUint32Payload(p uint32Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
