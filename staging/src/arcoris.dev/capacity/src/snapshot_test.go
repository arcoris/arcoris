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

package capacity_test

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestNewSnapshotDerivesAvailabilityAndDebt(t *testing.T) {
	t.Parallel()

	limits := vector(t, entry("memory_bytes", 8), entry("worker_slots", 4))
	reserved := vector(t, entry("memory_bytes", 10), entry("queue_slots", 2), entry("worker_slots", 1))

	snap := capacity.NewSnapshot(limits, reserved)
	if !snap.IsValid() || !snap.HasDebt() {
		t.Fatalf("snapshot valid=%v debt=%v, want true/true", snap.IsValid(), snap.HasDebt())
	}
	requireVector(t, snap.Available, entry("worker_slots", 3))
	requireVector(t, snap.Debt, entry("memory_bytes", 2), entry("queue_slots", 2))
}

func TestSnapshotVectorFieldsAreCopySafe(t *testing.T) {
	t.Parallel()

	snap := capacity.NewSnapshot(vector(t, entry("worker_slots", 2)), vector(t, entry("worker_slots", 1)))
	entries := snap.Available.Entries()
	entries[0] = entry("worker_slots", 99)

	if got := snap.AvailableFor(capacity.MustResource("worker_slots")); got != 1 {
		t.Fatalf("AvailableFor() = %d, want 1", got)
	}
}
