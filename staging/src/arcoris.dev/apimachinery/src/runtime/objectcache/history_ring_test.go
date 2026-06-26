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

func TestVersionRingEvictsOldestPerObject(t *testing.T) {
	ring := newVersionRing(2)
	key := testKey("system", 1)

	ring.append(liveVersion(1, testState(key, 1, "one")))
	ring.append(liveVersion(3, testState(key, 3, "three")))
	ring.append(liveVersion(7, testState(key, 7, "seven")))
	if ring.count() != 2 || ring.capacity() != 2 {
		t.Fatalf("count/capacity = %d/%d; want 2/2", ring.count(), ring.capacity())
	}

	var got []objectstore.Revision
	ring.oldestToNewest(func(version objectVersion) bool {
		got = append(got, version.Revision)
		return true
	})

	if len(got) != 2 || got[0] != 3 || got[1] != 7 {
		t.Fatalf("ring revisions = %v; want [3 7]", got)
	}
}

func TestVersionRingReturnsDetachedStates(t *testing.T) {
	ring := newVersionRing(1)
	key := testKey("system", 1)
	state := testState(key, 1, "stored")
	ring.append(liveVersion(1, state))

	ring.newestToOldest(func(version objectVersion) bool {
		mutateState(&version.State, "mutated")
		return false
	})

	ring.newestToOldest(func(version objectVersion) bool {
		if desired := desiredString(t, version.State); desired != "stored" {
			t.Fatalf("desired = %q; want stored", desired)
		}
		return false
	})
}
