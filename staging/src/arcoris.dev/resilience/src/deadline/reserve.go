// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	result := ReserveBudget(ctx, now, reserve)
	return result.Duration, result.Bounded, result.OK
}

// ReserveBudget subtracts reserve from ctx's remaining deadline budget at now.
//
// ReserveBudget is the named-result form of Reserve. It preserves the inspected
// Budget and local Reason so callers can branch without relying on tuple
// position. It does not choose fallback timeouts for unbounded contexts and does
// not create child contexts, timers, goroutines, waits, or retries.
func ReserveBudget(ctx context.Context, now time.Time, reserve time.Duration) ReserveResult {
	requireContext(ctx)
	requireNonNegativeDuration("reserve", reserve)

	budget := Inspect(ctx, now)
	if budget.Expired {
		return ReserveResult{
			Bounded: true,
			Reason:  ReasonExpired,
			Budget:  budget,
		}
	}

	if ctx.Err() != nil {
		return ReserveResult{
			Bounded: budget.HasDeadline,
			Reason:  ReasonContextDone,
			Budget:  budget,
		}
	}

	if !budget.HasDeadline {
		return ReserveResult{
			OK:     true,
			Reason: ReasonNoDeadline,
			Budget: budget,
		}
	}

	if budget.Remaining <= reserve {
		return ReserveResult{
			Bounded: true,
			Reason:  ReasonInsufficientBudget,
			Budget:  budget,
		}
	}

	return ReserveResult{
		Duration: budget.Remaining - reserve,
		Bounded:  true,
		OK:       true,
		Reason:   ReasonAllowed,
		Budget:   budget,
	}
}
