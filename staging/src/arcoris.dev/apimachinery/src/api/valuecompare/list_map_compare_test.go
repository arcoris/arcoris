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
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"errors"
	"slices"
	"testing"
)

func TestCompareListMapSameReorderedIsEmpty(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Degraded", "False"),
	)
	newValue := value.MustListValue(
		conditionValue("Degraded", "False"),
		conditionValue("Ready", "True"),
	)

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareListMapModifiedItemField(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func TestCompareListMapModifiedItemFieldUsesSelectorPath(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareListMapAddedItem(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Select(readySelector()).Field("type"), path.Select(readySelector()).Field("status")), nil, nil)
}

func TestCompareListMapAddedItemUsesSelectorFieldSet(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(conditionValue("Progressing", "True"))
	selectorPath := path.Select(progressingSelector())

	got, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(selectorPath.Field("type"), selectorPath.Field("status")), nil, nil)
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareListMapRemovedItem(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, value.MustListValue(), conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Select(readySelector()).Field("type"), path.Select(readySelector()).Field("status")), nil)
}

func TestCompareListMapRemovedItemUsesSelectorFieldSet(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Progressing", "True"))
	selectorPath := path.Select(progressingSelector())

	got, err := CompareAt(path, oldValue, value.MustListValue(), conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(selectorPath.Field("type"), selectorPath.Field("status")), nil)
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareListMapAddedRemovedModifiedKeepsBucketsDisjoint(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(
		conditionValue("Ready", "False"),
		conditionValue("Degraded", "False"),
	)
	newValue := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Progressing", "True"),
	)
	readyPath := path.Select(readySelector())
	degradedPath := path.Select(degradedSelector())
	progressingPath := path.Select(progressingSelector())

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got,
		paths(progressingPath.Field("type"), progressingPath.Field("status")),
		paths(degradedPath.Field("type"), degradedPath.Field("status")),
		paths(readyPath.Field("status")),
	)
}

func TestCompareListMapMultiKeySelector(t *testing.T) {
	path := rootField("routes")
	descriptor := types.ListOf(
		types.Object(
			types.Field("host").String().Required(),
			types.Field("port").Uint64().Required(),
			types.Field("backend").String().Required(),
		),
	).Map("host", "port").Descriptor()
	oldValue := value.MustListValue(routeValue("old"))
	newValue := value.MustListValue(routeValue("new"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(routeSelector()).Field("backend")))
}

func TestCompareListMapDuplicateOldSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "True"), conditionValue("Ready", "False"))

	_, err := CompareAt(path, oldValue, value.MustListValue(), conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
	requireErrorDetailContains(t, err, "first occurrence at $.conditions[0]")
	requireErrorDetailContains(t, err, "duplicate at $.conditions[1]")
}

func TestCompareListMapDuplicateNewSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(conditionValue("Ready", "True"), conditionValue("Ready", "False"))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
}

func TestCompareListMapMissingKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(valueObject("status", "True"))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareListMapNullKeyReturnsInvalidListKey(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(value.MustRecordValue(
		value.MustRecordMember("type", value.NullValue()),
		value.MustRecordMember("status", value.StringValue("True")),
	))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareListMapWrongKeyKindReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(value.MustRecordValue(
		value.MustRecordMember("type", value.BoolValue(true)),
		value.MustRecordMember("status", value.StringValue("True")),
	))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareListMapDoesNotUseIndexPathForSuccessfulSelector(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
	requireNoChangedPathContaining(t, got, "[0]")
}

func TestCompareListMapRefElement(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.dev.Condition": types.Define("example.dev.Condition", conditionExpr()),
	}
	descriptor := types.ListOf(types.Ref("example.dev.Condition")).Map("type").Descriptor()
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func TestCompareListMapRefKeyType(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.dev.ConditionType": types.Define("example.dev.ConditionType", types.String()),
	}
	descriptor := types.ListOf(
		types.Object(
			types.Field("type").Ref("example.dev.ConditionType").Required(),
			types.Field("status").String().Required(),
		),
	).Map("type").Descriptor()
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func routeValue(backend string) value.Value {
	return value.MustRecordValue(
		value.MustRecordMember("host", value.StringValue("api.example.com")),
		value.MustRecordMember("port", value.Uint64Value(443)),
		value.MustRecordMember("backend", value.StringValue(backend)),
	)
}
func TestListMapEntriesIndexesBySelector(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.AsList()
	listValue, _ := value.MustListValue(conditionValue("Ready", "True")).AsList()

	got, err := newComparer(Options{}).listMapEntries(rootField("conditions"), listValue, listView.Element(), listView.MapKeys())
	requireNoError(t, err)

	entry, ok := got[readySelector().String()]
	if !ok {
		t.Fatalf("ready selector missing")
	}
	if !entry.selector.Equal(readySelector()) {
		t.Fatalf("selector = %s, want %s", entry.selector, readySelector())
	}
}

func TestListMapEntriesRejectsEmptyKeys(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.AsList()
	listValue, _ := value.MustListValue(conditionValue("Ready", "True")).AsList()

	_, err := newComparer(Options{}).listMapEntries(rootField("conditions"), listValue, listView.Element(), nil)

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
func TestListMapEntryStoresSelectorIndexAndItem(t *testing.T) {
	item := conditionValue("Ready", "True")
	entry := listMapEntry{
		selector:  readySelector(),
		indexPath: rootField("conditions").Index(0),
		item:      item,
	}

	if !entry.selector.Equal(readySelector()) {
		t.Fatalf("selector = %s, want %s", entry.selector, readySelector())
	}
	if !entry.indexPath.Equal(rootField("conditions").Index(0)) {
		t.Fatalf("indexPath = %s", entry.indexPath)
	}
	equal, err := newComparer(Options{}).equalOpaqueValue(rootField("conditions").Index(0), entry.item, item)
	requireNoError(t, err)
	if !equal {
		t.Fatalf("item mismatch")
	}
}
func TestEqualListMapSameReorderedIsTrue(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.AsList()
	oldList, _ := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Degraded", "False"),
	).AsList()
	newList, _ := value.MustListValue(
		conditionValue("Degraded", "False"),
		conditionValue("Ready", "True"),
	).AsList()

	got, err := newComparer(Options{}).equalListMap(rootField("conditions"), oldList, newList, listView.Element(), listView.MapKeys(), 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalListMap() = false")
	}
}

func TestEqualListMapDifferentSelectorSetIsFalse(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.AsList()
	oldList, _ := value.MustListValue(conditionValue("Ready", "True")).AsList()
	newList, _ := value.MustListValue(conditionValue("Degraded", "False")).AsList()

	got, err := newComparer(Options{}).equalListMap(rootField("conditions"), oldList, newList, listView.Element(), listView.MapKeys(), 0)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalListMap() = true")
	}
}
func TestDuplicateListMapEntryErrorIncludesBothIndexes(t *testing.T) {
	err := duplicateListMapEntryError(
		rootField("conditions").Select(readySelector()),
		rootField("conditions").Index(0),
		rootField("conditions").Index(1),
	)

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorDetailContains(t, err, "$.conditions[0]")
	requireErrorDetailContains(t, err, "$.conditions[1]")
}

func TestCompareListMapKeyErrorMapsInternalError(t *testing.T) {
	err := compareListMapKeyError(&listmapkey.Error{
		Path:   rootField("conditions").Index(0).Field("type"),
		Kind:   listmapkey.FailureMissingKey,
		Detail: "missing",
	})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
}

func TestCompareListMapKeyErrorLeavesUnknownErrorAlone(t *testing.T) {
	cause := errors.New("plain")

	if got := compareListMapKeyError(cause); !errors.Is(got, cause) {
		t.Fatalf("compareListMapKeyError() did not preserve plain error")
	}
}
func TestUnionSortedListMapKeys(t *testing.T) {
	got := unionSortedListMapKeys(
		map[string]listMapEntry{"b": {}, "a": {}},
		map[string]listMapEntry{"c": {}, "a": {}},
	)

	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("unionSortedListMapKeys() = %#v, want %#v", got, want)
	}
}
func TestListMapOperandKeepsPresence(t *testing.T) {
	entry := listMapEntry{item: conditionValue("Ready", "True")}

	got := listMapOperand(entry, true)
	val, ok := got.ValueOK()
	equal, err := newComparer(Options{}).equalOpaqueValue(rootField("conditions").Index(0), val, entry.item)
	requireNoError(t, err)
	if !ok || !equal {
		t.Fatalf("listMapOperand(present) = %#v", got)
	}

	got = listMapOperand(entry, false)
	if got.Present() || !got.Value().IsZero() {
		t.Fatalf("listMapOperand(absent) = %#v", got)
	}
}
func TestExtractListMapSelectorSuccess(t *testing.T) {
	got, err := newComparer(Options{}).extractListMapSelector(
		rootField("conditions").Index(0),
		conditionValue("Ready", "True"),
		conditionDescriptor(),
		[]types.FieldName{"type"},
	)
	requireNoError(t, err)

	if !got.Equal(readySelector()) {
		t.Fatalf("selector = %s, want %s", got, readySelector())
	}
}

func TestExtractListMapSelectorMapsMissingKey(t *testing.T) {
	_, err := newComparer(Options{}).extractListMapSelector(
		rootField("conditions").Index(0),
		valueObject("status", "True"),
		conditionDescriptor(),
		[]types.FieldName{"type"},
	)

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
}
