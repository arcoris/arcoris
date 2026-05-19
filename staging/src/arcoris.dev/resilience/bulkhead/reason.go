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

import "errors"

// Reason classifies a bulkhead admission decision.
//
// Reasons are stable diagnostic values. They are intentionally small because the
// base limiter has only two outcomes: allowed admission or full capacity.
type Reason string

const (
	// ReasonUnknown is the zero reason and is invalid for completed decisions.
	ReasonUnknown Reason = ""

	// ReasonAllowed reports that a permit was acquired.
	ReasonAllowed Reason = "allowed"

	// ReasonFull reports that the limiter had no available capacity.
	ReasonFull Reason = "full"
)

var (
	// ErrFull reports a rejected admission because the bulkhead is full.
	ErrFull = errors.New("bulkhead: full")
)

// String returns the stable textual reason.
func (r Reason) String() string {
	return string(r)
}

// IsValid reports whether r is a known non-zero reason.
func (r Reason) IsValid() bool {
	switch r {
	case ReasonAllowed, ReasonFull:
		return true
	default:
		return false
	}
}

// IsDenied reports whether r describes a denied admission.
func (r Reason) IsDenied() bool {
	return r == ReasonFull
}

// Err returns the error associated with r.
//
// Allowed decisions have no error. Unknown reasons also return nil here because
// Decision.Err performs full decision validation before delegating to Reason.
func (r Reason) Err() error {
	if r == ReasonFull {
		return ErrFull
	}

	return nil
}
