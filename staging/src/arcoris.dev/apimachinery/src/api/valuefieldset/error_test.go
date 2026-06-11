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

package valuefieldset

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestErrorIsSentinel(t *testing.T) {
	_, err := ExtractOwnershipFields(value.StringValue("api"), types.Int64().Descriptor(), Options{})

	requireErrorIs(t, err, ErrKindMismatch)
}

func TestErrorIsBroadSentinels(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
	}{
		{
			name:   "invalid value",
			err:    errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "invalid"),
			target: ErrInvalidValue,
		},
		{
			name:   "invalid descriptor",
			err:    errorAt(fieldpath.Root(), ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "invalid"),
			target: ErrInvalidDescriptor,
		},
		{
			name:   "invalid path",
			err:    errorAt(fieldpath.Root(), ErrInvalidPath, ErrorReasonInvalidPath, "invalid"),
			target: ErrInvalidPath,
		},
		{
			name:   "kind mismatch",
			err:    errorAt(fieldpath.Root(), ErrKindMismatch, ErrorReasonKindMismatch, "invalid"),
			target: ErrKindMismatch,
		},
		{
			name:   "unknown field",
			err:    errorAt(fieldpath.Root(), ErrUnknownField, ErrorReasonUnknownField, "invalid"),
			target: ErrUnknownField,
		},
		{
			name:   "unresolved ref",
			err:    errorAt(fieldpath.Root(), ErrUnresolvedRef, ErrorReasonUnresolvedRef, "invalid"),
			target: ErrUnresolvedRef,
		},
		{
			name:   "reference cycle",
			err:    errorAt(fieldpath.Root(), ErrReferenceCycle, ErrorReasonReferenceCycle, "invalid"),
			target: ErrReferenceCycle,
		},
		{
			name:   "invalid list key",
			err:    errorAt(fieldpath.Root(), ErrInvalidListKey, ErrorReasonInvalidListKey, "invalid"),
			target: ErrInvalidListKey,
		},
		{
			name:   "duplicate list key",
			err:    errorAt(fieldpath.Root(), ErrDuplicateListKey, ErrorReasonDuplicateListKey, "invalid"),
			target: ErrDuplicateListKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireErrorIs(t, tt.err, tt.target)
		})
	}
}

func TestErrorAsValueFieldSetError(t *testing.T) {
	_, err := ExtractOwnershipFields(value.StringValue("api"), types.Int64().Descriptor(), Options{})

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if got.Path != "$" {
		t.Fatalf("Path = %q, want $", got.Path)
	}
	if got.Reason != ErrorReasonKindMismatch {
		t.Fatalf("Reason = %q, want %q", got.Reason, ErrorReasonKindMismatch)
	}
}

func TestErrorPath(t *testing.T) {
	path := rootField("spec", "replicas")

	_, err := ExtractOwnershipFieldsAt(path, value.StringValue("api"), types.Int64().Descriptor(), Options{})

	requireErrorPath(t, err, "$.spec.replicas")
}

func TestDuplicateListMapKeyErrorIncludesPhysicalOccurrences(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		readyConditionValue("True"),
		readyConditionValue("False"),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorDetailContains(t, err, "first occurrence at $.conditions[0]")
	requireErrorDetailContains(t, err, "duplicate at $.conditions[1]")
}
