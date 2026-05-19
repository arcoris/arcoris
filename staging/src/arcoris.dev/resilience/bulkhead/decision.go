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

package bulkhead

import (
	"errors"

	"arcoris.dev/snapshot"
)

// Decision describes the result of one TryAcquire call.
//
// A Decision is metadata. The acquired capability is returned separately as a
// Permit so admission state and caller-owned lifecycle ownership remain distinct.
type Decision struct {
	// Allowed reports whether a permit was acquired.
	Allowed bool

	// Reason classifies the admission result.
	Reason Reason

	// Snapshot is the limiter state published for this decision.
	Snapshot snapshot.Snapshot[Snapshot]
}

var (
	// ErrInvalidDecision reports an inconsistent Decision value.
	ErrInvalidDecision = errors.New("bulkhead: invalid decision")
)

// IsValid reports whether d is internally consistent.
func (d Decision) IsValid() bool {
	if !d.Reason.IsValid() {
		return false
	}
	if d.Snapshot.IsZeroRevision() || !d.Snapshot.Value.IsValid() {
		return false
	}

	if d.Allowed {
		return d.Reason == ReasonAllowed
	}

	return d.Reason.IsDenied()
}

// Err returns the error represented by a denied decision.
//
// Allowed decisions return nil. Invalid decisions return ErrInvalidDecision so
// callers do not accidentally treat malformed metadata as successful admission.
func (d Decision) Err() error {
	if !d.IsValid() {
		return ErrInvalidDecision
	}
	if d.Allowed {
		return nil
	}

	return d.Reason.Err()
}
