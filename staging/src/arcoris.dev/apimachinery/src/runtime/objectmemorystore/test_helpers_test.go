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

func testStore(t testing.TB, opts ...Option) *Store {
	t.Helper()

	store, err := New(opts...)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	return store
}

func testKey(index int) objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{Namespace: "system", Name: metaidentity.Name(fmt.Sprintf("worker-%d", index))},
	)
}

func testState(text string) objectstore.State {
	observed := value.StringValue("observed-" + text)

	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   "control.arcoris.dev",
				Version: "v1",
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      "worker",
				Namespace: "system",
			},
			value.StringValue(text),
			observed,
		),
		Ownership: objectownership.NewState(ownershipState(ownershipEntry("manager", "$.desired"))),
	}
}

func ownershipEntry(owner string, paths ...string) fieldownership.Entry {
	return fieldownership.MustEntry(fieldownership.MustOwner(owner), fieldSet(paths...))
}

func ownershipState(entries ...fieldownership.Entry) fieldownership.State {
	return fieldownership.MustState(entries...)
}

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

func requireDesiredString(t *testing.T, state objectstore.State, want string) {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok {
		t.Fatalf("desired value is not string")
	}
	if got != want {
		t.Fatalf("desired = %q; want %q", got, want)
	}
}

func createState(t *testing.T, store *Store, key objectstore.Key, text string) objectstore.State {
	t.Helper()

	created, err := store.Create(context.Background(), key, testState(text))
	requireNoError(t, err)

	return created
}
