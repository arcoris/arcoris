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
	"errors"
	"fmt"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestCatalogZeroValueUsableAndStableOrder(t *testing.T) {
	var catalog Catalog

	requireNoError(t, catalog.Register(types.Define("example.Name", types.String().MinLen(1))))
	requireNoError(t, catalog.Register(types.Define("example.Count", types.Int64().Min(0))))

	requireNames(t, &catalog, "example.Name", "example.Count")
	requireDefinitions(t, &catalog, "example.Name", "example.Count")
}

func TestCatalogRegisterValidDefinition(t *testing.T) {
	var catalog Catalog

	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	requireDefinition(t, &catalog, "example.Name")
}

func TestCatalogRejectsDuplicateExistingName(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	err := catalog.Register(types.Define("example.Name", types.Int64()))
	requireErrorIs(t, err, ErrDefinitionExists)
	requireErrorIs(t, err, types.ErrInvalidTypeReference)
}

func TestCatalogRejectsInvalidDefinition(t *testing.T) {
	var catalog Catalog

	requireErrorIs(t, catalog.Register(types.Define("bad", types.String())), types.ErrInvalidTypeReference)
	requireErrorIs(t, catalog.Register(types.Define("example.Bad", types.ListOf(types.TypeExpr(nil)))), types.ErrInvalidType)
}

func TestCatalogRegisterManyAtomicOnInvalidDefinition(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Existing", types.String())))

	err := catalog.RegisterMany(
		types.Define("example.Next", types.String()),
		types.Define("example.Bad", types.ListOf(types.TypeExpr(nil))),
	)
	requireErrorIs(t, err, types.ErrInvalidType)

	_, ok := catalog.ResolveType("example.Next")
	requireEqual(t, ok, false)
}

func TestCatalogRegisterManyAtomicOnExistingConflict(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Existing", types.String())))

	err := catalog.RegisterMany(
		types.Define("example.Next", types.String()),
		types.Define("example.Existing", types.Int64()),
	)
	requireErrorIs(t, err, ErrDefinitionExists)

	_, ok := catalog.ResolveType("example.Next")
	requireEqual(t, ok, false)
}

func TestCatalogRegisterManyAtomicOnBatchDuplicate(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String()),
		types.Define("example.Name", types.Int64()),
	)
	requireErrorIs(t, err, ErrDuplicateDefinition)
	requireEqual(t, errors.Is(err, types.ErrDuplicateField), false)

	_, ok := catalog.ResolveType("example.Name")
	requireEqual(t, ok, false)
}

func TestCatalogRegisterManyAllowsReferencesInsideBatch(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String().MinLen(1)),
		types.Define("example.NameList", types.ListOf(types.Ref("example.Name"))),
	)
	requireNoError(t, err)

	_, ok := catalog.ResolveType("example.NameList")
	requireEqual(t, ok, true)
}

func TestCatalogRegisterManyRejectsUnresolvedExternalRefs(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String()),
		types.Define("example.ExternalList", types.ListOf(types.Ref("example.Missing"))),
	)
	requireErrorIs(t, err, types.ErrUnknownTypeReference)

	_, ok := catalog.ResolveType("example.Name")
	requireEqual(t, ok, false)
}

func TestCatalogResolveTypeAndDefinitionsDetached(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String().Enum("alpha"))))

	def, ok := catalog.ResolveType("example.Name")
	requireEqual(t, ok, true)
	view, ok := def.Type().String()
	requireEqual(t, ok, true)
	enum := view.Enum()
	enum[0] = "changed"

	defAgain, ok := catalog.ResolveType("example.Name")
	requireEqual(t, ok, true)
	view, ok = defAgain.Type().String()
	requireEqual(t, ok, true)
	requireEqual(t, view.Enum()[0], "alpha")

	defs := catalog.Definitions()
	view, ok = defs[0].Type().String()
	requireEqual(t, ok, true)
	enum = view.Enum()
	enum[0] = "changed"
	view, ok = catalog.Definitions()[0].Type().String()
	requireEqual(t, ok, true)
	requireEqual(t, view.Enum()[0], "alpha")
}

func TestNilCatalogReadMethods(t *testing.T) {
	var catalog *Catalog

	_, ok := catalog.ResolveType("example.Name")
	requireEqual(t, ok, false)
	requireEqual(t, len(catalog.Names()), 0)
	requireEqual(t, len(catalog.Definitions()), 0)
}

func TestCatalogConcurrentAccessRaceFree(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				catalog.ResolveType("example.Name")
				catalog.Names()
				catalog.Definitions()
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			name := types.TypeName(fmt.Sprintf("example.Generated%d", j))
			_ = catalog.Register(types.Define(name, types.String()))
		}
	}()
	wg.Wait()
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error matching %v, got %v", target, err)
	}
}

func requireEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func requireDefinition(t *testing.T, c *Catalog, name types.TypeName) types.TypeDefinition {
	t.Helper()
	def, ok := c.ResolveType(name)
	if !ok {
		t.Fatalf("missing definition %s", name)
	}
	return def
}

func requireNames(t *testing.T, c *Catalog, want ...types.TypeName) {
	t.Helper()
	got := c.Names()
	requireEqual(t, len(got), len(want))
	for i := range want {
		requireEqual(t, got[i], want[i])
	}
}

func requireDefinitions(t *testing.T, c *Catalog, want ...types.TypeName) {
	t.Helper()
	got := c.Definitions()
	requireEqual(t, len(got), len(want))
	for i := range want {
		requireEqual(t, got[i].Name(), want[i])
	}
}
