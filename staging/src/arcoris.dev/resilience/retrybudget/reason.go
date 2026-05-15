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

// Reason explains a retry-budget admission decision.
type Reason uint8

const (
	// ReasonUnknown is the zero reason and is not valid for produced decisions.
	ReasonUnknown Reason = iota

	// ReasonAllowed means the retry attempt was admitted.
	ReasonAllowed

	// ReasonExhausted means no retry budget capacity was available.
	ReasonExhausted
)

// String returns the stable diagnostic name for r.
func (r Reason) String() string {
	switch r {
	case ReasonAllowed:
		return "allowed"
	case ReasonExhausted:
		return "exhausted"
	default:
		return "unknown"
	}
}

// IsValid reports whether r is a valid produced decision reason.
func (r Reason) IsValid() bool {
	switch r {
	case ReasonAllowed, ReasonExhausted:
		return true
	default:
		return false
	}
}

// IsAllowed reports whether r represents an admitted retry.
func (r Reason) IsAllowed() bool {
	return r == ReasonAllowed
}

// IsDenied reports whether r represents a denied retry.
func (r Reason) IsDenied() bool {
	return r == ReasonExhausted
}
