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
	"sync"
	"testing"
	"time"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

var fixedWindowTestNow = time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(now time.Time) *fakeClock {
	return &fakeClock{now: now}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *fakeClock) Set(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = now
}

func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func newTestLimiter(t *testing.T, opts ...Option) (*Limiter, *fakeClock) {
	t.Helper()
	clk := newFakeClock(fixedWindowTestNow)
	all := append([]Option{WithClock(clk)}, opts...)
	limiter, err := New(all...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return limiter, clk
}

func requireValidSnapshot(t *testing.T, snap snapshot.Snapshot[retrybudget.Snapshot]) {
	t.Helper()
	if snap.IsZeroRevision() {
		t.Fatalf("Snapshot revision is zero")
	}
	if !snap.Value.IsValid() {
		t.Fatalf("Snapshot value is invalid: %+v", snap.Value)
	}
}

func requireDecision(t *testing.T, got retrybudget.Decision, allowed bool, reason retrybudget.Reason) {
	t.Helper()
	if got.Allowed != allowed {
		t.Fatalf("Decision.Allowed = %v, want %v", got.Allowed, allowed)
	}
	if got.Reason != reason {
		t.Fatalf("Decision.Reason = %s, want %s", got.Reason, reason)
	}
	if !got.IsValid() {
		t.Fatalf("Decision is invalid: %+v", got)
	}
}
