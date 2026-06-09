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

var benchmarkDescriptorSink Descriptor

func BenchmarkStringDescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = String().
			MinBytes(1).
			MaxBytes(253).
			Pattern("^[a-z]+$").
			Enum("alpha", "beta").
			Descriptor()
	}
}

func BenchmarkInt8DescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = Int8().
			Range(0, 10).
			Enum(1, 2, 3).
			Descriptor()
	}
}

func BenchmarkInt64DescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = Int64().
			Range(1, 1000).
			Enum(1, 10, 100).
			Descriptor()
	}
}

func BenchmarkUint64DescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = Uint64().
			Range(0, 1000).
			Enum(1, 10, 100).
			Descriptor()
	}
}

func BenchmarkFloat64DescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = Float64().
			Range(0, 1).
			Enum(0.25, 0.5, 0.75).
			Descriptor()
	}
}

func BenchmarkObjectDescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = Object(
			Field("name").String().Required().MinBytes(1),
			Field("replicas").Int64().Optional().Min(1),
			Field("labels").MapOf(String()).Optional(),
		).UnknownFields(UnknownReject).Descriptor()
	}
}

func BenchmarkListDescriptorBuild(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDescriptorSink = ListOf(Object(
			Field("name").String().Required().MinBytes(1),
			Field("value").Int64().Required().Min(0),
		)).Map("name").Descriptor()
	}
}

func BenchmarkValidateSmallObject(b *testing.B) {
	desc := Object(
		Field("name").String().Required().MinBytes(1),
		Field("replicas").Int64().Optional().Min(1),
	).Descriptor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := ValidateLocal(desc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateNestedObject(b *testing.B) {
	desc := Object(
		Field("spec").Object(
			Field("name").String().Required().MinBytes(1),
			Field("ports").ListOf(Object(
				Field("name").String().Required().MinBytes(1),
				Field("port").Uint16().Required().Min(1),
			)).Optional().Map("name"),
		).Required().UnknownFields(UnknownReject),
	).Descriptor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := ValidateLocal(desc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateRefWithResolver(b *testing.B) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		if name == "example.Name" {
			return Define("example.Name", String().MinBytes(1)), true
		}
		return Definition{}, false
	})
	desc := Ref("example.Name").Descriptor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := ValidateResolved(desc, resolver); err != nil {
			b.Fatal(err)
		}
	}
}
