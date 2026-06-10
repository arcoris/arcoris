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
	"arcoris.dev/apimachinery/api/types"
)

func TestApplyForceExactConflictTakesOwnership(t *testing.T) {
	req := specRequest(owner("user"))
	req.Ownership = state(entry("other", imagePath()))

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	requireOwnersOf(t, result.Ownership, imagePath(), "user")
	requireStringMember(t, result.Value, "image", "api:v2")
}

func TestApplyForceOtherOwnerAncestorConflictUnsupported(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("spec", obj(member("image", str("old"))))),
		Applied:    obj(member("spec", obj(member("image", str("new"))))),
		Descriptor: typesObjectWithSpec(),
		Ownership:  state(entry("other", specPath())),
	}

	result, err := Apply(req, Options{Force: true})
	requireErrorIs(t, err, ErrUnsupportedTakeover)
	requireErrorIs(t, err, fieldownership.ErrConflict)

	if !result.Value.IsZero() {
		t.Fatalf("value was merged")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("ownership was updated")
	}
	requireSet(t, result.ChangedAppliedFields, "$.spec.image")
}

func TestApplyForceOtherOwnerDescendantConflictTakesOwnership(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("spec", obj(member("image", str("old"))))),
		Applied:    obj(member("spec", obj())),
		Descriptor: typesObjectWithSpec(),
		Ownership:  state(entry("other", specPath().Field(testFieldName("image")))),
	}

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	spec := requireMember(t, result.Value, "spec")
	requireNoMember(t, spec, "image")
	requireOwnersOf(t, result.Ownership, specPath(), "user")
	requireOwnersOf(t, result.Ownership, specPath().Field(testFieldName("image")))
}

func TestApplyForceListMapItemOwnerAndFieldAttemptUnsupported(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(readyCondition("False")),
		Applied:    list(readyCondition("True")),
		Descriptor: conditionsDescriptor(),
		Ownership:  state(entry("other", root().Select(readySelector()))),
	}

	_, err := Apply(req, Options{Force: true})

	requireErrorIs(t, err, ErrUnsupportedTakeover)
	requireErrorIs(t, err, fieldownership.ErrConflict)
}

func TestApplyUnsupportedForceTakeoverReturnsPartialResult(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("spec", obj(member("image", str("old"))))),
		Applied:    obj(member("spec", obj(member("image", str("new"))))),
		Descriptor: typesObjectWithSpec(),
		Ownership:  state(entry("other", specPath())),
	}

	result, err := Apply(req, Options{Force: true})
	requireErrorIs(t, err, ErrUnsupportedTakeover)

	requireSet(t, result.AppliedFields, "$.spec.image")
	requireSet(t, result.ChangedAppliedFields, "$.spec.image")
	if result.Conflicts.Len() != 1 {
		t.Fatalf("conflicts = %d; want 1", result.Conflicts.Len())
	}
}

func TestApplyUnsupportedForceTakeoverPreservesFieldOwnershipConflictCause(t *testing.T) {
	req := Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("spec", obj(member("image", str("old"))))),
		Applied:    obj(member("spec", obj(member("image", str("new"))))),
		Descriptor: typesObjectWithSpec(),
		Ownership:  state(entry("other", specPath())),
	}

	_, err := Apply(req, Options{Force: true})

	var conflictError *fieldownership.ConflictError
	if !errors.As(err, &conflictError) {
		t.Fatalf("error does not wrap fieldownership.ConflictError: %v", err)
	}
}

func TestUnsupportedTakeoverConflictsKeepsOnlyOwnedAncestors(t *testing.T) {
	conflicts := fieldownership.NewConflictSet(
		fieldownership.Conflict{
			Owner:         owner("ancestor"),
			OwnedPath:     path("$.spec"),
			AttemptedPath: path("$.spec.image"),
		},
		fieldownership.Conflict{
			Owner:         owner("exact"),
			OwnedPath:     path("$.spec.image"),
			AttemptedPath: path("$.spec.image"),
		},
		fieldownership.Conflict{
			Owner:         owner("descendant"),
			OwnedPath:     path("$.spec.image"),
			AttemptedPath: path("$.spec"),
		},
	)

	got := unsupportedTakeoverConflicts(conflicts)

	if got.Len() != 1 {
		t.Fatalf("unsupported conflicts = %d; want 1", got.Len())
	}
	requireSet(t, got.OwnedPaths(), "$.spec")
	requireSet(t, got.AttemptedPaths(), "$.spec.image")
}

func TestApplyForceListDescendantConflictTakesOwnership(t *testing.T) {
	req := Request{
		Path:    root(),
		Owner:   owner("user"),
		Live:    obj(member("args", list(str("old")))),
		Applied: obj(member("args", list(str("new")))),
		Descriptor: types.Object(
			types.Field("args").ListOf(types.String()).Atomic().Optional(),
		).Descriptor(),
		Ownership: state(entry("other", path("$.args[0]"))),
	}

	result, err := Apply(req, Options{Force: true})
	requireNoError(t, err)

	args := requireMember(t, result.Value, "args")
	requireListStrings(t, args, "new")
	requireOwnersOf(t, result.Ownership, path("$.args"), "user")
	requireOwnersOf(t, result.Ownership, path("$.args[0]"))
}
