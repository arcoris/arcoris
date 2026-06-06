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

// VectorSnapshot is a copyable multi-resource accounting read model.
//
// Limits and Reserved are source accounting state. Available and Debt are
// derived from their per-resource difference. Debt records overcommit after
// limit shrink and never revokes existing reservations.
type VectorSnapshot struct {
	// Limits is the configured resource capacity.
	Limits Vector

	// Reserved is capacity currently owned by live vector reservations.
	Reserved Vector

	// Available is per-resource capacity that can be reserved immediately.
	Available Vector

	// Debt is per-resource overcommit created by limit shrink.
	Debt Vector
}

// NewVectorSnapshot derives a read model from limits and reserved amounts.
func NewVectorSnapshot(limits Vector, reserved Vector) VectorSnapshot {
	available, debt := deriveVectorCapacity(limits.entries, reserved.entries)

	return VectorSnapshot{
		Limits:    vectorFromSorted(limits.entries),
		Reserved:  vectorFromSorted(reserved.entries),
		Available: vectorFromSorted(available),
		Debt:      vectorFromSorted(debt),
	}
}

// IsValid reports whether s is internally consistent and canonical.
func (s VectorSnapshot) IsValid() bool {
	if !s.Limits.IsValid() {
		return false
	}
	if !s.Reserved.IsValid() {
		return false
	}
	if !s.Available.IsValid() {
		return false
	}
	if !s.Debt.IsValid() {
		return false
	}

	expected := NewVectorSnapshot(s.Limits, s.Reserved)

	return s.Available.Equal(expected.Available) && s.Debt.Equal(expected.Debt)
}

// IsZero reports whether s contains no limits, reservations, availability, or debt.
func (s VectorSnapshot) IsZero() bool {
	return s.Limits.IsZero() && s.Reserved.IsZero() && s.Available.IsZero() && s.Debt.IsZero()
}

// HasDebt reports whether any resource is overcommitted.
func (s VectorSnapshot) HasDebt() bool {
	return !s.Debt.IsZero()
}

// CanReserve reports whether demand fits this read model.
func (s VectorSnapshot) CanReserve(demand Demand) bool {
	return s.Fit(demand).Fits()
}

// Fit evaluates demand against this read model without mutating state.
func (s VectorSnapshot) Fit(demand Demand) Fit {
	requireValidDemand("demand", demand)

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

			missing = append(missing, Entry{
				Resource: resource,
				Amount:   entry.Amount - available,
			})
		}
	}

	return fitFromVectors(missing, debt, s.Limits)
}

// AvailableFor returns available capacity for resource, or zero when absent.
func (s VectorSnapshot) AvailableFor(resource Resource) Amount {
	return s.Available.Amount(resource)
}

// DebtFor returns debt for resource, or zero when absent.
func (s VectorSnapshot) DebtFor(resource Resource) Amount {
	return s.Debt.Amount(resource)
}

// deriveVectorCapacity returns derived available and debt vectors.
func deriveVectorCapacity(limits []Entry, reserved []Entry) ([]Entry, []Entry) {
	available := make([]Entry, 0, len(limits))
	debt := make([]Entry, 0, len(reserved))

	limitIndex, reservedIndex := 0, 0
	for limitIndex < len(limits) && reservedIndex < len(reserved) {
		limit := limits[limitIndex]
		held := reserved[reservedIndex]

		switch {
		case limit.Resource < held.Resource:
			available = append(available, limit)
			limitIndex++

		case held.Resource < limit.Resource:
			debt = append(debt, held)
			reservedIndex++

		default:
			available, debt = appendVectorBalance(limit, held, available, debt)
			limitIndex++
			reservedIndex++
		}
	}

	available = append(available, limits[limitIndex:]...)
	debt = append(debt, reserved[reservedIndex:]...)

	return available, debt
}

// appendVectorBalance appends the non-zero difference for one shared resource.
func appendVectorBalance(limit Entry, held Entry, available []Entry, debt []Entry) ([]Entry, []Entry) {
	switch limit.Amount.Compare(held.Amount) {
	case 1:
		available = append(available, Entry{
			Resource: limit.Resource,
			Amount:   limit.Amount - held.Amount,
		})

	case -1:
		debt = append(debt, Entry{
			Resource: held.Resource,
			Amount:   held.Amount - limit.Amount,
		})
	}

	return available, debt
}

// fitFromVectors classifies demand diagnostics.
//
// Refusal precedence is unknown resource, debt, then ordinary insufficiency.
// Missing and Debt are preserved together so callers can inspect every local
// accounting problem even when one refusal value must summarize the fit.
func fitFromVectors(missing []Entry, debt []Entry, limits Vector) Fit {
	hasMissing := len(missing) > 0
	hasDebt := len(debt) > 0
	hasUnknown := false

	for _, entry := range missing {
		if !limits.Has(entry.Resource) {
			hasUnknown = true
			break
		}
	}

	switch {
	case hasUnknown:
		return Fit{
			Refusal: RefusalUnknownResource,
			Missing: vectorFromSorted(missing),
			Debt:    vectorFromSorted(debt),
		}

	case hasDebt:
		return Fit{
			Refusal: RefusalDebt,
			Missing: vectorFromSorted(missing),
			Debt:    vectorFromSorted(debt),
		}

	case hasMissing:
		return Fit{
			Refusal: RefusalInsufficient,
			Missing: vectorFromSorted(missing),
		}

	default:
		return Fit{Refusal: RefusalNone}
	}
}
