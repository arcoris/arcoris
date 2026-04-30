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

// TestUint32CounterDeltaZeroValue verifies zero deltas are empty value objects.
func TestUint32CounterDeltaZeroValue(t *testing.T) {
	t.Parallel()

	var delta Uint32CounterDelta

	if delta.Previous != 0 {
		t.Fatalf("zero-value Uint32CounterDelta.Previous = %d, want 0", delta.Previous)
	}

	if delta.Current != 0 {
		t.Fatalf("zero-value Uint32CounterDelta.Current = %d, want 0", delta.Current)
	}

	if delta.Value != 0 {
		t.Fatalf("zero-value Uint32CounterDelta.Value = %d, want 0", delta.Value)
	}

	if delta.Wrapped {
		t.Fatal("zero-value Uint32CounterDelta.Wrapped = true, want false")
	}

	if !delta.IsZero() {
		t.Fatal("zero-value Uint32CounterDelta.IsZero() = false, want true")
	}
}

// TestNewUint32CounterDeltaIncreasing verifies ordinary increasing samples.
func TestNewUint32CounterDeltaIncreasing(t *testing.T) {
	t.Parallel()

	delta := NewUint32CounterDelta(100, 175)

	if delta.Previous != 100 {
		t.Fatalf("delta.Previous = %d, want 100", delta.Previous)
	}

	if delta.Current != 175 {
		t.Fatalf("delta.Current = %d, want 175", delta.Current)
	}

	if delta.Value != 75 {
		t.Fatalf("delta.Value = %d, want 75", delta.Value)
	}

	if delta.Wrapped {
		t.Fatal("delta.Wrapped = true, want false")
	}

	if delta.IsZero() {
		t.Fatal("delta.IsZero() = true, want false")
	}
}

// TestNewUint32CounterDeltaSameValue verifies equal samples produce zero activity.
func TestNewUint32CounterDeltaSameValue(t *testing.T) {
	t.Parallel()

	delta := NewUint32CounterDelta(42, 42)

	if delta.Previous != 42 {
		t.Fatalf("delta.Previous = %d, want 42", delta.Previous)
	}

	if delta.Current != 42 {
		t.Fatalf("delta.Current = %d, want 42", delta.Current)
	}

	if delta.Value != 0 {
		t.Fatalf("delta.Value = %d, want 0", delta.Value)
	}

	if delta.Wrapped {
		t.Fatal("delta.Wrapped = true, want false")
	}

	if !delta.IsZero() {
		t.Fatal("delta.IsZero() = false, want true")
	}
}

// TestNewUint32CounterDeltaWrapped verifies one observed wrap uses modulo arithmetic.
func TestNewUint32CounterDeltaWrapped(t *testing.T) {
	t.Parallel()

	previous := testMaxUint32 - 2
	current := uint32(4)

	delta := NewUint32CounterDelta(previous, current)

	if delta.Previous != previous {
		t.Fatalf("delta.Previous = %d, want %d", delta.Previous, previous)
	}

	if delta.Current != current {
		t.Fatalf("delta.Current = %d, want %d", delta.Current, current)
	}

	if delta.Value != 7 {
		t.Fatalf("delta.Value = %d, want 7", delta.Value)
	}

	if !delta.Wrapped {
		t.Fatal("delta.Wrapped = false, want true")
	}

	if delta.IsZero() {
		t.Fatal("delta.IsZero() = true, want false")
	}
}

// TestNewUint32CounterDeltaFullRangeWrap verifies max-to-zero is one event under wrap semantics.
func TestNewUint32CounterDeltaFullRangeWrap(t *testing.T) {
	t.Parallel()

	delta := NewUint32CounterDelta(testMaxUint32, 0)

	if delta.Value != 1 {
		t.Fatalf("delta.Value = %d, want 1", delta.Value)
	}

	if !delta.Wrapped {
		t.Fatal("delta.Wrapped = false, want true")
	}
}
