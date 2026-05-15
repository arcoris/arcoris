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

// Decision reports whether work may start at one deadline decision boundary.
//
// Decision is value-only control data. It does not start work, sleep, create
// child contexts, classify operation errors, or update retry state.
type Decision struct {
	// Allowed reports whether the caller may start the work guarded by this
	// decision.
	Allowed bool

	// Remaining is the non-negative remaining budget observed for the parent
	// context.
	//
	// Remaining is zero for contexts without a deadline, expired deadlines,
	// already-done contexts without a deadline, and denied decisions without a
	// positive remaining budget.
	Remaining time.Duration

	// Reason explains why the decision was allowed or denied.
	Reason Reason
}

// IsAllowed reports whether the decision allows work to start.
func (d Decision) IsAllowed() bool {
	return d.Allowed
}

// IsDenied reports whether the decision rejects work at this boundary.
func (d Decision) IsDenied() bool {
	return !d.Allowed
}
