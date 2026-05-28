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

// int8Payload stores TypeInt8 constraints in their exact native width.
type int8Payload struct {
	// min stores the inclusive lower bound for TypeInt8.
	min int8Limit
	// max stores the inclusive upper bound for TypeInt8.
	max int8Limit
	// enum stores accepted int8 literals in declaration order.
	enum []int8
}

// cloneInt8Payload detaches int8 enum values.
func cloneInt8Payload(p int8Payload) int8Payload {
	p.enum = cloneInt8s(p.enum)
	return p
}
