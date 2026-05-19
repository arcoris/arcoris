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

// TryReserve attempts to reserve amount from l.
//
// TryReserve is a non-blocking check-and-reserve operation. It does not wait,
// enqueue callers, observe contexts, apply fairness policy, or retry.
//
// On success, TryReserve returns a live Reservation, the committed snapshot after
// the reservation, and true.
//
// On failure, TryReserve returns nil, the current snapshot observed while
// holding the ledger lock, and false. Failed reservation attempts do not
// advance the ledger revision.
func (l *Ledger) TryReserve(amount Amount) (*Reservation, snapshot.Snapshot[Snapshot], bool) {
	requirePositiveAmount(amount)
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	current := l.snapshotLocked()
	if !current.Value.CanReserve(amount) {
		return nil, current, false
	}

	l.reserved += amount
	l.revision = l.revision.Next()

	res := &Reservation{
		ledger: l,
		amount: amount,
	}

	return res, l.snapshotLocked(), true
}
