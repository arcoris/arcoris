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
	state.Ownership = objectownership.Document{
		Version: objectownership.VersionV1,
		Desired: objectownership.Surface{
			Entries: []objectownership.Entry{
				{Owner: "z", Fields: []objectownership.Path{"$.z"}},
				{Owner: "a", Fields: []objectownership.Path{"$.a"}},
			},
		},
	}

	prepared, err := PrepareInputState(state)
	requireNoError(t, err)

	if got, want := prepared.Ownership.Desired.Entries[0].Owner.String(), "a"; got != want {
		t.Fatalf("first owner = %q; want %q", got, want)
	}
}

func TestPrepareInputStateRejectsInvalidState(t *testing.T) {
	_, err := PrepareInputState(State{})
	requireErrorIs(t, err, ErrInvalidState)
}
