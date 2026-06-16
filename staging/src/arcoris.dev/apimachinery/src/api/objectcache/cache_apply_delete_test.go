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
