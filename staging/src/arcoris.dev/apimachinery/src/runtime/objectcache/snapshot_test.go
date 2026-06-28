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

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/snapshot"
)

var _ objectreflector.Sink = (*Cache)(nil)
var _ snapshot.SnapshotReader[objectstore.Revision, View] = (*Cache)(nil)

func TestReadSnapshotRejectsNilOrNotReadyCache(t *testing.T) {
	var nilCache *Cache
	_, err := nilCache.ReadSnapshot()
	requireErrorIs(t, err, ErrInvalidCache)

	cache, err := New(testCollection())
	requireNoError(t, err)
	_, err = cache.ReadSnapshot()
	requireErrorIs(t, err, ErrNotReady)
}

func TestReadSnapshotReturnsDetachedViewAndRevisionInvariant(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	cache := readyCache(t, 9, listItem(first, 1, "first"), listItem(second, 2, "second"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	if snap.Revision != snap.Value.Revision() {
		t.Fatalf("snapshot revision = %s, view revision = %s", snap.Revision, snap.Value.Revision())
	}
	if snap.Revision != 9 {
		t.Fatalf("snapshot revision = %s; want 9", snap.Revision)
	}
	if snap.Value.Len() != 2 {
		t.Fatalf("view Len() = %d; want 2", snap.Value.Len())
	}
	requireListOrder(t, snap.Value.List(), first, second)
}

func TestReadSnapshotReadyEmptyCollection(t *testing.T) {
	cache := readyCache(t, 0)

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	if snap.Revision != 0 || snap.Value.Revision() != 0 {
		t.Fatalf("snapshot = %#v; want revision 0", snap)
	}
	if snap.Value.Len() != 0 {
		t.Fatalf("view Len() = %d; want 0", snap.Value.Len())
	}
	result := snap.Value.List()
	if result.Revision != 0 || len(result.Items) != 0 || result.Items != nil {
		t.Fatalf("view List() = %#v; want nil items at revision 0", result)
	}
}

func TestReadSnapshotExcludesHistoryRecords(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "cached"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "cached"), 2),
	))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	if snap.Value.Len() != 0 {
		t.Fatalf("view Len() = %d; want no latest live items", snap.Value.Len())
	}
	if result := snap.Value.Query(objectquery.Predicate{}); result.Revision != 2 || len(result.Items) != 0 {
		t.Fatalf("view Query() = %#v; want no latest live items at revision 2", result)
	}

	historical, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if !historical.Found {
		t.Fatal("cache history was not retained for deleted key")
	}
}
