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

package objectmemorystore

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestDeleteRequiresExpectedRevision(t *testing.T) {
	store := testStore(t)

	_, err := store.Delete(context.Background(), testKey(1), 0)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestDeleteMissingReturnsNotFound(t *testing.T) {
	store := testStore(t)

	_, err := store.Delete(context.Background(), testKey(1), 1)
	requireErrorIs(t, err, objectstore.ErrNotFound)
}

func TestDeleteRejectsStaleRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	_, err := store.Delete(context.Background(), key, created.Revision+1)
	requireErrorIs(t, err, objectstore.ErrStaleRevision)

	var storeErr *objectstore.Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("error type = %T; want *objectstore.Error", err)
	}
	if storeErr.Reason != objectstore.ErrorReasonStaleRevision || storeErr.Actual != created.Revision {
		t.Fatalf("structured error = %#v", storeErr)
	}
}

func TestDeleteReturnsDeletedLiveState(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	deleted, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	if deleted.Deleted.Revision != created.Revision {
		t.Fatalf("deleted Revision = %v; want %v", deleted.Deleted.Revision, created.Revision)
	}
	requireDesiredString(t, deleted.Deleted, "created")
}

func TestDeleteExposesDeleteRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	deleted, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	if deleted.Deleted.Revision != created.Revision {
		t.Fatalf("deleted Revision = %v; want live revision %v", deleted.Deleted.Revision, created.Revision)
	}
	if !deleted.Revision.IsValid() || !created.Revision.Before(deleted.Revision) {
		t.Fatalf("delete revision = %v; want committed revision after %v", deleted.Revision, created.Revision)
	}

	current := store.shardFor(key).get(key).load()
	if current == nil || !current.deleted {
		t.Fatalf("current record = %#v; want tombstone", current)
	}
	if !current.deleteRevision.IsValid() || !created.Revision.Before(current.deleteRevision) {
		t.Fatalf("deleteRevision = %v; want committed revision after %v", current.deleteRevision, created.Revision)
	}
	if current.deleteRevision != deleted.Revision {
		t.Fatalf("tombstone revision = %v; want returned %v", current.deleteRevision, deleted.Revision)
	}
	if current.state.Revision != created.Revision {
		t.Fatalf("tombstone state revision = %v; want deleted live revision %v", current.state.Revision, created.Revision)
	}
}

func TestDeleteResultIsZero(t *testing.T) {
	var result objectstore.DeleteResult
	if !result.IsZero() {
		t.Fatalf("zero DeleteResult is not zero")
	}

	result = objectstore.DeleteResult{Deleted: validDeleteState(), Revision: 2}
	if result.IsZero() {
		t.Fatalf("populated DeleteResult is zero")
	}
}

func validDeleteState() objectstore.State {
	state := testState("deleted")
	state.Revision = 1

	return state
}

func TestDeleteMakesObjectInvisible(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	_, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	_, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if ok {
		t.Fatalf("deleted object is visible")
	}
}
