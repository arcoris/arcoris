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

package object

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/meta"
)

func TestObjectValidateMeta(t *testing.T) {
	obj := New[testDesired, testObserved](
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
	)

	requireNoError(t, obj.ValidateMeta())
}

func TestObjectValidateMetaRejectsInvalidMetadata(t *testing.T) {
	tests := []struct {
		name   string
		obj    Object[testDesired, testObserved]
		target error
		path   string
		reason ErrorReason
	}{
		{
			name: "type meta",
			obj: Object[testDesired, testObserved]{
				TypeMeta:   meta.TypeMeta{Kind: "Worker"},
				ObjectMeta: validObjectMeta(),
				Desired:    testDesired{Replicas: 3},
			},
			target: meta.ErrInvalidTypeMeta,
			path:   "object.typeMeta",
			reason: ErrorReasonInvalidTypeMeta,
		},
		{
			name: "object meta",
			obj: Object[testDesired, testObserved]{
				TypeMeta:   validTypeMeta(),
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
			err := tt.obj.ValidateMeta()
			requireErrorIs(t, err, ErrInvalidObject)
			requireErrorIs(t, err, tt.target)

			var objectErr *Error
			if !errors.As(err, &objectErr) {
				t.Fatalf("errors.As(%T) = false", objectErr)
			}
			if objectErr.Path != tt.path {
				t.Fatalf("Path = %q, want %q", objectErr.Path, tt.path)
			}
			if objectErr.Reason != tt.reason {
				t.Fatalf("Reason = %q", objectErr.Reason)
			}
			if objectErr.Cause == nil {
				t.Fatal("Cause = nil")
			}
			if objectErr.Detail == "" {
				t.Fatal("Detail = empty")
			}
		})
	}
}

func TestObjectValidateMetaDoesNotInspectPayloads(t *testing.T) {
	desiredCalled := false
	observedCalled := false
	obj := NewObserved(
		validTypeMeta(),
		validObjectMeta(),
		payloadWithValidate{Called: &desiredCalled},
		payloadWithValidate{Called: &observedCalled},
	)

	requireNoError(t, obj.ValidateMeta())
	if desiredCalled {
		t.Fatal("ValidateMeta called desired payload Validate")
	}
	if observedCalled {
		t.Fatal("ValidateMeta called observed payload Validate")
	}
}
