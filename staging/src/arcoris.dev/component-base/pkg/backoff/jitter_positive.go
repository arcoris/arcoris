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

// PositiveJitter returns a schedule that only extends each available child delay.
//
// Positive jitter keeps the child delay as the lower bound and adds a random
// non-negative delta bounded by factor. For a base delay d and factor f, the
// returned delay is in:
//
//	[d, d + d*f]
//
// The final value is capped by time.Duration's maximum representable value if
// the computed upper bound would overflow.
//
// PositiveJitter is useful for conservative polling and cooldown paths where
// waking earlier than the base cadence would violate caller policy. It spreads
// wake-ups without shortening the original delay.
//
// For example:
//
//	backoff.PositiveJitter(backoff.Fixed(time.Second), 0.2)
//
// may produce values in:
//
//	[1*time.Second, 1200*time.Millisecond]
//
// PositiveJitter is different from FullJitter and EqualJitter. It never shortens
// a positive child delay. That makes it safer for conservative intervals, but it
// may increase tail latency and is usually not the best default for high
// contention retry overload protection.
//
// PositiveJitter preserves child exhaustion. If the child sequence reports
// ok=false, the jittered sequence also reports ok=false.
//
// A child delay of zero remains zero. Negative child delays are invalid and are
// rejected by the shared jitter wrapper as Sequence contract violations.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates a fresh child
// sequence and wraps it in an independent jittered sequence.
//
// PositiveJitter does not sleep, create timers, observe context cancellation,
// execute operations, classify errors, retry failed work, log, trace, export
// metrics, rate limit callers, schedule queue items, or make domain decisions.
//
// PositiveJitter panics when schedule is nil, factor is negative, NaN, or
// infinite, or a random option is invalid. A factor of zero is valid and makes
// the wrapper return child delays unchanged.
func PositiveJitter(schedule Schedule, factor float64, opts ...RandomOption) Schedule {
	requireJitterFactor(factor)

	return newJitterSchedule(schedule, positiveJitterTransform(factor), opts...)
}

// positiveJitterTransform returns a transform that extends each base delay by a
// random non-negative delta bounded by factor.
//
// The returned transform assumes factor has already been validated by
// requireJitterFactor. The shared jitter wrapper guarantees that baseDelay is
// non-negative.
func positiveJitterTransform(factor float64) jitterTransform {
	return func(baseDelay time.Duration, random RandomGenerator) time.Duration {
		if baseDelay <= 0 || factor == 0 {
			return baseDelay
		}

		maxDelta := saturatingDurationMulFloat(baseDelay, factor)
		return saturatingDurationAdd(baseDelay, randomOffsetInclusive(random, maxDelta))
	}
}
