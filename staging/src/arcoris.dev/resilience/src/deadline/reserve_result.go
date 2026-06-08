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

import "time"

// ReserveResult describes one tail-budget reservation decision.
//
// OK reports whether caller-owned work may continue after leaving the requested
// reserve. Bounded reports whether Duration is derived from a parent deadline.
// Duration is meaningful only when OK and Bounded are both true. Budget is the
// passive deadline read model observed at the same decision boundary.
type ReserveResult struct {
	// Duration is the bounded caller-owned budget left after reserve is removed.
	Duration time.Duration

	// Bounded reports whether Duration comes from a parent context deadline.
	Bounded bool

	// OK reports whether the caller may continue after leaving reserve.
	OK bool

	// Reason explains the reservation outcome using local deadline diagnostics.
	Reason Reason

	// Budget is the inspected deadline budget at the reservation boundary.
	Budget Budget
}

// IsValid reports whether r satisfies the ReserveBudget result invariants.
func (r ReserveResult) IsValid() bool {
	if r.Duration < 0 {
		return false
	}
	if !r.Reason.IsValid() {
		return false
	}
	if !r.Budget.HasDeadline && (r.Budget.Remaining != 0 || r.Budget.Expired) {
		return false
	}
	if r.Budget.HasDeadline && r.Budget.Expired && r.Budget.Remaining != 0 {
		return false
	}

	if r.OK {
		return r.isValidAllowed()
	}
	return r.isValidDenied()
}

// isValidAllowed validates successful reservation shapes.
func (r ReserveResult) isValidAllowed() bool {
	switch r.Reason {
	case ReasonNoDeadline:
		return !r.Bounded &&
			r.Duration == 0 &&
			!r.Budget.HasDeadline &&
			!r.Budget.Expired
	case ReasonAllowed:
		return r.Bounded &&
			r.Duration > 0 &&
			r.Budget.HasDeadline &&
			!r.Budget.Expired
	default:
		return false
	}
}

// isValidDenied validates failed reservation shapes.
func (r ReserveResult) isValidDenied() bool {
	switch r.Reason {
	case ReasonExpired:
		return r.Bounded &&
			r.Duration == 0 &&
			r.Budget.HasDeadline &&
			r.Budget.Expired
	case ReasonContextDone:
		return r.Duration == 0 &&
			r.Bounded == r.Budget.HasDeadline &&
			!r.Budget.Expired
	case ReasonInsufficientBudget:
		return r.Bounded &&
			r.Duration == 0 &&
			r.Budget.HasDeadline &&
			!r.Budget.Expired
	default:
		return false
	}
}
