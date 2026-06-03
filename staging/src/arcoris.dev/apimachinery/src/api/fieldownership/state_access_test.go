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

func TestStateEntriesReturnsDetachedCopy(t *testing.T) {
	state := baseState()
	entries := state.Entries()

	entries[0] = entry("other", metadataPath())

	requireOwners(t, state.Owners(), "autoscaler", "health-controller", "user-cli")
}

func TestStateOwnersSorted(t *testing.T) {
	requireOwners(t, baseState().Owners(), "autoscaler", "health-controller", "user-cli")
}

func TestStateFieldsReturnsUnion(t *testing.T) {
	requireSet(
		t,
		baseState().Fields(),
		`$.conditions[{"type":"Ready"}].status`,
		"$.spec.image",
		"$.spec.replicas",
	)
}

func TestStateFieldsForOwner(t *testing.T) {
	requireSet(t, baseState().FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
}

func TestStateFieldsForUnknownOwnerIsEmpty(t *testing.T) {
	requireEqual(t, baseState().FieldsFor("missing").IsEmpty(), true)
}

func TestStateFieldsForInvalidOwnerIsEmpty(t *testing.T) {
	requireEqual(t, baseState().FieldsFor("").IsEmpty(), true)
}
