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

// CheckResult describes a pure accounting fit check.
//
// Missing reports per-resource shortage for insufficient or unknown-resource
// refusals. Debt reports per-resource overcommit relevant to the checked
// demand. Metadata, retry policy, admission reasons, and scheduling decisions
// intentionally live above capacity.
type CheckResult struct {
	// Status classifies the accounting result.
	Status ReserveStatus

	// Missing contains requested capacity that is unavailable or unknown.
	Missing Vector

	// Debt contains existing overcommit for demanded resources.
	Debt Vector
}

// Reserved reports whether the check succeeded.
func (r CheckResult) Reserved() bool {
	return r.Status.Reserved()
}

// Denied reports whether the check refused the demand.
func (r CheckResult) Denied() bool {
	return r.Status.Denied()
}

// IsValid reports whether r is internally consistent.
func (r CheckResult) IsValid() bool {
	if !r.Status.IsValid() || !r.Missing.IsValid() || !r.Debt.IsValid() {
		return false
	}
	switch r.Status {
	case ReserveStatusReserved:
		return r.Missing.IsZero() && r.Debt.IsZero()
	case ReserveStatusInsufficient, ReserveStatusUnknownResource:
		return !r.Missing.IsZero() && r.Debt.IsZero()
	case ReserveStatusDebt:
		return r.Missing.IsZero() && !r.Debt.IsZero()
	default:
		return false
	}
}
