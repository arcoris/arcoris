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

package clock

import (
	"runtime"
	"testing"
	"time"
)

const (
	// clockTestTimeout is a safety guard for tests that involve goroutines or
	// real runtime timer channels.
	//
	// Clock tests must not depend on real-time sleeping for correctness. This
	// timeout is used only to prevent a broken test or implementation from
	// hanging the test process indefinitely.
	clockTestTimeout = 500 * time.Millisecond
)

func fakeClockTestTime() time.Time {
	return time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
}

func waitUntil(t *testing.T, description string, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(clockTestTimeout)

	for time.Now().Before(deadline) {
		if condition() {
			return
		}

		runtime.Gosched()
	}

	t.Fatalf("condition did not become true before safety timeout %s: %s", clockTestTimeout, description)
}

func mustEqualTime(t *testing.T, name string, got, want time.Time) {
	t.Helper()

	if !got.Equal(want) {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}
