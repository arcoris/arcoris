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

// TryReserve attempts to reserve amount from l.
//
// A zero amount is a programmer error and panics. Ordinary accounting refusals
// are returned as ReserveStatus values.
func (l *ScalarLedger) TryReserve(amount Amount) ScalarReserveResult {
	l.requireNonNil()
	if amount.IsZero() {
		panicAt("amount", ErrZeroAmount, ErrorReasonZeroAmount, "amount must be positive")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	current := l.snapshotLocked()
	switch {
	case current.Value.Debt.IsPositive():
		return ScalarReserveResult{Status: ReserveStatusDebt, Snapshot: current}
	case current.Value.Available < amount:
		return ScalarReserveResult{Status: ReserveStatusInsufficient, Snapshot: current}
	}

	// Overflow is an accounting refusal, not a partial mutation.
	next, ok := l.reserved.CheckedAdd(amount)
	if !ok {
		return ScalarReserveResult{Status: ReserveStatusInsufficient, Snapshot: current}
	}

	l.reserved = next
	l.revision = l.revision.Next()

	// The reservation is created only after the scalar state is committed.
	reservation := &ScalarReservation{
		ledger: l,
		amount: amount,
	}

	return ScalarReserveResult{
		Status:      ReserveStatusReserved,
		Snapshot:    l.snapshotLocked(),
		Reservation: reservation,
	}
}
