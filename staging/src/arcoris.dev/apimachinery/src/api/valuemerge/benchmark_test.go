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
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func BenchmarkMergeScalarExact(b *testing.B) {
	benchmarkMerge(
		b,
		str("old"),
		str("new"),
		types.String().Descriptor(),
		pathSet(root()),
	)
}

func BenchmarkMergeRecordSmall(b *testing.B) {
	base, overlay, descriptor := benchmarkRecordMergeInputs(8)

	benchmarkMerge(
		b,
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("f003"))),
	)
}

func BenchmarkMergeRecordLargeSelectedFew(b *testing.B) {
	base, overlay, descriptor := benchmarkRecordMergeInputs(100)

	benchmarkMerge(
		b,
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("f010")), root().Field(testFieldName("f090"))),
	)
}

func BenchmarkMergeRecordLargeSelectedMany(b *testing.B) {
	base, overlay, descriptor := benchmarkRecordMergeInputs(100)
	paths := make([]fieldpath.Path, 0, 100)
	for i := 0; i < 100; i++ {
		paths = append(paths, root().Field(testFieldName(fmt.Sprintf("f%03d", i))))
	}

	benchmarkMerge(b, base, overlay, descriptor, pathSet(paths...))
}

func BenchmarkMergeMapLargeSelectedFew(b *testing.B) {
	base := benchmarkStringRecord(100, "old")
	overlay := benchmarkStringRecord(100, "new")
	descriptor := types.MapOf(types.String()).Descriptor()

	benchmarkMerge(
		b,
		base,
		overlay,
		descriptor,
		pathSet(root().Key(testMapKey("k010")), root().Key(testMapKey("k090"))),
	)
}

func BenchmarkMergeOrderedListLargeSelectedFew(b *testing.B) {
	base := benchmarkStringList(100, "old")
	overlay := benchmarkStringList(100, "new")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	benchmarkMerge(
		b,
		base,
		overlay,
		descriptor,
		pathSet(root().Index(10), root().Index(90)),
	)
}

func BenchmarkMergeListMapLargeSelectedFew(b *testing.B) {
	base := benchmarkConditionList(100, "False")
	overlay := benchmarkConditionList(100, "True")

	benchmarkMerge(
		b,
		base,
		overlay,
		conditionListDescriptor(),
		pathSet(
			benchmarkConditionPath("Ready-010").Field(testFieldName("status")),
			benchmarkConditionPath("Ready-090").Field(testFieldName("status")),
		),
	)
}

func BenchmarkMergeListMapLargeSelectedMany(b *testing.B) {
	base := benchmarkConditionList(100, "False")
	overlay := benchmarkConditionList(100, "True")
	paths := make([]fieldpath.Path, 0, 100)
	for i := 0; i < 100; i++ {
		paths = append(paths, benchmarkConditionPath(fmt.Sprintf("Ready-%03d", i)).Field(testFieldName("status")))
	}

	benchmarkMerge(b, base, overlay, conditionListDescriptor(), pathSet(paths...))
}

func benchmarkMerge(
	b *testing.B,
	base value.Value,
	overlay value.Value,
	descriptor types.Descriptor,
	fields fieldpath.Set,
) {
	b.Helper()

	for i := 0; i < b.N; i++ {
		if _, err := Merge(base, overlay, descriptor, fields, Options{}); err != nil {
			b.Fatalf("Merge returned error: %v", err)
		}
	}
}

func benchmarkRecordMergeInputs(count int) (value.Value, value.Value, types.Descriptor) {
	baseMembers := make([]value.RecordMember, 0, count)
	overlayMembers := make([]value.RecordMember, 0, count)
	fields := make([]types.FieldExpr, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("f%03d", i)
		baseMembers = append(baseMembers, member(name, str("old")))
		overlayMembers = append(overlayMembers, member(name, str("new")))
		fields = append(fields, types.Field(name).String().Optional())
	}

	return obj(baseMembers...), obj(overlayMembers...), types.Object(fields...).Descriptor()
}

func benchmarkStringRecord(count int, text string) value.Value {
	members := make([]value.RecordMember, 0, count)
	for i := 0; i < count; i++ {
		members = append(members, member(fmt.Sprintf("k%03d", i), str(text)))
	}

	return obj(members...)
}

func benchmarkStringList(count int, text string) value.Value {
	items := make([]value.Value, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, str(text))
	}

	return list(items...)
}

func benchmarkConditionList(count int, status string) value.Value {
	items := make([]value.Value, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, conditionItem(fmt.Sprintf("Ready-%03d", i), status))
	}

	return list(items...)
}

func benchmarkConditionPath(conditionType string) fieldpath.Path {
	return root().Select(
		fieldpath.MustSelector(
			fieldpath.NewSelectorEntry(testFieldName("type"), fieldpath.StringLiteral(conditionType)),
		),
	)
}
