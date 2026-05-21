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
	"testing"

	"arcoris.dev/chrono/clock"
)

func TestDefaultConfigHasClock(t *testing.T) {
	cfg := defaultConfig()
	if cfg.clock == nil {
		t.Fatal("default config clock is nil")
	}

	if _, ok := cfg.clock.(clock.RealClock); !ok {
		t.Fatalf("default config clock type = %T, want clock.RealClock", cfg.clock)
	}
}

func TestNewConfigAppliesOptionsInOrder(t *testing.T) {
	first := newTestClock()
	second := newTestClock()

	cfg := newConfig(WithClock(first), WithClock(second))
	if cfg.clock != second {
		t.Fatal("newConfig did not apply options in order")
	}
}

func TestNewConfigPanicsOnNilOption(t *testing.T) {
	requirePanicWith(t, "snapshot: nil option", func() {
		_ = newConfig(nil)
	})
}
