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

package objectstore

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestChangeIsZero(t *testing.T) {
	if !(Change{}).IsZero() {
		t.Fatalf("zero change is not zero")
	}
	if (Change{Kind: ChangeCreated}).IsZero() {
		t.Fatalf("non-zero change kind reported zero")
	}
}

func TestChangeIsValid(t *testing.T) {
	valid := Change{
		Kind:     ChangeCreated,
		Key:      validKey(),
		Revision: 1,
		After:    committedStateAt(1, "after"),
	}
	if !valid.IsValid() {
		t.Fatalf("valid change reported invalid")
	}
	if (Change{}).IsValid() {
		t.Fatalf("zero change reported valid")
	}
}

func TestChangeValidateMatchesValidateChange(t *testing.T) {
	change := Change{
		Kind:     ChangeUpdated,
		Key:      validKey(),
		Revision: 2,
		Before:   committedStateAt(1, "before"),
		After:    committedStateAt(2, "after"),
	}

	requireNoError(t, change.Validate())
	requireNoError(t, ValidateChange(change))
}

func TestChangeCloneDetachesStates(t *testing.T) {
	change := Change{
		Kind:     ChangeUpdated,
		Key:      validKey(),
		Revision: 2,
		Before:   committedStateAt(1, "before"),
		After:    committedStateAt(2, "after"),
	}

	cloned := change.Clone()
	cloned.Before.Object.Desired = value.StringValue("before-mutated")
	cloned.After.Object.Desired = value.StringValue("after-mutated")

	if cloned.Kind != change.Kind || !cloned.Key.Equal(change.Key) || cloned.Revision != change.Revision {
		t.Fatalf("clone identity = %#v; want %#v", cloned, change)
	}
	requireChangeDesired(t, change.Before, "before")
	requireChangeDesired(t, change.After, "after")
}

func TestChangeCloneZeroRemainsZero(t *testing.T) {
	if clone := (Change{}).Clone(); !clone.IsZero() {
		t.Fatalf("zero clone = %#v; want zero", clone)
	}
}

// committedStateAt constructs a valid committed state with a chosen revision.
func committedStateAt(revision Revision, desired string) State {
	state := validState()
	state.Object.Desired = value.StringValue(desired)
	state.Revision = revision

	return state
}

// requireChangeDesired checks the Desired string payload in a change state.
func requireChangeDesired(t *testing.T, state State, want string) {
	t.Helper()

	got, ok := state.Object.Desired.AsString()
	if !ok || got != want {
		t.Fatalf("desired = %q, %v; want %q, true", got, ok, want)
	}
}
