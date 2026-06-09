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
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

// VersionDefinition describes one API contract version of a resource family.
//
// Desired is the user-authored state surface. Observed is an optional
// system-computed state surface. Both surfaces are structural api/types
// descriptors and must be object-like when validated by this package.
type VersionDefinition struct {
	// version is the API version identity for this contract surface.
	version identity.Version

	// desired is the required user-authored state descriptor.
	desired types.Descriptor

	// observed is the optional system-computed state descriptor.
	observed types.Descriptor

	// flags stores exposed/canonical descriptor markers.
	flags versionFlags
}

// NewVersion constructs a versioned resource contract.
//
// The desired descriptor is required by validation. The constructor records the
// supplied values without validating them so declarations stay cheap and errors
// are reported consistently by ValidateDefinitionResolved.
func NewVersion(version identity.Version, desired types.Descriptor, options ...VersionOption) VersionDefinition {
	out := VersionDefinition{version: version, desired: desired}
	for _, option := range options {
		if option != nil {
			option(&out)
		}
	}
	return out
}

// IsZero reports whether v is the empty optional sentinel.
//
// A zero VersionDefinition is invalid as a complete API contract version, but it
// remains useful as a missing lookup result.
func (v VersionDefinition) IsZero() bool {
	return v.version.IsZero() &&
		v.desired.IsZero() &&
		v.observed.IsZero() &&
		v.flags == 0
}

// Version returns the API contract version.
func (v VersionDefinition) Version() identity.Version { return v.version }

// Desired returns the user-authored/requested state descriptor.
func (v VersionDefinition) Desired() types.Descriptor { return v.desired }

// Observed returns the optional system-computed/read state descriptor.
func (v VersionDefinition) Observed() (types.Descriptor, bool) {
	if v.observed.IsZero() {
		return types.Descriptor{}, false
	}
	return v.observed, true
}

// Exposed reports whether the version is part of the public API surface.
func (v VersionDefinition) Exposed() bool { return v.flags&versionFlagExposed != 0 }

// Canonical reports whether the version is the canonical family version.
func (v VersionDefinition) Canonical() bool { return v.flags&versionFlagCanonical != 0 }
