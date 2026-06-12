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

// Package objectstoretest provides reusable contract tests for object stores.
package objectstoretest

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// Factory constructs a fresh object store for one contract test.
type Factory func(testing.TB) objectstore.Store

// RunStoreContractTests runs the shared objectstore.Store behavioral contract.
func RunStoreContractTests(t *testing.T, newStore Factory) {
	t.Helper()

	tests := []struct {
		name string
		run  func(*testing.T, Factory)
	}{
		{name: "get missing", run: testGetMissing},
		{name: "create success", run: testCreateSuccess},
		{name: "create existing", run: testCreateExisting},
		{name: "update success", run: testUpdateSuccess},
		{name: "update stale", run: testUpdateStale},
		{name: "delete success", run: testDeleteSuccess},
		{name: "delete stale", run: testDeleteStale},
		{name: "recreate after delete", run: testRecreateAfterDelete},
		{name: "detachment", run: testStateDetachment},
		{name: "invalid inputs", run: testInvalidInputs},
		{name: "context", run: testContext},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t, newStore)
		})
	}
}

func testGetMissing(t *testing.T, newStore Factory) {
	_, ok, err := newStore(t).Get(context.Background(), key(1))
	requireNoError(t, err)
	if ok {
		t.Fatalf("Get ok = true; want false")
	}
}

func testCreateSuccess(t *testing.T, newStore Factory) {
	created, err := newStore(t).Create(context.Background(), key(1), rawState("created"))
	requireNoError(t, err)

	if !created.Revision.IsValid() {
		t.Fatalf("revision = %v; want valid", created.Revision)
	}
	requireDesired(t, created, "created")
	requireNormalizedOwnership(t, created)
}

func testCreateExisting(t *testing.T, newStore Factory) {
	store := newStore(t)
	_, err := store.Create(context.Background(), key(1), rawState("created"))
	requireNoError(t, err)

	_, err = store.Create(context.Background(), key(1), rawState("again"))
	requireErrorIs(t, err, objectstore.ErrAlreadyExists)
}

func testUpdateSuccess(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	updated, err := store.Update(context.Background(), key(1), created.Revision, rawState("updated"))
	requireNoError(t, err)

	if !created.Revision.Before(updated.Revision) {
		t.Fatalf("updated revision = %v; want after %v", updated.Revision, created.Revision)
	}
	requireDesired(t, updated, "updated")
	requireNormalizedOwnership(t, updated)
}

func testUpdateStale(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	_, err := store.Update(context.Background(), key(1), created.Revision+1, rawState("updated"))
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

func testDeleteSuccess(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	deleted, err := store.Delete(context.Background(), key(1), created.Revision)
	requireNoError(t, err)

	if deleted.Deleted.Revision != created.Revision {
		t.Fatalf("deleted revision = %v; want %v", deleted.Deleted.Revision, created.Revision)
	}
	if !created.Revision.Before(deleted.Revision) {
		t.Fatalf("delete revision = %v; want after %v", deleted.Revision, created.Revision)
	}
	requireDesired(t, deleted.Deleted, "created")
	requireNormalizedOwnership(t, deleted.Deleted)

	_, ok, err := store.Get(context.Background(), key(1))
	requireNoError(t, err)
	if ok {
		t.Fatalf("deleted object is visible")
	}
}

func testDeleteStale(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	_, err := store.Delete(context.Background(), key(1), created.Revision+1)
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

func testRecreateAfterDelete(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")
	deleted, err := store.Delete(context.Background(), key(1), created.Revision)
	requireNoError(t, err)

	recreated, err := store.Create(context.Background(), key(1), rawState("recreated"))
	requireNoError(t, err)

	if !deleted.Revision.Before(recreated.Revision) {
		t.Fatalf("recreated revision = %v; want after delete %v", recreated.Revision, deleted.Revision)
	}
	requireDesired(t, recreated, "recreated")
}

func testStateDetachment(t *testing.T, newStore Factory) {
	store := newStore(t)
	input := rawState("created")
	created, err := store.Create(context.Background(), key(1), input)
	requireNoError(t, err)

	mutateState(&input, "input-mutated")
	mutateState(&created, "return-mutated")

	got, ok, err := store.Get(context.Background(), key(1))
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireDesired(t, got, "created")
	requireNormalizedOwnership(t, got)

	mutateState(&got, "get-mutated")
	again, ok, err := store.Get(context.Background(), key(1))
	requireNoError(t, err)
	if !ok {
		t.Fatalf("Get ok = false")
	}
	requireDesired(t, again, "created")
}

func testInvalidInputs(t *testing.T, newStore Factory) {
	store := newStore(t)

	_, _, err := store.Get(context.Background(), objectstore.Key{})
	requireErrorIs(t, err, objectstore.ErrInvalidKey)

	_, err = store.Create(context.Background(), key(1), objectstore.State{})
	requireErrorIs(t, err, objectstore.ErrInvalidState)

	forged := rawState("forged")
	forged.Revision = 1
	_, err = store.Create(context.Background(), key(1), forged)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)

	_, err = store.Update(context.Background(), key(1), 0, rawState("updated"))
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)

	_, err = store.Delete(context.Background(), key(1), 0)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func testContext(t *testing.T, newStore Factory) {
	store := newStore(t)

	_, _, err := store.Get(nil, key(1))
	requireErrorIs(t, err, objectstore.ErrNilContext)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = store.Get(ctx, key(1))
	requireErrorIs(t, err, context.Canceled)
}

func create(t *testing.T, store objectstore.Store, key objectstore.Key, text string) objectstore.State {
	t.Helper()

	created, err := store.Create(context.Background(), key, rawState(text))
	requireNoError(t, err)

	return created
}

func key(index int) objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{
			Namespace: "system",
			Name:      metaidentity.Name(fmt.Sprintf("worker-%d", index)),
		},
	)
}

func rawState(text string) objectstore.State {
	observed := value.StringValue("observed-" + text)
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{Name: "worker", Namespace: "system"},
			value.StringValue(text),
			observed,
		),
		Ownership: objectownership.Document{
			Version: objectownership.DocumentVersionV1,
			Desired: objectownership.Surface{
				Entries: []objectownership.Entry{
					{Owner: fieldownership.MustOwner("z"), Fields: []objectownership.Path{"$.z"}},
					{Owner: fieldownership.MustOwner("a"), Fields: []objectownership.Path{"$.a", "$.a"}},
				},
			},
		},
	}
}

func mutateState(state *objectstore.State, text string) {
	state.Object.Desired = value.StringValue(text)
	if state.Object.Observed != nil {
		*state.Object.Observed = value.StringValue("observed-" + text)
	}
	if len(state.Ownership.Desired.Entries) > 0 && len(state.Ownership.Desired.Entries[0].Fields) > 0 {
		state.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"
	}
}

func requireDesired(t *testing.T, state objectstore.State, want string) {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok || got != want {
		t.Fatalf("desired = %q, %v; want %q, true", got, ok, want)
	}
}

func requireNormalizedOwnership(t *testing.T, state objectstore.State) {
	t.Helper()

	if err := objectownership.ValidateNormalized(state.Ownership); err != nil {
		t.Fatalf("ownership is not normalized: %v", err)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}
