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

import (
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
)

// Plan contains already-resolved dependencies for object contract validation.
//
// Resource is the contract to validate against. Resolver is passed through to
// surface validators so they can resolve structural TypeRef descriptors.
// DesiredValidator is required because desired is a required resource surface.
// ObservedValidator is needed only when the object carries observed data and
// the selected resource version defines an observed descriptor.
type Plan[D any, O any] struct {
	Resource resource.Definition
	Resolver types.Resolver

	DesiredValidator  SurfaceValidator[D]
	ObservedValidator SurfaceValidator[O]
}

// validatePlanShape checks only the dependencies needed before object work.
//
// Resource definitions are expected to be validated at catalog or construction
// boundaries. This package does not repeatedly validate the full resource
// descriptor graph for every object.
func validatePlanShape[D any, O any](plan Plan[D, O]) error {
	if plan.Resource.IsZero() {
		return errorf(
			pathPlanResource,
			ErrInvalidPlan,
			ErrorReasonInvalidPlan,
			"resource definition is required",
		)
	}

	if plan.DesiredValidator == nil {
		return missingValidator(pathPlanDesiredValidator, "desired surface validator is required")
	}

	return nil
}
