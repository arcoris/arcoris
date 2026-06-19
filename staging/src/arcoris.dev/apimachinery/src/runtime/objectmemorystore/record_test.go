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

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestRecordCarriesLiveOrTombstoneState(t *testing.T) {
	live := record{state: objectstore.AssignRevision(testState("live"), 1)}
	if live.deleted || live.deleteRevision != 0 {
		t.Fatalf("live record = %#v; want non-deleted", live)
	}

	tombstone := record{
		state:          objectstore.AssignRevision(testState("deleted"), 2),
		deleteRevision: 3,
		deleted:        true,
	}
	if !tombstone.deleted || tombstone.deleteRevision != 3 {
		t.Fatalf("tombstone record = %#v; want deleted at revision 3", tombstone)
	}
}
