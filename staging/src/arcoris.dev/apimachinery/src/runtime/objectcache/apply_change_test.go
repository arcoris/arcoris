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

func TestApplyChangeRejectsNilCache(t *testing.T) {
	var cache *Cache
	key := testKey("system", 1)

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(key, testState(key, 1, "created")),
	)

	requireErrorIs(t, err, ErrInvalidCache)
}

func TestApplyChangeHonorsContextCancellation(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := cache.ApplyChange(
		ctx,
		objectstore.MustCreatedChange(key, testState(key, 1, "created")),
	)

	requireErrorIs(t, err, context.Canceled)
	if cache.Len() != 0 {
		t.Fatalf("cache mutated after canceled ApplyChange")
	}
}

func TestApplyChangeBeforeReplaceReturnsNotReady(t *testing.T) {
	cache, err := New(testCollection())
	requireNoError(t, err)
	key := testKey("system", 1)

	err = cache.ApplyChange(context.Background(), objectstore.MustCreatedChange(key, testState(key, 1, "created")))

	requireErrorIs(t, err, ErrNotReady)
}

func TestApplyCreatedChangeAppendsItemAndAdvancesRevision(t *testing.T) {
	first := testKey("system", 1)
	created := testKey("system", 2)
	cache := readyCache(t, 1, listItem(first, 1, "first"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(created, testState(created, 2, "created")),
	))

	result, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if result.Revision != 2 {
		t.Fatalf("revision = %s; want 2", result.Revision)
	}
	requireListOrder(t, result, first, created)
}

func TestApplyUpdatedChangePreservesPositionAndAdvancesRevision(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	before := listItem(first, 1, "before")
	cache := readyCache(t, 2, before, listItem(second, 2, "second"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(first, before.State, testState(first, 3, "after")),
	))

	result, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if result.Revision != 3 {
		t.Fatalf("revision = %s; want 3", result.Revision)
	}
	requireListOrder(t, result, first, second)
	got, ok := cache.Get(first)
	if !ok {
		t.Fatalf("Get(updated) ok = false; want true")
	}
	if desired := desiredString(t, got); desired != "after" {
		t.Fatalf("desired = %q; want after", desired)
	}
}

func TestApplyDeletedChangeRemovesItemAndPreservesRemainingOrder(t *testing.T) {
	first := listItem(testKey("system", 1), 1, "first")
	second := listItem(testKey("system", 2), 2, "second")
	third := listItem(testKey("system", 3), 3, "third")
	cache := readyCache(t, 3, first, second, third)

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(second.Key, second.State, 4),
	))

	result, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if result.Revision != 4 {
		t.Fatalf("revision = %s; want 4", result.Revision)
	}
	requireListOrder(t, result, first.Key, third.Key)
	if _, ok := cache.Get(second.Key); ok {
		t.Fatalf("deleted key still present")
	}
}

func TestApplyChangeClonesRetainedState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 0)
	change := objectstore.MustCreatedChange(key, testState(key, 1, "created"))

	requireNoError(t, cache.ApplyChange(context.Background(), change))
	mutateState(&change.After, "mutated")

	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get() ok = false; want true")
	}
	if desired := desiredString(t, got); desired != "created" {
		t.Fatalf("desired = %q; want created", desired)
	}
}

func TestApplyChangeRejectsStaleRevision(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "current")
	cache := readyCache(t, 5, item)

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, item.State, testState(key, 4, "stale")),
	)

	requireErrorIs(t, err, ErrStaleChange)
}

func TestApplyChangeRejectsOutsideCollection(t *testing.T) {
	cache := readyCache(t, 0)
	key := otherResourceKey("system", 1)

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(key, testState(key, 1, "outside")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorIs(t, err, ErrOutsideCollection)
}

func TestApplyChangeRejectsInvalidChange(t *testing.T) {
	cache := readyCache(t, 0)

	err := cache.ApplyChange(context.Background(), objectstore.Change{})

	requireErrorIs(t, err, ErrInvalidChange)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
}

func TestApplyCreateRejectsExistingKey(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 5, listItem(key, 1, "current"))

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(key, testState(key, 6, "created")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestApplyUpdateRejectsMissingKey(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1)

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "before"), testState(key, 2, "after")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestApplyDeleteRejectsMissingKey(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1)

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "before"), 2),
	)

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestApplyUpdateRejectsMismatchedBeforeState(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "current")
	cache := readyCache(t, 1, item)
	before := testState(key, 1, "different")

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, before, testState(key, 2, "after")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestApplyDeleteRejectsMismatchedBeforeState(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "current")
	cache := readyCache(t, 1, item)
	before := testState(key, 1, "different")

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, before, 2),
	)

	requireErrorIs(t, err, ErrInvalidChange)
}

func TestFailedApplyChangeLeavesStateUnchanged(t *testing.T) {
	key := testKey("system", 1)
	item := listItem(key, 1, "current")
	cache := readyCache(t, 1, item)
	before, _ := cache.List()

	err := cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "different"), testState(key, 2, "after")),
	)

	requireErrorIs(t, err, ErrInvalidChange)
	after, _ := cache.List()
	if after.Revision != before.Revision || desiredString(t, after.Items[0].State) != "current" {
		t.Fatalf("failed ApplyChange mutated cache: before=%#v after=%#v", before, after)
	}
}
