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
	"errors"
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

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

func TestStoreUpdateReceivesIsolatedMutableReadModel(t *testing.T) {
	store := NewStore(mutableReadModelValue("initial-name", "initial-attr", "initial-tag"), cloneMutableReadModel)
	want := mutableReadModelValue("updated-name", "updated-attr", "updated-tag")
	var captured mutableReadModel

	_ = store.Update(func(v mutableReadModel) mutableReadModel {
		captured = v
		mutateMutableReadModel(&v, "updated-name", "updated-attr", "updated-tag")
		return v
	})
	mutateMutableReadModel(&captured, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
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

func TestStoreUpdateReturnsIsolatedMutableReadModel(t *testing.T) {
	store := NewStore(mutableReadModelValue("initial-name", "initial-attr", "initial-tag"), cloneMutableReadModel)
	want := mutableReadModelValue("updated-name", "updated-attr", "updated-tag")

	snap := store.Update(func(v mutableReadModel) mutableReadModel {
		mutateMutableReadModel(&v, "updated-name", "updated-attr", "updated-tag")
		return v
	})
	mutateMutableReadModel(&snap.Value, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
}

func TestStoreUpdateIncrementsOnce(t *testing.T) {
	store := NewStore(1, IdentityClone[int])

	snap := store.Update(func(v int) int {
		return v + 1
	})

	if got, want := snap.Revision, LocalRevision(2); got != want {
		t.Fatalf("Update revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, 2; got != want {
		t.Fatalf("Update value = %d, want %d", got, want)
	}
}

func TestStoreUpdateReturnedClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 4,
		clone: IdentityClone[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateStoredClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 3,
		clone: IdentityClone[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdatePanicsOnNilFunction(t *testing.T) {
	store := NewStore("value", IdentityClone[string])

	panicassert.RequireMessage(t, "snapshot: nil update function", func() {
		_ = store.Update(nil)
	})
}

func TestStoreUpdateErrCommitsOnNilError(t *testing.T) {
	store := NewStore(1, IdentityClone[int])

	snap, err := store.UpdateErr(func(v int) (int, error) {
		return v + 1, nil
	})
	if err != nil {
		t.Fatalf("UpdateErr() error = %v, want nil", err)
	}

	if got, want := snap.Revision, LocalRevision(2); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, 2; got != want {
		t.Fatalf("value = %d, want %d", got, want)
	}
}

func TestStoreUpdateErrLeavesValueUnchangedOnError(t *testing.T) {
	updateErr := errors.New("update failed")
	store := NewStore([]string{"initial"}, cloneStrings)

	_, err := store.UpdateErr(func(v []string) ([]string, error) {
		v[0] = "changed-working-copy"
		return v, updateErr
	})
	if !errors.Is(err, updateErr) {
		t.Fatalf("UpdateErr() error = %v, want %v", err, updateErr)
	}

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value[0], "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateStampedErrCommitsTimestampOnNilError(t *testing.T) {
	clk := newTestClock()
	clk.set(time.Unix(1, 0))
	store := NewStore("initial", IdentityClone[string], WithClock(clk))

	clk.set(time.Unix(2, 0))
	stamped, err := store.UpdateStampedErr(func(string) (string, error) {
		return "next", nil
	})
	if err != nil {
		t.Fatalf("UpdateStampedErr() error = %v, want nil", err)
	}

	if got, want := stamped.Revision, LocalRevision(2); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if !stamped.Updated.Equal(time.Unix(2, 0)) {
		t.Fatalf("updated = %s, want %s", stamped.Updated, time.Unix(2, 0))
	}
	if got, want := stamped.Value, "next"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateStampedErrLeavesValueUnchangedOnError(t *testing.T) {
	updateErr := errors.New("update failed")
	store := NewStore("initial", IdentityClone[string])

	_, err := store.UpdateStampedErr(func(string) (string, error) {
		return "next", updateErr
	})
	if !errors.Is(err, updateErr) {
		t.Fatalf("UpdateStampedErr() error = %v, want %v", err, updateErr)
	}

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateErrPanicsOnNilFunction(t *testing.T) {
	store := NewStore("value", IdentityClone[string])

	panicassert.RequireMessage(t, "snapshot: nil update function", func() {
		_, _ = store.UpdateErr(nil)
	})
}

func TestStoreUpdatePanicLeavesValueUnchanged(t *testing.T) {
	store := NewStore([]string{"initial"}, cloneStrings)

	panicassert.RequireMessage(t, "boom", func() {
		_ = store.Update(func(v []string) []string {
			v[0] = "mutated-working-copy"
			panic("boom")
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, LocalRevision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value[0], "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdatePanicsOnRevisionOverflowWithoutCommit(t *testing.T) {
	store := NewStore("initial", IdentityClone[string])
	store.mu.Lock()
	store.revision = ^LocalRevision(0)
	store.mu.Unlock()

	panicassert.RequireMessage(t, "snapshot: local revision overflow", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, ^LocalRevision(0); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}
