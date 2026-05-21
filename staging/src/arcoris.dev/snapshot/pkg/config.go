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

// config contains shared holder configuration.
//
// The configuration is internal so public construction remains small and the
// package can add holder-specific options later without exposing implementation
// detail.
type config struct {
	// clock provides local commit and publication timestamps for Stamped values.
	//
	// Store and Publisher only need read-only time, so the narrow PassiveClock
	// contract is intentionally used instead of the full clock.Clock interface.
	clock clock.PassiveClock
}

// defaultConfig returns the package default holder configuration.
//
// RealClock is the production default. Tests that need deterministic timestamps
// should pass clock.NewFakeClock through WithClock.
func defaultConfig() config {
	return config{
		clock: clock.RealClock{},
	}
}

// newConfig applies opts to the package defaults and returns the final config.
//
// Options are applied in order. Later options may override values from earlier
// options.
func newConfig(opts ...Option) config {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt == nil {
			panic("snapshot: nil option")
		}
		opt(&cfg)
	}
	return cfg
}
