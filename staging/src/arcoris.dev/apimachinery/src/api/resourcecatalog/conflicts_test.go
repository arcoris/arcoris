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

	"arcoris.dev/apimachinery/api/resource"
)

func TestCheckIncomingDefinitionsLockedRejectsDuplicateIdentities(t *testing.T) {
	cases := []struct {
		name   string
		defs   []resource.Definition
		path   string
		reason ErrorReason
	}{
		{
			name: "group resource",
			defs: []resource.Definition{
				validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
				validDefinition("WorkerCopy", "workers", objectVersion("v2", resource.Canonical())),
			},
			path:   "definitions[control.arcoris.dev:workers]",
			reason: ErrorReasonDuplicateResource,
		},
		{
			name: "group kind",
			defs: []resource.Definition{
				validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
				validDefinition("Worker", "workerobjects", objectVersion("v2", resource.Canonical())),
			},
			path:   "definitions[control.arcoris.dev#Worker]",
			reason: ErrorReasonDuplicateKind,
		},
		{
			name: "group version resource",
			defs: []resource.Definition{
				validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
				validDefinition("WorkerCopy", "workers", objectVersion("v1", resource.Canonical())),
			},
			path:   "definitions[control.arcoris.dev/v1:workers]",
			reason: ErrorReasonDuplicateVersionResource,
		},
		{
			name: "group version kind",
			defs: []resource.Definition{
				validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
				validDefinition("Worker", "workerobjects", objectVersion("v1", resource.Canonical())),
			},
			path:   "definitions[control.arcoris.dev/v1#Worker]",
			reason: ErrorReasonDuplicateVersionKind,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var catalog Catalog
			err := catalog.checkIncomingDefinitionsLocked(tc.defs)
			requireCatalogError(t, err, ErrDuplicateDefinition, tc.path, tc.reason)
		})
	}
}

func TestCheckIncomingDefinitionsLockedRejectsExistingIdentityConflicts(t *testing.T) {
	cases := []struct {
		name     string
		existing resource.Definition
		incoming resource.Definition
		path     string
	}{
		{
			name:     "group resource",
			existing: validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
			incoming: validDefinition("WorkerCopy", "workers", objectVersion("v2", resource.Canonical())),
			path:     "definitions[control.arcoris.dev:workers]",
		},
		{
			name:     "group kind",
			existing: validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
			incoming: validDefinition("Worker", "workerobjects", objectVersion("v2", resource.Canonical())),
			path:     "definitions[control.arcoris.dev#Worker]",
		},
		{
			name:     "group version resource",
			existing: validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
			incoming: validDefinition("WorkerCopy", "workers", objectVersion("v1", resource.Canonical())),
			path:     "definitions[control.arcoris.dev/v1:workers]",
		},
		{
			name:     "group version kind",
			existing: validDefinition("Worker", "workers", objectVersion("v1", resource.Canonical())),
			incoming: validDefinition("Worker", "workerobjects", objectVersion("v1", resource.Canonical())),
			path:     "definitions[control.arcoris.dev/v1#Worker]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var catalog Catalog
			catalog.storeLocked(tc.existing)

			err := catalog.checkIncomingDefinitionsLocked([]resource.Definition{tc.incoming})
			requireCatalogError(t, err, ErrDefinitionExists, tc.path, ErrorReasonDefinitionExists)
		})
	}
}

func TestCheckIncomingDefinitionsLockedLeavesSameDefinitionVersionDuplicatesToValidation(t *testing.T) {
	def := validDefinition(
		"Worker",
		"workers",
		objectVersion("v1", resource.Canonical()),
		objectVersion("v1", resource.Exposed()),
	)

	var catalog Catalog
	requireNoError(t, catalog.checkIncomingDefinitionsLocked([]resource.Definition{def}))
}
