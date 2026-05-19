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

// Snapshot returns the current revisioned ledger snapshot.
//
// The returned snapshot is assembled while holding the ledger lock, so its
// revision and value always describe the same committed ledger state.
func (l *Ledger) Snapshot() snapshot.Snapshot[Snapshot] {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	return l.snapshotLocked()
}

// Revision returns the latest committed ledger revision.
//
// Revision is local to this ledger. It advances only after successful
// reservations, releases, or limit changes that commit a different ledger
// state.
func (l *Ledger) Revision() snapshot.Revision {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	return l.revision
}

// snapshotLocked returns a revisioned snapshot from fields protected by l.mu.
//
// The caller must hold l.mu and must have checked that l is initialized.
func (l *Ledger) snapshotLocked() snapshot.Snapshot[Snapshot] {
	val := Snapshot{
		Limit:    l.limit,
		Reserved: l.reserved,
	}

	if l.reserved <= l.limit {
		val.Available = l.limit - l.reserved
		val.Debt = 0
	} else {
		val.Available = 0
		val.Debt = l.reserved - l.limit
	}

	return snapshot.Snapshot[Snapshot]{
		Revision: l.revision,
		Value:    val,
	}
}
