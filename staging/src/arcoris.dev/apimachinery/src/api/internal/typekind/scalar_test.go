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

package typekind

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestScalar(t *testing.T) {
	cases := []struct {
		name string
		code types.TypeCode
		want value.Kind
		ok   bool
	}{
		{name: "bool", code: types.TypeBool, want: value.KindBool, ok: true},
		{name: "string", code: types.TypeString, want: value.KindString, ok: true},
		{name: "bytes", code: types.TypeBytes, want: value.KindBytes, ok: true},
		{name: "int8", code: types.TypeInt8, want: value.KindInteger, ok: true},
		{name: "int16", code: types.TypeInt16, want: value.KindInteger, ok: true},
		{name: "int32", code: types.TypeInt32, want: value.KindInteger, ok: true},
		{name: "int64", code: types.TypeInt64, want: value.KindInteger, ok: true},
		{name: "uint8", code: types.TypeUint8, want: value.KindInteger, ok: true},
		{name: "uint16", code: types.TypeUint16, want: value.KindInteger, ok: true},
		{name: "uint32", code: types.TypeUint32, want: value.KindInteger, ok: true},
		{name: "uint64", code: types.TypeUint64, want: value.KindInteger, ok: true},
		{name: "float32", code: types.TypeFloat32, want: value.KindFloat, ok: true},
		{name: "float64", code: types.TypeFloat64, want: value.KindFloat, ok: true},
		{name: "decimal", code: types.TypeDecimal, want: value.KindDecimal, ok: true},
		{name: "timestamp", code: types.TypeTimestamp, want: value.KindTimestamp, ok: true},
		{name: "date", code: types.TypeDate, want: value.KindDate, ok: true},
		{name: "time", code: types.TypeTime, want: value.KindTimeOfDay, ok: true},
		{name: "duration", code: types.TypeDuration, want: value.KindDuration, ok: true},
		{name: "null", code: types.TypeNull, want: value.KindInvalid, ok: false},
		{name: "object", code: types.TypeObject, want: value.KindInvalid, ok: false},
		{name: "map", code: types.TypeMap, want: value.KindInvalid, ok: false},
		{name: "list", code: types.TypeList, want: value.KindInvalid, ok: false},
		{name: "ref", code: types.TypeRef, want: value.KindInvalid, ok: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Scalar(tt.code)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("Scalar(%s) = %s, %v; want %s, %v", tt.code, got, ok, tt.want, tt.ok)
			}
		})
	}
}
