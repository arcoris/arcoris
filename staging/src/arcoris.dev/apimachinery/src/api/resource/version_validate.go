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
	"fmt"

	"arcoris.dev/apimachinery/api/types"
)

// validateVersionDefinitionLocal checks one version descriptor locally.
//
// The version validator owns only version identity and Desired/Observed surface
// shape. Cross-version invariants such as exposure and canonical selection stay
// in definition_versions_validate.go.
func validateVersionDefinitionLocal(version VersionDefinition, path string) error {
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

	if err := validateSurfaceLocal(
		version.desired,
		path+".desired",
		ErrorReasonInvalidDesired,
		ErrorReasonDesiredNotObject,
		detailDesiredObjectLikeTemplate,
	); err != nil {
		return err
	}

	if !version.observed.IsZero() {
		if err := validateSurfaceLocal(
			version.observed,
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

// validateVersionDefinitionResolved checks one version descriptor with a
// resolver.
//
// This is the resolved counterpart to validateVersionDefinitionLocal: it owns
// version identity, Desired presence, and resolver-proven Desired/Observed root
// object shape, but still delegates full descriptor graph validation to
// api/types.
func validateVersionDefinitionResolved(version VersionDefinition, resolver types.Resolver, path string) error {
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

	if err := validateSurfaceResolved(
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
		if err := validateSurfaceResolved(
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

// invalidSurfaceDetail summarizes nested api/types failures for Desired/Observed.
//
// The nested error remains attached as Cause on the returned resource Error, so
// callers can still use errors.As for the original *types.DescriptorError. The detail
// string keeps only the most useful structured fields to avoid repeating the
// entire nested Error text in CLI/editor diagnostics.
func invalidSurfaceDetail(label string, err error) string {
	var typeErr *types.DescriptorError
	if errors.As(err, &typeErr) {
		if typeErr.Detail != "" {
			return fmt.Sprintf(
				"%s descriptor is structurally invalid at %s: %s: %s",
				label,
				typeErr.Path,
				typeErr.Reason,
				typeErr.Detail,
			)
		}

		if typeErr.Reason != "" {
			return fmt.Sprintf(
				"%s descriptor is structurally invalid at %s: %s",
				label,
				typeErr.Path,
				typeErr.Reason,
			)
		}

		if typeErr.Path != "" {
			return fmt.Sprintf("%s descriptor is structurally invalid at %s", label, typeErr.Path)
		}
	}

	return fmt.Sprintf("%s descriptor is structurally invalid: %v", label, err)
}
