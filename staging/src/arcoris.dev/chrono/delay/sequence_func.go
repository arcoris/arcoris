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

package delay

import "time"

const (
	// errNilSequenceFunc is the stable diagnostic text used when a SequenceFunc
	// method is called on a nil function value.
	//
	// A nil SequenceFunc cannot produce delay values and indicates invalid
	// package use rather than finite sequence exhaustion. The adapter panics
	// immediately with this message so the invalid sequence boundary is visible
	// at the point where the next delay is requested.
	errNilSequenceFunc = "delay: nil SequenceFunc"

	// errSequenceFuncReturnedNegativeDelay is the stable diagnostic text used
	// when a SequenceFunc returns a negative delay with ok=true.
	//
	// Negative delays violate the Sequence contract. A zero delay is valid and
	// means immediate continuation, but a negative delay has no meaningful
	// runtime interpretation for retry, polling, reconnect, or cooldown owners.
	// The adapter panics immediately instead of allowing the invalid value to
	// leak into timer, clock, retry, or wait code.
	errSequenceFuncReturnedNegativeDelay = "delay: SequenceFunc returned negative delay"
)

// SequenceFunc adapts a function into a Sequence.
//
// SequenceFunc is useful for small custom sequences, tests, and adapters that
// need to satisfy Sequence without declaring a named type. The wrapped function
// is called every time Next is invoked.
//
// Example:
//
//	sequence := delay.SequenceFunc(func() (time.Duration, bool) {
//		return time.Second, true
//	})
//	delay, ok := sequence.Next()
//	_ = delay
//	_ = ok
//
// SequenceFunc follows the same responsibility boundary as any other Sequence
// implementation:
//
//   - it may close over per-sequence iteration state;
//   - it may return ok=false to report finite sequence exhaustion;
//   - it must return non-negative delays when ok=true;
//   - it must not sleep, create timers, observe context cancellation, execute
//     operations, classify errors, retry work, log, trace, export metrics,
//     schedule queue items, rate limit callers, or make domain decisions.
//
// A nil SequenceFunc is a programming error. Next panics immediately instead of
// returning a delayed nil dereference from a runtime loop.
//
// A SequenceFunc that returns a negative delay with ok=true violates the
// Sequence contract. Next panics so the invalid adapter is detected before the
// value reaches clock, wait, retry, or controller code.
//
// When ok=false, the delay value is ignored by the Sequence contract. Callers
// should not inspect delay after exhaustion. Implementations may return zero for
// clarity, but the adapter does not require it.
//
// SequenceFunc does not recover panics raised by the wrapped function. Panic
// recovery, if required, belongs to the caller or to an explicit higher-level
// wrapper.
type SequenceFunc func() (delay time.Duration, ok bool)

// Next calls f and returns the delay produced by it.
//
// Next panics when f is nil. It also panics when f returns a negative delay with
// ok=true. Both cases are programming errors, not sequence exhaustion.
func (f SequenceFunc) Next() (delay time.Duration, ok bool) {
	if f == nil {
		panic(errNilSequenceFunc)
	}

	delay, ok = f()
	if ok && delay < 0 {
		panic(errSequenceFuncReturnedNegativeDelay)
	}

	return delay, ok
}
