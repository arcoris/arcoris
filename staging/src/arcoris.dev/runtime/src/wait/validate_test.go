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

package wait

import (
	"context"
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

// TestRequireContextAcceptsNonNilContext verifies the valid context path.
func TestRequireContextAcceptsNonNilContext(t *testing.T) {
	t.Parallel()

	requireContext(context.Background())
}

// TestRequireContextPanicsOnNilContext verifies nil-context validation.
func TestRequireContextPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	panicassert.RequireValue(t, errNilContext, func() {
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

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			panicassert.RequireValue(t, errNonPositiveInterval, func() {
				requirePositiveInterval(tc.interval)
			})
		})
	}
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

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			panicassert.RequireValue(t, errNilTimer, func() {
				tc.timer.requireUsable()
			})
		})
	}
}
