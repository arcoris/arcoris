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

import (
	"math"
	"time"
)

const (
	// maxDuration is the largest representable time.Duration value used by
	// backoff duration arithmetic.
	//
	// Backoff schedules produce runtime delay values. Negative or wrapped
	// durations have no valid interpretation for retry loops, polling loops,
	// reconnect loops, controller cooldowns, timers, or clocks. Helpers in this
	// file saturate at maxDuration when arithmetic would overflow.
	maxDuration time.Duration = 1<<63 - 1

	// maxDurationFloat is maxDuration represented as a float64 nanosecond count.
	//
	// Some schedules, such as exponential backoff and proportional jitter, are
	// configured by floating-point multipliers. This constant provides one stable
	// upper bound for converting floating-point duration arithmetic back into
	// time.Duration values.
	maxDurationFloat = float64(maxDuration)
)

// saturatingDurationAdd returns left+right capped at maxDuration.
//
// The helper assumes callers are working with runtime delay values. Negative
// inputs are treated defensively: a negative side contributes nothing to the
// result. Public constructors and Sequence implementations should still reject
// or avoid negative delays before they reach this helper.
//
// Saturation is used instead of overflow. If the exact sum cannot fit in a
// positive time.Duration, maxDuration is returned.
func saturatingDurationAdd(left, right time.Duration) time.Duration {
	if left <= 0 {
		if right <= 0 {
			return 0
		}
		return right
	}
	if right <= 0 {
		return left
	}

	max := uint64(maxDuration)
	l := uint64(left)
	r := uint64(right)

	if r > max-l {
		return maxDuration
	}

	return time.Duration(l + r)
}

// saturatingDurationSub returns left-right with a lower bound of zero.
//
// The helper is useful for computing non-negative jitter spans and remaining
// capacity. If right is greater than or equal to left, the result is zero.
// Negative inputs are handled defensively and do not produce negative output.
//
// Public schedule constructors should still validate their own invariants
// before relying on this helper.
func saturatingDurationSub(left, right time.Duration) time.Duration {
	if left <= 0 {
		return 0
	}
	if right <= 0 {
		return left
	}
	if right >= left {
		return 0
	}

	return left - right
}

// saturatingDurationMul returns duration*factor capped at maxDuration.
//
// The helper assumes duration is a runtime delay value and factor is a count-like
// multiplier. If duration or factor is zero, the result is zero. If the exact
// product cannot fit in a positive time.Duration, maxDuration is returned.
func saturatingDurationMul(duration time.Duration, factor uint64) time.Duration {
	if duration <= 0 || factor == 0 {
		return 0
	}

	value := uint64(duration)
	max := uint64(maxDuration)

	if factor > max/value {
		return maxDuration
	}

	return time.Duration(value * factor)
}

// saturatingDurationMulFloat returns duration*factor capped at maxDuration.
//
// The helper is intended for algorithms configured by floating-point
// multipliers, such as exponential backoff and proportional jitter. Invalid
// floating-point values do not propagate into negative or wrapped durations:
//
//   - factor <= 0 returns zero;
//   - NaN returns zero;
//   - positive infinity saturates to maxDuration;
//   - finite values above the duration range saturate to maxDuration.
//
// Public constructors should validate user-facing factors before constructing
// schedules. This helper remains defensive so internal arithmetic stays total.
func saturatingDurationMulFloat(duration time.Duration, factor float64) time.Duration {
	if duration <= 0 || factor <= 0 || math.IsNaN(factor) {
		return 0
	}
	if math.IsInf(factor, 1) {
		return maxDuration
	}

	value := float64(duration) * factor
	return durationFromFloat(value)
}

// durationFromFloat converts a floating-point nanosecond count into
// time.Duration using saturating semantics.
//
// Values less than or equal to zero and NaN become zero. Positive infinity and
// values greater than or equal to maxDuration become maxDuration. Finite values
// inside the representable range are truncated toward zero, matching ordinary
// time.Duration conversion from floating-point nanoseconds.
func durationFromFloat(value float64) time.Duration {
	if value <= 0 || math.IsNaN(value) {
		return 0
	}
	if math.IsInf(value, 1) || value >= maxDurationFloat {
		return maxDuration
	}

	return time.Duration(value)
}

// capDuration returns delay capped to maxDelay.
//
// The helper assumes maxDelay is non-negative. Negative delay values are treated
// defensively as zero, but public Sequence implementations should reject
// negative available delays before calling capDuration when negative values
// indicate a child contract violation.
func capDuration(delay, maxDelay time.Duration) time.Duration {
	if delay <= 0 {
		return 0
	}
	if maxDelay <= 0 {
		return 0
	}
	if delay > maxDelay {
		return maxDelay
	}

	return delay
}

// minDuration returns the smaller of left and right.
//
// The helper is package-local to keep duration comparisons readable in schedule
// implementations without introducing a public utility API.
func minDuration(left, right time.Duration) time.Duration {
	if left <= right {
		return left
	}

	return right
}

// maxDurationValue returns the larger of left and right.
//
// The name avoids colliding with the maxDuration constant. The helper is
// package-local to keep duration comparisons readable in schedule
// implementations without introducing a public utility API.
func maxDurationValue(left, right time.Duration) time.Duration {
	if left >= right {
		return left
	}

	return right
}

// isNegativeDuration reports whether d is below zero.
//
// The helper exists for readability at Sequence contract boundaries where
// negative child delays should be treated as programming errors.
func isNegativeDuration(d time.Duration) bool {
	return d < 0
}

// isNonNegativeDuration reports whether d is zero or positive.
//
// The helper exists for validation code that needs to accept immediate delays
// while rejecting invalid negative runtime durations.
func isNonNegativeDuration(d time.Duration) bool {
	return d >= 0
}

// isPositiveDuration reports whether d is strictly positive.
//
// The helper exists for validation code used by schedules such as Exponential
// and Fibonacci, where zero would collapse the schedule into an immediate or
// non-growing sequence.
func isPositiveDuration(d time.Duration) bool {
	return d > 0
}
