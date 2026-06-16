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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCacheGetNilReceiver(t *testing.T) {
	var cache *Cache

	got, ok := cache.Get(objectstore.Key{})

	if ok {
		t.Fatal("Get() ok = true; want false")
	}
	if !got.Key.Equal(objectstore.Key{}) {
		t.Fatalf("Get() item = %#v; want zero item", got)
	}
}

func TestCacheGetExistingAndMissingKey(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(31, items...))

	got, ok := cache.Get(items[1].Key)
	if !ok {
		t.Fatal("Get(existing) ok = false; want true")
	}
	requireItemOrder(t, []objectstore.ListItem{got}, itemRef{"system", "worker-2", 2})

	missing := testItem("system", "missing", 99, nil, nil).Key
	if got, ok := cache.Get(missing); ok || !got.Key.Equal(objectstore.Key{}) {
		t.Fatalf("Get(missing) = %#v, %v; want zero, false", got, ok)
	}
}
