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

	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateObservedPresence(t *testing.T) {
	t.Run("no descriptor nil observed", func(t *testing.T) {
		plan := validPlan()
		plan.Resource = resourceDefinition(
			resource.ScopeNamespaced,
			versionWithoutObserved("v1"),
		)

		requireNoError(t, Validate(validObjectWithoutObserved(), plan))
	})

	t.Run("no descriptor with observed", func(t *testing.T) {
		plan := validPlan()
		plan.Resource = resourceDefinition(
			resource.ScopeNamespaced,
			versionWithoutObserved("v1"),
		)

		err := Validate(validObject(), plan)
		requireValidationError(
			t,
			err,
			ErrObservedNotAllowed,
			pathObjectObserved,
			ErrorReasonObservedNotAllowed,
		)
	})

	t.Run("descriptor nil observed", func(t *testing.T) {
		requireNoError(t, Validate(validObjectWithoutObserved(), validPlan()))
	})
}

func TestValidateObservedSurface(t *testing.T) {
	plan, _, observed := validPlanWithSpies()
	resolver := &testResolver{name: "observed"}
	plan.Resolver = resolver
	plan.Resource = resourceDefinition(
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desiredDescriptor(), resource.Observed(alternateDescriptor())),
	)

	obj := validObject()
	requireNoError(t, Validate(obj, plan))

	if observed.called != 1 {
		t.Fatalf("observed validator called %d times, want 1", observed.called)
	}
	if observed.value != *obj.Observed {
		t.Fatalf("observed value = %#v, want %#v", observed.value, *obj.Observed)
	}
	requireTypeEqual(t, observed.typ, alternateDescriptor())
	if observed.resolver != resolver {
		t.Fatalf("resolver = %#v, want %#v", observed.resolver, resolver)
	}
}

func TestValidateObservedReceivesNilResolver(t *testing.T) {
	plan, _, observed := validPlanWithSpies()
	plan.Resolver = nil

	requireNoError(t, Validate(validObject(), plan))

	if observed.resolver != nil {
		t.Fatalf("observed resolver = %#v, want nil", observed.resolver)
	}
}

func TestValidateObservedNeedsValidatorOnlyWhenValueIsPresent(t *testing.T) {
	plan := validPlan()
	plan.ObservedValidator = nil

	requireNoError(t, Validate(validObjectWithoutObserved(), plan))

	err := Validate(validObject(), plan)
	requireValidationError(
		t,
		err,
		ErrMissingValidator,
		pathPlanObservedValidator,
		ErrorReasonMissingValidator,
	)
	requireErrorIs(t, err, ErrInvalidPlan)
	requireErrorNotIs(t, err, ErrInvalidObject)
}

func TestValidateWrapsObservedSurfaceError(t *testing.T) {
	cause := errors.New("observed failed")
	plan, _, observed := validPlanWithSpies()
	observed.err = cause

	err := Validate(validObject(), plan)
	validationErr := requireValidationError(
		t,
		err,
		ErrInvalidObserved,
		pathObjectObserved,
		ErrorReasonInvalidObserved,
	)
	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, cause)
	if validationErr.Cause != cause {
		t.Fatalf("Cause = %#v, want %#v", validationErr.Cause, cause)
	}
}
