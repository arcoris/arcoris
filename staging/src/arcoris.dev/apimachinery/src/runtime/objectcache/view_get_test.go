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

import "testing"

func TestViewGetMatchesCacheGetSemantics(t *testing.T) {
	key := testKey("system", 1)
	missing := testKey("system", 2)
	outside := testKey("other", 1)
	cache, err := New(namespaceCollection("system"))
	requireNoError(t, err)
	requireNoError(t, cache.Replace(nil, collectionRead(
		t,
		namespaceCollection("system"),
		7,
		listItem(key, 3, "cached"),
	)))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	view := snap.Value

	got, err := view.Get(key)
	requireNoError(t, err)
	if !got.Found || !got.Key.Equal(key) || got.Revision != 7 || desiredString(t, got.State) != "cached" {
		t.Fatalf("view Get(existing) = %#v; want cached object at revision 7", got)
	}

	absent, err := view.Get(missing)
	requireNoError(t, err)
	if absent.Found || !absent.Key.Equal(missing) || absent.Revision != 7 {
		t.Fatalf("view Get(missing) = %#v; want known absence at revision 7", absent)
	}

	_, err = view.Get(outside)
	requireErrorIs(t, err, ErrOutsideCollection)
}

func TestViewGetReturnsDetachedState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 3, listItem(key, 1, "cached"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	got, err := snap.Value.Get(key)
	requireNoError(t, err)
	mutateState(&got.State, "mutated")

	again, err := snap.Value.Get(key)
	requireNoError(t, err)
	if desired := desiredString(t, again.State); desired != "cached" {
		t.Fatalf("view desired = %q; want cached", desired)
	}

	cacheResult, err := cache.Get(key)
	requireNoError(t, err)
	if desired := desiredString(t, cacheResult.State); desired != "cached" {
		t.Fatalf("cache desired = %q; want cached", desired)
	}
}
