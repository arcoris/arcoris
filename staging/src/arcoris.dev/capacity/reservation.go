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

import "arcoris.dev/snapshot"

// Reservation owns capacity reserved from a Ledger.
//
// A Reservation is returned only by Ledger.TryReserve. It holds its amount until
// Release or TryRelease returns that amount to the ledger.
//
// Reservation must not be copied after creation. Copying a live reservation can
// split release ownership and corrupt accounting.
type Reservation struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// ledger is the owner that created this reservation.
	ledger *Ledger

	// amount is the immutable amount reserved from ledger.
	amount Amount

	// released records whether this reservation has already returned its amount.
	// It is protected by ledger.mu.
	released bool
}

// Amount returns the amount reserved by r.
//
// The amount is immutable for the lifetime of the reservation and remains
// observable after release so callers can log or audit what was owned.
func (r *Reservation) Amount() Amount {
	r.requireNonNil()
	return r.amount
}

// Released reports whether r has already been released.
//
// Released is serialized by the owning ledger lock so concurrent callers see a
// consistent ownership state.
func (r *Reservation) Released() bool {
	r.requireNonNil()

	r.ledger.mu.Lock()
	defer r.ledger.mu.Unlock()

	r.ledger.requireInitializedLocked()
	return r.released
}

// Release returns r's capacity to its ledger.
//
// Release panics if r has already been released. Double release is an ownership
// bug. Use TryRelease when idempotent cleanup is required.
func (r *Reservation) Release() snapshot.Snapshot[Snapshot] {
	snap, ok := r.TryRelease()
	if !ok {
		panic(errReservationAlreadyReleased)
	}

	return snap
}

// TryRelease returns r's capacity to its ledger if r is still live.
//
// On first release, TryRelease updates the ledger, advances the revision, and
// returns the resulting snapshot with true.
//
// If r has already been released, TryRelease leaves the ledger unchanged and
// returns the current snapshot with false.
func (r *Reservation) TryRelease() (snapshot.Snapshot[Snapshot], bool) {
	r.requireNonNil()

	l := r.ledger
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	if r.released {
		return l.snapshotLocked(), false
	}
	if l.reserved < r.amount {
		panic(errLedgerReservedUnderflow)
	}

	l.reserved -= r.amount
	r.released = true
	l.revision = l.revision.Next()

	return l.snapshotLocked(), true
}
