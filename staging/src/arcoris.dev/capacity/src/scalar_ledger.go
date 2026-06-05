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
	"sync"

	"arcoris.dev/snapshot"
)

// ScalarLedger owns optimized single-resource local capacity accounting.
//
// ScalarLedger mirrors Ledger semantics for one amount: debt after limit shrink,
// strict reservation ownership, failed-reserve no mutation, and revisioned
// snapshots. Use Ledger for multi-resource all-or-nothing accounting.
type ScalarLedger struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects revision, limit, reserved, and reservation release ownership.
	mu sync.Mutex

	// revision is the local version of the last committed scalar state.
	revision snapshot.Revision

	// limit is the configured scalar capacity.
	limit Amount

	// reserved is the amount held by live scalar reservations.
	reserved Amount
}

// NewScalarLedger returns a scalar ledger with initial limit.
func NewScalarLedger(limit Amount) *ScalarLedger {
	return &ScalarLedger{
		revision: snapshot.ZeroRevision.Next(),
		limit:    limit,
	}
}

// SetLimit replaces the scalar capacity limit.
//
// Existing reservations are never revoked. Lower limits may create debt.
func (l *ScalarLedger) SetLimit(limit Amount) snapshot.Snapshot[ScalarSnapshot] {
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

// Snapshot returns the current revisioned scalar ledger snapshot.
func (l *ScalarLedger) Snapshot() snapshot.Snapshot[ScalarSnapshot] {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	return l.snapshotLocked()
}

// Revision returns the latest committed scalar ledger revision.
func (l *ScalarLedger) Revision() snapshot.Revision {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	return l.revision
}

// snapshotLocked derives a revisioned snapshot for scalar state protected by l.mu.
func (l *ScalarLedger) snapshotLocked() snapshot.Snapshot[ScalarSnapshot] {
	return snapshot.Snapshot[ScalarSnapshot]{
		Revision: l.revision,
		Value:    NewScalarSnapshot(l.limit, l.reserved),
	}
}
