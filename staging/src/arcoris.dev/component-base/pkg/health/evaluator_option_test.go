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

package health

import (
	"errors"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
)

func TestDefaultEvaluatorConfig(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	if config.clock == nil {
		t.Fatal("default clock is nil")
	}
	if config.defaultTimeout != defaultCheckTimeout {
		t.Fatalf("default timeout = %s, want %s", config.defaultTimeout, defaultCheckTimeout)
	}
	if config.targetTimeouts == nil {
		t.Fatal("target timeouts map is nil")
	}
}

func TestApplyEvaluatorOptions(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	fakeClock := clock.NewFakeClock(testObserved)

	err := applyEvaluatorOptions(
		&config,
		WithDefaultTimeout(2*time.Second),
		WithTargetTimeout(TargetReady, 3*time.Second),
		WithClock(fakeClock),
	)
	if err != nil {
		t.Fatalf("applyEvaluatorOptions() = %v, want nil", err)
	}
	if config.defaultTimeout != 2*time.Second {
		t.Fatalf("default timeout = %s, want 2s", config.defaultTimeout)
	}
	if config.targetTimeouts[TargetReady] != 3*time.Second {
		t.Fatalf("ready timeout = %s, want 3s", config.targetTimeouts[TargetReady])
	}
	if config.clock != fakeClock {
		t.Fatal("clock was not configured")
	}
}

func TestEvaluatorOptionsRejectInvalidInputs(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()

	if err := applyEvaluatorOptions(&config, nil); !errors.Is(err, ErrNilEvaluatorOption) {
		t.Fatalf("nil option = %v, want ErrNilEvaluatorOption", err)
	}
	if err := WithClock(nil)(&config); !errors.Is(err, ErrNilClock) {
		t.Fatalf("WithClock(nil) = %v, want ErrNilClock", err)
	}
	if err := WithDefaultTimeout(-time.Second)(&config); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("WithDefaultTimeout(-1s) = %v, want ErrInvalidTimeout", err)
	}
	if err := WithTargetTimeout(TargetUnknown, time.Second)(&config); !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("WithTargetTimeout(invalid target) = %v, want ErrInvalidTarget", err)
	}
	if err := WithTargetTimeout(TargetReady, -time.Second)(&config); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("WithTargetTimeout(-1s) = %v, want ErrInvalidTimeout", err)
	}
}
