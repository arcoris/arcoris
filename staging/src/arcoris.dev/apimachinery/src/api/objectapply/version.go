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

package objectapply

import "arcoris.dev/apimachinery/api/resource"

// validateVersionCompatibility rejects cross-version apply.
//
// objectapply v1 has no conversion layer. A caller that accepts one version and
// stores another must convert before calling Apply.
func validateVersionCompatibility(live ValueObject, applied ValueObject) error {
	liveVersion := live.GroupVersionKind().Version
	appliedVersion := applied.GroupVersionKind().Version
	if liveVersion == appliedVersion {
		return nil
	}

	return errorfAt(
		pathObjectAppliedTypeMeta,
		ErrVersionMismatch,
		ErrorReasonVersionMismatch,
		"applied version %q does not match live version %q",
		appliedVersion,
		liveVersion,
	)
}

// selectVersion returns the live object's resource version descriptor.
//
// The returned descriptor supplies the Desired type passed to valueapply. The
// helper uses only the supplied resource.Definition and never consults a
// resource catalog.
func selectVersion(
	obj ValueObject,
	def resource.Definition,
) (resource.VersionDefinition, error) {
	version, ok := def.Version(obj.GroupVersionKind().Version)
	if !ok {
		return resource.VersionDefinition{}, errorfAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource %s does not define object API version %q",
			def.GroupKind(),
			obj.GroupVersionKind().Version,
		)
	}

	return version, nil
}
