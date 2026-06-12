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

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestPrepareKeyedChecksReceiverBeforeContext(t *testing.T) {
	var store *Store

	err := store.prepareKeyed(nil, objectstore.Key{})

	requireErrorIs(t, err, objectstore.ErrUninitializedStore)
}

func TestPrepareKeyedChecksContextBeforeKey(t *testing.T) {
	store := testStore(t)

	err := store.prepareKeyed(nil, objectstore.Key{})

	requireErrorIs(t, err, objectstore.ErrNilContext)
}

func TestPrepareRevisionedChecksRevisionAfterKey(t *testing.T) {
	store := testStore(t)

	err := store.prepareRevisioned(context.Background(), objectstore.Key{}, 0)

	requireErrorIs(t, err, objectstore.ErrInvalidKey)
}

func TestPrepareUpdateChecksRevisionBeforeState(t *testing.T) {
	store := testStore(t)
	state := testState("forged")
	state.Revision = 1

	_, err := store.prepareUpdate(context.Background(), testKey(1), 0, state)

	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestPrepareCreateReturnsDetachedInputState(t *testing.T) {
	store := testStore(t)
	state := testState("prepared")

	prepared, err := store.prepareCreate(context.Background(), testKey(1), state)
	requireNoError(t, err)

	state.Ownership.Desired.Entries[0].Owner = fieldownership.MustOwner("mutated")
	if prepared.Ownership.Desired.Entries[0].Owner != fieldownership.MustOwner("manager") {
		t.Fatalf("prepareCreate retained caller mutation")
	}
}
