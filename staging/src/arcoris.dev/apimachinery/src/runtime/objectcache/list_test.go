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

func TestListReturnsFalseWhenNilOrNotReady(t *testing.T) {
	var nilCache *Cache
	if result, ok := nilCache.List(); ok || !result.IsZero() {
		t.Fatalf("nil List() = %#v, %v; want zero, false", result, ok)
	}

	cache, err := New(testCollection())
	requireNoError(t, err)
	if result, ok := cache.List(); ok || !result.IsZero() {
		t.Fatalf("not-ready List() = %#v, %v; want zero, false", result, ok)
	}
}

func TestListReturnsDetachedItemsAndRevision(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	cache := readyCache(t, 9, listItem(first, 1, "first"), listItem(second, 2, "second"))

	result, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if result.Revision != 9 {
		t.Fatalf("revision = %s; want 9", result.Revision)
	}
	requireListOrder(t, result, first, second)
	mutateState(&result.Items[0].State, "mutated")

	again, ok := cache.List()
	if !ok {
		t.Fatalf("List() ok = false; want true")
	}
	if desired := desiredString(t, again.Items[0].State); desired != "first" {
		t.Fatalf("desired = %q; want first", desired)
	}
}
