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

// int32Payload stores TypeInt32 constraints in their exact native width.
type int32Payload struct {
	// min stores the inclusive lower bound for TypeInt32.
	min limit[int32]
	// max stores the inclusive upper bound for TypeInt32.
	max limit[int32]
	// enum stores accepted int32 literals in declaration order.
	enum []int32
}

// cloneInt32Payload detaches int32 enum values.
func cloneInt32Payload(p int32Payload) int32Payload {
	p.enum = cloneSlice(p.enum)

	return p
}

// emptyInt32Payload reports whether p has no configured TypeInt32 state.
func emptyInt32Payload(p int32Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
