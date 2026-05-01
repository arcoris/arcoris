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

import (
	"math"
	"time"
)

const (
	// errNonPositiveExponentialInitialDelay is the stable diagnostic text used
	// when Exponential receives a non-positive initial delay.
	//
	// Exponential schedules model growing positive runtime durations between
	// owner-controlled loop steps. A zero or negative initial delay cannot grow
	// into a meaningful positive exponential sequence. Callers that need
	// immediate continuation should compose an explicit immediate delay before an
	// exponential schedule instead of using a zero initial value.
	errNonPositiveExponentialInitialDelay = "backoff: non-positive exponential initial delay"

	// errInvalidExponentialMultiplier is the stable diagnostic text used when
	// Exponential receives an invalid multiplier.
	//
	// The multiplier must be finite and greater than one. A multiplier less than
	// or equal to one does not describe exponential growth. NaN and infinite
	// values cannot produce a stable runtime delay sequence. The constructor
	// panics immediately so invalid configuration is detected at the schedule
	// boundary instead of leaking into retry, polling, reconnect, or cooldown
	// loop code.
	errInvalidExponentialMultiplier = "backoff: invalid exponential multiplier"
)

// Exponential returns a schedule whose delays grow by a multiplier each time.
//
// Every sequence created by the returned Schedule starts at initial and then
// multiplies the previous mathematical delay by multiplier after each produced
// delay. The first call to Next returns initial. The second call returns roughly
// initial*multiplier. The third call returns roughly initial*multiplier^2, and
// so on.
//
// For example:
//
//	backoff.Exponential(100*time.Millisecond, 2.0)
//
// produces:
//
//	100*time.Millisecond
//	200*time.Millisecond
//	400*time.Millisecond
//	800*time.Millisecond
//	...
//
// Exponential is useful for:
//
//   - retry paths against remote dependencies;
//   - reconnect loops;
//   - controller loops that should back away quickly after repeated failures;
//   - background operations where repeated immediate attempts would increase
//     pressure on a dependency;
//   - composing with Cap to bound the maximum delay;
//   - composing with jitter wrappers to desynchronize otherwise identical loops.
//
// Exponential does not cap itself. Without an explicit Cap wrapper, long-running
// sequences eventually saturate at the largest representable time.Duration.
// Production retry policies should usually compose Exponential with Cap and a
// jitter wrapper once those wrappers are available.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence with its own mathematical delay state.
// Sequence values are single-owner by default and should not be shared across
// unrelated runtime loops.
//
// Exponential does not sleep, create timers, observe context cancellation,
// execute operations, classify errors, retry failed work, log, trace, export
// metrics, rate limit callers, schedule queue items, or make domain decisions.
//
// Exponential panics when initial is not positive or when multiplier is not a
// finite value greater than one.
//
// If the mathematical delay would overflow time.Duration, the sequence
// saturates at the largest representable duration instead of returning a wrapped
// negative value.
func Exponential(initial time.Duration, multiplier float64) Schedule {
	requirePositiveDuration(initial, errNonPositiveExponentialInitialDelay)
	requireFloatGreaterThanOne(multiplier, errInvalidExponentialMultiplier)

	return exponentialSchedule{
		initial:    initial,
		multiplier: multiplier,
	}
}

// exponentialSchedule is the reusable recipe behind Exponential.
//
// The type stores only the validated initial delay and multiplier. It does not
// store iteration state, previous failures, random state, timer state, context
// state, retry state, or ownership metadata. Per-owner iteration is represented
// by exponentialSequence values returned from NewSequence.
type exponentialSchedule struct {
	// initial is the first delay returned by each sequence created from this
	// schedule.
	//
	// The value is strictly positive. Non-positive values are rejected by
	// Exponential before an exponentialSchedule can be constructed.
	initial time.Duration

	// multiplier is the finite factor applied to the mathematical delay after
	// each produced delay.
	//
	// The value is strictly greater than one. Invalid floating-point values are
	// rejected by Exponential before an exponentialSchedule can be constructed.
	multiplier float64
}

// NewSequence returns an independent exponential delay sequence.
//
// The returned sequence starts at the schedule's initial delay and advances by
// the schedule's multiplier on each call to Next. Multiple sequences created
// from the same schedule advance independently because each sequence owns its
// own mathematical delay state.
func (s exponentialSchedule) NewSequence() Sequence {
	return &exponentialSequence{
		next:       float64(s.initial),
		multiplier: s.multiplier,
	}
}

// exponentialSequence is the per-owner delay stream produced by Exponential.
//
// The sequence stores the next mathematical delay as a floating-point nanosecond
// value. This is intentional. If the sequence multiplied already-truncated
// time.Duration values, small durations with small multipliers could stop
// growing forever because each step would round back to the same integer
// nanosecond value.
//
// Returned delays are still ordinary time.Duration values. The floating-point
// state is only an internal mechanism for preserving exponential growth between
// returned integer durations.
type exponentialSequence struct {
	// next is the next mathematical delay in nanoseconds.
	//
	// The value is positive until the sequence saturates. It is converted to a
	// time.Duration on each Next call.
	next float64

	// multiplier is the finite factor applied after each produced delay.
	multiplier float64

	// saturated records that the sequence has reached the maximum representable
	// time.Duration.
	//
	// Once saturated, later calls return maxDuration forever.
	saturated bool
}

// Next returns the next exponential delay and reports that the sequence is still
// available.
//
// Exponential sequences are intentionally infinite. Exhaustion, if required,
// should be provided by a higher-level owner such as retry max attempts or by a
// finite wrapper such as Limit.
//
// The returned delay is computed with saturating arithmetic. Once the calculated
// value reaches the maximum representable time.Duration, later calls continue to
// return that maximum value.
func (s *exponentialSequence) Next() (time.Duration, bool) {
	if s.saturated {
		return maxDuration, true
	}

	delay := durationFromFloat(s.next)
	s.advance()

	return delay, true
}

// advance moves the sequence to the next mathematical exponential delay.
//
// If the next mathematical value cannot be represented as a finite duration, the
// sequence enters saturated state. Saturation is sticky: once the sequence
// reaches maxDuration, every later call returns that value.
func (s *exponentialSequence) advance() {
	next := s.next * s.multiplier
	if math.IsInf(next, 0) || math.IsNaN(next) || next >= maxDurationFloat {
		s.saturated = true
		s.next = maxDurationFloat
		return
	}

	s.next = next
}
