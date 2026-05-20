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

package bulkhead

import (
	"arcoris.dev/capacity"
	"arcoris.dev/snapshot"
)

// Bulkhead bounds the number of operations concurrently executing inside one
// protected section.
//
// Bulkhead intentionally owns no admission queue, waiter lifecycle, fairness
// policy, retry behavior, health integration, metrics hooks, logging hooks,
// tracing hooks, or worker pool. It is a small resilience-domain wrapper around
// capacity.Ledger: acquiring a Lease reserves local scalar capacity, and
// releasing that Lease returns the same amount.
//
// The wrapped capacity.Ledger owns all low-level scalar accounting, including
// revisioned snapshots, limit changes, release ownership, and debt semantics
// after a limit is reduced below active leases. Bulkhead owns only the
// execution-protection meaning of that accounting: bounded in-flight isolation.
//
// Bulkhead is safe for concurrent use. A Bulkhead must be created with New and
// must not be copied after first use.
type Bulkhead struct {
	// noCopy lets go vet report accidental Bulkhead copies after first use.
	//
	// Copying a Bulkhead would copy the pointer to the same ledger and make
	// ownership harder to reason about. Keeping this marker matches the ARCORIS
	// pattern for values that represent mutable lifecycle state.
	noCopy noCopy

	// ledger owns the low-level live capacity accounting.
	//
	// All mutation, synchronization, revisioning, and lease/release ownership
	// checks are delegated to this ledger. Bulkhead methods do not maintain
	// parallel counters; duplicating that state here would make the accounting
	// easier to skew.
	ledger *capacity.Ledger
}

// Snapshot returns the current revisioned capacity state.
//
// The returned value is the underlying capacity.Ledger snapshot. It is safe to
// store or compare as a value. It describes local in-flight capacity only; it
// does not include health, routing, scheduling, metrics, or distributed state.
func (b *Bulkhead) Snapshot() snapshot.Snapshot[Snapshot] {
	b.requireReady()
	return b.ledger.Snapshot()
}

// Revision returns the latest committed bulkhead capacity revision.
//
// Revisions are source-local to this Bulkhead. They are useful for cheap change
// detection by consumers observing the same Bulkhead, but they are not a global
// ordering across components.
func (b *Bulkhead) Revision() snapshot.Revision {
	b.requireReady()
	return b.ledger.Revision()
}
