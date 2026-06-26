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

package objectcache

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/value"
)

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
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

func testState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	labelSet, err := labels.FromStrings(map[string]string{"env": desired})
	if err != nil {
		panic(err)
	}

	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      key.Object.Name,
				Namespace: key.Object.Namespace,
				Labels:    labelSet,
			},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

func listItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{Key: key, State: testState(key, revision, desired)}
}

func collectionRead(
	t testing.TB,
	collection objectstore.ListRequest,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(collection, objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func readyCache(t testing.TB, revision objectstore.Revision, items ...objectstore.ListItem) *Cache {
	t.Helper()

	cache, err := New(testCollection())
	requireNoError(t, err)
	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), revision, items...)))

	return cache
}

func readyHistoryCache(
	t testing.TB,
	retained int,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) *Cache {
	t.Helper()

	cache, err := New(
		testCollection(),
		WithHistory(HistoryPolicy{RetainedVersionsPerObject: retained}),
	)
	requireNoError(t, err)
	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), revision, items...)))

	return cache
}

func desiredString(t testing.TB, state objectstore.State) string {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok {
		t.Fatalf("desired value is not string: %#v", state.Object.Desired)
	}

	return got
}

func mutateState(state *objectstore.State, desired string) {
	state.Object.Desired = value.StringValue(desired)
	if state.Object.ObjectMeta.Labels == nil {
		state.Object.ObjectMeta.Labels = labels.Set{}
	}
	state.Object.ObjectMeta.Labels[labels.Key("env")] = labels.Value(desired)
}

func requireListOrder(t testing.TB, result objectstore.ListResult, keys ...objectstore.Key) {
	t.Helper()

	if len(result.Items) != len(keys) {
		t.Fatalf("items = %d; want %d", len(result.Items), len(keys))
	}
	for i, key := range keys {
		if !result.Items[i].Key.Equal(key) {
			t.Fatalf("item[%d] key = %#v; want %#v", i, result.Items[i].Key, key)
		}
	}
}

func waitUntil(t testing.TB, ctx context.Context, condition func() bool) {
	t.Helper()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatalf("condition not met before context done: %v", ctx.Err())
		case <-ticker.C:
		}
	}
}
