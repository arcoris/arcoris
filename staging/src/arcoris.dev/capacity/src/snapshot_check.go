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

// CanReserve reports whether demand fits this read model.
func (s Snapshot) CanReserve(demand Demand) bool {
	return s.Check(demand).Reserved()
}

// Check evaluates demand against this read model without mutating state.
func (s Snapshot) Check(demand Demand) CheckResult {
	if !demand.IsValid() {
		panicAt(
			"demand",
			ErrInvalidDemand,
			ErrorReasonInvalidDemand,
			"demand must be non-empty and canonical",
		)
	}

	missing := make([]Entry, 0)
	debt := make([]Entry, 0)

	for _, entry := range demand.vector.entries {
		resource := entry.Resource

		switch {
		case !s.Limits.Has(resource):
			missing = append(missing, entry)

		case s.Debt.Has(resource):
			debt = append(debt, Entry{
				Resource: resource,
				Amount:   s.Debt.Amount(resource),
			})

		default:
			available := s.Available.Amount(resource)
			if available >= entry.Amount {
				continue
			}

			shortage := entry.Amount - available
			missing = append(missing, Entry{Resource: resource, Amount: shortage})
		}
	}

	return checkResultFromDiagnostics(missing, debt, s.Limits)
}

// checkResultFromDiagnostics classifies the vectors assembled for one demand.
func checkResultFromDiagnostics(missing []Entry, debt []Entry, limits Vector) CheckResult {
	switch {
	case len(missing) > 0:
		status := ReserveStatusInsufficient

		for _, entry := range missing {
			if !limits.Has(entry.Resource) {
				status = ReserveStatusUnknownResource
				break
			}
		}

		return CheckResult{Status: status, Missing: vectorFromSorted(missing)}

	case len(debt) > 0:
		return CheckResult{Status: ReserveStatusDebt, Debt: vectorFromSorted(debt)}

	default:
		return CheckResult{Status: ReserveStatusReserved}
	}
}
