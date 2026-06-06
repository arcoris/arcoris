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

// Fit describes whether a demand fits a vector accounting state.
//
// Fit is pure accounting diagnostics. It does not grant ownership, enqueue
// work, choose retry policy, or make an admission decision.
type Fit struct {
	// Refusal classifies why the demand does not fit.
	Refusal Refusal

	// Missing reports per-resource shortage for insufficient or unknown resources.
	Missing Vector

	// Debt reports per-resource debt that blocks the demanded resources.
	Debt Vector
}

// Fits reports whether the checked demand can be reserved.
func (f Fit) Fits() bool {
	return f.Refusal == RefusalNone
}

// Refused reports whether the checked demand was refused by local accounting.
func (f Fit) Refused() bool {
	return f.Refusal.Refused()
}

// IsValid reports whether f is internally consistent.
func (f Fit) IsValid() bool {
	if !f.Refusal.IsValid() {
		return false
	}
	if !f.Missing.IsValid() || !f.Debt.IsValid() {
		return false
	}

	switch f.Refusal {
	case RefusalNone:
		return f.Missing.IsZero() && f.Debt.IsZero()
	case RefusalInsufficient, RefusalUnknownResource:
		return !f.Missing.IsZero() && f.Debt.IsZero()
	case RefusalDebt:
		return f.Missing.IsZero() && !f.Debt.IsZero()
	default:
		return false
	}
}
