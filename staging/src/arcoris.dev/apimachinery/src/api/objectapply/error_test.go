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

package objectapply

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valueapply"
)

func TestErrorIsSentinel(t *testing.T) {
	err := errorAt(pathObject, ErrInvalidRequest, ErrorReasonInvalidRequest, "bad request")

	requireErrorIs(t, err, ErrInvalidRequest)
}

func TestErrorAsObjectApplyError(t *testing.T) {
	err := errorAt(pathObject, ErrInvalidRequest, ErrorReasonInvalidRequest, "bad request")

	var applyErr *Error
	if !errors.As(err, &applyErr) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if applyErr.Reason != ErrorReasonInvalidRequest {
		t.Fatalf("reason = %q; want %q", applyErr.Reason, ErrorReasonInvalidRequest)
	}
}

func TestApplyUnsupportedObservedApplyError(t *testing.T) {
	req := testRequest()
	req.Applied = req.Applied.WithObserved(obj(member("ready", str("false"))))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedObservedApply)
}

func TestApplyUnsupportedMetadataChangeError(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.ResourceVersion = "rv-applied"

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrUnsupportedMetadataChange)
}

func TestApplyIdentityMismatchError(t *testing.T) {
	req := testRequest()
	req.Applied.ObjectMeta.Name = "other"

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrIdentityMismatch)
}

func TestApplyVersionMismatchError(t *testing.T) {
	req := testRequest()
	req.Applied.TypeMeta = testTypeMeta("v2")

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrVersionMismatch)
}

func TestApplyDesiredApplyFailureWrapsDesiredApplyFailed(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", fieldpath.RootPath()))

	_, err := Apply(req, Options{Force: true})

	requireErrorIs(t, err, ErrDesiredApplyFailed)
	requireErrorIs(t, err, valueapply.ErrUnsupportedTakeover)
}

func TestApplyDesiredConflictPreservesLowerLevelCauses(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrConflict)
	requireErrorIs(t, err, valueapply.ErrConflict)
	requireErrorIs(t, err, fieldownership.ErrConflict)
}

func TestNilErrorString(t *testing.T) {
	var err *Error

	if err.Error() != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", err.Error())
	}
}

func TestNilErrorUnwrap(t *testing.T) {
	var err *Error

	if err.Unwrap() != nil {
		t.Fatalf("Unwrap() is not nil")
	}
}
