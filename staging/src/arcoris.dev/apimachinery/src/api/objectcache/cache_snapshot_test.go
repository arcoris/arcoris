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

func TestCacheSnapshotReturnsDetachedImmutableSnapshot(t *testing.T) {
	cache := mustCache(t, testListResult(
		60,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	snapshot := cache.Snapshot()
	if snapshot.Revision() != 60 {
		t.Fatalf("Snapshot Revision() = %v; want 60", snapshot.Revision())
	}
	items := snapshot.Items()
	mutateItem(&items[0], "mutated")

	again := snapshot.Items()
	if got := desiredString(t, again[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
}

func TestCacheSnapshotUnaffectedByLaterApply(t *testing.T) {
	before := testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil)
	cache := mustCache(t, testListResult(1, before))
	snapshot := cache.Snapshot()
	after := testItem("system", "worker-1", 2, labelsMap("env", "qa"), nil)

	requireNoError(t, cache.Apply(objectstore.MustUpdatedChange(before.Key, before.State, after.State)))

	requireItemOrder(t, snapshot.Items(), itemRef{"system", "worker-1", 1})
	if got := labelValue(t, snapshot.Items()[0], "env"); got != "prod" {
		t.Fatalf("snapshot label env = %q; want prod", got)
	}
	requireQueryOrder(t, cache, objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "qa")),
	}, itemRef{"system", "worker-1", 2})
}

func TestCacheSnapshotUnaffectedByLaterReplace(t *testing.T) {
	cache := mustCache(t, testListResult(1, testItems()...))
	snapshot := cache.Snapshot()

	requireNoError(t, cache.Replace(testListResult(
		70,
		testItem("system", "replacement", 70, nil, nil),
	)))

	requireItemOrder(
		t,
		snapshot.Items(),
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"system", "worker-4", 4},
	)
	requireItemOrder(t, cache.Items(), itemRef{"system", "replacement", 70})
}
