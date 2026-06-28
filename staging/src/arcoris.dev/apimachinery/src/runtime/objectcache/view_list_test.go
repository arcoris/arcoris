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

func TestViewListReturnsLatestLiveItemsInOrder(t *testing.T) {
	first := testKey("system", 2)
	second := testKey("system", 1)
	cache := readyCache(t, 11, listItem(first, 1, "first"), listItem(second, 2, "second"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.List()

	if result.Revision != 11 {
		t.Fatalf("view List() revision = %s; want 11", result.Revision)
	}
	requireListOrder(t, result, first, second)
}

func TestViewListExcludesDeletedHistoryRecords(t *testing.T) {
	deleted := testKey("system", 1)
	live := testKey("system", 2)
	cache := readyHistoryCache(
		t,
		2,
		1,
		listItem(deleted, 1, "deleted"),
		listItem(live, 1, "live"),
	)

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(deleted, testState(deleted, 1, "deleted"), 2),
	))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.List()

	if result.Revision != 2 {
		t.Fatalf("view List() revision = %s; want 2", result.Revision)
	}
	requireListOrder(t, result, live)

	historical, err := cache.GetAt(deleted, 1)
	requireNoError(t, err)
	if !historical.Found {
		t.Fatal("cache history was not retained for deleted key")
	}
}

func TestViewListReturnsDetachedItems(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 3, listItem(key, 1, "cached"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result := snap.Value.List()
	mutateState(&result.Items[0].State, "mutated")

	again := snap.Value.List()
	if desired := desiredString(t, again.Items[0].State); desired != "cached" {
		t.Fatalf("view desired = %q; want cached", desired)
	}

	cacheResult, err := cache.List()
	requireNoError(t, err)
	if desired := desiredString(t, cacheResult.Items[0].State); desired != "cached" {
		t.Fatalf("cache desired = %q; want cached", desired)
	}
}
