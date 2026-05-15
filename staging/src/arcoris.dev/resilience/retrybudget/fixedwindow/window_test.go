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

package fixedwindow

import (
	"testing"
	"time"
)

func TestLimiterRotateLocked(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithWindow(time.Minute), WithMinRetries(0))

	limiter.mu.Lock()
	limiter.original = 3
	limiter.retries = 2
	start := limiter.windowStart

	if got := limiter.rotateLocked(start.Add(30 * time.Second)); got {
		t.Fatal("rotate before window end returned true")
	}
	if limiter.original != 3 || limiter.retries != 2 || !limiter.windowStart.Equal(start) {
		t.Fatalf("state changed before window end: original=%d retries=%d start=%s", limiter.original, limiter.retries, limiter.windowStart)
	}

	if got := limiter.rotateLocked(start.Add(-time.Second)); got {
		t.Fatal("rotate with backwards time returned true")
	}
	if limiter.original != 3 || limiter.retries != 2 || !limiter.windowStart.Equal(start) {
		t.Fatalf("state changed on backwards time: original=%d retries=%d start=%s", limiter.original, limiter.retries, limiter.windowStart)
	}

	next := start.Add(time.Minute)
	if got := limiter.rotateLocked(next); !got {
		t.Fatal("rotate at exact window end returned false")
	}
	if limiter.original != 0 || limiter.retries != 0 || !limiter.windowStart.Equal(next) {
		t.Fatalf("state after rotation: original=%d retries=%d start=%s", limiter.original, limiter.retries, limiter.windowStart)
	}
	limiter.mu.Unlock()
}

func TestLimiterWindowEndLocked(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithWindow(45*time.Second))

	limiter.mu.Lock()
	got := limiter.windowEndLocked()
	want := limiter.windowStart.Add(45 * time.Second)
	limiter.mu.Unlock()

	if !got.Equal(want) {
		t.Fatalf("windowEndLocked() = %s, want %s", got, want)
	}
}
