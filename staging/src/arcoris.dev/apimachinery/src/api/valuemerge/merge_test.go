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
)

func TestMergeEmptyFieldsPreservesBase(t *testing.T) {
	base := str("old")
	overlay := str("new")

	got, err := Merge(base, overlay, types.String().Descriptor(), fieldpath.EmptySet(), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, base)
}

func TestMergeEmptyFieldsDoesNotRequireValidOverlay(t *testing.T) {
	base := str("old")

	got, err := Merge(
		base,
		value.Value{},
		types.String().Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, base)
}

func TestMergeEmptyFieldsStillRejectsInvalidBase(t *testing.T) {
	_, err := Merge(
		value.Value{},
		value.Value{},
		types.String().Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestMergeEmptyFieldsClonesBase(t *testing.T) {
	base := obj(member("name", str("old")))

	got, err := Merge(
		base,
		value.Value{},
		types.Object(types.Field("name").String().Optional()).Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, base)
	if got.IsZero() {
		t.Fatalf("got is zero")
	}
}

func TestMergeAtEmptyFieldsRejectsInvalidBasePath(t *testing.T) {
	invalidPath := root().Append(fieldpath.Element{})

	_, err := MergeAt(
		invalidPath,
		str("old"),
		value.Value{},
		types.String().Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestMergeAtEmptyFieldsPreservesNestedBase(t *testing.T) {
	path := root().Field(testFieldName("spec"))
	base := obj(member("name", str("old")))

	got, err := MergeAt(
		path,
		base,
		value.Value{},
		types.Object(types.Field("name").String().Optional()).Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)
	if err != nil {
		t.Fatalf("MergeAt returned error: %v", err)
	}

	requireValue(t, got, base)
}

func TestMergeAtRejectsSelectedRootOutsideNestedBase(t *testing.T) {
	path := root().Field(testFieldName("spec"))

	_, err := MergeAt(
		path,
		str("old"),
		str("new"),
		types.String().Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestMergeAtRejectsUnrelatedSelectedPath(t *testing.T) {
	path := root().Field(testFieldName("spec"))

	_, err := MergeAt(
		path,
		str("old"),
		str("new"),
		types.String().Descriptor(),
		pathSet(root().Field(testFieldName("metadata"))),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestMergeAtAcceptsExactBasePath(t *testing.T) {
	path := root().Field(testFieldName("spec"))

	got, err := MergeAt(
		path,
		str("old"),
		str("new"),
		types.String().Descriptor(),
		pathSet(path),
		Options{},
	)
	if err != nil {
		t.Fatalf("MergeAt returned error: %v", err)
	}

	requireValue(t, got, str("new"))
}

func TestMergeAtAcceptsDescendantOfBasePath(t *testing.T) {
	path := root().Field(testFieldName("spec"))
	descriptor := types.Object(types.Field("name").String().Optional()).Descriptor()

	got, err := MergeAt(
		path,
		obj(member("name", str("old"))),
		obj(member("name", str("new"))),
		descriptor,
		pathSet(path.Field(testFieldName("name"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("MergeAt returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
}

func TestMergeDoesNotMutateBase(t *testing.T) {
	base := obj(member("name", str("old")))
	overlay := obj(member("name", str("new")))
	descriptor := types.Object(types.Field("name").String().Optional()).Descriptor()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
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
	descriptor := types.Object(types.Field("name").String().Optional()).Descriptor()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
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
		types.String().Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestMergeInvalidOverlayValueReturnsInvalidValue(t *testing.T) {
	_, err := Merge(
		str("old"),
		value.Value{},
		types.String().Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidValue)
}

func TestMergeKindMismatchReturnsKindMismatch(t *testing.T) {
	_, err := Merge(
		str("old"),
		boolValue(true),
		types.String().Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrKindMismatch)
}
