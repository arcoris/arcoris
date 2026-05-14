/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

import (
	"slices"
	"testing"
	"time"
)

func TestNewStoreClonesInitialValue(t *testing.T) {
	initial := []string{"a", "b"}
	store := NewStore(initial, cloneStrings)

	initial[0] = "changed"

	snap := store.Snapshot()
	if got, want := snap.Value[0], "a"; got != want {
		t.Fatalf("Snapshot value[0] = %q, want %q", got, want)
	}
}

func TestNewStoreInitialRevisionIsOne(t *testing.T) {
	store := NewStore("value", Identity[string])

	if got, want := store.Revision(), Revision(1); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestStoreSnapshotClonesValue(t *testing.T) {
	store := NewStore([]string{"a", "b"}, cloneStrings)

	snap := store.Snapshot()
	snap.Value[0] = "changed"

	next := store.Snapshot()
	if got, want := next.Value[0], "a"; got != want {
		t.Fatalf("internal value changed through snapshot: got %q, want %q", got, want)
	}
}

func TestStoreStampedUsesClock(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(100, 0))
	store := NewStore("value", Identity[string], WithClock(clk))

	stamped := store.Stamped()
	if !stamped.Updated.Equal(time.Unix(100, 0)) {
		t.Fatalf("Updated = %s, want %s", stamped.Updated, time.Unix(100, 0))
	}
}

func TestStoreReplaceClonesInput(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)
	next := []string{"next"}

	snap := store.Replace(next)
	next[0] = "changed"

	if got, want := snap.Revision, Revision(2); got != want {
		t.Fatalf("Replace revision = %d, want %d", got, want)
	}

	loaded := store.Snapshot()
	if got, want := loaded.Value[0], "next"; got != want {
		t.Fatalf("stored value = %q, want %q", got, want)
	}
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

func TestStoreReplaceStampedUpdatesTime(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(1, 0))
	store := NewStore("initial", Identity[string], WithClock(clk))

	clk.set(time.Unix(2, 0))
	stamped := store.ReplaceStamped("next")

	if !stamped.Updated.Equal(time.Unix(2, 0)) {
		t.Fatalf("Updated = %s, want %s", stamped.Updated, time.Unix(2, 0))
	}
}

func TestStoreReplacePanicsOnRevisionOverflowWithoutCommit(t *testing.T) {
	store := NewStore("initial", Identity[string])
	store.mu.Lock()
	store.revision = ^Revision(0)
	store.mu.Unlock()

	requirePanicWith(t, "snapshot: revision overflow", func() {
		_ = store.Replace("next")
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, ^Revision(0); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateReceivesIsolatedValue(t *testing.T) {
	store := NewStore([]string{"a"}, cloneStrings)
	var captured []string

	store.Update(func(v []string) []string {
		captured = v
		v[0] = "updated"
		return v
	})

	captured[0] = "changed-after-update"

	snap := store.Snapshot()
	if got, want := snap.Value[0], "updated"; got != want {
		t.Fatalf("stored value = %q, want %q", got, want)
	}
}

func TestStoreUpdateReturnsIsolatedValue(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)

	snap := store.Update(func(v []string) []string {
		v[0] = "updated"
		return v
	})
	snap.Value[0] = "changed-through-result"

	loaded := store.Snapshot()
	if got, want := loaded.Value[0], "updated"; got != want {
		t.Fatalf("stored value = %q, want %q", got, want)
	}
}

func TestStoreUpdateIncrementsOnce(t *testing.T) {
	store := NewStore(1, Identity[int])

	snap := store.Update(func(v int) int {
		return v + 1
	})

	if got, want := snap.Revision, Revision(2); got != want {
		t.Fatalf("Update revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, 2; got != want {
		t.Fatalf("Update value = %d, want %d", got, want)
	}
}

func TestStoreUpdatePanicsOnNilFunction(t *testing.T) {
	store := NewStore("value", Identity[string])

	requirePanicWith(t, "snapshot: nil update function", func() {
		_ = store.Update(nil)
	})
}

func TestStoreUpdatePanicLeavesValueUnchanged(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)

	requirePanicWith(t, "boom", func() {
		_ = store.Update(func(v []string) []string {
			v[0] = "mutated-working-copy"
			panic("boom")
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, Revision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value[0], "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdatePanicsOnRevisionOverflowWithoutCommit(t *testing.T) {
	store := NewStore("initial", Identity[string])
	store.mu.Lock()
	store.revision = ^Revision(0)
	store.mu.Unlock()

	requirePanicWith(t, "snapshot: revision overflow", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, ^Revision(0); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestNewStorePanicsOnNilClone(t *testing.T) {
	requirePanicWith(t, "snapshot: nil clone function", func() {
		_ = NewStore("value", nil)
	})
}

func TestNewStorePanicsOnNilOption(t *testing.T) {
	requirePanicWith(t, "snapshot: nil option", func() {
		_ = NewStore("value", Identity[string], nil)
	})
}

func cloneStrings(v []string) []string {
	return slices.Clone(v)
}
