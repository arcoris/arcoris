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
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

// checkIncomingDefinitionsLocked checks incoming identities before mutation.
//
// Version indexes are checked before family indexes so callers get precise
// GVR/GVK diagnostics when the concrete version identity is the conflicting
// surface. Resource and kind conflicts are still detected independently when
// versions differ.
func (c *Catalog) checkIncomingDefinitionsLocked(defs []resource.Definition) error {
	seenResources := make(map[identity.GroupResource]int, len(defs))
	seenKinds := make(map[identity.GroupKind]int, len(defs))
	seenVersionResources := make(map[identity.GroupVersionResource]int)
	seenVersionKinds := make(map[identity.GroupVersionKind]int)

	for i, def := range defs {
		if err := c.checkVersionResourceConflictsLocked(i, def, seenVersionResources); err != nil {
			return err
		}
		if err := c.checkVersionKindConflictsLocked(i, def, seenVersionKinds); err != nil {
			return err
		}
		if err := c.checkResourceConflictLocked(i, def, seenResources); err != nil {
			return err
		}
		if err := c.checkKindConflictLocked(i, def, seenKinds); err != nil {
			return err
		}
	}

	return nil
}

// checkResourceConflictLocked checks the primary GroupResource identity.
func (c *Catalog) checkResourceConflictLocked(
	index int,
	def resource.Definition,
	seen map[identity.GroupResource]int,
) error {
	key := groupResourceOf(def)
	path := resourcePath(key)

	if previous, ok := seen[key]; ok && previous != index {
		return catalogErrorf(
			path,
			ErrDuplicateDefinition,
			ErrorReasonDuplicateResource,
			"resource %q is declared more than once in the batch",
			key,
		)
	}
	seen[key] = index

	if _, ok := c.defsByResource[key]; ok {
		return catalogErrorf(
			path,
			ErrDefinitionExists,
			ErrorReasonDefinitionExists,
			"resource %q is already registered",
			key,
		)
	}

	return nil
}

// checkKindConflictLocked checks the version-independent GroupKind identity.
func (c *Catalog) checkKindConflictLocked(
	index int,
	def resource.Definition,
	seen map[identity.GroupKind]int,
) error {
	key := groupKindOf(def)
	path := kindPath(key)

	if previous, ok := seen[key]; ok && previous != index {
		return catalogErrorf(
			path,
			ErrDuplicateDefinition,
			ErrorReasonDuplicateKind,
			"kind %q is declared more than once in the batch",
			key,
		)
	}
	seen[key] = index

	if _, ok := c.resourceByKind[key]; ok {
		return catalogErrorf(
			path,
			ErrDefinitionExists,
			ErrorReasonDefinitionExists,
			"kind %q is already registered",
			key,
		)
	}

	return nil
}

// checkVersionResourceConflictsLocked checks all concrete GVR identities.
func (c *Catalog) checkVersionResourceConflictsLocked(
	index int,
	def resource.Definition,
	seen map[identity.GroupVersionResource]int,
) error {
	for _, key := range versionResourceKeys(def) {
		path := versionResourcePath(key)
		if previous, ok := seen[key]; ok {
			if previous == index {
				continue
			}
			return catalogErrorf(
				path,
				ErrDuplicateDefinition,
				ErrorReasonDuplicateVersionResource,
				"version resource %q is declared more than once in the batch",
				key,
			)
		}
		seen[key] = index

		if _, ok := c.versionByResource[key]; ok {
			return catalogErrorf(
				path,
				ErrDefinitionExists,
				ErrorReasonDefinitionExists,
				"version resource %q is already registered",
				key,
			)
		}
	}

	return nil
}

// checkVersionKindConflictsLocked checks all concrete GVK identities.
func (c *Catalog) checkVersionKindConflictsLocked(
	index int,
	def resource.Definition,
	seen map[identity.GroupVersionKind]int,
) error {
	for _, key := range versionKindKeys(def) {
		path := versionKindPath(key)
		if previous, ok := seen[key]; ok {
			if previous == index {
				continue
			}
			return catalogErrorf(
				path,
				ErrDuplicateDefinition,
				ErrorReasonDuplicateVersionKind,
				"version kind %q is declared more than once in the batch",
				key,
			)
		}
		seen[key] = index

		if _, ok := c.versionByKind[key]; ok {
			return catalogErrorf(
				path,
				ErrDefinitionExists,
				ErrorReasonDefinitionExists,
				"version kind %q is already registered",
				key,
			)
		}
	}

	return nil
}
