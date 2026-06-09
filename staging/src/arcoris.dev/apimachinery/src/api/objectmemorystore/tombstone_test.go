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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestUpdateTombstoneReturnsNotFound(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	_, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	_, err = store.Update(context.Background(), key, created.Revision, testState("updated"))
	requireErrorIs(t, err, objectstore.ErrNotFound)
}

func TestDeleteTombstoneReturnsNotFound(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	_, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	_, err = store.Delete(context.Background(), key, created.Revision)
	requireErrorIs(t, err, objectstore.ErrNotFound)
}

func TestCreateAfterDeleteRecreatesObject(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "first")
	_, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	recreated, err := store.Create(context.Background(), key, testState("second"))
	requireNoError(t, err)

	if !created.Revision.Before(recreated.Revision) {
		t.Fatalf("recreated revision %v did not advance from %v", recreated.Revision, created.Revision)
	}
	requireDesiredString(t, recreated, "second")
}
