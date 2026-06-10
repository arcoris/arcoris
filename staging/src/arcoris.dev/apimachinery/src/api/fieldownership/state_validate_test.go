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

package fieldownership

import "testing"

func TestZeroStateValidateStructure(t *testing.T) {
	requireNoError(t, State{}.ValidateStructure())
}

func TestStateValidateStructureAcceptsNormalizedState(t *testing.T) {
	requireNoError(t, baseState().ValidateStructure())
}

func TestStateValidateStructureRejectsUnsortedEntries(t *testing.T) {
	state := State{entries: []Entry{
		entry("z", imagePath()),
		entry("a", replicasPath()),
	}}

	err := state.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidState)
	requireFieldOwnershipError(t, err, "entries[1].owner", ErrorReasonUnsortedStateEntries)
}

func TestStateValidateStructureRejectsDuplicateOwners(t *testing.T) {
	state := State{entries: []Entry{
		entry("a", imagePath()),
		entry("a", replicasPath()),
	}}

	err := state.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidState)
	requireFieldOwnershipError(t, err, "entries[1].owner", ErrorReasonDuplicateStateOwner)
}

func TestStateValidateStructureRejectsEmptyEntry(t *testing.T) {
	state := State{entries: []Entry{
		emptyEntry("a"),
	}}

	err := state.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidState)
	requireFieldOwnershipError(t, err, "entries[0].fields", ErrorReasonEmptyStateEntry)
}

func TestStateValidateStructureRejectsInvalidEntryOwner(t *testing.T) {
	state := State{entries: []Entry{
		{owner: Owner{text: ""}, fields: set(imagePath())},
	}}

	err := state.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidState)
	requireErrorIs(t, err, ErrInvalidOwner)
	requireFieldOwnershipError(t, err, "entries[0].owner", ErrorReasonInvalidEntryOwner)
}
