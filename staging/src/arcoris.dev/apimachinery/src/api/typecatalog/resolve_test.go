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

func TestResolveReturnsDefinition(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	def, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, true)
	requireEqual(t, def.Name(), types.TypeName("example.Name"))
}

func TestResolveReturnsDetachedDefinition(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String().Enum("alpha"))))

	def, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, true)
	view, ok := def.Descriptor().AsString()
	requireEqual(t, ok, true)
	enum := view.Enum()
	enum[0] = "changed"

	defAgain, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, true)
	view, ok = defAgain.Descriptor().AsString()
	requireEqual(t, ok, true)
	requireEqual(t, view.Enum()[0], "alpha")
}

func TestResolveNilCatalogBehavesLikeEmptyCatalog(t *testing.T) {
	var catalog *Catalog

	_, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, false)
}
