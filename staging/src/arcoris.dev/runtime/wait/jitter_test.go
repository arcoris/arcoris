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
	"fmt"
	"math"
	"testing"
	"time"
)

// TestJitterReturnsOriginalDurationForZeroFactor verifies that factor zero is a
// valid disabled-jitter configuration.
func TestJitterReturnsOriginalDurationForZeroFactor(t *testing.T) {
	t.Parallel()

	if got := Jitter(time.Second, 0); got != time.Second {
		t.Fatalf("Jitter(1s, 0) = %v, want %v", got, time.Second)
	}
}

// TestJitterReturnsOriginalDurationForNonPositiveDurations verifies that jitter
// does not turn immediate waits into positive waits.
func TestJitterReturnsOriginalDurationForNonPositiveDurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "zero",
			duration: 0,
		},
		{
			name:     "negative",
			duration: -time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := Jitter(tc.duration, 1); got != tc.duration {
				t.Fatalf("Jitter(%v, 1) = %v, want %v", tc.duration, got, tc.duration)
			}
		})
	}
}

// TestJitterWithDrawUsesLowerBound verifies the deterministic lower jitter
// bound without depending on random output.
func TestJitterWithDrawUsesLowerBound(t *testing.T) {
	t.Parallel()

	got := jitterWithDraw(10*time.Second, 0.5, func(maxDelta time.Duration) time.Duration {
		if maxDelta != 5*time.Second {
			t.Fatalf("maxDelta = %v, want %v", maxDelta, 5*time.Second)
		}
		return 0
	})

	if got != 10*time.Second {
		t.Fatalf("jitterWithDraw lower bound = %v, want %v", got, 10*time.Second)
	}
}

// TestJitterWithDrawUsesUpperBound verifies the deterministic upper jitter
// bound without depending on random output.
func TestJitterWithDrawUsesUpperBound(t *testing.T) {
	t.Parallel()

	got := jitterWithDraw(10*time.Second, 0.5, func(maxDelta time.Duration) time.Duration {
		if maxDelta != 5*time.Second {
			t.Fatalf("maxDelta = %v, want %v", maxDelta, 5*time.Second)
		}
		return maxDelta
	})

	if got != 15*time.Second {
		t.Fatalf("jitterWithDraw upper bound = %v, want %v", got, 15*time.Second)
	}
}

// TestJitterWithDrawUsesIntermediateDelta verifies that jitter composition adds
// the selected delta to the base duration exactly once.
func TestJitterWithDrawUsesIntermediateDelta(t *testing.T) {
	t.Parallel()

	got := jitterWithDraw(10*time.Second, 0.5, func(maxDelta time.Duration) time.Duration {
		return maxDelta / 2
	})

	if got != 12500*time.Millisecond {
		t.Fatalf("jitterWithDraw intermediate = %v, want %v", got, 12500*time.Millisecond)
	}
}

// TestJitterReturnsValueWithinBounds verifies public random jitter bounds over
// repeated calls without asserting any specific random sequence.
func TestJitterReturnsValueWithinBounds(t *testing.T) {
	t.Parallel()

	base := 100 * time.Millisecond
	max := 125 * time.Millisecond

	for i := 0; i < 128; i++ {
		got := Jitter(base, 0.25)
		if got < base || got > max {
			t.Fatalf("Jitter(%v, 0.25) = %v, want value in [%v, %v]", base, got, base, max)
		}
	}
}

// TestMaxJitterDelta verifies deterministic maximum-delta calculation,
// rounding, and overflow saturation.
func TestMaxJitterDelta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
		factor   float64
		want     time.Duration
	}{
		{
			name:     "zero duration",
			duration: 0,
			factor:   1,
			want:     0,
		},
		{
			name:     "negative duration",
			duration: -time.Second,
			factor:   1,
			want:     0,
		},
		{
			name:     "zero factor",
			duration: time.Second,
			factor:   0,
			want:     0,
		},
		{
			name:     "fractional factor",
			duration: 10 * time.Second,
			factor:   0.25,
			want:     2500 * time.Millisecond,
		},
		{
			name:     "sub nanosecond delta rounds down",
			duration: time.Nanosecond,
			factor:   0.5,
			want:     0,
		},
		{
			name:     "saturates before duration overflow",
			duration: maxDuration - time.Second,
			factor:   2,
			want:     time.Second,
		},
		{
			name:     "max duration has no remaining headroom",
			duration: maxDuration,
			factor:   1,
			want:     0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := maxJitterDelta(tc.duration, tc.factor); got != tc.want {
				t.Fatalf("maxJitterDelta(%v, %v) = %v, want %v", tc.duration, tc.factor, got, tc.want)
			}
		})
	}
}

// TestJitterSaturatesAtMaxDuration verifies that public jitter never overflows
// when the requested maximum delta exceeds representable duration headroom.
func TestJitterSaturatesAtMaxDuration(t *testing.T) {
	t.Parallel()

	base := maxDuration - time.Second

	got := jitterWithDraw(base, 2, func(maxDelta time.Duration) time.Duration {
		if maxDelta != time.Second {
			t.Fatalf("maxDelta = %v, want %v", maxDelta, time.Second)
		}
		return maxDelta
	})

	if got != maxDuration {
		t.Fatalf("saturated jitter = %v, want %v", got, maxDuration)
	}
}

// TestJitterPanicsOnInvalidFactor verifies validation of caller-provided jitter
// factors.
func TestJitterPanicsOnInvalidFactor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		factor float64
		panic  any
	}{
		{
			name:   "negative",
			factor: -0.1,
			panic:  errNegativeJitterFactor,
		},
		{
			name:   "nan",
			factor: math.NaN(),
			panic:  errNonFiniteJitterFactor,
		},
		{
			name:   "positive infinity",
			factor: math.Inf(1),
			panic:  errNonFiniteJitterFactor,
		},
		{
			name:   "negative infinity",
			factor: math.Inf(-1),
			panic:  errNonFiniteJitterFactor,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustJitterPanicWith(t, tc.panic, func() {
				_ = Jitter(time.Second, tc.factor)
			})
		})
	}
}

// mustJitterPanicWith fails the test unless fn panics with want.
func mustJitterPanicWith(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		if got != want {
			t.Fatalf("panic = %s, want %s", fmt.Sprint(got), fmt.Sprint(want))
		}
	}()

	fn()
}
