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

package types

import "testing"

func TestTypeZeroValueInvalid(t *testing.T) {
	var typ Type

	requireEqual(t, typ.IsZero(), true)
	requireEqual(t, typ.IsValid(), false)
	requireEqual(t, typ.Code(), TypeInvalid)
	requireErrorIs(t, ValidateType(typ, nil), ErrInvalidTypeCode)
}

func TestConstructorsCreateExpectedTypeCodes(t *testing.T) {
	tests := []struct {
		name string
		typ  Type
		code TypeCode
	}{
		{"null", Null().Type(), TypeNull},
		{"bool", Bool().Type(), TypeBool},
		{"string", String().Type(), TypeString},
		{"bytes", Bytes().Type(), TypeBytes},
		{"int8", Int8().Type(), TypeInt8},
		{"int16", Int16().Type(), TypeInt16},
		{"int32", Int32().Type(), TypeInt32},
		{"int64", Int64().Type(), TypeInt64},
		{"uint8", Uint8().Type(), TypeUint8},
		{"uint16", Uint16().Type(), TypeUint16},
		{"uint32", Uint32().Type(), TypeUint32},
		{"uint64", Uint64().Type(), TypeUint64},
		{"float32", Float32().Type(), TypeFloat32},
		{"float64", Float64().Type(), TypeFloat64},
		{"decimal", Decimal().Type(), TypeDecimal},
		{"timestamp", Timestamp().Type(), TypeTimestamp},
		{"date", Date().Type(), TypeDate},
		{"time", Time().Type(), TypeTime},
		{"duration", Duration().Type(), TypeDuration},
		{"object", Object().Type(), TypeObject},
		{"list", ListOf(String()).Type(), TypeList},
		{"map", MapOf(String()).Type(), TypeMap},
		{"ref", Ref("arcoris.meta.Name").Type(), TypeRef},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.typ.Code(), tt.code)
			requireEqual(t, tt.typ.IsValid(), true)
		})
	}
}

func TestNullableFlagBehavior(t *testing.T) {
	tests := []Type{
		Bool().Nullable().Type(),
		String().Nullable().Type(),
		Int64().Nullable().Type(),
		Uint64().Nullable().Type(),
		Float64().Nullable().Type(),
		Decimal().Nullable().Type(),
		Timestamp().Nullable().Type(),
		Object().Nullable().Type(),
		ListOf(String()).Nullable().Type(),
		MapOf(String()).Nullable().Type(),
		Ref("arcoris.meta.Name").Nullable().Type(),
	}
	for _, typ := range tests {
		requireEqual(t, typ.Nullable(), true)
	}
	requireEqual(t, Null().Type().Nullable(), false)
}

func TestTypeNullCannotBeNullable(t *testing.T) {
	typ := Null().Type()
	typ.flags = typeFlagNullable

	requireErrorIs(t, ValidateType(typ, nil), ErrInvalidType)
}
