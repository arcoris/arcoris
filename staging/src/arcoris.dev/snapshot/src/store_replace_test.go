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

	panicassert "arcoris.dev/testutil/panic"
)

func TestStoreReplaceClonesInput(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)
	next := []string{"next"}

	snap := store.Replace(next)
	next[0] = "changed"

	if got, want := snap.Revision, LocalRevision(2); got != want {
		t.Fatalf("Replace revision = %d, want %d", got, want)
	}

	loaded := store.Snapshot()
	if got, want := loaded.Value[0], "next"; got != want {
		t.Fatalf("stored value = %q, want %q", got, want)
	}
}

func TestStoreReplaceClonesMutableReadModelInput(t *testing.T) {
	store := NewStore(mutableReadModelValue("initial-name", "initial-attr", "initial-tag"), cloneMutableReadModel)
	next := mutableReadModelValue("next-name", "next-attr", "next-tag")
	want := cloneMutableReadModel(next)

	_ = store.Replace(next)
	mutateMutableReadModel(&next, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
}

func TestStoreReplaceReturnsIsolatedValue(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)

	snap := store.Replace([]string{"next"})
	snap.Value[0] = "changed-through-result"

	loaded := store.Snapshot()
	if got, want := loaded.Value[0], "next"; got != want {
		t.Fatalf("stored value = %q, want %q", got, want)
	}
}

func TestStoreReplaceReturnsIsolatedMutableReadModel(t *testing.T) {
	store := NewStore(mutableReadModelValue("initial-name", "initial-attr", "initial-tag"), cloneMutableReadModel)
	want := mutableReadModelValue("next-name", "next-attr", "next-tag")

	snap := store.Replace(want)
	mutateMutableReadModel(&snap.Value, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
}

func TestStoreReplaceStampedUpdatesTime(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(1, 0))
	store := NewStore("initial", IdentityClone[string], WithClock(clk))

	clk.set(time.Unix(2, 0))
	stamped := store.ReplaceStamped("next")

	if !stamped.Updated.Equal(time.Unix(2, 0)) {
		t.Fatalf("Updated = %s, want %s", stamped.Updated, time.Unix(2, 0))
	}
}

func TestStoreReplacePanicsOnRevisionOverflowWithoutCommit(t *testing.T) {
	store := NewStore("initial", IdentityClone[string])
	store.mu.Lock()
	store.revision = ^LocalRevision(0)
	store.mu.Unlock()

	panicassert.RequireMessage(t, "snapshot: local revision overflow", func() {
		_ = store.Replace("next")
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, ^LocalRevision(0); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreReplaceStoredClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 2,
		clone: IdentityClone[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Replace("next")
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreReplaceReturnedClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 3,
		clone: IdentityClone[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Replace("next")
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}
