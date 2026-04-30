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

// Snapshot returns an immutable point-in-time sample of the counter.
//
// A snapshot is a plain value object. It does not retain a pointer to the
// counter, so later counter updates cannot change the sample.
//
// Snapshot observes exactly one atomic value. It is safe to call concurrently
// with Load, Add, and Inc.
//
// The returned value is intentionally copyable. Snapshot and delta values do not
// contain noCopy because they do not own atomic state and do not participate in
// synchronization.
//
// Uint32Counter snapshots are intentionally separated from the mutable
// Uint32Counter implementation. This file owns the copyable sampling contract;
// counter_uint32.go owns the mutable counter state.
func (c *Uint32Counter) Snapshot() Uint32CounterSnapshot {
	return Uint32CounterSnapshot{
		Value: c.value.Load(),
	}
}

// Uint32CounterSnapshot is an immutable point-in-time sample of a
// Uint32Counter.
//
// Uint32CounterSnapshot belongs to the sampling layer, not to the mutable
// counter state layer. It records one observed lifetime counter value and is
// safe to copy, store, compare, and pass by value.
//
// Snapshots are used to compute activity over a window without resetting or
// mutating the source lifetime counter:
//
//	previous := counter.Snapshot()
//	// Workload executes.
//	current := counter.Snapshot()
//	delta := current.DeltaSince(previous)
//
// A snapshot records only one counter value. A group of snapshots taken from
// several counters is not globally atomic unless the caller provides additional
// synchronization around the whole sampling operation.
//
// Uint32CounterSnapshot does not store wall-clock time or monotonic time. Time
// windows, rates, and sampling cadence belong to the caller or to a higher-level
// metrics/control package, not to atomicx.
//
// Uint32 counters wrap much sooner than uint64 counters. Callers that compute
// deltas from Uint32CounterSnapshot values must sample frequently enough for the
// expected event rate.
type Uint32CounterSnapshot struct {
	// Value is the observed lifetime counter value at sampling time.
	Value uint32
}

// DeltaSince returns the monotonic counter delta from previous to s.
//
// The receiver must be the newer sample. The previous argument must be the older
// sample. The method cannot verify sample ordering because snapshots do not
// carry time metadata.
//
// The returned delta is wrap-aware for one uint32 wrap:
//
//   - if s.Value >= previous.Value, the delta value is s.Value - previous.Value;
//   - if s.Value < previous.Value, the counter is interpreted as having wrapped
//     once and the delta value is computed with uint32 modulo arithmetic.
//
// Multiple wraps between two samples cannot be detected from two uint32 values
// alone. Callers that rely on accurate activity windows must sample frequently
// enough to make multiple wraps impossible in practice.
//
// Prefer Uint64Counter for long-lived or high-rate counters unless a 32-bit
// counter is part of a deliberate bounded state model or external protocol
// boundary.
func (s Uint32CounterSnapshot) DeltaSince(previous Uint32CounterSnapshot) Uint32CounterDelta {
	return NewUint32CounterDelta(previous.Value, s.Value)
}
