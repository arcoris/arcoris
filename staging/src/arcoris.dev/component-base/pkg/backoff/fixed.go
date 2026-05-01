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
	// errNegativeFixedDelay is the stable diagnostic text used when Fixed
	// receives a negative delay.
	//
	// Fixed delays model concrete runtime durations between owner-controlled
	// loop steps. A zero delay is valid and means immediate continuation, but a
	// negative delay has no meaningful timer, clock, retry, polling, reconnect,
	// or cooldown interpretation. The constructor panics immediately so invalid
	// configuration is detected at the schedule boundary instead of leaking into
	// runtime loop code.
	errNegativeFixedDelay = "backoff: negative fixed delay"
)

// Fixed returns a schedule that produces an infinite stream of one delay.
//
// Fixed is the simplest non-growing delay schedule. Every sequence created by
// the returned Schedule reports the same delay and ok=true for every call to
// Next. The delay may be zero, in which case the owner may continue immediately
// on every step.
//
// Example:
//
//	schedule := backoff.Fixed(250 * time.Millisecond)
//	sequence := schedule.NewSequence()
//	delay, ok := sequence.Next()
//	_ = delay
//	_ = ok
//
// Fixed is useful for:
//
//   - simple polling loops;
//   - readiness checks;
//   - interactive retry paths with predictable cadence;
//   - tests that need deterministic delay values;
//   - composing with Limit to produce a finite number of identical delays;
//   - composing with jitter wrappers to desynchronize otherwise identical loops.
//
// Fixed does not provide progressive overload relief by itself. If many clients
// use the same fixed delay without jitter, their wake-ups may remain
// synchronized. Distributed retry or polling code should usually combine Fixed
// with a jitter wrapper or use an exponential schedule when repeated failures
// should reduce pressure on a dependency.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence value. The concrete sequence is stateless apart
// from the fixed delay value, but callers should still follow the package-wide
// single-owner Sequence model and avoid sharing one Sequence across unrelated
// runtime loops.
//
// Fixed does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Fixed panics when delay is negative.
func Fixed(delay time.Duration) Schedule {
	requireNonNegativeDuration(delay, errNegativeFixedDelay)

	return fixedSchedule{delay: delay}
}

// fixedSchedule is the reusable recipe behind Fixed.
//
// The type stores only the configured delay. It does not store iteration state,
// previous failures, attempt counters, random state, timers, or ownership
// metadata. Per-owner iteration is represented by fixedSequence values returned
// from NewSequence.
type fixedSchedule struct {
	// delay is the non-negative duration returned by every sequence produced by
	// this schedule.
	//
	// A zero value is valid and represents immediate continuation. Negative
	// values are rejected by Fixed before a fixedSchedule can be constructed.
	delay time.Duration
}

// NewSequence returns an independent fixed-delay sequence.
//
// The returned sequence produces the same non-negative delay forever. Since the
// sequence does not mutate internal state, many equivalent sequences can be
// created cheaply from the same fixedSchedule.
func (s fixedSchedule) NewSequence() Sequence {
	return fixedSequence{delay: s.delay}
}

// fixedSequence is the per-owner delay stream produced by Fixed.
//
// The sequence returns one stable delay forever. It is modeled as a Sequence
// rather than being special-cased by retry or polling code so fixed delays
// compose uniformly with future wrappers such as Limit, Cap, FullJitter,
// PositiveJitter, and ProportionalJitter.
type fixedSequence struct {
	// delay is the non-negative duration returned from every Next call.
	//
	// The value is copied from fixedSchedule at sequence creation time. This
	// keeps the sequence independent from the schedule value and preserves the
	// Schedule/Sequence ownership boundary.
	delay time.Duration
}

// Next returns the configured fixed delay and reports that the sequence is still
// available.
//
// The sequence is intentionally infinite. Exhaustion, if required, should be
// provided by a higher-level owner such as retry max attempts or by a finite
// wrapper such as Limit.
func (s fixedSequence) Next() (time.Duration, bool) {
	return s.delay, true
}
