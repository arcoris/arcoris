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

package valuemerge

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
)

func TestMergeListMapKeyErrorKindMissingKey(t *testing.T) {
	err, reason := mergeListMapKeyErrorKind(listmapkey.FailureMissingKey)

	if err != ErrInvalidListKey {
		t.Fatalf("err = %v; want %v", err, ErrInvalidListKey)
	}
	if reason != ErrorReasonMissingListKey {
		t.Fatalf("reason = %q; want %q", reason, ErrorReasonMissingListKey)
	}
}

func TestMergeListMapKeyErrorKindReferenceCycle(t *testing.T) {
	err, reason := mergeListMapKeyErrorKind(listmapkey.FailureReferenceCycle)

	if err != ErrReferenceCycle {
		t.Fatalf("err = %v; want %v", err, ErrReferenceCycle)
	}
	if reason != ErrorReasonReferenceCycle {
		t.Fatalf("reason = %q; want %q", reason, ErrorReasonReferenceCycle)
	}
}

func TestMergeListMapKeyUnexpectedErrorWrapped(t *testing.T) {
	err := mergeListMapKeyError(
		root().Index(0),
		errors.New("unexpected failure"),
	)

	requireErrorIs(t, err, ErrInvalidListKey)

	var mergeError *Error
	if !errors.As(err, &mergeError) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if mergeError.Path != root().Index(0).String() {
		t.Fatalf("path = %s; want %s", mergeError.Path, root().Index(0))
	}
}

func TestMergeListMapKeyErrorUsesSharedErrorPath(t *testing.T) {
	err := mergeListMapKeyError(
		root().Index(0),
		&listmapkey.Error{
			Path:   root().Index(1),
			Kind:   listmapkey.FailureMissingKey,
			Detail: "missing key",
		},
	)

	requireErrorIs(t, err, ErrInvalidListKey)

	var mergeError *Error
	if !errors.As(err, &mergeError) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if mergeError.Path != root().Index(1).String() {
		t.Fatalf("path = %s; want %s", mergeError.Path, root().Index(1))
	}
}

func TestMergeListMapKeyUnexpectedErrorAcceptsAnyPath(t *testing.T) {
	err := mergeListMapKeyError(
		fieldpath.Root(),
		errors.New("unexpected failure"),
	)

	requireErrorIs(t, err, ErrInvalidListKey)
}
