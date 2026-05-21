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
	"errors"
	"testing"
	"time"

	"arcoris.dev/resilience/retrybudget"
)

func TestNew(t *testing.T) {
	l, _ := newTestLimiter(t, WithWindow(30*time.Second), WithRatio(0.25), WithMinRetries(4))

	snap := l.Snapshot()
	requireValidSnapshot(t, snap)
	if snap.Value.Window.Duration != 30*time.Second {
		t.Fatalf("Window.Duration = %s, want 30s", snap.Value.Window.Duration)
	}
	if !snap.Value.Window.StartedAt.Equal(fixedWindowTestNow) {
		t.Fatalf("Window.StartedAt = %s, want %s", snap.Value.Window.StartedAt, fixedWindowTestNow)
	}
	if snap.Value.Policy.Ratio != 0.25 || snap.Value.Policy.Minimum != 4 {
		t.Fatalf("Policy = %+v, want ratio=0.25 minimum=4", snap.Value.Policy)
	}
}

func TestNewPublishesValidInitialSnapshotWithZeroMinimum(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(0))

	snap := l.Snapshot()
	requireValidSnapshot(t, snap)

	if snap.IsZeroRevision() {
		t.Fatalf("Snapshot revision is zero")
	}
	if snap.Value.Kind != retrybudget.KindFixedWindow {
		t.Fatalf("Kind = %s, want %s", snap.Value.Kind, retrybudget.KindFixedWindow)
	}
	if snap.Value.Attempts.Original != 0 || snap.Value.Attempts.Retry != 0 {
		t.Fatalf("Attempts = %+v, want zero attempts", snap.Value.Attempts)
	}
	if snap.Value.Capacity.Allowed != 0 {
		t.Fatalf("Capacity.Allowed = %d, want 0", snap.Value.Capacity.Allowed)
	}
	if snap.Value.Capacity.Available != 0 {
		t.Fatalf("Capacity.Available = %d, want 0", snap.Value.Capacity.Available)
	}
	if !snap.Value.Capacity.Exhausted {
		t.Fatalf("Capacity.Exhausted = false, want true")
	}
	if !snap.Value.Window.Bounded {
		t.Fatalf("Window.Bounded = false, want true")
	}
	if !snap.Value.Policy.Bounded {
		t.Fatalf("Policy.Bounded = false, want true")
	}
	if snap.Value.Policy.Minimum != 0 {
		t.Fatalf("Policy.Minimum = %d, want 0", snap.Value.Policy.Minimum)
	}
}

func TestNewBindsPublisherClock(t *testing.T) {
	l, _ := newTestLimiter(t)

	stamped := l.published.Stamped()
	if !stamped.Updated.Equal(fixedWindowTestNow) {
		t.Fatalf("Stamped.Updated = %s, want %s", stamped.Updated, fixedWindowTestNow)
	}
}

func TestNewReturnsConfigError(t *testing.T) {
	_, err := New(WithWindow(0))
	if !errors.Is(err, ErrInvalidWindow) {
		t.Fatalf("New() error = %v, want %v", err, ErrInvalidWindow)
	}
}
