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

	"arcoris.dev/apimachinery/api/value"
)

func TestMergeListMapSelectedItemField(t *testing.T) {
	got, err := mergeConditions(
		list(conditionItem("Ready", "False"), conditionItem("Degraded", "False")),
		list(conditionItem("Ready", "True")),
		pathSet(conditionPath("Ready").Field(testFieldName("status"))),
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireConditionStatus(t, got, "Ready", "True")
	requireConditionStatus(t, got, "Degraded", "False")
}

func TestMergeListMapSameReorderedIsPreserved(t *testing.T) {
	base := list(conditionItem("Ready", "True"), conditionItem("Degraded", "False"))
	overlay := list(conditionItem("Degraded", "False"), conditionItem("Ready", "True"))

	got, err := mergeConditions(
		base,
		overlay,
		pathSet(conditionPath("Ready").Field(testFieldName("status"))),
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireConditionOrder(t, got, "Ready", "Degraded")
}

func TestMergeListMapAddsSelectedItem(t *testing.T) {
	got, err := mergeConditions(
		list(conditionItem("Ready", "True")),
		list(conditionItem("Ready", "True"), conditionItem("Progressing", "True")),
		pathSet(conditionPath("Progressing")),
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireConditionOrder(t, got, "Ready", "Progressing")
	requireConditionStatus(t, got, "Progressing", "True")
}

func TestMergeListMapRemovesSelectedItemAbsentFromOverlay(t *testing.T) {
	got, err := mergeConditions(
		list(conditionItem("Ready", "True"), conditionItem("Degraded", "False")),
		list(conditionItem("Ready", "True")),
		pathSet(conditionPath("Degraded")),
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireConditionOrder(t, got, "Ready")
}

func TestMergeListMapDuplicateBaseSelectorReturnsDuplicateListKey(t *testing.T) {
	_, err := mergeConditions(
		list(conditionItem("Ready", "True"), conditionItem("Ready", "False")),
		list(conditionItem("Ready", "True")),
		pathSet(conditionPath("Ready")),
	)

	requireErrorIs(t, err, ErrDuplicateListKey)
}

func TestMergeListMapDuplicateOverlaySelectorReturnsDuplicateListKey(t *testing.T) {
	_, err := mergeConditions(
		list(conditionItem("Ready", "True")),
		list(conditionItem("Ready", "True"), conditionItem("Ready", "False")),
		pathSet(conditionPath("Ready")),
	)

	requireErrorIs(t, err, ErrDuplicateListKey)
}

func TestMergeListMapMissingKeyReturnsInvalidListKey(t *testing.T) {
	_, err := mergeConditions(
		list(obj(member("status", str("True")))),
		list(conditionItem("Ready", "True")),
		pathSet(conditionPath("Ready")),
	)

	requireErrorIs(t, err, ErrInvalidListKey)
}

func TestMergeListMapWrongKeyKindReturnsInvalidListKey(t *testing.T) {
	_, err := mergeConditions(
		list(obj(member("type", intValue(1)), member("status", str("True")))),
		list(conditionItem("Ready", "True")),
		pathSet(conditionPath("Ready")),
	)

	requireErrorIs(t, err, ErrInvalidListKey)
}

func TestMergeListMapPreservesBaseOrderAndAppendsAddedOverlayItems(t *testing.T) {
	got, err := mergeConditions(
		list(conditionItem("Ready", "True"), conditionItem("Degraded", "False")),
		list(conditionItem("Progressing", "True"), conditionItem("Ready", "True")),
		pathSet(conditionPath("Progressing")),
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireConditionOrder(t, got, "Ready", "Degraded", "Progressing")
}

func requireConditionStatus(t *testing.T, listValue value.Value, conditionType string, want string) {
	t.Helper()

	item := requireCondition(t, listValue, conditionType)
	requireStringMember(t, item, "status", want)
}

func requireConditionOrder(t *testing.T, listValue value.Value, want ...string) {
	t.Helper()

	view, ok := listValue.AsList()
	if !ok {
		t.Fatalf("value kind = %s; want list", listValue.Kind())
	}
	if view.Len() != len(want) {
		t.Fatalf("list length = %d; want %d", view.Len(), len(want))
	}

	for i, conditionType := range want {
		item, _ := view.At(i)
		requireStringMember(t, item, "type", conditionType)
	}
}

func requireCondition(t *testing.T, listValue value.Value, conditionType string) value.Value {
	t.Helper()

	view, ok := listValue.AsList()
	if !ok {
		t.Fatalf("value kind = %s; want list", listValue.Kind())
	}

	for i := 0; i < view.Len(); i++ {
		item, _ := view.At(i)
		got, ok := requireMember(t, item, "type").AsString()
		if ok && got == conditionType {
			return item
		}
	}

	t.Fatalf("condition %q is absent", conditionType)
	return value.Value{}
}
