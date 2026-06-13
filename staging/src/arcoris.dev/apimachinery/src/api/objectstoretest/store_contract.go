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
	"arcoris.dev/apimachinery/api/fieldpath"
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
		{name: "list success", run: testListSuccess},
		{name: "list invalid request", run: testListInvalidRequest},
		{name: "list context", run: testListContext},
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

// testGetMissing verifies absent keys return ok=false without an error.
func testGetMissing(t *testing.T, newStore Factory) {
	_, ok, err := newStore(t).Get(context.Background(), key(1))
	requireNoError(t, err)
	if ok {
		t.Fatalf("Get ok = true; want false")
	}
}

// testCreateSuccess verifies create assigns a committed revision and normalizes ownership.
func testCreateSuccess(t *testing.T, newStore Factory) {
	created, err := newStore(t).Create(context.Background(), key(1), rawState("created"))
	requireNoError(t, err)

	if !created.Revision.IsValid() {
		t.Fatalf("revision = %v; want valid", created.Revision)
	}
	requireDesired(t, created, "created")
	requireNormalizedOwnership(t, created)
}

// testCreateExisting verifies create rejects an already-live key.
func testCreateExisting(t *testing.T, newStore Factory) {
	store := newStore(t)
	_, err := store.Create(context.Background(), key(1), rawState("created"))
	requireNoError(t, err)

	_, err = store.Create(context.Background(), key(1), rawState("again"))
	requireErrorIs(t, err, objectstore.ErrAlreadyExists)
}

// testUpdateSuccess verifies update commits a newer revision and normalized ownership.
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

// testUpdateStale verifies update enforces optimistic revision matching.
func testUpdateStale(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	_, err := store.Update(context.Background(), key(1), created.Revision+1, rawState("updated"))
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

// testDeleteSuccess verifies delete tombstones live state and exposes both revisions.
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

// testDeleteStale verifies delete enforces optimistic revision matching.
func testDeleteStale(t *testing.T, newStore Factory) {
	store := newStore(t)
	created := create(t, store, key(1), "created")

	_, err := store.Delete(context.Background(), key(1), created.Revision+1)
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

// testRecreateAfterDelete verifies tombstoned slots can accept a later create.
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

// testListSuccess verifies the store-level live collection read contract.
func testListSuccess(t *testing.T, newStore Factory) {
	store := newStore(t)
	firstKey := keyWith("workers", "alpha", "one")
	secondKey := keyWith("workers", "beta", "two")
	otherResourceKey := keyWith("jobs", "alpha", "skip")

	first := create(t, store, firstKey, "first")
	second := create(t, store, secondKey, "second")
	create(t, store, otherResourceKey, "other")
	deleted := create(t, store, keyWith("workers", "alpha", "deleted"), "deleted")
	_, err := store.Delete(context.Background(), keyWith("workers", "alpha", "deleted"), deleted.Revision)
	requireNoError(t, err)

	all, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: firstKey.Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	if all.Len() != 2 {
		t.Fatalf("all len = %d; want 2", all.Len())
	}
	requireListItem(t, all.Items[0], firstKey, first.Revision, "first")
	requireListItem(t, all.Items[1], secondKey, second.Revision, "second")
	if all.Revision.IsZero() {
		t.Fatalf("list revision is zero")
	}

	namespace, err := objectstore.InNamespace("alpha")
	requireNoError(t, err)
	namespaced, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: firstKey.Resource,
		Scope:    namespace,
	})
	requireNoError(t, err)
	if namespaced.Len() != 1 {
		t.Fatalf("namespace len = %d; want 1", namespaced.Len())
	}
	requireListItem(t, namespaced.Items[0], firstKey, first.Revision, "first")

	namespaced.Items[0].State.Object.Desired = value.StringValue("mutated")
	again, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: firstKey.Resource,
		Scope:    namespace,
	})
	requireNoError(t, err)
	requireListItem(t, again.Items[0], firstKey, first.Revision, "first")
}

