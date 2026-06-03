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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valueapply"
)

func TestApplyDesiredForce(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
	requireOwners(t, result.Ownership.Desired().OwnersOf(path("$.image")), "user")
}

func TestApplyDesiredForceKeepsPreForceConflictsInResult(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	if result.Desired.Conflicts.Len() != 1 {
		t.Fatalf("conflicts = %d; want 1", result.Desired.Conflicts.Len())
	}
}

func TestApplyDesiredForceUpdatesDesiredOwnership(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.Desired().OwnersOf(path("$.image")), "user")
}

func TestApplyDesiredUnsupportedForceTakeoverPropagatesValueApplyCause(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", fieldpath.RootPath()))

	_, err := Apply(req, Options{Force: true})

	requireErrorIs(t, err, ErrDesiredApplyFailed)
	requireErrorIs(t, err, valueapply.ErrUnsupportedTakeover)
}
