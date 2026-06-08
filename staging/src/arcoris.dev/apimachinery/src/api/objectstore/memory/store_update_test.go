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

func TestUpdateRequiresExpectedRevision(t *testing.T) {
	store := testStore(t)

	_, err := store.Update(context.Background(), testKey(1), 0, testState("next"))
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestUpdateMissingReturnsNotFound(t *testing.T) {
	store := testStore(t)

	_, err := store.Update(context.Background(), testKey(1), 1, testState("next"))
	requireErrorIs(t, err, objectstore.ErrNotFound)

	var storeErr *objectstore.Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("error type = %T; want *objectstore.Error", err)
	}
	if storeErr.Reason != objectstore.ErrorReasonNotFound {
		t.Fatalf("Reason = %v; want %v", storeErr.Reason, objectstore.ErrorReasonNotFound)
	}
}

func TestUpdateRejectsStaleRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "first")

	_, err := store.Update(context.Background(), key, created.Revision+1, testState("next"))
	requireErrorIs(t, err, objectstore.ErrStaleRevision)

	var storeErr *objectstore.Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("error type = %T; want *objectstore.Error", err)
	}
	if storeErr.Actual != created.Revision {
		t.Fatalf("Actual = %v; want %v", storeErr.Actual, created.Revision)
	}
}

func TestUpdateCommitsNewRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "first")

	updated, err := store.Update(context.Background(), key, created.Revision, testState("second"))
	requireNoError(t, err)

	if !created.Revision.Before(updated.Revision) {
		t.Fatalf("updated revision %v did not advance from %v", updated.Revision, created.Revision)
	}
	requireDesiredString(t, updated, "second")
}

func TestUpdateRejectsForgedStateRevision(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "first")
	next := testState("second")
	next.Revision = 99

	_, err := store.Update(context.Background(), key, created.Revision, next)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}
