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

// GuardError describes a transition rejected by a TransitionGuard.
//
// GuardError is used after a candidate transition has passed transition-table
// validation and required payload validation, but before Controller commits the
// transition. It carries the candidate transition and preserves the underlying
// guard-specific rejection reason.
//
// GuardError matches ErrGuardRejected through errors.Is while still unwrapping
// to the underlying guard error. This lets callers distinguish guard rejection as
// a lifecycle error class and also match domain-specific guard causes.
type GuardError struct {
	// Transition is the table-valid candidate transition rejected by a guard.
	//
	// The transition is not committed. It usually does not carry Revision or At
	// because those fields are assigned only after guards allow the transition.
	Transition Transition

	// Err is the underlying guard rejection reason.
	//
	// Err SHOULD be non-nil when a guard provides a domain-specific reason. If Err
	// is nil, the error still classifies as ErrGuardRejected.
	Err error
}

// Error returns a Go-style diagnostic message for the guard rejection.
func (e *GuardError) Error() string {
	if e == nil {
		return ErrGuardRejected.Error()
	}

	message := ErrGuardRejected.Error() + ": " + e.Transition.String()

	if e.Err == nil || errors.Is(e.Err, ErrGuardRejected) {
		return message
	}

	return message + ": " + e.Err.Error()
}

// Unwrap returns the underlying guard-specific rejection reason.
//
// ErrGuardRejected is matched through Is rather than Unwrap so the underlying
// guard error remains visible to errors.Is.
func (e *GuardError) Unwrap() error {
	if e == nil {
		return nil
	}

	if e.Err == nil || errors.Is(e.Err, ErrGuardRejected) {
		return nil
	}

	return e.Err
}

// Is reports whether target matches the guard rejection class.
func (e *GuardError) Is(target error) bool {
	return target == ErrGuardRejected
}

// newGuardError constructs a GuardError.
//
// The helper is package-local because Controller should be the normal producer
// of guard errors. It preserves the underlying guard error while ensuring the
// resulting error also matches ErrGuardRejected through errors.Is.
func newGuardError(transition Transition, err error) *GuardError {
	return &GuardError{
		Transition: transition,
		Err:        err,
	}
}
