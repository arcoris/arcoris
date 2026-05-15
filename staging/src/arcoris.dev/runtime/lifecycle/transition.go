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

import "time"

// Transition describes one lifecycle state change.
//
// A Transition is a runtime value object. It records the state that a lifecycle
// instance moved from, the state it moved to, and the event that caused the
// movement. It may also carry commit metadata assigned by Controller after the
// transition is accepted and committed.
//
// Transition deliberately does not execute lifecycle work. It does not start a
// component, stop goroutines, close resources, retry failed operations, emit
// metrics, or notify observers. Those responsibilities belong to controller,
// run, retry, health, and observability layers.
//
// The zero value is not a valid transition.
type Transition struct {
	// From is the lifecycle state observed before the transition was applied.
	//
	// From MUST be a valid State. A transition from a terminal state is not valid
	// in the default lifecycle model because StateStopped and StateFailed end the
	// lifecycle instance.
	From State

	// To is the lifecycle state produced by applying Event to From.
	//
	// To MUST be the exact state selected by the lifecycle transition table. A
	// caller must not construct arbitrary state movements by setting To directly.
	// Use NextState or reduceTransition to derive the target state.
	To State

	// Event is the lifecycle input that caused this transition.
	//
	// Event describes what happened or what the owner is attempting to record. It
	// is not merely an alias for To. For example, EventBeginStop may move a
	// lifecycle from StateNew directly to StateStopped, or from StateRunning to
	// StateStopping, depending on the current state.
	Event Event

	// Revision is the monotonic commit sequence number assigned by Controller.
	//
	// A zero Revision means the transition has not been committed by a controller
	// yet. Candidate transitions produced during validation may have zero
	// Revision. Committed transitions SHOULD use one-based revisions so the zero
	// value remains reserved for "not committed".
	Revision uint64

	// At is the time at which Controller committed the transition.
	//
	// A zero At means the transition has not been committed by a controller yet.
	// Candidate transitions produced by the pure transition reducer may leave At
	// unset. Controller is responsible for assigning time through the configured
	// clock source.
	At time.Time

	// Cause carries the failure cause for transitions that mark the lifecycle as
	// failed.
	//
	// EventMarkFailed transitions MUST carry a non-nil Cause. Normal lifecycle
	// transitions should not rely on Cause for control flow. If startup, runtime,
	// or shutdown cannot complete successfully, the owner should record
	// EventMarkFailed instead of attaching an error to a successful transition.
	Cause error
}

// String returns a compact diagnostic representation of t.
//
// The returned value is intended for logs, tests, diagnostics, and error
// messages. It is not a stable serialization format.
func (t Transition) String() string {
	return t.From.String() + " --" + t.Event.String() + "--> " + t.To.String()
}

// IsValid reports whether t is a valid lifecycle transition.
//
// IsValid checks both table-level validity and runtime payload validity:
//
//   - From is a known state;
//   - To is a known state;
//   - Event is a known event;
//   - applying Event to From produces To;
//   - EventMarkFailed carries a non-nil Cause.
//
// IsValid does not require Revision or At to be set. Candidate transitions may
// be valid before they are committed by Controller.
func (t Transition) IsValid() bool {
	if !t.IsTableValid() {
		return false
	}

	if t.Event.RequiresCause() && t.Cause == nil {
		return false
	}

	return true
}

// IsZero reports whether t is the zero transition value.
//
// The check is field-by-field so it remains safe when Cause holds an error whose
// dynamic type is not comparable.
func (t Transition) IsZero() bool {
	return t.From == 0 &&
		t.To == 0 &&
		t.Event == 0 &&
		t.Revision == 0 &&
		t.At.IsZero() &&
		t.Cause == nil
}

// IsTableValid reports whether t matches the lifecycle transition table without
// checking runtime payload requirements such as failure cause.
//
// This method is useful when code needs to distinguish "the state/event pair is
// not allowed" from "the transition shape is allowed, but the controller must
// still reject it because required runtime data is missing".
func (t Transition) IsTableValid() bool {
	if !t.From.IsValid() || !t.To.IsValid() || !t.Event.IsValid() {
		return false
	}

	to, ok := NextState(t.From, t.Event)
	return ok && to == t.To
}

// IsCommitted reports whether t carries controller-assigned commit metadata.
//
// Candidate transitions produced by the pure transition reducer are not
// committed. A committed transition SHOULD have a non-zero Revision and a
// non-zero commit time.
func (t Transition) IsCommitted() bool {
	return t.Revision != 0 && !t.At.IsZero()
}

// IsTerminal reports whether t moves the lifecycle into a terminal state.
//
// Terminal transitions end the lifecycle instance. The default terminal states
// are StateStopped and StateFailed.
func (t Transition) IsTerminal() bool {
	return t.To.IsTerminal()
}

// IsFailure reports whether t is a failure transition.
//
// A failure transition is strict: it must be caused by EventMarkFailed and must
// move the lifecycle to StateFailed. Partially constructed or inconsistent
// transitions should not be treated as failure transitions merely because one
// field happens to mention failure.
func (t Transition) IsFailure() bool {
	return t.Event == EventMarkFailed && t.To == StateFailed
}

// IsStartup reports whether t belongs to the normal startup side of the
// lifecycle.
//
// Startup transitions begin with EventBeginStart and complete with
// EventMarkRunning. A failure during startup is still represented as a failure
// transition, not as a startup transition.
func (t Transition) IsStartup() bool {
	return t.Event.IsStartEvent()
}

// IsShutdown reports whether t belongs to the normal shutdown side of the
// lifecycle.
//
// Shutdown transitions begin with EventBeginStop and complete with
// EventMarkStopped. A failure during shutdown is still represented as a failure
// transition, not as a normal shutdown transition.
func (t Transition) IsShutdown() bool {
	return t.Event.IsStopEvent()
}

// withCommitMetadata returns a copy of t with controller-assigned commit
// metadata.
//
// Controller should call this only after the transition has passed table
// validation, failure-cause validation, and guard validation. The method is
// intentionally package-local so ordinary callers cannot fabricate committed
// transitions.
func (t Transition) withCommitMetadata(rev uint64, at time.Time) Transition {
	t.Revision = rev
	t.At = at
	return t
}
