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

	"arcoris.dev/apimachinery/api/value"
)

// opaqueScalarValuesEqual compares same-kind primitive values without descriptor context.
func opaqueScalarValuesEqual(oldValue value.Value, newValue value.Value) bool {
	switch oldValue.Kind() {
	case value.KindNull:
		return true
	case value.KindBool:
		oldBool, _ := oldValue.Bool()
		newBool, _ := newValue.Bool()
		return oldBool == newBool
	case value.KindString:
		oldString, _ := oldValue.String()
		newString, _ := newValue.String()
		return oldString == newString
	case value.KindBytes:
		oldBytes, _ := oldValue.Bytes()
		newBytes, _ := newValue.Bytes()
		return bytes.Equal(oldBytes, newBytes)
	case value.KindInteger:
		oldInteger, _ := oldValue.Integer()
		newInteger, _ := newValue.Integer()
		return oldInteger.Equal(newInteger)
	case value.KindFloat:
		oldFloat, _ := oldValue.Float()
		newFloat, _ := newValue.Float()
		return oldFloat == newFloat
	case value.KindDecimal:
		oldDecimal, _ := oldValue.Decimal()
		newDecimal, _ := newValue.Decimal()
		return oldDecimal.Compare(newDecimal) == 0
	case value.KindTimestamp:
		oldTimestamp, _ := oldValue.Timestamp()
		newTimestamp, _ := newValue.Timestamp()
		return oldTimestamp.Equal(newTimestamp)
	case value.KindDate:
		oldDate, _ := oldValue.Date()
		newDate, _ := newValue.Date()
		return oldDate.Equal(newDate)
	case value.KindTimeOfDay:
		oldTime, _ := oldValue.TimeOfDay()
		newTime, _ := newValue.TimeOfDay()
		return oldTime.Equal(newTime)
	case value.KindDuration:
		oldDuration, _ := oldValue.Duration()
		newDuration, _ := newValue.Duration()
		return oldDuration == newDuration
	default:
		return false
	}
}
