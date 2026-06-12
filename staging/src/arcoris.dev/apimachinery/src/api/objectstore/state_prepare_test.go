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
	state.Ownership = ownershipWithSurfaces()

	prepared, err := PrepareInputState(state)
	requireNoError(t, err)

	requireNoError(t, objectownership.ValidateNormalized(prepared.Ownership))
	requireOwnershipField(t, prepared.Ownership.Desired(), "manager", "$.spec")
	requireOwnershipField(t, prepared.Ownership.Observed(), "controller", "$.ready")
	requireOwnershipField(t, prepared.Ownership.Metadata().Labels(), "labeler", `$["app"]`)
	requireOwnershipField(t, prepared.Ownership.Metadata().Annotations(), "annotator", `$["scheduler.arcoris.dev/mode"]`)
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
	state.Ownership = ownershipWithSurfaces()

	_, err := PrepareInputState(state)
	requireNoError(t, err)

	requireOwnershipField(t, state.Ownership.Metadata().Annotations(), "annotator", `$["scheduler.arcoris.dev/mode"]`)
}

func TestPrepareInputStateRejectsInvalidState(t *testing.T) {
	_, err := PrepareInputState(State{})
	requireErrorIs(t, err, ErrInvalidState)
}

func TestAssignRevisionDetachesState(t *testing.T) {
	state := validState()
	state.Ownership = ownershipWithEntry()

	committed := AssignRevision(state, 9)

	if committed.Revision != 9 {
		t.Fatalf("revision = %v; want 9", committed.Revision)
	}
	requireOwnershipField(t, committed.Ownership.Desired(), "manager", "$.spec")
}
