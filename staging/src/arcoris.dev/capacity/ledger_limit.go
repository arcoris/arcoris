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

// SetLimit replaces the ledger capacity limit and returns the resulting
// snapshot.
//
// SetLimit never revokes existing reservations. If the new limit is lower than
// currently reserved capacity, the ledger enters capacity debt: Available
// becomes zero and Debt reports the amount that must be released before new
// reservations can succeed.
//
// Setting the same limit is a successful no-op and does not advance the
// revision. A zero limit is valid and simply prevents further successful
// reservations until capacity is restored by SetLimit or existing reservations
// are released.
func (l *Ledger) SetLimit(limit Amount) snapshot.Snapshot[Snapshot] {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	if l.limit == limit {
		return l.snapshotLocked()
	}

	l.limit = limit
	l.revision = l.revision.Next()

	return l.snapshotLocked()
}
