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

	"arcoris.dev/chrono/delay"
)

func TestWithInterval(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig()
	err := WithInterval(2 * time.Second)(&cfg)
	if err != nil {
		t.Fatalf("WithInterval() = %v, want nil", err)
	}
	if d := firstScheduleDelay(t, cfg.schedule); d != 2*time.Second {
		t.Fatalf("schedule delay = %s, want 2s", d)
	}
}

func TestWithIntervalRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	tests := []time.Duration{0, -time.Nanosecond}

	for _, interval := range tests {
		t.Run(interval.String(), func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig()
			err := WithInterval(interval)(&cfg)

			if !errors.Is(err, ErrInvalidInterval) {
				t.Fatalf("WithInterval(%s) = %v, want ErrInvalidInterval", interval, err)
			}
		})
	}
}

func TestWithSchedule(t *testing.T) {
	t.Parallel()

	schedule := delay.Delays(10 * time.Second)
	cfg := defaultConfig()
	err := WithSchedule(schedule)(&cfg)
	if err != nil {
		t.Fatalf("WithSchedule() = %v, want nil", err)
	}
	if d := firstScheduleDelay(t, cfg.schedule); d != 10*time.Second {
		t.Fatalf("schedule delay = %s, want 10s", d)
	}
}

func TestWithScheduleRejectsNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		schedule delay.Schedule
	}{
		{
			name:     "nil interface",
			schedule: nil,
		},
		{
			name:     "typed nil",
			schedule: (*nilSequenceSchedule)(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig()
			err := WithSchedule(tc.schedule)(&cfg)

			if !errors.Is(err, ErrNilSchedule) {
				t.Fatalf("WithSchedule() = %v, want ErrNilSchedule", err)
			}
		})
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
