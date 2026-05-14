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

import (
	"math"
	"time"
)

const (
	// maxDuration is the largest representable non-negative time.Duration.
	//
	// Package delay uses this value for deterministic duration arithmetic that
	// must never wrap into a negative delay. Randomization-specific arithmetic
	// lives in package jitter.
	maxDuration time.Duration = 1<<63 - 1

	// maxDurationFloat is maxDuration represented as a float64 nanosecond count.
	//
	// Exponential keeps floating-point mathematical state so small durations with
	// fractional multipliers can still grow after integer truncation. This
	// constant provides the stable saturation boundary for that conversion.
	maxDurationFloat = float64(maxDuration)
)

// saturatingDurationAdd returns l+r capped at maxDuration.
//
// The helper is intentionally defensive: non-positive inputs do not make the
// result negative, and overflow saturates. Public constructors and Sequence
// wrappers still validate their own contracts before depending on this helper.
func saturatingDurationAdd(l, r time.Duration) time.Duration {
	if l <= 0 {
		if r <= 0 {
			return 0
		}
		return r
	}
	if r <= 0 {
		return l
	}

	m := uint64(maxDuration)
	lu := uint64(l)
	ru := uint64(r)
	if ru > m-lu {
		return maxDuration
	}

	return time.Duration(lu + ru)
}

// saturatingDurationSub returns l-r with a lower bound of zero.
//
// The helper is useful for deterministic range and remaining-capacity arithmetic
// where negative durations have no runtime meaning. Negative inputs are handled
// defensively so callers receive a non-negative result.
func saturatingDurationSub(l, r time.Duration) time.Duration {
	if l <= 0 {
		return 0
	}
	if r <= 0 {
		return l
	}
	if r >= l {
		return 0
	}

	return l - r
}

// saturatingDurationMul returns d*n capped at maxDuration.
//
// Linear uses this helper to compute step*index without allowing a large index
// to overflow. Constructors validate schedule inputs before sequences call this
// helper, but the helper remains defensive for package-local arithmetic.
func saturatingDurationMul(d time.Duration, n uint64) time.Duration {
	if d <= 0 || n == 0 {
		return 0
	}

	v := uint64(d)
	m := uint64(maxDuration)
	if n > m/v {
		return maxDuration
	}

	return time.Duration(v * n)
}

// durationFromFloat converts a floating-point nanosecond count into a saturated
// time.Duration.
//
// Exponential calls this when exposing its mathematical state as a concrete
// delay value. NaN and non-positive values collapse to zero, while positive
// infinity and values beyond the duration range saturate at maxDuration.
func durationFromFloat(v float64) time.Duration {
	if v <= 0 || math.IsNaN(v) {
		return 0
	}
	if math.IsInf(v, 1) || v >= maxDurationFloat {
		return maxDuration
	}

	return time.Duration(v)
}

// capDuration returns d capped to m.
//
// The helper assumes m was accepted by a public constructor and is therefore
// non-negative. Negative delay values are still treated defensively as zero, but
// wrapper code should reject negative available child delays before calling
// capDuration because they violate the Sequence contract.
func capDuration(d, m time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	if m <= 0 {
		return 0
	}
	if d > m {
		return m
	}

	return d
}

// isNegativeDuration reports whether d is below zero.
//
// The helper exists for tests and validation code that need to name the
// Sequence contract boundary explicitly: available delay values must never be
// negative.
func isNegativeDuration(d time.Duration) bool {
	return d < 0
}

// isNonNegativeDuration reports whether d is zero or positive.
//
// Zero is a valid delay value in package delay and means immediate continuation.
// This helper names that common validation predicate without exporting utility
// API.
func isNonNegativeDuration(d time.Duration) bool {
	return d >= 0
}
