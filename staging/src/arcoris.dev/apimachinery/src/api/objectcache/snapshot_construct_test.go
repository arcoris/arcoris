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

	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestNewSnapshotZero(t *testing.T) {
	snapshot, err := NewSnapshot(objectstore.ListResult{})
	requireNoError(t, err)

	if !snapshot.IsZero() {
		t.Fatal("IsZero() = false; want true")
	}
	if got := snapshot.Items(); got != nil {
		t.Fatalf("Items() = %#v; want nil", got)
	}
}

func TestNewSnapshotEmpty(t *testing.T) {
	snapshot, err := NewSnapshot(objectstore.ListResult{
		Items:    []objectstore.ListItem{},
		Revision: 7,
	})
	requireNoError(t, err)

	if snapshot.IsZero() {
		t.Fatal("IsZero() = true; want false because revision is set")
	}
	if got := snapshot.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
	if got := snapshot.Revision(); got != 7 {
		t.Fatalf("Revision() = %v; want 7", got)
	}
	if got := snapshot.Items(); got != nil {
		t.Fatalf("Items() = %#v; want nil", got)
	}
}

func TestNewSnapshotPreservesRevisionAndOrder(t *testing.T) {
	items := testItems()
	snapshot := mustSnapshot(t, testListResult(11, items...))

	if got := snapshot.Revision(); got != 11 {
		t.Fatalf("Revision() = %v; want 11", got)
	}
	requireItemOrder(
		t,
		snapshot.Items(),
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
		itemRef{"other", "worker-3", 3},
		itemRef{"system", "worker-4", 4},
	)
}

func TestNewSnapshotRejectsDuplicateKeys(t *testing.T) {
	item := testItem("system", "worker-1", 1, nil, nil)
	duplicate := item
	duplicate.State.Revision = 2

	_, err := NewSnapshot(testListResult(3, item, duplicate))

	requireErrorIs(t, err, ErrInvalidSnapshot)
	requireErrorIs(t, err, ErrDuplicateKey)
}

func TestNewSnapshotClonesInputResult(t *testing.T) {
	result := testListResult(
		10,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
	)
	snapshot := mustSnapshot(t, result)

	result.Items[0] = testItem("system", "changed", 99, nil, nil)

	requireItemOrder(
		t,
		snapshot.Items(),
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-2", 2},
	)
}

func TestNewSnapshotDoesNotExposeInputState(t *testing.T) {
	result := testListResult(
		10,
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	)
	snapshot := mustSnapshot(t, result)

	mutateItem(&result.Items[0], "mutated")

	got, ok := snapshot.Get(result.Items[0].Key)
	if !ok {
		t.Fatal("Get() ok = false; want true")
	}
	if got := desiredString(t, got); got != "worker-1" {
		t.Fatalf("Desired = %q; want worker-1", got)
	}
	if got := labelValue(t, got, "env"); got != "prod" {
		t.Fatalf("label env = %q; want prod", got)
	}
	if result.Items[0].State.Object.ObjectMeta.Labels[labels.Key("env")] != labels.Value("mutated") {
		t.Fatal("test mutation did not mutate input label")
	}
}
