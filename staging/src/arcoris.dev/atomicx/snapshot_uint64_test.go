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

package atomicx

import "testing"

// TestUint64CounterSnapshotZeroValue verifies snapshot value objects start at zero.
func TestUint64CounterSnapshotZeroValue(t *testing.T) {
	t.Parallel()

	var snapshot Uint64CounterSnapshot

	if snapshot.Value != 0 {
		t.Fatalf("zero-value Uint64CounterSnapshot.Value = %d, want 0", snapshot.Value)
	}
}

// TestUint64CounterSnapshotCapturesPointInTimeValue verifies Snapshot reads the current counter value.
func TestUint64CounterSnapshotCapturesPointInTimeValue(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(64)
	first := counter.Snapshot()

	counter.Add(32)
	second := counter.Snapshot()

	if first.Value != 64 {
		t.Fatalf("first Uint64CounterSnapshot.Value = %d, want 64", first.Value)
	}

	if second.Value != 96 {
		t.Fatalf("second Uint64CounterSnapshot.Value = %d, want 96", second.Value)
	}
}

// TestUint64CounterSnapshotIsIndependentFromLaterCounterUpdates verifies snapshots are immutable samples.
func TestUint64CounterSnapshotIsIndependentFromLaterCounterUpdates(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(10)
	snapshot := counter.Snapshot()

	counter.Add(90)

	if snapshot.Value != 10 {
		t.Fatalf("snapshot Value changed after counter update: got %d, want 10", snapshot.Value)
	}

	if got := counter.Load(); got != 100 {
		t.Fatalf("Uint64Counter.Load() = %d, want 100", got)
	}
}

// TestUint64CounterSnapshotDeltaSince verifies snapshots delegate ordinary deltas correctly.
func TestUint64CounterSnapshotDeltaSince(t *testing.T) {
	t.Parallel()

	previous := Uint64CounterSnapshot{Value: 100}
	current := Uint64CounterSnapshot{Value: 175}

	delta := current.DeltaSince(previous)

	if delta.Previous != previous.Value {
		t.Fatalf("delta.Previous = %d, want %d", delta.Previous, previous.Value)
	}

	if delta.Current != current.Value {
		t.Fatalf("delta.Current = %d, want %d", delta.Current, current.Value)
	}

	if delta.Value != 75 {
		t.Fatalf("delta.Value = %d, want 75", delta.Value)
	}

	if delta.Wrapped {
		t.Fatal("delta.Wrapped = true, want false")
	}
}

// TestUint64CounterSnapshotDeltaSinceWrapped verifies snapshot deltas preserve single-wrap semantics.
func TestUint64CounterSnapshotDeltaSinceWrapped(t *testing.T) {
	t.Parallel()

	previous := Uint64CounterSnapshot{Value: testMaxUint64 - 2}
	current := Uint64CounterSnapshot{Value: 4}

	delta := current.DeltaSince(previous)

	if delta.Previous != previous.Value {
		t.Fatalf("delta.Previous = %d, want %d", delta.Previous, previous.Value)
	}

	if delta.Current != current.Value {
		t.Fatalf("delta.Current = %d, want %d", delta.Current, current.Value)
	}

	if delta.Value != 7 {
		t.Fatalf("delta.Value = %d, want 7", delta.Value)
	}

	if !delta.Wrapped {
		t.Fatal("delta.Wrapped = false, want true")
	}
}
