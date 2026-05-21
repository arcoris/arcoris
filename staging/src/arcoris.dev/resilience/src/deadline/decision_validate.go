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

// IsValid reports whether d satisfies deadline decision invariants.
//
// The method validates only the value shape produced by deadline decision
// boundaries. It does not re-read a context, re-run deadline math, or decide
// whether the original caller should have used a different minimum budget.
//
// Zero Decision is invalid: ReasonAllowed is the zero reason, but it is valid
// only for an allowed decision with positive remaining budget. Expired denials
// have no positive remaining budget, insufficient-budget denials preserve the
// positive budget that was too small for the caller, and context-done denials
// may preserve positive remaining budget when cancellation happens before a
// future deadline.
func (d Decision) IsValid() bool {
	if d.Remaining < 0 {
		return false
	}

	if d.Allowed {
		return d.isValidAllowed()
	}

	return d.isValidDenied()
}

// isValidAllowed validates the reason/budget shape for allowed decisions.
func (d Decision) isValidAllowed() bool {
	switch d.Reason {
	case ReasonAllowed:
		return d.Remaining > 0
	case ReasonNoDeadline:
		return d.Remaining == 0
	default:
		return false
	}
}

// isValidDenied validates the reason/budget shape for denied decisions.
func (d Decision) isValidDenied() bool {
	switch d.Reason {
	case ReasonContextDone:
		return true
	case ReasonExpired:
		return d.Remaining == 0
	case ReasonInsufficientBudget:
		return d.Remaining > 0
	default:
		return false
	}
}
