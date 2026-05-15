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

func TestNewReturnsConfigError(t *testing.T) {
	_, err := New(WithWindow(0))
	if !errors.Is(err, ErrInvalidWindow) {
		t.Fatalf("New() error = %v, want %v", err, ErrInvalidWindow)
	}
}
