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

// validateSurface checks one resource API surface descriptor.
//
// Resource versions only accept object-like root surfaces: either a direct
// object descriptor or a reference that the supplied resolver can prove resolves
// to an object. This keeps resource definitions structural and avoids sneaking
// value validation, codecs, runtime object machinery, or resource metadata into
// this package.
func validateSurface(
	typ types.Type,
	resolver types.Resolver,
	path string,
	invalidReason ErrorReason,
	objectReason ErrorReason,
	label string,
) error {
	if err := types.ValidateType(typ, resolver); err != nil {
		return nestedVersionError(
			path,
			invalidReason,
			invalidSurfaceDetail(label, err),
			err,
		)
	}

	return requireObjectLike(
		typ,
		resolver,
		path,
		objectReason,
		label,
	)
}
