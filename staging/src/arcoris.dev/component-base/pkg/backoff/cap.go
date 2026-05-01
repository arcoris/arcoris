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
	// errNilCapSchedule is the stable diagnostic text used when Cap receives a
	// nil child schedule.
	//
	// Cap is a wrapper schedule. It cannot produce delay values without a child
	// schedule to wrap. A nil child indicates invalid package use rather than
	// finite schedule exhaustion. The constructor panics immediately so invalid
	// configuration is detected at the schedule boundary instead of failing later
	// in retry, polling, reconnect, or cooldown loop code.
	errNilCapSchedule = "backoff: nil cap schedule"

	// errNegativeCapMaxDelay is the stable diagnostic text used when Cap receives
	// a negative maximum delay.
	//
	// Cap limits concrete runtime delays. A zero maximum is valid and forces all
	// available child delays to immediate continuation, but a negative maximum
	// delay has no meaningful timer, clock, retry, polling, reconnect, or
	// cooldown interpretation. The constructor panics immediately instead of
	// allowing invalid delay bounds to leak into runtime loop code.
	errNegativeCapMaxDelay = "backoff: negative cap maximum delay"

	// errCapScheduleReturnedNilSequence is the stable diagnostic text used when
	// the child schedule wrapped by Cap returns a nil Sequence.
	//
	// Schedule.NewSequence must return a usable Sequence. Returning nil violates
	// the Schedule contract and would otherwise move a construction-time
	// programming error into a later runtime path. Cap checks this boundary when
	// creating its own sequence.
	errCapScheduleReturnedNilSequence = "backoff: cap schedule returned nil Sequence"

	// errCapScheduleReturnedNegativeDelay is the stable diagnostic text used when
	// the child sequence wrapped by Cap returns a negative delay with ok=true.
	//
	// Sequence implementations must return non-negative delays when ok=true.
	// Cap cannot safely interpret a negative child delay as a runtime duration,
	// even though it could mathematically clamp it. The wrapper panics
	// immediately so the child contract violation is visible at the backoff
	// boundary instead of being silently hidden.
	errCapScheduleReturnedNegativeDelay = "backoff: cap schedule returned negative delay"
)

// Cap returns a schedule that limits every delay produced by schedule.
//
// Cap is a wrapper schedule. It creates a child sequence from the wrapped
// schedule and returns the child's delay values with an upper bound applied.
// Values smaller than or equal to maxDelay pass through unchanged. Values larger
// than maxDelay are returned as maxDelay.
//
// For example:
//
//	backoff.Cap(
//	    backoff.Exponential(100*time.Millisecond, 2.0),
//	    time.Second,
//	)
//
// may produce:
//
//	100*time.Millisecond
//	200*time.Millisecond
//	400*time.Millisecond
//	800*time.Millisecond
//	time.Second
//	time.Second
//	...
//
// Cap is useful for:
//
//   - bounding exponential, linear, Fibonacci, or random schedules;
//   - enforcing a hard maximum retry delay;
//   - preventing long-running sequences from growing beyond operational limits;
//   - composing with jitter wrappers where a final hard bound is required;
//   - tests that need predictable upper-bound behavior.
//
// Cap preserves child exhaustion. If the child sequence returns ok=false, the
// capped sequence also returns ok=false. Cap does not turn a finite child
// schedule into an infinite one and does not turn an infinite child schedule into
// a finite one.
//
// Composition order matters. Cap limits the output of the schedule it wraps. If
// a jitter wrapper is applied after Cap, that later wrapper may increase the
// final delay unless the jitter algorithm itself guarantees otherwise. Use Cap
// as the outermost wrapper when maxDelay must be a hard final bound.
//
// A zero maxDelay is valid. It makes every available child delay return as zero
// while still preserving child exhaustion. This is different from Immediate:
// Cap(schedule, 0) keeps the child sequence's finite or infinite availability
// semantics but clamps all available delays to immediate continuation.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates a fresh child
// sequence and wraps it in an independent capped sequence.
//
// Cap does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Cap panics when schedule is nil or maxDelay is negative.
func Cap(schedule Schedule, maxDelay time.Duration) Schedule {
	requireSchedule(schedule, errNilCapSchedule)
	requireNonNegativeDuration(maxDelay, errNegativeCapMaxDelay)

	return capSchedule{
		schedule: schedule,
		maxDelay: maxDelay,
	}
}

// capSchedule is the reusable recipe behind Cap.
//
// The type stores the child schedule and the validated maximum delay. It does
// not store iteration state, previous failures, random state, timer state,
// context state, retry state, or ownership metadata. Per-owner iteration is
// represented by capSequence values returned from NewSequence.
type capSchedule struct {
	// schedule is the child schedule whose delay values are capped.
	//
	// The value is non-nil. Cap validates this before constructing capSchedule.
	schedule Schedule

	// maxDelay is the inclusive upper bound applied to every child delay.
	//
	// A zero value is valid and clamps every available child delay to immediate
	// continuation. Negative values are rejected by Cap before a capSchedule can
	// be constructed.
	maxDelay time.Duration
}

// NewSequence returns an independent capped delay sequence.
//
// The returned sequence wraps a fresh child sequence created from the child
// schedule. Multiple capped sequences created from the same capSchedule advance
// independently because each one owns its own child sequence.
func (s capSchedule) NewSequence() Sequence {
	child := s.schedule.NewSequence()
	requireSequence(child, errCapScheduleReturnedNilSequence)

	return capSequence{
		child:    child,
		maxDelay: s.maxDelay,
	}
}

// capSequence is the per-owner delay stream produced by Cap.
//
// The sequence delegates availability and raw delay generation to child. It
// applies maxDelay only to available child delays. Exhaustion is passed through
// unchanged.
type capSequence struct {
	// child is the wrapped per-owner sequence.
	//
	// The value is non-nil. capSchedule.NewSequence validates this before
	// constructing capSequence.
	child Sequence

	// maxDelay is the inclusive upper bound applied to child delays.
	maxDelay time.Duration
}

// Next returns the next child delay capped to maxDelay.
//
// If the child sequence is exhausted, Next returns delay=0, ok=false. If the
// child returns a negative delay with ok=true, Next panics because the child has
// violated the Sequence contract.
func (s capSequence) Next() (time.Duration, bool) {
	delay, ok := s.child.Next()
	if !ok {
		return 0, false
	}
	requireNonNegativeSequenceDelay(delay, ok, errCapScheduleReturnedNegativeDelay)

	return capDuration(delay, s.maxDelay), true
}
