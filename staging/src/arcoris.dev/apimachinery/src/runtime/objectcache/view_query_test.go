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
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestViewQueryZeroAllAndNonePredicates(t *testing.T) {
	first := testKey("system", 2)
	second := testKey("system", 1)
	cache := readyCache(t, 8, listItem(first, 1, "first"), listItem(second, 2, "second"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	view := snap.Value

	zero := view.Query(objectquery.Predicate{})
	if zero.Revision != 8 {
		t.Fatalf("zero predicate revision = %s; want 8", zero.Revision)
	}
	requireListOrder(t, zero, first, second)

	all := view.Query(mustPredicate(t, objectquery.All()))
	requireListOrder(t, all, first, second)

	none := view.Query(mustPredicate(t, objectquery.None()))
	if none.Revision != 8 || len(none.Items) != 0 || none.Items != nil {
		t.Fatalf("none predicate result = %#v; want nil items at revision 8", none)
	}
}

func TestViewQueryFiltersLatestItemsInViewOrder(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	third := testKey("other", 1)
	cache := readyCache(
		t,
		9,
		listItem(first, 1, "match"),
		listItem(second, 2, "skip"),
		listItem(third, 3, "match"),
	)

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "match")
	}))

	if result.Revision != 9 {
		t.Fatalf("view Query() revision = %s; want 9", result.Revision)
	}
	requireListOrder(t, result, first, third)
}

func TestViewQueryIgnoresHistoryAndDeletedTombstones(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "match"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "match"), 2),
	))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "match")
	}))

	if result.Revision != 2 || len(result.Items) != 0 {
		t.Fatalf("view Query() = %#v; want no latest matches at revision 2", result)
	}

	historical, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if !historical.Found {
		t.Fatal("cache history was not retained for deleted key")
	}
}

func TestViewQueryReturnsDetachedItems(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 3, listItem(key, 1, "cached"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.Query(objectquery.Predicate{})
	mutateState(&result.Items[0].State, "mutated")

	again := snap.Value.Query(objectquery.Predicate{})
	if desired := desiredString(t, again.Items[0].State); desired != "cached" {
		t.Fatalf("view desired = %q; want cached", desired)
	}

	cacheResult, err := cache.Query(objectquery.Predicate{})
	requireNoError(t, err)
	if desired := desiredString(t, cacheResult.Items[0].State); desired != "cached" {
		t.Fatalf("cache desired = %q; want cached", desired)
	}
}
