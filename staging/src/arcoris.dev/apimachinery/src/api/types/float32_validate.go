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

import (
	"fmt"
	"math"
)

// validateFloat32 checks TypeFloat32 finite bounds, enum uniqueness, and enum membership.
func validateFloat32(t Type, path string) error {
	if t.float32.min.set && invalidFloat32(t.float32.min.value) {
		return typeErrorf(
			path+".min",
			ErrInvalidType,
			TypeErrorReasonNonFiniteValue,
			"float32 minimum must be finite, got %s",
			float32Diagnostic(t.float32.min.value),
		)
	}

	if t.float32.max.set && invalidFloat32(t.float32.max.value) {
		return typeErrorf(
			path+".max",
			ErrInvalidType,
			TypeErrorReasonNonFiniteValue,
			"float32 maximum must be finite, got %s",
			float32Diagnostic(t.float32.max.value),
		)
	}

	if err := validateRangeRule(path, "float32", t.float32.min, t.float32.max); err != nil {
		return err
	}

	for i, value := range t.float32.enum {
		if invalidFloat32(value) {
			return typeErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidType,
				TypeErrorReasonNonFiniteValue,
				"float32 enum value must be finite, got %s",
				float32Diagnostic(value),
			)
		}
	}

	// NaN and infinities are rejected before duplicate checks. Descriptor enum
	// identity then follows Go numeric equality, not float bit-pattern identity.
	return validateEnumRules(path, "float32", t.float32.enum, t.float32.min, t.float32.max)
}

// invalidFloat32 reports whether value is not a finite portable float32 rule.
func invalidFloat32(value float32) bool {
	return math.IsNaN(float64(value)) || math.IsInf(float64(value), 0)
}

// float32Diagnostic names non-finite values for validation details.
func float32Diagnostic(value float32) string {
	switch {
	case math.IsNaN(float64(value)):
		return "NaN"
	case math.IsInf(float64(value), 1):
		return "+Inf"
	case math.IsInf(float64(value), -1):
		return "-Inf"
	default:
		return fmt.Sprintf("%v", value)
	}
}
