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

func TestDescriptorKindValidityAndString(t *testing.T) {
	tests := []struct {
		code  DescriptorKind
		valid bool
		text  string
	}{
		{DescriptorInvalid, false, "invalid"},
		{DescriptorNull, true, "null"},
		{DescriptorBool, true, "bool"},
		{DescriptorString, true, "string"},
		{DescriptorBytes, true, "bytes"},
		{DescriptorInt8, true, "int8"},
		{DescriptorInt16, true, "int16"},
		{DescriptorInt32, true, "int32"},
		{DescriptorInt64, true, "int64"},
		{DescriptorUint8, true, "uint8"},
		{DescriptorUint16, true, "uint16"},
		{DescriptorUint32, true, "uint32"},
		{DescriptorUint64, true, "uint64"},
		{DescriptorFloat32, true, "float32"},
		{DescriptorFloat64, true, "float64"},
		{DescriptorDecimal, true, "decimal"},
		{DescriptorTimestamp, true, "timestamp"},
		{DescriptorDate, true, "date"},
		{DescriptorTime, true, "time"},
		{DescriptorDuration, true, "duration"},
		{DescriptorObject, true, "object"},
		{DescriptorList, true, "list"},
		{DescriptorMap, true, "map"},
		{DescriptorRef, true, "ref"},
		{DescriptorKind(255), false, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			requireEqual(t, tt.code.IsValid(), tt.valid)
			requireEqual(t, tt.code.String(), tt.text)
		})
	}
}

func TestDescriptorKindFamilyClassification(t *testing.T) {
	tests := []struct {
		code      DescriptorKind
		primitive bool
		number    bool
		integer   bool
		intCode   bool
		uintCode  bool
		float     bool
		temporal  bool
		composite bool
	}{
		{DescriptorNull, true, false, false, false, false, false, false, false},
		{DescriptorBool, true, false, false, false, false, false, false, false},
		{DescriptorString, true, false, false, false, false, false, false, false},
		{DescriptorBytes, true, false, false, false, false, false, false, false},
		{DescriptorInt8, true, true, true, true, false, false, false, false},
		{DescriptorInt16, true, true, true, true, false, false, false, false},
		{DescriptorInt32, true, true, true, true, false, false, false, false},
		{DescriptorInt64, true, true, true, true, false, false, false, false},
		{DescriptorUint8, true, true, true, false, true, false, false, false},
		{DescriptorUint16, true, true, true, false, true, false, false, false},
		{DescriptorUint32, true, true, true, false, true, false, false, false},
		{DescriptorUint64, true, true, true, false, true, false, false, false},
		{DescriptorFloat32, true, true, false, false, false, true, false, false},
		{DescriptorFloat64, true, true, false, false, false, true, false, false},
		{DescriptorDecimal, true, true, false, false, false, false, false, false},
		{DescriptorTimestamp, true, false, false, false, false, false, true, false},
		{DescriptorDate, true, false, false, false, false, false, true, false},
		{DescriptorTime, true, false, false, false, false, false, true, false},
		{DescriptorDuration, true, false, false, false, false, false, true, false},
		{DescriptorObject, false, false, false, false, false, false, false, true},
		{DescriptorList, false, false, false, false, false, false, false, true},
		{DescriptorMap, false, false, false, false, false, false, false, true},
		{DescriptorRef, false, false, false, false, false, false, false, false},
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
