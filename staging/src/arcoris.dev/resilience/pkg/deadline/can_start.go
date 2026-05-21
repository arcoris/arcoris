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

// CanStart decides whether work requiring at least min budget may start at now.
//
// CanStart is an operational decision, not a pure deadline inspection. It denies
// work when ctx is already done, when the observed deadline is expired, or when
// the remaining budget is smaller than min.
func CanStart(ctx context.Context, now time.Time, min time.Duration) Decision {
	requireContext(ctx)
	requireNonNegativeDuration("min", min)

	budget := Inspect(ctx, now)
	if budget.Expired {
		return Decision{Reason: ReasonExpired}
	}

	if ctx.Err() != nil {
		return Decision{
			Remaining: budget.Remaining,
			Reason:    ReasonContextDone,
		}
	}

	if !budget.HasDeadline {
		return Decision{
			Allowed: true,
			Reason:  ReasonNoDeadline,
		}
	}

	if budget.Remaining < min {
		return Decision{
			Remaining: budget.Remaining,
			Reason:    ReasonInsufficientBudget,
		}
	}

	return Decision{
		Allowed:   true,
		Remaining: budget.Remaining,
		Reason:    ReasonAllowed,
	}
}
