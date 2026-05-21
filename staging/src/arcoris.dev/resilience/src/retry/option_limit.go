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

import "time"

// WithMaxAttempts configures the maximum number of operation attempts.
//
// The limit includes the initial operation call:
//
//   - WithMaxAttempts(1) allows only the initial call and no retry attempts;
//   - WithMaxAttempts(2) allows the initial call and one retry attempt;
//   - WithMaxAttempts(n) allows at most n total operation calls.
//
// Max-attempt exhaustion is retry-owned exhaustion and is represented by
// ErrExhausted with StopReasonMaxAttempts.
//
// WithMaxAttempts panics when n is zero.
func WithMaxAttempts(n uint) Option {
	requireMaxAttempts(n)

	return func(cfg *config) {
		cfg.maxAttempts = n
	}
}

// WithMaxElapsed configures the maximum elapsed runtime for one retry execution.
//
// A zero duration disables elapsed-time limiting. A positive duration limits the
// total elapsed runtime measured by the configured retry clock.
//
// Retry should check this boundary before sleeping for a delay that would exceed
// the configured elapsed limit. Max-elapsed exhaustion is retry-owned exhaustion
// and is represented by ErrExhausted with StopReasonMaxElapsed.
//
// WithMaxElapsed panics when d is negative.
func WithMaxElapsed(d time.Duration) Option {
	requireMaxElapsed(d)

	return func(cfg *config) {
		cfg.maxElapsed = d
	}
}
