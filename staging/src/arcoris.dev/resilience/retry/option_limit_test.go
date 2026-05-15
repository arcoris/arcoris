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

package retry

import (
	"testing"
	"time"
)

func TestWithMaxAttempts(t *testing.T) {
	cfg := configOf(WithMaxAttempts(3))

	if cfg.maxAttempts != 3 {
		t.Fatalf("maxAttempts = %d, want 3", cfg.maxAttempts)
	}
}

func TestWithMaxAttemptsLastWins(t *testing.T) {
	cfg := configOf(
		WithMaxAttempts(2),
		WithMaxAttempts(5),
	)

	if cfg.maxAttempts != 5 {
		t.Fatalf("maxAttempts = %d, want 5", cfg.maxAttempts)
	}
}

func TestWithMaxAttemptsPanicsOnZero(t *testing.T) {
	expectPanic(t, panicZeroMaxAttempts, func() {
		_ = WithMaxAttempts(0)
	})
}

func TestWithMaxElapsed(t *testing.T) {
	cfg := configOf(WithMaxElapsed(5 * time.Second))

	if cfg.maxElapsed != 5*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", cfg.maxElapsed, 5*time.Second)
	}
}

func TestWithMaxElapsedAllowsZero(t *testing.T) {
	cfg := configOf(WithMaxElapsed(0))

	if cfg.maxElapsed != 0 {
		t.Fatalf("maxElapsed = %s, want 0", cfg.maxElapsed)
	}
}

func TestWithMaxElapsedLastWins(t *testing.T) {
	cfg := configOf(
		WithMaxElapsed(time.Second),
		WithMaxElapsed(2*time.Second),
	)

	if cfg.maxElapsed != 2*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", cfg.maxElapsed, 2*time.Second)
	}
}

func TestWithMaxElapsedPanicsOnNegative(t *testing.T) {
	expectPanic(t, panicNegativeMaxElapsed, func() {
		_ = WithMaxElapsed(-time.Nanosecond)
	})
}
