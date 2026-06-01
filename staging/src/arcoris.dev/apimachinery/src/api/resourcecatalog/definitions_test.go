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

func TestDefinitionsAndResourcesStableOrder(t *testing.T) {
	first := validDefinition("Worker", "workers")
	second := validDefinition("Job", "jobs")

	var catalog Catalog
	requireNoError(t, catalog.RegisterMany(first, second))

	defs := catalog.Definitions()
	requireEqual(t, len(defs), 2)
	requireEqual(t, defs[0].GroupResource(), first.GroupResource())
	requireEqual(t, defs[1].GroupResource(), second.GroupResource())

	requireSliceEqual(
		t,
		catalog.Resources(),
		[]identity.GroupResource{first.GroupResource(), second.GroupResource()},
	)
}

func TestDefinitionsAndResourcesReturnDetachedSlices(t *testing.T) {
	first := validDefinition("Worker", "workers")
	second := validDefinition("Job", "jobs")

	var catalog Catalog
	requireNoError(t, catalog.RegisterMany(first, second))

	defs := catalog.Definitions()
	defs[0] = second
	requireEqual(t, catalog.Definitions()[0].GroupResource(), first.GroupResource())

	resources := catalog.Resources()
	resources[0] = second.GroupResource()
	requireEqual(t, catalog.Resources()[0], first.GroupResource())
}

func TestDefinitionsAndResourcesNilCatalog(t *testing.T) {
	var catalog *Catalog

	if catalog.Definitions() != nil {
		t.Fatalf("nil Definitions() must return nil")
	}
	if catalog.Resources() != nil {
		t.Fatalf("nil Resources() must return nil")
	}
}
