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

package backoff

import "time"

const (
	// errNonPositiveFibonacciBaseDelay is the stable diagnostic text used when
	// Fibonacci receives a non-positive base delay.
	//
	// Fibonacci schedules model growing positive runtime durations between
	// owner-controlled loop steps. A zero or negative base delay cannot produce a
	// meaningful positive Fibonacci backoff sequence. Callers that need an
	// immediate first step should compose an explicit zero delay before a
	// Fibonacci schedule instead of using a zero base value.
	errNonPositiveFibonacciBaseDelay = "backoff: non-positive fibonacci base delay"
)

// Fibonacci returns a schedule whose delays grow by Fibonacci numbers.
//
// Every sequence created by the returned Schedule starts with base and then
// follows the Fibonacci growth pattern:
//
//	base
//	base
//	2*base
//	3*base
//	5*base
//	8*base
//	...
//
// Fibonacci backoff grows faster than linear backoff but usually slower than
// exponential backoff. It is useful when fixed or linear intervals are too
// aggressive, but exponential growth backs off too quickly for the owner.
//
// Fibonacci is useful for:
//
//   - reconnect loops that should slow down progressively but not too sharply;
//   - retry paths where linear growth is too weak and exponential growth is too
//     steep;
//   - background operations with moderate recovery windows;
//   - tests that need deterministic non-linear delay growth;
//   - composing with Cap to bound the maximum delay;
//   - composing with jitter wrappers to desynchronize otherwise identical loops.
//
// Fibonacci does not cap itself. Without an explicit Cap wrapper, long-running
// sequences eventually saturate at the largest representable time.Duration.
// Production retry policies should usually compose Fibonacci with Cap and,
// where many owners may retry concurrently, a jitter wrapper.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence with its own Fibonacci state. Sequence values
// are single-owner by default and should not be shared across unrelated runtime
// loops.
//
// Fibonacci does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Fibonacci panics when base is not positive.
//
// If the next Fibonacci delay would overflow time.Duration, the sequence
// saturates at the largest representable duration instead of returning a wrapped
// negative value.
func Fibonacci(base time.Duration) Schedule {
	requirePositiveDuration(base, errNonPositiveFibonacciBaseDelay)

	return fibonacciSchedule{base: base}
}

// fibonacciSchedule is the reusable recipe behind Fibonacci.
//
// The type stores only the validated base delay. It does not store iteration
// state, previous failures, random state, timer state, context state, retry
// state, or ownership metadata. Per-owner iteration is represented by
// fibonacciSequence values returned from NewSequence.
type fibonacciSchedule struct {
	// base is the positive delay used as the first two Fibonacci values and as
	// the unit multiplied by later Fibonacci numbers.
	//
	// Non-positive values are rejected by Fibonacci before a fibonacciSchedule can
	// be constructed.
	base time.Duration
}

// NewSequence returns an independent Fibonacci delay sequence.
//
// The returned sequence starts with base, base, 2*base, and then continues by
// adding the two previous delay values. Multiple sequences created from the same
// schedule advance independently because each sequence owns its own Fibonacci
// state.
func (s fibonacciSchedule) NewSequence() Sequence {
	return &fibonacciSequence{
		previous: 0,
		current:  s.base,
	}
}

// fibonacciSequence is the per-owner delay stream produced by Fibonacci.
//
// The sequence stores the two values needed to produce the next Fibonacci delay.
// State is stored as time.Duration values rather than as integer Fibonacci
// multipliers so the implementation can use saturating duration addition
// directly and avoid multiplier overflow.
//
// The sequence starts in a pre-first state:
//
//	previous = 0
//	current  = base
//
// The first call returns current, then advances to:
//
//	previous = base
//	current  = base
//
// The second call returns base again, then advances to:
//
//	previous = base
//	current  = 2*base
//
// This produces the expected multiplier sequence:
//
//	1, 1, 2, 3, 5, 8, ...
type fibonacciSequence struct {
	// previous is the delay that preceded current in the Fibonacci stream.
	//
	// The initial zero value is an internal bootstrap value. It is not returned
	// by the sequence. It allows the first advance to preserve current=base for
	// the second returned value.
	previous time.Duration

	// current is the delay returned by the next call to Next.
	//
	// The value is always positive until saturation. Once current reaches
	// maxDuration, it remains saturated.
	current time.Duration

	// saturated records that the sequence has reached the maximum representable
	// time.Duration.
	//
	// Once saturated, later calls return maxDuration forever.
	saturated bool
}

// Next returns the next Fibonacci delay and reports that the sequence is still
// available.
//
// Fibonacci sequences are intentionally infinite. Exhaustion, if required,
// should be provided by a higher-level owner such as retry max attempts or by a
// finite wrapper such as Limit.
//
// The returned delay is computed with saturating arithmetic. Once the calculated
// value reaches the maximum representable time.Duration, later calls continue to
// return that maximum value.
func (s *fibonacciSequence) Next() (time.Duration, bool) {
	if s.saturated {
		return maxDuration, true
	}

	delay := s.current
	s.advance()

	return delay, true
}

// advance moves the sequence to the next Fibonacci delay.
//
// The transition is:
//
//	next     = previous + current
//	previous = current
//	current  = next
//
// Addition uses saturating duration arithmetic. Saturation is sticky: once the
// sequence reaches maxDuration, every later call returns that value.
func (s *fibonacciSequence) advance() {
	next := saturatingDurationAdd(s.previous, s.current)

	s.previous = s.current
	s.current = next

	if next == maxDuration {
		s.saturated = true
	}
}
