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

// validateFloat64 checks TypeFloat64 finite bounds, enum uniqueness, and enum membership.
func validateFloat64(t Type, path string) error {
	if t.float64.min.set && invalidFloat64(t.float64.min.value) {
		return typeErrorf(
			path+".min",
			ErrInvalidType,
			TypeErrorReasonNonFiniteValue,
			"float64 minimum must be finite, got %s",
			float64Diagnostic(t.float64.min.value),
		)
	}
	if t.float64.max.set && invalidFloat64(t.float64.max.value) {
		return typeErrorf(
			path+".max",
			ErrInvalidType,
			TypeErrorReasonNonFiniteValue,
			"float64 maximum must be finite, got %s",
			float64Diagnostic(t.float64.max.value),
		)
	}
	if err := validateRangeRule(path, "float64", t.float64.min, t.float64.max); err != nil {
		return err
	}
	for i, value := range t.float64.enum {
		if invalidFloat64(value) {
			return typeErrorf(
				fmt.Sprintf("%s.enum[%d]", path, i),
				ErrInvalidType,
				TypeErrorReasonNonFiniteValue,
				"float64 enum value must be finite, got %s",
				float64Diagnostic(value),
			)
		}
	}
	// NaN and infinities are rejected before duplicate checks. Descriptor enum
	// identity then follows Go numeric equality, not float bit-pattern identity.
	return validateEnumRules(path, "float64", t.float64.enum, t.float64.min, t.float64.max)
}

// invalidFloat64 reports whether value is not a finite portable float64 rule.
func invalidFloat64(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}

// float64Diagnostic names non-finite values for validation details.
func float64Diagnostic(value float64) string {
	switch {
	case math.IsNaN(value):
		return "NaN"
	case math.IsInf(value, 1):
		return "+Inf"
	case math.IsInf(value, -1):
		return "-Inf"
	default:
		return fmt.Sprintf("%v", value)
	}
}
