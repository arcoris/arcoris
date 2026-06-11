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

package valuecompare

import (
	"strconv"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func BenchmarkCompareScalar(b *testing.B) {
	descriptor := types.String().Descriptor()
	oldValue := value.StringValue("old")
	newValue := value.StringValue("new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareRecordSmall(b *testing.B) {
	descriptor := typesObject("name", "image")
	oldValue := valueRecord("name", "app", "image", "v1")
	newValue := valueRecord("name", "app", "image", "v2")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareRecordLarge(b *testing.B) {
	descriptor := largeObjectDescriptor(24)
	oldValue := largeRecordValue(24, "old")
	newValue := largeRecordValue(24, "new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareMapLarge(b *testing.B) {
	descriptor := types.MapOf(types.String()).Descriptor()
	oldValue := largeMapValue(32, "old")
	newValue := largeMapValue(32, "new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareOrderedListLarge(b *testing.B) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()
	oldValue := largeListValue(32, "old")
	newValue := largeListValue(32, "new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareListMapLarge(b *testing.B) {
	descriptor := types.ListOf(conditionExpr()).Map("type").Descriptor()
	oldValue := largeConditionList("False")
	newValue := largeConditionList("True")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareAddedSubtreeLarge(b *testing.B) {
	descriptor := largeObjectDescriptor(24)
	newValue := largeRecordValue(24, "new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(valueRecord(), newValue, descriptor, Options{})
	}
}

func BenchmarkCompareRemovedSubtreeLarge(b *testing.B) {
	descriptor := largeObjectDescriptor(24)
	oldValue := largeRecordValue(24, "old")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, valueRecord(), descriptor, Options{})
	}
}

func BenchmarkCompareUnknownPreserveOpaque(b *testing.B) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("extra", largeRecordValue(12, "old")))
	newValue := value.MustRecordValue(value.MustRecordMember("extra", largeRecordValue(12, "new")))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareAtomicList(b *testing.B) {
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()
	oldValue := value.MustListValue(value.StringValue("old"))
	newValue := value.MustListValue(value.StringValue("new"))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func BenchmarkCompareListSet(b *testing.B) {
	descriptor := types.ListOf(types.String()).Set().Descriptor()
	oldValue := largeListValue(16, "old")
	newValue := largeListValue(16, "new")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Compare(oldValue, newValue, descriptor, Options{})
	}
}

func largeObjectDescriptor(fields int) types.Descriptor {
	names := make([]string, 0, fields)
	for i := 0; i < fields; i++ {
		names = append(names, benchFieldName(i))
	}

	return typesObject(names...)
}

func largeRecordValue(fields int, suffix string) value.Value {
	members := make([]string, 0, fields*2)
	for i := 0; i < fields; i++ {
		name := benchFieldName(i)
		members = append(members, name, name+"-"+suffix)
	}

	return valueRecord(members...)
}

func largeMapValue(fields int, suffix string) value.Value {
	return largeRecordValue(fields, suffix)
}

func largeListValue(items int, suffix string) value.Value {
	values := make([]value.Value, 0, items)
	for i := 0; i < items; i++ {
		values = append(values, value.StringValue(benchFieldName(i)+"-"+suffix))
	}

	return value.MustListValue(values...)
}

func largeConditionList(status string) value.Value {
	items := make([]value.Value, 0, 16)
	for i := 0; i < 16; i++ {
		items = append(items, conditionValue(benchFieldName(i), status))
	}

	return value.MustListValue(items...)
}

func benchFieldName(i int) string {
	return "field-" + strconv.Itoa(i)
}
