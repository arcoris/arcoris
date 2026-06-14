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
)

func TestApplySameOwnerSameFieldSucceeds(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("user", path("$.image")))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
	requireOwnersOf(t, result.Ownership.Desired(), path("$.image"), "user")
}

func TestApplyDifferentOwnerDisjointFieldSucceeds(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.replicas")))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
	requireStringMember(t, result.Object.Desired, "replicas", "3")
	requireOwnersOf(t, result.Ownership.Desired(), path("$.image"), "user")
	requireOwnersOf(t, result.Ownership.Desired(), path("$.replicas"), "other")
}

func TestApplyDifferentOwnerSameChangedFieldConflictDetailsAreStable(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrConflict)
	requireSet(t, result.Desired.Conflicts.AttemptedPaths(), "$.image")

	var conflictErr *fieldownership.ConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("error does not wrap fieldownership.ConflictError: %v", err)
	}
	conflicts := conflictErr.Conflicts().Conflicts()
	if len(conflicts) != 1 {
		t.Fatalf("conflict count = %d; want 1", len(conflicts))
	}
	if conflicts[0].Owner != owner("other") {
		t.Fatalf("conflict owner = %q; want other", conflicts[0].Owner)
	}
	if !conflicts[0].OwnedPath.Equal(path("$.image")) {
		t.Fatalf("owned path = %s; want $.image", conflicts[0].OwnedPath)
	}
	if !conflicts[0].AttemptedPath.Equal(path("$.image")) {
		t.Fatalf("attempted path = %s; want $.image", conflicts[0].AttemptedPath)
	}
}
