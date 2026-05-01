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

// EqualJitter returns a schedule that randomizes each available child delay in
// the upper half of the child delay range.
//
// Equal jitter keeps part of the child delay as a minimum cooldown and
// randomizes the remaining part. For a base delay d, the returned delay is in:
//
//	[d/2, d]
//
// More precisely, for integer nanosecond durations, the lower bound is d/2
// rounded down, and the random span is d-d/2. This keeps the original base delay
// reachable even when d is an odd number of nanoseconds.
//
// EqualJitter is useful when FullJitter is too aggressive because it may return
// very small delays. EqualJitter still desynchronizes owners, but it preserves a
// minimum delay floor for each child value.
//
// For example, if the child schedule produces:
//
//	100*time.Millisecond
//	200*time.Millisecond
//	400*time.Millisecond
//
// EqualJitter may produce values in:
//
//	[50*time.Millisecond, 100*time.Millisecond]
//	[100*time.Millisecond, 200*time.Millisecond]
//	[200*time.Millisecond, 400*time.Millisecond]
//
// EqualJitter preserves child exhaustion. If the child sequence reports
// ok=false, the jittered sequence also reports ok=false. EqualJitter does not
// make a finite child schedule infinite and does not make an infinite child
// schedule finite.
//
// A child delay of zero remains zero. Negative child delays are invalid and are
// rejected by the shared jitter wrapper as Sequence contract violations.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates a fresh child
// sequence and wraps it in an independent jittered sequence.
//
// EqualJitter does not sleep, create timers, observe context cancellation,
// execute operations, classify errors, retry failed work, log, trace, export
// metrics, rate limit callers, schedule queue items, or make domain decisions.
//
// EqualJitter panics when schedule is nil or a random option is invalid.
func EqualJitter(schedule Schedule, opts ...RandomOption) Schedule {
	return newJitterSchedule(schedule, equalJitterTransform, opts...)
}

// equalJitterTransform applies equal jitter to one available base delay.
//
// The transform returns a random delay in [baseDelay/2, baseDelay]. The base
// delay is guaranteed non-negative by the shared jitter wrapper.
func equalJitterTransform(baseDelay time.Duration, random RandomGenerator) time.Duration {
	if baseDelay <= 0 {
		return 0
	}

	lower := baseDelay / 2
	span := saturatingDurationSub(baseDelay, lower)

	return saturatingDurationAdd(lower, randomOffsetInclusive(random, span))
}
