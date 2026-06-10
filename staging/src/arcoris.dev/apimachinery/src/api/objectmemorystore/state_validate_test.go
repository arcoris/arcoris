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

package objectmemorystore

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestPrepareInputStateRejectsInvalidState(t *testing.T) {
	state := testState("invalid")
	state.Revision = 1

	_, err := prepareInputState(state)
	requireErrorIs(t, err, objectstore.ErrInvalidRevision)
}

func TestPrepareInputStateReturnsDetachedState(t *testing.T) {
	state := testState("prepared")

	prepared, err := prepareInputState(state)
	requireNoError(t, err)

	state.Ownership.Desired.Entries[0].Owner = fieldownership.MustOwner("mutated")
	if prepared.Ownership.Desired.Entries[0].Owner != fieldownership.MustOwner("manager") {
		t.Fatalf("prepareInputState retained caller mutation")
	}
}
