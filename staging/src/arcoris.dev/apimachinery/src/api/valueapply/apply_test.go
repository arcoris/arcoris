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

import "testing"

func TestApplySelectedFieldNoConflict(t *testing.T) {
	result, err := Apply(specRequest(owner("user")), Options{})
	requireNoError(t, err)

	requireSet(t, result.Conflicts.AttemptedPaths())
	requireStringMember(t, result.Value, "image", "api:v2")
	requireStringMember(t, result.Value, "replicas", "3")
}

func TestApplyUpdatesValue(t *testing.T) {
	result, err := Apply(specRequest(owner("user")), Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Value, "image", "api:v2")
	requireSet(t, result.MergeFields, "$.image")
}

func TestApplyUpdatesOwnerFieldsToAppliedFields(t *testing.T) {
	result, err := Apply(specRequest(owner("user")), Options{})
	requireNoError(t, err)

	requireSet(t, result.Ownership.FieldsFor(owner("user")), "$.image")
}

func TestApplyDoesNotMutateLive(t *testing.T) {
	req := specRequest(owner("user"))

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, req.Live, "image", "api:v1")
}

func TestApplyDoesNotMutateApplied(t *testing.T) {
	req := specRequest(owner("user"))

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, req.Applied, "image", "api:v2")
}

func TestApplyDoesNotMutateOwnership(t *testing.T) {
	req := specRequest(owner("user"))

	_, err := Apply(req, Options{})
	requireNoError(t, err)

	requireSet(t, req.Ownership.FieldsFor(owner("user")))
}

func TestApplySameValueOwnedByOtherDoesNotConflict(t *testing.T) {
	req := specRequest(owner("user"))
	req.Applied = obj(member("image", str("api:v1")))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireSet(t, result.ChangedAppliedFields)
	requireSet(t, result.Conflicts.AttemptedPaths())
}

func TestApplySameValueOwnedByOtherCreatesSharedOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Applied = obj(member("image", str("api:v1")))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.OwnersOf(imagePath()), "other", "user")
}

func TestApplySameValueOwnedByOtherDoesNotRemoveOtherOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Applied = obj(member("image", str("api:v1")))
	req.Ownership = state(entry("other", imagePath(), replicasPath()))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireOwners(t, result.Ownership.OwnersOf(imagePath()), "other", "user")
	requireOwners(t, result.Ownership.OwnersOf(replicasPath()), "other")
}

func TestApplySameValueMapKeyCreatesSharedOwnership(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("app", str("same"))),
		Applied:    obj(member("app", str("same"))),
		Descriptor: mapDescriptor(),
		Ownership:  state(entry("other", labelPath())),
	}, Options{})
	requireNoError(t, err)

	requireSet(t, result.ChangedAppliedFields)
	requireOwners(t, result.Ownership.OwnersOf(labelPath()), "other", "user")
}

func TestApplyEmptyAppliedFieldsReleasesOwner(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("x", str("old"))),
		Applied:    obj(member("x", str("new"))),
		Descriptor: typesUnknownPruneObject(),
		Ownership:  state(entry("user", path("$.x"))),
	}

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireSet(t, result.AppliedFields)
	requireSet(t, result.Ownership.FieldsFor(owner("user")))
}
