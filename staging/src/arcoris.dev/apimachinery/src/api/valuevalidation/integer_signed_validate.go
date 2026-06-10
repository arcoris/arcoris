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

// validateSignedInteger checks signed integer descriptors against exact value integers.
func (v *validator) validateSignedInteger(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) {
	if !v.requireKind(path, val, value.KindInteger, descriptor.Code()) {
		return
	}

	integer, _ := val.AsInteger()
	got, ok := integer.Int64()
	if !ok {
		v.add(
			path,
			ErrValueOutOfRange,
			ErrorReasonAboveMaximum,
			"integer does not fit int64",
		)
		return
	}

	switch descriptor.Code() {
	case types.DescriptorInt8:
		v.validateInt8(path, got, descriptor)
	case types.DescriptorInt16:
		v.validateInt16(path, got, descriptor)
	case types.DescriptorInt32:
		v.validateInt32(path, got, descriptor)
	case types.DescriptorInt64:
		v.validateInt64(path, got, descriptor)
	}
}

// validateInt8 checks int8 width, descriptor bounds, and enum rules.
func (v *validator) validateInt8(path fieldpath.Path, got int64, descriptor types.Descriptor) {
	view, ok := descriptor.AsInt8()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not int8",
		)
		return
	}

	validateIntegerLimits(v, path, got, signedWidthLimits(math.MinInt8, math.MaxInt8))
	validateIntegerLimits(v, path, got, signedDescriptorLimits[int8](view.Min, view.Max))
	v.validateSignedEnum(path, got, signedEnum[int8](view.Enum()))
}

// validateInt16 checks int16 width, descriptor bounds, and enum rules.
func (v *validator) validateInt16(path fieldpath.Path, got int64, descriptor types.Descriptor) {
	view, ok := descriptor.AsInt16()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not int16",
		)
		return
	}

	validateIntegerLimits(v, path, got, signedWidthLimits(math.MinInt16, math.MaxInt16))
	validateIntegerLimits(v, path, got, signedDescriptorLimits[int16](view.Min, view.Max))
	v.validateSignedEnum(path, got, signedEnum[int16](view.Enum()))
}

// validateInt32 checks int32 width, descriptor bounds, and enum rules.
func (v *validator) validateInt32(path fieldpath.Path, got int64, descriptor types.Descriptor) {
	view, ok := descriptor.AsInt32()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not int32",
		)
		return
	}

	validateIntegerLimits(v, path, got, signedWidthLimits(math.MinInt32, math.MaxInt32))
	validateIntegerLimits(v, path, got, signedDescriptorLimits[int32](view.Min, view.Max))
	v.validateSignedEnum(path, got, signedEnum[int32](view.Enum()))
}

// validateInt64 checks int64 descriptor bounds and enum rules.
func (v *validator) validateInt64(path fieldpath.Path, got int64, descriptor types.Descriptor) {
	view, ok := descriptor.AsInt64()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not int64",
		)
		return
	}

	validateIntegerLimits(v, path, got, exactIntegerLimits[int64](view.Min, view.Max))
	v.validateSignedEnum(path, got, view.Enum())
}
