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

func TestGetReturnsFalseWhenNilNotReadyOutsideOrMissing(t *testing.T) {
	var nilCache *Cache
	if _, ok := nilCache.Get(testKey("system", 1)); ok {
		t.Fatalf("nil Get() ok = true; want false")
	}

	cache, err := New(namespaceCollection("system"))
	requireNoError(t, err)
	if _, ok := cache.Get(testKey("system", 1)); ok {
		t.Fatalf("not-ready Get() ok = true; want false")
	}
	if _, ok := cache.Get(testKey("other", 1)); ok {
		t.Fatalf("outside Get() ok = true; want false")
	}
	requireNoError(t, cache.Replace(nil, collectionRead(t, namespaceCollection("system"), 0)))
	if _, ok := cache.Get(testKey("system", 1)); ok {
		t.Fatalf("missing Get() ok = true; want false")
	}
}

func TestGetReturnsDetachedState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "cached"))

	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get() ok = false; want true")
	}
	mutateState(&got, "mutated")

	again, ok := cache.Get(key)
	if !ok {
		t.Fatalf("Get() ok = false; want true")
	}
	if desired := desiredString(t, again); desired != "cached" {
		t.Fatalf("desired = %q; want cached", desired)
	}
}
