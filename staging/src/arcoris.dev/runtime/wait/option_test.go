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

// TestWithJitterStoresFactor verifies that WithJitter updates only the jitter
// domain of the private normalized configuration.
func TestWithJitterStoresFactor(t *testing.T) {
	t.Parallel()

	cfg := optionsOf(WithJitter(0.25))

	if cfg.jitterFactor != 0.25 {
		t.Fatalf("jitterFactor = %v, want 0.25", cfg.jitterFactor)
	}
}

// TestWithJitterZeroDisablesJitter verifies that an explicit zero jitter option
// is valid and preserves the base interval exactly.
func TestWithJitterZeroDisablesJitter(t *testing.T) {
	t.Parallel()

	cfg := optionsOf(WithJitter(0))

	mustEqualDuration(t, "zero-jitter interval", cfg.interval(time.Second), time.Second)
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tc.panic, func() {
				_ = WithJitter(tc.factor)
			})
		})
	}
}
