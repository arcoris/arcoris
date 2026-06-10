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

package objectapply

import "testing"

func TestApplyDesiredField(t *testing.T) {
	result, err := Apply(testRequest(), Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "image", "api:v2")
	requireStringMember(t, result.Object.Desired, "replicas", "3")
	requireSet(t, result.Desired.AppliedFields, "$.image")
}

func TestApplyDesiredMapKey(t *testing.T) {
	result, err := Apply(Request{
		Owner:    owner("user"),
		Live:     testObject(obj(member("app", str("old")))),
		Applied:  appliedObject(obj(member("app", str("new")))),
		Resource: testResource(mapDesiredDescriptor()),
	}, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "app", "new")
	requireSet(t, result.Desired.AppliedFields, `$["app"]`)
}

func TestApplyDesiredListMapConditionStatus(t *testing.T) {
	result, err := Apply(Request{
		Owner:    owner("user"),
		Live:     testObject(list(readyCondition("False"))),
		Applied:  appliedObject(list(readyCondition("True"))),
		Resource: testResource(conditionsDescriptor()),
	}, Options{})
	requireNoError(t, err)

	view, ok := result.Object.Desired.AsList()
	if !ok {
		t.Fatalf("desired kind = %s; want list", result.Object.Desired.Kind())
	}
	item, ok := view.At(0)
	if !ok {
		t.Fatalf("list item 0 is absent")
	}
	requireStringMember(t, item, "status", "True")
	requireSet(t, result.Desired.AppliedFields, `$[{"type":"Ready"}].status`, `$[{"type":"Ready"}].type`)
}

func TestApplyDesiredDroppedFieldDeletion(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("user", path("$.image"), path("$.replicas")))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireNoMember(t, result.Object.Desired, "replicas")
	requireSet(t, result.Desired.DeletedFields, "$.replicas")
}

func TestApplyDesiredDroppedFieldPreservedByOtherOwner(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(
		entry("user", path("$.image"), path("$.replicas")),
		entry("other", path("$.replicas")),
	)

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Object.Desired, "replicas", "3")
	requireSet(t, result.Desired.DeletedFields)
}

func TestApplyDesiredSameValueSharedOwnership(t *testing.T) {
	req := testRequest()
	req.Applied = appliedObject(obj(member("image", str("api:v1"))))
	req.Ownership = desiredOwnership(entry("other", path("$.image")))

	result, err := Apply(req, Options{})
	requireNoError(t, err)

	requireSet(t, result.Desired.ChangedAppliedFields)
	requireOwnersOf(t, result.Ownership.Desired(), path("$.image"), "other", "user")
}
