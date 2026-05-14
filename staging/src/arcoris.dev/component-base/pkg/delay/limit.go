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
	// errNilLimitSchedule is the stable diagnostic text used when Limit receives
	// a nil child schedule.
	//
	// Limit is a wrapper schedule. It cannot produce delay values without a child
	// schedule to wrap. A nil child indicates invalid package use rather than
	// finite schedule exhaustion. The constructor panics immediately so invalid
	// configuration is detected at the schedule boundary instead of failing later
	// in retry, polling, reconnect, or cooldown loop code.
	errNilLimitSchedule = "delay: nil limit schedule"

	// errNegativeLimitMaxDelays is the stable diagnostic text used when Limit
	// receives a negative delay count.
	//
	// Limit bounds how many delay values a child sequence may expose. A zero
	// limit is valid and means immediate exhaustion, but a negative count has no
	// meaningful sequence interpretation. The constructor panics immediately
	// instead of allowing invalid ownership policy to leak into runtime loop code.
	errNegativeLimitMaxDelays = "delay: negative limit maximum delays"

	// errLimitScheduleReturnedNilSequence is the stable diagnostic text used when
	// the child schedule wrapped by Limit returns a nil Sequence.
	//
	// Schedule.NewSequence must return a usable Sequence. Returning nil violates
	// the Schedule contract and would otherwise move a construction-time
	// programming error into a later runtime path. Limit checks this boundary
	// when creating a non-empty limited sequence.
	errLimitScheduleReturnedNilSequence = "delay: limit schedule returned nil Sequence"

	// errLimitScheduleReturnedNegativeDelay is the stable diagnostic text used
	// when the child sequence wrapped by Limit returns a negative delay with
	// ok=true.
	//
	// Sequence implementations must return non-negative delays when ok=true.
	// Limit cannot safely interpret a negative child delay as a runtime duration.
	// The wrapper panics immediately so the child contract violation is visible
	// at the delay boundary instead of being silently hidden.
	errLimitScheduleReturnedNegativeDelay = "delay: limit schedule returned negative delay"
)

// Limit returns a schedule that exposes at most maxDelays values from schedule.
//
// Limit is a wrapper schedule. It creates a child sequence from the wrapped
// schedule and forwards at most maxDelays available delay values. After that
// limit is reached, the sequence reports exhaustion by returning ok=false.
//
// For example:
//
//	delay.Limit(
//	    delay.Fixed(100*time.Millisecond),
//	    3,
//	)
//
// produces:
//
//	100*time.Millisecond, true
//	100*time.Millisecond, true
//	100*time.Millisecond, true
//	0, false
//	0, false
//	...
//
// Limit is useful for:
//
//   - making infinite schedules finite;
//   - bounding retry delay budgets at the delay layer;
//   - tests that need deterministic exhaustion;
//   - exposing only the first N values of exponential, linear, Fibonacci, fixed,
//     random, or jittered schedules;
//   - combining with retry max attempts when retry should distinguish attempt
//     exhaustion from schedule exhaustion.
//
// Limit preserves child exhaustion. If the child sequence reports ok=false
// before maxDelays values have been returned, the limited sequence also reports
// ok=false. Limit does not manufacture additional values after the child is
// exhausted.
//
// A zero maxDelays is valid. It creates a schedule whose sequences are
// exhausted immediately. In that case NewSequence does not create a child
// sequence because no child delay can ever be observed.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates an independent
// limited sequence with its own remaining-count state.
//
// Limit does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Limit panics when sched is nil or n is negative.
func Limit(sched Schedule, n int) Schedule {
	requireSchedule(sched, errNilLimitSchedule)
	requireNonNegativeCount(n, errNegativeLimitMaxDelays)

	return limitSchedule{
		schedule:  sched,
		maxDelays: n,
	}
}

// limitSchedule is the reusable recipe behind Limit.
//
// The type stores the child schedule and the validated maximum number of delay
// values to expose. It does not store iteration state, previous failures, random
// state, timer state, context state, retry state, or ownership metadata.
// Per-owner iteration is represented by limitSequence values returned from
// NewSequence.
type limitSchedule struct {
	// schedule is the child schedule whose delay values are exposed up to
	// maxDelays times.
	//
	// The value is non-nil. Limit validates this before constructing
	// limitSchedule.
	schedule Schedule

	// maxDelays is the maximum number of available child delays to expose.
	//
	// A zero value is valid and makes every sequence created from this schedule
	// exhausted immediately. Negative values are rejected by Limit before a
	// limitSchedule can be constructed.
	maxDelays int
}

// NewSequence returns an independent limited delay sequence.
//
// If maxDelays is zero, NewSequence returns an already exhausted sequence and
// does not ask the child schedule to create a Sequence. This preserves the
// meaning of a zero limit: no child delay can be observed.
//
// If maxDelays is positive, NewSequence creates a fresh child sequence and wraps
// it with independent remaining-count state.
func (s limitSchedule) NewSequence() Sequence {
	if s.maxDelays == 0 {
		return &limitSequence{
			remaining: 0,
		}
	}

	child := s.schedule.NewSequence()
	requireSequence(child, errLimitScheduleReturnedNilSequence)

	return &limitSequence{
		child:     child,
		remaining: s.maxDelays,
	}
}

// limitSequence is the per-owner delay stream produced by Limit.
//
// The sequence delegates raw delay generation to child while remaining is
// positive. Once remaining reaches zero, Next reports exhaustion without calling
// the child sequence again.
type limitSequence struct {
	// child is the wrapped per-owner sequence.
	//
	// The value is nil only for zero-limit sequences. Next checks remaining
	// before using child, so a zero-limit sequence never dereferences child.
	child Sequence

	// remaining is the number of child delay values that may still be exposed.
	//
	// The value belongs to this sequence instance. Independent retry, polling,
	// reconnect, or cooldown loops created from the same limitSchedule do not
	// share remaining-count state.
	remaining int
}

// Next returns the next child delay while the limit has remaining capacity.
//
// If the limit has already been reached, Next returns delay=0, ok=false without
// calling the child sequence. If the child sequence is exhausted before the limit
// is reached, Next also returns delay=0, ok=false.
//
// If the child returns a negative delay with ok=true, Next panics because the
// child has violated the Sequence contract.
func (s *limitSequence) Next() (time.Duration, bool) {
	if s.remaining <= 0 {
		return 0, false
	}

	d, ok := s.child.Next()
	if !ok {
		s.remaining = 0
		return 0, false
	}
	requireNonNegativeSequenceDelay(d, ok, errLimitScheduleReturnedNegativeDelay)

	s.remaining--

	return d, true
}
