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

	"arcoris.dev/apimachinery/api/types"
)

func TestMergeOrderedListSelectedIndex(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a"), str("b"), str("c")),
		list(str("x"), str("B"), str("z")),
		descriptor,
		pathSet(root().Index(1)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "B", "c")
}

func TestMergeOrderedListAddsSelectedIndex(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a")),
		list(str("a"), str("b")),
		descriptor,
		pathSet(root().Index(1)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "b")
}

func TestMergeOrderedListSparseAppendUnsupported(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	_, err := Merge(
		list(str("a")),
		list(str("a"), str("b"), str("c")),
		descriptor,
		pathSet(root().Index(2)),
		Options{},
	)

	requireErrorIs(t, err, ErrUnsupportedMerge)
}

func TestMergeOrderedListContiguousAppendSelectedIndexes(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a")),
		list(str("a"), str("b"), str("c")),
		descriptor,
		pathSet(root().Index(1), root().Index(2)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "b", "c")
}

func TestMergeOrderedListAppendFirstMissingIndexAllowed(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a")),
		list(str("a"), str("b"), str("c")),
		descriptor,
		pathSet(root().Index(1)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "b")
}

func TestMergeOrderedListExistingIndexSelectionStillWorks(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a"), str("b")),
		list(str("x"), str("B"), str("c")),
		descriptor,
		pathSet(root().Index(1)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "B")
}

func TestMergeOrderedListRemovesSelectedIndex(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a"), str("b"), str("c")),
		list(str("a")),
		descriptor,
		pathSet(root().Index(1), root().Index(2)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a")
}

func TestMergeOrderedListRemovesMultipleIndexes(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := Merge(
		list(str("a"), str("b"), str("c"), str("d")),
		list(str("a")),
		descriptor,
		pathSet(root().Index(1), root().Index(3)),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireListStrings(t, got, "a", "c")
}

func TestMergeOrderedListNestedObjectItemField(t *testing.T) {
	descriptor := types.ListOf(
		types.Object(
			types.Field("image").String().Optional(),
			types.Field("name").String().Optional(),
		),
	).Ordered().Descriptor()
	base := list(obj(member("name", str("api")), member("image", str("api:v1"))))
	overlay := list(obj(member("name", str("ignored")), member("image", str("api:v2"))))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Index(0).Field("image")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	view, _ := got.List()
	item, _ := view.At(0)
	requireStringMember(t, item, "name", "api")
	requireStringMember(t, item, "image", "api:v2")
}

func TestMergeOrderedListExactSelectedReplacesWholeList(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()
	base := list(str("old"))
	overlay := list(str("new"), str("next"))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}
