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

// uint32Payload stores TypeUint32 constraints in their exact native width.
type uint32Payload struct {
	// min stores the inclusive lower bound for TypeUint32.
	min uint32Limit
	// max stores the inclusive upper bound for TypeUint32.
	max uint32Limit
	// enum stores accepted uint32 literals in declaration order.
	enum []uint32
}

// cloneUint32Payload detaches uint32 enum values.
func cloneUint32Payload(p uint32Payload) uint32Payload {
	p.enum = cloneUint32s(p.enum)
	return p
}
