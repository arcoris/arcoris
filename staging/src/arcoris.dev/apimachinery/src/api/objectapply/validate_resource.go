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

// validateResource checks the caller-supplied resource contract.
//
// objectapply never performs catalog lookup, but the single supplied
// resource.Definition must still have usable family and version identities
// before objectvalidation and valueapply rely on it. This deliberately stops
// short of full resource.ValidateResolved descriptor graph checks.
func (a Applier) validateResource(def resource.Definition) error {
	if err := def.Group().Validate(); err != nil {
		return wrapAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource group is invalid",
			err,
		)
	}
	if err := def.Kind().Validate(); err != nil {
		return wrapAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource kind is invalid",
			err,
		)
	}
	if err := def.Resource().Validate(); err != nil {
		return wrapAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource name is invalid",
			err,
		)
	}
	if err := def.Scope().Validate(); err != nil {
		return wrapAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource scope is invalid",
			err,
		)
	}

	return validateResourceVersions(def.Versions())
}

// validateResourceVersions checks only version-set invariants objectapply needs.
//
// Full resource descriptor validation remains a registration/catalog concern.
// Here we only need a version table that can select a Desired descriptor for
// the live object's API version. Resource surface validation remains owned by
// api/resource and registration/catalog trust boundaries.
func validateResourceVersions(versions []resource.VersionDefinition) error {
	if len(versions) == 0 {
		return wrapAt(
			pathRequestResource,
			ErrInvalidResource,
			ErrorReasonInvalidResource,
			"resource definition must define at least one version",
			resource.ErrInvalidDefinition,
		)
	}

	seen := make(map[string]struct{}, len(versions))
	for _, version := range versions {
		if err := version.Version().Validate(); err != nil {
			return wrapAt(
				pathRequestResource,
				ErrInvalidResource,
				ErrorReasonInvalidResource,
				"resource version is invalid",
				err,
			)
		}

		key := version.Version().String()
		if _, ok := seen[key]; ok {
			return wrapAt(
				pathRequestResource,
				ErrInvalidResource,
				ErrorReasonInvalidResource,
				"resource definition contains duplicate versions",
				resource.ErrInvalidDefinition,
			)
		}
		seen[key] = struct{}{}

		if version.Desired().IsZero() {
			return wrapAt(
				pathRequestResource,
				ErrInvalidResource,
				ErrorReasonInvalidResource,
				"resource version desired descriptor is required",
				resource.ErrInvalidVersion,
			)
		}
	}

	return nil
}
