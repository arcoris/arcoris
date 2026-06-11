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
	"arcoris.dev/apimachinery/api/value"
)

func TestMergeUnknownRejectSelectedUnknownReturnsUnknownField(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownReject).Descriptor()

	_, err := Merge(
		obj(member("xExtra", str("old"))),
		obj(member("xExtra", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)

	requireErrorIs(t, err, ErrUnknownField)
}

func TestMergeUnknownRejectUnselectedBaseUnknownReturnsUnknownField(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownReject).Descriptor()

	_, err := Merge(
		obj(member("name", str("old")), member("xExtra", str("old"))),
		obj(member("name", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)

	requireErrorIs(t, err, ErrUnknownField)
}

func TestMergeUnknownRejectUnselectedOverlayUnknownReturnsUnknownField(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownReject).Descriptor()

	_, err := Merge(
		obj(member("name", str("old"))),
		obj(member("name", str("new")), member("xExtra", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)

	requireErrorIs(t, err, ErrUnknownField)
}

func TestMergeUnknownPruneIgnoresUnknown(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownPrune).Descriptor()
	base := obj(member("name", str("old")), member("xExtra", str("old")))
	overlay := obj(member("xExtra", str("new")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "old")
	requireNoMember(t, got, "xExtra")
}

func TestMergeUnknownPruneUnselectedBaseUnknownPruned(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Merge(
		obj(member("name", str("old")), member("xExtra", str("old"))),
		obj(member("name", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
	requireNoMember(t, got, "xExtra")
}

func TestMergeUnknownPruneSelectedUnknownIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Merge(
		obj(member("xExtra", str("old"))),
		obj(member("xExtra", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "xExtra")
}

func TestMergeUnknownPruneUnselectedOverlayUnknownPruned(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Merge(
		obj(member("name", str("old"))),
		obj(member("name", str("new")), member("xExtra", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
	requireNoMember(t, got, "xExtra")
}

func TestMergeUnknownPreserveOpaqueUnselectedBaseUnknownPreserveOpaque(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Merge(
		obj(member("name", str("old")), member("xExtra", str("old"))),
		obj(member("name", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "name", "new")
	requireStringMember(t, got, "xExtra", "old")
}

func TestMergeUnknownPreserveOpaqueExactCopiesOpaque(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()
	overlayExtra := obj(member("nested", str("new")))

	got, err := Merge(
		obj(member("xExtra", obj(member("nested", str("old"))))),
		obj(member("xExtra", overlayExtra)),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	view, _ := got.AsRecord()
	extra, ok := view.Get(value.MemberName("xExtra"))
	if !ok {
		t.Fatalf("xExtra is absent")
	}
	requireValue(t, extra, overlayExtra)
}

func TestMergeUnknownPreserveOpaqueRemovesOpaqueWhenOverlayAbsent(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Merge(
		obj(member("xExtra", str("old"))),
		obj(),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "xExtra")
}

func TestMergeUnknownPreserveOpaqueDescendantSelectionUnsupported(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	_, err := Merge(
		obj(member("xExtra", obj(member("nested", str("old"))))),
		obj(member("xExtra", obj(member("nested", str("new"))))),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra")).Field(testFieldName("nested"))),
		Options{},
	)

	requireErrorIs(t, err, ErrUnsupportedMerge)
}

func TestMergeInvalidUnknownPolicyUnselectedUnknownReturnsInvalidDescriptor(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
	).UnknownFields(types.UnknownFieldPolicy(255)).Descriptor()

	_, err := Merge(
		obj(member("name", str("old")), member("xExtra", str("old"))),
		obj(member("name", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("name"))),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidDescriptor)
}

func TestMergeInvalidUnknownPolicySelectedUnknownReturnsInvalidDescriptor(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownFieldPolicy(255)).Descriptor()

	_, err := Merge(
		obj(member("xExtra", str("old"))),
		obj(member("xExtra", str("new"))),
		descriptor,
		pathSet(root().Field(testFieldName("xExtra"))),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
