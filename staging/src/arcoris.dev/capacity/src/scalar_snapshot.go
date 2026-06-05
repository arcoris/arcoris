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

// ScalarSnapshot is a copyable read model for ScalarLedger.
type ScalarSnapshot struct {
	// Limit is the configured scalar capacity.
	Limit Amount

	// Reserved is capacity currently held by live scalar reservations.
	Reserved Amount

	// Available is immediately reservable capacity when Debt is zero.
	Available Amount

	// Debt is overcommit created when Limit shrinks below Reserved.
	Debt Amount
}

// NewScalarSnapshot derives scalar availability and debt.
func NewScalarSnapshot(limit Amount, reserved Amount) ScalarSnapshot {
	if reserved <= limit {
		return ScalarSnapshot{
			Limit:     limit,
			Reserved:  reserved,
			Available: limit - reserved,
		}
	}
	return ScalarSnapshot{
		Limit:    limit,
		Reserved: reserved,
		Debt:     reserved - limit,
	}
}

// IsValid reports whether s is internally consistent.
func (s ScalarSnapshot) IsValid() bool {
	return s == NewScalarSnapshot(s.Limit, s.Reserved)
}

// HasDebt reports whether s is overcommitted.
func (s ScalarSnapshot) HasDebt() bool {
	return s.Debt.IsPositive()
}

// CanReserve reports whether amount fits this read model.
func (s ScalarSnapshot) CanReserve(amount Amount) bool {
	return s.IsValid() && amount.IsPositive() && s.Debt.IsZero() && s.Available >= amount
}
