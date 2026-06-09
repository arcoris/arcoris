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

package valuevalidation

import (
	"math"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateFloat32 checks float32 bounds, range, and enum constraints.
func (v *validator) validateFloat32(path fieldpath.Path, val value.Value, descriptor types.Descriptor) {
	if !v.requireKind(path, val, value.KindFloat, descriptor.Code()) {
		return
	}

	got, _ := val.Float()
	if math.Abs(got) > math.MaxFloat32 {
		v.add(path, ErrValueOutOfRange, ErrorReasonAboveMaximum, "float does not fit float32")
		return
	}

	view, ok := descriptor.AsFloat32()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not float32")
		return
	}

	if minValue, ok := view.Min(); ok && got < float64(minValue) {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonBelowMinimum,
			"float %v is below minimum %v",
			got,
			minValue,
		)
	}
	if maxValue, ok := view.Max(); ok && got > float64(maxValue) {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"float %v is above maximum %v",
			got,
			maxValue,
		)
	}
	if enum := view.Enum(); len(enum) > 0 && !containsFloat32(enum, got) {
		v.add(path, ErrEnumMismatch, ErrorReasonEnumMismatch, "float value is not in enum")
	}
}

// validateFloat64 checks float64 bounds and enum constraints.
func (v *validator) validateFloat64(path fieldpath.Path, val value.Value, descriptor types.Descriptor) {
	if !v.requireKind(path, val, value.KindFloat, descriptor.Code()) {
		return
	}

	got, _ := val.Float()
	view, ok := descriptor.AsFloat64()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not float64")
		return
	}

	if minValue, ok := view.Min(); ok && got < minValue {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonBelowMinimum,
			"float %v is below minimum %v",
			got,
			minValue,
		)
	}
	if maxValue, ok := view.Max(); ok && got > maxValue {
		v.addf(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"float %v is above maximum %v",
			got,
			maxValue,
		)
	}
	if enum := view.Enum(); len(enum) > 0 && !containsFloat64(enum, got) {
		v.add(path, ErrEnumMismatch, ErrorReasonEnumMismatch, "float value is not in enum")
	}
}

// containsFloat32 reports whether a float64 payload matches a float32 enum value.
func containsFloat32(values []float32, target float64) bool {
	for _, candidate := range values {
		if float64(candidate) == target {
			return true
		}
	}

	return false
}

// containsFloat64 reports whether values contains target.
func containsFloat64(values []float64, target float64) bool {
	for _, candidate := range values {
		if candidate == target {
			return true
		}
	}

	return false
}
