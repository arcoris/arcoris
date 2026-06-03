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
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareScalar records the current path when scalar payloads differ.
func (c *comparer) compareScalar(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
) (Result, error) {
	equal, err := c.equalValue(path, oldValue, newValue, descriptor, 0)
	if err != nil {
		return Result{}, err
	}
	if equal {
		return EmptyResult(), nil
	}

	return EmptyResult().withModified(path)
}

// equalScalar compares scalar payloads under scalar descriptor semantics.
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

// scalarKind maps a scalar descriptor code to its concrete value kind.
func scalarKind(code types.TypeCode) (value.Kind, bool) {
	switch code {
	case types.TypeBool:
		return value.KindBool, true
	case types.TypeString:
		return value.KindString, true
	case types.TypeBytes:
		return value.KindBytes, true
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		return value.KindInteger, true
	case types.TypeFloat32,
		types.TypeFloat64:
		return value.KindFloat, true
	case types.TypeDecimal:
		return value.KindDecimal, true
	case types.TypeTimestamp:
		return value.KindTimestamp, true
	case types.TypeDate:
		return value.KindDate, true
	case types.TypeTime:
		return value.KindTimeOfDay, true
	case types.TypeDuration:
		return value.KindDuration, true
	default:
		return value.KindInvalid, false
	}
}

// scalarValuesEqual compares payloads after the caller checked kind compatibility.
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
