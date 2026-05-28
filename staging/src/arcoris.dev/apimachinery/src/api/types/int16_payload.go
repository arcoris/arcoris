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

// int16Payload stores TypeInt16 constraints in their exact native width.
type int16Payload struct {
	// min stores the inclusive lower bound for TypeInt16.
	min int16Limit
	// max stores the inclusive upper bound for TypeInt16.
	max int16Limit
	// enum stores accepted int16 literals in declaration order.
	enum []int16
}

// cloneInt16Payload detaches int16 enum values.
func cloneInt16Payload(p int16Payload) int16Payload {
	p.enum = cloneInt16s(p.enum)
	return p
}
