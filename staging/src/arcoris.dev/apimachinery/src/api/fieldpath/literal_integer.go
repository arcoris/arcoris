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

package fieldpath

import (
	"math"
	"strconv"
)

// integer stores one exact integer from the int64 ∪ uint64 domain.
type integer struct {
	negative  bool
	magnitude uint64
}

// newInt64 stores v without losing math.MinInt64.
func newInt64(v int64) integer {
	if v >= 0 {
		return integer{magnitude: uint64(v)}
	}

	return integer{
		negative:  true,
		magnitude: uint64(-(v + 1)) + 1,
	}
}

// newUint64 stores v as a non-negative integer.
func newUint64(v uint64) integer {
	return integer{magnitude: v}
}

// validate reports whether i satisfies the internal sign/magnitude invariant.
func (i integer) validate() error {
	if i.negative && i.magnitude == 0 {
		return newError(
			ErrInvalidLiteral,
			ErrorReasonInvalidLiteral,
			"integer literal uses an invalid negative zero representation",
		)
	}

	return nil
}

// fitsInt64 reports whether i can be represented as int64.
func (i integer) fitsInt64() bool {
	if i.negative {
		return i.magnitude <= uint64(math.MaxInt64)+1
	}

	return i.magnitude <= uint64(math.MaxInt64)
}

// int64Value returns i as int64 when the value is representable.
func (i integer) int64Value() (int64, bool) {
	if !i.fitsInt64() {
		return 0, false
	}

	if i.negative {
		if i.magnitude == uint64(math.MaxInt64)+1 {
			return math.MinInt64, true
		}

		return -int64(i.magnitude), true
	}

	return int64(i.magnitude), true
}

// uint64Value returns i as uint64 when the value is non-negative.
func (i integer) uint64Value() (uint64, bool) {
	if i.negative {
		return 0, false
	}

	return i.magnitude, true
}

// compare imposes numeric ordering across the full int64 ∪ uint64 domain.
func (i integer) compare(other integer) int {
	switch {
	case i.negative && !other.negative:
		return -1
	case !i.negative && other.negative:
		return 1
	case i.negative && other.negative:
		switch {
		case i.magnitude > other.magnitude:
			return -1
		case i.magnitude < other.magnitude:
			return 1
		default:
			return 0
		}
	default:
		switch {
		case i.magnitude < other.magnitude:
			return -1
		case i.magnitude > other.magnitude:
			return 1
		default:
			return 0
		}
	}
}

// string returns the canonical decimal diagnostic form of i.
func (i integer) string() string {
	if i.negative {
		return "-" + strconv.FormatUint(i.magnitude, 10)
	}

	return strconv.FormatUint(i.magnitude, 10)
}
