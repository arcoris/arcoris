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

package resource

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

var (
	benchmarkDefinitionSink Definition
	benchmarkScopeSink      Scope
	benchmarkBytesSink      []byte
)

func BenchmarkDefinitionBuildSmall(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkDefinitionSink = NewDefinition(
			identity.Group("control.arcoris.dev"),
			identity.Kind("Worker"),
			identity.Resource("workers"),
			ScopeNamespaced,
			NewVersion(
				identity.Version("v1"),
				types.Object().Descriptor(),
				Observed(types.Object().Descriptor()),
				Exposed(),
				Canonical(),
			),
		)
	}
}

func BenchmarkNewDefinitionPreparedVersion(b *testing.B) {
	version := validVersion()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkDefinitionSink = NewDefinition(
			identity.Group("control.arcoris.dev"),
			identity.Kind("Worker"),
			identity.Resource("workers"),
			ScopeNamespaced,
			version,
		)
	}
}

func BenchmarkValidateDefinition(b *testing.B) {
	def := validDefinition()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := ValidateDefinitionLocal(def); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkScopeMarshalJSON(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data, err := ScopeNamespaced.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytesSink = data
	}
}

func BenchmarkScopeUnmarshalJSON(b *testing.B) {
	data, err := ScopeNamespaced.MarshalJSON()
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var scope Scope
		if err := scope.UnmarshalJSON(data); err != nil {
			b.Fatal(err)
		}
		benchmarkScopeSink = scope
	}
}
