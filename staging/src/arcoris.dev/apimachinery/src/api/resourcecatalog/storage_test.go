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
)

func TestEnsureStorageLocked(t *testing.T) {
	var catalog Catalog
	catalog.ensureStorageLocked()

	if catalog.defsByResource == nil ||
		catalog.resourceByKind == nil ||
		catalog.versionByResource == nil ||
		catalog.versionByKind == nil {
		t.Fatalf("ensureStorageLocked did not initialize all maps")
	}
}

func TestStoreLockedPopulatesIndexes(t *testing.T) {
	def := validDefinition("Worker", "workers")

	var catalog Catalog
	catalog.storeLocked(def)

	requireEqual(t, catalog.defsByResource[def.GroupResource()].GroupResource(), def.GroupResource())
	requireEqual(t, catalog.resourceByKind[def.GroupKind()], def.GroupResource())

	_, ok := catalog.versionByResource[identity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v1",
		Resource: "workers",
	}]
	requireEqual(t, ok, true)

	_, ok = catalog.versionByKind[identity.GroupVersionKind{
		Group:   testGroup,
		Version: "v1",
		Kind:    "Worker",
	}]
	requireEqual(t, ok, true)
}

func TestCloneLockedDetachesStorage(t *testing.T) {
	first := validDefinition("Worker", "workers")
	second := validDefinition("Job", "jobs")

	var catalog Catalog
	catalog.storeLocked(first)

	clone := catalog.cloneLocked()
	clone.storeLocked(second)

	if _, ok := catalog.ResolveResource(second.GroupResource()); ok {
		t.Fatalf("clone mutation leaked into original catalog")
	}
	requireEqual(t, len(catalog.order), 1)
	requireEqual(t, len(clone.order), 2)
}
