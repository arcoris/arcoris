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
	"time"

	"arcoris.dev/health"
)

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(&cfg, nil)

	if !errors.Is(err, ErrNilOption) {
		t.Fatalf("applyOptions(nil) = %v, want ErrNilOption", err)
	}
}

func TestApplyOptionsAppliesInOrder(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := applyOptions(
		&cfg,
		WithInterval(time.Second),
		WithInterval(2*time.Second),
		WithTargets(health.TargetLive),
		WithTargets(health.TargetReady),
		WithInitialProbe(false),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if d := firstScheduleDelay(t, cfg.schedule); d != 2*time.Second {
		t.Fatalf("schedule delay = %s, want 2s", d)
	}
	if len(cfg.targets) != 1 || cfg.targets[0] != health.TargetReady {
		t.Fatalf("targets = %v, want [ready]", cfg.targets)
	}
	if cfg.initialProbe {
		t.Fatal("initialProbe = true, want false")
	}
}
