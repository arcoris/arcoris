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

import "arcoris.dev/component-base/pkg/clock"

// WithClock configures the clock used by Runner.
//
// Runner requires clock.Clock rather than clock.PassiveClock because healthprobe
// owns a ticker loop. The configured clock provides observation time for cache
// updates, stale calculations through Since, and the ticker that drives Run.
//
// A nil clock is rejected with ErrNilClock. Custom clock implementations used by
// Runner must be safe for the ordinary ownership pattern where Run owns ticker
// lifecycle and concurrent readers call Snapshot or Snapshots while Run updates
// the cache.
func WithClock(clk clock.Clock) Option {
	return func(config *config) error {
		if nilClock(clk) {
			return ErrNilClock
		}

		config.clock = clk
		return nil
	}
}
