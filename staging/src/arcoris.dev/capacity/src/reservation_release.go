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

import "arcoris.dev/snapshot"

// Release returns r's amount to its ledger.
//
// Release panics on double release. Use TryRelease for idempotent cleanup.
func (r *Reservation) Release() {
	if !r.TryRelease() {
		panicAt("reservation", ErrReservationReleased, "reservation has already been released")
	}
}

// TryRelease returns r's amount to its ledger if r is still live.
//
// A previously released reservation is an expected cleanup path and returns
// false. Nil, detached, or internally corrupted reservations panic because they
// violate owner-state invariants.
func (r *Reservation) TryRelease() bool {
	return r.release()
}

// ReleaseObserved releases r and returns the post-release snapshot.
func (r *Reservation) ReleaseObserved() snapshot.Snapshot[Snapshot] {
	snap, ok := r.TryReleaseObserved()
	if !ok {
		panicAt("reservation", ErrReservationReleased, "reservation has already been released")
	}

	return snap
}

// TryReleaseObserved releases r and returns the current snapshot and outcome.
func (r *Reservation) TryReleaseObserved() (snapshot.Snapshot[Snapshot], bool) {
	ok := r.release()

	return r.ledger.Snapshot(), ok
}

// release performs the idempotent ownership transition and ledger subtraction.
func (r *Reservation) release() bool {
	r.requireValid()

	if !r.released.CompareAndSwap(false, true) {
		return false
	}

	if !r.ledger.releaseAmount(r.amount) {
		panicAt(
			"ledger.reserved",
			ErrReservedUnderflow,
			"reserved amount no longer covers reservation",
		)
	}

	return true
}
