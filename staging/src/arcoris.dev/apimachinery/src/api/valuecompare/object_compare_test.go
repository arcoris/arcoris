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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"slices"
	"testing"
)

func TestCompareObjectSameIsEmpty(t *testing.T) {
	descriptor := typesObject("name")
	oldValue := valueObject("name", "app")

	got, err := Compare(oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectModifiedField(t *testing.T) {
	got, err := Compare(valueObject("name", "a"), valueObject("name", "b"), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("name")))
}

func TestCompareObjectAddedField(t *testing.T) {
	got, err := Compare(valueObject(), valueObject("name", "app"), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("name")), nil, nil)
}

func TestCompareObjectRemovedField(t *testing.T) {
	got, err := Compare(valueObject("name", "app"), valueObject(), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("name")), nil)
}

func TestCompareObjectNestedModifiedField(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
		).Optional(),
	).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v1")))
	newValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v2")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec", "image")))
}

func TestCompareObjectMissingBothFieldIsEmpty(t *testing.T) {
	got, err := Compare(valueObject(), valueObject(), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectAddedNullField(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional().Nullable()).Type()
	newValue := value.MustObjectValue(value.ObjectMember("name", value.NullValue()))

	got, err := Compare(valueObject(), newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("name")), nil, nil)
}

func TestCompareObjectRemovedNullField(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional().Nullable()).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("name", value.NullValue()))

	got, err := Compare(oldValue, valueObject(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("name")), nil)
}

func TestCompareObjectNullToSubtreeKeepsBucketsDisjoint(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional().Nullable(),
	).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("spec", value.NullValue()))
	newValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v1")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec")))
}

func TestCompareObjectSubtreeToNullKeepsBucketsDisjoint(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional().Nullable(),
	).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v1")))
	newValue := value.MustObjectValue(value.ObjectMember("spec", value.NullValue()))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec")))
}

func TestCompareObjectEmptyToNonEmpty(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional(),
	).Type()
	newValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v1")))

	got, err := Compare(valueObject(), newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("spec", "image")), nil, nil)
}

func TestCompareObjectNonEmptyToEmpty(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional(),
	).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("spec", valueObject("image", "v1")))

	got, err := Compare(oldValue, valueObject(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("spec", "image")), nil)
}

func TestCompareObjectUnknownRejectedReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueObject("extra", "old"), valueObject(), types.Object().Type(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.extra")
}

func TestCompareObjectUnknownPreservedSameOpaqueIsEmpty(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()

	got, err := Compare(valueObject("extra", "same"), valueObject("extra", "same"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectUnknownPreservedChangedOpaqueIsModified(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()

	got, err := Compare(valueObject("extra", "old"), valueObject("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("extra")))
}

func TestCompareObjectUnknownPreservedAddedOpaqueIsAdded(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()

	got, err := Compare(valueObject(), valueObject("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("extra")), nil, nil)
}

func TestCompareObjectUnknownPreservedRemovedOpaqueIsRemoved(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()

	got, err := Compare(valueObject("extra", "old"), valueObject(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("extra")), nil)
}

func TestCompareObjectUnknownPrunedIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Type()

	got, err := Compare(valueObject("extra", "old"), valueObject("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectOnlyPrunedUnknownFieldsIsEmpty(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Type()

	got, err := Compare(valueObject("extra", "old"), valueObject(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}
func TestObjectFieldsByNameBuildsLookup(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
		types.Field("image").String().Optional(),
	).Type()
	objectView, _ := descriptor.Object()

	got := objectFieldsByName(objectView.Fields())

	if string(got["name"].Name()) != "name" || string(got["image"].Name()) != "image" {
		t.Fatalf("objectFieldsByName() = %#v", got)
	}
}
func TestCompareUnknownObjectMembersPruneIgnoresUnknowns(t *testing.T) {
	oldObject, _ := valueObject("extra", "old").Object()
	newObject, _ := valueObject("extra", "new").Object()

	got, err := newComparer(Options{}).compareUnknownObjectMembers(rootField("spec"), oldObject, newObject, nil, types.UnknownPrune)
	requireNoError(t, err)

	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectUnknownPrunedOldIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Type()

	got, err := Compare(valueObject("x-extra", "old"), valueObject(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareObjectUnknownPrunedNewIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Type()

	got, err := Compare(valueObject(), valueObject("x-extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareUnknownObjectMembersInvalidPolicy(t *testing.T) {
	oldObject, _ := valueObject().Object()
	newObject, _ := valueObject().Object()

	_, err := newComparer(Options{}).compareUnknownObjectMembers(rootField("spec"), oldObject, newObject, nil, types.UnknownFieldPolicy(99))

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
func TestUnknownMemberNamesReturnsSortedUndeclaredNames(t *testing.T) {
	descriptor := types.Object(types.Field("known").String().Optional()).Type()
	objectView, _ := descriptor.Object()
	declared := objectFieldsByName(objectView.Fields())
	oldObject, _ := valueObject("known", "x", "zeta", "old").Object()
	newObject, _ := valueObject("alpha", "new").Object()

	got := unknownMemberNames(oldObject, newObject, declared)

	if want := []string{"alpha", "zeta"}; !slices.Equal(got, want) {
		t.Fatalf("unknownMemberNames() = %#v, want %#v", got, want)
	}
}
func TestComparePreservedUnknownObjectMemberAdded(t *testing.T) {
	oldObject, _ := valueObject().Object()
	newObject, _ := valueObject("extra", "new").Object()

	got, err := newComparer(Options{}).comparePreservedUnknownObjectMember(rootField("extra"), oldObject, newObject, "extra")
	requireNoError(t, err)

	requireResult(t, got, paths(rootField("extra")), nil, nil)
}

func TestCompareOpaqueLeafModified(t *testing.T) {
	got, err := newComparer(Options{}).compareOpaqueLeaf(rootField("extra"), valueObject("nested", "old"), valueObject("nested", "new"))
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(rootField("extra")))
}

func TestCompareObjectUnknownPreservedDoesNotDescendIntoNestedObject(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("x-extra", valueObject("nested", "old")))
	newValue := value.MustObjectValue(value.ObjectMember("x-extra", valueObject("nested", "new")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("x-extra")))
}

func TestCompareObjectUnknownPreservedDoesNotDescendIntoNestedList(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserve).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("x-extra", value.MustListValue(value.StringValue("old"))))
	newValue := value.MustObjectValue(value.ObjectMember("x-extra", value.MustListValue(value.StringValue("new"))))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("x-extra")))
}
func TestRejectUnknownObjectMembersReturnsUnknownField(t *testing.T) {
	oldObject, _ := valueObject("extra", "old").Object()
	newObject, _ := valueObject().Object()

	_, err := newComparer(Options{}).rejectUnknownObjectMembers(rootField("spec"), oldObject, newObject, nil)

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorPath(t, err, "$.spec.extra")
}

func TestCompareObjectUnknownRejectedOldReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueObject("x-extra", "old"), valueObject(), types.Object().Type(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, rootField("x-extra").String())
}

func TestCompareObjectUnknownRejectedNewReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueObject(), valueObject("x-extra", "new"), types.Object().Type(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, rootField("x-extra").String())
}
