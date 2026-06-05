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

// Release returns r's demand to its ledger.
//
// Release panics on double release. Use TryRelease for idempotent cleanup.
func (r *Reservation) Release() snapshot.Snapshot[Snapshot] {
	snap, ok := r.TryRelease()
	if !ok {
		panicAt(
			"reservation",
			ErrReservationReleased,
			ErrorReasonReservationReleased,
			"reservation has already been released",
		)
	}

	return snap
}

// TryRelease returns r's demand to its ledger if r is still live.
//
// A previously released reservation is an expected cleanup path and returns the
// current snapshot with false. Nil, detached, or internally corrupted
// reservations panic because they violate owner-state invariants.
func (r *Reservation) TryRelease() (snapshot.Snapshot[Snapshot], bool) {
	r.requireNonNil()

	l := r.ledger
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	if r.released {
		return l.snapshotLocked(), false
	}

	next, ok := l.state.Release(r.demand)
	if !ok {
		panicAt(
			"ledger.reserved",
			ErrReservedUnderflow,
			ErrorReasonReservedUnderflow,
			"reserved state no longer covers reservation demand",
		)
	}

	l.state = next
	r.released = true
	l.revision = l.revision.Next()

	return l.snapshotLocked(), true
}
