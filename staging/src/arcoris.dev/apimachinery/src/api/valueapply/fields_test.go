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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valuecompare"
)

func TestChangedAppliedFieldsExactMatch(t *testing.T) {
	changes := valuecompare.MustResult(
		fields(path("$.image")),
		fields(path("$.old")),
		fields(path("$.replicas")),
	)

	got := changedAppliedFields(fields(path("$.replicas"), path("$.same")), changes)

	requireSet(t, got, "$.replicas")
}

func TestChangedAppliedFieldsAppliedAncestorOfChanged(t *testing.T) {
	changes := valuecompare.MustResult(fieldpath.EmptySet(), fieldpath.EmptySet(), fields(path("$.spec.replicas")))

	got := changedAppliedFields(fields(path("$.spec")), changes)

	requireSet(t, got, "$.spec")
}

func TestChangedAppliedFieldsAppliedDescendantOfChanged(t *testing.T) {
	changes := valuecompare.MustResult(fieldpath.EmptySet(), fieldpath.EmptySet(), fields(path("$.spec")))

	got := changedAppliedFields(fields(path("$.spec.replicas")), changes)

	requireSet(t, got, "$.spec.replicas")
}

func TestChangedAppliedFieldsSiblingIgnored(t *testing.T) {
	changes := valuecompare.MustResult(fieldpath.EmptySet(), fieldpath.EmptySet(), fields(path("$.spec.image")))

	got := changedAppliedFields(fields(path("$.metadata.name")), changes)

	requireSet(t, got)
}

func TestChangedAppliedFieldsListMapItemAndFieldOverlap(t *testing.T) {
	itemPath := root().Select(readySelector())
	changes := valuecompare.MustResult(fieldpath.EmptySet(), fieldpath.EmptySet(), fields(readyStatusPath()))

	got := changedAppliedFields(fields(itemPath), changes)

	requireSet(t, got, `$[{"type":"Ready"}]`)
}

func TestChangedAppliedFieldsDoesNotReturnUnappliedChangedPath(t *testing.T) {
	changes := valuecompare.MustResult(fieldpath.EmptySet(), fieldpath.EmptySet(), fields(path("$.spec.image")))

	got := changedAppliedFields(fields(path("$.metadata.name")), changes)

	requireSet(t, got)
}

func TestDroppedFieldsExactAppliedRetainsOld(t *testing.T) {
	got := droppedFields(
		fields(path("$.image"), path("$.replicas")),
		fields(path("$.replicas")),
	)

	requireSet(t, got, "$.image")
}

func TestDroppedFieldsAppliedAncestorRetainsOldDescendant(t *testing.T) {
	got := droppedFields(
		fields(path("$.spec.image")),
		fields(path("$.spec")),
	)

	requireSet(t, got)
}

func TestDroppedFieldsAppliedDescendantDoesNotRetainOldAncestor(t *testing.T) {
	got := droppedFields(
		fields(path("$.spec")),
		fields(path("$.spec.image")),
	)

	requireSet(t, got, "$.spec")
}

func TestDroppedFieldsSiblingDoesNotRetainOld(t *testing.T) {
	got := droppedFields(
		fields(path("$.spec.image")),
		fields(path("$.metadata.name")),
	)

	requireSet(t, got, "$.spec.image")
}

func TestDroppedFieldsListMapItemRetainsItemField(t *testing.T) {
	itemPath := root().Select(readySelector())

	got := droppedFields(
		fields(readyStatusPath()),
		fields(itemPath),
	)

	requireSet(t, got)
}

func TestDroppedFieldsListMapItemFieldDoesNotRetainWholeItem(t *testing.T) {
	itemPath := root().Select(readySelector())

	got := droppedFields(
		fields(itemPath),
		fields(readyStatusPath()),
	)

	requireSet(t, got, `$[{"type":"Ready"}]`)
}

func TestMergeFieldsAppliedUnionDeleted(t *testing.T) {
	got := mergeFields(
		fields(path("$.replicas")),
		fields(path("$.image")),
	)

	requireSet(t, got, "$.image", "$.replicas")
}
