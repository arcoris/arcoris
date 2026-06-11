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

package valuefieldset

import (
	"strconv"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

var benchmarkSet fieldpath.Set

func BenchmarkExtractOwnershipFieldsScalar(b *testing.B) {
	benchmarkExtract(b, value.StringValue("api"), types.String().Descriptor())
}

func BenchmarkExtractOwnershipFieldsObjectSmall(b *testing.B) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
		types.Field("image").String().Optional(),
		types.Field("replicas").Int64().Optional(),
	).Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("name", value.StringValue("api")),
		value.MustRecordMember("image", value.StringValue("api:v1")),
		value.MustRecordMember("replicas", value.Int64Value(3)),
	)

	benchmarkExtract(b, val, descriptor)
}

func BenchmarkExtractOwnershipFieldsObjectLarge(b *testing.B) {
	val, descriptor := benchmarkObject(100)
	benchmarkExtract(b, val, descriptor)
}

func BenchmarkExtractOwnershipFieldsMapLarge(b *testing.B) {
	val := benchmarkRecord("key", 100)
	benchmarkExtract(b, val, types.MapOf(types.String()).Descriptor())
}

func BenchmarkExtractOwnershipFieldsOrderedListLarge(b *testing.B) {
	val := benchmarkList(100)
	benchmarkExtract(b, val, types.ListOf(types.String()).Ordered().Descriptor())
}

func BenchmarkExtractOwnershipFieldsListMapLarge(b *testing.B) {
	descriptor := types.ListOf(types.Object(
		types.Field("type").String().Required(),
		types.Field("status").String().Optional(),
	)).Map("type").Descriptor()
	items := make([]value.Value, 100)
	for i := range items {
		items[i] = value.MustRecordValue(
			value.MustRecordMember("type", value.StringValue("Type"+strconv.Itoa(i))),
			value.MustRecordMember("status", value.StringValue("True")),
		)
	}
	val := value.MustListValue(items...)

	benchmarkExtract(b, val, descriptor)
}

func BenchmarkExtractOwnershipFieldsUnknownPreserveOpaque(b *testing.B) {
	val := benchmarkRecord("unknown", 100)
	descriptor := types.Object().
		UnknownFields(types.UnknownPreserveOpaque).
		Descriptor()

	benchmarkExtract(b, val, descriptor)
}

func BenchmarkExtractOwnershipFieldsDeepNested(b *testing.B) {
	descriptor := types.Object(
		types.Field("a").Object(
			types.Field("b").Object(
				types.Field("c").Object(
					types.Field("d").String().Optional(),
				).Optional(),
			).Optional(),
		).Optional(),
	).Descriptor()
	val := value.MustRecordValue(value.MustRecordMember("a",
		value.MustRecordValue(value.MustRecordMember("b",
			value.MustRecordValue(value.MustRecordMember("c",
				value.MustRecordValue(value.MustRecordMember("d", value.StringValue("value"))),
			)),
		)),
	))

	benchmarkExtract(b, val, descriptor)
}

func benchmarkExtract(b *testing.B, val value.Value, descriptor types.Descriptor) {
	b.Helper()
	b.ReportAllocs()

	var got fieldpath.Set
	for i := 0; i < b.N; i++ {
		set, err := ExtractOwnershipFields(val, descriptor, Options{})
		if err != nil {
			b.Fatalf("ExtractOwnershipFields returned error: %v", err)
		}
		got = set
	}
	benchmarkSet = got
}

func benchmarkObject(size int) (value.Value, types.Descriptor) {
	fields := make([]types.FieldExpr, size)
	members := make([]value.RecordMember, size)
	for i := 0; i < size; i++ {
		name := "field" + strconv.Itoa(i)
		fields[i] = types.Field(name).String().Optional()
		members[i] = value.MustRecordMember(name, value.StringValue("value"))
	}

	return value.MustRecordValue(members...), types.Object(fields...).Descriptor()
}

func benchmarkRecord(prefix string, size int) value.Value {
	members := make([]value.RecordMember, size)
	for i := 0; i < size; i++ {
		members[i] = value.MustRecordMember(prefix+strconv.Itoa(i), value.StringValue("value"))
	}

	return value.MustRecordValue(members...)
}

func benchmarkList(size int) value.Value {
	items := make([]value.Value, size)
	for i := 0; i < size; i++ {
		items[i] = value.StringValue("value")
	}

	return value.MustListValue(items...)
}
