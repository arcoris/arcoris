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

import "testing"

func TestLiveRecordAssignsRevision(t *testing.T) {
	rec := liveRecord(testState("live"), 9)

	if rec.deleted {
		t.Fatalf("live record marked deleted")
	}
	if rec.state.Revision != 9 {
		t.Fatalf("Revision = %v; want 9", rec.state.Revision)
	}
	requireDesiredString(t, rec.state, "live")
}

func TestTombstoneRecordPreservesPreviousStateAndDeleteRevision(t *testing.T) {
	previous := liveRecord(testState("deleted"), 4).state
	rec := tombstoneRecord(previous, 5)

	if !rec.deleted {
		t.Fatalf("tombstone record not marked deleted")
	}
	if rec.deleteRevision != 5 {
		t.Fatalf("deleteRevision = %v; want 5", rec.deleteRevision)
	}
	if rec.state.Revision != 4 {
		t.Fatalf("state revision = %v; want previous revision 4", rec.state.Revision)
	}
}
