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
	"arcoris.dev/apimachinery/api/value"
)

// equalOpaqueValue compares preserved unknown payloads without descriptor descent.
//
// Unknown-preserved members are semantic leaves for comparison output. Their
// nested payload is inspected only to decide whether that leaf changed.
func (c *comparer) equalOpaqueValue(path fieldpath.Path, oldValue value.Value, newValue value.Value) (bool, error) {
	if oldValue.IsZero() || newValue.IsZero() {
		return false, errorAt(path, ErrInvalidValue, ErrorReasonInvalidZero, "value is the invalid zero Value")
	}
	if oldValue.Kind() != newValue.Kind() {
		return false, nil
	}

	if oldValue.Kind().IsScalar() {
		return opaqueScalarValuesEqual(oldValue, newValue), nil
	}

	switch oldValue.Kind() {
	case value.KindRecord:
		return c.equalOpaqueObject(path, oldValue, newValue)
	case value.KindList:
		return c.equalOpaqueList(path, oldValue, newValue)
	default:
		return false, errorAt(path, ErrInvalidValue, ErrorReasonInvalidZero, "value has invalid kind")
	}
}

// opaqueScalarValuesEqual compares same-kind primitive payloads without a descriptor.
func opaqueScalarValuesEqual(oldValue value.Value, newValue value.Value) bool {
	switch oldValue.Kind() {
	case value.KindNull:
		return true
	case value.KindBool:
		oldBool, _ := oldValue.AsBool()
		newBool, _ := newValue.AsBool()
		return oldBool == newBool
	case value.KindString:
		oldString, _ := oldValue.AsString()
		newString, _ := newValue.AsString()
		return oldString == newString
	case value.KindBytes:
		oldBytes, _ := oldValue.AsBytes()
		newBytes, _ := newValue.AsBytes()
		return bytes.Equal(oldBytes, newBytes)
	case value.KindInteger:
		oldInteger, _ := oldValue.AsInteger()
		newInteger, _ := newValue.AsInteger()
		return oldInteger.Equal(newInteger)
	case value.KindFloat:
		oldFloat, _ := oldValue.AsFloat()
		newFloat, _ := newValue.AsFloat()
		return oldFloat == newFloat
	case value.KindDecimal:
		oldDecimal, _ := oldValue.AsDecimal()
		newDecimal, _ := newValue.AsDecimal()
		return oldDecimal.Compare(newDecimal) == 0
	case value.KindTimestamp:
		oldTimestamp, _ := oldValue.AsTimestamp()
		newTimestamp, _ := newValue.AsTimestamp()
		return oldTimestamp.Equal(newTimestamp)
	case value.KindDate:
		oldDate, _ := oldValue.AsDate()
		newDate, _ := newValue.AsDate()
		return oldDate.Equal(newDate)
	case value.KindTimeOfDay:
		oldTime, _ := oldValue.AsTimeOfDay()
		newTime, _ := newValue.AsTimeOfDay()
		return oldTime.Equal(newTime)
	case value.KindDuration:
		oldDuration, _ := oldValue.AsDuration()
		newDuration, _ := newValue.AsDuration()
		return oldDuration == newDuration
	default:
		return false
	}
}

// equalOpaqueList compares unknown list payloads by exact physical order.
func (c *comparer) equalOpaqueList(path fieldpath.Path, oldValue value.Value, newValue value.Value) (bool, error) {
	oldList, _ := oldValue.AsList()
	newList, _ := newValue.AsList()
	n := oldList.Len()
	if n != newList.Len() {
		return false, nil
	}

	for i := 0; i < n; i++ {
		oldItem, _ := oldList.At(i)
		newItem, _ := newList.At(i)
		equal, err := c.equalOpaqueValue(path.Index(i), oldItem, newItem)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

// equalOpaqueObject compares unknown record payloads by member name.
//
// Unknown-preserved record members have no descriptor, so equality is purely
// structural and exact. Missing or extra members are enough to mark the opaque
// parent leaf as changed.
func (c *comparer) equalOpaqueObject(path fieldpath.Path, oldValue value.Value, newValue value.Value) (bool, error) {
	oldObject, _ := oldValue.AsRecord()
	newObject, _ := newValue.AsRecord()
	if oldObject.Len() != newObject.Len() {
		return false, nil
	}

	oldMembers := membersByName(oldObject.Members())
	for _, newMember := range newObject.Members() {
		name := newMember.Name.String()
		oldMember, found := oldMembers[name]
		if !found {
			return false, nil
		}

		equal, err := c.equalOpaqueValue(path.Field(fieldpath.MustFieldName(name)), oldMember, newMember.Value)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}
