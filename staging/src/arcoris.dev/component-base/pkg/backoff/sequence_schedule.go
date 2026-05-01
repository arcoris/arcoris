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
	// errNegativeSequenceDelay is the stable diagnostic text used when Delays
	// receives a negative delay.
	//
	// Explicit delay schedules model concrete runtime durations between
	// owner-controlled loop steps. A zero delay is valid and means immediate
	// continuation, but a negative delay has no meaningful timer, clock, retry,
	// polling, reconnect, or cooldown interpretation. The constructor panics
	// immediately so invalid configuration is detected at the schedule boundary
	// instead of leaking into runtime loop code.
	errNegativeSequenceDelay = "backoff: negative sequence delay"
)

// Delays returns a finite schedule backed by an explicit list of delays.
//
// Delays is useful when callers need a precise, non-formulaic delay profile.
// Each Sequence created by the returned Schedule walks the configured delays in
// order and then reports exhaustion by returning ok=false from Next.
//
// For example:
//
//	schedule := backoff.Delays(
//	    0,
//	    10*time.Millisecond,
//	    100*time.Millisecond,
//	)
//
// The produced sequence yields:
//
//	0, true
//	10*time.Millisecond, true
//	100*time.Millisecond, true
//	0, false
//
// Delays is useful for:
//
//   - tests that need exact deterministic delay values;
//   - retry profiles with an immediate first retry and slower later retries;
//   - migration from hand-written retry loops;
//   - protocol adapters with a small fixed sequence of delays;
//   - composing finite delay streams with wrappers such as jitter or cap.
//
// Delays copies the input slice before storing it. Mutating the caller-owned
// slice after construction does not affect the returned Schedule. This preserves
// the Schedule contract: schedule values are reusable delay recipes, not views
// over caller-owned mutable memory.
//
// An empty delay list is valid. It creates a schedule whose sequences are
// exhausted immediately. This is useful in tests and in higher-level policies
// that want an explicit "no delay values available" schedule.
//
// A zero delay is valid and means immediate continuation. A negative delay is a
// programming error. Delays panics when any configured delay is negative.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence with its own cursor. Sequence values are
// single-owner by default and should not be shared across unrelated runtime
// loops.
//
// Delays does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
func Delays(delays ...time.Duration) Schedule {
	copied := make([]time.Duration, len(delays))
	copy(copied, delays)

	for _, delay := range copied {
		requireNonNegativeDuration(delay, errNegativeSequenceDelay)
	}

	return sequenceSchedule{delays: copied}
}

// sequenceSchedule is the reusable recipe behind Delays.
//
// The type stores a copied immutable list of non-negative delays. It does not
// store cursor state, previous delay state, random state, timer state, context
// state, retry state, or ownership metadata. Per-owner iteration is represented
// by sequenceScheduleSequence values returned from NewSequence.
type sequenceSchedule struct {
	// delays is the copied list of delay values returned by sequences created
	// from this schedule.
	//
	// The slice MUST NOT be mutated after construction. The public constructor
	// copies caller input before storing it, so caller-owned mutation cannot
	// affect this value.
	delays []time.Duration
}

// NewSequence returns an independent finite delay sequence.
//
// The returned sequence starts at the first configured delay. Multiple sequences
// created from the same schedule advance independently because each sequence has
// its own cursor.
func (s sequenceSchedule) NewSequence() Sequence {
	return &sequenceScheduleSequence{
		delays: s.delays,
		next:   0,
	}
}

// sequenceScheduleSequence is the per-owner delay stream produced by Delays.
//
// The sequence holds only a cursor into the schedule's copied delay slice. It
// does not mutate the slice. It is finite: after all delays have been returned,
// Next reports exhaustion with ok=false.
type sequenceScheduleSequence struct {
	// delays is the immutable copied delay list owned by the parent schedule.
	//
	// The sequence reads this slice but does not modify it. Sharing the slice
	// between many independent sequences is safe as long as the schedule keeps it
	// immutable after construction.
	delays []time.Duration

	// next is the index of the next delay to return.
	//
	// next belongs to this sequence instance. It is intentionally per-sequence
	// state so independent retry, polling, reconnect, or cooldown loops do not
	// interfere with one another.
	next int
}

// Next returns the next configured delay or reports sequence exhaustion.
//
// When a delay is available, Next returns that delay with ok=true and advances
// the sequence cursor. When all configured delays have been consumed, Next
// returns delay=0, ok=false.
//
// The delay returned with ok=false is intentionally not meaningful. Owners
// should inspect ok before using delay.
func (s *sequenceScheduleSequence) Next() (time.Duration, bool) {
	if s.next >= len(s.delays) {
		return 0, false
	}

	delay := s.delays[s.next]
	s.next++

	return delay, true
}
