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

package resource

import (
	"fmt"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

// validateDefinitionVersions checks the version set owned by a Definition.
//
// Version-set rules are intentionally small: at least one version, unique
// version identifiers, at least one exposed version, and exactly one exposed
// canonical version. The package does not define conversion, storage versions,
// rollout policy, or version negotiation.
func validateDefinitionVersions(versions []VersionDefinition, resolver types.Resolver) error {
	if len(versions) == 0 {
		return definitionError(
			pathDefinitionVersions,
			ErrorReasonMissingVersion,
			detailDefinitionNeedsVersion,
		)
	}

	summary := versionSetSummary{
		seen: make(map[identity.Version]struct{}, len(versions)),
	}

	for i, version := range versions {
		path := versionPath(i, version.version)

		if err := validateVersionDefinition(version, resolver, path); err != nil {
			return err
		}

		if err := summary.record(path, version); err != nil {
			return err
		}
	}

	return summary.validate()
}

// versionSetSummary accumulates cross-version invariants for a resource family.
type versionSetSummary struct {
	// seen tracks duplicate version identities while preserving declaration
	// order validation.
	seen map[identity.Version]struct{}

	// exposed counts versions marked for the public API surface.
	exposed int

	// canonical counts versions marked as the family canonical descriptor.
	canonical int

	// canonicalExposed records whether the single canonical descriptor is also
	// exposed.
	canonicalExposed bool
}

// record adds one version to the cross-version summary.
func (s *versionSetSummary) record(path string, version VersionDefinition) error {
	if _, ok := s.seen[version.version]; ok {
		return definitionErrorf(
			path+".version",
			ErrorReasonDuplicateVersion,
			"version %q is declared more than once",
			version.version,
		)
	}

	s.seen[version.version] = struct{}{}

	if version.Exposed() {
		s.exposed++
	}

	if version.Canonical() {
		s.canonical++
		s.canonicalExposed = version.Exposed()
	}

	return nil
}

// validate checks cross-version summary invariants after all versions are
// structurally valid.
func (s versionSetSummary) validate() error {
	if s.exposed == 0 {
		return definitionError(
			pathDefinitionVersions,
			ErrorReasonNoExposedVersion,
			detailDefinitionNeedsExposed,
		)
	}

	if s.canonical == 0 {
		return definitionError(
			pathDefinitionVersions,
			ErrorReasonNoCanonicalVersion,
			detailDefinitionNeedsCanonical,
		)
	}

	if s.canonical > 1 {
		return definitionError(
			pathDefinitionVersions,
			ErrorReasonMultipleCanonical,
			detailDefinitionNeedsCanonical,
		)
	}

	if !s.canonicalExposed {
		return definitionError(
			pathDefinitionVersions,
			ErrorReasonCanonicalNotExposed,
			detailDefinitionCanonicalExposed,
		)
	}

	return nil
}

// versionPath returns the most specific stable descriptor path available for a
// version entry.
func versionPath(index int, version identity.Version) string {
	if err := version.Validate(); err == nil {
		return fmt.Sprintf("%s[%s]", pathDefinitionVersions, version)
	}

	return fmt.Sprintf("%s[%d]", pathDefinitionVersions, index)
}
