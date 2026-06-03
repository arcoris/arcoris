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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestCompareSameStringIsEmpty(t *testing.T) {
	got, err := Compare(value.StringValue("a"), value.StringValue("a"), types.String().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareDifferentStringIsModified(t *testing.T) {
	got, err := Compare(value.StringValue("a"), value.StringValue("b"), types.String().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareDifferentBoolIsModified(t *testing.T) {
	got, err := Compare(value.BoolValue(false), value.BoolValue(true), types.Bool().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareDifferentBytesIsModified(t *testing.T) {
	got, err := Compare(value.BytesValue([]byte("a")), value.BytesValue([]byte("b")), types.Bytes().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareDifferentIntegerIsModified(t *testing.T) {
	got, err := Compare(value.Int64Value(1), value.Int64Value(2), types.Int64().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareDifferentFloatIsModified(t *testing.T) {
	got, err := Compare(mustFloat(t, 1.25), mustFloat(t, 2.5), types.Float64().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareDecimalEqualDifferentScaleIsEmpty(t *testing.T) {
	got, err := Compare(mustDecimal(t, "1.0"), mustDecimal(t, "1.00"), types.Decimal().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareDecimalDifferentValueIsModified(t *testing.T) {
	got, err := Compare(mustDecimal(t, "1.0"), mustDecimal(t, "1.01"), types.Decimal().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareNullToNullIsEmpty(t *testing.T) {
	got, err := Compare(value.NullValue(), value.NullValue(), types.String().Nullable().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareNullToScalarIsModified(t *testing.T) {
	got, err := Compare(value.NullValue(), value.StringValue("x"), types.String().Nullable().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}

func TestCompareScalarToNullIsModified(t *testing.T) {
	got, err := Compare(value.StringValue("x"), value.NullValue(), types.String().Nullable().Type(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(fieldpath.RootPath()))
}
