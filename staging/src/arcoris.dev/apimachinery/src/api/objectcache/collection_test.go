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
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCollectionZeroAccessors(t *testing.T) {
	var col collection

	if !col.isZero() {
		t.Fatal("isZero() = false; want true")
	}
	if got := col.len(); got != 0 {
		t.Fatalf("len() = %d; want 0", got)
	}
	if got, ok := col.item(objectstore.Key{}); ok || !got.Key.Equal(objectstore.Key{}) {
		t.Fatalf("item() = %#v, %v; want zero, false", got, ok)
	}
	if got := col.listItems(); got != nil {
		t.Fatalf("listItems() = %#v; want nil", got)
	}
}

func TestCollectionItemReturnsClone(t *testing.T) {
	col := mustCollection(t, testListResult(
		10,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	got, ok := col.item(col.order[0])
	if !ok {
		t.Fatal("item() ok = false; want true")
	}
	mutateItem(&got, "mutated")

	again, ok := col.item(col.order[0])
	if !ok {
		t.Fatal("second item() ok = false; want true")
	}
	if got := desiredString(t, again); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
}

func mustCollection(t *testing.T, result objectstore.ListResult) collection {
	t.Helper()

	col, err := buildCollection(result, ErrInvalidCache)
	requireNoError(t, err)

	return col
}

func assertCollectionListEquivalent(t *testing.T, col collection, query objectquery.Query) {
	t.Helper()

	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)
	want := predicate.Filter(col.listItems())

	got := col.list(predicate)

	requireSameItems(t, got, want)
}

func assertCollectionInvariants(t *testing.T, col collection) {
	t.Helper()

	seen := map[objectstore.Key]struct{}{}
	for _, key := range col.order {
		if _, exists := seen[key]; exists {
			t.Fatalf("duplicate ordered key: %s", key)
		}
		seen[key] = struct{}{}
		if _, exists := col.items[key]; !exists {
			t.Fatalf("ordered key %s missing from items", key)
		}
	}
	for key := range col.items {
		if _, exists := seen[key]; !exists {
			t.Fatalf("item key %s missing from order", key)
		}
	}
	for _, query := range representativeQueries(t) {
		assertCollectionListEquivalent(t, col, query)
	}
}
