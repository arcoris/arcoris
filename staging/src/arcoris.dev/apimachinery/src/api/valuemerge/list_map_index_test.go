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

	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/value"
)

func TestListMapIndexPreservesPhysicalOrder(t *testing.T) {
	index, err := newMerger(Options{}).listMapIndex(
		root().Field("conditions"),
		valuepresence.Present(list(
			conditionItem("Ready", "True"),
			conditionItem("Degraded", "False"),
		)),
		conditionElementDescriptor(),
		conditionKeys(),
	)
	if err != nil {
		t.Fatalf("listMapIndex returned error: %v", err)
	}

	if len(index.order) != 2 {
		t.Fatalf("order length = %d; want 2", len(index.order))
	}
	if index.lookup[index.order[0]].item.Kind() != value.KindRecord {
		t.Fatalf("first item kind = %s; want record", index.lookup[index.order[0]].item.Kind())
	}
}

func TestListMapIndexRejectsDuplicateSelector(t *testing.T) {
	_, err := newMerger(Options{}).listMapIndex(
		root().Field("conditions"),
		valuepresence.Present(list(
			conditionItem("Ready", "True"),
			conditionItem("Ready", "False"),
		)),
		conditionElementDescriptor(),
		conditionKeys(),
	)

	requireErrorIs(t, err, ErrDuplicateListKey)
}
