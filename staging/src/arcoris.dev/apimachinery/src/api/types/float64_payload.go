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

// float64Payload stores TypeFloat64 constraints in their exact native width.
type float64Payload struct {
	// min stores the inclusive finite lower bound for TypeFloat64.
	min limit[float64]
	// max stores the inclusive finite upper bound for TypeFloat64.
	max limit[float64]
	// enum stores accepted finite float64 literals in declaration order.
	enum []float64
}

// cloneFloat64Payload detaches float64 enum values.
func cloneFloat64Payload(p float64Payload) float64Payload {
	p.enum = cloneSlice(p.enum)

	return p
}

// emptyFloat64Payload reports whether p has no configured TypeFloat64 state.
func emptyFloat64Payload(p float64Payload) bool {
	return !p.min.set && !p.max.set && len(p.enum) == 0
}
