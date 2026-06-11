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

func TestCompareRecordSameIsEmpty(t *testing.T) {
	descriptor := typesObject("name")
	oldValue := valueRecord("name", "app")

	got, err := Compare(oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordModifiedField(t *testing.T) {
	got, err := Compare(valueRecord("name", "a"), valueRecord("name", "b"), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("name")))
}

func TestCompareRecordAddedField(t *testing.T) {
	got, err := Compare(valueRecord(), valueRecord("name", "app"), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("name")), nil, nil)
}

func TestCompareRecordRemovedField(t *testing.T) {
	got, err := Compare(valueRecord("name", "app"), valueRecord(), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("name")), nil)
}

func TestCompareRecordNestedModifiedField(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
		).Optional(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v1")))
	newValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v2")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec", "image")))
}

func TestCompareRecordMissingBothFieldIsEmpty(t *testing.T) {
	got, err := Compare(valueRecord(), valueRecord(), typesObject("name"), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordAddedNullField(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional().Nullable()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("name", value.NullValue()))

	got, err := Compare(valueRecord(), newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("name")), nil, nil)
}

func TestCompareRecordRemovedNullField(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional().Nullable()).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("name", value.NullValue()))

	got, err := Compare(oldValue, valueRecord(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("name")), nil)
}

func TestCompareRecordNullToSubtreeKeepsBucketsDisjoint(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional().Nullable(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("spec", value.NullValue()))
	newValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v1")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec")))
}

func TestCompareRecordSubtreeToNullKeepsBucketsDisjoint(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional().Nullable(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v1")))
	newValue := value.MustRecordValue(value.MustRecordMember("spec", value.NullValue()))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("spec")))
}

func TestCompareRecordEmptyToNonEmpty(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional(),
	).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v1")))

	got, err := Compare(valueRecord(), newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("spec", "image")), nil, nil)
}

func TestCompareRecordNonEmptyToEmpty(t *testing.T) {
	descriptor := types.Object(
		types.Field("spec").Object(types.Field("image").String().Optional()).Optional(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("spec", valueRecord("image", "v1")))

	got, err := Compare(oldValue, valueRecord(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("spec", "image")), nil)
}

func TestCompareRecordUnknownRejectedReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueRecord("extra", "old"), valueRecord(), types.Object().Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.extra")
}

