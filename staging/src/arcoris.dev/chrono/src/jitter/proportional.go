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

// Proportional returns a schedule that randomizes each available child
// delay around the child delay by ratio.
//
// For a base delay d and ratio r, the returned delay is in:
//
//	[d * (1-r), d * (1+r)]
//
// Ratio must be in [0, 1]. This keeps the lower bound non-negative. A ratio of
// zero disables randomization and returns child delays unchanged. A ratio of one
// allows the lower bound to reach zero and the upper bound to reach roughly
// twice the base delay, subject to time.Duration saturation.
//
// Proportional jitter is useful when Full is too aggressive but positive
// one-sided jitter is too conservative. It desynchronizes owners while keeping
// returned values close to the child schedule's original cadence.
//
// For example:
//
//	jitter.Proportional(delay.Fixed(time.Second), 0.2)
//
// may produce values in:
//
//	[800*time.Millisecond, 1200*time.Millisecond]
//
// Proportional preserves child exhaustion. If the child sequence reports
// ok=false, the jittered sequence also reports ok=false.
//
// A child delay of zero remains zero. Negative child delays are invalid and are
// rejected by the shared jitter wrapper as Sequence contract violations.
//
// The returned Schedule is immutable and safe to reuse as long as the wrapped
// schedule is safe to reuse. Each call to NewSequence creates a fresh child
// sequence and wraps it in an independent jittered sequence.
//
// Proportional does not sleep, create timers, observe context
// cancellation, execute operations, classify errors, retry failed work, log,
// trace, export metrics, rate limit callers, schedule queue items, or make
// domain decisions.
//
// Proportional panics when schedule is nil, ratio is outside [0, 1], NaN,
// or infinite, or a random option is invalid.
func Proportional(sched delay.Schedule, r float64, opts ...RandomOption) delay.Schedule {
	requireJitterRatio(r)

	return newJitterSchedule(sched, proportionalJitterTransform(r), opts...)
}

// proportionalJitterTransform returns a transform that randomizes each base
// delay in [base*(1-ratio), base*(1+ratio)].
//
// The returned transform assumes ratio has already been validated by
// requireJitterRatio. The shared jitter wrapper guarantees that baseDelay is
// non-negative.
func proportionalJitterTransform(r float64) jitterTransform {
	return func(base time.Duration, random RandomGenerator) time.Duration {
		if base <= 0 || r == 0 {
			return base
		}

		delta := saturatingDurationMulFloat(base, r)

		lower := saturatingDurationSub(base, delta)
		upper := saturatingDurationAdd(base, delta)
		span := saturatingDurationSub(upper, lower)

		return saturatingDurationAdd(lower, randomOffsetInclusive(random, span))
	}
}
