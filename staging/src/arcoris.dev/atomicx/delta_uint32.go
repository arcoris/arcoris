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

// Uint32CounterDelta describes the difference between two Uint32Counter samples.
//
// Uint32CounterDelta is a copyable value object. It intentionally does not
// contain noCopy because it does not own atomic state, does not synchronize
// access to anything, and is safe to pass by value.
//
// The delta uses unsigned modulo arithmetic:
//
//   - if Current >= Previous, Value is Current - Previous;
//   - if Current < Previous, the counter is interpreted as having wrapped once,
//     and Value is still computed as Current - Previous using uint32 arithmetic.
//
// Wrapped reports whether Current was lower than Previous and the delta was
// interpreted as a single uint32 wrap.
//
// Uint32 counters wrap much sooner than uint64 counters. This type is therefore
// appropriate only when the 32-bit counter range is deliberate and the caller can
// guarantee a sampling cadence that makes multiple wraps impossible in practice.
// Prefer Uint64CounterDelta for general long-running or high-rate event
// accounting.
type Uint32CounterDelta struct {
	// Previous is the older lifetime counter value.
	Previous uint32

	// Current is the newer lifetime counter value.
	Current uint32

	// Value is the computed delta between Previous and Current.
	//
	// Value is computed with uint32 modulo arithmetic, so it remains correct for
	// one wrap between Previous and Current.
	Value uint32

	// Wrapped reports whether Current was lower than Previous and the delta was
	// interpreted as a single uint32 wrap.
	Wrapped bool
}

// NewUint32CounterDelta computes a wrap-aware delta between two uint32 lifetime
// counter values.
//
// The previous value must be the older sample and current must be the newer
// sample. The function cannot verify temporal ordering because it receives only
// raw counter values.
//
// uint32 subtraction already has the desired modulo behavior for a single wrap,
// so current - previous works for both ordinary and wrapped cases. Wrapped is
// recorded separately because diagnostics, tests, and higher-level reporting may
// need to know that a wrap was observed.
//
// Multiple wraps cannot be detected from previous and current alone. Because
// uint32 wraps relatively quickly under high event rates, callers must choose the
// sampling cadence carefully.
func NewUint32CounterDelta(prev, cur uint32) Uint32CounterDelta {
	return Uint32CounterDelta{
		Previous: prev,
		Current:  cur,
		Value:    cur - prev,
		Wrapped:  cur < prev,
	}
}

// IsZero reports whether the computed delta value is zero.
//
// IsZero means no progress was visible between the two samples under the
// single-wrap model. It does not prove that no increments occurred if the
// counter wrapped exactly back to the previous value between samples. Such
// multiple-wrap cases must be prevented by reasonable sampling cadence.
func (d Uint32CounterDelta) IsZero() bool {
	return d.Value == 0
}
