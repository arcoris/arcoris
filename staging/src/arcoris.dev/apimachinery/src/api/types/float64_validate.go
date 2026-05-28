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

import "math"

// validateFloat64 checks TypeFloat64 finite bounds, enum uniqueness, and enum membership.
func validateFloat64(t Type, path string) error {
	if t.float64.min.set && invalidFloat64(t.float64.min.value) {
		return typeError(path+".range", ErrInvalidType)
	}
	if t.float64.max.set && invalidFloat64(t.float64.max.value) {
		return typeError(path+".range", ErrInvalidType)
	}
	if invalidRange(t.float64.min, t.float64.max) {
		return typeError(path+".range", ErrInvalidType)
	}
	for _, value := range t.float64.enum {
		if invalidFloat64(value) {
			return typeError(path+".enum", ErrInvalidType)
		}
	}
	if hasDuplicates(t.float64.enum) || enumBelowMin(t.float64.enum, t.float64.min) || enumAboveMax(t.float64.enum, t.float64.max) {
		return typeError(path+".enum", ErrInvalidType)
	}
	return nil
}

// invalidFloat64 reports whether value is not a finite portable float64 rule.
func invalidFloat64(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}
