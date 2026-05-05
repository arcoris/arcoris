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

import "time"

// WithInterval configures the fixed interval between probe cycles.
//
// The interval must be positive. Runner uses one fixed interval for all
// configured targets in v1. Target-specific intervals, jitter, adaptive backoff,
// and failure-specific scheduling are intentionally outside the base runner
// contract.
func WithInterval(interval time.Duration) Option {
	return func(config *config) error {
		if err := validateInterval(interval); err != nil {
			return err
		}

		config.interval = interval
		return nil
	}
}

// WithInitialProbe configures whether Runner performs a probe cycle immediately
// when Run starts.
//
// Initial probing is enabled by default so Snapshot and Snapshots can become
// useful before the first ticker interval elapses. Disabling the initial probe is
// useful in tests or in components that want the first observation to align with
// the periodic cadence exactly.
func WithInitialProbe(enabled bool) Option {
	return func(config *config) error {
		config.initialProbe = enabled
		return nil
	}
}
