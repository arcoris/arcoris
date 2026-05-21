/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package capacity

// Amount is a scalar number of local capacity units.
//
// Amount intentionally does not define what one unit means. A component may map
// one unit to one worker slot, one request-cost unit, one buffer class unit, or
// another component-local resource model. Multi-resource accounting, request
// cost estimation, tenant weighting, and scheduling policy belong to higher
// layers.
//
// The zero Amount is valid as a ledger limit. It is invalid as a reservation
// request because reserving zero units would create ownership without resource
// movement.
type Amount uint64

// IsZero reports whether a is zero.
func (a Amount) IsZero() bool {
	return a == 0
}

// IsPositive reports whether a can represent a reservation request.
func (a Amount) IsPositive() bool {
	return a > 0
}

// Uint64 returns a as uint64.
func (a Amount) Uint64() uint64 {
	return uint64(a)
}
