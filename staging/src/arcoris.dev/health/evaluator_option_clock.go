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

package health

import "arcoris.dev/chrono/clock"

// WithClock configures the time source used by Evaluator.
//
// Evaluator accepts clock.PassiveClock rather than clock.Clock because it only
// reads time and measures elapsed duration. It does not create timers, tickers,
// sleeps, retry waits, periodic probes, or background loops.
//
// The configured clock is used for:
//
//   - Report.Observed;
//   - Report.Duration;
//   - Result.Observed normalization;
//   - Result.Duration normalization.
//
// When parallel execution is enabled, Evaluator may call the configured clock
// concurrently from multiple goroutines. Custom clock implementations used with
// parallel execution MUST be safe for concurrent Now and Since calls.
//
// Passing nil returns ErrNilClock.
func WithClock(source clock.PassiveClock) EvaluatorOption {
	return func(cfg *evaluatorConfig) error {
		if source == nil {
			return ErrNilClock
		}

		cfg.clock = source
		return nil
	}
}
