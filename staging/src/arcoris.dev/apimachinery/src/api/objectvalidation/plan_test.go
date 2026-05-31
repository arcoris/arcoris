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
	"testing"

	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateRejectsInvalidPlanShape(t *testing.T) {
	tests := []struct {
		name   string
		plan   Plan[testDesired, testObserved]
		target error
		path   string
		reason ErrorReason
	}{
		{
			name:   "zero resource",
			plan:   Plan[testDesired, testObserved]{DesiredValidator: &spySurfaceValidator[testDesired]{}},
			target: ErrInvalidPlan,
			path:   pathPlanResource,
			reason: ErrorReasonInvalidPlan,
		},
		{
			name: "missing desired validator",
			plan: Plan[testDesired, testObserved]{
				Resource: resourceDefinition(resource.ScopeNamespaced),
			},
			target: ErrMissingValidator,
			path:   pathPlanDesiredValidator,
			reason: ErrorReasonMissingValidator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(validObject(), tt.plan)
			requireValidationError(t, err, tt.target, tt.path, tt.reason)
			if tt.target == ErrMissingValidator {
				requireErrorIs(t, err, ErrInvalidPlan)
			}
		})
	}
}

func TestValidateAllowsNilObservedValidatorUntilObservedValidationIsNeeded(t *testing.T) {
	plan := validPlan()
	plan.ObservedValidator = nil

	requireNoError(t, Validate(validObjectWithoutObserved(), plan))
}
