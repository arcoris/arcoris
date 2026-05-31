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
)

func TestErrorfBuildsStructuredError(t *testing.T) {
	err := errorf(
		pathResourceVersions,
		ErrVersionNotDefined,
		ErrorReasonVersionNotDefined,
		"resource version %q is not defined",
		"v2",
	)

	validationErr := requireValidationError(
		t,
		err,
		ErrVersionNotDefined,
		pathResourceVersions,
		ErrorReasonVersionNotDefined,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	if validationErr.Detail != `resource version "v2" is not defined` {
		t.Fatalf("Error.Detail = %q", validationErr.Detail)
	}
	if validationErr.Cause != nil {
		t.Fatalf("Error.Cause = %#v, want nil", validationErr.Cause)
	}
}

func TestNestedBuildsStructuredErrorWithCause(t *testing.T) {
	cause := errors.New("metadata failure")

	err := nested(
		pathObjectTypeMeta,
		ErrInvalidMetadata,
		ErrorReasonInvalidMetadata,
		"object type metadata is invalid",
		cause,
	)

	validationErr := requireValidationError(
		t,
		err,
		ErrInvalidMetadata,
		pathObjectTypeMeta,
		ErrorReasonInvalidMetadata,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, cause)
	if validationErr.Detail != "object type metadata is invalid" {
		t.Fatalf("Error.Detail = %q", validationErr.Detail)
	}
	if validationErr.Cause != cause {
		t.Fatalf("Error.Cause = %#v, want %#v", validationErr.Cause, cause)
	}
}

func TestMissingValidatorBuildsPlanError(t *testing.T) {
	err := missingValidator(
		pathPlanObservedValidator,
		"observed surface validator is required",
	)

	validationErr := requireValidationError(
		t,
		err,
		ErrMissingValidator,
		pathPlanObservedValidator,
		ErrorReasonMissingValidator,
	)

	requireErrorIs(t, err, ErrInvalidPlan)
	requireErrorNotIs(t, err, ErrInvalidObject)
	if validationErr.Detail != "observed surface validator is required" {
		t.Fatalf("Error.Detail = %q", validationErr.Detail)
	}
	if validationErr.Cause != nil {
		t.Fatalf("Error.Cause = %#v, want nil", validationErr.Cause)
	}
}