// testListInvalidRequest verifies list requests do not reuse invalid-key errors.
func testListInvalidRequest(t *testing.T, newStore Factory) {
	_, err := newStore(t).List(context.Background(), objectstore.ListRequest{})
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

// testListContext verifies list follows the common store context contract.
func testListContext(t *testing.T, newStore Factory) {
	store := newStore(t)

	_, err := store.List(nil, objectstore.ListRequest{
		Resource: key(1).Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireErrorIs(t, err, objectstore.ErrNilContext)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = store.List(ctx, objectstore.ListRequest{
		Resource: key(1).Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireErrorIs(t, err, context.Canceled)
}

// testStateDetachment verifies stores detach caller input and returned states.
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

// testInvalidInputs verifies common key, state, and revision validation.
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

// testContext verifies nil and canceled contexts are classified consistently.
func testContext(t *testing.T, newStore Factory) {
	store := newStore(t)

	_, _, err := store.Get(nil, key(1))
	requireErrorIs(t, err, objectstore.ErrNilContext)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = store.Get(ctx, key(1))
	requireErrorIs(t, err, context.Canceled)
}

// create commits a test state and fails the current test on error.
func create(t *testing.T, store objectstore.Store, key objectstore.Key, text string) objectstore.State {
	t.Helper()

	created, err := store.Create(context.Background(), key, rawState(text))
	requireNoError(t, err)

	return created
}

// key constructs the default worker key for reusable store contract tests.
func key(index int) objectstore.Key {
	return keyWith("workers", "system", fmt.Sprintf("worker-%d", index))
}

// keyWith constructs a validated key for reusable store contract tests.
func keyWith(resourceName, namespace, name string) objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: apiidentity.Resource(resourceName),
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)
}

// rawState constructs valid uncommitted state with intentionally raw ownership order.
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
		Ownership: objectownership.NewState(fieldownership.MustState(
			fieldownership.MustEntry(fieldownership.MustOwner("z"), fieldSet("$.z")),
			fieldownership.MustEntry(fieldownership.MustOwner("a"), fieldSet("$.a")),
		)),
	}
}

// mutateState changes visible payload values to test store detachment.
func mutateState(state *objectstore.State, text string) {
	state.Object.Desired = value.StringValue(text)
	if state.Object.Observed != nil {
		*state.Object.Observed = value.StringValue("observed-" + text)
	}
}

// fieldSet parses canonical field paths for test ownership entries.
func fieldSet(paths ...string) fieldpath.Set {
	parsed := make([]fieldpath.Path, 0, len(paths))
	for _, text := range paths {
		path, err := fieldpath.ParseCanonical(text)
		if err != nil {
			panic(err)
		}
		parsed = append(parsed, path)
	}

	return fieldpath.MustSet(parsed...)
}

// requireDesired checks the committed Desired string payload.
func requireDesired(t *testing.T, state objectstore.State, want string) {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok || got != want {
		t.Fatalf("desired = %q, %v; want %q, true", got, ok, want)
	}
}

// requireNormalizedOwnership checks committed ownership canonicality.
func requireNormalizedOwnership(t *testing.T, state objectstore.State) {
	t.Helper()

	if err := objectownership.ValidateNormalized(state.Ownership); err != nil {
		t.Fatalf("ownership is not normalized: %v", err)
	}
}

// requireListItem checks the visible storage identity and detached committed state.
func requireListItem(
	t *testing.T,
	item objectstore.ListItem,
	key objectstore.Key,
	revision objectstore.Revision,
	desired string,
) {
	t.Helper()

	if !item.Key.Equal(key) {
		t.Fatalf("item key = %s; want %s", item.Key, key)
	}
	if item.State.Revision != revision {
		t.Fatalf("item revision = %v; want %v", item.State.Revision, revision)
	}
	requireDesired(t, item.State, desired)
	requireNormalizedOwnership(t, item.State)
}

// requireNoError fails the test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs checks sentinel preservation through wrapping.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}
