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
	"testing"
	"time"
)

// TestDefaultOptionsPreserveBaseInterval verifies that the zero wait
// configuration keeps fixed intervals exact and does not add optional policy.
func TestDefaultOptionsPreserveBaseInterval(t *testing.T) {
	t.Parallel()

	cfg := defaultOptions()

	mustEqualDuration(t, "default interval", cfg.interval(time.Second), time.Second)
}

// TestOptionsOfWithoutOptionsReturnsDefault verifies option normalization for
// callers that do not request optional wait behavior.
func TestOptionsOfWithoutOptionsReturnsDefault(t *testing.T) {
	t.Parallel()

	cfg := optionsOf()

	mustEqualDuration(t, "normalized default interval", cfg.interval(time.Second), time.Second)
}

// TestOptionsApplyInOrder verifies deterministic last-option-wins behavior for
// option domains configured more than once.
func TestOptionsApplyInOrder(t *testing.T) {
	t.Parallel()

	cfg := optionsOf(
		WithJitter(0.50),
		WithJitter(0.25),
	)

	if cfg.jitterFactor != 0.25 {
		t.Fatalf("jitterFactor = %v, want last configured value 0.25", cfg.jitterFactor)
	}
}

// TestOptionsLaterJitterCanDisableEarlierJitter verifies that last-option-wins
// semantics can intentionally return a configuration to the exact-interval mode.
func TestOptionsLaterJitterCanDisableEarlierJitter(t *testing.T) {
	t.Parallel()

	cfg := optionsOf(
		WithJitter(0.50),
		WithJitter(0),
	)

	mustEqualDuration(t, "disabled jitter interval", cfg.interval(time.Second), time.Second)
}

// TestOptionsIntervalAppliesJitterWithinBounds verifies that normalized options
// delegate interval spreading to the package jitter primitive.
func TestOptionsIntervalAppliesJitterWithinBounds(t *testing.T) {
	t.Parallel()

	base := time.Second
	factor := 0.25
	cfg := optionsOf(WithJitter(factor))

	got := cfg.interval(base)
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
