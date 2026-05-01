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
	config := configOf(WithMaxAttempts(3))

	if config.maxAttempts != 3 {
		t.Fatalf("maxAttempts = %d, want 3", config.maxAttempts)
	}
}

func TestWithMaxAttemptsLastWins(t *testing.T) {
	config := configOf(
		WithMaxAttempts(2),
		WithMaxAttempts(5),
	)

	if config.maxAttempts != 5 {
		t.Fatalf("maxAttempts = %d, want 5", config.maxAttempts)
	}
}

func TestWithMaxAttemptsPanicsOnZero(t *testing.T) {
	expectPanic(t, panicZeroMaxAttempts, func() {
		_ = WithMaxAttempts(0)
	})
}

func TestWithMaxElapsed(t *testing.T) {
	config := configOf(WithMaxElapsed(5 * time.Second))

	if config.maxElapsed != 5*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", config.maxElapsed, 5*time.Second)
	}
}

func TestWithMaxElapsedAllowsZero(t *testing.T) {
	config := configOf(WithMaxElapsed(0))

	if config.maxElapsed != 0 {
		t.Fatalf("maxElapsed = %s, want 0", config.maxElapsed)
	}
}

func TestWithMaxElapsedLastWins(t *testing.T) {
	config := configOf(
		WithMaxElapsed(time.Second),
		WithMaxElapsed(2*time.Second),
	)

	if config.maxElapsed != 2*time.Second {
		t.Fatalf("maxElapsed = %s, want %s", config.maxElapsed, 2*time.Second)
	}
}

func TestWithMaxElapsedPanicsOnNegative(t *testing.T) {
	expectPanic(t, panicNegativeMaxElapsed, func() {
		_ = WithMaxElapsed(-time.Nanosecond)
	})
}
