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
	"arcoris.dev/apimachinery/api/valueapply"
)

func TestApplyDesiredConflict(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyDesiredConflictReturnsPartialDesiredResult(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrConflict)

	requireSet(t, result.Desired.AppliedFields, "$.image")
	requireSet(t, result.Desired.ChangedAppliedFields, "$.image")
	if result.Desired.Conflicts.Len() != 1 {
		t.Fatalf("conflicts = %d; want 1", result.Desired.Conflicts.Len())
	}
	if !result.Object.Desired.IsZero() {
		t.Fatalf("object was built")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("ownership was updated")
	}
}

func TestApplyDesiredConflictWrapsObjectApplyConflict(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyDesiredConflictWrapsValueApplyConflict(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, valueapply.ErrConflict)
}

func TestApplyDesiredConflictWrapsFieldOwnershipConflict(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, fieldownership.ErrConflict)

	var conflictErr *fieldownership.ConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("error does not wrap fieldownership.ConflictError: %v", err)
	}
}
