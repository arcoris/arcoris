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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

func TestRegisterValidInlineDefinition(t *testing.T) {
	var catalog Catalog
	def := validDefinition("Worker", "workers")

	requireNoError(t, catalog.Register(def))

	resolved, ok := catalog.ResolveResource(def.GroupResource())
	if !ok {
		t.Fatalf("ResolveResource() ok = false")
	}
	requireEqual(t, resolved.GroupResource(), def.GroupResource())
}

func TestRegisterRejectsNilCatalog(t *testing.T) {
	var catalog *Catalog

	err := catalog.RegisterMany(validDefinition("Worker", "workers"))
	requireCatalogError(t, err, ErrNilCatalog, "catalog", ErrorReasonNilCatalog)
}

func TestRegisterManyAcceptsDistinctDefinitions(t *testing.T) {
	var catalog Catalog

	requireNoError(
		t,
		catalog.RegisterMany(
			validDefinition("Worker", "workers"),
			validDefinition("Job", "jobs"),
		),
	)

	requireEqual(t, len(catalog.Definitions()), 2)
}

func TestRegisterManyIsAtomicOnInvalidDefinition(t *testing.T) {
	var catalog Catalog
	valid := validDefinition("Worker", "workers")

	err := catalog.RegisterMany(valid, invalidDefinition())
	requireCatalogError(t, err, ErrInvalidCatalog, "definitions[1]", ErrorReasonInvalidDefinition)

	if _, ok := catalog.ResolveResource(valid.GroupResource()); ok {
		t.Fatalf("RegisterMany stored partial state")
	}
}

func TestRegisterManyRejectsInvalidDefinition(t *testing.T) {
	var catalog Catalog

	err := catalog.RegisterMany(invalidDefinition())
	requireCatalogError(t, err, ErrInvalidCatalog, "definitions[0]", ErrorReasonInvalidDefinition)
	requireErrorIs(t, err, resource.ErrInvalidDefinition)
}

func TestRegisterManyRejectsTypeRefRootsWithoutResolver(t *testing.T) {
	var catalog Catalog

	err := catalog.Register(refDefinition())
	requireCatalogError(t, err, ErrInvalidCatalog, "definitions[0]", ErrorReasonInvalidDefinition)
	requireErrorIs(t, err, resource.ErrInvalidVersion)
}

func TestRegisterManyRejectsDuplicateIdentitiesInsideBatch(t *testing.T) {
	var catalog Catalog
	err := catalog.RegisterMany(
		validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
		validDefinition("WorkerCopy", "workers", objectVersion("v2", resource.Canonical())),
	)
	requireCatalogError(
		t,
		err,
		ErrDuplicateDefinition,
		"definitions[control.arcoris.dev:workers]",
		ErrorReasonDuplicateResource,
	)
}

func TestRegisterManyRejectsExistingIdentityConflicts(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(validDefinition("Worker", "workers")))

	err := catalog.Register(validDefinition("WorkerCopy", "workers"))
	requireCatalogError(
		t,
		err,
		ErrDefinitionExists,
		"definitions[control.arcoris.dev/v1:workers]",
		ErrorReasonDefinitionExists,
	)
}

func TestRegisterManyPreservesNestedResourceErrors(t *testing.T) {
	var catalog Catalog
	def := validDefinition(
		"Worker",
		"workers",
		objectVersion(identity.Version("v0"), resource.Canonical()),
	)

	err := catalog.Register(def)
	requireCatalogError(t, err, ErrInvalidCatalog, "definitions[0]", ErrorReasonInvalidDefinition)
	requireErrorIs(t, err, resource.ErrInvalidVersion)

	var resourceErr *resource.Error
	if !errors.As(err, &resourceErr) {
		t.Fatalf("expected nested *resource.Error, got %T", err)
	}
}
