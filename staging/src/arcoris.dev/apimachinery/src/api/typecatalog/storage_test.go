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

func TestStoreLockedInitializesMapAndOrder(t *testing.T) {
	var catalog Catalog

	catalog.storeLocked(types.Define("example.dev.Name", types.String()))

	requireEqual(t, len(catalog.defs), 1)
	requireEqual(t, catalog.order[0], types.TypeName("example.dev.Name"))
}

func TestCloneLockedDetachesOrderAndDefinitions(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.dev.Name", types.String().Enum("alpha"))))

	clone := catalog.cloneLocked()
	clone.order[0] = "example.dev.Changed"

	requireNames(t, &catalog, "example.dev.Name")

	view, ok := clone.defs["example.dev.Name"].Descriptor().AsString()
	requireEqual(t, ok, true)
	enum := view.Enum()
	enum[0] = "changed"

	view, ok = catalog.defs["example.dev.Name"].Descriptor().AsString()
	requireEqual(t, ok, true)
	requireEqual(t, view.Enum()[0], "alpha")
}
