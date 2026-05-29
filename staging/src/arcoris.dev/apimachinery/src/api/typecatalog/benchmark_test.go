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

package typecatalog

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

var (
	benchmarkDefinitionSink  types.TypeDefinition
	benchmarkNamesSink       []types.TypeName
	benchmarkCatalogSizeSink int
)

func BenchmarkCatalogResolveType(b *testing.B) {
	var catalog Catalog
	if err := catalog.Register(types.Define("example.Name", types.String().MinLen(1))); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		def, ok := catalog.ResolveType("example.Name")
		if !ok {
			b.Fatal("missing definition")
		}
		benchmarkDefinitionSink = def
	}
}

func BenchmarkCatalogNames(b *testing.B) {
	var catalog Catalog
	if err := catalog.RegisterMany(
		types.Define("example.Name", types.String().MinLen(1)),
		types.Define("example.Count", types.Int64().Min(0)),
		types.Define("example.Enabled", types.Bool()),
	); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkNamesSink = catalog.Names()
	}
}

func BenchmarkCatalogRegisterManySmallBatch(b *testing.B) {
	defs := []types.TypeDefinition{
		types.Define("example.Name", types.String().MinLen(1)),
		types.Define("example.Count", types.Int64().Min(0)),
		types.Define("example.Enabled", types.Bool()),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var catalog Catalog
		if err := catalog.RegisterMany(defs...); err != nil {
			b.Fatal(err)
		}
		benchmarkCatalogSizeSink = len(catalog.defs)
	}
}
