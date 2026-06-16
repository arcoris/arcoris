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
