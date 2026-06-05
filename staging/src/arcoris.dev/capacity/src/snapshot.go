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

// Snapshot is a copyable multi-resource accounting read model.
//
// Snapshot contains no locks, pointers, reservations, waiters, or references
// back to a ledger. Limits and Reserved are source accounting state. Available
// and Debt are derived from their per-resource difference.
type Snapshot struct {
	// Limits is the configured resource capacity.
	Limits Vector

	// Reserved is capacity currently owned by live reservations.
	Reserved Vector

	// Available is per-resource capacity that can be reserved immediately.
	Available Vector

	// Debt is per-resource overcommit created when limits shrink below live
	// reservations. Debt never revokes an existing reservation.
	Debt Vector
}

// NewSnapshot derives a read model from limits and reserved amounts.
func NewSnapshot(limits Vector, reserved Vector) Snapshot {
	available, debt := deriveCapacityVectors(limits.entries, reserved.entries)

	return Snapshot{
		Limits:    vectorFromSorted(limits.entries),
		Reserved:  vectorFromSorted(reserved.entries),
		Available: vectorFromSorted(available),
		Debt:      vectorFromSorted(debt),
	}
}

// IsValid reports whether s is internally consistent and canonical.
func (s Snapshot) IsValid() bool {
	if !s.Limits.IsValid() || !s.Reserved.IsValid() || !s.Available.IsValid() || !s.Debt.IsValid() {
		return false
	}
	expected := NewSnapshot(s.Limits, s.Reserved)
	return s.Available.Equal(expected.Available) && s.Debt.Equal(expected.Debt)
}

// IsZero reports whether s contains no limits, reservations, availability, or debt.
func (s Snapshot) IsZero() bool {
	return s.Limits.IsZero() && s.Reserved.IsZero() && s.Available.IsZero() && s.Debt.IsZero()
}

// HasDebt reports whether any resource is overcommitted.
func (s Snapshot) HasDebt() bool {
	return !s.Debt.IsZero()
}

// AvailableFor returns available capacity for resource, or zero when absent.
func (s Snapshot) AvailableFor(resource Resource) Amount {
	return s.Available.Amount(resource)
}

// DebtFor returns debt for resource, or zero when absent.
func (s Snapshot) DebtFor(resource Resource) Amount {
	return s.Debt.Amount(resource)
}

// deriveCapacityVectors compares canonical limit and reservation vectors and
// returns the derived available and debt vectors.
func deriveCapacityVectors(limits []Entry, reserved []Entry) ([]Entry, []Entry) {
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
			available, debt = appendResourceBalance(limit, held, available, debt)
			limitIndex++
			reservedIndex++
		}
	}

	available = append(available, limits[limitIndex:]...)
	debt = append(debt, reserved[reservedIndex:]...)

	return available, debt
}

// appendResourceBalance appends the non-zero difference for one shared resource.
func appendResourceBalance(limit Entry, held Entry, available []Entry, debt []Entry) ([]Entry, []Entry) {
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
