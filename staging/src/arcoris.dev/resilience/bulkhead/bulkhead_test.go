/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package bulkhead

import "testing"

func TestBulkheadSnapshotReturnsValidCapacityState(t *testing.T) {
	t.Parallel()

	b := New(2)
	snap := b.Snapshot()
	if !snap.Value.IsValid() {
		t.Fatalf("snapshot is invalid: %+v", snap.Value)
	}
	requireSnapshotValue(t, snap, 2, 0, 2, 0)
}

func TestBulkheadRevisionMatchesSnapshotRevision(t *testing.T) {
	t.Parallel()

	b := New(1)
	if got, want := b.Revision(), b.Snapshot().Revision; got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}
