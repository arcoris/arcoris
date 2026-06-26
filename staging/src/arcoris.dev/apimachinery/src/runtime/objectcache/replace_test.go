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

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

func TestReplaceRejectsNilCache(t *testing.T) {
	var cache *Cache

	err := cache.Replace(context.Background(), collectionRead(t, testCollection(), 0))

	requireErrorIs(t, err, ErrInvalidCache)
}

func TestReplaceHonorsContextCancellation(t *testing.T) {
	cache, err := New(testCollection())
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = cache.Replace(ctx, collectionRead(t, testCollection(), 0))

	requireErrorIs(t, err, context.Canceled)
	if cache.Ready() {
		t.Fatalf("cache became ready after canceled Replace")
	}
}

func TestReplaceInstallsCollectionReadAndPreservesOrder(t *testing.T) {
	first := testKey("system", 3)
	second := testKey("system", 1)
	cache, err := New(testCollection())
	requireNoError(t, err)

	requireNoError(t, cache.Replace(context.Background(), collectionRead(
		t,
		testCollection(),
		7,
		listItem(first, 3, "first"),
		listItem(second, 4, "second"),
	)))

	result, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if result.Revision != 7 {
		t.Fatalf("revision = %s; want 7", result.Revision)
	}
	requireListOrder(t, result, first, second)
}

func TestReplaceStoresDetachedState(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "original")
	read := collectionRead(t, testCollection(), 1, item)
	cache, err := New(testCollection())
	requireNoError(t, err)

	requireNoError(t, cache.Replace(context.Background(), read))
	mutateState(&item.State, "mutated")

	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get() ok = false; want true")
	}
	if desired := desiredString(t, got); desired != "original" {
		t.Fatalf("desired = %q; want original", desired)
	}
}

func TestReplaceRejectsDuplicateKeys(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 5, listItem(testKey("system", 2), 2, "existing"))
	before, _ := cache.List()

	err := cache.Replace(context.Background(), collectionRead(
		t,
		testCollection(),
		6,
		listItem(key, 1, "first"),
		listItem(key, 2, "second"),
	))

	requireErrorIs(t, err, ErrInvalidRead)
	requireErrorIs(t, err, ErrDuplicateKey)
	after, _ := cache.List()
	if after.Revision != before.Revision || len(after.Items) != len(before.Items) {
		t.Fatalf("failed Replace mutated cache: before=%#v after=%#v", before, after)
	}
}

func TestReplaceRejectsOutsideCollection(t *testing.T) {
	cache, err := New(namespaceCollection("system"))
	requireNoError(t, err)

	err = cache.Replace(context.Background(), collectionRead(
		t,
		namespaceCollection("other"),
		1,
		listItem(testKey("other", 1), 1, "outside"),
	))

	requireErrorIs(t, err, ErrInvalidRead)
	requireErrorIs(t, err, ErrOutsideCollection)
	if cache.Ready() {
		t.Fatalf("cache became ready after rejected Replace")
	}
}

func TestReplaceRejectsStaleReadAndAcceptsEqualRevision(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 10, listItem(key, 1, "old"))

	err := cache.Replace(context.Background(), collectionRead(t, testCollection(), 9))
	requireErrorIs(t, err, ErrStaleRead)

	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), 10, listItem(key, 1, "refresh"))))
	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get() ok = false; want true")
	}
	if desired := desiredString(t, got); desired != "refresh" {
		t.Fatalf("desired = %q; want refresh", desired)
	}
}

func TestReplaceRejectsInvalidRead(t *testing.T) {
	cache, err := New(testCollection())
	requireNoError(t, err)

	err = cache.Replace(context.Background(), storewatchapi.CollectionRead{})
	requireErrorIs(t, err, ErrInvalidRead)
	requireErrorIs(t, err, storewatchapi.ErrInvalidCollectionRead)
}
