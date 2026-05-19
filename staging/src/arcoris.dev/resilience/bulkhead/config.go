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

import "arcoris.dev/chrono/clock"

// config contains construction-time limiter configuration.
//
// The first bulkhead implementation keeps capacity fixed for the lifetime of the
// limiter. Dynamic policy belongs to a later layer that can define what happens
// when a new limit is lower than the current in-flight count.
type config struct {
	// limit is the maximum number of concurrently held permits.
	limit uint64

	// clock provides local timestamps for snapshot publication.
	clock clock.PassiveClock
}

// newConfig builds config from the required limit and optional overrides.
func newConfig(limit uint64, opts ...Option) config {
	cfg := config{
		limit: limit,
		clock: clock.RealClock{},
	}

	for _, opt := range opts {
		if opt == nil {
			panic(ErrNilOption)
		}
		opt(&cfg)
	}

	return cfg
}
