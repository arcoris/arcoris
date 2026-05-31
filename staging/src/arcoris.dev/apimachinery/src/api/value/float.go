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

import "math"

// NewFloat constructs a finite binary floating-point value.
//
// NaN and infinities are rejected because they are not portable API payload
// values unless a future descriptor rule explicitly admits them. The package
// stores float64 only; float32 portability is a descriptor constraint, not a
// separate value kind.
func NewFloat(v float64) (Value, error) {
	if math.IsNaN(v) {
		return Value{}, errorf(
			pathFloat,
			ErrInvalidFloat,
			ErrorReasonInvalidFloat,
			"NaN is not a portable API float value",
		)
	}

	if math.IsInf(v, 0) {
		return Value{}, errorf(
			pathFloat,
			ErrInvalidFloat,
			ErrorReasonInvalidFloat,
			"infinity is not a portable API float value",
		)
	}

	return Value{kind: KindFloat, floatValue: v}, nil
}

// MustFloat constructs a float Value or panics when v is not finite.
//
// It is intended for package-level fixtures and tests where invalid input is a
// programmer error. Runtime parsing paths should use NewFloat and handle the
// structured error.
func MustFloat(v float64) Value {
	value, err := NewFloat(v)
	if err != nil {
		panic(err)
	}

	return value
}
