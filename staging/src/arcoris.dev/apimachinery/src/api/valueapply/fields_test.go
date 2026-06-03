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

package valueapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/valuecompare"
)

func TestChangedAppliedFieldsUsesAppliedIntersectionWithChanges(t *testing.T) {
	changes := valuecompare.Result{
		Added:    fields(path("$.image")),
		Modified: fields(path("$.replicas")),
		Removed:  fields(path("$.old")),
	}

	got := changedAppliedFields(fields(path("$.replicas"), path("$.same")), changes)

	requireSet(t, got, "$.replicas")
}

func TestDroppedFieldsOldOwnerMinusApplied(t *testing.T) {
	got := droppedFields(
		fields(path("$.image"), path("$.replicas")),
		fields(path("$.replicas")),
	)

	requireSet(t, got, "$.image")
}

func TestMergeFieldsAppliedUnionDeleted(t *testing.T) {
	got := mergeFields(
		fields(path("$.replicas")),
		fields(path("$.image")),
	)

	requireSet(t, got, "$.image", "$.replicas")
}
