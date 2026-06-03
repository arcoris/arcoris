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

func TestDeletableDroppedFieldsNoOtherOwner(t *testing.T) {
	got := deletableDroppedFields(
		state(entry("user", path("$.spec.image"))),
		owner("user"),
		fields(path("$.spec.image")),
	)

	requireSet(t, got, "$.spec.image")
}

func TestDeletableDroppedFieldsOtherExactOwnerPreserves(t *testing.T) {
	got := deletableDroppedFields(
		state(
			entry("user", path("$.spec.image")),
			entry("other", path("$.spec.image")),
		),
		owner("user"),
		fields(path("$.spec.image")),
	)

	requireSet(t, got)
}

func TestDeletableDroppedFieldsOtherAncestorOwnerPreserves(t *testing.T) {
	got := deletableDroppedFields(
		state(
			entry("user", path("$.spec.image")),
			entry("other", path("$.spec")),
		),
		owner("user"),
		fields(path("$.spec.image")),
	)

	requireSet(t, got)
}

func TestDeletableDroppedFieldsOtherDescendantOwnerPreserves(t *testing.T) {
	got := deletableDroppedFields(
		state(
			entry("user", path("$.spec")),
			entry("other", path("$.spec.image")),
		),
		owner("user"),
		fields(path("$.spec")),
	)

	requireSet(t, got)
}

func TestApplyDroppedOwnedFieldDeletesWhenUnownedByOthers(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("user", imagePath(), replicasPath()))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireNoMember(t, result.Value, "replicas")
	requireSet(t, result.DroppedFields, "$.replicas")
	requireSet(t, result.DeletedFields, "$.replicas")
}

func TestApplyDroppedOwnedFieldPreservedWhenOwnedByOther(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(
		entry("user", imagePath(), replicasPath()),
		entry("other", replicasPath()),
	)

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Value, "replicas", "3")
	requireSet(t, result.DeletedFields)
	requireOwners(t, result.Ownership.OwnersOf(replicasPath()), "other")
}

func TestApplyDroppedParentPreservedWhenOtherOwnerOwnsDescendant(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("name", str("old")), member("spec", obj(member("image", str("old"))))),
		Applied:    obj(member("name", str("new"))),
		Descriptor: typesObjectWithSpec(),
		Ownership: state(
			entry("user", specPath()),
			entry("other", specPath().Field("image")),
		),
	}

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	spec := requireMember(t, result.Value, "spec")
	requireStringMember(t, spec, "image", "old")
	requireSet(t, result.DeletedFields)
}

func TestApplyDroppedFieldDoesNotConflict(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(
		entry("user", imagePath(), replicasPath()),
		entry("other", replicasPath()),
	)

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireSet(t, result.Conflicts.AttemptedPaths())
	requireSet(t, result.DroppedFields, "$.replicas")
}
