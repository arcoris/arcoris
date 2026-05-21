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

package snapshot

import (
	"sync"
	"testing"
	"time"
)

// testClock is a deterministic PassiveClock used by snapshot tests.
//
// Production code uses arcoris.dev/chrono/clock implementations. The test helper
// keeps snapshot tests independent from richer fake timer/ticker behavior because
// snapshot only requires Now and Since.
type testClock struct {
	mu  sync.Mutex
	now time.Time
}

// newTestClock creates a test clock initialized to the Unix epoch.
func newTestClock() *testClock {
	return &testClock{now: time.Unix(0, 0)}
}

// Now returns the current test time.
func (c *testClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

// Since returns the elapsed duration since t according to the current test time.
func (c *testClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// set replaces the current test time.
func (c *testClock) set(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = now
}

// requirePanicWith runs fn and requires it to panic with want.
//
// Snapshot uses stable package-local diagnostic strings for programmer errors,
// so tests should assert the exact panic value instead of only checking that a
// panic happened.
func requirePanicWith(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		got, ok := recovered.(string)
		if !ok {
			t.Fatalf("panic = %T(%v), want string %q", recovered, recovered, want)
		}
		if got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	fn()
}
