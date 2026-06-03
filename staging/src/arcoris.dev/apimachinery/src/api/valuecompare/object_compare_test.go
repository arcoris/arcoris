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
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
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
