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

	"arcoris.dev/resilience/retrybudget"
)

func TestLimiterSnapshotValueLocked(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithWindow(time.Minute), WithRatio(0.5), WithMinRetries(2))

	limiter.mu.Lock()
	limiter.original = 5
	limiter.retries = 3
	val := limiter.snapshotValueLocked()
	limiter.mu.Unlock()

	if !val.IsValid() {
		t.Fatalf("snapshot value is invalid: %+v", val)
	}
	if val.Kind != retrybudget.KindFixedWindow {
		t.Fatalf("Kind = %s, want %s", val.Kind, retrybudget.KindFixedWindow)
	}
	if val.Attempts.Original != 5 || val.Attempts.Retry != 3 {
		t.Fatalf("Attempts = %+v, want original=5 retry=3", val.Attempts)
	}
	if val.Capacity.Allowed != 4 || val.Capacity.Available != 1 || val.Capacity.Exhausted {
		t.Fatalf("Capacity = %+v, want allowed=4 available=1 exhausted=false", val.Capacity)
	}
	if !val.Window.Bounded || val.Window.Duration != time.Minute {
		t.Fatalf("Window = %+v, want bounded 1m", val.Window)
	}
	if !val.Policy.Bounded || val.Policy.Ratio != 0.5 || val.Policy.Minimum != 2 {
		t.Fatalf("Policy = %+v, want ratio=0.5 minimum=2", val.Policy)
	}
}

func TestLimiterSnapshotAndRevision(t *testing.T) {
	limiter, _ := newTestLimiter(t)

	snap := limiter.Snapshot()
	rev := limiter.Revision()

	requireValidSnapshot(t, snap)
	if snap.Revision != rev {
		t.Fatalf("Revision() = %d, want snapshot revision %d", rev, snap.Revision)
	}
	if snap.Value.Kind != retrybudget.KindFixedWindow {
		t.Fatalf("Kind = %s, want %s", snap.Value.Kind, retrybudget.KindFixedWindow)
	}
}

func TestLimiterPublishLockedAdvancesRevision(t *testing.T) {
	limiter, _ := newTestLimiter(t)
	prev := limiter.Revision()

	limiter.mu.Lock()
	snap := limiter.publishLocked()
	limiter.mu.Unlock()

	if snap.Revision == prev {
		t.Fatalf("publishLocked revision = %d, want different from %d", snap.Revision, prev)
	}
	if limiter.Revision() != snap.Revision {
		t.Fatalf("Revision() = %d, want %d", limiter.Revision(), snap.Revision)
	}
}
