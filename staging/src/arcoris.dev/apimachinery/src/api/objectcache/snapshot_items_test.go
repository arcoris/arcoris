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

func TestSnapshotItemsReturnsAllInOrder(t *testing.T) {
	snapshot := mustSnapshot(t, testListResult(9, testItems()...))

	requireItemOrder(
		t,
		snapshot.Items(),
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"system", "worker-4", 4},
	)
}

func TestSnapshotItemsReturnsDetachedCopies(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(9, items...))

	got := snapshot.Items()
	mutateItem(&got[0], "mutated")
	got = append(got, testItem("system", "extra", 99, nil, nil))

	again := snapshot.Items()
	requireItemOrder(t, again, itemRef{"system", "worker-1", 1})
	if got := desiredString(t, again[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, again[0], "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestSnapshotUnaffectedByReturnedItemsMutation(t *testing.T) {
	snapshot := mustSnapshot(t, testListResult(
		9,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	))

	items := snapshot.Items()
	mutateItem(&items[0], "changed")

	got, ok := snapshot.Get(items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	if got := desiredString(t, got); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
}
