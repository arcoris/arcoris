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

// float32Payload stores TypeFloat32 constraints in their exact native width.
type float32Payload struct {
	// min stores the inclusive finite lower bound for TypeFloat32.
	min float32Limit
	// max stores the inclusive finite upper bound for TypeFloat32.
	max float32Limit
	// enum stores accepted finite float32 literals in declaration order.
	enum []float32
}

// cloneFloat32Payload detaches float32 enum values.
func cloneFloat32Payload(p float32Payload) float32Payload {
	p.enum = cloneFloat32s(p.enum)
	return p
}
