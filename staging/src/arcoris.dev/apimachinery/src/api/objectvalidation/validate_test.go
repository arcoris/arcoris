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

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestValidatorMatchesPackageValidate(t *testing.T) {
	obj := validObject()
	plan := validPlan()

	errFromFunction := Validate(obj, plan)
	errFromValidator := New(plan).Validate(obj)

	requireNoError(t, errFromFunction)
	requireNoError(t, errFromValidator)
}

func TestValidatorCanBeReused(t *testing.T) {
	plan, desired, observed := validPlanWithSpies()
	validator := New(plan)

	requireNoError(t, validator.Validate(validObject()))
	requireNoError(t, validator.Validate(validObjectWithoutObserved()))

	if desired.called != 2 {
		t.Fatalf("desired validator called %d times, want 2", desired.called)
	}
	if observed.called != 1 {
		t.Fatalf("observed validator called %d times, want 1", observed.called)
	}
}

func TestValidatorStoresPlanByValue(t *testing.T) {
	plan := validPlan()
	validator := New(plan)
	plan.Resource = mismatchedResourceDefinition()

	requireNoError(t, validator.Validate(validObject()))
}

func TestValidatorDoesNotMutateObject(t *testing.T) {
	obj := validObject()
	observedPtr := obj.Observed
	observedValue := *obj.Observed

	requireNoError(t, New(validPlan()).Validate(obj))

	if obj.TypeMeta != validTypeMeta("v1") {
		t.Fatalf("TypeMeta = %#v, want unchanged", obj.TypeMeta)
	}
	if obj.ObjectMeta.Namespace != metaidentity.Namespace("system") ||
		obj.ObjectMeta.Name != "worker" ||
		obj.ObjectMeta.UID != "uid-1" {
		t.Fatalf("ObjectMeta = %#v, want unchanged", obj.ObjectMeta)
	}
	if obj.Desired != (testDesired{Replicas: 3}) {
		t.Fatalf("Desired = %#v, want unchanged", obj.Desired)
	}
	if obj.Observed != observedPtr {
		t.Fatalf("Observed pointer = %#v, want %#v", obj.Observed, observedPtr)
	}
	if *obj.Observed != observedValue {
		t.Fatalf("Observed value = %#v, want %#v", *obj.Observed, observedValue)
	}
}

func TestValidateEndToEnd(t *testing.T) {
	plan, desired, observed := validPlanWithSpies()

	requireNoError(t, Validate(validObject(), plan))
	if desired.called != 1 {
		t.Fatalf("desired validator called %d times, want 1", desired.called)
	}
	if observed.called != 1 {
		t.Fatalf("observed validator called %d times, want 1", observed.called)
	}
}

func TestValidateHappyPathCallOrder(t *testing.T) {
	calls := []string{}

	desired := &spySurfaceValidator[testDesired]{
		name:      "desired",
		callOrder: &calls,
	}
	observed := &spySurfaceValidator[testObserved]{
		name:      "observed",
		callOrder: &calls,
	}

	plan := validPlan()
	plan.DesiredValidator = desired
	plan.ObservedValidator = observed

	requireNoError(t, Validate(validObject(), plan))
	requireCallOrder(t, calls, "desired", "observed")
}

func TestValidateFailureOrder(t *testing.T) {
	calls := []string{}
	cause := errors.New("desired failed")

	desired := &spySurfaceValidator[testDesired]{
		name:      "desired",
		callOrder: &calls,
		err:       cause,
	}
	observed := &spySurfaceValidator[testObserved]{
		name:      "observed",
		callOrder: &calls,
	}
	plan := validPlan()
	plan.DesiredValidator = desired
	plan.ObservedValidator = observed

	err := Validate(validObject(), plan)
	requireValidationError(
		t,
		err,
		ErrInvalidDesired,
		pathObjectDesired,
		ErrorReasonInvalidDesired,
	)
	requireErrorIs(t, err, cause)

	requireCallOrder(t, calls, "desired")
}
