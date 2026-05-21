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

package eval

import (
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

func TestApplyEvaluatorOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultEvaluatorConfig()
	fakeClock := clock.NewFakeClock(testObserved)

	err := applyEvaluatorOptions(
		&cfg,
		WithDefaultTimeout(2*time.Second),
		WithTargetTimeout(health.TargetReady, 3*time.Second),
		WithClock(fakeClock),
	)
	if err != nil {
		t.Fatalf("applyEvaluatorOptions() = %v, want nil", err)
	}
	if cfg.defaultTimeout != 2*time.Second {
		t.Fatalf("default timeout = %s, want 2s", cfg.defaultTimeout)
	}
	if cfg.targetTimeouts[health.TargetReady] != 3*time.Second {
		t.Fatalf("ready timeout = %s, want 3s", cfg.targetTimeouts[health.TargetReady])
	}
	if cfg.clock != fakeClock {
		t.Fatal("clock was not configured")
	}
}

func TestEvaluatorOptionsRejectInvalidInputs(t *testing.T) {
	t.Parallel()

	cfg := defaultEvaluatorConfig()

	if err := applyEvaluatorOptions(&cfg, nil); !errors.Is(err, ErrNilEvaluatorOption) {
		t.Fatalf("nil option = %v, want ErrNilEvaluatorOption", err)
	}
	if err := WithClock(nil)(&cfg); !errors.Is(err, ErrNilClock) {
		t.Fatalf("WithClock(nil) = %v, want ErrNilClock", err)
	}
	if err := WithDefaultTimeout(-time.Second)(&cfg); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("WithDefaultTimeout(-1s) = %v, want ErrInvalidTimeout", err)
	}
	if err := WithTargetTimeout(health.TargetUnknown, time.Second)(&cfg); !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("WithTargetTimeout(invalid target) = %v, want health.ErrInvalidTarget", err)
	}
	if err := WithTargetTimeout(health.TargetReady, -time.Second)(&cfg); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("WithTargetTimeout(-1s) = %v, want ErrInvalidTimeout", err)
	}
}
