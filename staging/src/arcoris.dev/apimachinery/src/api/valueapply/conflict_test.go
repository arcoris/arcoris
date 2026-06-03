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

package valueapply

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
)

func TestApplyChangedFieldOwnedByOtherConflicts(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	_, err := Apply(req, Options{})

	requireErrorIs(t, err, ErrConflict)
	requireErrorIs(t, err, fieldownership.ErrConflict)
}

func TestApplyConflictReturnsPartialResult(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrConflict)

	requireSet(t, result.AppliedFields, "$.image")
	requireSet(t, result.ChangedAppliedFields, "$.image")
	if result.Conflicts.Len() != 1 {
		t.Fatalf("conflicts = %d; want 1", result.Conflicts.Len())
	}
}

func TestApplyConflictDoesNotMerge(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrConflict)

	if !result.Value.IsZero() {
		t.Fatalf("value was merged")
	}
}

func TestApplyConflictDoesNotUpdateOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrConflict)

	if !result.Ownership.IsEmpty() {
		t.Fatalf("ownership was updated")
	}
}

func TestApplyConflictErrorWrapsFieldOwnershipConflict(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	_, err := Apply(req, Options{})

	var conflictError *fieldownership.ConflictError
	if !errors.As(err, &conflictError) {
		t.Fatalf("error does not wrap fieldownership.ConflictError: %v", err)
	}
}

func TestApplyChangedFieldOwnedByOtherForceTakesOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.OwnersOf(imagePath()), "user")
	requireStringMember(t, result.Value, "image", "api:v2")
}

func TestApplyForceRemovesOnlyConflictingOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(
		entry("other", imagePath(), replicasPath()),
	)

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.OwnersOf(imagePath()), "user")
	requireOwners(t, result.Ownership.OwnersOf(replicasPath()), "other")
}

func TestApplyForceDoesNotRemoveUnrelatedOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(
		entry("other", replicasPath()),
	)

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.OwnersOf(replicasPath()), "other")
}

func TestOwnershipConflictsWrapsInvalidOwner(t *testing.T) {
	req := specRequest(owner(" "))

	_, err := ownershipConflicts(req, fields(imagePath()))

	requireErrorIs(t, err, ErrConflict)
}
