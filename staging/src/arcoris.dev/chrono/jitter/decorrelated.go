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

package jitter

import (
	"time"

	"arcoris.dev/chrono/delay"
)

const (
	// errNonPositiveDecorrelatedInitialDelay is the stable diagnostic text used
	// when Decorrelated receives a non-positive initial delay.
	//
	// Decorrelated jitter schedules model positive runtime durations. A zero or
	// negative initial delay cannot provide a meaningful lower bound for the
	// random delay range. Callers that need an immediate first step should
	// compose an explicit zero delay before a decorrelated jitter schedule.
	errNonPositiveDecorrelatedInitialDelay = "jitter: non-positive decorrelated initial delay"

	// errDecorrelatedMaxDelayBeforeInitialDelay is the stable diagnostic text
	// used when Decorrelated receives a maximum delay smaller than the
	// initial delay.
	//
	// Decorrelated jitter draws delays from ranges whose lower bound is the
	// initial delay. A maximum delay smaller than that lower bound cannot produce
	// a valid non-empty delay range and is rejected at the schedule boundary.
	errDecorrelatedMaxDelayBeforeInitialDelay = "jitter: decorrelated maximum delay before initial delay"

	// errInvalidDecorrelatedMultiplier is the stable diagnostic text used when
	// Decorrelated receives an invalid multiplier.
	//
	// The multiplier must be finite and greater than one. A multiplier less than
	// or equal to one does not provide growth room above the previous delay. NaN
	// and infinite values cannot produce a stable runtime delay sequence.
	errInvalidDecorrelatedMultiplier = "jitter: invalid decorrelated multiplier"

	// errNilDecorrelatedSource is the stable diagnostic text used when the
	// package-local decorrelated jitter constructor receives a nil RandomSource.
	errNilDecorrelatedSource = "jitter: nil decorrelated source"
)

// Decorrelated returns an infinite decorrelated-jitter schedule.
//
// Decorrelated jitter is a stateful randomized delay algorithm. Each produced
// delay influences the range used for the next delay. For initial delay b,
// maximum delay cap, multiplier m, and previous delay p, the next delay is:
//
//	random(b, min(cap, p*m))
//
// The first previous delay is initial, so the first returned delay is in:
//
//	[initial, min(maxDelay, initial*multiplier)]
//
// Decorrelated is useful for advanced retry and reconnect paths where
// deterministic exponential steps are likely to synchronize many owners. It
// keeps a stable lower bound while allowing the upper bound to move based on the
// previously selected delay.
//
// Decorrelated is not a wrapper around another Schedule. It owns its own
// sequence state because its algorithm depends on the previously returned delay.
//
// Decorrelated does not sleep, create timers, observe context
// cancellation, execute operations, classify errors, retry failed work, log,
// trace, export metrics, rate limit callers, schedule queue items, or make
// domain decisions.
//
// The returned Schedule is immutable and safe to reuse. Each call to NewSequence
// returns an independent Sequence with its own previous-delay state.
//
// Decorrelated panics when initial is not positive, maxDelay is smaller
// than initial, or multiplier is not finite and greater than one.
//
// If previous*multiplier would overflow time.Duration, the upper bound saturates
// at maxDelay. Returned delays never exceed maxDelay.
func Decorrelated(
	initial time.Duration,
	m time.Duration,
	f float64,
	opts ...RandomOption,
) delay.Schedule {
	config := randomOptionsOf(opts...)
	return decorrelatedJitterWithSource(initial, m, f, config.source)
}

// decorrelatedJitterWithSource returns a decorrelated-jitter schedule using
// source for per-sequence random values.
func decorrelatedJitterWithSource(
	initial time.Duration,
	m time.Duration,
	f float64,
	src RandomSource,
) delay.Schedule {
	requirePositiveDuration(initial, errNonPositiveDecorrelatedInitialDelay)
	requireDurationNotBefore(m, initial, errDecorrelatedMaxDelayBeforeInitialDelay)
	requireFloatGreaterThanOne(f, errInvalidDecorrelatedMultiplier)
	requireRandomSource(src, errNilDecorrelatedSource)

	return decorrelatedJitterSchedule{
		initial:    initial,
		maxDelay:   m,
		multiplier: f,
		source:     src,
	}
}

// decorrelatedJitterSchedule is the reusable recipe behind Decorrelated.
//
// The type stores only validated configuration and the random draw function. It
// does not store per-owner previous-delay state. Per-owner iteration is
// represented by decorrelatedJitterSequence values returned from NewSequence.
type decorrelatedJitterSchedule struct {
	// initial is the positive lower bound for every random delay range.
	initial time.Duration

	// maxDelay is the inclusive upper bound for every returned delay.
	maxDelay time.Duration

	// multiplier controls how far above the previous delay the next random range
	// may extend before maxDelay is applied.
	multiplier float64

	// source creates the per-sequence random generator used for offset draws.
	source RandomSource
}

// NewSequence returns an independent decorrelated-jitter delay sequence.
//
// The returned sequence starts with previous=initial. Multiple sequences created
// from the same schedule advance independently because each one owns its own
// previous-delay state.
func (s decorrelatedJitterSchedule) NewSequence() delay.Sequence {
	random := s.source.NewRandom()
	requireRandom(random, errNilRandom)

	return &decorrelatedJitterSequence{
		initial:    s.initial,
		maxDelay:   s.maxDelay,
		multiplier: s.multiplier,
		previous:   s.initial,
		random:     random,
	}
}

// decorrelatedJitterSequence is the per-owner delay stream produced by
// Decorrelated.
//
// The sequence stores the previous returned delay because decorrelated jitter
// uses it to compute the next upper bound.
type decorrelatedJitterSequence struct {
	// initial is the positive lower bound for every random delay range.
	initial time.Duration

	// maxDelay is the inclusive upper bound for every returned delay.
	maxDelay time.Duration

	// multiplier controls the upper range expansion from previous.
	multiplier float64

	// previous is the previously returned delay.
	//
	// The first previous value is initial. After each Next call, previous is set
	// to the returned delay.
	previous time.Duration

	// random is owned by this sequence and used for inclusive offset draws.
	random RandomGenerator
}

// Next returns the next decorrelated-jitter delay and reports that the sequence
// is still available.
//
// Decorrelated-jitter sequences are intentionally infinite. Exhaustion, if
// required, should be provided by a higher-level owner such as retry max attempts
// or by a finite wrapper such as Limit.
func (s *decorrelatedJitterSequence) Next() (time.Duration, bool) {
	upper := s.upperBound()
	span := saturatingDurationSub(upper, s.initial)

	delay := s.initial
	if span > 0 {
		delay = saturatingDurationAdd(delay, randomOffsetInclusive(s.random, span))
	}

	s.previous = delay

	return delay, true
}

// upperBound returns min(maxDelay, previous*multiplier) using saturating
// arithmetic.
//
// The result is always in [initial, maxDelay]. Constructor validation guarantees
// initial <= maxDelay and multiplier > 1.
func (s *decorrelatedJitterSequence) upperBound() time.Duration {
	scaled := saturatingDurationMulFloat(s.previous, s.multiplier)
	if scaled < s.initial {
		return s.initial
	}

	return minDuration(scaled, s.maxDelay)
}
