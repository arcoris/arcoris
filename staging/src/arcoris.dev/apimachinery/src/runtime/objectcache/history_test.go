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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestLiveVersionUsesObservationRevisionAndClonesState(t *testing.T) {
	key := testKey("system", 1)
	state := testState(key, 3, "object-revision-three")

	version := liveVersion(10, state)
	mutateState(&state, "mutated")

	if !version.Live || version.Revision != 10 || version.State.Revision != 3 {
		t.Fatalf("liveVersion() = %#v; want live observation revision 10 with object revision 3", version)
	}
	if desired := desiredString(t, version.State); desired != "object-revision-three" {
		t.Fatalf("desired = %q; want object-revision-three", desired)
	}
}

func TestTombstoneVersionCarriesOnlyObservationRevision(t *testing.T) {
	version := tombstoneVersion(7)

	if version.Live || version.Revision != 7 || !sameState(version.State, objectstore.State{}) {
		t.Fatalf("tombstoneVersion() = %#v; want tombstone at revision 7 with no state", version)
	}
}

func TestObjectVersionCloneDetachesLiveState(t *testing.T) {
	key := testKey("system", 1)
	version := liveVersion(5, testState(key, 5, "stored"))

	clone := version.clone()
	mutateState(&clone.State, "mutated")

	if desired := desiredString(t, version.State); desired != "stored" {
		t.Fatalf("desired = %q; want stored", desired)
	}
}
