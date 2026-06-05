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

package capacity

// Amount is an exact unsigned local capacity quantity.
//
// Amount intentionally does not define what one unit means. A Resource or
// higher-level package may map one unit to a worker slot, byte, request-cost
// unit, retry-budget unit, or another local accounting dimension.
//
// The zero Amount is representable and useful in snapshots. It is rejected in
// vector and demand entries because canonical resource vectors omit zeroes.
type Amount uint64

// IsZero reports whether a is zero.
func (a Amount) IsZero() bool {
	return a == 0
}

// IsPositive reports whether a is greater than zero.
func (a Amount) IsPositive() bool {
	return a > 0
}

// Uint64 returns a as its raw architecture-independent integer value.
func (a Amount) Uint64() uint64 {
	return uint64(a)
}

// Compare compares a and b.
func (a Amount) Compare(b Amount) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// CheckedAdd returns a+b and whether the sum fits in Amount.
func (a Amount) CheckedAdd(b Amount) (Amount, bool) {
	sum := a + b
	return sum, sum >= a
}

// CheckedSub returns a-b and whether a covers b.
func (a Amount) CheckedSub(b Amount) (Amount, bool) {
	if a < b {
		return 0, false
	}
	return a - b, true
}

// SaturatingSub returns a-b, clamped at zero.
func (a Amount) SaturatingSub(b Amount) Amount {
	if a < b {
		return 0
	}
	return a - b
}
