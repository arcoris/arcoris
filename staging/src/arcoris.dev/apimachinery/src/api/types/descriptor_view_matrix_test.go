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

func TestDescriptorExactViewWrongKindMatrix(t *testing.T) {
	tests := []struct {
		name       string
		descriptor Descriptor
		exact      func(Descriptor) bool
		wrong      []func(Descriptor) bool
	}{
		{
			name:       "string",
			descriptor: String().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsBytes(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt8(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsRef(); return ok },
			},
		},
		{
			name:       "bytes",
			descriptor: Bytes().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsBytes(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsUint8(); return ok },
			},
		},
		{
			name:       "int8",
			descriptor: Int8().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsInt8(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsInt16(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsUint8(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsFloat32(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
			},
		},
		{
			name:       "int64",
			descriptor: Int64().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsInt64(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsInt8(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsUint64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
			},
		},
		{
			name:       "uint8",
			descriptor: Uint8().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsUint8(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsUint16(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt8(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsFloat32(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
			},
		},
		{
			name:       "uint64",
			descriptor: Uint64().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsUint64(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsUint32(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsMap(); return ok },
			},
		},
		{
			name:       "float32",
			descriptor: Float32().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsFloat32(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt32(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDecimal(); return ok },
			},
		},
		{
			name:       "float64",
			descriptor: Float64().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsFloat32(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsUint64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDecimal(); return ok },
			},
		},
		{
			name:       "decimal",
			descriptor: Decimal().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsDecimal(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsInt64(); return ok },
			},
		},
		{
			name:       "timestamp",
			descriptor: Timestamp().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsTimestamp(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsDate(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsTime(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDuration(); return ok },
			},
		},
		{
			name:       "date",
			descriptor: Date().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsDate(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsTimestamp(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsTime(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDuration(); return ok },
			},
		},
		{
			name:       "time",
			descriptor: Time().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsTime(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsTimestamp(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDate(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDuration(); return ok },
			},
		},
		{
			name:       "duration",
			descriptor: Duration().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsDuration(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsTimestamp(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsDate(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsTime(); return ok },
			},
		},
		{
			name:       "object",
			descriptor: Object().Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsList(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsMap(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsRef(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
			},
		},
		{
			name:       "list",
			descriptor: ListOf(String()).Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsList(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsMap(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsRef(); return ok },
			},
		},
		{
			name:       "map",
			descriptor: MapOf(String()).Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsMap(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsList(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsRef(); return ok },
			},
		},
		{
			name:       "ref",
			descriptor: Ref("example.Name").Descriptor(),
			exact:      func(desc Descriptor) bool { _, ok := desc.AsRef(); return ok },
			wrong: []func(Descriptor) bool{
				func(desc Descriptor) bool { _, ok := desc.AsObject(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsList(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsMap(); return ok },
				func(desc Descriptor) bool { _, ok := desc.AsString(); return ok },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.exact(tt.descriptor), true)
			for _, wrong := range tt.wrong {
				requireEqual(t, wrong(tt.descriptor), false)
			}
		})
	}
}

func TestDescriptorViewWrongKindReturnsZeroView(t *testing.T) {
	desc := String().Descriptor()
	desc.int8.enum = []int8{1}
	int8View, ok := desc.AsInt8()
	requireEqual(t, ok, false)
	requireEqual(t, len(int8View.Enum()), 0)

	elem := String().Enum("nested").Descriptor()
	desc.list.elem = &elem
	desc.list.mapKeys = []FieldName{"name"}
	listView, ok := desc.AsList()
	requireEqual(t, ok, false)
	requireEqual(t, listView.Element().IsZero(), true)
	requireEqual(t, len(listView.MapKeys()), 0)

	value := String().Enum("map").Descriptor()
	desc.mapType.value = &value
	mapView, ok := desc.AsMap()
	requireEqual(t, ok, false)
	requireEqual(t, mapView.Value().IsZero(), true)

	desc.ref.name = "example.Name"
	refView, ok := desc.AsRef()
	requireEqual(t, ok, false)
	requireEqual(t, refView.Name(), TypeName(""))
}

func TestDescriptorViewDetachedNestedDescriptors(t *testing.T) {
	list := ListOf(String().Enum("a")).Map("name").Descriptor()
	elem := requireListView(t, list).Element()
	elem.string.enum[0] = "b"
	requireEqual(t, requireStringView(t, requireListView(t, list).Element()).Enum()[0], "a")

	keys := requireListView(t, list).MapKeys()
	keys[0] = "changed"
	requireEqual(t, requireListView(t, list).MapKeys()[0], FieldName("name"))

	mapping := MapOf(String().Enum("a")).Descriptor()
	value := requireMapView(t, mapping).Value()
	value.string.enum[0] = "b"
	requireEqual(t, requireStringView(t, requireMapView(t, mapping).Value()).Enum()[0], "a")
}
