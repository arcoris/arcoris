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

package fieldpath

import "testing"

var (
	benchmarkPathSink Path
	benchmarkSetSink  Set
	benchmarkBoolSink bool
)

func BenchmarkParseCanonicalSelectorPath(b *testing.B) {
	text := `$.routes[{"host":"api.example.com","port":443}].backend`

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		benchmarkPathSink, err = ParseCanonical(text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPathAppend(b *testing.B) {
	base := Root().Field(MustFieldName("spec")).Field(MustFieldName("template"))
	element := MustFieldElement("containers")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkPathSink = base.Append(element)
	}
}

func BenchmarkNewSet(b *testing.B) {
	paths := []Path{
		Root().Field(MustFieldName("status")),
		Root().Field(MustFieldName("spec")).Field(MustFieldName("replicas")),
		Root().Field(MustFieldName("metadata")).Field(MustFieldName("labels")).Key(MustMapKey("app")),
		Root().Field(MustFieldName("spec")),
		Root().Field(MustFieldName("status")),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		benchmarkSetSink, err = NewSet(paths...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetHas(b *testing.B) {
	set := MustSet(
		Root().Field(MustFieldName("metadata")).Field(MustFieldName("labels")).Key(MustMapKey("app")),
		Root().Field(MustFieldName("spec")),
		Root().Field(MustFieldName("spec")).Field(MustFieldName("replicas")),
		Root().Field(MustFieldName("status")),
	)
	target := Root().Field(MustFieldName("spec")).Field(MustFieldName("replicas"))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkBoolSink = set.Has(target)
	}
}
