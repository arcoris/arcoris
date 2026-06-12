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

package objectmemorystore_test

import (
	"context"
	"errors"
	"fmt"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func ExampleStore_createGetUpdateDelete() {
	ctx := context.Background()
	store, _ := objectmemorystore.New()
	key := exampleKey()

	created, _ := store.Create(ctx, key, exampleState("first"))
	got, _, _ := store.Get(ctx, key)
	updated, _ := store.Update(ctx, key, got.Revision, exampleState("second"))
	deleted, _ := store.Delete(ctx, key, updated.Revision)

	fmt.Println(created.Revision.IsValid())
	fmt.Println(got.Revision == created.Revision)
	fmt.Println(deleted.Deleted.Revision == updated.Revision)
	fmt.Println(updated.Revision.Before(deleted.Revision))

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleStore_staleRevision() {
	ctx := context.Background()
	store, _ := objectmemorystore.New()
	key := exampleKey()

	created, _ := store.Create(ctx, key, exampleState("first"))
	_, _ = store.Update(ctx, key, created.Revision, exampleState("second"))
	_, err := store.Update(ctx, key, created.Revision, exampleState("third"))

	fmt.Println(errors.Is(err, objectstore.ErrStaleRevision))

	// Output:
	// true
}

func ExampleStore_recreateAfterDelete() {
	ctx := context.Background()
	store, _ := objectmemorystore.New()
	key := exampleKey()

	created, _ := store.Create(ctx, key, exampleState("first"))
	_, _ = store.Delete(ctx, key, created.Revision)
	recreated, _ := store.Create(ctx, key, exampleState("second"))

	fmt.Println(created.Revision.Before(recreated.Revision))

	// Output:
	// true
}

func exampleKey() objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{Namespace: "system", Name: "main"},
	)
}

func exampleState(text string) objectstore.State {
	return objectstore.State{
		Object: object.New[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{Name: "main", Namespace: "system"},
			value.StringValue(text),
		),
		Ownership: objectownership.EmptyState(),
	}
}
