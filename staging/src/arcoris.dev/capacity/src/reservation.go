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

// Reservation owns a multi-resource demand reserved from a Ledger.
//
// Reservation is returned only by a successful Ledger.TryReserve. It must not
// be copied after creation. Demand remains observable after release for audit
// and cleanup diagnostics.
type Reservation struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// ledger is the owner that created this reservation.
	ledger *Ledger

	// demand is the immutable resource demand owned by this reservation.
	demand Demand

	// released records whether this reservation has already returned demand.
	// It is protected by ledger.mu.
	released bool
}

// Demand returns the immutable demand owned by r.
func (r *Reservation) Demand() Demand {
	r.requireNonNil()
	return Demand{vector: r.demand.Vector()}
}

// Released reports whether r has already been released.
func (r *Reservation) Released() bool {
	r.requireNonNil()

	r.ledger.mu.Lock()
	defer r.ledger.mu.Unlock()

	r.ledger.requireInitializedLocked()
	return r.released
}
