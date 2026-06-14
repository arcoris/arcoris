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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/meta"
	apiobject "arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateDeterministicFailureOrderStopsLaterStages(t *testing.T) {
	desiredCause := errors.New("desired failed")

	tests := []struct {
		name       string
		obj        apiobject.Object[testDesired, testObserved]
		plan       Plan[testDesired, testObserved]
		desiredErr error
		wantTarget error
		wantPath   string
		wantReason ErrorReason
		wantCalls  []string
	}{
		{
			name: "invalid plan",
			obj:  validObject(),
			plan: Plan[testDesired, testObserved]{
				DesiredValidator: &spySurfaceValidator[testDesired]{},
			},
			wantTarget: ErrInvalidPlan,
			wantPath:   pathPlanResource,
			wantReason: ErrorReasonInvalidPlan,
		},
		{
			name: "invalid metadata",
			obj: apiobject.Object[testDesired, testObserved]{
				TypeMeta: validTypeMeta("v1"),
				ObjectMeta: meta.ObjectMeta{
					Name: "Worker",
				},
				Desired: testDesired{Replicas: 3},
			},
			plan:       validPlan(),
			wantTarget: ErrInvalidMetadata,
			wantPath:   "object.metadata",
			wantReason: ErrorReasonInvalidObjectMeta,
		},
		{
			name:       "resource mismatch",
			obj:        validObject(),
			plan:       Plan[testDesired, testObserved]{Resource: mismatchedResourceDefinition()},
			wantTarget: ErrResourceMismatch,
			wantPath:   pathObjectTypeMeta,
			wantReason: ErrorReasonResourceMismatch,
		},
		{
			name: "unknown version",
			obj:  validObject(),
			plan: Plan[testDesired, testObserved]{
				Resource: resourceDefinition(resource.ScopeNamespaced, versionWithObserved("v2")),
			},
			wantTarget: ErrVersionNotDefined,
			wantPath:   pathResourceVersions,
			wantReason: ErrorReasonVersionNotDefined,
		},
		{
			name:       "invalid scope",
			obj:        validObject(),
			plan:       Plan[testDesired, testObserved]{Resource: resourceDefinition(resource.ScopeGlobal)},
			wantTarget: ErrInvalidScope,
			wantPath:   pathObjectNamespace,
			wantReason: ErrorReasonInvalidScope,
		},
		{
			name:       "desired failure",
			obj:        validObject(),
			plan:       validPlan(),
			desiredErr: desiredCause,
			wantTarget: ErrInvalidDesired,
			wantPath:   pathObjectDesired,
			wantReason: ErrorReasonInvalidDesired,
			wantCalls:  []string{"desired"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := []string{}
			desired := &spySurfaceValidator[testDesired]{
				name:      "desired",
				callOrder: &calls,
				err:       tt.desiredErr,
			}
			observed := &spySurfaceValidator[testObserved]{
				name:      "observed",
				callOrder: &calls,
			}
			tt.plan.DesiredValidator = desired
			tt.plan.ObservedValidator = observed

			err := Validate(tt.obj, tt.plan)
			requireValidationError(t, err, tt.wantTarget, tt.wantPath, tt.wantReason)
			requireCallOrder(t, calls, tt.wantCalls...)
			if tt.desiredErr != nil {
				requireErrorIs(t, err, tt.desiredErr)
			}
		})
	}
}
