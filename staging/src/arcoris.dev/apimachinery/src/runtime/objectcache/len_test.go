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

func TestLenReflectsReadyItems(t *testing.T) {
	var nilCache *Cache
	if got := nilCache.Len(); got != 0 {
		t.Fatalf("nil Len() = %d; want 0", got)
	}

	cache, err := New(testCollection())
	requireNoError(t, err)
	if got := cache.Len(); got != 0 {
		t.Fatalf("not-ready Len() = %d; want 0", got)
	}

	requireNoError(t, cache.Replace(nil, collectionRead(
		t,
		testCollection(),
		2,
		listItem(testKey("system", 1), 1, "one"),
		listItem(testKey("system", 2), 2, "two"),
	)))
	if got := cache.Len(); got != 2 {
		t.Fatalf("ready Len() = %d; want 2", got)
	}
}
