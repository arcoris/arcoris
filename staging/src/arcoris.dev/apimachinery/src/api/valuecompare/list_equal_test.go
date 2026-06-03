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

func TestEqualListOrderedDifferentLengthIsFalse(t *testing.T) {
	descriptor := types.ListOf(types.String()).Ordered().Type()

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
