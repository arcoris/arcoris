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

import "testing"

func TestHistoricalReadHelpersBuildDetachedResults(t *testing.T) {
	key := testKey("system", 1)
	state := testState(key, 4, "stored")

	absence := knownAbsenceAt(key, 10)
	if absence.Found || absence.Revision != 10 || !absence.Key.Equal(key) {
		t.Fatalf("knownAbsenceAt() = %#v; want absence at revision 10", absence)
	}

	live := liveResultAt(key, 10, state)
	mutateState(&state, "mutated")
	if !live.Found || live.Revision != 10 || desiredString(t, live.State) != "stored" {
		t.Fatalf("liveResultAt() = %#v; want detached live state at revision 10", live)
	}
}
