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

package wait

import (
	"context"
	"math"
	"time"
)

const (
	// errNilContext is the panic value used when a public wait primitive receives
	// a nil context.
	//
	// A nil context is a caller programming error. Wait primitives require an
	// explicit cancellation scope so cancellation, timeout, shutdown, and
	// ownership semantics are visible at the API boundary. Callers that do not
	// need a narrower cancellation scope should pass context.Background.
	errNilContext = "wait: nil context"

	// errNonPositiveInterval is the panic value used when a fixed-cadence wait
	// loop receives a zero or negative interval.
	//
	// Fixed-cadence loops require a strictly positive interval. A non-positive
	// interval would either produce a busy loop or define surprising immediate
	// re-evaluation semantics, so it is rejected at the loop API boundary.
	errNonPositiveInterval = "wait: non-positive interval"

	// errNegativeJitterFactor is the panic value used when jitter receives a
	// negative factor.
	//
	// This package models jitter as positive extra delay only. A negative factor
	// would shorten the base duration and silently change caller-owned cadence
	// policy, so it is rejected as invalid configuration.
	errNegativeJitterFactor = "wait: negative jitter factor"

	// errNonFiniteJitterFactor is the panic value used when jitter receives NaN or
	// an infinite factor.
	//
	// Jitter factors participate in duration arithmetic. Non-finite values do not
	// describe a bounded runtime delay and are rejected before any calculation is
	// attempted.
	errNonFiniteJitterFactor = "wait: non-finite jitter factor"

	// errNilOption is the panic value used when a public wait primitive receives a
	// nil functional option.
	//
	// Nil options are programming errors. Accepting them silently would hide a
	// broken option construction path and make the final wait configuration depend
	// on accidental nil values.
	errNilOption = "wait: nil option"

	// errNilTimer is the panic value used when a Timer method is called on a nil
	// receiver or on a zero-value Timer.
	//
	// Timer values must be created with NewTimer so the wrapper has explicit
	// ownership over an initialized runtime timer. A nil or zero-value Timer is a
	// construction-time programming error, not a recoverable runtime condition.
	errNilTimer = "wait: nil timer"
)

// requireContext panics when ctx is nil.
//
// Nil contexts are rejected at public wait primitive boundaries. Lower-level
// helpers may then pass the context through without repeating nil checks on
// every internal step. This keeps validation close to caller input while keeping
// hot condition-evaluation paths small.
func requireContext(ctx context.Context) {
	if ctx == nil {
		panic(errNilContext)
	}
}

// requirePositiveInterval panics when interval is zero or negative.
//
// The validation belongs to fixed-cadence loops such as Until. It intentionally
// does not apply to Delay or Timer construction, where a non-positive duration
// has useful immediate-wait semantics.
func requirePositiveInterval(interval time.Duration) {
	if interval <= 0 {
		panic(errNonPositiveInterval)
	}
}

// requireJitterFactor panics when factor is not a valid positive-jitter factor.
//
// A valid jitter factor is finite and non-negative. The value may be zero, which
// means that no extra jitter is applied and the base duration is returned
// unchanged.
func requireJitterFactor(factor float64) {
	if math.IsNaN(factor) || math.IsInf(factor, 0) {
		panic(errNonFiniteJitterFactor)
	}
	if factor < 0 {
		panic(errNegativeJitterFactor)
	}
}

// requireOption panics when opt is nil.
//
// Nil options are rejected before any condition evaluation, delay, timer
// allocation, or runtime loop side effect occurs. This makes invalid option
// construction fail at the wait API boundary rather than inside loop execution.
func requireOption(opt Option) {
	if opt == nil {
		panic(errNilOption)
	}
}

// requireUsable verifies that t owns an initialized runtime timer.
//
// Timer's zero value is intentionally invalid. Requiring construction through
// NewTimer keeps ownership explicit and avoids silently creating timers with an
// unclear duration or lifecycle state.
func (t *Timer) requireUsable() {
	if t == nil || t.timer == nil {
		panic(errNilTimer)
	}
}
