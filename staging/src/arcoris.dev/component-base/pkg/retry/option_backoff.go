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

import "arcoris.dev/component-base/pkg/backoff"

// WithBackoff configures the backoff schedule used between retry attempts.
//
// Retry stores the reusable Schedule, not a mutable Sequence. Each Do or DoValue
// execution creates and owns its own Sequence by calling schedule.NewSequence.
// This preserves the backoff package ownership model:
//
//   - Schedule is a reusable recipe;
//   - Sequence is a mutable per-owner delay stream.
//
// The supplied schedule is consumed only after an operation attempt fails with a
// retryable error and retry-owned limits still allow another attempt. If the
// resulting Sequence returns ok=false, retry treats that as retry-owned
// exhaustion with StopReasonBackoffExhausted.
//
// A zero delay produced by the schedule is valid and means immediate retry after
// retry-owned context checks. A negative delay with ok=true violates the
// backoff.Sequence contract and is rejected at the retry boundary.
//
// WithBackoff panics when schedule is nil.
func WithBackoff(schedule backoff.Schedule) Option {
	requireBackoff(schedule)

	return func(config *config) {
		config.backoff = schedule
	}
}
