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

func TestValidateTypeRejectsInactivePayloads(t *testing.T) {
	tests := []Type{
		func() Type { typ := String().Type(); typ.bytes.minLen = limit[int]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.int8.min = limit[int8]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.int16.min = limit[int16]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.int32.min = limit[int32]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.int64.min = limit[int64]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.uint8.min = limit[uint8]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.uint16.min = limit[uint16]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.uint32.min = limit[uint32]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.uint64.min = limit[uint64]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.float32.min = limit[float32]{value: 1, set: true}; return typ }(),
		func() Type { typ := String().Type(); typ.float64.min = limit[float64]{value: 1, set: true}; return typ }(),
		func() Type {
			typ := String().Type()
			typ.decimal.precision = limit[int]{value: 1, set: true}
			return typ
		}(),
		func() Type {
			typ := String().Type()
			typ.object.fields = []FieldDescriptor{Field("name").String().Required().Field()}
			return typ
		}(),
		func() Type { typ := String().Type(); elem := Bool().Type(); typ.list.elem = &elem; return typ }(),
		func() Type { typ := String().Type(); value := Bool().Type(); typ.mapType.value = &value; return typ }(),
		func() Type { typ := String().Type(); typ.ref.name = "example.Name"; return typ }(),
		func() Type { typ := Int8().Type(); typ.int16.min = limit[int16]{value: 1, set: true}; return typ }(),
		func() Type { typ := Int8().Type(); typ.string.minLen = limit[int]{value: 1, set: true}; return typ }(),
		func() Type { typ := Uint64().Type(); typ.int64.min = limit[int64]{value: 1, set: true}; return typ }(),
		func() Type {
			typ := Float32().Type()
			typ.float64.min = limit[float64]{value: 1, set: true}
			return typ
		}(),
		func() Type { typ := Object().Type(); typ.ref.name = "example.Name"; return typ }(),
		func() Type {
			typ := Ref("example.Name").Type()
			typ.object.fields = []FieldDescriptor{Field("name").String().Required().Field()}
			return typ
		}(),
	}
	for _, typ := range tests {
		requireInvalidType(t, typ, nil, ErrInvalidType)
	}
}

func TestValidateInactivePayloadsAcceptsExactTemporalSlots(t *testing.T) {
	requireValidType(t, Timestamp().Type(), nil)
	requireValidType(t, Date().Type(), nil)
	requireValidType(t, Time().Type(), nil)
	requireValidType(t, Duration().Type(), nil)
}
