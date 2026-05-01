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
	// errNegativeRandomMinDelay is the stable diagnostic text used when Random
	// receives a negative minimum delay.
	//
	// Random schedules model concrete runtime durations between owner-controlled
	// loop steps. A zero minimum delay is valid and means the lower bound may
	// allow immediate continuation, but a negative minimum has no meaningful
	// timer, clock, retry, polling, reconnect, or cooldown interpretation. The
	// constructor panics immediately so invalid configuration is detected at the
	// schedule boundary instead of leaking into runtime loop code.
	errNegativeRandomMinDelay = "backoff: negative random minimum delay"

	// errRandomMaxDelayBeforeMinDelay is the stable diagnostic text used when
	// Random receives a maximum delay smaller than the minimum delay.
	//
	// Random schedules generate delays from a closed interval. An interval whose
	// upper bound is smaller than its lower bound is invalid and cannot be
	// interpreted as a runtime delay range. The constructor panics immediately
	// instead of silently swapping bounds or producing surprising delay values.
	errRandomMaxDelayBeforeMinDelay = "backoff: random maximum delay before minimum delay"

	// errNilRandomScheduleSource is the stable diagnostic text used when the
	// package-local random constructor receives a nil RandomSource.
	//
	// Public callers normally configure randomness through RandomOption values.
	// The diagnostic keeps the internal source boundary explicit for tests and
	// package-owned construction paths.
	errNilRandomScheduleSource = "backoff: nil random schedule source"
)

// Random returns a schedule that produces an infinite stream of random delays.
//
// Every sequence created by the returned Schedule reports delay values in the
// closed interval [minDelay, maxDelay]. Both bounds are inclusive. A range with
// equal bounds is valid and behaves like a fixed delay schedule with that value.
//
// Example:
//
//	schedule := backoff.Random(time.Second, 5*time.Second, backoff.WithSeed(42))
//	sequence := schedule.NewSequence()
//	delay, ok := sequence.Next()
//	_ = delay
//	_ = ok
//
// Random is useful for:
//
//   - desynchronizing polling loops;
//   - spreading reconnect attempts across time;
//   - tests and adapters that need a simple random delay range;
//   - composing with Limit to produce a finite random delay stream;
//   - lightweight retry policies where progressive backoff growth is not needed.
//
// Random does not provide progressive overload relief by itself. The delay range
// does not grow after repeated failures. For remote dependencies or overload
// sensitive paths, callers should usually prefer capped exponential schedules
// with jitter once those schedules are available.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence value. The concrete sequence is stateless apart
// from the configured range and draw function, but callers should still follow
// the package-wide single-owner Sequence model and avoid sharing one Sequence
// across unrelated runtime loops.
//
// By default Random uses the package default non-cryptographic source. Options
// may provide deterministic or custom pseudo-random sources. The source is
// stored on the reusable Schedule, while each Sequence owns the random generator
// returned by the source.
//
// Random does not sleep, create timers, observe context cancellation, execute
// operations, classify errors, retry failed work, log, trace, export metrics,
// rate limit callers, schedule queue items, or make domain decisions.
//
// Random panics when minDelay is negative, maxDelay is smaller than minDelay, a
// random option is nil, or the configured random source is nil.
func Random(minDelay, maxDelay time.Duration, opts ...RandomOption) Schedule {
	config := randomOptionsOf(opts...)
	return randomWithSource(minDelay, maxDelay, config.source)
}

// randomWithSource returns a random delay schedule using source for per-sequence
// pseudo-random values.
//
// The helper keeps Random option parsing separate from range validation. It is
// package-local so tests and other package-owned constructors can verify the
// RandomSource boundary directly.
func randomWithSource(
	minDelay time.Duration,
	maxDelay time.Duration,
	source RandomSource,
) Schedule {
	requireNonNegativeDuration(minDelay, errNegativeRandomMinDelay)
	requireDurationNotBefore(maxDelay, minDelay, errRandomMaxDelayBeforeMinDelay)
	requireRandomSource(source, errNilRandomScheduleSource)

	return randomSchedule{
		minDelay: minDelay,
		span:     saturatingDurationSub(maxDelay, minDelay),
		source:   source,
	}
}

// randomSchedule is the reusable recipe behind Random.
//
// The type stores only the validated lower bound, the non-negative span between
// bounds, and the source used to create per-sequence random values. It does not
// store iteration state, previous failures, attempt counters, timers, contexts,
// retry state, or ownership metadata.
type randomSchedule struct {
	// minDelay is the inclusive lower bound returned by sequences created from
	// this schedule.
	//
	// A zero value is valid and allows immediate continuation when the random
	// draw returns zero.
	minDelay time.Duration

	// span is the non-negative distance between the inclusive maximum delay and
	// minDelay.
	//
	// The actual returned delay is minDelay plus an offset in [0, span]. A zero
	// span makes the schedule behave like a fixed delay schedule.
	span time.Duration

	// source creates the per-sequence random generator used for offset draws.
	source RandomSource
}

// NewSequence returns an independent random delay sequence.
//
// The returned sequence produces random values from the schedule's configured
// range forever. Since randomSequence has no cursor and does not mutate the
// schedule, many sequences can be created cheaply from the same randomSchedule.
func (s randomSchedule) NewSequence() Sequence {
	random := s.source.NewRandom()
	requireRandom(random, errNilRandom)

	return &randomSequence{
		minDelay: s.minDelay,
		span:     s.span,
		random:   random,
	}
}

// randomSequence is the per-owner delay stream produced by Random.
//
// The sequence returns a delay in [minDelay, minDelay+span] on every call to
// Next. It is infinite. Exhaustion, if required, should be provided by a
// higher-level owner such as retry max attempts or by a finite wrapper such as
// Limit.
type randomSequence struct {
	// minDelay is the inclusive lower bound returned by this sequence.
	minDelay time.Duration

	// span is the non-negative inclusive offset range above minDelay.
	span time.Duration

	// random is owned by this sequence and used for inclusive offset draws.
	random RandomGenerator
}

// Next returns the next random delay and reports that the sequence is still
// available.
//
// The returned delay is always in [minDelay, minDelay+span]. If span is zero,
// Next returns minDelay without consulting the draw function.
func (s *randomSequence) Next() (time.Duration, bool) {
	if s.span == 0 {
		return s.minDelay, true
	}

	return saturatingDurationAdd(s.minDelay, randomDurationInclusive(s.random, s.span)), true
}
