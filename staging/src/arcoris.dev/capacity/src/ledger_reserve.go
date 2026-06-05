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

// TryReserve attempts to reserve demand from l.
//
// The operation is non-blocking and all-or-nothing. On refusal it leaves ledger
// state and revision unchanged and returns no Reservation.
func (l *Ledger) TryReserve(demand Demand) ReserveResult {
	l.requireNonNil()
	if !demand.IsValid() {
		panicAt(
			"demand",
			ErrInvalidDemand,
			ErrorReasonInvalidDemand,
			"demand must be non-empty and canonical",
		)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	// Check against the committed state before mutating anything. A refusal must
	// leave both reserved amounts and revision unchanged.
	check := l.state.Check(demand)
	if check.Denied() {
		return ReserveResult{
			Status:   check.Status,
			Snapshot: l.snapshotLocked(),
			Missing:  check.Missing,
			Debt:     check.Debt,
		}
	}

	// Reserve on the pure value state first so the owner commit below is a
	// single all-or-nothing assignment.
	next, result := l.state.Reserve(demand)
	if result.Denied() {
		return ReserveResult{
			Status:   result.Status,
			Snapshot: l.snapshotLocked(),
			Missing:  result.Missing,
			Debt:     result.Debt,
		}
	}

	l.state = next
	l.revision = l.revision.Next()

	// The reservation is created only after the ledger state is committed.
	reservation := &Reservation{
		ledger: l,
		demand: demand,
	}

	return ReserveResult{
		Status:      ReserveStatusReserved,
		Snapshot:    l.snapshotLocked(),
		Reservation: reservation,
	}
}
