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

import (
	"sync"

	"arcoris.dev/snapshot"
)

// Ledger owns local scalar capacity state.
//
// Ledger serializes limit changes, reservations, releases, and snapshot reads
// with one internal mutex so each returned snapshot is internally consistent.
// Failed reservation attempts do not mutate ledger state and do not advance the
// revision.
//
// The zero Ledger is not usable. Use NewLedger to create a ledger with an
// explicit initial limit. A Ledger must not be copied after first use.
type Ledger struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects revision, limit, and reserved.
	mu sync.Mutex

	// revision is the source-local revision of the last committed ledger state.
	revision snapshot.Revision

	// limit is the current configured capacity limit.
	limit Amount

	// reserved is the amount currently owned by live reservations.
	reserved Amount
}

// NewLedger returns a Ledger with initial limit.
//
// A zero limit is valid. A zero-limit ledger refuses reservations until its limit
// is increased with SetLimit.
func NewLedger(limit Amount) *Ledger {
	return &Ledger{
		revision: snapshot.ZeroRevision.Next(),
		limit:    limit,
	}
}
