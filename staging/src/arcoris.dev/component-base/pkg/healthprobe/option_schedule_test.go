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

package healthprobe

import (
	"errors"
	"testing"
	"time"
)

func TestWithInterval(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithInterval(2 * time.Second)(&cfg)
	if err != nil {
		t.Fatalf("WithInterval() = %v, want nil", err)
	}
	if cfg.interval != 2*time.Second {
		t.Fatalf("interval = %s, want 2s", cfg.interval)
	}
}

func TestWithIntervalRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithInterval(0)(&cfg)

	if !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("WithInterval(0) = %v, want ErrInvalidInterval", err)
	}
}

func TestWithInitialProbe(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithInitialProbe(false)(&cfg)
	if err != nil {
		t.Fatalf("WithInitialProbe(false) = %v, want nil", err)
	}
	if cfg.initialProbe {
		t.Fatal("initialProbe = true, want false")
	}
}
