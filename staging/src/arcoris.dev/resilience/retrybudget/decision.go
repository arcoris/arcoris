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

package retrybudget

import "arcoris.dev/snapshot"

// Decision is the result of a retry-budget admission attempt.
//
// Snapshot carries the implementation state observed while producing the
// decision. Implementations should return a consistent snapshot from the same
// critical section or publication boundary that admitted or denied the retry.
type Decision struct {
	// Allowed reports whether the retry attempt was admitted.
	Allowed bool

	// Reason explains why the retry attempt was admitted or denied.
	Reason Reason

	// Snapshot is the revisioned retry-budget state associated with the decision.
	Snapshot snapshot.Snapshot[Snapshot]
}

// IsAllowed reports whether d admits the retry attempt.
func (d Decision) IsAllowed() bool {
	return d.Allowed
}

// IsDenied reports whether d denies the retry attempt.
func (d Decision) IsDenied() bool {
	return !d.Allowed
}

// IsValid reports whether d is an internally consistent retry-budget decision.
//
// The zero Decision is invalid. Allowed decisions must use ReasonAllowed. Denied
// decisions must use a denied reason such as ReasonExhausted. The associated
// domain snapshot value must also be valid.
func (d Decision) IsValid() bool {
	if !d.Reason.IsValid() {
		return false
	}
	if d.Snapshot.IsZeroRevision() {
		return false
	}
	if !d.Snapshot.Value.IsValid() {
		return false
	}
	if d.Allowed {
		return d.Reason == ReasonAllowed
	}
	return d.Reason.IsDenied()
}
