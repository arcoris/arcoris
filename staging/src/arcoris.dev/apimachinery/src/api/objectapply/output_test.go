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

	"arcoris.dev/apimachinery/api/meta/stamp"
)

func TestApplyOutputUsesLiveMetadata(t *testing.T) {
	req := testRequest()
	req.Live.ObjectMeta.ResourceVersion = stamp.ResourceVersion("rv-1")

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	if result.Object.ObjectMeta.ResourceVersion != "rv-1" {
		t.Fatalf("ResourceVersion = %q; want live metadata", result.Object.ObjectMeta.ResourceVersion)
	}
}

func TestApplyOutputUsesMergedDesired(t *testing.T) {
	result, err := Apply(testRequest(), Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
}

func TestApplyOutputPreservesLiveObserved(t *testing.T) {
	req := testRequest()
	req.Live = testObjectObserved(req.Live.Desired, obj(member("ready", str("true"))))
	req.Resource = testResourceWithObserved(desiredDescriptor())

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, *result.Object.Observed, "ready", "true")
}

func TestApplyOutputDoesNotMutateLive(t *testing.T) {
	req := testRequest()

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, req.Live.Desired, "image", "api:v1")
}

func TestApplyOutputDoesNotMutateApplied(t *testing.T) {
	req := testRequest()

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, req.Applied.Desired, "image", "api:v2")
}

func TestApplyOutputDoesNotMutateOwnership(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("old", path("$.replicas")))

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireOwners(t, req.Ownership.Desired().OwnersOf(path("$.replicas")), "old")
}

func TestBuildOutputObjectPreservesDesiredAndObservedValues(t *testing.T) {
	live := testObjectObserved(obj(member("image", str("old"))), obj(member("ready", str("true"))))
	desired := obj(member("image", str("new")))

	out := buildOutputObject(live, desired)

	requireStringMember(t, out.Desired, "image", "new")
	requireStringMember(t, *out.Observed, "ready", "true")
}
