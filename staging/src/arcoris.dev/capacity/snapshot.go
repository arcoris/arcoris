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

package capacity

// Snapshot is a copyable read model of one Ledger state.
//
// Snapshot contains no locks, pointers, reservations, waiters, or references
// back to the ledger. It is safe to store, compare, publish, and pass across
// component boundaries as a value.
//
// In normal state, Reserved is less than or equal to Limit, Available is
// Limit-Reserved, and Debt is zero. If a ledger limit is reduced below already
// reserved capacity, existing reservations remain valid. In that overcommitted
// state, Reserved is greater than Limit, Available is zero, and Debt is
// Reserved-Limit.
type Snapshot struct {
	// Limit is the current configured capacity limit.
	Limit Amount

	// Reserved is the amount currently held by live reservations.
	Reserved Amount

	// Available is the amount that may still be reserved immediately.
	Available Amount

	// Debt is the amount by which Reserved exceeds Limit after a limit reduction.
	Debt Amount
}

// IsValid reports whether s is internally consistent.
func (s Snapshot) IsValid() bool {
	if s.Reserved <= s.Limit {
		return s.Available == s.Limit-s.Reserved && s.Debt == 0
	}

	return s.Available == 0 && s.Debt == s.Reserved-s.Limit
}

// Exhausted reports whether no additional capacity can be reserved now.
func (s Snapshot) Exhausted() bool {
	return s.Available == 0
}

// Overcommitted reports whether the snapshot contains capacity debt.
func (s Snapshot) Overcommitted() bool {
	return s.Debt > 0
}

// CanReserve reports whether amount can be reserved from this snapshot.
//
// CanReserve is a read-model helper only. It does not reserve capacity and must
// not be used as a substitute for Ledger.TryReserve in concurrent code.
func (s Snapshot) CanReserve(amount Amount) bool {
	return amount > 0 && s.Debt == 0 && s.Available >= amount
}
