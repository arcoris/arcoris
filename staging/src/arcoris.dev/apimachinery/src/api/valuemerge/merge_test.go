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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuecompare"
)

func TestMergeEmptyFieldsPreservesBase(t *testing.T) {
	base := str("old")
	overlay := str("new")

	got, err := Merge(base, overlay, types.String().Type(), fieldpath.EmptySet(), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, base)
}

func TestMergeDoesNotMutateBase(t *testing.T) {
	base := obj(member("name", str("old")))
	overlay := obj(member("name", str("new")))
	descriptor := types.Object(types.Field("name").String().Optional()).Type()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("name")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
	requireStringMember(t, base, "name", "old")
}

func TestMergeDoesNotMutateOverlay(t *testing.T) {
	base := obj(member("name", str("old")))
	overlay := obj(member("name", str("new")))
	descriptor := types.Object(types.Field("name").String().Optional()).Type()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("name")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
	requireStringMember(t, overlay, "name", "new")
}

func TestMergeInvalidBaseValueReturnsInvalidValue(t *testing.T) {
	_, err := Merge(
		value.Value{},
		str("new"),
		types.String().Type(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestMergeInvalidOverlayValueReturnsInvalidValue(t *testing.T) {
	_, err := Merge(
		str("old"),
		value.Value{},
		types.String().Type(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestMergeKindMismatchReturnsKindMismatch(t *testing.T) {
	_, err := Merge(
		str("old"),
		boolValue(true),
		types.String().Type(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrKindMismatch)
}

func TestMergeValueCompareModifiedSet(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Type()
	base := obj(member("spec", obj(member("image", str("api:v1")), member("replicas", intValue(3)))))
	overlay := obj(member("spec", obj(member("image", str("api:v1")), member("replicas", intValue(5)))))

	changes, err := valuecompare.Compare(base, overlay, descriptor, valuecompare.Options{})
	if err != nil {
		t.Fatalf("Compare returned error: %v", err)
	}

	got, err := Merge(base, overlay, descriptor, changes.Modified, Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	view, _ := got.Object()
	spec, _ := view.Get("spec")
	requireStringMember(t, spec, "image", "api:v1")
	requireIntegerMember(t, spec, "replicas", 5)
}
