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

package objectworkqueue

import "testing"

func TestEntryStateValuesAreNonZero(t *testing.T) {
	if stateQueued == 0 || stateProcessing == 0 {
		t.Fatalf("entry state zero value must remain invalid")
	}
}

func TestEntryTracksItemStateDirtyAndLinks(t *testing.T) {
	item := testItem(1)
	entry := &entry{item: item, state: stateProcessing, dirty: true}

	if !entry.item.Key.Equal(item.Key) {
		t.Fatalf("entry item = %s; want %s", entry.item.Key, item.Key)
	}
	if entry.state != stateProcessing {
		t.Fatalf("entry state = %d; want processing", entry.state)
	}
	if !entry.dirty {
		t.Fatalf("entry dirty = false; want true")
	}
	if entry.prev != nil || entry.next != nil {
		t.Fatalf("entry links = %p,%p; want nil", entry.prev, entry.next)
	}
}
