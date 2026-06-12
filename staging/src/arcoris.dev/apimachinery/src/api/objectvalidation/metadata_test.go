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
)

func TestValidateWrapsMetadataErrors(t *testing.T) {
	tests := []struct {
		name   string
		obj    apiobject.Object[testDesired, testObserved]
		target error
		path   string
		reason ErrorReason
	}{
		{
			name: "type meta",
			obj: apiobject.Object[testDesired, testObserved]{
				TypeMeta:   meta.TypeMeta{Kind: "Worker"},
				ObjectMeta: validObjectMeta("system"),
				Desired:    testDesired{Replicas: 3},
			},
			target: meta.ErrInvalidTypeMeta,
			path:   pathObjectTypeMeta,
			reason: ErrorReasonInvalidTypeMeta,
		},
		{
			name: "object meta",
			obj: apiobject.Object[testDesired, testObserved]{
				TypeMeta:   validTypeMeta("v1"),
				ObjectMeta: meta.ObjectMeta{Name: "Worker"},
				Desired:    testDesired{Replicas: 3},
			},
			target: meta.ErrInvalidObjectMeta,
			path:   "object.metadata",
			reason: ErrorReasonInvalidObjectMeta,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.obj, validPlan())
			validationErr := requireValidationError(
				t,
				err,
				ErrInvalidMetadata,
				tt.path,
				tt.reason,
			)
			requireErrorIs(t, err, ErrInvalidObject)
			requireErrorIs(t, err, apiobject.ErrInvalidObject)
			requireErrorIs(t, err, tt.target)
			if validationErr.Cause == nil {
				t.Fatal("Error.Cause is nil")
			}
		})
	}
}

func TestMetadataDiagnosticFallsBackForNonObjectErrors(t *testing.T) {
	path, reason, detail := metadataDiagnostic(errors.New("metadata failed"))

	if path != pathObject {
		t.Fatalf("path = %q, want %q", path, pathObject)
	}
	if reason != ErrorReasonInvalidMetadata {
		t.Fatalf("reason = %q, want %q", reason, ErrorReasonInvalidMetadata)
	}
	if detail != "object metadata is invalid" {
		t.Fatalf("detail = %q, want %q", detail, "object metadata is invalid")
	}
}

func TestValidateValidMetadataReachesResourceMatch(t *testing.T) {
	plan := validPlan()
	plan.Resource = mismatchedResourceDefinition()

	err := Validate(validObject(), plan)
	requireValidationError(
		t,
		err,
		ErrResourceMismatch,
		pathObjectTypeMeta,
		ErrorReasonResourceMismatch,
	)
}
