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

package snapshot

import (
	"testing"
	"time"
)

func TestStoreSnapshotClonesValue(t *testing.T) {
	store := NewStore([]string{"a", "b"}, cloneStrings)

	snap := store.Snapshot()
	snap.Value[0] = "changed"

	next := store.Snapshot()
	if got, want := next.Value[0], "a"; got != want {
		t.Fatalf("internal value changed through snapshot: got %q, want %q", got, want)
	}
}

func TestStoreSnapshotClonesMutableReadModel(t *testing.T) {
	want := mutableReadModelValue("initial-name", "initial-attr", "initial-tag")
	store := NewStore(want, cloneMutableReadModel)

	snap := store.Snapshot()
	mutateMutableReadModel(&snap.Value, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
}

func TestStoreStampedClonesValue(t *testing.T) {
	store := NewStore([]string{"a", "b"}, cloneStrings)

	stamped := store.Stamped()
	stamped.Value[0] = "changed"

	next := store.Stamped()
	if got, want := next.Value[0], "a"; got != want {
		t.Fatalf("internal value changed through stamped snapshot: got %q, want %q", got, want)
	}
}

func TestStoreStampedUsesClock(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(100, 0))
	store := NewStore("value", IdentityClone[string], WithClock(clk))

	stamped := store.Stamped()
	if !stamped.Updated.Equal(time.Unix(100, 0)) {
		t.Fatalf("Updated = %s, want %s", stamped.Updated, time.Unix(100, 0))
	}
}

func TestStoreBadCloneCanLeakMutableState(t *testing.T) {
	store := NewStore([]string{"a"}, IdentityClone[[]string])

	snap := store.Snapshot()
	snap.Value[0] = "changed"

	next := store.Snapshot()
	if got, want := next.Value[0], "changed"; got != want {
		t.Fatalf("bad clone test expected aliasing to be visible: got %q, want %q", got, want)
	}
}
