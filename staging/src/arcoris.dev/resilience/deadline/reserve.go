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

// Reserve subtracts reserve from ctx's remaining budget at now.
//
// Reserve is useful when a caller must leave tail budget for cleanup, response
// writing, release paths, or parent-owned coordination. When ctx has no deadline,
// Reserve returns zero, true to indicate that no bounded duration was derived
// but the operation is not deadline-limited.
func Reserve(ctx context.Context, now time.Time, reserve time.Duration) (duration time.Duration, ok bool) {
	requireContext(ctx)
	requireNonNegativeDuration("reserve", reserve)

	budget, ok := activeBudget(ctx, now)
	if !ok {
		return 0, false
	}

	if !budget.HasDeadline {
		return 0, true
	}

	if budget.Remaining <= reserve {
		return 0, false
	}

	return budget.Remaining - reserve, true
}
