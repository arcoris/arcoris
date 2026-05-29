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

package resourcecatalog

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

var (
	benchmarkDefinitionSink        resource.Definition
	benchmarkVersionDefinitionSink resource.VersionDefinition
	benchmarkResourcesSink         []identity.GroupResource
	benchmarkCatalogSizeSink       int
)

func BenchmarkRegisterOne(b *testing.B) {
	def := validDefinition("Worker", "workers")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var catalog Catalog
		if err := catalog.Register(def); err != nil {
			b.Fatal(err)
		}
		benchmarkCatalogSizeSink = len(catalog.defsByResource)
	}
}

func BenchmarkRegisterMany(b *testing.B) {
	defs := []resource.Definition{
		validDefinition("Worker", "workers"),
		validDefinition("Job", "jobs"),
		validDefinition("Queue", "queues"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var catalog Catalog
		if err := catalog.RegisterMany(defs...); err != nil {
			b.Fatal(err)
		}
		benchmarkCatalogSizeSink = len(catalog.defsByResource)
	}
}

func BenchmarkResolveResource(b *testing.B) {
	def := validDefinition("Worker", "workers")
	var catalog Catalog
	if err := catalog.Register(def); err != nil {
		b.Fatal(err)
	}

	key := def.GroupResource()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resolved, ok := catalog.ResolveResource(key)
		if !ok {
			b.Fatal("missing definition")
		}
		benchmarkDefinitionSink = resolved
	}
}

func BenchmarkResolveVersionKind(b *testing.B) {
	def := validDefinition("Worker", "workers")
	var catalog Catalog
	if err := catalog.Register(def); err != nil {
		b.Fatal(err)
	}

	key := identity.GroupVersionKind{
		Group:   testGroup,
		Version: "v1",
		Kind:    "Worker",
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resolved, version, ok := catalog.ResolveVersionKind(key)
		if !ok {
			b.Fatal("missing version")
		}
		benchmarkDefinitionSink = resolved
		benchmarkVersionDefinitionSink = version
	}
}

func BenchmarkDefinitions(b *testing.B) {
	var catalog Catalog
	if err := catalog.RegisterMany(
		validDefinition("Worker", "workers"),
		validDefinition("Job", "jobs"),
		validDefinition("Queue", "queues"),
	); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkResourcesSink = catalog.Resources()
		benchmarkCatalogSizeSink = len(catalog.Definitions())
	}
}
