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

func TestCacheApplyCreateAddsItemIncrementally(t *testing.T) {
	cache := mustCache(t, testListResult(4, testItems()...))
	created := testItem("system", "worker-5", 5, labelsMap("env", "new"), annotationsMap("team", "new"))
	change := objectstore.MustCreatedChange(created.Key, created.State)

	requireNoError(t, cache.Apply(change))

	if got := cache.Revision(); got != 5 {
		t.Fatalf("Revision() = %v; want 5", got)
	}
	requireItemOrder(
		t,
		cache.Items(),
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"system", "worker-4", 4},
		itemRef{"system", "worker-5", 5},
	)
	assertCacheListEquivalent(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "new")),
	})
	assertCacheListEquivalent(t, cache, objectquery.Query{
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "new")),
	})
	assertCacheInvariants(t, cache)
}

func TestCacheApplyMultipleCreatesPreserveChangeOrder(t *testing.T) {
	cache := mustCache(t, testListResult(1))
	first := testItem("system", "worker-1", 2, nil, nil)
	second := testItem("system", "worker-2", 3, nil, nil)

	requireNoError(t, cache.Apply(objectstore.MustCreatedChange(first.Key, first.State)))
	requireNoError(t, cache.Apply(objectstore.MustCreatedChange(second.Key, second.State)))

	requireItemOrder(
		t,
		cache.Items(),
		itemRef{"system", "worker-1", 2},
		itemRef{"system", "worker-2", 3},
	)
	assertCacheInvariants(t, cache)
}

func TestCacheApplyCreateRejectsExistingKeyAndLeavesStateUnchanged(t *testing.T) {
	cache := mustCache(t, testListResult(4, testItems()...))
	before := captureCacheView(t, cache)
	created := testItem("system", "worker-1", 5, labelsMap("env", "new"), nil)
	change := objectstore.MustCreatedChange(created.Key, created.State)

	err := cache.Apply(change)

	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, before)
}

func TestCacheApplyUpdateMovesLabelAndAnnotationIndexes(t *testing.T) {
	before := testItem("system", "worker-1", 1, labelsMap("env", "prod"), annotationsMap("team", "core"))
	other := testItem("system", "worker-2", 2, labelsMap("env", "prod"), annotationsMap("team", "core"))
	cache := mustCache(t, testListResult(4, before, other))
	after := testItem("system", "worker-1", 5, map[string]string{
		"env":  "qa",
		"tier": "worker",
	}, map[string]string{"team": "tools"})
	change := objectstore.MustUpdatedChange(before.Key, before.State, after.State)

	requireNoError(t, cache.Apply(change))

	if got := cache.Revision(); got != 5 {
		t.Fatalf("Revision() = %v; want 5", got)
	}
	requireItemOrder(
		t,
		cache.Items(),
		itemRef{"system", "worker-1", 5},
		itemRef{"system", "worker-2", 2},
	)
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}, itemRef{"system", "worker-2", 2})
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "qa")),
	}, itemRef{"system", "worker-1", 5})
	requireQueryOrder(t, cache, objectquery.Query{
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "tools")),
	}, itemRef{"system", "worker-1", 5})
	requireQueryOrder(t, cache, objectquery.Query{
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "core")),
	}, itemRef{"system", "worker-2", 2})
	assertCacheInvariants(t, cache)
}

func TestCacheApplyUpdateRejectsMissingOrMismatchedBefore(t *testing.T) {
	cache := mustCache(t, testListResult(4, testItems()...))
	beforeView := captureCacheView(t, cache)
	before := testItem("system", "missing", 1, nil, nil)
	after := testItem("system", "missing", 5, nil, nil)
	err := cache.Apply(objectstore.MustUpdatedChange(before.Key, before.State, after.State))
	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, beforeView)

	beforeView = captureCacheView(t, cache)
	existing := testItems()[0]
	wrongBefore := existing
	wrongBefore.State.Revision = 3
	next := testItem("system", "worker-1", 5, nil, nil)
	err = cache.Apply(objectstore.MustUpdatedChange(existing.Key, wrongBefore.State, next.State))
	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, beforeView)
}

