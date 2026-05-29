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
	"arcoris.dev/apimachinery/api/types"
)

func TestValidateVersionDefinitionRejectsInvalidVersionAndMissingDesired(t *testing.T) {
	cases := []struct {
		name    string
		version VersionDefinition
		path    string
		reason  ErrorReason
	}{
		{
			name:    "invalid version",
			version: NewVersion(identity.Version("v0"), objectType(), Exposed(), Canonical()),
			path:    "definition.versions[0].version",
			reason:  ErrorReasonInvalidVersion,
		},
		{
			name:    "missing desired",
			version: NewVersion(identity.Version("v1"), types.Type{}, Exposed(), Canonical()),
			path:    "definition.versions[v1].desired",
			reason:  ErrorReasonMissingDesired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateVersionDefinition(tc.version, nil, versionPath(0, tc.version.version))
			requireResourceError(t, err, ErrInvalidVersion, tc.path, tc.reason)
		})
	}
}

func TestValidateVersionDefinitionAcceptsValidVersionWithoutObserved(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Exposed(),
		Canonical(),
	)

	requireNoError(t, validateVersionDefinition(version, nil, versionPath(0, version.version)))
}

func TestValidateVersionDefinitionRejectsStructurallyInvalidDesired(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		types.String().MinLen(2).MaxLen(1).Type(),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinition(
		version,
		nil,
		versionPath(0, version.version),
	)
	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
	)
}

func TestValidateVersionDefinitionRejectsInvalidObserved(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Observed(stringType()),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinition(
		version,
		nil,
		versionPath(0, version.version),
	)
	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].observed",
		ErrorReasonObservedNotObject,
	)
}

func TestValidateVersionDefinitionRejectsStructurallyInvalidObserved(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Observed(types.String().MinLen(2).MaxLen(1).Type()),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinition(
		version,
		nil,
		versionPath(0, version.version),
	)
	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].observed",
		ErrorReasonInvalidObserved,
	)
}

func TestInvalidSurfaceDetail(t *testing.T) {
	err := invalidSurfaceDetail("desired", types.ErrInvalidType)
	requireEqual(t, err, "desired descriptor is structurally invalid: invalid type")
}
