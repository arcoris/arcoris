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

package capacity

// ReserveStatus classifies a local accounting reservation check.
//
// The status is policy-free. Higher layers may map it to admission outcomes,
// retry behavior, logs, or metrics, but capacity itself only reports accounting
// facts.
type ReserveStatus uint8

// Reserve statuses are the closed local accounting result vocabulary.
const (
	// ReserveStatusReserved means the requested demand fit and may be reserved.
	ReserveStatusReserved ReserveStatus = iota + 1

	// ReserveStatusInsufficient means all resources are known but at least one
	// demanded resource lacks enough available capacity.
	ReserveStatusInsufficient

	// ReserveStatusDebt means at least one demanded resource is already in debt.
	ReserveStatusDebt

	// ReserveStatusUnknownResource means the demand references a resource absent
	// from the configured limits.
	ReserveStatusUnknownResource
)

// IsValid reports whether s is one of the closed status values.
func (s ReserveStatus) IsValid() bool {
	switch s {
	case ReserveStatusReserved, ReserveStatusInsufficient, ReserveStatusDebt, ReserveStatusUnknownResource:
		return true
	default:
		return false
	}
}

// Reserved reports whether s represents successful accounting fit.
func (s ReserveStatus) Reserved() bool {
	return s == ReserveStatusReserved
}

// Denied reports whether s represents an accounting refusal.
func (s ReserveStatus) Denied() bool {
	return s == ReserveStatusInsufficient || s == ReserveStatusDebt || s == ReserveStatusUnknownResource
}

// String returns the stable diagnostic spelling for s.
func (s ReserveStatus) String() string {
	switch s {
	case ReserveStatusReserved:
		return "reserved"
	case ReserveStatusInsufficient:
		return "insufficient"
	case ReserveStatusDebt:
		return "debt"
	case ReserveStatusUnknownResource:
		return "unknown_resource"
	default:
		return "invalid"
	}
}
