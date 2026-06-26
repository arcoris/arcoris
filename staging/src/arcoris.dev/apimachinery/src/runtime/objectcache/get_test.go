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

func TestGetReportsReceiverReadinessCollectionAndAbsence(t *testing.T) {
	var nilCache *Cache
	_, err := nilCache.Get(testKey("system", 1))
	requireErrorIs(t, err, ErrInvalidCache)

	cache, err := New(namespaceCollection("system"))
	requireNoError(t, err)
	_, err = cache.Get(testKey("system", 1))
	requireErrorIs(t, err, ErrNotReady)
	_, err = cache.Get(testKey("other", 1))
	requireErrorIs(t, err, ErrOutsideCollection)
	requireNoError(t, cache.Replace(nil, collectionRead(t, namespaceCollection("system"), 0)))
	result, err := cache.Get(testKey("system", 1))
	requireNoError(t, err)
	if result.Found || result.Revision != 0 {
		t.Fatalf("missing Get() = %#v; want known absence at revision 0", result)
	}
}

func TestGetReturnsDetachedState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "cached"))

	got, err := cache.Get(key)
	requireNoError(t, err)
	if !got.Found {
		t.Fatalf("Get() Found = false; want true")
	}
	mutateState(&got.State, "mutated")

	again, err := cache.Get(key)
	requireNoError(t, err)
	if !again.Found {
		t.Fatalf("Get() Found = false; want true")
	}
	if desired := desiredString(t, again.State); desired != "cached" {
		t.Fatalf("desired = %q; want cached", desired)
	}
}
