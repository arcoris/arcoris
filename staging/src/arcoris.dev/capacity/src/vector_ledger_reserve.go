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
// The operation is non-blocking and all-or-nothing. On refusal it leaves vector
// state and revision unchanged and returns nil, false. TryReserve does not build
// a snapshot; use TryReserveObserved when diagnostics are required.
func (l *VectorLedger) TryReserve(demand Demand) (*VectorReservation, bool) {
	requireValidDemand("demand", demand)
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	reservation, _, ok := l.tryReserveLocked(demand)

	return reservation, ok
}

// TryReserveObserved attempts to reserve demand and returns vector diagnostics.
func (l *VectorLedger) TryReserveObserved(demand Demand) (*VectorReservation, VectorObservation, bool) {
	requireValidDemand("demand", demand)
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	reservation, fit, ok := l.tryReserveLocked(demand)

	if !ok {
		return nil, VectorObservation{
			Snapshot: l.snapshotLocked(),
			Refusal:  fit.Refusal,
			Missing:  fit.Missing,
			Debt:     fit.Debt,
		}, false
	}

	return reservation, VectorObservation{
		Snapshot: l.snapshotLocked(),
		Refusal:  RefusalNone,
	}, true
}

// tryReserveLocked reserves demand while l.mu is held.
func (l *VectorLedger) tryReserveLocked(demand Demand) (*VectorReservation, Fit, bool) {
	next, fit := l.state.WithReserved(demand)
	if fit.Refused() {
		return nil, fit, false
	}

	l.state = next
	l.revision = l.revision.Next()

	return &VectorReservation{
		ledger: l,
		demand: demand,
	}, fit, true
}
