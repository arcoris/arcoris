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

func TestValidateDesiredSurface(t *testing.T) {
	plan, desired, _ := validPlanWithSpies()
	resolver := &testResolver{name: "desired"}
	plan.Resolver = resolver
	plan.Resource = resourceDefinition(
		resource.ScopeNamespaced,
		resource.NewVersion("v1", alternateDescriptor(), resource.Observed(observedDescriptor())),
	)

	obj := validObject()
	requireNoError(t, Validate(obj, plan))

	if desired.called != 1 {
		t.Fatalf("desired validator called %d times, want 1", desired.called)
	}
	if desired.value != obj.Desired {
		t.Fatalf("desired value = %#v, want %#v", desired.value, obj.Desired)
	}
	requireTypeEqual(t, desired.typ, alternateDescriptor())
	if desired.resolver != resolver {
		t.Fatalf("resolver = %#v, want %#v", desired.resolver, resolver)
	}
}

func TestValidateWrapsDesiredSurfaceError(t *testing.T) {
	cause := errors.New("desired failed")
	plan, desired, observed := validPlanWithSpies()
	desired.err = cause

	err := Validate(validObject(), plan)
	validationErr := requireValidationError(
		t,
		err,
		ErrInvalidDesired,
		pathObjectDesired,
		ErrorReasonInvalidDesired,
	)
	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, cause)
	if validationErr.Cause != cause {
		t.Fatalf("Cause = %#v, want %#v", validationErr.Cause, cause)
	}
	if observed.called != 0 {
		t.Fatalf("observed validator called after desired failure: %d", observed.called)
	}
}

func TestValidateDesiredValidatorNotCalledAfterEarlierFailure(t *testing.T) {
	tests := []struct {
		name string
		obj  apiobject.Object[testDesired, testObserved]
		plan Plan[testDesired, testObserved]
	}{
		{
			name: "metadata",
			obj: apiobject.Object[testDesired, testObserved]{
				TypeMeta: validTypeMeta("v1"),
				ObjectMeta: meta.ObjectMeta{
					Name: "Worker",
				},
				Desired: testDesired{Replicas: 3},
			},
			plan: validPlan(),
		},
		{
			name: "resource mismatch",
			obj:  validObject(),
			plan: Plan[testDesired, testObserved]{
				Resource:         mismatchedResourceDefinition(),
				DesiredValidator: &spySurfaceValidator[testDesired]{},
			},
		},
		{
			name: "missing version",
			obj:  validObject(),
			plan: Plan[testDesired, testObserved]{
				Resource:         resourceDefinition(resource.ScopeNamespaced, versionWithObserved("v2")),
				DesiredValidator: &spySurfaceValidator[testDesired]{},
			},
		},
		{
			name: "scope",
			obj:  validObject(),
			plan: Plan[testDesired, testObserved]{
				Resource:         resourceDefinition(resource.ScopeGlobal),
				DesiredValidator: &spySurfaceValidator[testDesired]{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desired := &spySurfaceValidator[testDesired]{}
			tt.plan.DesiredValidator = desired

			err := Validate(tt.obj, tt.plan)
			if err == nil {
				t.Fatal("Validate() = nil")
			}
			if desired.called != 0 {
				t.Fatalf("desired validator called after earlier failure: %d", desired.called)
			}
		})
	}
}
