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

// int64Payload stores TypeInt64 constraints in their exact native width.
type int64Payload struct {
	// min stores the inclusive lower bound for TypeInt64.
	min limit[int64]
	// max stores the inclusive upper bound for TypeInt64.
	max limit[int64]
	// enum stores accepted int64 literals in declaration order.
	enum []int64
}

// cloneInt64Payload detaches int64 enum values.
func cloneInt64Payload(p int64Payload) int64Payload {
	p.enum = cloneSlice(p.enum)
	return p
}

// emptyInt64Payload reports whether p has no configured TypeInt64 state.
func emptyInt64Payload(p int64Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
