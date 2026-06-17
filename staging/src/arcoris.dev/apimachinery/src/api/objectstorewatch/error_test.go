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

package objectstorewatch

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestErrorUnwrapPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{Path: "collection_read.result", Reason: ErrorReasonInvalidCollectionRead, Cause: cause}

	requireErrorIs(t, err, cause)
}

func TestErrorStringHandlesNilCause(t *testing.T) {
	tests := []*Error{
		nil,
		{},
		{Reason: ErrorReasonInvalidCollectionRead},
		{Path: "collection_read.result", Reason: ErrorReasonInvalidCollectionRead},
	}

	for _, err := range tests {
		if got := err.Error(); got == "" {
			t.Fatalf("Error() returned empty string")
		}
	}
}

func TestErrorForExposesSentinelCauseAndDiagnostic(t *testing.T) {
	err := errorFor(
		"collection_read.result",
		ErrorReasonInvalidCollectionRead,
		ErrInvalidCollectionRead,
		objectstore.ErrInvalidListResult,
	)

	requireErrorIs(t, err, ErrInvalidCollectionRead)
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
	requireWatchError(t, err, ErrorReasonInvalidCollectionRead, "collection_read.result")
}

func TestErrorForDoesNotDuplicatePackagePrefix(t *testing.T) {
	err := errorFor(
		"collection_read.collection",
		ErrorReasonInvalidCollectionRead,
		ErrInvalidCollectionRead,
		objectstore.ErrInvalidListRequest,
	)
	text := err.Error()

	duplicatedPrefix := "objectstorewatch: " + "objectstorewatch."
	if strings.Contains(text, duplicatedPrefix) {
		t.Fatalf("Error() = %q; contains duplicated package path", text)
	}
	if !strings.Contains(text, "objectstorewatch: collection_read.collection") {
		t.Fatalf("Error() = %q; want package prefix plus local path", text)
	}
}

func TestErrorForFallsBackToSentinelCause(t *testing.T) {
	err := errorFor("boundary.collection", ErrorReasonInvalidBoundary, ErrInvalidBoundary, nil)

	requireErrorIs(t, err, ErrInvalidBoundary)
	requireWatchError(t, err, ErrorReasonInvalidBoundary, "boundary.collection")
}
