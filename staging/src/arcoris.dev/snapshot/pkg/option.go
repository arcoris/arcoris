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
	"arcoris.dev/chrono/clock"
)

// Option configures Store and Publisher construction.
//
// Options are intentionally small and shared between holders. The package should
// not grow broad configuration policy; domain-specific behavior belongs in the
// component that owns the state.
type Option func(*config)

// WithClock configures the passive clock used for stamped snapshots.
//
// The supplied clock is used only to set Stamped.Updated when Store commits a
// value or Publisher publishes a value. Snapshot does not own timers, tickers,
// sleeps, retry loops, or other runtime waiting behavior, so it depends only on
// clock.PassiveClock.
//
// WithClock panics when clk is nil.
func WithClock(clk clock.PassiveClock) Option {
	if clk == nil {
		panic("snapshot: nil clock")
	}

	return func(cfg *config) {
		cfg.clock = clk
	}
}
