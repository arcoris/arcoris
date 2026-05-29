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

func TestResolveMethodsReturnRegisteredDefinitions(t *testing.T) {
	def := validDefinition(
		"Worker",
		"workers",
		objectVersion("v1alpha1"),
		objectVersion("v1", resource.Canonical()),
	)
	var catalog Catalog
	requireNoError(t, catalog.Register(def))

	byResource, ok := catalog.ResolveResource(def.GroupResource())
	requireEqual(t, ok, true)
	requireEqual(t, byResource.GroupResource(), def.GroupResource())

	byKind, ok := catalog.ResolveKind(def.GroupKind())
	requireEqual(t, ok, true)
	requireEqual(t, byKind.GroupKind(), def.GroupKind())

	byGVR, version, ok := catalog.ResolveVersionResource(identity.GroupVersionResource{
		Group:    testGroup,
		Version:  identity.Version("v1alpha1"),
		Resource: identity.Resource("workers"),
	})
	requireEqual(t, ok, true)
	requireEqual(t, byGVR.GroupResource(), def.GroupResource())
	requireEqual(t, version.Version(), identity.Version("v1alpha1"))

	byGVK, version, ok := catalog.ResolveVersionKind(identity.GroupVersionKind{
		Group:   testGroup,
		Version: identity.Version("v1alpha1"),
		Kind:    identity.Kind("Worker"),
	})
	requireEqual(t, ok, true)
	requireEqual(t, byGVK.GroupKind(), def.GroupKind())
	requireEqual(t, version.Version(), identity.Version("v1alpha1"))
}

func TestResolveMethodsReturnFalseForUnknownKeys(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(validDefinition("Worker", "workers")))

	if _, ok := catalog.ResolveResource(identity.GroupResource{Group: testGroup, Resource: "missing"}); ok {
		t.Fatalf("ResolveResource unknown ok = true")
	}
	if _, ok := catalog.ResolveKind(identity.GroupKind{Group: testGroup, Kind: "Missing"}); ok {
		t.Fatalf("ResolveKind unknown ok = true")
	}
	if _, _, ok := catalog.ResolveVersionResource(identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v1",
		Resource: "missing",
	}); ok {
		t.Fatalf("ResolveVersionResource unknown ok = true")
	}
	if _, _, ok := catalog.ResolveVersionKind(identity.GroupVersionKind{
		Group:   testGroup,
		Version: "v1",
		Kind:    "Missing",
	}); ok {
		t.Fatalf("ResolveVersionKind unknown ok = true")
	}
}

func TestResolveMethodsNilCatalog(t *testing.T) {
	var catalog *Catalog

	if _, ok := catalog.ResolveResource(identity.GroupResource{}); ok {
		t.Fatalf("nil ResolveResource ok = true")
	}
	if _, ok := catalog.ResolveKind(identity.GroupKind{}); ok {
		t.Fatalf("nil ResolveKind ok = true")
	}
	if _, _, ok := catalog.ResolveVersionResource(identity.GroupVersionResource{}); ok {
		t.Fatalf("nil ResolveVersionResource ok = true")
	}
	if _, _, ok := catalog.ResolveVersionKind(identity.GroupVersionKind{}); ok {
		t.Fatalf("nil ResolveVersionKind ok = true")
	}
}

func TestResolveKindFailsSafelyOnInconsistentIndex(t *testing.T) {
	var catalog Catalog
	catalog.ensureStorageLocked()
	catalog.resourceByKind[identity.GroupKind{
		Group: testGroup,
		Kind:  "Worker",
	}] = identity.GroupResource{
		Group:    testGroup,
		Resource: "workers",
	}

	if _, ok := catalog.ResolveKind(identity.GroupKind{
		Group: testGroup,
		Kind:  "Worker",
	}); ok {
		t.Fatalf("inconsistent kind index ok = true")
	}
}

func TestResolveVersionFailsSafelyOnInconsistentIndex(t *testing.T) {
	var catalog Catalog
	catalog.ensureStorageLocked()
	catalog.versionByResource[identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v1",
		Resource: "workers",
	}] = versionRef{
		resource: identity.GroupResource{Group: testGroup, Resource: "workers"},
		version:  "v1",
	}

	if _, _, ok := catalog.ResolveVersionResource(identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v1",
		Resource: "workers",
	}); ok {
		t.Fatalf("inconsistent version index ok = true")
	}
}

func TestResolveVersionFailsSafelyWhenVersionMissing(t *testing.T) {
	def := validDefinition("Worker", "workers")
	var catalog Catalog
	catalog.storeLocked(def)
	catalog.versionByResource[identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v2",
		Resource: "workers",
	}] = versionRef{
		resource: def.GroupResource(),
		version:  "v2",
	}

	if _, _, ok := catalog.ResolveVersionResource(identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v2",
		Resource: "workers",
	}); ok {
		t.Fatalf("missing version index ok = true")
	}
}