func TestCacheApplyDeleteRemovesItemIndexesAndOrder(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(4, items...))
	change := objectstore.MustDeletedChange(items[0].Key, items[0].State, 5)

	requireNoError(t, cache.Apply(change))

	if got := cache.Revision(); got != 5 {
		t.Fatalf("Revision() = %v; want 5", got)
	}
	if _, ok := cache.Get(items[0].Key); ok {
		t.Fatal("deleted key is still present")
	}
	requireItemOrder(
		t,
		cache.Items(),
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"system", "worker-4", 4},
	)
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}, itemRef{"other", "worker-3", 3})
	requireQueryOrder(t, cache, objectquery.Query{
		Identity: objectquery.IdentitySelector{
			Name: mustNameEquals(t, "worker-1"),
		},
	})
	assertCacheInvariants(t, cache)
}

func TestCacheApplyDeleteRejectsMissingOrMismatchedBefore(t *testing.T) {
	cache := mustCache(t, testListResult(4, testItems()...))
	beforeView := captureCacheView(t, cache)
	missing := testItem("system", "missing", 1, nil, nil)
	err := cache.Apply(objectstore.MustDeletedChange(missing.Key, missing.State, 5))
	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, beforeView)

	beforeView = captureCacheView(t, cache)
	existing := testItems()[0]
	wrongBefore := existing
	wrongBefore.State.Revision = 3
	err = cache.Apply(objectstore.MustDeletedChange(existing.Key, wrongBefore.State, 5))
	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, beforeView)
}

func TestCacheApplyRejectsStaleChangesAndLeavesStateUnchanged(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(10, items...))
	before := captureCacheView(t, cache)
	next := testItem("system", "worker-1", 5, labelsMap("env", "qa"), nil)

	err := cache.Apply(objectstore.MustUpdatedChange(items[0].Key, items[0].State, next.State))
	requireErrorIs(t, err, ErrStaleChange)

	equalRevision := testItem("system", "worker-1", 10, labelsMap("env", "qa"), nil)
	err = cache.Apply(objectstore.MustUpdatedChange(items[0].Key, items[0].State, equalRevision.State))
	requireErrorIs(t, err, ErrStaleChange)

	requireCacheUnchanged(t, cache, before)
}

func TestCacheApplyRejectsMalformedChange(t *testing.T) {
	cache := mustCache(t, testListResult(4, testItems()...))
	before := captureCacheView(t, cache)

	err := cache.Apply(objectstore.Change{})

	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, before)
}

func TestCacheApplyFailureLeavesStateIndexesAndRevisionUnchanged(t *testing.T) {
	cache := mustCache(t, testListResult(10, testItems()...))
	before := captureCacheView(t, cache)
	existing := testItems()[0]
	wrongBefore := existing
	wrongBefore.State.Revision = 7
	next := testItem("system", "worker-1", 11, labelsMap("env", "new"), annotationsMap("team", "new"))

	err := cache.Apply(objectstore.MustUpdatedChange(existing.Key, wrongBefore.State, next.State))

	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorNotIs(t, err, ErrStaleChange)
	requireCacheUnchanged(t, cache, before)
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}, itemRef{"system", "worker-1", 1}, itemRef{"other", "worker-3", 3})
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "new")),
	})
	requireQueryOrder(t, cache, objectquery.Query{
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "new")),
	})
}

func TestCacheListEquivalentAfterCreateUpdateDelete(t *testing.T) {
	item := testItem("system", "worker-1", 1, labelsMap("env", "prod"), annotationsMap("team", "core"))
	cache := mustCache(t, testListResult(1, item))

	created := testItem("system", "worker-2", 2, labelsMap("env", "qa"), annotationsMap("team", "tools"))
	requireNoError(t, cache.Apply(objectstore.MustCreatedChange(created.Key, created.State)))
	assertCacheInvariants(t, cache)

	updated := testItem("system", "worker-1", 3, labelsMap("env", "qa"), annotationsMap("team", "core"))
	requireNoError(t, cache.Apply(objectstore.MustUpdatedChange(item.Key, item.State, updated.State)))
	assertCacheInvariants(t, cache)

	requireNoError(t, cache.Apply(objectstore.MustDeletedChange(created.Key, created.State, 4)))
	assertCacheInvariants(t, cache)
}

func requireQueryOrder(t *testing.T, cache *Cache, query objectquery.Query, want ...itemRef) {
	t.Helper()

	result, err := cache.List(query)
	requireNoError(t, err)
	requireItemOrder(t, result.Items, want...)
}
