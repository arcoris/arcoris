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

import (
	"context"
	"time"
)

// Clamp bounds requested to ctx's remaining budget at now.
//
// Clamp returns requested unchanged when ctx has no deadline. It returns the
// remaining budget when requested is larger than the remaining budget. It returns
// ok=false when ctx is already done or the observed deadline is expired.
func Clamp(ctx context.Context, now time.Time, requested time.Duration) (duration time.Duration, ok bool) {
	requireContext(ctx)
	requireNonNegativeDuration("requested", requested)

	budget, ok := activeBudget(ctx, now)
	if !ok {
		return 0, false
	}

	if !budget.HasDeadline {
		return requested, true
	}

	if requested > budget.Remaining {
		return budget.Remaining, true
	}

	return requested, true
}
