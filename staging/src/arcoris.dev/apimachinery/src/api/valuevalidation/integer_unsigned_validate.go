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

// validateUnsignedInteger checks unsigned integer descriptors against exact value integers.
func (v *validator) validateUnsignedInteger(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
) {
	if !v.requireKind(path, val, value.KindInteger, descriptor.Code()) {
		return
	}

	integer, _ := val.AsInteger()
	got, ok := integer.Uint64()
	if !ok {
		v.add(
			path,
			ErrValueOutOfRange,
			ErrorReasonBelowMinimum,
			"integer is negative",
		)
		return
	}

	switch descriptor.Code() {
	case types.DescriptorUint8:
		v.validateUint8(path, got, descriptor)
	case types.DescriptorUint16:
		v.validateUint16(path, got, descriptor)
	case types.DescriptorUint32:
		v.validateUint32(path, got, descriptor)
	case types.DescriptorUint64:
		v.validateUint64(path, got, descriptor)
	}
}

// validateUint8 checks uint8 width, descriptor bounds, and enum rules.
func (v *validator) validateUint8(path fieldpath.Path, got uint64, descriptor types.Descriptor) {
	view, ok := descriptor.AsUint8()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not uint8",
		)
		return
	}

	validateIntegerLimits(v, path, got, unsignedWidthLimits(math.MaxUint8))
	validateIntegerLimits(v, path, got, unsignedDescriptorLimits[uint8](view.Min, view.Max))
	v.validateUnsignedEnum(path, got, unsignedEnum[uint8](view.Enum()))
}

// validateUint16 checks uint16 width, descriptor bounds, and enum rules.
func (v *validator) validateUint16(path fieldpath.Path, got uint64, descriptor types.Descriptor) {
	view, ok := descriptor.AsUint16()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not uint16",
		)
		return
	}

	validateIntegerLimits(v, path, got, unsignedWidthLimits(math.MaxUint16))
	validateIntegerLimits(v, path, got, unsignedDescriptorLimits[uint16](view.Min, view.Max))
	v.validateUnsignedEnum(path, got, unsignedEnum[uint16](view.Enum()))
}

// validateUint32 checks uint32 width, descriptor bounds, and enum rules.
func (v *validator) validateUint32(path fieldpath.Path, got uint64, descriptor types.Descriptor) {
	view, ok := descriptor.AsUint32()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not uint32",
		)
		return
	}

	validateIntegerLimits(v, path, got, unsignedWidthLimits(math.MaxUint32))
	validateIntegerLimits(v, path, got, unsignedDescriptorLimits[uint32](view.Min, view.Max))
	v.validateUnsignedEnum(path, got, unsignedEnum[uint32](view.Enum()))
}

// validateUint64 checks uint64 descriptor bounds and enum rules.
func (v *validator) validateUint64(path fieldpath.Path, got uint64, descriptor types.Descriptor) {
	view, ok := descriptor.AsUint64()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not uint64",
		)
		return
	}

	validateIntegerLimits(v, path, got, exactIntegerLimits[uint64](view.Min, view.Max))
	v.validateUnsignedEnum(path, got, view.Enum())
}
