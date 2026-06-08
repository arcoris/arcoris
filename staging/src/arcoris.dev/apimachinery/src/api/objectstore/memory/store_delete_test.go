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

package memory

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

	if deleted.Revision != created.Revision {
		t.Fatalf("deleted Revision = %v; want %v", deleted.Revision, created.Revision)
	}
	requireDesiredString(t, deleted, "created")
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
