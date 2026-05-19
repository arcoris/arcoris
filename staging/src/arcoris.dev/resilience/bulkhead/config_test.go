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

package bulkhead

import (
	"testing"
	"time"
)

type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time {
	return f.now
}

func (f fakeClock) Since(t time.Time) time.Duration {
	return f.now.Sub(t)
}

func TestNewConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := newConfig(7)
	if cfg.limit != 7 {
		t.Fatalf("limit = %d, want 7", cfg.limit)
	}
	if cfg.clock == nil {
		t.Fatal("clock is nil")
	}
}

func TestNewConfigAppliesOptionsInOrder(t *testing.T) {
	t.Parallel()

	first := fakeClock{now: time.Unix(1, 0)}
	second := fakeClock{now: time.Unix(2, 0)}

	cfg := newConfig(1, WithClock(first), WithClock(second))
	if got := cfg.clock.Now(); !got.Equal(second.now) {
		t.Fatalf("clock.Now() = %v, want %v", got, second.now)
	}
}
