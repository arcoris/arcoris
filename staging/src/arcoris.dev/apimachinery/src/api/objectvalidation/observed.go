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

package objectvalidation

import "arcoris.dev/apimachinery/api/resource"

// validateObserved enforces observed descriptor presence and delegates values.
//
// Observed may be absent even when the resource defines an observed descriptor.
// Requiring observed state is a persistence/runtime concern, not baseline
// resource contract validation.
func validateObserved[D any, O any](
	value *O,
	version resource.VersionDefinition,
	plan Plan[D, O],
) error {
	descriptor, ok := version.Observed()
	if !ok {
		if value != nil {
			return errorf(
				pathObjectObserved,
				ErrObservedNotAllowed,
				ErrorReasonObservedNotAllowed,
				"resource version %q does not define an observed surface",
				version.Version(),
			)
		}

		return nil
	}

	if value == nil {
		return nil
	}

	if plan.ObservedValidator == nil {
		return missingValidator(pathPlanObservedValidator, "observed surface validator is required")
	}

	if err := plan.ObservedValidator.ValidateSurface(*value, descriptor, plan.Resolver); err != nil {
		return nested(
			pathObjectObserved,
			ErrInvalidObserved,
			ErrorReasonInvalidObserved,
			"observed surface is invalid",
			err,
		)
	}

	return nil
}
