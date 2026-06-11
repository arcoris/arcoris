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

func TestMergeRecordSelectedField(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("image"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v2")
	requireIntegerMember(t, got, "replicas", 3)
}

func TestMergeRecordUnselectedFieldPreserved(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireIntegerMember(t, got, "replicas", 5)
}

func TestMergeRecordAddsSelectedField(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireIntegerMember(t, got, "replicas", 5)
}

func TestMergeRecordUnselectedOverlayOnlyFieldNotInspected(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")))
	overlay := obj(member("image", str("api:v2")), member("replicas", str("wrong-kind")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("image"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v2")
	requireNoMember(t, got, "replicas")
}

func TestMergeRecordRemovesSelectedFieldAbsentFromOverlay(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireNoMember(t, got, "replicas")
}

func TestMergeRecordExactSelectedReplacesWholeObject(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}

func TestMergeRecordNestedSelectedField(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Descriptor()
	base := obj(member("spec", obj(member("image", str("api:v1")), member("replicas", intValue(3)))))
	overlay := obj(member("spec", obj(member("image", str("api:v2")))))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("spec")).Field(testFieldName("image"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	specView, _ := got.AsRecord()
	spec, _ := specView.Get("spec")
	requireStringMember(t, spec, "image", "api:v2")
	requireIntegerMember(t, spec, "replicas", 3)
}

func TestMergeDescendantIntoWrongKindBaseReturnsKindMismatch(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Descriptor()

	_, err := Merge(
		obj(member("spec", str("invalid"))),
		obj(),
		descriptor,
		pathSet(root().Field(testFieldName("spec")).Field(testFieldName("replicas"))),
		Options{},
	)

	requireErrorIs(t, err, ErrKindMismatch)
}

func TestMergeDescendantIntoAbsentBaseAndAbsentOverlayNoops(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Descriptor()

	got, err := Merge(
		obj(),
		obj(),
		descriptor,
		pathSet(root().Field(testFieldName("spec")).Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "spec")
}

func TestMergeDescendantIntoNullBaseAndAbsentOverlayPreservesNull(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Descriptor()

	got, err := Merge(
		obj(member("spec", value.NullValue())),
		obj(),
		descriptor,
		pathSet(root().Field(testFieldName("spec")).Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	spec := requireMember(t, got, "spec")
	if !spec.IsNull() {
		t.Fatalf("spec is not null")
	}
}

func TestMergeDescendantIntoAbsentBaseCreatesFromOverlayContainer(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Descriptor()

	got, err := Merge(
		obj(),
		obj(member("spec", obj(member("image", str("ignored")), member("replicas", intValue(5))))),
		descriptor,
		pathSet(root().Field(testFieldName("spec")).Field(testFieldName("replicas"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	spec := requireMember(t, got, "spec")
	requireIntegerMember(t, spec, "replicas", 5)
	requireNoMember(t, spec, "image")
}

func TestMergeRecordRequiredFieldAbsentIsNotValidationError(t *testing.T) {
	descriptor := types.Object(
		types.Field("image").String().Required(),
		types.Field("replicas").Int64().Optional(),
	).Descriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("image"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "image")
	requireIntegerMember(t, got, "replicas", 3)
}

func simpleSpecDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").Int64().Optional(),
	).Descriptor()
}
