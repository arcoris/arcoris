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

package objectreflector

import (
	"fmt"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
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
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Worker",
			}),
			meta.ObjectMeta{Name: key.Object.Name, Namespace: key.Object.Namespace},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

func testRead(t testing.TB, revision objectstore.Revision, items ...objectstore.ListItem) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func testReadForCollection(
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

func listItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{Key: key, State: testState(key, revision, desired)}
}

func TestValidateCollectionReadAcceptsOwnedCollection(t *testing.T) {
	read := testRead(t, 1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	err := reflector.validateCollectionRead(read)

	requireNoError(t, err)
}

func TestValidateCollectionReadRejectsInvalidRead(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	err := reflector.validateCollectionRead(storewatchapi.CollectionRead{})

	requireErrorIs(t, err, storewatchapi.ErrInvalidCollectionRead)
}

func TestValidateCollectionReadRejectsDifferentCollection(t *testing.T) {
	read := testReadForCollection(t, namespaceCollection("other"), 1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	err := reflector.validateCollectionRead(read)

	requireErrorIs(t, err, ErrSourceContractViolation)
}
