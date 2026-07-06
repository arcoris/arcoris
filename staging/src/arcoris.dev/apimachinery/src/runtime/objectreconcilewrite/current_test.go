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

package objectreconcilewrite

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

func TestFromItemAcceptsValidListItem(t *testing.T) {
	key := testKey("task-1")
	item := testItem(key, 7, "desired")

	current, err := FromItem(item)

	requireNoError(t, err)
	if !current.Key().Equal(key) {
		t.Fatalf("key = %#v; want %#v", current.Key(), key)
	}
	if current.Revision() != 7 {
		t.Fatalf("revision = %s; want 7", current.Revision())
	}
}

func TestFromItemRejectsInvalidListItem(t *testing.T) {
	_, err := FromItem(objectstore.ListItem{})

	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
}

func TestFromItemRejectsZeroStateRevision(t *testing.T) {
	item := testItem(testKey("task-1"), 0, "desired")

	_, err := FromItem(item)

	requireErrorIs(t, err, ErrMissingRevision)
}

func TestFromItemDetachesStoredState(t *testing.T) {
	key := testKey("task-1")
	item := testItem(key, 7, "desired")

	current, err := FromItem(item)
	requireNoError(t, err)
	item.State.Object.ObjectMeta.Labels["app"] = "mutated"

	got, ok := current.State().Object.ObjectMeta.Labels.Get("app")
	if !ok {
		t.Fatal("label app is missing")
	}
	if got.String() != "original" {
		t.Fatalf("label app = %q; want original", got)
	}
}

func TestCurrentAccessorsReturnExpectedValues(t *testing.T) {
	key := testKey("task-1")
	current := newCurrent(t, key, 7, "desired")

	if !current.Key().Equal(key) {
		t.Fatalf("key = %#v; want %#v", current.Key(), key)
	}
	if current.Revision() != 7 {
		t.Fatalf("revision = %s; want 7", current.Revision())
	}
	state := current.State()
	if state.Revision != 7 {
		t.Fatalf("state revision = %s; want 7", state.Revision)
	}
}

func TestZeroCurrentIsRejectedByRequestMethods(t *testing.T) {
	var current Current

	_, err := current.UpdateObserved(value.StringValue("observed"), testOwner())
	requireErrorIs(t, err, ErrInvalidCurrent)

	_, err = current.PatchMetadata(nil, nil, testOwner())
	requireErrorIs(t, err, ErrInvalidCurrent)

	_, err = current.Delete()
	requireErrorIs(t, err, ErrInvalidCurrent)
}

func newCurrent(t testing.TB, key objectstore.Key, revision objectstore.Revision, desired string) Current {
	t.Helper()

	current, err := FromItem(testItem(key, revision, desired))
	requireNoError(t, err)

	return current
}

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: testResource(),
		Scope:    objectstore.AllNamespaces(),
	}
}

func testResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "tasks",
	}
}

func testOtherResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func testKey(name string) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: "default",
		Name:      metaidentity.Name(name),
	})
}

func testOutsideKey(name string) objectstore.Key {
	return objectstore.MustKey(testOtherResource(), metaidentity.ObjectName{
		Namespace: "default",
		Name:      metaidentity.Name(name),
	})
}

func testItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{
		Key:   key,
		State: testState(key, revision, desired),
	}
}

func testState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Task",
			}),
			meta.ObjectMeta{
				Name:      key.Object.Name,
				Namespace: key.Object.Namespace,
				Labels: labels.Set{
					"app": "original",
				},
			},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Revision: revision,
	}
}

func testSnapshot(
	t testing.TB,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) objectreconciler.Snapshot {
	t.Helper()

	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)
	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)
	requireNoError(t, cache.Replace(context.Background(), read))
	snapshot, err := cache.ReadSnapshot()
	requireNoError(t, err)

	return objectreconciler.Snapshot{
		Revision: snapshot.Revision,
		View:     snapshot.Value,
	}
}

func testOwner() fieldownership.Owner {
	return fieldownership.MustOwner("test-controller")
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
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}
