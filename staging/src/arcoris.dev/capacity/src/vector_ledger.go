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

// VectorLedger owns strict multi-resource capacity accounting state.
//
// VectorLedger is explicit because multi-resource all-or-nothing accounting is
// heavier than scalar hot-path accounting. It uses one mutex to serialize vector
// limit changes, reservations, releases, and snapshot reads so every returned
// revisioned snapshot describes one committed vector state.
//
// The zero VectorLedger is invalid. Use NewVectorLedger. A VectorLedger must not
// be copied after first use.
type VectorLedger struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu protects revision, state, and vector reservation release ownership.
	mu sync.Mutex

	// revision is the local version of the last committed vector state.
	revision snapshot.Revision

	// state is the committed vector accounting source state.
	state VectorState
}

// VectorObservation is returned by vector observed methods.
type VectorObservation struct {
	// Snapshot is post-mutation on success and current unchanged state on refusal.
	Snapshot snapshot.Snapshot[VectorSnapshot]

	// Refusal classifies why an observed vector reservation did not acquire capacity.
	Refusal Refusal

	// Missing reports per-resource shortage for insufficient or unknown resources.
	Missing Vector

	// Debt reports per-resource debt blocking demanded resources.
	Debt Vector
}

// NewVectorLedger returns a multi-resource ledger with initial limits.
func NewVectorLedger(limits Vector) *VectorLedger {
	requireValidVector("limits", limits)

	return &VectorLedger{
		revision: snapshot.ZeroRevision.Next(),
		state: VectorState{
			Limits: limits,
		},
	}
}

// SetLimits replaces the configured vector limits.
//
// Existing vector reservations are never revoked. Lower limits may create
// per-resource debt. Setting identical limits is a no-op and does not advance
// the revision.
func (l *VectorLedger) SetLimits(limits Vector) {
	_ = l.setLimits(limits)
}

// SetLimitsObserved replaces vector limits and returns the resulting snapshot.
func (l *VectorLedger) SetLimitsObserved(limits Vector) snapshot.Snapshot[VectorSnapshot] {
	return l.setLimits(limits)
}

// Snapshot returns the current revisioned vector ledger snapshot.
func (l *VectorLedger) Snapshot() snapshot.Snapshot[VectorSnapshot] {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	return l.snapshotLocked()
}

// Revision returns the latest committed vector ledger revision.
func (l *VectorLedger) Revision() snapshot.Revision {
	l.requireNonNil()

	l.mu.Lock()
	defer l.mu.Unlock()

	l.requireInitializedLocked()

	return l.revision
}

// setLimits commits new vector limits while holding l.mu.
func (l *VectorLedger) setLimits(limits Vector) snapshot.Snapshot[VectorSnapshot] {
	requireValidVector("limits", limits)
	l.requireNonNil()

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

// snapshotLocked derives a revisioned vector snapshot from protected state.
func (l *VectorLedger) snapshotLocked() snapshot.Snapshot[VectorSnapshot] {
	return snapshot.Snapshot[VectorSnapshot]{
		Revision: l.revision,
		Value:    l.state.Snapshot(),
	}
}

// requireNonNil panics when l is nil.
func (l *VectorLedger) requireNonNil() {
	if l == nil {
		panicAt("vector_ledger", ErrNilLedger, "vector ledger receiver is nil")
	}
}

// requireInitializedLocked panics when l is a zero-value VectorLedger.
func (l *VectorLedger) requireInitializedLocked() {
	if l.revision.IsZero() {
		panicAt("vector_ledger", ErrUninitializedLedger, "vector ledger must be created with NewVectorLedger")
	}
}
