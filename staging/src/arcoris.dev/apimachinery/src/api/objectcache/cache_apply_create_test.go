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
