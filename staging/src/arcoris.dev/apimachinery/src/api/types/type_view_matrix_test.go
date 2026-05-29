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

func TestTypeExactViewWrongCodeMatrix(t *testing.T) {
	tests := []struct {
		name  string
		typ   Type
		exact func(Type) bool
		wrong []func(Type) bool
	}{
		{
			name:  "string",
			typ:   String().Type(),
			exact: func(typ Type) bool { _, ok := typ.String(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Bytes(); return ok },
				func(typ Type) bool { _, ok := typ.Object(); return ok },
				func(typ Type) bool { _, ok := typ.Int8(); return ok },
				func(typ Type) bool { _, ok := typ.Ref(); return ok },
			},
		},
		{
			name:  "bytes",
			typ:   Bytes().Type(),
			exact: func(typ Type) bool { _, ok := typ.Bytes(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.String(); return ok },
				func(typ Type) bool { _, ok := typ.Object(); return ok },
				func(typ Type) bool { _, ok := typ.Uint8(); return ok },
			},
		},
		{
			name:  "int8",
			typ:   Int8().Type(),
			exact: func(typ Type) bool { _, ok := typ.Int8(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Int16(); return ok },
				func(typ Type) bool { _, ok := typ.Int64(); return ok },
				func(typ Type) bool { _, ok := typ.Uint8(); return ok },
				func(typ Type) bool { _, ok := typ.Float32(); return ok },
				func(typ Type) bool { _, ok := typ.String(); return ok },
			},
		},
		{
			name:  "int64",
			typ:   Int64().Type(),
			exact: func(typ Type) bool { _, ok := typ.Int64(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Int8(); return ok },
				func(typ Type) bool { _, ok := typ.Uint64(); return ok },
				func(typ Type) bool { _, ok := typ.Float64(); return ok },
				func(typ Type) bool { _, ok := typ.Object(); return ok },
			},
		},
		{
			name:  "uint8",
			typ:   Uint8().Type(),
			exact: func(typ Type) bool { _, ok := typ.Uint8(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Uint16(); return ok },
				func(typ Type) bool { _, ok := typ.Int8(); return ok },
				func(typ Type) bool { _, ok := typ.Float32(); return ok },
				func(typ Type) bool { _, ok := typ.String(); return ok },
			},
		},
		{
			name:  "uint64",
			typ:   Uint64().Type(),
			exact: func(typ Type) bool { _, ok := typ.Uint64(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Uint32(); return ok },
				func(typ Type) bool { _, ok := typ.Int64(); return ok },
				func(typ Type) bool { _, ok := typ.Float64(); return ok },
				func(typ Type) bool { _, ok := typ.Map(); return ok },
			},
		},
		{
			name:  "float32",
			typ:   Float32().Type(),
			exact: func(typ Type) bool { _, ok := typ.Float32(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Float64(); return ok },
				func(typ Type) bool { _, ok := typ.Int32(); return ok },
				func(typ Type) bool { _, ok := typ.Decimal(); return ok },
			},
		},
		{
			name:  "float64",
			typ:   Float64().Type(),
			exact: func(typ Type) bool { _, ok := typ.Float64(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Float32(); return ok },
				func(typ Type) bool { _, ok := typ.Uint64(); return ok },
				func(typ Type) bool { _, ok := typ.Decimal(); return ok },
			},
		},
		{
			name:  "decimal",
			typ:   Decimal().Type(),
			exact: func(typ Type) bool { _, ok := typ.Decimal(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Float64(); return ok },
				func(typ Type) bool { _, ok := typ.Int64(); return ok },
			},
		},
		{
			name:  "timestamp",
			typ:   Timestamp().Type(),
			exact: func(typ Type) bool { _, ok := typ.Timestamp(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Date(); return ok },
				func(typ Type) bool { _, ok := typ.Time(); return ok },
				func(typ Type) bool { _, ok := typ.Duration(); return ok },
			},
		},
		{
			name:  "date",
			typ:   Date().Type(),
			exact: func(typ Type) bool { _, ok := typ.Date(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Timestamp(); return ok },
				func(typ Type) bool { _, ok := typ.Time(); return ok },
				func(typ Type) bool { _, ok := typ.Duration(); return ok },
			},
		},
		{
			name:  "time",
			typ:   Time().Type(),
			exact: func(typ Type) bool { _, ok := typ.Time(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Timestamp(); return ok },
				func(typ Type) bool { _, ok := typ.Date(); return ok },
				func(typ Type) bool { _, ok := typ.Duration(); return ok },
			},
		},
		{
			name:  "duration",
			typ:   Duration().Type(),
			exact: func(typ Type) bool { _, ok := typ.Duration(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Timestamp(); return ok },
				func(typ Type) bool { _, ok := typ.Date(); return ok },
				func(typ Type) bool { _, ok := typ.Time(); return ok },
			},
		},
		{
			name:  "object",
			typ:   Object().Type(),
			exact: func(typ Type) bool { _, ok := typ.Object(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.List(); return ok },
				func(typ Type) bool { _, ok := typ.Map(); return ok },
				func(typ Type) bool { _, ok := typ.Ref(); return ok },
				func(typ Type) bool { _, ok := typ.String(); return ok },
			},
		},
		{
			name:  "list",
			typ:   ListOf(String()).Type(),
			exact: func(typ Type) bool { _, ok := typ.List(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Object(); return ok },
				func(typ Type) bool { _, ok := typ.Map(); return ok },
				func(typ Type) bool { _, ok := typ.Ref(); return ok },
			},
		},
		{
			name:  "map",
			typ:   MapOf(String()).Type(),
			exact: func(typ Type) bool { _, ok := typ.Map(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Object(); return ok },
				func(typ Type) bool { _, ok := typ.List(); return ok },
				func(typ Type) bool { _, ok := typ.Ref(); return ok },
			},
		},
		{
			name:  "ref",
			typ:   Ref("example.Name").Type(),
			exact: func(typ Type) bool { _, ok := typ.Ref(); return ok },
			wrong: []func(Type) bool{
				func(typ Type) bool { _, ok := typ.Object(); return ok },
				func(typ Type) bool { _, ok := typ.List(); return ok },
				func(typ Type) bool { _, ok := typ.Map(); return ok },
				func(typ Type) bool { _, ok := typ.String(); return ok },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.exact(tt.typ), true)
			for _, wrong := range tt.wrong {
				requireEqual(t, wrong(tt.typ), false)
			}
		})
	}
}

