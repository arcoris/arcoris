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

// uint8Payload stores TypeUint8 constraints in their exact native width.
type uint8Payload struct {
	// min stores the inclusive lower bound for TypeUint8.
	min limit[uint8]
	// max stores the inclusive upper bound for TypeUint8.
	max limit[uint8]
	// enum stores accepted uint8 literals in declaration order.
	enum []uint8
}

// cloneUint8Payload detaches uint8 enum values.
func cloneUint8Payload(p uint8Payload) uint8Payload {
	p.enum = cloneSlice(p.enum)
	return p
}

// emptyUint8Payload reports whether p has no configured TypeUint8 state.
func emptyUint8Payload(p uint8Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
