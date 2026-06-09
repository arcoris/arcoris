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
	tests := []Descriptor{
		func() Descriptor {
			desc := String().Descriptor()
			desc.bytes.minBytes = limit[int]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.int8.min = limit[int8]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.int16.min = limit[int16]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.int32.min = limit[int32]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.int64.min = limit[int64]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.uint8.min = limit[uint8]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.uint16.min = limit[uint16]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.uint32.min = limit[uint32]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.uint64.min = limit[uint64]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.float32.min = limit[float32]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.float64.min = limit[float64]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.decimal.precision = limit[int]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			desc.object.fields = []FieldDescriptor{Field("name").String().Required().Field()}
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			elem := Bool().Descriptor()
			desc.list.elem = &elem
			return desc
		}(),
		func() Descriptor {
			desc := String().Descriptor()
			value := Bool().Descriptor()
			desc.mapType.value = &value
			return desc
		}(),
		func() Descriptor { desc := String().Descriptor(); desc.ref.name = "example.Name"; return desc }(),
		func() Descriptor {
			desc := Int8().Descriptor()
			desc.int16.min = limit[int16]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := Int8().Descriptor()
			desc.string.minBytes = limit[int]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := Uint64().Descriptor()
			desc.int64.min = limit[int64]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor {
			desc := Float32().Descriptor()
			desc.float64.min = limit[float64]{value: 1, set: true}
			return desc
		}(),
		func() Descriptor { desc := Object().Descriptor(); desc.ref.name = "example.Name"; return desc }(),
		func() Descriptor {
			desc := Ref("example.Name").Descriptor()
			desc.object.fields = []FieldDescriptor{Field("name").String().Required().Field()}
			return desc
		}(),
	}
	for _, desc := range tests {
		requireInvalidDescriptor(t, desc, nil, ErrInvalidDescriptor)
	}
}

func TestValidateInactivePayloadsAcceptsExactTemporalSlots(t *testing.T) {
	requireValidDescriptor(t, Timestamp().Descriptor(), nil)
	requireValidDescriptor(t, Date().Descriptor(), nil)
	requireValidDescriptor(t, Time().Descriptor(), nil)
	requireValidDescriptor(t, Duration().Descriptor(), nil)
}
