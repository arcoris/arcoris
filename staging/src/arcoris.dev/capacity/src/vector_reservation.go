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

// VectorReservation owns a demand reserved from a VectorLedger.
//
// VectorReservation is returned only by a successful VectorLedger reservation.
// It must not be copied after creation. Demand remains observable after release
// for cleanup diagnostics.
type VectorReservation struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// ledger is the owner that created this reservation.
	ledger *VectorLedger

	// demand is the immutable resource demand owned by this reservation.
	demand Demand

	// released records whether this reservation has returned demand.
	released bool
}

// Demand returns the immutable demand owned by r.
func (r *VectorReservation) Demand() Demand {
	r.requireValid()

	return Demand{vector: r.demand.Vector()}
}

// Released reports whether r has already been released.
func (r *VectorReservation) Released() bool {
	r.requireValid()

	r.ledger.mu.Lock()
	defer r.ledger.mu.Unlock()

	r.ledger.requireInitializedLocked()

	return r.released
}

// requireValid panics when r is nil, detached, or zero-valued.
func (r *VectorReservation) requireValid() {
	if r == nil {
		panicAt("vector_reservation", ErrNilReservation, "vector reservation receiver is nil")
	}
	if r.ledger == nil || !r.demand.IsValid() {
		panicAt("vector_reservation", ErrInvalidReservation, "reservation must be created by VectorLedger.TryReserve")
	}
}
