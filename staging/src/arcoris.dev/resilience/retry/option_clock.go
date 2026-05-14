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

package retry

import "arcoris.dev/chrono/clock"

// WithClock configures the clock used by retry execution.
//
// Retry uses the configured clock for:
//
//   - attempt start timestamps;
//   - terminal outcome timestamps;
//   - elapsed-time checks;
//   - delay timers between retry attempts.
//
// Use clock.RealClock{} for production runtime behavior and clock.FakeClock for
// deterministic tests. Callers that only need to configure timestamp generation
// cannot provide a narrower dependency because retry also owns timer-based
// delays.
//
// WithClock panics when c is nil.
func WithClock(c clock.Clock) Option {
	requireClock(c)

	return func(config *config) {
		config.clock = c
	}
}
