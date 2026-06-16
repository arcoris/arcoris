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

func TestSnapshotGetExistingKey(t *testing.T) {
	items := testItems()
	snapshot := mustSnapshot(t, testListResult(8, items...))

	got, ok := snapshot.Get(items[1].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	requireItemOrder(t, []objectstore.ListItem{got}, itemRef{"system", "worker-2", 2})
}

func TestSnapshotGetMissingKey(t *testing.T) {
	snapshot := mustSnapshot(t, testListResult(8, testItems()...))
	missing := testItem("system", "missing", 99, nil, nil).Key

	got, ok := snapshot.Get(missing)
	if ok {
		t.Fatal("Get() ok = true; want false")
	}
	if !got.Key.Equal(objectstore.Key{}) || !got.State.Revision.IsZero() {
		t.Fatalf("Get() item = %#v; want zero item", got)
	}
}

func TestSnapshotGetReturnsDetachedCopy(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(8, items...))

	got, ok := snapshot.Get(items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	mutateItem(&got, "mutated")

	again, ok := snapshot.Get(items[0].Key)
	if !ok {
		t.Fatal("Get() after mutation ok = false; want true")
	}
	if got := desiredString(t, again); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, again, "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
}

func TestSnapshotGetDoesNotMutateSnapshot(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	}
	snapshot := mustSnapshot(t, testListResult(8, items...))

	got, ok := snapshot.Get(items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	got.State.Revision = 99
	mutateItem(&got, "changed")

	all := snapshot.Items()
	requireItemOrder(t, all, itemRef{"system", "worker-1", 1})
	if got := desiredString(t, all[0]); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
}
