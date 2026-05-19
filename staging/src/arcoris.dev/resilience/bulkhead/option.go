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

package bulkhead

import (
	"errors"

	"arcoris.dev/chrono/clock"
)

// Option configures a Limiter at construction time.
//
// Options are applied in order by New. Passing nil as an explicit option or
// option argument panics because it is a programming error; omitting an option
// entirely selects the package default.
type Option func(*config)

var (
	// ErrNilOption reports an explicit nil Option passed to New.
	ErrNilOption = errors.New("bulkhead: nil option")

	// ErrNilClock reports an attempt to configure a nil clock.
	ErrNilClock = errors.New("bulkhead: nil clock")
)

// WithClock configures the local timestamp source used for published snapshots.
//
// The clock is used only when a snapshot is published and Stamped.Updated is set.
// It does not control blocking waits, retries, timers, deadlines, or background
// scheduling because the base bulkhead does not own those behaviors.
func WithClock(clk clock.PassiveClock) Option {
	if clk == nil {
		panic(ErrNilClock)
	}

	return func(cfg *config) {
		cfg.clock = clk
	}
}
