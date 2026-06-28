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

func TestNewStoreClonesInitialMutableReadModel(t *testing.T) {
	initial := mutableReadModelValue("initial-name", "initial-attr", "initial-tag")
	want := cloneMutableReadModel(initial)
	store := NewStore(initial, cloneMutableReadModel)

	mutateMutableReadModel(&initial, "changed-name", "changed-attr", "changed-tag")

	snap := store.Snapshot()
	assertMutableReadModel(t, snap.Value, want)
}

func TestNewStoreInitialRevisionIsOne(t *testing.T) {
	store := NewStore("value", IdentityClone[string])

	if got, want := store.Revision(), LocalRevision(1); got != want {
		t.Fatalf("LocalRevision() = %d, want %d", got, want)
	}
}

func TestNewStorePanicsOnNilClone(t *testing.T) {
	panicassert.RequireMessage(t, "snapshot: nil clone function", func() {
		_ = NewStore("value", nil)
	})
}

func TestNewStorePanicsOnNilOption(t *testing.T) {
	panicassert.RequireMessage(t, "snapshot: nil option", func() {
		_ = NewStore("value", IdentityClone[string], nil)
	})
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
