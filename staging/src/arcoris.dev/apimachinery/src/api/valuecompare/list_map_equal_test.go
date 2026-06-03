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

	"arcoris.dev/apimachinery/api/value"
)

func TestEqualListMapSameReorderedIsTrue(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.List()
	oldList, _ := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Degraded", "False"),
	).List()
	newList, _ := value.MustListValue(
		conditionValue("Degraded", "False"),
		conditionValue("Ready", "True"),
	).List()

	got, err := newComparer(Options{}).equalListMap(rootField("conditions"), oldList, newList, listView.Element(), listView.MapKeys(), 0)
	requireNoError(t, err)

	if !got {
		t.Fatalf("equalListMap() = false")
	}
}

func TestEqualListMapDifferentSelectorSetIsFalse(t *testing.T) {
	descriptor := conditionsDescriptor()
	listView, _ := descriptor.List()
	oldList, _ := value.MustListValue(conditionValue("Ready", "True")).List()
	newList, _ := value.MustListValue(conditionValue("Degraded", "False")).List()

	got, err := newComparer(Options{}).equalListMap(rootField("conditions"), oldList, newList, listView.Element(), listView.MapKeys(), 0)
	requireNoError(t, err)

	if got {
		t.Fatalf("equalListMap() = true")
	}
}
