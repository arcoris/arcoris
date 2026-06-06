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

// Snapshot is a copyable scalar read model for Ledger.
type Snapshot struct {
	// Limit is the configured scalar capacity.
	Limit Amount

	// Reserved is capacity currently held by live reservations.
	Reserved Amount

	// Available is immediately reservable capacity when Debt is zero.
	Available Amount

	// Debt is overcommit created when Limit shrinks below Reserved.
	Debt Amount
}

// NewSnapshot derives scalar availability and debt.
func NewSnapshot(limit Amount, reserved Amount) Snapshot {
	if reserved <= limit {
		return Snapshot{
			Limit:     limit,
			Reserved:  reserved,
			Available: limit - reserved,
		}
	}

	return Snapshot{
		Limit:    limit,
		Reserved: reserved,
		Debt:     reserved - limit,
	}
}

// IsValid reports whether s is internally consistent.
func (s Snapshot) IsValid() bool {
	return s == NewSnapshot(s.Limit, s.Reserved)
}

// HasDebt reports whether s is overcommitted.
func (s Snapshot) HasDebt() bool {
	return s.Debt.IsPositive()
}

// CanReserve reports whether amount fits this scalar read model.
func (s Snapshot) CanReserve(amount Amount) bool {
	if !s.IsValid() {
		return false
	}
	if amount.IsZero() {
		return false
	}
	if s.Debt.IsPositive() {
		return false
	}

	return s.Available >= amount
}
