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
