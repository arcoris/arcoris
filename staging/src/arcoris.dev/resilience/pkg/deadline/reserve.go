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

// Reserve subtracts reserve from ctx's remaining deadline budget at now.
//
// Reserve is useful when a caller must leave tail budget for cleanup, response
// writing, release paths, or parent-owned coordination before starting or
// continuing caller-owned work.
//
// The returned duration is meaningful only when both bounded and ok are true. In
// that case, duration is the finite budget remaining after reserve has been
// subtracted from the parent context deadline budget.
//
// bounded reports whether ctx had an explicit deadline at the observation
// boundary. When bounded is true and ok is false, the context was deadline-bound,
// but no usable caller budget remains. When bounded is false, ctx has no
// deadline and duration is not a finite budget.
//
// ok reports whether the observed deadline and runtime context state still allow
// caller-owned work to continue.
//
// When ctx has no deadline and is still active, Reserve returns zero, false,
// true. Reserve does not choose fallback timeouts for unbounded contexts. Callers
// that require a finite child budget must apply their own timeout policy when
// bounded is false.
func Reserve(ctx context.Context, now time.Time, reserve time.Duration) (duration time.Duration, bounded bool, ok bool) {
	requireContext(ctx)
	requireNonNegativeDuration("reserve", reserve)

	budget, active := activeBudget(ctx, now)
	if !active {
		return 0, budget.HasDeadline, false
	}

	if !budget.HasDeadline {
		return 0, false, true
	}

	if budget.Remaining <= reserve {
		return 0, true, false
	}

	return budget.Remaining - reserve, true, true
}
