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

// TransitionError describes a rejected lifecycle transition request.
//
// TransitionError is used when a Controller cannot apply an Event from the
// current State. It carries the state/event pair for diagnostics and wraps a
// stable error class such as ErrInvalidTransition, ErrTerminalState, or
// ErrFailureCauseRequired.
//
// TransitionError intentionally does not contain a target State. For invalid
// transitions there may be no valid target state. The authoritative target for
// allowed transitions belongs to the transition table and to Transition values
// created by reduceTransition.
type TransitionError struct {
	// From is the state from which the controller attempted to apply Event.
	From State

	// Event is the lifecycle input that could not be applied.
	Event Event

	// Err is the stable cause of the rejection.
	//
	// Err SHOULD be one of the package sentinels:
	//
	//   - ErrInvalidTransition;
	//   - ErrTerminalState;
	//   - ErrFailureCauseRequired.
	//
	// If Err is nil, the error is treated as ErrInvalidTransition.
	Err error
}

// Error returns a Go-style diagnostic message for the rejected transition.
func (e *TransitionError) Error() string {
	cause := transitionErrorCause(e)

	if e == nil {
		return cause.Error()
	}

	return cause.Error() +
		": state " +
		e.From.String() +
		", event " +
		e.Event.String()
}

// Unwrap returns the stable transition rejection cause.
//
// This preserves ordinary errors.Is behavior:
//
//	errors.Is(err, lifecycle.ErrInvalidTransition)
//	errors.Is(err, lifecycle.ErrTerminalState)
//	errors.Is(err, lifecycle.ErrFailureCauseRequired)
func (e *TransitionError) Unwrap() error {
	return transitionErrorCause(e)
}

// Is reports whether target matches the transition rejection class.
//
// ErrTerminalState is a more specific form of ErrInvalidTransition. This means a
// terminal-state transition error matches both ErrTerminalState and
// ErrInvalidTransition.
func (e *TransitionError) Is(target error) bool {
	cause := transitionErrorCause(e)

	if target == ErrInvalidTransition {
		return errors.Is(cause, ErrInvalidTransition) ||
			errors.Is(cause, ErrTerminalState)
	}

	return errors.Is(cause, target)
}

// newTransitionError constructs a TransitionError with a stable non-nil cause.
//
// The helper is package-local because Controller should be the normal producer
// of transition errors. Public callers can still inspect TransitionError values
// returned by Controller.
func newTransitionError(from State, event Event, err error) *TransitionError {
	if err == nil {
		err = ErrInvalidTransition
	}

	return &TransitionError{
		From:  from,
		Event: event,
		Err:   err,
	}
}

// transitionErrorCause returns the effective cause of e.
func transitionErrorCause(e *TransitionError) error {
	if e == nil || e.Err == nil {
		return ErrInvalidTransition
	}

	return e.Err
}
