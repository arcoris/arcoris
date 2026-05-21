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
	"testing"

	"arcoris.dev/health"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()

	if nilClock(cfg.clock) {
		t.Fatal("default clock is nil")
	}
	if cfg.schedule == nil {
		t.Fatal("default schedule is nil")
	}
	if d := firstScheduleDelay(t, cfg.schedule); d != defaultInterval {
		t.Fatalf("default schedule delay = %s, want %s", d, defaultInterval)
	}
	if cfg.staleAfter != defaultStaleAfter {
		t.Fatalf("staleAfter = %s, want %s", cfg.staleAfter, defaultStaleAfter)
	}
	if len(cfg.targets) != 0 {
		t.Fatalf("targets = %v, want empty", cfg.targets)
	}
	if !cfg.initialProbe {
		t.Fatal("initialProbe = false, want true")
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	cfg.targets = []health.Target{health.TargetReady}

	if err := cfg.validate(); err != nil {
		t.Fatalf("validate() = %v, want nil", err)
	}
}
