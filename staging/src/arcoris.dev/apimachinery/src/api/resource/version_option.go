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

import "arcoris.dev/apimachinery/api/types"

// VersionOption mutates VersionDefinition construction state.
//
// Options are deliberately narrow descriptor switches. They do not define
// conversion strategy, storage versioning, defaulting, or admission behavior.
type VersionOption func(*VersionDefinition)

// Exposed marks a resource version as part of the public API surface.
func Exposed() VersionOption {
	return func(v *VersionDefinition) {
		v.flags |= versionFlagExposed
	}
}

// Canonical marks a resource version as the canonical family version.
//
// Canonical does not define storage behavior. Storage and persistence semantics
// belong to future runtime/storage layers.
func Canonical() VersionOption {
	return func(v *VersionDefinition) {
		v.flags |= versionFlagCanonical
	}
}

// Observed sets the optional system-computed/read state descriptor.
//
// Observed describes a structural read surface. It does not introduce status
// subresources, condition conventions, persistence rules, or object value
// validation.
func Observed(observed types.Descriptor) VersionOption {
	return func(v *VersionDefinition) {
		v.observed = observed
	}
}
