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

package probe

import (
	"errors"
	"testing"

	"arcoris.dev/chrono/clock"
)

func TestWithClockRejectsNilClock(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithClock(nil)(&cfg)

	if !errors.Is(err, ErrNilClock) {
		t.Fatalf("WithClock(nil) = %v, want ErrNilClock", err)
	}
}

func TestWithClockRejectsTypedNilClock(t *testing.T) {
	t.Parallel()

	var clk *clock.FakeClock
	cfg := defaultConfig()
	err := WithClock(clk)(&cfg)

	if !errors.Is(err, ErrNilClock) {
		t.Fatalf("WithClock(typed nil) = %v, want ErrNilClock", err)
	}
}

func TestWithClockAppliesClock(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	cfg := defaultConfig()

	if err := WithClock(clk)(&cfg); err != nil {
		t.Fatalf("WithClock() = %v, want nil", err)
	}
	if cfg.clock != clk {
		t.Fatal("clock was not applied")
	}
}
