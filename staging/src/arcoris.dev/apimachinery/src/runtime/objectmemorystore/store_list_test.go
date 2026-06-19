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
	"fmt"
	"sync"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func TestListEmptyStoreReturnsEmptyResult(t *testing.T) {
	store := testStore(t)

	result, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	if !result.IsZero() {
		t.Fatalf("result = %#v; want zero empty result", result)
	}
}

func TestListFiltersResourceAndScope(t *testing.T) {
	store := testStore(t)
	workerAlpha := listKey("workers", "alpha", "one")
	workerBeta := listKey("workers", "beta", "two")
	jobAlpha := listKey("jobs", "alpha", "skip")
	createState(t, store, workerAlpha, "alpha")
	createState(t, store, workerBeta, "beta")
	createState(t, store, jobAlpha, "job")

	all, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	requireListKeys(t, all, workerAlpha, workerBeta)

	scope, err := objectstore.InNamespace("alpha")
	requireNoError(t, err)
	namespaced, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    scope,
	})
	requireNoError(t, err)
	requireListKeys(t, namespaced, workerAlpha)
}

func TestListExcludesTombstonesAndIncludesRecreatedObjects(t *testing.T) {
	store := testStore(t)
	key := listKey("workers", "system", "main")
	created := createState(t, store, key, "created")

	deleted, err := store.Delete(context.Background(), key, created.Revision)
	requireNoError(t, err)

	empty, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: key.Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	if empty.Len() != 0 {
		t.Fatalf("len after delete = %d; want 0", empty.Len())
	}

	recreated := createState(t, store, key, "recreated")
	if !deleted.Revision.Before(recreated.Revision) {
		t.Fatalf("recreated revision = %v; want after delete %v", recreated.Revision, deleted.Revision)
	}

	result, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: key.Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	requireListKeys(t, result, key)
	requireDesiredString(t, result.Items[0].State, "recreated")
}

func TestListReturnsDeterministicKeyOrder(t *testing.T) {
	store := testStore(t)
	second := listKey("workers", "beta", "b")
	first := listKey("workers", "alpha", "z")
	third := listKey("workers", "beta", "c")
	createState(t, store, third, "third")
	createState(t, store, second, "second")
	createState(t, store, first, "first")

	result, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	requireListKeys(t, result, first, second, third)
}

func TestListReturnsDetachedResult(t *testing.T) {
	store := testStore(t)
	key := listKey("workers", "system", "main")
	createState(t, store, key, "created")

	first, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: key.Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)
	first.Items[0].State.Object.Desired = valueForListTest("mutated")
	first.Items[0] = objectstore.ListItem{}

	second, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: key.Resource,
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	requireListKeys(t, second, key)
	requireDesiredString(t, second.Items[0].State, "created")
}

func TestListReportsCurrentRevision(t *testing.T) {
	store := testStore(t)
	createState(t, store, listKey("workers", "system", "one"), "one")
	second := createState(t, store, listKey("workers", "system", "two"), "two")

	result, err := store.List(context.Background(), objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	})
	requireNoError(t, err)

	if result.Revision != second.Revision {
		t.Fatalf("list revision = %v; want current store revision %v", result.Revision, second.Revision)
	}
}

func TestListRejectsInvalidRequestAndCanceledContext(t *testing.T) {
	store := testStore(t)

	_, err := store.List(context.Background(), objectstore.ListRequest{})
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = store.List(ctx, objectstore.ListRequest{
		Resource: listResource("workers"),
		Scope:    objectstore.AllNamespaces(),
	})
	requireErrorIs(t, err, context.Canceled)
}

func TestListConcurrentOperations(t *testing.T) {
	store := testStore(t, WithShardCount(4))
	ctx := context.Background()
	keys := []objectstore.Key{
		listKey("workers", "alpha", "one"),
		listKey("workers", "alpha", "two"),
		listKey("workers", "beta", "one"),
		listKey("workers", "beta", "two"),
	}
	for i, key := range keys {
		createState(t, store, key, fmt.Sprintf("initial-%d", i))
	}

	var wg sync.WaitGroup
	errs := make(chan error, 64)
	for worker := 0; worker < 4; worker++ {
		worker := worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 64; i++ {
				key := keys[(worker+i)%len(keys)]
				state, ok, err := store.Get(ctx, key)
				if err != nil {
					errs <- err
					continue
				}
				if !ok {
					if _, err := store.Create(ctx, key, testState(fmt.Sprintf("create-%d-%d", worker, i))); !allowedConcurrentListError(err) {
						errs <- err
					}
					continue
				}
				if i%7 == 0 {
					if _, err := store.Delete(ctx, key, state.Revision); !allowedConcurrentListError(err) {
						errs <- err
					}
					continue
				}
				if _, err := store.Update(ctx, key, state.Revision, testState(fmt.Sprintf("update-%d-%d", worker, i))); !allowedConcurrentListError(err) {
					errs <- err
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 128; i++ {
			result, err := store.List(ctx, objectstore.ListRequest{
				Resource: listResource("workers"),
				Scope:    objectstore.AllNamespaces(),
			})
			if err != nil {
				errs <- err
				continue
			}
			for _, item := range result.Items {
				if err := objectstore.ValidateCommittedState(item.State); err != nil {
					errs <- err
				}
				if item.Key.Resource != listResource("workers") {
					errs <- fmt.Errorf("unexpected resource %s", item.Key.Resource)
				}
			}
		}
	}()

	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent operation error: %v", err)
		}
	}
}

// listResource constructs the shared GVR used by memory-store list tests.
func listResource(name string) apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: apiidentity.Resource(name),
	}
}

// listKey constructs a validated key for memory-store list tests.
func listKey(resourceName, namespace, name string) objectstore.Key {
	return objectstore.MustKey(
		listResource(resourceName),
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)
}

// requireListKeys checks both result ordering and committed-state validity.
func requireListKeys(t *testing.T, result objectstore.ListResult, keys ...objectstore.Key) {
	t.Helper()

	if result.Len() != len(keys) {
		t.Fatalf("len = %d; want %d", result.Len(), len(keys))
	}
	for i, key := range keys {
		if !result.Items[i].Key.Equal(key) {
			t.Fatalf("item[%d] key = %s; want %s", i, result.Items[i].Key, key)
		}
		if err := objectstore.ValidateCommittedState(result.Items[i].State); err != nil {
			t.Fatalf("item[%d] state invalid: %v", i, err)
		}
	}
}

// valueForListTest keeps detachment mutations visually tied to list tests.
func valueForListTest(text string) value.Value {
	return value.StringValue(text)
}

// allowedConcurrentListError reports expected optimistic races in the stress test.
func allowedConcurrentListError(err error) bool {
	return err == nil ||
		errors.Is(err, objectstore.ErrAlreadyExists) ||
		errors.Is(err, objectstore.ErrConflict) ||
		errors.Is(err, objectstore.ErrNotFound) ||
		errors.Is(err, objectstore.ErrStaleRevision)
}
