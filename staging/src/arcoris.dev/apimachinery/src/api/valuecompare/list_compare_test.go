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
	"testing"
)

func TestCompareAtomicListSameIsEmpty(t *testing.T) {
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()
	oldValue := value.MustListValue(value.StringValue("one"))

	got, err := CompareAt(rootField("args"), oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareAtomicListChangedItemIsModifiedAtListPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("one")),
		value.MustListValue(value.StringValue("two")),
		descriptor,
		Options{},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareAtomicListAddedItemIsModifiedAtListPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()

	got, err := CompareAt(path, value.MustListValue(), value.MustListValue(value.StringValue("one")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareAtomicListRemovedItemIsModifiedAtListPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()

	got, err := CompareAt(path, value.MustListValue(value.StringValue("one")), value.MustListValue(), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareAtomicListNeverEmitsIndexPaths(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Atomic().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("one")),
		value.MustListValue(value.StringValue("two")),
		descriptor,
		Options{},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareSetListSameIsEmpty(t *testing.T) {
	descriptor := types.ListOf(types.String()).Set().Descriptor()
	oldValue := value.MustListValue(value.StringValue("one"))

	got, err := CompareAt(rootField("tags"), oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareSetListDifferentOrderIsModifiedAtListPathForNow(t *testing.T) {
	path := rootField("tags")
	descriptor := types.ListOf(types.String()).Set().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("a"), value.StringValue("b")),
		value.MustListValue(value.StringValue("b"), value.StringValue("a")),
		descriptor,
		Options{},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareSetListNeverEmitsIndexPaths(t *testing.T) {
	path := rootField("tags")
	descriptor := types.ListOf(types.String()).Set().Descriptor()

	got, err := CompareAt(path, value.MustListValue(), value.MustListValue(value.StringValue("a")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareSetListAddedItemIsModifiedAtListPath(t *testing.T) {
	path := rootField("tags")
	descriptor := types.ListOf(types.String()).Set().Descriptor()

	got, err := CompareAt(path, value.MustListValue(), value.MustListValue(value.StringValue("a")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path))
}

func TestCompareOrderedListSameIsEmpty(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()
	oldValue := value.MustListValue(value.StringValue("one"))

	got, err := CompareAt(rootField("args"), oldValue, oldValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareOrderedListModifiedItem(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("one"), value.StringValue("two")),
		value.MustListValue(value.StringValue("one"), value.StringValue("three")),
		descriptor,
		Options{},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Index(1)))
}

func TestCompareOrderedListModifiedItemUsesIndexPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(
		path,
		value.MustListValue(value.StringValue("one"), value.StringValue("two")),
		value.MustListValue(value.StringValue("one"), value.StringValue("three")),
		descriptor,
		Options{},
	)
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Index(1)))
}

func TestCompareOrderedListAddedItem(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(value.StringValue("one")), value.MustListValue(value.StringValue("one"), value.StringValue("two")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Index(1)), nil, nil)
}

func TestCompareOrderedListAddedItemUsesIndexPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(value.StringValue("one")), value.MustListValue(value.StringValue("one"), value.StringValue("two")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Index(1)), nil, nil)
}

func TestCompareOrderedListRemovedItem(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(value.StringValue("one"), value.StringValue("two")), value.MustListValue(value.StringValue("one")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Index(1)), nil)
}

func TestCompareOrderedListRemovedItemUsesIndexPath(t *testing.T) {
	path := rootField("args")
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(value.StringValue("one"), value.StringValue("two")), value.MustListValue(value.StringValue("one")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Index(1)), nil)
}

func TestCompareOrderedListObjectItemModifiedField(t *testing.T) {
	path := rootField("containers")
	descriptor := types.ListOf(
		types.Object(
			types.Field("name").String().Optional(),
			types.Field("image").String().Optional(),
		),
	).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(imageContainer("v1")), value.MustListValue(imageContainer("v2")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Index(0).Field("image")))
}

func TestCompareOrderedListObjectItemModifiedFieldUsesIndexedFieldPath(t *testing.T) {
	path := rootField("containers")
	descriptor := types.ListOf(
		types.Object(
			types.Field("name").String().Optional(),
			types.Field("image").String().Optional(),
		),
	).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(imageContainer("v1")), value.MustListValue(imageContainer("v2")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Index(0).Field("image")))
}

func TestCompareOrderedListObjectItemAddedField(t *testing.T) {
	path := rootField("containers")
	descriptor := types.ListOf(
		types.Object(types.Field("image").String().Optional()),
	).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(valueObject()), value.MustListValue(valueObject("image", "v1")), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Index(0).Field("image")), nil, nil)
}

func TestCompareOrderedListObjectItemRemovedField(t *testing.T) {
	path := rootField("containers")
	descriptor := types.ListOf(
		types.Object(types.Field("image").String().Optional()),
	).Ordered().Descriptor()

	got, err := CompareAt(path, value.MustListValue(valueObject("image", "v1")), value.MustListValue(valueObject()), descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Index(0).Field("image")), nil)
}
func TestEqualListOrderedDifferentLengthIsFalse(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Descriptor()

	got, err := newComparer(Options{}).equalList(
		rootField("args"),
		value.MustListValue(value.StringValue("one")),
		value.MustListValue(value.StringValue("one"), value.StringValue("two")),
		descriptor,
		0,
	)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalList() = true")
	}
}

func TestEqualListMapIgnoresPhysicalOrder(t *testing.T) {
	oldValue := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Degraded", "False"),
	)
	newValue := value.MustListValue(
		conditionValue("Degraded", "False"),
		conditionValue("Ready", "True"),
	)

	got, err := newComparer(Options{}).equalList(rootField("conditions"), oldValue, newValue, conditionsDescriptor(), 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalList() = false")
	}
}
func TestEqualListByIndexSameItemsIsTrue(t *testing.T) {
	oldList, _ := value.MustListValue(value.StringValue("one")).List()
	newList, _ := value.MustListValue(value.StringValue("one")).List()

	got, err := newComparer(Options{}).equalListByIndex(rootField("args"), oldList, newList, types.String().Descriptor(), 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalListByIndex() = false")
	}
}

func TestEqualListByIndexInvalidElementDescriptorReturnsError(t *testing.T) {
	oldList, _ := value.MustListValue(value.StringValue("one")).List()
	newList, _ := value.MustListValue(value.StringValue("one")).List()

	_, err := newComparer(Options{}).equalListByIndex(rootField("args"), oldList, newList, types.Descriptor{}, 0)

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
