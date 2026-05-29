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
	"testing"

	"arcoris.dev/apimachinery/api/identity"
)

func TestValidateDefinitionVersionsRejectsInvalidVersionSet(t *testing.T) {
	cases := []struct {
		name   string
		input  []VersionDefinition
		path   string
		reason ErrorReason
	}{
		{
			name:   "missing versions",
			path:   pathDefinitionVersions,
			reason: ErrorReasonMissingVersion,
		},
		{
			name: "duplicate version",
			input: []VersionDefinition{
				NewVersion(identity.Version("v1"), objectType(), Exposed(), Canonical()),
				NewVersion(identity.Version("v1"), objectType(), Exposed()),
			},
			path:   "definition.versions[v1].version",
			reason: ErrorReasonDuplicateVersion,
		},
		{
			name: "no exposed",
			input: []VersionDefinition{
				NewVersion(identity.Version("v1"), objectType(), Canonical()),
			},
			path:   pathDefinitionVersions,
			reason: ErrorReasonNoExposedVersion,
		},
		{
			name: "no canonical",
			input: []VersionDefinition{
				NewVersion(identity.Version("v1"), objectType(), Exposed()),
			},
			path:   pathDefinitionVersions,
			reason: ErrorReasonNoCanonicalVersion,
		},
		{
			name: "multiple canonical",
			input: []VersionDefinition{
				NewVersion(identity.Version("v1alpha1"), objectType(), Exposed(), Canonical()),
				NewVersion(identity.Version("v1"), objectType(), Exposed(), Canonical()),
			},
			path:   pathDefinitionVersions,
			reason: ErrorReasonMultipleCanonical,
		},
		{
			name: "canonical not exposed",
			input: []VersionDefinition{
				NewVersion(identity.Version("v1alpha1"), objectType(), Exposed()),
				NewVersion(identity.Version("v1"), objectType(), Canonical()),
			},
			path:   pathDefinitionVersions,
			reason: ErrorReasonCanonicalNotExposed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDefinitionVersions(tc.input, nil)
			requireResourceError(t, err, ErrInvalidDefinition, tc.path, tc.reason)
		})
	}
}

func TestValidateDefinitionVersionsPropagatesVersionErrors(t *testing.T) {
	err := validateDefinitionVersions(
		[]VersionDefinition{
			NewVersion(identity.Version("v0"), objectType(), Exposed(), Canonical()),
		},
		nil,
	)

	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[0].version",
		ErrorReasonInvalidVersion,
	)
}

func TestVersionPath(t *testing.T) {
	requireEqual(t, versionPath(0, identity.Version("v1")), "definition.versions[v1]")
	requireEqual(t, versionPath(2, identity.Version("v0")), "definition.versions[2]")
}
