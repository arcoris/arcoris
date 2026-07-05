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

package objectindex

import (
	"errors"
	"fmt"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/value"
)

func newDesiredIndex(t testing.TB) *Index {
	t.Helper()

	index, err := New(Definition{Name: "desired", Extractor: desiredExtractor()})
	requireNoError(t, err)

	return index
}

func desiredExtractor() Extractor {
	return ExtractorFunc(func(item objectstore.ListItem, emit EmitFunc) error {
		desired, ok := item.State.Object.Desired.AsString()
		if !ok {
			return errors.New("desired is not string")
		}

		return emit(Value(desired))
	})
}

func nameExtractor() Extractor {
	return ExtractorFunc(func(item objectstore.ListItem, emit EmitFunc) error {
		return emit(Value(item.Key.Object.Name))
	})
}

func fixedExtractor(values ...Value) Extractor {
	return ExtractorFunc(func(_ objectstore.ListItem, emit EmitFunc) error {
		for _, value := range values {
			if err := emit(value); err != nil {
				return err
			}
		}

		return nil
	})
}

func errorExtractor(err error) Extractor {
	return ExtractorFunc(func(objectstore.ListItem, EmitFunc) error {
		return err
	})
}

func requireLookupKeys(t testing.TB, index *Index, name Name, value Value, want ...objectstore.Key) {
	t.Helper()

	keys, err := index.Lookup(name, value)
	requireNoError(t, err)
	requireKeys(t, keys, want...)
}

func requireKeys(t testing.TB, got []objectstore.Key, want ...objectstore.Key) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("keys = %d; want %d: %#v", len(got), len(want), got)
	}
	for i, key := range want {
		if !got[i].Equal(key) {
			t.Fatalf("key %d = %#v; want %#v", i, got[i], key)
		}
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
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
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

func testKey(name string) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
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
			},
			value.StringValue(desired),
			value.StringValue(fmt.Sprintf("observed-%s", desired)),
		),
		Revision: revision,
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

func createdChange(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.Change {
	return objectstore.MustCreatedChange(key, testState(key, revision, desired))
}

func updatedChange(key objectstore.Key, beforeRevision objectstore.Revision, before string, afterRevision objectstore.Revision, after string) objectstore.Change {
	return objectstore.MustUpdatedChange(key, testState(key, beforeRevision, before), testState(key, afterRevision, after))
}

func deletedChange(key objectstore.Key, beforeRevision objectstore.Revision, before string, revision objectstore.Revision) objectstore.Change {
	return objectstore.MustDeletedChange(key, testState(key, beforeRevision, before), revision)
}
