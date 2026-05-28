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

func TestTypeCodeValidityAndString(t *testing.T) {
	tests := []struct {
		code  TypeCode
		valid bool
		text  string
	}{
		{TypeInvalid, false, "invalid"},
		{TypeNull, true, "null"},
		{TypeBool, true, "bool"},
		{TypeString, true, "string"},
		{TypeBytes, true, "bytes"},
		{TypeInt8, true, "int8"},
		{TypeInt16, true, "int16"},
		{TypeInt32, true, "int32"},
		{TypeInt64, true, "int64"},
		{TypeUint8, true, "uint8"},
		{TypeUint16, true, "uint16"},
		{TypeUint32, true, "uint32"},
		{TypeUint64, true, "uint64"},
		{TypeFloat32, true, "float32"},
		{TypeFloat64, true, "float64"},
		{TypeDecimal, true, "decimal"},
		{TypeTimestamp, true, "timestamp"},
		{TypeDate, true, "date"},
		{TypeTime, true, "time"},
		{TypeDuration, true, "duration"},
		{TypeObject, true, "object"},
		{TypeList, true, "list"},
		{TypeMap, true, "map"},
		{TypeRef, true, "ref"},
		{TypeCode(255), false, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			requireEqual(t, tt.code.IsValid(), tt.valid)
			requireEqual(t, tt.code.String(), tt.text)
		})
	}
}

func TestTypeCodeFamilyClassification(t *testing.T) {
	tests := []struct {
		code      TypeCode
		primitive bool
		number    bool
		integer   bool
		intCode   bool
		uintCode  bool
		float     bool
		temporal  bool
		composite bool
	}{
		{TypeNull, true, false, false, false, false, false, false, false},
		{TypeBool, true, false, false, false, false, false, false, false},
		{TypeString, true, false, false, false, false, false, false, false},
		{TypeBytes, true, false, false, false, false, false, false, false},
		{TypeInt8, true, true, true, true, false, false, false, false},
		{TypeInt16, true, true, true, true, false, false, false, false},
		{TypeInt32, true, true, true, true, false, false, false, false},
		{TypeInt64, true, true, true, true, false, false, false, false},
		{TypeUint8, true, true, true, false, true, false, false, false},
		{TypeUint16, true, true, true, false, true, false, false, false},
		{TypeUint32, true, true, true, false, true, false, false, false},
		{TypeUint64, true, true, true, false, true, false, false, false},
		{TypeFloat32, true, true, false, false, false, true, false, false},
		{TypeFloat64, true, true, false, false, false, true, false, false},
		{TypeDecimal, true, true, false, false, false, false, false, false},
		{TypeTimestamp, true, false, false, false, false, false, true, false},
		{TypeDate, true, false, false, false, false, false, true, false},
		{TypeTime, true, false, false, false, false, false, true, false},
		{TypeDuration, true, false, false, false, false, false, true, false},
		{TypeObject, false, false, false, false, false, false, false, true},
		{TypeList, false, false, false, false, false, false, false, true},
		{TypeMap, false, false, false, false, false, false, false, true},
		{TypeRef, false, false, false, false, false, false, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			requireEqual(t, tt.code.IsPrimitive(), tt.primitive)
			requireEqual(t, tt.code.IsNumber(), tt.number)
			requireEqual(t, tt.code.IsInteger(), tt.integer)
			requireEqual(t, tt.code.IsInt(), tt.intCode)
			requireEqual(t, tt.code.IsUint(), tt.uintCode)
			requireEqual(t, tt.code.IsFloat(), tt.float)
			requireEqual(t, tt.code.IsTemporal(), tt.temporal)
			requireEqual(t, tt.code.IsComposite(), tt.composite)
		})
	}
}
