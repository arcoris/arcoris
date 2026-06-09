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
	"arcoris.dev/apimachinery/api/value"
)

func TestCreateDoesNotRetainInputStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	input := testState("created")

	_, err := store.Create(context.Background(), key, input)
	requireNoError(t, err)

	input.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	*input.Object.Observed = value.StringValue("mutated")

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	if got.Ownership.Desired.Entries[0].Fields[0] != objectownership.Path("$.desired") {
		t.Fatalf("stored ownership was mutated: %#v", got.Ownership.Desired.Entries)
	}
	observed, ok := got.Object.Observed.String()
	if !ok || observed != "observed-created" {
		t.Fatalf("observed = %q, %v; want observed-created, true", observed, ok)
	}
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
	got.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	*got.Object.Observed = value.StringValue("mutated")

	again, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	if again.Ownership.Desired.Entries[0].Fields[0] != objectownership.Path("$.desired") {
		t.Fatalf("stored ownership was mutated: %#v", again.Ownership.Desired.Entries)
	}
	observed, ok := again.Object.Observed.String()
	if !ok || observed != "observed-created" {
		t.Fatalf("observed = %q, %v; want observed-created, true", observed, ok)
	}
}

func TestUpdateDoesNotRetainInputStateAliases(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	next := testState("updated")

	_, err := store.Update(context.Background(), key, created.Revision, next)
	requireNoError(t, err)
	next.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"

	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	if got.Ownership.Desired.Entries[0].Fields[0] != objectownership.Path("$.desired") {
		t.Fatalf("stored ownership was mutated: %#v", got.Ownership.Desired.Entries)
	}
}
