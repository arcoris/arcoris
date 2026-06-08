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

package wait

import (
	"context"
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
