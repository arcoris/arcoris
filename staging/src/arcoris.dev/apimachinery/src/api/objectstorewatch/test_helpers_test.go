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
	"strings"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
)

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
	return objectstore.ListRequest{
		Resource: testResource(),
		Scope:    objectstore.AllNamespaces(),
	}
}

func testNamespaceCollection(namespace metaidentity.Namespace) objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: testResource(),
		Scope:    objectstore.MustNamespace(namespace),
	}
}

func testKey(namespace metaidentity.Namespace, name metaidentity.Name) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: namespace,
		Name:      name,
	})
}

func testListItem(namespace metaidentity.Namespace, name metaidentity.Name, revision objectstore.Revision) objectstore.ListItem {
	return objectstore.ListItem{
		Key:   testKey(namespace, name),
		State: testCommittedState(namespace, name, revision, "desired"),
	}
}

func testCommittedState(
	namespace metaidentity.Namespace,
	name metaidentity.Name,
	revision objectstore.Revision,
	desired string,
) objectstore.State {
	return objectstore.State{
		Object: object.New[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			value.StringValue(desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

func testListResult(items []objectstore.ListItem, revision objectstore.Revision) objectstore.ListResult {
	return objectstore.ListResult{Items: items, Revision: revision}
}

func testSnapshot(t testing.TB) Snapshot {
	t.Helper()

	snapshot, err := NewSnapshot(
		testCollection(),
		testListResult([]objectstore.ListItem{testListItem("system", "main", 1)}, 1),
	)
	requireNoError(t, err)

	return snapshot
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

func requireWatchError(t testing.TB, err error, reason ErrorReason, pathPart string) {
	t.Helper()

	var watchErr *Error
	if !errors.As(err, &watchErr) {
		t.Fatalf("errors.As(%v, *Error) = false", err)
	}
	if watchErr.Reason != reason {
		t.Fatalf("reason = %s; want %s", watchErr.Reason, reason)
	}
	if pathPart != "" && !strings.Contains(watchErr.Path, pathPart) {
		t.Fatalf("path = %q; want to contain %q", watchErr.Path, pathPart)
	}
}

type fakeSnapshotter struct{}

func (*fakeSnapshotter) Snapshot(context.Context, objectstore.ListRequest) (Snapshot, error) {
	return Snapshot{}, nil
}

type fakeListerWatcher struct {
	fakeSnapshotter
}

func (*fakeListerWatcher) Watch(context.Context, objectwatch.Request) (objectwatch.Stream, error) {
	return nil, nil
}

type fakeStore struct {
	fakeListerWatcher
}

func (*fakeStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	return objectstore.State{}, false, nil
}

func (*fakeStore) Create(context.Context, objectstore.Key, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, nil
}

func (*fakeStore) Update(context.Context, objectstore.Key, objectstore.Revision, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, nil
}

func (*fakeStore) Delete(context.Context, objectstore.Key, objectstore.Revision) (objectstore.DeleteResult, error) {
	return objectstore.DeleteResult{}, nil
}

func (*fakeStore) List(context.Context, objectstore.ListRequest) (objectstore.ListResult, error) {
	return objectstore.ListResult{}, nil
}

type fakeCapableStore struct {
	fakeStore
}

func (*fakeCapableStore) WatchCapabilities() objectwatch.Capabilities {
	return objectwatch.Capabilities{HistoricalStart: true}
}
