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

package value

import (
	"math"
	"strconv"
)

// minInt64Magnitude is abs(math.MinInt64), which cannot be represented as a
// positive int64 but must remain exactly representable by Integer.
const minInt64Magnitude = uint64(math.MaxInt64) + 1

// Integer stores one exact integer in the union of int64 and uint64 domains.
//
// The representation keeps sign and magnitude separate so math.MinInt64 and
// math.MaxUint64 are both represented without loss. Width-specific checks such
// as int32 or uint16 are descriptor validation concerns and intentionally do
// not live in this package.
//
// The magnitude is the absolute value. Negative zero is never constructed by
// public constructors: negative is always false when magnitude is zero.
type Integer struct {
	// negative records the sign for non-zero magnitudes.
	negative bool
	// magnitude stores the absolute value.
	magnitude uint64
}

// IsNegative reports whether i is strictly less than zero.
//
// A malformed private value with negative=true and magnitude=0 is treated as
// zero so externally observed behavior remains canonical.
func (i Integer) IsNegative() bool {
	return i.negative && i.magnitude != 0
}

// Magnitude returns the absolute value of i.
//
// The magnitude is useful for descriptor-aware validation that needs to compare
// against numeric bounds without losing unsigned range.
func (i Integer) Magnitude() uint64 {
	return i.magnitude
}

// FitsInt64 reports whether i can be represented as int64 exactly.
//
// The negative side admits abs(math.MinInt64), while the positive side is
// limited by math.MaxInt64.
func (i Integer) FitsInt64() bool {
	if i.IsNegative() {
		return i.magnitude <= minInt64Magnitude
	}

	return i.magnitude <= uint64(math.MaxInt64)
}

// Int64 returns i as int64 when it fits exactly.
//
// When i does not fit, the returned integer is zero and ok=false. Callers must
// not treat zero as a conversion result unless ok is true.
func (i Integer) Int64() (int64, bool) {
	if !i.FitsInt64() {
		return 0, false
	}

	if !i.IsNegative() {
		return int64(i.magnitude), true
	}

	if i.magnitude == minInt64Magnitude {
		return math.MinInt64, true
	}

	return -int64(i.magnitude), true
}

// FitsUint64 reports whether i can be represented as uint64 exactly.
//
// Every non-negative Integer fits because magnitude is already uint64.
func (i Integer) FitsUint64() bool {
	return !i.IsNegative()
}

// Uint64 returns i as uint64 when it fits exactly.
//
// Negative values return zero and ok=false. This mirrors Int64 and prevents a
// caller from accidentally interpreting a failed conversion as payload data.
func (i Integer) Uint64() (uint64, bool) {
	if !i.FitsUint64() {
		return 0, false
	}

	return i.magnitude, true
}

// String returns the canonical diagnostic decimal text for i.
//
// The text is intended for diagnostics and tests, not as a package-level codec
// contract.
func (i Integer) String() string {
	if i.IsNegative() {
		return "-" + strconv.FormatUint(i.magnitude, 10)
	}

	return strconv.FormatUint(i.magnitude, 10)
}

// Compare returns -1, 0, or 1 when i is less than, equal to, or greater than other.
//
// The comparison works directly on sign and magnitude so no conversion to int64
// or uint64 can overflow.
func (i Integer) Compare(other Integer) int {
	if i.IsNegative() != other.IsNegative() {
		if i.IsNegative() {
			return -1
		}

		return 1
	}

	switch {
	case i.magnitude == other.magnitude:
		return 0
	case i.IsNegative():
		if i.magnitude > other.magnitude {
			return -1
		}
		return 1
	case i.magnitude < other.magnitude:
		return -1
	default:
		return 1
	}
}

// Equal reports whether i and other represent the same integer value.
//
// Equal uses Compare so canonical zero behavior is shared with ordering.
func (i Integer) Equal(other Integer) bool {
	return i.Compare(other) == 0
}
