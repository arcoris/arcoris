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
	descriptor types.Type,
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
	descriptor types.Type,
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
func scalarKind(code types.TypeCode) (value.Kind, bool) {
	return typekind.Scalar(code)
}

// scalarValuesEqual compares same-kind scalar payloads without descriptor redispatch.
func scalarValuesEqual(oldValue value.Value, newValue value.Value, code types.TypeCode) (bool, bool) {
	switch code {
	case types.TypeBool:
		oldBool, _ := oldValue.Bool()
		newBool, _ := newValue.Bool()
		return oldBool == newBool, true
	case types.TypeString:
		oldString, _ := oldValue.String()
		newString, _ := newValue.String()
		return oldString == newString, true
	case types.TypeBytes:
		oldBytes, _ := oldValue.Bytes()
		newBytes, _ := newValue.Bytes()
		return bytes.Equal(oldBytes, newBytes), true
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		oldInteger, _ := oldValue.Integer()
		newInteger, _ := newValue.Integer()
		return oldInteger.Equal(newInteger), true
	case types.TypeFloat32,
		types.TypeFloat64:
		oldFloat, _ := oldValue.Float()
		newFloat, _ := newValue.Float()
		return oldFloat == newFloat, true
	case types.TypeDecimal:
		oldDecimal, _ := oldValue.Decimal()
		newDecimal, _ := newValue.Decimal()
		return oldDecimal.Compare(newDecimal) == 0, true
	case types.TypeTimestamp:
		oldTimestamp, _ := oldValue.Timestamp()
		newTimestamp, _ := newValue.Timestamp()
		return oldTimestamp.Equal(newTimestamp), true
	case types.TypeDate:
		oldDate, _ := oldValue.Date()
		newDate, _ := newValue.Date()
		return oldDate.Equal(newDate), true
	case types.TypeTime:
		oldTime, _ := oldValue.TimeOfDay()
		newTime, _ := newValue.TimeOfDay()
		return oldTime.Equal(newTime), true
	case types.TypeDuration:
		oldDuration, _ := oldValue.Duration()
		newDuration, _ := newValue.Duration()
		return oldDuration == newDuration, true
	default:
		return false, false
	}
}
