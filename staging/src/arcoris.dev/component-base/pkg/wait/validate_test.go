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
	"context"
	"math"
	"testing"
	"time"
)

// TestRequireContextAcceptsNonNilContext verifies the valid context path.
func TestRequireContextAcceptsNonNilContext(t *testing.T) {
	t.Parallel()

	requireContext(context.Background())
}

// TestRequireContextPanicsOnNilContext verifies nil-context validation.
func TestRequireContextPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilContext, func() {
		requireContext(nil)
	})
}

// TestRequirePositiveIntervalAcceptsPositiveDuration verifies valid fixed-loop
// interval validation.
func TestRequirePositiveIntervalAcceptsPositiveDuration(t *testing.T) {
	t.Parallel()

	requirePositiveInterval(time.Nanosecond)
}

// TestRequirePositiveIntervalPanicsOnNonPositiveDuration verifies that fixed
// cadence loops reject busy-loop intervals.
func TestRequirePositiveIntervalPanicsOnNonPositiveDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
	}{
		{
			name:     "zero",
			interval: 0,
		},
		{
			name:     "negative",
			interval: -time.Nanosecond,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNonPositiveInterval, func() {
				requirePositiveInterval(tt.interval)
			})
		})
	}
}

// TestRequireJitterFactorAcceptsValidFactors verifies valid positive-jitter
// factor validation.
func TestRequireJitterFactorAcceptsValidFactors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		factor float64
	}{
		{
			name:   "zero",
			factor: 0,
		},
		{
			name:   "fraction",
			factor: 0.25,
		},
		{
			name:   "large finite",
			factor: math.MaxFloat64,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			requireJitterFactor(tt.factor)
		})
	}
}

// TestRequireJitterFactorPanicsOnInvalidFactors verifies invalid jitter-factor
// validation.
func TestRequireJitterFactorPanicsOnInvalidFactors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		factor float64
		panic  string
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

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tt.panic, func() {
				requireJitterFactor(tt.factor)
			})
		})
	}
}

// TestRequireOptionAcceptsNonNilOption verifies valid option validation.
func TestRequireOptionAcceptsNonNilOption(t *testing.T) {
	t.Parallel()

	requireOption(WithJitter(0))
}

// TestRequireOptionPanicsOnNilOption verifies nil-option validation.
func TestRequireOptionPanicsOnNilOption(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilOption, func() {
		requireOption(nil)
	})
}

// TestTimerRequireUsableAcceptsNewTimer verifies valid Timer receiver
// validation.
func TestTimerRequireUsableAcceptsNewTimer(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	defer timer.StopAndDrain()

	timer.requireUsable()
}

// TestTimerRequireUsablePanicsOnInvalidTimer verifies nil and zero-value Timer
// receiver validation.
func TestTimerRequireUsablePanicsOnInvalidTimer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		timer *Timer
	}{
		{
			name:  "nil receiver",
			timer: nil,
		},
		{
			name:  "zero value",
			timer: &Timer{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilTimer, func() {
				tt.timer.requireUsable()
			})
		})
	}
}
