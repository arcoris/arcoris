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

// Ledger owns local multi-resource capacity accounting state.
//
// Ledger serializes limit changes, reservations, releases, and snapshot reads
// with one mutex so every returned revisioned snapshot describes one committed
// state. It performs all-or-nothing reservations and records debt instead of
// revoking existing reservations when limits shrink.
//
// The zero Ledger is invalid. Use NewLedger. A Ledger must not be copied after
// first use.
type Ledger struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects revision, state, and reservation release ownership.
	mu sync.Mutex

	// revision is the local version of the last committed state.
	revision snapshot.Revision

	// state is the committed accounting source state.
	state State
}

// NewLedger returns a multi-resource ledger with initial limits.
func NewLedger(limits Vector) *Ledger {
	if !limits.IsValid() {
		panicAt("limits", ErrInvalidVector, ErrorReasonInvalidVector, "limits vector must be canonical")
	}

	return &Ledger{
		revision: snapshot.ZeroRevision.Next(),
		state: State{
			Limits: limits,
		},
	}
}

// SetLimits replaces the ledger's configured resource limits.
//
// Existing reservations are never revoked. Lower limits may create per-resource
// debt. Setting identical limits is a no-op and does not advance the revision.
func (l *Ledger) SetLimits(limits Vector) snapshot.Snapshot[Snapshot] {
	l.requireNonNil()
	if !limits.IsValid() {
		panicAt("limits", ErrInvalidVector, ErrorReasonInvalidVector, "limits vector must be canonical")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	if l.state.Limits.Equal(limits) {
		return l.snapshotLocked()
	}

	l.state.Limits = limits
	l.revision = l.revision.Next()
	return l.snapshotLocked()
}

// Snapshot returns the current revisioned multi-resource ledger snapshot.
func (l *Ledger) Snapshot() snapshot.Snapshot[Snapshot] {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()
	return l.snapshotLocked()
}

// Revision returns the latest committed ledger revision.
func (l *Ledger) Revision() snapshot.Revision {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	return l.revision
}

// snapshotLocked derives a revisioned snapshot for state protected by l.mu.
func (l *Ledger) snapshotLocked() snapshot.Snapshot[Snapshot] {
	return snapshot.Snapshot[Snapshot]{
		Revision: l.revision,
		Value:    l.state.Snapshot(),
	}
}