func TestTypeViewWrongCodeReturnsZeroView(t *testing.T) {
	typ := String().Type()
	typ.int8.enum = []int8{1}
	int8View, ok := typ.Int8()
	requireEqual(t, ok, false)
	requireEqual(t, len(int8View.Enum()), 0)

	elem := String().Enum("nested").Type()
	typ.list.elem = &elem
	typ.list.mapKeys = []FieldName{"name"}
	listView, ok := typ.List()
	requireEqual(t, ok, false)
	requireEqual(t, listView.Element().IsZero(), true)
	requireEqual(t, len(listView.MapKeys()), 0)

	value := String().Enum("map").Type()
	typ.mapType.value = &value
	mapView, ok := typ.Map()
	requireEqual(t, ok, false)
	requireEqual(t, mapView.Value().IsZero(), true)

	typ.ref.name = "example.Name"
	refView, ok := typ.Ref()
	requireEqual(t, ok, false)
	requireEqual(t, refView.Name(), TypeName(""))
}

func TestTypeViewDetachedNestedDescriptors(t *testing.T) {
	list := ListOf(String().Enum("a")).Map("name").Type()
	elem := requireListView(t, list).Element()
	elem.string.enum[0] = "b"
	requireEqual(t, requireStringView(t, requireListView(t, list).Element()).Enum()[0], "a")

	keys := requireListView(t, list).MapKeys()
	keys[0] = "changed"
	requireEqual(t, requireListView(t, list).MapKeys()[0], FieldName("name"))

	mapping := MapOf(String().Enum("a")).Type()
	value := requireMapView(t, mapping).Value()
	value.string.enum[0] = "b"
	requireEqual(t, requireStringView(t, requireMapView(t, mapping).Value()).Enum()[0], "a")
}
