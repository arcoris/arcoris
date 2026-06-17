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

func TestNewCollectionReadAcceptsValidCollectionResult(t *testing.T) {
	collection := testCollection()
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 1)}, 1)

	read, err := NewCollectionRead(collection, result)
	requireNoError(t, err)

	if read.Collection() != collection {
		t.Fatalf("Collection() = %#v; want %#v", read.Collection(), collection)
	}
	if read.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", read.Len())
	}
	if read.Revision() != 1 {
		t.Fatalf("Revision() = %s; want 1", read.Revision())
	}
}

func TestNewCollectionReadAcceptsEmptyZeroRevisionResult(t *testing.T) {
	read, err := NewCollectionRead(testCollection(), objectstore.ListResult{})
	requireNoError(t, err)

	if read.Len() != 0 || !read.Revision().IsZero() {
		t.Fatalf("collection read len/revision = %d/%s; want 0/0", read.Len(), read.Revision())
	}
}

func TestCollectionReadResultReturnsDetachedClone(t *testing.T) {
	read := testCollectionRead(t)

	result := read.Result()
	result.Items[0].State.Object.Desired = value.StringValue("mutated")
	result.Items = append(result.Items, objectstore.ListItem{})

	fresh := read.Result()
	if len(fresh.Items) != 1 {
		t.Fatalf("fresh result len = %d; want 1", len(fresh.Items))
	}
	got, ok := fresh.Items[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("fresh desired = %q, %v; want desired, true", got, ok)
	}
}

func TestCollectionReadItemsReturnsDetachedItems(t *testing.T) {
	read := testCollectionRead(t)

	items := read.Items()
	items[0].State.Object.Desired = value.StringValue("mutated")

	fresh := read.Items()
	got, ok := fresh[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("fresh desired = %q, %v; want desired, true", got, ok)
	}
}

func TestNewCollectionReadClonesInputResult(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 1)}, 1)

	read, err := NewCollectionRead(testCollection(), result)
	requireNoError(t, err)

	result.Items[0].State.Object.Desired = value.StringValue("mutated")
	got, ok := read.Items()[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("collection read desired = %q, %v; want desired, true", got, ok)
	}
}

func TestCollectionReadBoundaryPreservesCollectionAndRevision(t *testing.T) {
	read := testCollectionRead(t)

	boundary := read.Boundary()
	if boundary.Collection() != read.Collection() {
		t.Fatalf("boundary collection = %#v; want %#v", boundary.Collection(), read.Collection())
	}
	if boundary.Revision() != read.Revision() {
		t.Fatalf("boundary revision = %s; want %s", boundary.Revision(), read.Revision())
	}
	if !boundary.IsValid() {
		t.Fatalf("boundary from collection read is invalid")
	}
}

func TestCollectionReadCloneDetachesListData(t *testing.T) {
	read := testCollectionRead(t)

	cloned := read.Clone()
	items := cloned.Items()
	items[0].State.Object.Desired = value.StringValue("mutated")

	got, ok := read.Items()[0].State.Object.Desired.AsString()
	if !ok || got != "desired" {
		t.Fatalf("source desired = %q, %v; want desired, true", got, ok)
	}
}

func TestCollectionReadIsZero(t *testing.T) {
	if !(CollectionRead{}).IsZero() {
		t.Fatalf("zero CollectionRead reported non-zero")
	}
	if testCollectionRead(t).IsZero() {
		t.Fatalf("non-zero CollectionRead reported zero")
	}
}

func TestCollectionReadIsValid(t *testing.T) {
	if !testCollectionRead(t).IsValid() {
		t.Fatalf("valid collection read reported invalid")
	}
	if (CollectionRead{}).IsValid() {
		t.Fatalf("zero collection read reported valid")
	}
}

func TestCollectionReadBoundaryForInvalidCollectionReadIsZero(t *testing.T) {
	if boundary := (CollectionRead{}).Boundary(); !boundary.IsZero() {
		t.Fatalf("Boundary() = %#v; want zero", boundary)
	}
}

func TestNewCollectionReadRejectsInvalidCollection(t *testing.T) {
	_, err := NewCollectionRead(objectstore.ListRequest{}, objectstore.ListResult{})

	requireErrorIs(t, err, ErrInvalidCollectionRead)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
	requireWatchError(t, err, ErrorReasonInvalidCollectionRead, "collection_read.collection")
}

func TestNewCollectionReadRejectsInvalidListResult(t *testing.T) {
	result := testListResult([]objectstore.ListItem{{Key: testKey("system", "main")}}, 1)

	_, err := NewCollectionRead(testCollection(), result)

	requireErrorIs(t, err, ErrInvalidCollectionRead)
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
	requireWatchError(t, err, ErrorReasonInvalidCollectionRead, "collection_read.result")
}

func TestNewCollectionReadRejectsListItemOutsideCollection(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("other", "main", 1)}, 1)

	_, err := NewCollectionRead(testNamespaceCollection("system"), result)

	requireErrorIs(t, err, ErrInvalidCollectionRead)
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
}

func TestNewCollectionReadRejectsItemNewerThanResultWatermark(t *testing.T) {
	result := testListResult([]objectstore.ListItem{testListItem("system", "main", 2)}, 1)

	_, err := NewCollectionRead(testCollection(), result)

	requireErrorIs(t, err, ErrInvalidCollectionRead)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestNewCollectionReadPreservesUnderlyingObjectstoreError(t *testing.T) {
	_, err := NewCollectionRead(testCollection(), testListResult([]objectstore.ListItem{{}}, 1))

	if !errors.Is(err, objectstore.ErrInvalidListResult) {
		t.Fatalf("errors.Is(%v, objectstore.ErrInvalidListResult) = false", err)
	}
}
