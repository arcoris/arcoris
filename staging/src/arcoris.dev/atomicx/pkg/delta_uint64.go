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

// Uint64CounterDelta describes the difference between two Uint64Counter samples.
//
// Uint64CounterDelta is a copyable value object. It intentionally does not
// contain noCopy because it does not own atomic state, does not synchronize
// access to anything, and is safe to pass by value.
//
// The delta uses unsigned modulo arithmetic:
//
//   - if Current >= Previous, Value is Current - Previous;
//   - if Current < Previous, the counter is interpreted as having wrapped once,
//     and Value is still computed as Current - Previous using uint64 arithmetic.
//
// Wrapped reports whether Current was lower than Previous and the delta was
// interpreted as a single uint64 wrap.
//
// Multiple wraps between two samples cannot be detected from two uint64 values
// alone. Callers that rely on accurate activity windows must sample frequently
// enough to make multiple wraps impossible in practice.
type Uint64CounterDelta struct {
	// Previous is the older lifetime counter value.
	Previous uint64

	// Current is the newer lifetime counter value.
	Current uint64

	// Value is the computed delta between Previous and Current.
	//
	// Value is computed with uint64 modulo arithmetic, so it remains correct for
	// one wrap between Previous and Current.
	Value uint64

	// Wrapped reports whether Current was lower than Previous and the delta was
	// interpreted as a single uint64 wrap.
	Wrapped bool
}

// NewUint64CounterDelta computes a wrap-aware delta between two uint64 lifetime
// counter values.
//
// The previous value must be the older sample and current must be the newer
// sample. The function cannot verify temporal ordering because it receives only
// raw counter values.
//
// uint64 subtraction already has the desired modulo behavior for a single wrap,
// so current - previous works for both ordinary and wrapped cases. Wrapped is
// recorded separately because diagnostics, tests, and higher-level reporting may
// need to know that a wrap was observed.
//
// Multiple wraps cannot be detected from previous and current alone.
func NewUint64CounterDelta(prev, cur uint64) Uint64CounterDelta {
	return Uint64CounterDelta{
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
func (d Uint64CounterDelta) IsZero() bool {
	return d.Value == 0
}
