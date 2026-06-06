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

import (
	"sync/atomic"

	"arcoris.dev/snapshot"
)

// Ledger owns optimized scalar local capacity accounting.
//
// Ledger is the default hot-path capacity owner for single-resource limits such
// as bulkhead slots, active requests, queue slots, and worker slots. Raw reserve
// and release methods update atomic counters directly and do not allocate, build
// snapshots, or construct diagnostics. Callers that want capacity-owned release
// ownership can use TryAcquire, which returns a Reservation token.
//
// Ledger is safe for concurrent use. Successful mutations advance the source
// revision once after the accounting counter has changed. Snapshot is an
// explicit observation built from atomic values; under concurrent mutation it is
// internally valid and monotonic enough for diagnostics, but it is not a global
// serialization point for reserve and release operations.
//
// The zero Ledger is invalid. Use NewLedger. A Ledger must not be copied after
// first use.
type Ledger struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// limit is the configured scalar capacity.
	limit atomic.Uint64

	// reserved is the amount currently held by raw or owned reservations.
	reserved atomic.Uint64

	// revision is the local version of the last completed mutation.
	revision atomic.Uint64
}

// NewLedger returns a scalar ledger with initial limit.
func NewLedger(limit Amount) *Ledger {
	ledger := &Ledger{}
	ledger.limit.Store(limit.Uint64())
	ledger.revision.Store(uint64(snapshot.ZeroRevision.Next()))

	return ledger
}

// SetLimit replaces the scalar capacity limit.
//
// Existing reservations are never revoked. Lower limits may create debt. Setting
// the same limit is a no-op and does not advance the revision.
func (l *Ledger) SetLimit(limit Amount) {
	l.requireReady()
	for {
		current := l.limit.Load()
		if current == limit.Uint64() {
			return
		}
		if l.limit.CompareAndSwap(current, limit.Uint64()) {
			l.advanceRevision()
			return
		}
	}
}

// SetLimitObserved replaces the scalar capacity limit and then reads a snapshot.
func (l *Ledger) SetLimitObserved(limit Amount) snapshot.Snapshot[Snapshot] {
	l.SetLimit(limit)

	return l.Snapshot()
}

// Snapshot returns the current revisioned scalar ledger snapshot.
func (l *Ledger) Snapshot() snapshot.Snapshot[Snapshot] {
	l.requireReady()

	return snapshot.Snapshot[Snapshot]{
		Revision: l.Revision(),
		Value: NewSnapshot(
			Amount(l.limit.Load()),
			Amount(l.reserved.Load()),
		),
	}
}

// Revision returns the latest committed ledger revision.
func (l *Ledger) Revision() snapshot.Revision {
	l.requireReady()

	return snapshot.Revision(l.revision.Load())
}
