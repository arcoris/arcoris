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

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestReplaceResetsHistoryAndRebuildsAtCollectionRevision(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 10, listItem(key, 3, "old-observed-at-ten"))
	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 3, "old-observed-at-ten"), testState(key, 11, "eleven")),
	))

	requireNoError(t, cache.Replace(
		context.Background(),
		collectionRead(t, testCollection(), 30, listItem(key, 17, "replacement")),
	))

	_, err := cache.GetAt(key, 29)
	requireErrorIs(t, err, ErrHistoryUnavailable)
	result, err := cache.GetAt(key, 30)
	requireNoError(t, err)
	if !result.Found || result.Revision != 30 || result.State.Revision != 17 {
		t.Fatalf("GetAt(30) = %#v; want replacement observed at 30 with object revision 17", result)
	}
}

func TestReplaceEmptyCollectionClearsLatestAndHistory(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 5, listItem(key, 5, "live"))

	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), 8)))

	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
	result, err := cache.Get(key)
	requireNoError(t, err)
	if result.Found || result.Revision != 8 {
		t.Fatalf("Get() = %#v; want known current absence at revision 8", result)
	}
	_, err = cache.GetAt(key, 7)
	requireErrorIs(t, err, ErrHistoryUnavailable)
	current, err := cache.GetAt(key, 8)
	requireNoError(t, err)
	if current.Found || current.Revision != 8 {
		t.Fatalf("GetAt(8) = %#v; want known absence at replacement revision", current)
	}
}

func TestHotObjectHistoryDoesNotEvictAnotherObjectHistory(t *testing.T) {
	hot := testKey("system", 1)
	cool := testKey("system", 2)
	cache := readyHistoryCache(
		t,
		2,
		1,
		listItem(hot, 1, "hot-1"),
		listItem(cool, 1, "cool-1"),
	)

	beforeHot := testState(hot, 1, "hot-1")
	for revision := objectstore.Revision(2); revision <= 4; revision++ {
		after := testState(hot, revision, "hot-"+revision.String())
		requireNoError(t, cache.ApplyChange(
			context.Background(),
			objectstore.MustUpdatedChange(hot, beforeHot, after),
		))
		beforeHot = after
	}

	_, err := cache.GetAt(hot, 1)
	requireErrorIs(t, err, ErrHistoryUnavailable)
	coolResult, err := cache.GetAt(cool, 1)
	requireNoError(t, err)
	if !coolResult.Found || desiredString(t, coolResult.State) != "cool-1" {
		t.Fatalf("GetAt(cool, 1) = %#v; want cool history retained", coolResult)
	}
}

func TestApplyChangeRecordsTombstoneButRemovesLatestItem(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 2),
	))

	list, err := cache.List()
	requireNoError(t, err)
	if len(list.Items) != 0 {
		t.Fatalf("List() items = %d; want 0 after delete", len(list.Items))
	}
	deleted, err := cache.GetAt(key, 2)
	requireNoError(t, err)
	if deleted.Found {
		t.Fatalf("GetAt(delete) Found = true; want tombstone absence")
	}
}

func TestFailedApplyChangeLeavesHistoryUnchanged(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "current"))
	before, err := cache.GetAt(key, 1)
	requireNoError(t, err)

	err = cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "different"), testState(key, 2, "after")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
	after, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if !after.Found || desiredString(t, after.State) != desiredString(t, before.State) {
		t.Fatalf("failed ApplyChange mutated history: before=%#v after=%#v", before, after)
	}
	_, err = cache.GetAt(key, 2)
	requireErrorIs(t, err, ErrFutureRevision)
}

func TestFailedReplaceLeavesHistoryUnchanged(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 5, listItem(key, 5, "current"))
	before, err := cache.GetAt(key, 5)
	requireNoError(t, err)

	err = cache.Replace(context.Background(), collectionRead(
		t,
		testCollection(),
		6,
		listItem(key, 5, "first"),
		listItem(key, 5, "duplicate"),
	))

	requireErrorIs(t, err, ErrDuplicateKey)
	after, err := cache.GetAt(key, 5)
	requireNoError(t, err)
	if !after.Found || desiredString(t, after.State) != desiredString(t, before.State) {
		t.Fatalf("failed Replace mutated history: before=%#v after=%#v", before, after)
	}
}
