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

func TestCreateAssignsRevision(t *testing.T) {
	store := testStore(t)

	created, err := store.Create(context.Background(), testKey(1), testState("created"))
	requireNoError(t, err)

	if !created.Revision.IsValid() {
		t.Fatalf("created revision is invalid")
	}
	requireDesiredString(t, created, "created")
}

func TestCreateRejectsInvalidKey(t *testing.T) {
	store := testStore(t)

	_, err := store.Create(context.Background(), objectstore.Key{}, testState("created"))
	requireErrorIs(t, err, objectstore.ErrInvalidKey)
}

func TestCreateRejectsInvalidState(t *testing.T) {
	store := testStore(t)

	_, err := store.Create(context.Background(), testKey(1), objectstore.State{})
	requireErrorIs(t, err, objectstore.ErrInvalidState)
}

func TestCreateRejectsForgedRevision(t *testing.T) {
	store := testStore(t)
	state := testState("created")
	state.Revision = 99

	_, err := store.Create(context.Background(), testKey(1), state)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestCreateExistingReturnsAlreadyExists(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	createState(t, store, key, "first")

	_, err := store.Create(context.Background(), key, testState("second"))
	requireErrorIs(t, err, objectstore.ErrAlreadyExists)

	var storeErr *objectstore.Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("error type = %T; want *objectstore.Error", err)
	}
	if storeErr.Reason != objectstore.ErrorReasonAlreadyExists || !storeErr.Key.Equal(key) {
		t.Fatalf("structured error = %#v", storeErr)
	}
}
