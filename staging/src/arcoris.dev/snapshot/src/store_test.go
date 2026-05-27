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
	"maps"
	"slices"
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
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

func TestStoreZeroValuePanicsOnSnapshot(t *testing.T) {
	var store Store[string]

	panicassert.Require(t, func() {
		_ = store.Snapshot()
	})
}

func TestStoreZeroValuePanicsOnStamped(t *testing.T) {
	var store Store[string]

	panicassert.Require(t, func() {
		_ = store.Stamped()
	})
}

func TestStoreZeroValuePanicsOnReplace(t *testing.T) {
	var store Store[string]

	panicassert.Require(t, func() {
		_ = store.Replace("next")
	})
}

func TestStoreZeroValuePanicsOnUpdate(t *testing.T) {
	var store Store[string]

	panicassert.Require(t, func() {
		_ = store.Update(func(v string) string {
			return v
		})
	})
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

	panicassert.RequireMessage(t, "snapshot: revision overflow", func() {
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

func TestStoreReplaceClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 3,
		clone: Identity[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Replace("next")
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, Revision(1); got != want {
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

func TestStoreClonesInitialMutableReadModel(t *testing.T) {
	initial := mutableReadModelValue("initial-name", "initial-attr", "initial-tag")
	want := cloneMutableReadModel(initial)
	store := NewStore(initial, cloneMutableReadModel)

	mutateMutableReadModel(&initial, "changed-name", "changed-attr", "changed-tag")

	snap := store.Snapshot()
	assertMutableReadModel(t, snap.Value, want)
}

func TestStoreSnapshotClonesMutableReadModel(t *testing.T) {
	want := mutableReadModelValue("initial-name", "initial-attr", "initial-tag")
	store := NewStore(want, cloneMutableReadModel)

	snap := store.Snapshot()
	mutateMutableReadModel(&snap.Value, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
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

func TestStoreReplaceReturnsIsolatedMutableReadModel(t *testing.T) {
	store := NewStore(mutableReadModelValue("initial-name", "initial-attr", "initial-tag"), cloneMutableReadModel)
	want := mutableReadModelValue("next-name", "next-attr", "next-tag")

	snap := store.Replace(want)
	mutateMutableReadModel(&snap.Value, "changed-name", "changed-attr", "changed-tag")

	loaded := store.Snapshot()
	assertMutableReadModel(t, loaded.Value, want)
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

func TestStoreBadCloneCanLeakMutableState(t *testing.T) {
	store := NewStore([]string{"a"}, Identity[[]string])

	snap := store.Snapshot()
	snap.Value[0] = "changed"

	// This documents caller-owned misuse: Identity is not a valid CloneFunc for
	// mutable values that require isolation.
	next := store.Snapshot()
	if got, want := next.Value[0], "changed"; got != want {
		t.Fatalf("bad clone test expected aliasing to be visible: got %q, want %q", got, want)
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

func TestStoreUpdateReturnedClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 4,
		clone: Identity[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, Revision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdateStoredClonePanicLeavesValueUnchanged(t *testing.T) {
	cloner := &clonePanicAfter[string]{
		after: 3,
		clone: Identity[string],
	}
	store := NewStore("initial", cloner.Clone)

	panicassert.RequireMessage(t, "clone failed", func() {
		_ = store.Update(func(string) string {
			return "next"
		})
	})

	snap := store.Snapshot()
	if got, want := snap.Revision, Revision(1); got != want {
		t.Fatalf("revision = %d, want %d", got, want)
	}
	if got, want := snap.Value, "initial"; got != want {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestStoreUpdatePanicsOnNilFunction(t *testing.T) {
	store := NewStore("value", Identity[string])

	panicassert.RequireMessage(t, "snapshot: nil update function", func() {
		_ = store.Update(nil)
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

	panicassert.RequireMessage(t, "snapshot: revision overflow", func() {
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
	panicassert.RequireMessage(t, "snapshot: nil clone function", func() {
		_ = NewStore("value", nil)
	})
}

func TestNewStorePanicsOnNilOption(t *testing.T) {
	panicassert.RequireMessage(t, "snapshot: nil option", func() {
		_ = NewStore("value", Identity[string], nil)
	})
}

type mutableReadModel struct {
	Names  []string
	Attrs  map[string]string
	Nested *nestedReadModel
}

type nestedReadModel struct {
	Tags []string
}

func mutableReadModelValue(name, attr, tag string) mutableReadModel {
	return mutableReadModel{
		Names: []string{name},
		Attrs: map[string]string{
			"role": attr,
		},
		Nested: &nestedReadModel{
			Tags: []string{tag},
		},
	}
}

func cloneMutableReadModel(v mutableReadModel) mutableReadModel {
	out := mutableReadModel{
		Names: slices.Clone(v.Names),
		Attrs: maps.Clone(v.Attrs),
	}
	if v.Nested != nil {
		out.Nested = &nestedReadModel{
			Tags: slices.Clone(v.Nested.Tags),
		}
	}
	return out
}

func mutateMutableReadModel(v *mutableReadModel, name, attr, tag string) {
	v.Names[0] = name
	v.Attrs["role"] = attr
	v.Nested.Tags[0] = tag
}

func assertMutableReadModel(t *testing.T, got, want mutableReadModel) {
	t.Helper()

	if !slices.Equal(got.Names, want.Names) {
		t.Fatalf("Names = %#v, want %#v", got.Names, want.Names)
	}
	if !maps.Equal(got.Attrs, want.Attrs) {
		t.Fatalf("Attrs = %#v, want %#v", got.Attrs, want.Attrs)
	}
	if got.Nested == nil || want.Nested == nil {
		if got.Nested != want.Nested {
			t.Fatalf("Nested = %#v, want %#v", got.Nested, want.Nested)
		}
		return
	}
	if !slices.Equal(got.Nested.Tags, want.Nested.Tags) {
		t.Fatalf("Nested.Tags = %#v, want %#v", got.Nested.Tags, want.Nested.Tags)
	}
}

type clonePanicAfter[T any] struct {
	after int
	calls int
	clone func(T) T
}

func (c *clonePanicAfter[T]) Clone(v T) T {
	c.calls++
	if c.calls == c.after {
		panic("clone failed")
	}
	return c.clone(v)
}

func cloneStrings(v []string) []string {
	return slices.Clone(v)
}
