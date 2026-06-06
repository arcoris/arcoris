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

// Release returns r's demand to its vector ledger.
//
// Release panics on double release. Use TryRelease for idempotent cleanup.
func (r *VectorReservation) Release() {
	if !r.TryRelease() {
		panicAt(
			"vector_reservation",
			ErrReservationReleased,
			"vector reservation has already been released",
		)
	}
}

// TryRelease returns r's demand to its vector ledger if r is still live.
func (r *VectorReservation) TryRelease() bool {
	return r.release()
}

// ReleaseObserved releases r and returns the post-release vector snapshot.
func (r *VectorReservation) ReleaseObserved() snapshot.Snapshot[VectorSnapshot] {
	snap, ok := r.TryReleaseObserved()
	if !ok {
		panicAt(
			"vector_reservation",
			ErrReservationReleased,
			"vector reservation has already been released",
		)
	}

	return snap
}

// TryReleaseObserved releases r and returns the current vector snapshot and outcome.
func (r *VectorReservation) TryReleaseObserved() (snapshot.Snapshot[VectorSnapshot], bool) {
	r.requireValid()

	l := r.ledger
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	ok := r.releaseLocked()

	return l.snapshotLocked(), ok
}

// release performs the idempotent vector ownership transition and subtraction.
func (r *VectorReservation) release() bool {
	r.requireValid()

	l := r.ledger
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	return r.releaseLocked()
}

// releaseLocked releases r while its ledger mutex is held.
func (r *VectorReservation) releaseLocked() bool {
	if r.released {
		return false
	}

	l := r.ledger
	next, ok := l.state.WithoutReserved(r.demand)
	if !ok {
		panicAt(
			"vector_ledger.reserved",
			ErrReservedUnderflow,
			"reserved vector no longer covers reservation demand",
		)
	}

	l.state = next
	r.released = true
	l.revision = l.revision.Next()

	return true
}
