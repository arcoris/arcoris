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

package value

import "testing"

func TestValueCoreState(t *testing.T) {
	var zero Value

	requireEqual(t, zero.Kind(), KindInvalid)
	requireEqual(t, zero.IsZero(), true)
	requireEqual(t, zero.IsNull(), false)
	requireEqual(t, zero.IsScalar(), false)
	requireEqual(t, zero.IsComposite(), false)

	requireEqual(t, StringValue("worker").IsScalar(), true)
	requireEqual(t, mustList(t, NullValue()).IsComposite(), true)
}

func TestWrongKindAccessorsReturnFalse(t *testing.T) {
	value := StringValue("worker")

	if _, ok := value.Bool(); ok {
		t.Fatal("BoolValue() ok = true for string value")
	}
	if _, ok := value.Bytes(); ok {
		t.Fatal("BytesValue() ok = true for string value")
	}
	if _, ok := value.Integer(); ok {
		t.Fatal("Integer() ok = true for string value")
	}
	if _, ok := value.Float(); ok {
		t.Fatal("Float() ok = true for string value")
	}
	if _, ok := value.Decimal(); ok {
		t.Fatal("Decimal() ok = true for string value")
	}
	if _, ok := value.Timestamp(); ok {
		t.Fatal("TimestampValue() ok = true for string value")
	}
	if _, ok := value.Date(); ok {
		t.Fatal("Date() ok = true for string value")
	}
	if _, ok := value.TimeOfDay(); ok {
		t.Fatal("TimeOfDay() ok = true for string value")
	}
	if _, ok := value.Duration(); ok {
		t.Fatal("DurationValue() ok = true for string value")
	}
	if _, ok := value.Object(); ok {
		t.Fatal("Object() ok = true for string value")
	}
	if _, ok := value.List(); ok {
		t.Fatal("List() ok = true for string value")
	}
}
