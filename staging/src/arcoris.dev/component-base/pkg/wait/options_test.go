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
	"math"
	"testing"
	"time"
)

// TestDefaultOptionsPreserveBaseInterval verifies that the zero wait
// configuration keeps fixed intervals exact and does not add optional policy.
func TestDefaultOptionsPreserveBaseInterval(t *testing.T) {
	t.Parallel()

	config := defaultOptions()

	mustEqualDuration(t, "default interval", config.interval(time.Second), time.Second)
}

// TestOptionsOfWithoutOptionsReturnsDefault verifies option normalization for
// callers that do not request optional wait behavior.
func TestOptionsOfWithoutOptionsReturnsDefault(t *testing.T) {
	t.Parallel()

	config := optionsOf()

	mustEqualDuration(t, "normalized default interval", config.interval(time.Second), time.Second)
}

// TestWithJitterStoresFactor verifies that WithJitter updates only the jitter
// domain of the private normalized configuration.
func TestWithJitterStoresFactor(t *testing.T) {
	t.Parallel()

	config := optionsOf(WithJitter(0.25))

	if config.jitterFactor != 0.25 {
		t.Fatalf("jitterFactor = %v, want 0.25", config.jitterFactor)
	}
}

// TestWithJitterZeroDisablesJitter verifies that an explicit zero jitter option
// is valid and preserves the base interval exactly.
func TestWithJitterZeroDisablesJitter(t *testing.T) {
	t.Parallel()

	config := optionsOf(WithJitter(0))

	mustEqualDuration(t, "zero-jitter interval", config.interval(time.Second), time.Second)
}

// TestOptionsApplyInOrder verifies deterministic last-option-wins behavior for
// option domains configured more than once.
func TestOptionsApplyInOrder(t *testing.T) {
	t.Parallel()

	config := optionsOf(
		WithJitter(0.50),
		WithJitter(0.25),
	)

	if config.jitterFactor != 0.25 {
		t.Fatalf("jitterFactor = %v, want last configured value 0.25", config.jitterFactor)
	}
}

// TestOptionsLaterJitterCanDisableEarlierJitter verifies that last-option-wins
// semantics can intentionally return a configuration to the exact-interval mode.
func TestOptionsLaterJitterCanDisableEarlierJitter(t *testing.T) {
	t.Parallel()

	config := optionsOf(
		WithJitter(0.50),
		WithJitter(0),
	)

	mustEqualDuration(t, "disabled jitter interval", config.interval(time.Second), time.Second)
}

// TestOptionsIntervalAppliesJitterWithinBounds verifies that normalized options
// delegate interval spreading to the package jitter primitive.
func TestOptionsIntervalAppliesJitterWithinBounds(t *testing.T) {
	t.Parallel()

	base := time.Second
	factor := 0.25
	config := optionsOf(WithJitter(factor))

	got := config.interval(base)
	min := base
	max := base + maxJitterDelta(base, factor)

	if got < min || got > max {
		t.Fatalf("jittered interval = %v, want in [%v, %v]", got, min, max)
	}
}

// TestOptionsOfPanicsOnNilOption verifies that invalid option lists fail before
// a wait primitive evaluates conditions or starts sleeping.
func TestOptionsOfPanicsOnNilOption(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilOption, func() {
		_ = optionsOf(nil)
	})
}

// TestWithJitterPanicsOnInvalidFactor verifies constructor-time validation for
// jitter factors supplied through options.
func TestWithJitterPanicsOnInvalidFactor(t *testing.T) {
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
				_ = WithJitter(tt.factor)
			})
		})
	}
}
