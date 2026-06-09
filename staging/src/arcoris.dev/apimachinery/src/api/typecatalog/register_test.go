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
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestRegisterStoresValidDefinition(t *testing.T) {
	var catalog Catalog

	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	requireDefinition(t, &catalog, "example.Name")
}

func TestRegisterRejectsDuplicateExistingName(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	err := catalog.Register(types.Define("example.Name", types.Int64()))
	requireErrorIs(t, err, ErrDefinitionExists)
	requireErrorNotIs(t, err, types.ErrInvalidDescriptorReference)
}

func TestRegisterRejectsInvalidDefinition(t *testing.T) {
	var catalog Catalog

	requireErrorIs(
		t,
		catalog.Register(types.Define("bad", types.String())),
		types.ErrInvalidDescriptorReference,
	)

	requireErrorIs(
		t,
		catalog.Register(types.Define("example.Bad", types.ListOf(types.DescriptorExpr(nil)))),
		types.ErrInvalidDescriptor,
	)
}

func TestRegisterManyAtomicOnInvalidDefinition(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Existing", types.String())))

	err := catalog.RegisterMany(
		types.Define("example.Next", types.String()),
		types.Define("example.Bad", types.ListOf(types.DescriptorExpr(nil))),
	)
	requireErrorIs(t, err, types.ErrInvalidDescriptor)

	_, ok := catalog.Resolve("example.Next")
	requireEqual(t, ok, false)
}

func TestRegisterManyAtomicOnExistingConflict(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Existing", types.String())))

	err := catalog.RegisterMany(
		types.Define("example.Next", types.String()),
		types.Define("example.Existing", types.Int64()),
	)
	requireErrorIs(t, err, ErrDefinitionExists)
	requireErrorNotIs(t, err, types.ErrInvalidDescriptorReference)

	_, ok := catalog.Resolve("example.Next")
	requireEqual(t, ok, false)
}

func TestRegisterManyAtomicOnBatchDuplicate(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String()),
		types.Define("example.Name", types.Int64()),
	)
	requireErrorIs(t, err, ErrDuplicateDefinition)
	requireEqual(t, errors.Is(err, types.ErrDuplicateField), false)
	requireErrorNotIs(t, err, types.ErrInvalidDescriptorReference)

	_, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, false)
}

func TestRegisterManyAllowsReferencesInsideBatch(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String().MinBytes(1)),
		types.Define("example.NameList", types.ListOf(types.Ref("example.Name"))),
	)
	requireNoError(t, err)

	_, ok := catalog.Resolve("example.NameList")
	requireEqual(t, ok, true)
}

func TestRegisterManyRejectsUnresolvedExternalRefs(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(
		types.Define("example.Name", types.String()),
		types.Define("example.ExternalList", types.ListOf(types.Ref("example.Missing"))),
	)
	requireErrorIs(t, err, types.ErrUnresolvedDescriptorReference)

	_, ok := catalog.Resolve("example.Name")
	requireEqual(t, ok, false)
}
