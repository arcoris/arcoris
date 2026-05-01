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

// FullJitter returns a schedule that randomizes each available child delay in
// the range [0, baseDelay].
//
// Full jitter is useful for distributed retry and polling paths where many
// owners may otherwise wake up at the same deterministic delay boundaries. It
// spreads each child delay across the full range from immediate continuation to
// the child-provided delay.
//
// For example, if the child schedule produces:
//
//	100*time.Millisecond
//	200*time.Millisecond
//	400*time.Millisecond
//
// FullJitter may produce values in:
//
//	[0, 100*time.Millisecond]
//	[0, 200*time.Millisecond]
//	[0, 400*time.Millisecond]
//
// FullJitter is commonly useful with capped exponential schedules:
//
//	backoff.FullJitter(
//	    backoff.Cap(
//	        backoff.Exponential(100*time.Millisecond, 2.0),
//	        2*time.Second,
//	    ),
//	)
//
// FullJitter preserves child exhaustion. If the child sequence reports ok=false,
// the jittered sequence also reports ok=false. FullJitter does not make a finite
// child schedule infinite and does not make an infinite child schedule finite.
//
// A child delay of zero remains zero. Negative child delays are invalid and are
// rejected by the shared jitter wrapper as Sequence contract violations.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates a fresh child
// sequence and wraps it in an independent jittered sequence.
//
// FullJitter does not sleep, create timers, observe context cancellation,
// execute operations, classify errors, retry failed work, log, trace, export
// metrics, rate limit callers, schedule queue items, or make domain decisions.
//
// FullJitter panics when schedule is nil or a random option is invalid.
func FullJitter(schedule Schedule, opts ...RandomOption) Schedule {
	return newJitterSchedule(schedule, fullJitterTransform, opts...)
}

// fullJitterTransform applies full jitter to one available base delay.
//
// The transform returns a random delay in [0, baseDelay]. The base delay is
// guaranteed non-negative by the shared jitter wrapper.
func fullJitterTransform(baseDelay time.Duration, random RandomGenerator) time.Duration {
	if baseDelay <= 0 {
		return 0
	}

	return randomDurationInclusive(random, baseDelay)
}
