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

	"arcoris.dev/apimachinery/api/objectownership"
)

func TestPrepareInputStateNormalizesOwnership(t *testing.T) {
	state := validState()
	state.Ownership = rawUnsortedOwnership()

	prepared, err := PrepareInputState(state)
	requireNoError(t, err)

	if got, want := prepared.Ownership.Desired.Entries[0].Owner.String(), "a"; got != want {
		t.Fatalf("first owner = %q; want %q", got, want)
	}
	if len(prepared.Ownership.Desired.Entries[0].Fields) != 1 {
		t.Fatalf("fields were not deduplicated: %#v", prepared.Ownership.Desired.Entries[0].Fields)
	}
	requireNoError(t, objectownership.ValidateNormalized(prepared.Ownership))
}

func TestPrepareInputStatePreservesZeroRevision(t *testing.T) {
	prepared, err := PrepareInputState(validState())
	requireNoError(t, err)

	if !prepared.Revision.IsZero() {
		t.Fatalf("revision = %v; want zero", prepared.Revision)
	}
}

func TestPrepareInputStateDoesNotMutateInput(t *testing.T) {
	state := validState()
	state.Ownership = rawUnsortedOwnership()

	_, err := PrepareInputState(state)
	requireNoError(t, err)

	if got := len(state.Ownership.Desired.Entries[1].Fields); got != 2 {
		t.Fatalf("input fields were mutated, len = %d; want 2", got)
	}
}

func TestPrepareInputStateRejectsInvalidState(t *testing.T) {
	_, err := PrepareInputState(State{})
	requireErrorIs(t, err, ErrInvalidState)
}

func TestAssignRevisionDetachesState(t *testing.T) {
	state := validState()
	state.Ownership = ownershipWithEntry()

	committed := AssignRevision(state, 9)
	state.Ownership.Desired.Entries[0].Fields[0] = "$.mutated"

	if committed.Revision != 9 {
		t.Fatalf("revision = %v; want 9", committed.Revision)
	}
	if got := committed.Ownership.Desired.Entries[0].Fields[0]; got != objectownership.Path("$.spec") {
		t.Fatalf("committed field = %q; want $.spec", got)
	}
}
