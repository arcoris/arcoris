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
	"time"

	"arcoris.dev/chrono/delay"
)

// WithInterval configures a fixed schedule between probe cycles.
//
// The interval must be positive. WithInterval is a convenience wrapper around
// WithSchedule(delay.Fixed(interval)) for callers that need a simple fixed
// cadence.
func WithInterval(interval time.Duration) Option {
	return func(cfg *config) error {
		if err := validateInterval(interval); err != nil {
			return err
		}

		cfg.schedule = delay.Fixed(interval)
		return nil
	}
}

// WithSchedule configures the delay schedule between probe cycles.
//
// Each Run call creates its own independent sequence from schedule. Zero delays
// produced by the sequence are valid and run the next cycle immediately. Finite
// sequence exhaustion stops Run with ErrExhaustedSchedule. Target-specific
// schedules, result-dependent backoff, retry behavior, and jitter-specific
// options remain outside the base runner contract; callers can compose those at
// the delay.Schedule layer.
func WithSchedule(sched delay.Schedule) Option {
	return func(cfg *config) error {
		if err := validateSchedule(sched); err != nil {
			return err
		}

		cfg.schedule = sched
		return nil
	}
}

// WithInitialProbe configures whether Runner performs a probe cycle immediately
// when Run starts.
//
// Initial probing is enabled by default so Snapshot and Snapshots can become
// useful before the first schedule delay elapses. Disabling the initial probe is
// useful in tests or in components that want the first observation to align with
// the configured schedule exactly.
func WithInitialProbe(enabled bool) Option {
	return func(cfg *config) error {
		cfg.initialProbe = enabled
		return nil
	}
}
