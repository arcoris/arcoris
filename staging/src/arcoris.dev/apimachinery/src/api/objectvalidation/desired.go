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

// validateDesired delegates desired payload validation to the typed validator.
//
// api/objectvalidation passes the exact desired value, selected version desired
// descriptor, and structural resolver through without reflection, defaulting,
// pruning, conversion, or mutation.
func validateDesired[D any, O any](
	value D,
	version resource.VersionDefinition,
	plan Plan[D, O],
) error {
	if plan.DesiredValidator == nil {
		return missingValidator(pathPlanDesiredValidator, "desired surface validator is required")
	}

	descriptor := version.Desired()
	if err := plan.DesiredValidator.ValidateSurface(value, descriptor, plan.Resolver); err != nil {
		return nested(
			pathObjectDesired,
			ErrInvalidDesired,
			ErrorReasonInvalidDesired,
			"desired surface is invalid",
			err,
		)
	}

	return nil
}
