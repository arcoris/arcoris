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
)

func TestCacheGetItemsAndListReturnDetachedCopies(t *testing.T) {
	cache := mustCache(t, testListResult(
		32,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	got, ok := cache.Get(cache.Items()[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	mutateItem(&got, "get-mutated")

	items := cache.Items()
	mutateItem(&items[0], "items-mutated")

	listed, err := cache.List(objectquery.Query{})
	requireNoError(t, err)
	mutateItem(&listed.Items[0], "list-mutated")

	again := cache.Items()
	requireItemOrder(t, again, itemRef{"system", "worker-1", 1})
	if got := desiredString(t, again[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, again[0], "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}
