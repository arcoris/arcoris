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

	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestCreateDoesNotRetainInputStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	input := testState("created")

	_, err := store.Create(context.Background(), key, input)
	requireNoError(t, err)

	input.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	input.Object.Desired = value.StringValue("mutated")
	*input.Object.Observed = value.StringValue("mutated")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireStateAliasesIntact(t, got, "created")
}

func TestCreateReturnDoesNotExposeStoredStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)

	created, err := store.Create(context.Background(), key, testState("created"))
	requireNoError(t, err)
	mutateStateAliases(&created)

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireStateAliasesIntact(t, got, "created")
}

func TestGetDoesNotExposeStoredStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	createState(t, store, key, "created")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	mutateStateAliases(&got)

	again, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireStateAliasesIntact(t, again, "created")
}

func TestUpdateDoesNotRetainInputStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	next := testState("updated")

	_, err := store.Update(context.Background(), key, created.Revision, next)
	requireNoError(t, err)
	next.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	next.Object.Desired = value.StringValue("mutated")
	*next.Object.Observed = value.StringValue("mutated")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireStateAliasesIntact(t, got, "updated")
}

func TestUpdateReturnDoesNotExposeStoredStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	updated, err := store.Update(context.Background(), key, created.Revision, testState("updated"))
	requireNoError(t, err)
	mutateStateAliases(&updated)

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireStateAliasesIntact(t, got, "updated")
}

func TestDeleteReturnDoesNotExposeTombstoneStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")

	deleted, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)
	mutateStateAliases(&deleted)

	current := store.shardFor(key).get(key).load()
	if current == nil || !current.deleted {
		t.Fatalf("current record = %#v; want tombstone", current)
	}
	requireStateAliasesIntact(t, current.visibleState(), "created")
}

func mutateStateAliases(state *objectstore.State) {
	state.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	state.Object.Desired = value.StringValue("mutated")
	if state.Object.Observed != nil {
		*state.Object.Observed = value.StringValue("mutated")
	}
}

func requireStateAliasesIntact(t *testing.T, state objectstore.State, text string) {
	t.Helper()

	if state.Ownership.Desired.Entries[0].Fields[0] != objectownership.Path("$.desired") {
		t.Fatalf("stored ownership was mutated: %#v", state.Ownership.Desired.Entries)
	}
	requireDesiredString(t, state, text)

	observed, ok := state.Object.Observed.String()
	if !ok || observed != "observed-"+text {
		t.Fatalf("observed = %q, %v; want observed-%s, true", observed, ok, text)
	}
}
