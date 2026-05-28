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

// validateFloat32 checks TypeFloat32 finite bounds, enum uniqueness, and enum membership.
func validateFloat32(t Type, path string) error {
	if t.float32.min.set && invalidFloat32(t.float32.min.value) {
		return typeError(path+".range", ErrInvalidType)
	}
	if t.float32.max.set && invalidFloat32(t.float32.max.value) {
		return typeError(path+".range", ErrInvalidType)
	}
	if invalidRange(t.float32.min, t.float32.max) {
		return typeError(path+".range", ErrInvalidType)
	}
	for _, value := range t.float32.enum {
		if invalidFloat32(value) {
			return typeError(path+".enum", ErrInvalidType)
		}
	}
	if hasDuplicates(t.float32.enum) || enumBelowMin(t.float32.enum, t.float32.min) || enumAboveMax(t.float32.enum, t.float32.max) {
		return typeError(path+".enum", ErrInvalidType)
	}
	return nil
}

// invalidFloat32 reports whether value is not a finite portable float32 rule.
func invalidFloat32(value float32) bool {
	return math.IsNaN(float64(value)) || math.IsInf(float64(value), 0)
}
