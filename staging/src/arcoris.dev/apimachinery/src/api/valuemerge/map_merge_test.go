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

func TestMergeMapSelectedKey(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	base := obj(member("app", str("old")), member("tier", str("backend")))
	overlay := obj(member("app", str("new")), member("tier", str("frontend")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Key(testMapKey("app"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "app", "new")
	requireStringMember(t, got, "tier", "backend")
}

func TestMergeMapAddsSelectedKey(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	base := obj(member("tier", str("backend")))
	overlay := obj(member("app", str("api")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Key(testMapKey("app"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "tier", "backend")
	requireStringMember(t, got, "app", "api")
}

func TestMergeMapRemovesSelectedKeyAbsentFromOverlay(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	base := obj(member("app", str("api")), member("tier", str("backend")))
	overlay := obj(member("tier", str("frontend")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Key(testMapKey("app"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireNoMember(t, got, "app")
	requireStringMember(t, got, "tier", "backend")
}

func TestMergeMapExactSelectedReplacesWholeMap(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	base := obj(member("app", str("old")))
	overlay := obj(member("tier", str("backend")))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}

func TestMergeMapNestedObjectValueSelectedField(t *testing.T) {
	descriptor := types.MapOf(
		types.Object(
			types.Field("image").String().Optional(),
			types.Field("replicas").Int64().Optional(),
		),
	).Descriptor()
	base := obj(member("api", obj(member("image", str("api:v1")), member("replicas", intValue(3)))))
	overlay := obj(member("api", obj(member("image", str("api:v2")))))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Key(testMapKey("api")).Field(testFieldName("image"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	view, _ := got.AsRecord()
	api, _ := view.Get(value.MemberName("api"))
	requireStringMember(t, api, "image", "api:v2")
	requireIntegerMember(t, api, "replicas", 3)
}

func TestMergeMapUsesKeyPathNotFieldPath(t *testing.T) {
	descriptor := types.MapOf(types.String()).Descriptor()
	base := obj(member("app", str("old")))
	overlay := obj(member("app", str("new")))

	got, err := Merge(
		base,
		overlay,
		descriptor,
		pathSet(root().Field(testFieldName("app"))),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireStringMember(t, got, "app", "old")
}
