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

package deadline

import "time"

// Budget describes the local time budget derived from a context deadline at one
// observation time.
//
// Budget is a read model. It does not observe future context cancellation, does
// not reserve time, and does not create child contexts. Callers that need an
// operational allow/deny decision should use CanStart.
type Budget struct {
	// Deadline is the parent context deadline when HasDeadline is true.
	//
	// Deadline is the zero time when HasDeadline is false.
	Deadline time.Time

	// Remaining is the non-negative duration between the observation time and the
	// deadline.
	//
	// Remaining is zero when HasDeadline is false or Expired is true. The package
	// never exposes negative remaining durations.
	Remaining time.Duration

	// HasDeadline reports whether the inspected context had an explicit deadline.
	HasDeadline bool

	// Expired reports whether the inspected context deadline was at or before the
	// observation time.
	Expired bool
}

// HasBudget reports whether the budget has a non-expired deadline with positive
// remaining time.
func (b Budget) HasBudget() bool {
	return b.HasDeadline && !b.Expired && b.Remaining > 0
}
