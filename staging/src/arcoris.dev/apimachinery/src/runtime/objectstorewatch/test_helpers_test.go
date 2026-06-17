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

package objectstorewatch

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
)

func testRuntimeStore(t testing.TB, options ...Option) *Store {
	t.Helper()

	store, err := New(testBackend(t), options...)
	requireNoError(t, err)

	return store
}

func testBackend(t testing.TB) objectstore.Store {
	t.Helper()

	backend, err := objectmemorystore.New()
	requireNoError(t, err)

	return backend
}

func testResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func otherResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "tasks",
	}
}

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.AllNamespaces()}
}

func namespaceCollection(namespace metaidentity.Namespace) objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.MustNamespace(namespace)}
}

func testKey(namespace metaidentity.Namespace, index int) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: namespace,
		Name:      metaidentity.Name(fmt.Sprintf("worker-%d", index)),
	})
}

func otherResourceKey(namespace metaidentity.Namespace, index int) objectstore.Key {
	return objectstore.MustKey(otherResource(), metaidentity.ObjectName{
		Namespace: namespace,
		Name:      metaidentity.Name(fmt.Sprintf("task-%d", index)),
	})
}

func testState(namespace metaidentity.Namespace, name metaidentity.Name, desired string) objectstore.State {
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{Name: name, Namespace: namespace},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
	}
}

func stateForKey(key objectstore.Key, desired string) objectstore.State {
	return testState(key.Object.Namespace, key.Object.Name, desired)
}

func createObject(t testing.TB, store *Store, key objectstore.Key, desired string) objectstore.State {
	t.Helper()

	state, err := store.Create(context.Background(), key, stateForKey(key, desired))
	requireNoError(t, err)

	return state
}

func updateObject(t testing.TB, store *Store, key objectstore.Key, expected objectstore.Revision, desired string) objectstore.State {
	t.Helper()

	state, err := store.Update(context.Background(), key, expected, stateForKey(key, desired))
	requireNoError(t, err)

	return state
}

func deleteObject(t testing.TB, store *Store, key objectstore.Key, expected objectstore.Revision) objectstore.DeleteResult {
	t.Helper()

	result, err := store.Delete(context.Background(), key, expected)
	requireNoError(t, err)

	return result
}

func watchAfter(t testing.TB, store *Store, collection objectstore.ListRequest, revision objectstore.Revision) objectwatch.Stream {
	t.Helper()

	start, err := objectwatch.AfterRevision(revision)
	requireNoError(t, err)
	stream, err := store.Watch(context.Background(), objectwatch.Request{Collection: collection, Start: start})
	requireNoError(t, err)

	return stream
}

func watchAtCurrent(t testing.TB, store *Store, collection objectstore.ListRequest) objectwatch.Stream {
	t.Helper()

	stream, err := store.Watch(context.Background(), objectwatch.Request{
		Collection: collection,
		Start:      objectwatch.AtCurrent(),
	})
	requireNoError(t, err)

	return stream
}

func watchRequestAfter(t testing.TB, collection objectstore.ListRequest, revision objectstore.Revision) objectwatch.Request {
	t.Helper()

	start, err := objectwatch.AfterRevision(revision)
	requireNoError(t, err)

	return objectwatch.Request{Collection: collection, Start: start}
}

func nextEvent(t testing.TB, stream objectwatch.Stream) objectwatch.Event {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	event, err := stream.Next(ctx)
	requireNoError(t, err)

	return event
}

func requireNoEvent(t testing.TB, stream objectwatch.Stream) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := stream.Next(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Next() error = %v; want context deadline", err)
	}
}

func requireNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireDesiredString(t testing.TB, state objectstore.State, want string) {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok || got != want {
		t.Fatalf("desired = %q, %v; want %q, true", got, ok, want)
	}
}

func requireChangedEvent(
	t testing.TB,
	event objectwatch.Event,
	request objectwatch.Request,
	kind objectstore.ChangeKind,
	revision objectstore.Revision,
) {
	t.Helper()

	if event.Kind != objectwatch.EventChanged {
		t.Fatalf("event kind = %s; want changed", event.Kind)
	}
	if event.Change.Kind != kind {
		t.Fatalf("change kind = %s; want %s", event.Change.Kind, kind)
	}
	if event.Revision != revision || event.Change.Revision != revision {
		t.Fatalf("revision = event %s change %s; want %s", event.Revision, event.Change.Revision, revision)
	}
	validator, err := objectwatch.NewValidator(request)
	requireNoError(t, err)
	requireNoError(t, validator.Accept(event))
}
