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

// Sample returns an immutable point-in-time observation of the counter.
//
// A sample is a plain value object. It does not retain a pointer to the counter,
// so later counter updates cannot change the sampled value. Sample observes
// exactly one atomic value and is safe to call concurrently with Load, Add, and
// Inc.
//
// Counter samples and deltas are intentionally copyable. They do not contain
// noCopy markers because they do not own mutable atomic state and do not
// participate in synchronization.
func (c *Uint64Counter) Sample() Uint64CounterSample {
	return Uint64CounterSample{
		Value: c.value.Load(),
	}
}

// Uint64CounterSample is an immutable point-in-time observation of a
// Uint64Counter.
//
// Uint64CounterSample belongs to the sampling layer, not to the mutable counter
// state layer. It records one observed lifetime counter value and is safe to
// copy, store, compare, and pass by value.
//
// Samples are used to compute activity over a window without resetting or
// mutating the source lifetime counter:
//
//	previous := counter.Sample()
//	// Workload executes.
//	current := counter.Sample()
//	delta := current.DeltaSince(previous)
//
// A sample records only one counter value. A group of samples taken from several
// counters is not globally atomic unless the caller provides additional
// synchronization around the whole sampling operation.
//
// Uint64CounterSample does not store wall-clock time or monotonic time. Time
// windows, rates, and sampling cadence belong to the caller or to a higher-level
// metrics/control package, not to atomicx.
type Uint64CounterSample struct {
	// Value is the observed lifetime counter value at sampling time.
	Value uint64
}

// DeltaSince returns the monotonic counter delta from previous to s.
//
// The receiver must be the newer sample. The previous argument must be the older
// sample. The method cannot verify sample ordering because samples do not carry
// time metadata.
//
// The returned delta is wrap-aware for one uint64 wrap:
//
//   - if s.Value >= previous.Value, the delta value is s.Value - previous.Value;
//   - if s.Value < previous.Value, the counter is interpreted as having wrapped
//     once and the delta value is computed with uint64 modulo arithmetic.
//
// Multiple wraps between two samples cannot be detected from two uint64 values
// alone. Callers that rely on accurate activity windows must sample frequently
// enough to make multiple wraps impossible in practice.
func (s Uint64CounterSample) DeltaSince(previous Uint64CounterSample) Uint64CounterDelta {
	return NewUint64CounterDelta(previous.Value, s.Value)
}
