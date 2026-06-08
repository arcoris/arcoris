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

package retrybudget

import (
	"errors"
	"math/bits"
	"strconv"
)

// ErrInvalidRatio reports an invalid retry-budget ratio.
var ErrInvalidRatio = errors.New("retrybudget: invalid ratio")

// Ratio is an exact retry allowance ratio in the closed range [0, 1].
//
// Ratio is immutable and safe to copy by value. Values are created with NewRatio
// or MustRatio so numerator and denominator are reduced to a canonical form.
// The zero value is intentionally invalid; use RatioZero for a valid 0% ratio.
type Ratio struct {
	// numerator is the reduced non-negative numerator. It is always less than or
	// equal to denominator on valid Ratio values.
	numerator uint64

	// denominator is the reduced positive denominator. A zero denominator marks an
	// invalid or unset Ratio.
	denominator uint64
}

var (
	// RatioZero is the exact 0% retry allowance ratio.
	RatioZero = Ratio{numerator: 0, denominator: 1}

	// RatioOne is the exact 100% retry allowance ratio.
	RatioOne = Ratio{numerator: 1, denominator: 1}
)

// NewRatio returns the canonical exact ratio numerator/denominator.
//
// The denominator must be non-zero and the numerator must not exceed the
// denominator. Ratios above one are rejected because retry budgets should remain
// conservative by default.
func NewRatio(numerator, denominator uint64) (Ratio, error) {
	if denominator == 0 || numerator > denominator {
		return Ratio{}, ErrInvalidRatio
	}
	if numerator == 0 {
		return RatioZero, nil
	}
	divisor := gcd(numerator, denominator)
	return Ratio{
		numerator:   numerator / divisor,
		denominator: denominator / divisor,
	}, nil
}

// MustRatio returns the canonical exact ratio numerator/denominator.
//
// MustRatio panics with ErrInvalidRatio when NewRatio would return an error. It
// is intended for package defaults, examples, and tests where invalid ratios are
// programmer mistakes.
func MustRatio(numerator, denominator uint64) Ratio {
	ratio, err := NewRatio(numerator, denominator)
	if err != nil {
		panic(err)
	}
	return ratio
}

// Numerator returns r's reduced numerator.
func (r Ratio) Numerator() uint64 {
	return r.numerator
}

// Denominator returns r's reduced denominator.
func (r Ratio) Denominator() uint64 {
	return r.denominator
}

// IsValid reports whether r is an exact ratio in the closed range [0, 1].
func (r Ratio) IsValid() bool {
	return r.denominator != 0 && r.numerator <= r.denominator
}

// IsZero reports whether r is the valid 0% ratio.
func (r Ratio) IsZero() bool {
	return r.IsValid() && r.numerator == 0
}

// IsOne reports whether r is the valid 100% ratio.
func (r Ratio) IsOne() bool {
	return r.IsValid() && r.numerator == r.denominator
}

// ScaleFloor returns floor(value * r).
//
// The calculation is exact for the full uint64 range. Invalid ratios panic with
// ErrInvalidRatio because scale operations are programmer errors after
// configuration validation.
func (r Ratio) ScaleFloor(value uint64) uint64 {
	if !r.IsValid() {
		panic(ErrInvalidRatio)
	}
	if value == 0 || r.numerator == 0 {
		return 0
	}
	hi, lo := bits.Mul64(value, r.numerator)
	scaled, _ := bits.Div64(hi, lo, r.denominator)
	return scaled
}

// String returns a stable human-readable representation of r.
func (r Ratio) String() string {
	if !r.IsValid() {
		return "invalid"
	}
	if r.IsZero() {
		return "0"
	}
	if r.IsOne() {
		return "1"
	}
	return strconv.FormatUint(r.numerator, 10) + "/" +
		strconv.FormatUint(r.denominator, 10)
}

// gcd returns the greatest common divisor of a and b.
//
// NewRatio calls gcd only after validating b is non-zero. The loop is the
// standard Euclidean algorithm and keeps Ratio construction allocation-free.
func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
