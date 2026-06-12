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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
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
			version: NewVersion(identity.Version("v1"), types.Descriptor{}, Exposed(), Canonical()),
			path:    "definition.versions[v1].desired",
			reason:  ErrorReasonMissingDesired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateVersionDefinitionLocal(tc.version, versionPath(0, tc.version.version))
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

	requireNoError(t, validateVersionDefinitionLocal(version, versionPath(0, version.version)))
}

func TestValidateVersionDefinitionRejectsStructurallyInvalidDesired(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		types.String().MinBytes(2).MaxBytes(1).Descriptor(),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinitionLocal(
		version,
		versionPath(0, version.version),
	)
	resourceErr := requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
	)

	var typeErr *types.DescriptorError
	if !errors.As(resourceErr.Cause, &typeErr) {
		t.Fatalf("expected nested *types.DescriptorError, got %T", resourceErr.Cause)
	}

	requireDetailContains(t, err, "desired descriptor is structurally invalid at")
}

func TestValidateVersionDefinitionRejectsInvalidObserved(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		objectType(),
		Observed(stringType()),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinitionLocal(
		version,
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
		Observed(types.String().MinBytes(2).MaxBytes(1).Descriptor()),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinitionLocal(
		version,
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
	err := invalidSurfaceDetail("desired", types.ErrInvalidDescriptor)
	requireEqual(t, err, "desired descriptor is structurally invalid: invalid descriptor")
}

func TestInvalidSurfaceDetailUsesStructuredTypeError(t *testing.T) {
	cases := []struct {
		name string
		err  *types.DescriptorError
		want string
	}{
		{
			name: "path reason detail",
			err: &types.DescriptorError{
				Record: diagnostic.NewRecord(
					"string.range",
					types.ErrInvalidDescriptor,
					types.DescriptorErrorReasonInvalidRange,
					"minimum must be <= maximum",
				),
			},
			want: "desired descriptor is structurally invalid at string.range: invalid_range: minimum must be <= maximum",
		},
		{
			name: "path reason",
			err: &types.DescriptorError{
				Record: diagnostic.NewRecord(
					"string.range",
					types.ErrInvalidDescriptor,
					types.DescriptorErrorReasonInvalidRange,
					"",
				),
			},
			want: "desired descriptor is structurally invalid at string.range: invalid_range",
		},
		{
			name: "path only",
			err: &types.DescriptorError{
				Record: diagnostic.NewRecord(
					"string.range",
					types.ErrInvalidDescriptor,
					types.DescriptorErrorReason(""),
					"",
				),
			},
			want: "desired descriptor is structurally invalid at string.range",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireEqual(t, invalidSurfaceDetail("desired", tc.err), tc.want)
		})
	}
}

func TestValidateVersionDefinitionLocalAcceptsRootRefs(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		refType("control.arcoris.dev.WorkerDesired"),
		Observed(refType("control.arcoris.dev.WorkerObserved")),
		Exposed(),
		Canonical(),
	)

	requireNoError(t, validateVersionDefinitionLocal(version, versionPath(0, version.version)))
}

func TestValidateVersionDefinitionLocalRejectsInvalidRootRefSyntax(t *testing.T) {
	version := NewVersion(
		identity.Version("v1"),
		refType("bad"),
		Exposed(),
		Canonical(),
	)

	err := validateVersionDefinitionLocal(version, versionPath(0, version.version))
	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
	)
	requireDetailContains(t, err, "invalid_reference_name")
}

func TestValidateVersionDefinitionResolvedHandlesRootRefCases(t *testing.T) {
	cases := []struct {
		name     string
		version  VersionDefinition
		resolver fakeResolver
		path     string
		reason   ErrorReason
		detail   string
	}{
		{
			name: "desired unresolved",
			version: NewVersion(
				identity.Version("v1"),
				refType("control.arcoris.dev.MissingDesired"),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{},
			path:     "definition.versions[v1].desired",
			reason:   ErrorReasonInvalidDesired,
			detail:   "was not found",
		},
		{
			name: "desired recursive",
			version: NewVersion(
				identity.Version("v1"),
				refType("control.arcoris.dev.LoopDesired"),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{
				types.TypeName("control.arcoris.dev.LoopDesired"): types.Define(
					"control.arcoris.dev.LoopDesired",
					types.Ref("control.arcoris.dev.LoopDesired"),
				),
			},
			path:   "definition.versions[v1].desired",
			reason: ErrorReasonInvalidDesired,
			detail: "recursive",
		},
		{
			name: "desired scalar",
			version: NewVersion(
				identity.Version("v1"),
				refType("control.arcoris.dev.TextDesired"),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{
				types.TypeName("control.arcoris.dev.TextDesired"): types.Define(
					"control.arcoris.dev.TextDesired",
					types.String(),
				),
			},
			path:   "definition.versions[v1].desired",
			reason: ErrorReasonDesiredNotObject,
			detail: "got string",
		},
		{
			name: "observed unresolved",
			version: NewVersion(
				identity.Version("v1"),
				objectType(),
				Observed(refType("control.arcoris.dev.MissingObserved")),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{},
			path:     "definition.versions[v1].observed",
			reason:   ErrorReasonInvalidObserved,
			detail:   "was not found",
		},
		{
			name: "observed recursive",
			version: NewVersion(
				identity.Version("v1"),
				objectType(),
				Observed(refType("control.arcoris.dev.LoopObserved")),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{
				types.TypeName("control.arcoris.dev.LoopObserved"): types.Define(
					"control.arcoris.dev.LoopObserved",
					types.Ref("control.arcoris.dev.LoopObserved"),
				),
			},
			path:   "definition.versions[v1].observed",
			reason: ErrorReasonInvalidObserved,
			detail: "recursive",
		},
		{
			name: "observed scalar",
			version: NewVersion(
				identity.Version("v1"),
				objectType(),
				Observed(refType("control.arcoris.dev.TextObserved")),
				Exposed(),
				Canonical(),
			),
			resolver: fakeResolver{
				types.TypeName("control.arcoris.dev.TextObserved"): types.Define(
					"control.arcoris.dev.TextObserved",
					types.String(),
				),
			},
			path:   "definition.versions[v1].observed",
			reason: ErrorReasonObservedNotObject,
			detail: "got string",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateVersionDefinitionResolved(tc.version, tc.resolver, versionPath(0, tc.version.version))
			requireResourceError(t, err, ErrInvalidVersion, tc.path, tc.reason)
			requireDetailContains(t, err, tc.detail)
		})
	}
}
