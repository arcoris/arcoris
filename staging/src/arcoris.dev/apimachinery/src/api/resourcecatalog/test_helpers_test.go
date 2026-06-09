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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/typecatalog"
	"arcoris.dev/apimachinery/api/types"
)

const testGroup = identity.Group("control.arcoris.dev")

func objectType() types.Descriptor {
	return types.Object().Descriptor()
}

func objectVersion(version identity.Version, options ...resource.VersionOption) resource.VersionDefinition {
	opts := append([]resource.VersionOption{resource.Exposed()}, options...)
	return resource.NewVersion(version, objectType(), opts...)
}

func validDefinition(
	kind identity.Kind,
	res identity.Resource,
	versions ...resource.VersionDefinition,
) resource.Definition {
	if len(versions) == 0 {
		versions = []resource.VersionDefinition{
			objectVersion(identity.Version("v1"), resource.Canonical()),
		}
	}

	return resource.NewDefinition(
		testGroup,
		kind,
		res,
		resource.ScopeNamespaced,
		versions...,
	)
}

func invalidDefinition() resource.Definition {
	return resource.NewDefinition(
		testGroup,
		identity.Kind("Broken"),
		identity.Resource("brokens"),
		resource.ScopeNamespaced,
	)
}

func refDefinition() resource.Definition {
	return resource.NewDefinition(
		testGroup,
		identity.Kind("Worker"),
		identity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion(
			identity.Version("v1"),
			types.Ref("control.arcoris.dev.WorkerDesired").Descriptor(),
			resource.Observed(types.Ref("control.arcoris.dev.WorkerObserved").Descriptor()),
			resource.Exposed(),
			resource.Canonical(),
		),
	)
}

func realTypeCatalog(t *testing.T) *typecatalog.Catalog {
	t.Helper()

	var catalog typecatalog.Catalog
	requireNoError(
		t,
		catalog.RegisterMany(
			types.Define("control.arcoris.dev.WorkerDesired", types.Object()),
			types.Define("control.arcoris.dev.WorkerObserved", types.Object()),
		),
	)
	return &catalog
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireCatalogError(t *testing.T, err error, target error, path string, reason ErrorReason) *Error {
	t.Helper()
	requireErrorIs(t, err, target)

	var catalogErr *Error
	if !errors.As(err, &catalogErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if catalogErr.Path != path {
		t.Fatalf("Error.Path = %q, want %q", catalogErr.Path, path)
	}
	if catalogErr.Reason != reason {
		t.Fatalf("Error.Reason = %q, want %q", catalogErr.Reason, reason)
	}
	if catalogErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
	return catalogErr
}

func requireDetailContains(t *testing.T, err error, text string) {
	t.Helper()

	var catalogErr *Error
	if !errors.As(err, &catalogErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if !strings.Contains(catalogErr.Detail, text) {
		t.Fatalf("Error.Detail = %q, want to contain %q", catalogErr.Detail, text)
	}
}

func requireEqual[T comparable](t *testing.T, got T, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func requireSliceEqual[T comparable](t *testing.T, got []T, want []T) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d; got %#v, want %#v", len(got), len(want), got, want)
	}

	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %#v, want %#v; got %#v, want %#v", i, got[i], want[i], got, want)
		}
	}
}
