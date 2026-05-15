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

package wait

import (
	"math/rand"
	"time"
)

const (
	// maxDuration is the largest representable time.Duration value.
	//
	// The constant is used by jitter arithmetic to saturate oversized positive
	// jitter instead of overflowing duration addition.
	maxDuration time.Duration = 1<<63 - 1
)

// Jitter returns duration plus a random positive jitter bounded by factor.
//
// Jitter is a low-level delay randomization primitive. It is intended to
// desynchronize otherwise identical wait loops so many components do not wake up
// at exactly the same interval boundary. It does not implement retry policy,
// backoff growth, retry limits, deadline ownership, metrics, tracing, logging,
// or scheduler policy.
//
// The factor argument is a fractional upper bound for additional delay. For
// example, Jitter(time.Second, 0.2) returns a value in the closed interval
// [1s, 1.2s], rounded to whole nanoseconds. A factor of zero returns duration
// unchanged.
//
// Non-positive durations are returned unchanged after factor validation. This
// preserves the immediate-delay semantics used by Delay and prevents jitter from
// turning an explicit immediate wait into a positive wait.
//
// Very large factors are saturated at the largest representable time.Duration.
// Saturation is preferable to overflow because jitter is a mechanical runtime
// primitive and must not produce negative or wrapped durations.
//
// Jitter uses the standard library pseudo-random generator. It is suitable for
// runtime desynchronization and load spreading, but it is not cryptographic
// randomness and MUST NOT be used for security decisions.
//
// Jitter panics when factor is negative, NaN, or infinite.
func Jitter(d time.Duration, factor float64) time.Duration {
	return jitterWithDraw(d, factor, randomJitterDelta)
}

// jitterWithDraw returns duration plus a jitter delta selected by draw.
//
// The helper owns the deterministic part of Jitter so boundary behavior can be
// tested without depending on random output. The draw function receives the
// maximum allowed delta and MUST return a value in the closed interval
// [0, maxDelta]. Public callers cannot provide draw; the invariant is enforced
// by the package-owned randomJitterDelta implementation.
func jitterWithDraw(
	d time.Duration,
	factor float64,
	draw func(maxDelta time.Duration) time.Duration,
) time.Duration {
	requireJitterFactor(factor)

	maxDelta := maxJitterDelta(d, factor)
	if maxDelta <= 0 {
		return d
	}

	return d + draw(maxDelta)
}

// randomJitterDelta returns a pseudo-random jitter delta in [0, maxDelta].
//
// The package-level math/rand functions are safe for concurrent use by multiple
// goroutines. The returned value is rounded to whole nanoseconds because
// time.Duration is an integer nanosecond count.
func randomJitterDelta(maxDelta time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(maxDelta) + 1))
}

// maxJitterDelta returns the largest positive delta that may be added to
// duration for factor.
//
// The result is rounded down to whole nanoseconds and capped so
// duration+result never overflows time.Duration. Non-positive durations and zero
// factors return zero because positive jitter is not meaningful for immediate or
// disabled waits.
func maxJitterDelta(d time.Duration, factor float64) time.Duration {
	if d <= 0 || factor == 0 {
		return 0
	}

	available := maxDuration - d
	if available <= 0 {
		return 0
	}

	requested := float64(d) * factor
	if requested <= 0 {
		return 0
	}
	if requested >= float64(available) {
		return available
	}

	return time.Duration(requested)
}
