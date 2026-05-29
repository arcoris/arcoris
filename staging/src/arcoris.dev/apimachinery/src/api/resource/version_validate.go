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

	"arcoris.dev/apimachinery/api/types"
)

// validateVersionDefinition checks one version descriptor.
//
// The version validator owns only version identity and Desired/Observed surface
// shape. Cross-version invariants such as exposure and canonical selection stay
// in definition_versions_validate.go.
func validateVersionDefinition(version VersionDefinition, resolver types.Resolver, path string) error {
	if err := version.version.Validate(); err != nil {
		return nestedVersionError(
			path+".version",
			ErrorReasonInvalidVersion,
			fmt.Sprintf("version %q is invalid", version.version),
			err,
		)
	}
	if version.desired.IsZero() {
		return versionError(
			path+".desired",
			ErrorReasonMissingDesired,
			detailVersionDesiredRequired,
		)
	}

	if err := validateSurface(
		version.desired,
		resolver,
		path+".desired",
		ErrorReasonInvalidDesired,
		ErrorReasonDesiredNotObject,
		detailDesiredObjectLikeTemplate,
	); err != nil {
		return err
	}

	if !version.observed.IsZero() {
		if err := validateSurface(
			version.observed,
			resolver,
			path+".observed",
			ErrorReasonInvalidObserved,
			ErrorReasonObservedNotObject,
			detailObservedObjectLikeTemplate,
		); err != nil {
			return err
		}
	}

	return nil
}

// invalidSurfaceDetail formats nested api/types failures for Desired/Observed.
func invalidSurfaceDetail(label string, err error) string {
	return fmt.Sprintf("%s descriptor is structurally invalid: %v", label, err)
}
