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

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestNewEmptyCache(t *testing.T) {
	cache, err := New(objectstore.ListResult{})
	requireNoError(t, err)
	if cache == nil {
		t.Fatal("New() cache = nil; want cache")
	}
	if !cache.IsZero() {
		t.Fatal("IsZero() = false; want true")
	}
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
	if got := cache.Items(); got != nil {
		t.Fatalf("Items() = %#v; want nil", got)
	}
}

func TestNewCacheFromListResult(t *testing.T) {
	items := testItems()
	cache := mustCache(t, testListResult(30, items...))

	if cache.IsZero() {
		t.Fatal("IsZero() = true; want false")
	}
	if got := cache.Len(); got != len(items) {
		t.Fatalf("Len() = %d; want %d", got, len(items))
	}
	if got := cache.Revision(); got != 30 {
		t.Fatalf("Revision() = %v; want 30", got)
	}
	requireSameItems(t, cache.Items(), items)
	assertCacheInvariants(t, cache)
}

func TestNewCacheRejectsDuplicateKeys(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)
	duplicate := item
	duplicate.State.Revision = 2

	cache, err := New(testListResult(3, item, duplicate))

	if cache != nil {
		t.Fatalf("cache = %#v; want nil", cache)
	}
	requireErrorIs(t, err, ErrInvalidCache)
	requireErrorIs(t, err, ErrDuplicateKey)
}

func TestNewCacheClonesInput(t *testing.T) {
	result := testListResult(
		30,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	)
	cache := mustCache(t, result)

	mutateItem(&result.Items[0], "mutated")

	got, ok := cache.Get(result.Items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	if got := desiredString(t, got); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, got, "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestNewCacheBuildsIndexes(t *testing.T) {
	cache := mustCache(t, testListResult(30, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := cache.List(query)
	requireNoError(t, err)

	requireItemOrder(
		t,
		got.Items,
		itemRef{"system", "worker-1", 1},
		itemRef{"other", "worker-3", 3},
	)
}
