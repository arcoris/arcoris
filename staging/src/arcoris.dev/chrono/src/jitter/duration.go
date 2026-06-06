// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jitter

import (
	"math"
	"time"
)

const (
	// maxDuration is the largest representable non-negative time.Duration.
	//
	// Jitter transforms saturate at this value instead of wrapping randomized
	// upper bounds into invalid negative durations.
	maxDuration time.Duration = 1<<63 - 1

	// maxDurationFloat is maxDuration represented as a float64 nanosecond count.
	//
	// Factor and ratio based jitter algorithms use floating-point arithmetic for
	// caller-facing configuration. This constant provides the stable conversion
	// boundary back into time.Duration.
	maxDurationFloat = float64(maxDuration)

	// int63Domain is the number of distinct values a valid Int63 draw can
	// produce.
	int63Domain = uint64(1) << 63
)

// saturatingDurationAdd returns l+r capped at maxDuration.
//
// Positive and proportional jitter use this helper when adding randomized
// offsets to a lower bound. Non-positive inputs are treated defensively so the
// result remains a valid non-negative delay.
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
// Jitter algorithms use this helper to derive non-negative spans between lower
// and upper bounds. If right is greater than left, the span collapses to zero
// instead of producing an invalid negative duration.
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

// saturatingDurationMulFloat returns d*f capped at maxDuration.
//
// Positive, proportional, and decorrelated jitter use this helper for
// caller-configured factors and ratios. Invalid floating-point values are mapped
// to zero or maxDuration defensively; public constructors still reject invalid
// user input before schedules are created.
func saturatingDurationMulFloat(d time.Duration, f float64) time.Duration {
	if d <= 0 || f <= 0 || math.IsNaN(f) {
		return 0
	}
	if math.IsInf(f, 1) {
		return maxDuration
	}

	v := float64(d) * f
	if v <= 0 || math.IsNaN(v) {
		return 0
	}
	if math.IsInf(v, 1) || v >= maxDurationFloat {
		return maxDuration
	}

	return time.Duration(v)
}

// minDuration returns the smaller duration.
//
// Decorrelated jitter uses this helper to apply maxDelay to a computed upper
// bound without introducing a broad public utility API.
func minDuration(l, r time.Duration) time.Duration {
	if l <= r {
		return l
	}

	return r
}

// randomUint63 returns a validated Int63 draw as uint64.
//
// RandomGenerator implementations must match math/rand.Int63 and return values
// in [0, 1<<63). Negative values are rejected instead of being converted through
// uint64 because that would hide a broken custom generator.
func randomUint63(r RandomGenerator) uint64 {
	requireRandom(r, errNilRandom)

	v := r.Int63()
	if v < 0 {
		panic(errRandomReturnedNegative)
	}

	return uint64(v)
}

// randomDurationInclusive returns a pseudo-random duration in [0, maxOffset].
//
// The helper maps RandomGenerator.Int63 values to a closed integer-nanosecond
// duration range. It uses rejection sampling to avoid modulo bias while
// remaining non-cryptographic and deterministic for caller-supplied generators.
// The caller owns range validation; maxOffset <= 0 returns zero without drawing.
func randomDurationInclusive(r RandomGenerator, maxOffset time.Duration) time.Duration {
	requireRandom(r, errNilRandom)
	if maxOffset <= 0 {
		return 0
	}

	bound := uint64(maxOffset) + 1
	limit := int63Domain - int63Domain%bound

	for {
		v := randomUint63(r)
		if v < limit {
			return time.Duration(v % bound)
		}
	}
}

// randomOffsetInclusive returns a pseudo-random offset in [0, maxOffset].
//
// The helper is a semantic alias for randomDurationInclusive used by transforms
// that add an offset to a non-zero lower bound.
func randomOffsetInclusive(r RandomGenerator, maxOffset time.Duration) time.Duration {
	return randomDurationInclusive(r, maxOffset)
}
