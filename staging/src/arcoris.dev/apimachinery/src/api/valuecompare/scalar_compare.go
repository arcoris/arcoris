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

package valuecompare

import (
	"bytes"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/typekind"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareScalar records the current path when descriptor-compatible scalars differ.
func (c *comparer) compareScalar(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
) (Result, error) {
	equal, err := c.equalScalar(path, oldValue, newValue, descriptor)
	if err != nil {
		return Result{}, err
	}
	if equal {
		return EmptyResult(), nil
	}

	return EmptyResult().withModified(path)
}

// equalScalar compares scalar payloads after enforcing descriptor-compatible kind.
//
// Decimal equality uses value.Decimal.Compare so numerically equal values with
// different scales compare equal. Bytes compare by byte content.
func (c *comparer) equalScalar(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
) (bool, error) {
	expected, ok := scalarKind(descriptor.Code())
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a scalar type")
	}
	if err := requireKind(path, oldValue, expected, descriptor.Code()); err != nil {
		return false, err
	}
	if err := requireKind(path, newValue, expected, descriptor.Code()); err != nil {
		return false, err
	}

	equal, ok := scalarValuesEqual(oldValue, newValue, descriptor.Code())
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a scalar type")
	}

	return equal, nil
}

// scalarKind maps scalar descriptor codes to their concrete payload kind.
func scalarKind(code types.DescriptorKind) (value.Kind, bool) {
	return typekind.Scalar(code)
}

// scalarValuesEqual compares same-kind scalar payloads without descriptor redispatch.
func scalarValuesEqual(oldValue value.Value, newValue value.Value, code types.DescriptorKind) (bool, bool) {
	switch code {
	case types.DescriptorBool:
		oldBool, _ := oldValue.AsBool()
		newBool, _ := newValue.AsBool()
		return oldBool == newBool, true
	case types.DescriptorString:
		oldString, _ := oldValue.AsString()
		newString, _ := newValue.AsString()
		return oldString == newString, true
	case types.DescriptorBytes:
		oldBytes, _ := oldValue.AsBytes()
		newBytes, _ := newValue.AsBytes()
		return bytes.Equal(oldBytes, newBytes), true
	case types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64,
		types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64:
		oldInteger, _ := oldValue.AsInteger()
		newInteger, _ := newValue.AsInteger()
		return oldInteger.Equal(newInteger), true
	case types.DescriptorFloat32,
		types.DescriptorFloat64:
		oldFloat, _ := oldValue.AsFloat()
		newFloat, _ := newValue.AsFloat()
		return oldFloat == newFloat, true
	case types.DescriptorDecimal:
		oldDecimal, _ := oldValue.AsDecimal()
		newDecimal, _ := newValue.AsDecimal()
		return oldDecimal.Compare(newDecimal) == 0, true
	case types.DescriptorTimestamp:
		oldTimestamp, _ := oldValue.AsTimestamp()
		newTimestamp, _ := newValue.AsTimestamp()
		return oldTimestamp.Equal(newTimestamp), true
	case types.DescriptorDate:
		oldDate, _ := oldValue.AsDate()
		newDate, _ := newValue.AsDate()
		return oldDate.Equal(newDate), true
	case types.DescriptorTime:
		oldTime, _ := oldValue.AsTimeOfDay()
		newTime, _ := newValue.AsTimeOfDay()
		return oldTime.Equal(newTime), true
	case types.DescriptorDuration:
		oldDuration, _ := oldValue.AsDuration()
		newDuration, _ := newValue.AsDuration()
		return oldDuration == newDuration, true
	default:
		return false, false
	}
}
