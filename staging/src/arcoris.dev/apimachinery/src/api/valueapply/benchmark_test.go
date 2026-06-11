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

package valueapply

import (
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func BenchmarkApplyScalarNoConflict(b *testing.B) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       str("old"),
		Applied:    str("new"),
		Descriptor: types.String().Descriptor(),
		Ownership:  state(),
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyRecordSmall(b *testing.B) {
	req := specRequest(owner("user"))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyRecordLarge(b *testing.B) {
	req := benchmarkLargeRecordRequest(24)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyMapLarge(b *testing.B) {
	req := benchmarkLargeMapRequest(24)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyListMapLarge(b *testing.B) {
	req := benchmarkLargeListMapRequest(24)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyConflictDetectionLargeOwnership(b *testing.B) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func BenchmarkApplyDroppedFieldsLargeOwnership(b *testing.B) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("user", imagePath(), replicasPath(), path("$.extra")))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Apply(req, Options{})
	}
}

func benchmarkLargeRecordRequest(size int) Request {
	fields := make([]types.FieldExpr, 0, size)
	liveMembers := make([]value.RecordMember, 0, size)
	appliedMembers := make([]value.RecordMember, 0, size/2+1)

	for i := 0; i < size; i++ {
		name := fmt.Sprintf("field%02d", i)
		fields = append(fields, types.Field(name).String().Optional())
		liveMembers = append(liveMembers, member(name, str(fmt.Sprintf("old-%02d", i))))
		if i%2 == 0 {
			appliedMembers = append(appliedMembers, member(name, str(fmt.Sprintf("new-%02d", i))))
		}
	}

	return Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(liveMembers...),
		Applied:    obj(appliedMembers...),
		Descriptor: types.Object(fields...).Descriptor(),
		Ownership:  state(),
	}
}

func benchmarkLargeMapRequest(size int) Request {
	fields := make([]value.RecordMember, 0, size)
	for i := 0; i < size; i++ {
		name := fmt.Sprintf("key%02d", i)
		fields = append(fields, member(name, str(fmt.Sprintf("new-%02d", i))))
	}

	live := make([]value.RecordMember, 0, size)
	for i := 0; i < size; i++ {
		name := fmt.Sprintf("key%02d", i)
		live = append(live, member(name, str(fmt.Sprintf("old-%02d", i))))
	}

	return Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(live...),
		Applied:    obj(fields...),
		Descriptor: mapDescriptor(),
		Ownership:  state(),
	}
}

func benchmarkLargeListMapRequest(size int) Request {
	live := make([]value.Value, 0, size)
	applied := make([]value.Value, 0, size)
	for i := 0; i < size; i++ {
		live = append(live, benchmarkCondition(fmt.Sprintf("T%02d", i), "old"))
		applied = append(applied, benchmarkCondition(fmt.Sprintf("T%02d", i), "new"))
	}

	return Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(live...),
		Applied:    list(applied...),
		Descriptor: benchmarkConditionListDescriptor(),
		Ownership:  state(),
	}
}

func benchmarkCondition(status string, value string) value.Value {
	return obj(member("type", str(status)), member("status", str(value)))
}

func benchmarkConditionListDescriptor() types.Descriptor {
	return types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().Optional(),
		),
	).Map("type").Descriptor()
}
