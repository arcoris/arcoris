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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

type resolverFunc func(types.TypeName) (types.Definition, bool)

func (f resolverFunc) Resolve(name types.TypeName) (types.Definition, bool) {
	return f(name)
}

func pathSet(paths ...fieldpath.Path) fieldpath.Set {
	return fieldpath.MustSet(paths...)
}

func root() fieldpath.Path {
	return fieldpath.Root()
}

func testFieldName(name string) fieldpath.FieldName {
	return fieldpath.MustFieldName(name)
}

func testMapKey(key string) fieldpath.MapKey {
	return fieldpath.MustMapKey(key)
}

func str(text string) value.Value {
	return value.StringValue(text)
}

func boolValue(v bool) value.Value {
	return value.BoolValue(v)
}

func intValue(v int64) value.Value {
	return value.IntegerValue(value.NewIntegerFromInt64(v))
}

func decimalValue(text string) value.Value {
	decimal, err := value.ParseDecimal(text)
	if err != nil {
		panic(err)
	}

	return value.DecimalValue(decimal)
}

func obj(members ...value.RecordMember) value.Value {
	return value.MustRecordValue(members...)
}

func list(items ...value.Value) value.Value {
	return value.MustListValue(items...)
}

func member(name string, v value.Value) value.RecordMember {
	return value.MustRecordMember(name, v)
}
