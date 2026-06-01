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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/meta"
)

func TestErrorFormatting(t *testing.T) {
	err := &Error{
		Record: diagnostic.NewRecord(
			"object.metadata",
			ErrInvalidObject,
			ErrorReasonInvalidMetadata,
			"metadata is invalid",
		),
	}

	got := err.Error()
	for _, want := range []string{
		"object",
		"object.metadata",
		ErrInvalidObject.Error(),
		string(ErrorReasonInvalidMetadata),
		"metadata is invalid",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Error() = %q, want segment %q", got, want)
		}
	}
}

func TestNestedErrorDetailDoesNotRepeatCause(t *testing.T) {
	cause := meta.TypeMeta{Kind: "Worker"}.Validate()
	err := nested("object.typeMeta", ErrInvalidObject, cause)

	var objectErr *Error
	if !errors.As(err, &objectErr) {
		t.Fatalf("errors.As(%T) = false", objectErr)
	}
	if objectErr.Detail != "metadata is invalid" {
		t.Fatalf("Detail = %q", objectErr.Detail)
	}
	if strings.Contains(objectErr.Detail, cause.Error()) {
		t.Fatalf("Detail repeats nested cause: %q", objectErr.Detail)
	}
}

func TestErrorUnwrapPreservesSentinelAndCause(t *testing.T) {
	cause := meta.TypeMeta{Kind: "Worker"}.Validate()
	err := nested("object.typeMeta", ErrInvalidObject, cause)

	if !errors.Is(err, ErrInvalidObject) {
		t.Fatalf("errors.Is(%v, ErrInvalidObject) = false", err)
	}
	if !errors.Is(err, meta.ErrInvalidTypeMeta) {
		t.Fatalf("errors.Is(%v, meta.ErrInvalidTypeMeta) = false", err)
	}
}
