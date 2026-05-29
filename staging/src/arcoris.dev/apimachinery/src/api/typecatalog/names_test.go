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

func TestNamesReturnsStableRegistrationOrder(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))
	requireNoError(t, catalog.Register(types.Define("example.Count", types.Int64())))

	requireNames(t, &catalog, "example.Name", "example.Count")
}

func TestNamesReturnsDetachedSlice(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	names := catalog.Names()
	names[0] = "example.Changed"

	requireNames(t, &catalog, "example.Name")
}

func TestNamesNilCatalogReturnsEmptySlice(t *testing.T) {
	var catalog *Catalog

	requireEqual(t, len(catalog.Names()), 0)
}
