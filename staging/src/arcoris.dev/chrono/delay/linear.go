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
	// errNegativeLinearInitialDelay is the stable diagnostic text used when
	// Linear receives a negative initial delay.
	//
	// Linear schedules model concrete runtime durations between owner-controlled
	// loop steps. A zero initial delay is valid and means the first step may
	// continue immediately, but a negative initial delay has no meaningful timer,
	// clock, retry, polling, reconnect, or cooldown interpretation.
	errNegativeLinearInitialDelay = "delay: negative linear initial delay"

	// errNegativeLinearStep is the stable diagnostic text used when Linear
	// receives a negative step.
	//
	// Linear growth is monotonic by construction: each produced delay is the
	// initial delay plus a non-negative multiple of step. A negative step would
	// make later delays smaller and may eventually produce negative runtime
	// durations.
	errNegativeLinearStep = "delay: negative linear step"
)

// Linear returns a schedule whose delays grow by a fixed step each time.
//
// Every sequence created by the returned Schedule starts at initial and adds
// step after each produced delay. The first call to Next returns initial. The
// second call returns initial+step. The third call returns initial+2*step, and
// so on.
//
// For example:
//
//	delay.Linear(100*time.Millisecond, 50*time.Millisecond)
//
// produces:
//
//	100*time.Millisecond
//	150*time.Millisecond
//	200*time.Millisecond
//	250*time.Millisecond
//	...
//
// Linear is useful for:
//
//   - retry paths where fixed intervals are too aggressive;
//   - interactive or near-interactive loops where exponential growth is too
//     steep;
//   - reconnect or polling loops that should slow down predictably;
//   - tests that need deterministic increasing delay values;
//   - composing with Cap to bound the maximum delay;
//   - composing with jitter wrappers to desynchronize otherwise identical loops.
//
// Linear does not provide as much overload relief as exponential growth. For
// remote dependencies under sustained failure, callers should usually prefer a
// capped exponential schedule with jitter.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence with its own index. Sequence values are
// single-owner by default and should not be shared across unrelated runtime
// loops.
//
// Linear does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Linear panics when initial or step is negative. A zero initial delay is valid.
// A zero step is valid and makes the schedule behave like Fixed(initial).
//
// If initial+step*index would overflow time.Duration, the sequence saturates at
// the largest representable duration instead of returning a wrapped negative
// value.
func Linear(initial, step time.Duration) Schedule {
	requireNonNegativeDuration(initial, errNegativeLinearInitialDelay)
	requireNonNegativeDuration(step, errNegativeLinearStep)

	return linearSchedule{
		initial: initial,
		step:    step,
	}
}

// linearSchedule is the reusable recipe behind Linear.
//
// The type stores only the validated initial delay and step. It does not store
// iteration state, previous failures, random state, timer state, context state,
// retry state, or ownership metadata. Per-owner iteration is represented by
// linearSequence values returned from NewSequence.
type linearSchedule struct {
	// initial is the first delay returned by each sequence created from this
	// schedule.
	//
	// A zero value is valid and represents immediate continuation on the first
	// step. Negative values are rejected by Linear before a linearSchedule can be
	// constructed.
	initial time.Duration

	// step is the non-negative amount added after each produced delay.
	//
	// A zero value is valid and makes the sequence return initial forever.
	// Negative values are rejected by Linear before a linearSchedule can be
	// constructed.
	step time.Duration
}

// NewSequence returns an independent linear delay sequence.
//
// The returned sequence starts at the schedule's initial delay and advances by
// the schedule's step on each call to Next. Multiple sequences created from the
// same schedule advance independently because each sequence owns its own index.
func (s linearSchedule) NewSequence() Sequence {
	return &linearSequence{
		initial: s.initial,
		step:    s.step,
		next:    0,
	}
}

// linearSequence is the per-owner delay stream produced by Linear.
//
// The sequence stores the immutable schedule parameters plus the next index to
// use. The index belongs to this sequence instance so independent retry,
// polling, reconnect, or cooldown loops do not interfere with one another.
type linearSequence struct {
	// initial is the first delay returned by this sequence.
	initial time.Duration

	// step is the non-negative amount multiplied by next and added to initial.
	step time.Duration

	// next is the zero-based index used to compute the next delay.
	//
	// The first call uses next=0, the second call uses next=1, and so on. The
	// value is sequence-owned mutable state and is not safe for concurrent
	// mutation unless the caller provides external synchronization.
	next uint64
}

// Next returns the next linear delay and reports that the sequence is still
// available.
//
// Linear sequences are intentionally infinite. Exhaustion, if required, should
// be provided by a higher-level owner such as retry max attempts or by a finite
// wrapper such as Limit.
//
// The returned delay is computed with saturating arithmetic. Once the calculated
// value reaches the maximum representable time.Duration, later calls continue to
// return that maximum value.
func (s *linearSequence) Next() (time.Duration, bool) {
	d := linearDelay(s.initial, s.step, s.next)

	if s.next < ^uint64(0) {
		s.next++
	}

	return d, true
}

// linearDelay returns initial + step*index using saturating duration arithmetic.
//
// The helper assumes initial and step are non-negative. Linear enforces that
// invariant before constructing schedules. Keeping the arithmetic in a helper
// makes overflow behavior explicit and testable without depending on a long
// runtime sequence.
func linearDelay(initial, step time.Duration, index uint64) time.Duration {
	if step == 0 || index == 0 {
		return initial
	}

	product := saturatingDurationMul(step, index)
	return saturatingDurationAdd(initial, product)
}
