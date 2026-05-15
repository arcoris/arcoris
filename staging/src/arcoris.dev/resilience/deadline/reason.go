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

// Reason describes why a deadline decision allowed or denied work.
//
// Reason values are intended for local diagnostics, tests, and higher-level
// policy branching. They are not wire-format status codes and must not be used
// as distributed compatibility contracts.
type Reason uint8

const (
	// ReasonAllowed means the operation is allowed within the observed budget.
	ReasonAllowed Reason = iota

	// ReasonContextDone means the parent context was already canceled or expired
	// according to ctx.Err at the decision boundary.
	ReasonContextDone

	// ReasonNoDeadline means the context has no deadline and the operation is
	// allowed because no local time budget limits it.
	ReasonNoDeadline

	// ReasonExpired means the context deadline is at or before the observation
	// time.
	ReasonExpired

	// ReasonInsufficientBudget means the context has a positive remaining budget,
	// but it is smaller than the minimum required by the caller.
	ReasonInsufficientBudget
)

// String returns a stable diagnostic name for r.
func (r Reason) String() string {
	switch r {
	case ReasonAllowed:
		return "allowed"
	case ReasonContextDone:
		return "context_done"
	case ReasonNoDeadline:
		return "no_deadline"
	case ReasonExpired:
		return "expired"
	case ReasonInsufficientBudget:
		return "insufficient_budget"
	default:
		return "unknown"
	}
}
