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

// uint64Payload stores TypeUint64 constraints in their exact native width.
type uint64Payload struct {
	// min stores the inclusive lower bound for TypeUint64.
	min uint64Limit
	// max stores the inclusive upper bound for TypeUint64.
	max uint64Limit
	// enum stores accepted uint64 literals in declaration order.
	enum []uint64
}

// cloneUint64Payload detaches uint64 enum values.
func cloneUint64Payload(p uint64Payload) uint64Payload {
	p.enum = cloneUint64s(p.enum)
	return p
}
