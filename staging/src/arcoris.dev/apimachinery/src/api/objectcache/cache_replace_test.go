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

func TestCacheReplaceReplacesItemsRevisionAndIndexes(t *testing.T) {
	cache := mustCache(t, testListResult(10, testItems()...))
	replacement := []objectstore.ListItem{
		testItem("other", "worker-9", 9, labelsMap("env", "prod"), annotationsMap("team", "new")),
		testItem("system", "worker-8", 8, labelsMap("env", "qa"), annotationsMap("team", "new")),
	}

	requireNoError(t, cache.Replace(testListResult(50, replacement...)))

	if got := cache.Revision(); got != 50 {
		t.Fatalf("Revision() = %v; want 50", got)
	}
	requireItemOrder(
		t,
		cache.Items(),
		itemRef{"other", "worker-9", 9},
		itemRef{"system", "worker-8", 8},
	)
	assertCacheListEquivalent(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	})
	assertCacheListEquivalent(t, cache, objectquery.Query{
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "new")),
	})
}

func TestCacheReplaceRejectsDuplicateKeysAndLeavesCacheUnchanged(t *testing.T) {
	cache := mustCache(t, testListResult(10, testItems()...))
	before := cache.Items()
	item := testItem("system", "worker-1", 20, nil, nil)
	duplicate := item
	duplicate.State.Revision = 21

	err := cache.Replace(testListResult(50, item, duplicate))

	requireErrorIs(t, err, ErrInvalidCache)
	requireErrorIs(t, err, ErrDuplicateKey)
	if got := cache.Revision(); got != 10 {
		t.Fatalf("Revision() = %v; want 10", got)
	}
	requireSameItems(t, cache.Items(), before)
	assertCacheInvariants(t, cache)
}

func TestCacheReplaceClonesInput(t *testing.T) {
	cache := mustCache(t, testListResult(10, testItems()...))
	replacement := testListResult(
		50,
		testItem("system", "worker-1", 20, labelsMap("env", "prod"), nil),
	)

	requireNoError(t, cache.Replace(replacement))
	mutateItem(&replacement.Items[0], "mutated")

	got, ok := cache.Get(replacement.Items[0].Key)
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
