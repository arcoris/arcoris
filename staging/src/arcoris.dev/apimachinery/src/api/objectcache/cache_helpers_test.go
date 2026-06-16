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

type cacheView struct {
	revision objectstore.Revision
	items    []objectstore.ListItem
}

func mustCache(t *testing.T, result objectstore.ListResult) *Cache {
	t.Helper()

	cache, err := New(result)
	requireNoError(t, err)

	return cache
}

func assertCacheListEquivalent(t *testing.T, cache *Cache, query objectquery.Query) {
	t.Helper()

	source := cache.Items()
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)
	want := predicate.Filter(source)

	got, err := cache.List(query)
	requireNoError(t, err)

	if got.Revision != cache.Revision() {
		t.Fatalf("Revision = %v; want %v", got.Revision, cache.Revision())
	}
	requireSameItems(t, got.Items, want)
}

func assertCacheInvariants(t *testing.T, cache *Cache) {
	t.Helper()

	items := cache.Items()
	seen := map[objectstore.Key]struct{}{}
	for _, item := range items {
		if _, exists := seen[item.Key]; exists {
			t.Fatalf("duplicate key in Items(): %s", item.Key)
		}
		seen[item.Key] = struct{}{}
		got, ok := cache.Get(item.Key)
		if !ok {
			t.Fatalf("Get(%s) ok = false; want true", item.Key)
		}
		requireSameItems(t, []objectstore.ListItem{got}, []objectstore.ListItem{item})
	}
	for _, query := range representativeQueries(t) {
		assertCacheListEquivalent(t, cache, query)
	}
}

func captureCacheView(t *testing.T, cache *Cache) cacheView {
	t.Helper()

	return cacheView{
		revision: cache.Revision(),
		items:    cache.Items(),
	}
}

func requireCacheUnchanged(t *testing.T, cache *Cache, before cacheView) {
	t.Helper()

	if got := cache.Revision(); got != before.revision {
		t.Fatalf("Revision() = %v; want %v", got, before.revision)
	}
	requireSameItems(t, cache.Items(), before.items)
	assertCacheInvariants(t, cache)
}

func requireQueryOrder(t *testing.T, cache *Cache, query objectquery.Query, want ...itemRef) {
	t.Helper()

	result, err := cache.List(query)
	requireNoError(t, err)
	requireItemOrder(t, result.Items, want...)
}
