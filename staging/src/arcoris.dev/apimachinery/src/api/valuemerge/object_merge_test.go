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

func TestMergeObjectSelectedField(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("image")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v2")
	requireIntegerMember(t, got, "replicas", 3)
}

func TestMergeObjectUnselectedFieldPreserved(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("replicas")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireIntegerMember(t, got, "replicas", 5)
}

func TestMergeObjectAddsSelectedField(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")))
	overlay := obj(member("image", str("api:v2")), member("replicas", intValue(5)))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("replicas")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireIntegerMember(t, got, "replicas", 5)
}

func TestMergeObjectRemovesSelectedFieldAbsentFromOverlay(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("replicas")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "image", "api:v1")
	requireNoMember(t, got, "replicas")
}

func TestMergeObjectExactSelectedReplacesWholeObject(t *testing.T) {
	descriptor := simpleSpecDescriptor()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj(member("image", str("api:v2")))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}

func TestMergeObjectNestedSelectedField(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
			types.Field("replicas").Int64().Optional(),
		).Optional(),
	).Type()
	base := obj(member("spec", obj(member("image", str("api:v1")), member("replicas", intValue(3)))))
	overlay := obj(member("spec", obj(member("image", str("api:v2")))))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("spec").Field("image")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	specView, _ := got.Object()
	spec, _ := specView.Get("spec")
	requireStringMember(t, spec, "image", "api:v2")
	requireIntegerMember(t, spec, "replicas", 3)
}

func TestMergeObjectRequiredFieldAbsentIsNotValidationError(t *testing.T) {
	descriptor := types.Object(
		types.Field("image").String().Required(),
		types.Field("replicas").Int64().Optional(),
	).Type()
	base := obj(member("image", str("api:v1")), member("replicas", intValue(3)))
	overlay := obj()

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field("image")),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "image")
	requireIntegerMember(t, got, "replicas", 3)
}

func simpleSpecDescriptor() types.Type {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").Int64().Optional(),
	).Type()
}
