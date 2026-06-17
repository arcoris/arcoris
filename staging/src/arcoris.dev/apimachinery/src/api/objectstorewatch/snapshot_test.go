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

package objectstorewatch

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestNewSnapshotAcceptsValidCollectionResult(t *testing.T) {
	collection := testCollection()
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 1)}, 1)

	snapshot, err := NewSnapshot(collection, result)
	requireNoError(t, err)

	if snapshot.Collection() != collection {
		t.Fatalf("Collection() = %#v; want %#v", snapshot.Collection(), collection)
	}
	if snapshot.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", snapshot.Len())
	}
	if snapshot.Revision() != 1 {
		t.Fatalf("Revision() = %s; want 1", snapshot.Revision())
	}
}

func TestNewSnapshotAcceptsEmptyZeroRevisionResult(t *testing.T) {
	snapshot, err := NewSnapshot(testCollection(), objectstore.ListResult{})
	requireNoError(t, err)

	if snapshot.Len() != 0 || !snapshot.Revision().IsZero() {
		t.Fatalf("snapshot len/revision = %d/%s; want 0/0", snapshot.Len(), snapshot.Revision())
	}
}

func TestSnapshotResultReturnsDetachedClone(t *testing.T) {
	snapshot := testSnapshot(t)

	result := snapshot.Result()
	result.Items[0].State.Object.Desired = value.StringValue("mutated")
	result.Items = append(result.Items, objectstore.ListItem{})

	fresh := snapshot.Result()
	if len(fresh.Items) != 1 {
		t.Fatalf("fresh result len = %d; want 1", len(fresh.Items))
	}
	got, ok := fresh.Items[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("fresh desired = %q, %v; want desired, true", got, ok)
	}
}

func TestSnapshotItemsReturnsDetachedItems(t *testing.T) {
	snapshot := testSnapshot(t)

	items := snapshot.Items()
	items[0].State.Object.Desired = value.StringValue("mutated")

	fresh := snapshot.Items()
	got, ok := fresh[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("fresh desired = %q, %v; want desired, true", got, ok)
	}
}

func TestNewSnapshotClonesInputResult(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 1)}, 1)

	snapshot, err := NewSnapshot(testCollection(), result)
	requireNoError(t, err)

	result.Items[0].State.Object.Desired = value.StringValue("mutated")
	got, ok := snapshot.Items()[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("snapshot desired = %q, %v; want desired, true", got, ok)
	}
}

func TestSnapshotBoundaryPreservesCollectionAndRevision(t *testing.T) {
	snapshot := testSnapshot(t)

	boundary := snapshot.Boundary()
	if boundary.Collection() != snapshot.Collection() {
		t.Fatalf("boundary collection = %#v; want %#v", boundary.Collection(), snapshot.Collection())
	}
	if boundary.Revision() != snapshot.Revision() {
		t.Fatalf("boundary revision = %s; want %s", boundary.Revision(), snapshot.Revision())
	}
	if !boundary.IsValid() {
		t.Fatalf("boundary from snapshot is invalid")
	}
}

func TestSnapshotCloneDetachesListData(t *testing.T) {
	snapshot := testSnapshot(t)

	cloned := snapshot.Clone()
	items := cloned.Items()
	items[0].State.Object.Desired = value.StringValue("mutated")

	got, ok := snapshot.Items()[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("source desired = %q, %v; want desired, true", got, ok)
	}
}

func TestSnapshotIsZero(t *testing.T) {
	if !(Snapshot{}).IsZero() {
		t.Fatalf("zero Snapshot reported non-zero")
	}
	if testSnapshot(t).IsZero() {
		t.Fatalf("non-zero Snapshot reported zero")
	}
}

func TestSnapshotIsValid(t *testing.T) {
	if !testSnapshot(t).IsValid() {
		t.Fatalf("valid snapshot reported invalid")
	}
	if (Snapshot{}).IsValid() {
		t.Fatalf("zero snapshot reported valid")
	}
}

func TestSnapshotBoundaryForInvalidSnapshotIsZero(t *testing.T) {
	if boundary := (Snapshot{}).Boundary(); !boundary.IsZero() {
		t.Fatalf("Boundary() = %#v; want zero", boundary)
	}
}

func TestNewSnapshotRejectsInvalidCollection(t *testing.T) {
	_, err := NewSnapshot(objectstore.ListRequest{}, objectstore.ListResult{})

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	requireWatchError(t, err, ErrorReasonInvalidSnapshot, "snapshot.collection")
}

func TestNewSnapshotRejectsInvalidListResult(t *testing.T) {
	result := testListResult([]objectstore.ListItem{{Key: testKey("system", "main")}}, 1)

	_, err := NewSnapshot(testCollection(), result)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
	requireWatchError(t, err, ErrorReasonInvalidSnapshot, "snapshot.result")
}

func TestNewSnapshotRejectsListItemOutsideCollection(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("other", "main", 1)}, 1)

	_, err := NewSnapshot(testNamespaceCollection("system"), result)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
}

func TestNewSnapshotRejectsItemNewerThanResultWatermark(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 2)}, 1)

	_, err := NewSnapshot(testCollection(), result)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestNewSnapshotPreservesUnderlyingObjectstoreError(t *testing.T) {
	_, err := NewSnapshot(testCollection(), testListResult([]objectstore.ListItem{{}}, 1))

	if !errors.Is(err, objectstore.ErrInvalidListResult) {
		t.Fatalf("errors.Is(%v, objectstore.ErrInvalidListResult) = false", err)
	}
}
