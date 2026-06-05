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

package atomicx

import "testing"

// TestUint64CounterSampleZeroValue verifies sample value objects start at zero.
func TestUint64CounterSampleZeroValue(t *testing.T) {
	t.Parallel()

	var sample Uint64CounterSample

	if sample.Value != 0 {
		t.Fatalf("zero-value Uint64CounterSample.Value = %d, want 0", sample.Value)
	}
}

// TestUint64CounterSampleCapturesPointInTimeValue verifies Sample reads the current counter value.
func TestUint64CounterSampleCapturesPointInTimeValue(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(64)
	first := counter.Sample()

	counter.Add(32)
	second := counter.Sample()

	if first.Value != 64 {
		t.Fatalf("first Uint64CounterSample.Value = %d, want 64", first.Value)
	}

	if second.Value != 96 {
		t.Fatalf("second Uint64CounterSample.Value = %d, want 96", second.Value)
	}
}

// TestUint64CounterSampleIsIndependentFromLaterCounterUpdates verifies samples are immutable.
func TestUint64CounterSampleIsIndependentFromLaterCounterUpdates(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(10)
	sample := counter.Sample()

	counter.Add(90)

	if sample.Value != 10 {
		t.Fatalf("sample Value changed after counter update: got %d, want 10", sample.Value)
	}

	if got := counter.Load(); got != 100 {
		t.Fatalf("Uint64Counter.Load() = %d, want 100", got)
	}
}

// TestUint64CounterSampleDeltaSince verifies samples delegate ordinary deltas correctly.
func TestUint64CounterSampleDeltaSince(t *testing.T) {
	t.Parallel()

	previous := Uint64CounterSample{Value: 100}
	current := Uint64CounterSample{Value: 175}

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

// TestUint64CounterSampleDeltaSinceWrapped verifies sample deltas preserve single-wrap semantics.
func TestUint64CounterSampleDeltaSinceWrapped(t *testing.T) {
	t.Parallel()

	previous := Uint64CounterSample{Value: testMaxUint64 - 2}
	current := Uint64CounterSample{Value: 4}

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
