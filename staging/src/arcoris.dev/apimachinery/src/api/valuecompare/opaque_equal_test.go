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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
	"testing"
)

func TestEqualOpaqueValueDifferentKindsIsFalse(t *testing.T) {
	got, err := newComparer(Options{}).equalOpaqueValue(fieldpath.RootPath(), value.StringValue("x"), value.BoolValue(true))
	requireNoError(t, err)

	if got {
		t.Fatalf("equalOpaqueValue() = true")
	}
}

func TestEqualOpaqueValueRejectsZeroValue(t *testing.T) {
	_, err := newComparer(Options{}).equalOpaqueValue(fieldpath.RootPath(), value.Value{}, value.StringValue("x"))

	requireErrorIs(t, err, ErrInvalidValue)
}
func TestEqualOpaqueListSameItemsIsTrue(t *testing.T) {
	oldValue := value.MustListValue(value.StringValue("a"))
	newValue := value.MustListValue(value.StringValue("a"))

	got, err := newComparer(Options{}).equalOpaqueList(rootField("items"), oldValue, newValue)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalOpaqueList() = false")
	}
}

func TestEqualOpaqueListDifferentLengthIsFalse(t *testing.T) {
	oldValue := value.MustListValue(value.StringValue("a"))
	newValue := value.MustListValue(value.StringValue("a"), value.StringValue("b"))

	got, err := newComparer(Options{}).equalOpaqueList(rootField("items"), oldValue, newValue)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalOpaqueList() = true")
	}
}
func TestEqualOpaqueObjectComparesStructure(t *testing.T) {
	oldValue := valueObject("nested", "one")
	newValue := valueObject("nested", "two")

	got, err := newComparer(Options{}).equalOpaqueObject(rootField("extra"), oldValue, newValue)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalOpaqueObject() = true")
	}
}

func TestEqualOpaqueObjectMissingMemberIsFalse(t *testing.T) {
	got, err := newComparer(Options{}).equalOpaqueObject(rootField("extra"), valueObject("nested", "one"), valueObject())
	requireNoError(t, err)

	if got {
		t.Fatalf("equalOpaqueObject() = true")
	}
}
func TestOpaqueScalarValuesEqualBytes(t *testing.T) {
	if !opaqueScalarValuesEqual(value.BytesValue([]byte("a")), value.BytesValue([]byte("a"))) {
		t.Fatalf("opaqueScalarValuesEqual(bytes) = false")
	}
	if opaqueScalarValuesEqual(value.BytesValue([]byte("a")), value.BytesValue([]byte("b"))) {
		t.Fatalf("opaqueScalarValuesEqual(bytes) = true")
	}
}

func TestOpaqueScalarValuesEqualDecimalUsesNumericCompare(t *testing.T) {
	if !opaqueScalarValuesEqual(mustDecimal(t, "1.0"), mustDecimal(t, "1.00")) {
		t.Fatalf("opaqueScalarValuesEqual(decimal) = false")
	}
}
