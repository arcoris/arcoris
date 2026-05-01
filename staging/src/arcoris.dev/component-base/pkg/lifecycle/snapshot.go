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

package lifecycle

import "strconv"

// Snapshot is an immutable point-in-time view of a lifecycle controller.
//
// A Snapshot is a read model, not mutable lifecycle state. It records the current
// lifecycle state together with the last committed transition metadata needed for
// diagnostics, waits, observers, tests, and controller integrations.
//
// Snapshot is intentionally copyable. It does not contain locks, waiters,
// pointers back to Controller, observers, guards, runtime handles, or ownership
// markers. Copying a Snapshot copies an observation; it does not split or mutate
// controller state.
//
// Controller is responsible for creating snapshots under its internal
// synchronization so State, Revision, LastTransition, and FailureCause describe
// one consistent point in the lifecycle timeline.
type Snapshot struct {
	// State is the current lifecycle state observed by the snapshot.
	//
	// State MUST be valid. The zero value is StateNew, which makes the zero
	// Snapshot a valid initial lifecycle snapshot.
	State State

	// Revision is the monotonic committed-transition revision observed by the
	// snapshot.
	//
	// Revision is zero before the first committed transition. After each committed
	// transition, Controller increments Revision and stores the same value in the
	// committed Transition.
	Revision uint64

	// LastTransition is the most recent committed lifecycle transition.
	//
	// LastTransition is the zero value before the first committed transition. When
	// Revision is non-zero, LastTransition SHOULD be a committed transition whose
	// Revision equals Snapshot.Revision and whose To state equals Snapshot.State.
	LastTransition Transition

	// FailureCause is the terminal failure cause for a failed lifecycle.
	//
	// FailureCause MUST be non-nil when State is StateFailed. It MUST be nil for
	// all non-failed states. The cause is stored separately from LastTransition so
	// callers waiting on terminal state can inspect failure without unpacking the
	// last transition.
	FailureCause error
}

// String returns a compact diagnostic representation of s.
//
// The returned value is intended for logs, tests, diagnostics, and error
// messages. It is not a stable serialization format.
func (s Snapshot) String() string {
	return s.State.String() + "@" + strconv.FormatUint(s.Revision, 10)
}

// IsValid reports whether s satisfies the lifecycle snapshot invariants.
//
// The zero Snapshot is valid and represents the initial StateNew snapshot before
// any transition has been committed.
func (s Snapshot) IsValid() bool {
	if !s.State.IsValid() {
		return false
	}

	if s.Revision == 0 {
		return s.State == StateNew &&
			s.FailureCause == nil &&
			s.LastTransition == (Transition{})
	}

	if s.State == StateFailed {
		if s.FailureCause == nil {
			return false
		}
	} else if s.FailureCause != nil {
		return false
	}

	if !s.LastTransition.IsValid() {
		return false
	}

	if !s.LastTransition.IsCommitted() {
		return false
	}

	if s.LastTransition.Revision != s.Revision {
		return false
	}

	if s.LastTransition.To != s.State {
		return false
	}

	if s.State == StateFailed {
		if !s.LastTransition.IsFailure() {
			return false
		}

		if s.LastTransition.Cause == nil {
			return false
		}
	}

	return true
}

// HasTransition reports whether the snapshot includes at least one committed
// lifecycle transition.
//
// A snapshot with Revision zero is the initial lifecycle snapshot. A snapshot
// with a non-zero Revision represents a lifecycle that has committed at least one
// transition.
func (s Snapshot) HasTransition() bool {
	return s.Revision != 0
}

// IsTerminal reports whether the snapshot state ends the lifecycle instance.
func (s Snapshot) IsTerminal() bool {
	return s.State.IsTerminal()
}

// IsActive reports whether the snapshot state belongs to a lifecycle that has
// started but has not reached a terminal state.
func (s Snapshot) IsActive() bool {
	return s.State.IsActive()
}

// AcceptsWork reports whether a component in the snapshot state may accept
// normal workload by default.
//
// AcceptsWork delegates to State.AcceptsWork. A component with stricter
// admission, readiness, health, or workload-specific policy may still reject work
// while its lifecycle state is StateRunning.
func (s Snapshot) AcceptsWork() bool {
	return s.State.AcceptsWork()
}

// IsFailed reports whether the snapshot describes an unsuccessful terminal
// lifecycle.
func (s Snapshot) IsFailed() bool {
	return s.State == StateFailed
}

// IsStopped reports whether the snapshot describes a successful terminal
// lifecycle.
func (s Snapshot) IsStopped() bool {
	return s.State == StateStopped
}