func TestCompareRecordUnknownPreserveOpaqueSameOpaqueIsEmpty(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Compare(valueRecord("extra", "same"), valueRecord("extra", "same"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordUnknownPreserveOpaqueChangedOpaqueIsModified(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Compare(valueRecord("extra", "old"), valueRecord("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("extra")))
}

func TestCompareRecordUnknownPreserveOpaqueAddedOpaqueIsAdded(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Compare(valueRecord(), valueRecord("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("extra")), nil, nil)
}

func TestCompareRecordUnknownPreserveOpaqueRemovedOpaqueIsRemoved(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	got, err := Compare(valueRecord("extra", "old"), valueRecord(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(rootField("extra")), nil)
}

func TestCompareRecordUnknownPrunedIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Compare(valueRecord("extra", "old"), valueRecord("extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordOnlyPrunedUnknownFieldsIsEmpty(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Compare(valueRecord("extra", "old"), valueRecord(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}
func TestObjectFieldsByNameBuildsLookup(t *testing.T) {
	descriptor := types.Object(
		types.Field("name").String().Optional(),
		types.Field("image").String().Optional(),
	).Descriptor()
	objectView, _ := descriptor.AsObject()

	got := recordFieldsByName(objectView.Fields())

	if string(got["name"].Name()) != "name" || string(got["image"].Name()) != "image" {
		t.Fatalf("recordFieldsByName() = %#v", got)
	}
}
func TestCompareUnknownRecordMembersPruneIgnoresUnknowns(t *testing.T) {
	oldRecord, _ := valueRecord("extra", "old").AsRecord()
	newRecord, _ := valueRecord("extra", "new").AsRecord()

	got, err := newComparer(Options{}).compareUnknownRecordMembers(rootField("spec"), oldRecord, newRecord, nil, types.UnknownPrune)
	requireNoError(t, err)

	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordUnknownPrunedOldIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Compare(valueRecord("x-extra", "old"), valueRecord(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareRecordUnknownPrunedNewIgnored(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPrune).Descriptor()

	got, err := Compare(valueRecord(), valueRecord("x-extra", "new"), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareUnknownRecordMembersInvalidPolicy(t *testing.T) {
	oldRecord, _ := valueRecord().AsRecord()
	newRecord, _ := valueRecord().AsRecord()

	_, err := newComparer(Options{}).compareUnknownRecordMembers(rootField("spec"), oldRecord, newRecord, nil, types.UnknownFieldPolicy(99))

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
func TestUnknownMemberNamesReturnsSortedUndeclaredNames(t *testing.T) {
	descriptor := types.Object(types.Field("known").String().Optional()).Descriptor()
	objectView, _ := descriptor.AsObject()
	declared := recordFieldsByName(objectView.Fields())
	oldRecord, _ := valueRecord("known", "x", "zeta", "old").AsRecord()
	newRecord, _ := valueRecord("alpha", "new").AsRecord()

	got := unknownMemberNames(oldRecord, newRecord, declared)

	if want := []string{"alpha", "zeta"}; !slices.Equal(got, want) {
		t.Fatalf("unknownMemberNames() = %#v, want %#v", got, want)
	}
}
func TestComparePreservedUnknownRecordMemberAdded(t *testing.T) {
	oldRecord, _ := valueRecord().AsRecord()
	newRecord, _ := valueRecord("extra", "new").AsRecord()

	got, err := newComparer(Options{}).comparePreservedUnknownRecordMember(rootField("extra"), oldRecord, newRecord, "extra")
	requireNoError(t, err)

	requireResult(t, got, paths(rootField("extra")), nil, nil)
}

func TestCompareOpaqueLeafModified(t *testing.T) {
	got, err := newComparer(Options{}).compareOpaqueLeaf(rootField("extra"), valueRecord("nested", "old"), valueRecord("nested", "new"))
	requireNoError(t, err)

	requireResult(t, got, nil, nil, paths(rootField("extra")))
}

func TestCompareRecordUnknownPreserveOpaqueDoesNotDescendIntoNestedObject(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("x-extra", valueRecord("nested", "old")))
	newValue := value.MustRecordValue(value.MustRecordMember("x-extra", valueRecord("nested", "new")))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("x-extra")))
}

func TestCompareRecordUnknownPreserveOpaqueDoesNotDescendIntoNestedList(t *testing.T) {
	descriptor := types.Object().UnknownFields(types.UnknownPreserveOpaque).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("x-extra", value.MustListValue(value.StringValue("old"))))
	newValue := value.MustRecordValue(value.MustRecordMember("x-extra", value.MustListValue(value.StringValue("new"))))

	got, err := Compare(oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(rootField("x-extra")))
}
func TestRejectUnknownRecordMembersReturnsUnknownField(t *testing.T) {
	oldRecord, _ := valueRecord("extra", "old").AsRecord()
	newRecord, _ := valueRecord().AsRecord()

	_, err := newComparer(Options{}).rejectUnknownRecordMembers(rootField("spec"), oldRecord, newRecord, nil)

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorPath(t, err, "$.spec.extra")
}

func TestCompareRecordUnknownRejectedOldReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueRecord("x-extra", "old"), valueRecord(), types.Object().Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, rootField("x-extra").String())
}

func TestCompareRecordUnknownRejectedNewReturnsUnknownField(t *testing.T) {
	_, err := Compare(valueRecord(), valueRecord("x-extra", "new"), types.Object().Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, rootField("x-extra").String())
}
