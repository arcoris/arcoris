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

import "errors"

var (
	// ErrInvalidTransition reports that a lifecycle event cannot be applied to a
	// source state according to the lifecycle transition table.
	//
	// This is a table-level error. It means the requested state/event pair is not
	// part of the lifecycle graph. It does not mean that a guard rejected the
	// transition or that a valid transition was missing required runtime payload.
	ErrInvalidTransition = errors.New("lifecycle: invalid transition")

	// ErrTerminalState reports that a lifecycle transition was attempted from a
	// terminal state.
	//
	// StateStopped and StateFailed are terminal. A lifecycle instance that reaches
	// either state MUST NOT be reused. Restart orchestration belongs to an owner
	// or supervisor that creates a fresh lifecycle instance.
	//
	// ErrTerminalState is treated as a more specific form of invalid transition
	// by TransitionError.
	ErrTerminalState = errors.New("lifecycle: terminal state")

	// ErrFailureCauseRequired reports that a failed lifecycle transition was
	// requested without a non-nil failure cause.
	//
	// EventMarkFailed transitions MUST carry a cause so a failed lifecycle cannot
	// lose the reason that made it terminal.
	ErrFailureCauseRequired = errors.New("lifecycle: failure cause required")

	// ErrGuardRejected reports that a transition guard rejected an otherwise
	// table-valid candidate transition.
	//
	// Guard rejection is distinct from invalid transition. The lifecycle graph may
	// allow the transition, but an owner-specific precondition prevented the
	// controller from committing it.
	ErrGuardRejected = errors.New("lifecycle: guard rejected")

	// ErrInvalidWaitPredicate reports that a wait operation received an invalid
	// predicate.
	//
	// Wait predicates are caller-provided read-side conditions. A nil predicate
	// cannot be evaluated and should be rejected before registering waiters.
	ErrInvalidWaitPredicate = errors.New("lifecycle: invalid wait predicate")

	// ErrInvalidWaitTarget reports that a wait operation received an invalid
	// target state.
	//
	// WaitState should use this error when the requested target is not one of the
	// lifecycle states defined by this package.
	ErrInvalidWaitTarget = errors.New("lifecycle: invalid wait target")

	// ErrWaitTargetUnreachable reports that a wait target can no longer be reached
	// from the current lifecycle state.
	//
	// The common case is waiting for an active state such as StateRunning while
	// the lifecycle reaches StateStopped or StateFailed first.
	ErrWaitTargetUnreachable = errors.New("lifecycle: wait target unreachable")
)
