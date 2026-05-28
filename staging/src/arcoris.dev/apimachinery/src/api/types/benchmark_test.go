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

func BenchmarkStringTypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = String().MinLen(1).MaxLen(253).Pattern("^[a-z]+$").Enum("alpha", "beta").Type()
	}
}

func BenchmarkInt8TypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Int8().Range(0, 10).Enum(1, 2, 3).Type()
	}
}

func BenchmarkInt64TypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Int64().Range(1, 1000).Enum(1, 10, 100).Type()
	}
}

func BenchmarkUint64TypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Uint64().Range(0, 1000).Enum(1, 10, 100).Type()
	}
}

func BenchmarkFloat64TypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Float64().Range(0, 1).Enum(0.25, 0.5, 0.75).Type()
	}
}

func BenchmarkObjectTypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Object(
			Field("name").String().Required().MinLen(1),
			Field("replicas").Int64().Optional().Min(1),
			Field("labels").MapOf(String()).Optional(),
		).UnknownFields(UnknownReject).Type()
	}
}

func BenchmarkListTypeBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ListOf(Object(
			Field("name").String().Required().MinLen(1),
			Field("value").Int64().Required().Min(0),
		)).Map("name").Type()
	}
}

func BenchmarkValidateSmallObject(b *testing.B) {
	tp := Object(
		Field("name").String().Required().MinLen(1),
		Field("replicas").Int64().Optional().Min(1),
	).Type()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := ValidateType(tp, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateNestedObject(b *testing.B) {
	tp := Object(
		Field("spec").Object(
			Field("name").String().Required().MinLen(1),
			Field("ports").ListOf(Object(
				Field("name").String().Required().MinLen(1),
				Field("port").Uint16().Required().Min(1),
			)).Optional().Map("name"),
		).Required().UnknownFields(UnknownReject),
	).Type()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := ValidateType(tp, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateRefWithResolver(b *testing.B) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		if name == "example.Name" {
			return Define("example.Name", String().MinLen(1)), true
		}
		return TypeDefinition{}, false
	})
	tp := Ref("example.Name").Type()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := ValidateType(tp, resolver); err != nil {
			b.Fatal(err)
		}
	}
}
