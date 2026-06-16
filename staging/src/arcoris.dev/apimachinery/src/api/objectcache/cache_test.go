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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestNewCacheFromListResult(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(30, items...))

	if cache.IsZero() {
		t.Fatal("IsZero() = true; want false")
	}
	if got := cache.Len(); got != len(items) {
		t.Fatalf("Len() = %d; want %d", got, len(items))
	}
	if got := cache.Revision(); got != 30 {
		t.Fatalf("Revision() = %v; want 30", got)
	}
	requireSameItems(t, cache.Items(), items)
	assertCacheInvariants(t, cache)
}

func TestNewCacheRejectsDuplicateKeys(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)
	duplicate := item
	duplicate.State.Revision = 2

	cache, err := New(testListResult(3, item, duplicate))

	if cache != nil {
		t.Fatalf("cache = %#v; want nil", cache)
	}
	requireErrorIs(t, err, ErrInvalidCache)
	requireErrorIs(t, err, ErrDuplicateKey)
}

func TestNewCacheClonesInput(t *testing.T) {
	result := testListResult(
		30,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	)
	cache := mustCache(t, result)

	mutateItem(&result.Items[0], "mutated")

	got, ok := cache.Get(result.Items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	if got := desiredString(t, got); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, got, "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestNewCacheBuildsIndexes(t *testing.T) {
	cache := mustCache(t, testListResult(30, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := cache.List(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got.Items,
		itemRef{"system", "worker-1", 1},
		itemRef{"other", "worker-3", 3},
	)
}

func TestCacheGetExistingAndMissingKey(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(31, items...))

	got, ok := cache.Get(items[1].Key)
	if !ok {
		t.Fatal("Get(existing) ok = false; want true")
	}
	requireItemOrder(t, []objectstore.ListItem{got}, itemRef{"system", "worker-2", 2})

	missing := testItem("system", "missing", 99, nil, nil).Key
	if got, ok := cache.Get(missing); ok || !got.Key.Equal(objectstore.Key{}) {
		t.Fatalf("Get(missing) = %#v, %v; want zero, false", got, ok)
	}
}

func TestCacheGetItemsAndListReturnDetachedCopies(t *testing.T) {
	cache := mustCache(t, testListResult(
		32,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	got, ok := cache.Get(cache.Items()[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	mutateItem(&got, "get-mutated")

	items := cache.Items()
	mutateItem(&items[0], "items-mutated")

	listed, err := cache.List(objectquery.Query{})
	requireNoError(t, err)
	mutateItem(&listed.Items[0], "list-mutated")

	again := cache.Items()
	requireItemOrder(t, again, itemRef{"system", "worker-1", 1})
	if got := desiredString(t, again[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, again[0], "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestCacheListUsesQueryIndexesAndPredicate(t *testing.T) {
	source := testItems()
	cache := mustCache(t, testListResult(33, source...))
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "identity namespace",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
			},
		},
		{
			name: "identity name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Name: mustNameEquals(t, "worker-3"),
				},
			},
		},
		{
			name: "label positive negative",
			query: objectquery.Query{
				Labels: mustLabelSelector(
					t,
					mustLabelExists(t, "env"),
					mustLabelNotEquals(t, "tier", "frontend"),
				),
			},
		},
		{
			name: "annotation in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "combined",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
				Labels:      mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
				Annotations: mustAnnotationSelector(t, mustAnnotationNotIn(t, "zone", "west")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCacheListEquivalent(t, cache, tt.query)
		})
	}
}

func TestCacheListPreservesOrderAndRevision(t *testing.T) {
	source := []objectstore.ListItem{
		testItem("system", "worker-3", 3, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
		testItem("system", "worker-4", 4, labelsMap("env", "prod"), nil),
	}
	cache := mustCache(t, testListResult(34, source...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := cache.List(query)
	requireNoError(t, err)

	if got.Revision != 34 {
		t.Fatalf("Revision = %v; want 34", got.Revision)
	}
	requireItemOrder(
		t,
		got.Items,
		itemRef{"system", "worker-3", 3},
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-4", 4},
	)
}

func TestCacheConcurrentReaders(t *testing.T) {
	cache := mustCache(t, testListResult(40, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = cache.Get(testItems()[0].Key)
				_ = cache.Items()
				_, _ = cache.List(query)
				_ = cache.Snapshot()
				_ = cache.Revision()
				_ = cache.Len()
			}
		}()
	}
	wg.Wait()
}

func TestCacheConcurrentReadWhileApply(t *testing.T) {
	before := testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil)
	cache := mustCache(t, testListResult(1, before))

	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_, _ = cache.Get(before.Key)
				_, _ = cache.List(objectquery.Query{})
			}
		}
	}()

	current := before
	for rev := objectstore.Revision(2); rev < 40; rev++ {
		next := testItem("system", "worker-1", rev, labelsMap("env", "prod"), nil)
		change := objectstore.MustUpdatedChange(current.Key, current.State, next.State)
		requireNoError(t, cache.Apply(change))
		current = next
	}
	close(stop)
	wg.Wait()
}

func TestCacheConcurrentReadWhileReplace(t *testing.T) {
	cache := mustCache(t, testListResult(1, testItems()...))

	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_ = cache.Items()
				_, _ = cache.List(objectquery.Query{})
			}
		}
	}()

	for rev := objectstore.Revision(2); rev < 40; rev++ {
		items := []objectstore.ListItem{
			testItem("system", "worker-1", rev, labelsMap("env", "prod"), nil),
			testItem("system", "worker-2", rev+100, labelsMap("env", "qa"), nil),
		}
		requireNoError(t, cache.Replace(testListResult(rev, items...)))
	}
	close(stop)
	wg.Wait()
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

type cacheView struct {
	revision objectstore.Revision
	items    []objectstore.ListItem
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
