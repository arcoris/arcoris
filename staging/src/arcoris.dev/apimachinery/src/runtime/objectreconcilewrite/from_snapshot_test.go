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

package objectreconcilewrite

import (
	"testing"

	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

func TestFromSnapshotReturnsCurrentWhenRequestKeyExists(t *testing.T) {
	key := testKey("task-1")
	snapshot := testSnapshot(t, 100, testItem(key, 7, "desired"))

	current, found, err := FromSnapshot(objectreconciler.Request{Key: key}, snapshot)

	requireNoError(t, err)
	if !found {
		t.Fatal("found = false; want true")
	}
	if !current.Key().Equal(key) {
		t.Fatalf("key = %#v; want %#v", current.Key(), key)
	}
	if current.Revision() != 7 {
		t.Fatalf("revision = %s; want 7", current.Revision())
	}
}

func TestFromSnapshotReturnsNotFoundForAbsentKey(t *testing.T) {
	key := testKey("task-1")
	snapshot := testSnapshot(t, 100)

	_, found, err := FromSnapshot(objectreconciler.Request{Key: key}, snapshot)

	requireNoError(t, err)
	if found {
		t.Fatal("found = true; want false")
	}
}

func TestFromSnapshotRejectsInvalidRequest(t *testing.T) {
	_, _, err := FromSnapshot(objectreconciler.Request{}, objectreconciler.Snapshot{})

	requireErrorIs(t, err, ErrInvalidRequest)
	requireErrorIs(t, err, objectreconciler.ErrInvalidRequest)
}

func TestFromSnapshotReturnsViewGetErrorForOutsideCollection(t *testing.T) {
	snapshot := testSnapshot(t, 100)

	_, _, err := FromSnapshot(objectreconciler.Request{Key: testOutsideKey("worker-1")}, snapshot)

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, objectcache.ErrOutsideCollection)
}

func TestFromSnapshotUsesObjectStateRevisionNotSnapshotRevision(t *testing.T) {
	key := testKey("task-1")
	snapshot := testSnapshot(t, 100, testItem(key, 7, "desired"))

	current, found, err := FromSnapshot(objectreconciler.Request{Key: key}, snapshot)
	requireNoError(t, err)
	if !found {
		t.Fatal("found = false; want true")
	}
	req, err := current.Delete()
	requireNoError(t, err)

	if snapshot.Revision != 100 {
		t.Fatalf("snapshot revision = %s; want 100", snapshot.Revision)
	}
	if req.Expected != 7 {
		t.Fatalf("delete expected = %s; want object state revision 7", req.Expected)
	}
}
