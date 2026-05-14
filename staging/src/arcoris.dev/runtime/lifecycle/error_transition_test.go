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

import (
	"errors"
	"testing"
)

func TestTransitionErrorError(t *testing.T) {
	t.Parallel()

	err := &TransitionError{From: StateRunning, Event: EventBeginStart, Err: ErrInvalidTransition}
	if got, want := err.Error(), "lifecycle: invalid transition: state running, event begin_start"; got != want {
		t.Fatalf("TransitionError.Error() = %q, want %q", got, want)
	}
}

func TestTransitionErrorNilReceiver(t *testing.T) {
	t.Parallel()

	var err *TransitionError
	if got := err.Error(); got != ErrInvalidTransition.Error() {
		t.Fatalf("nil Error() = %q, want %q", got, ErrInvalidTransition.Error())
	}
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatal("nil TransitionError does not match ErrInvalidTransition")
	}
}

func TestTransitionErrorUnwrap(t *testing.T) {
	t.Parallel()

	err := &TransitionError{Err: ErrFailureCauseRequired}
	if got := err.Unwrap(); got != ErrFailureCauseRequired {
		t.Fatalf("Unwrap = %v, want %v", got, ErrFailureCauseRequired)
	}
}

func TestTransitionErrorIs(t *testing.T) {
	t.Parallel()

	err := &TransitionError{Err: ErrFailureCauseRequired}
	mustMatch(t, err, ErrFailureCauseRequired)
	mustNotMatch(t, err, ErrInvalidTransition)
}

func TestTransitionErrorTerminalMatchesInvalidTransition(t *testing.T) {
	t.Parallel()

	// Terminal state is a specific invalid transition: callers can match either
	// the precise terminal condition or the broader invalid-transition class.
	err := &TransitionError{Err: ErrTerminalState}
	mustMatch(t, err, ErrTerminalState)
	mustMatch(t, err, ErrInvalidTransition)
}

func TestTransitionErrorNilErrDefaultsInvalidTransition(t *testing.T) {
	t.Parallel()

	err := &TransitionError{}
	if got := err.Unwrap(); got != ErrInvalidTransition {
		t.Fatalf("Unwrap = %v, want invalid transition", got)
	}
	mustMatch(t, err, ErrInvalidTransition)
}

func TestNewTransitionErrorDefaultsNilCause(t *testing.T) {
	t.Parallel()

	err := newTransitionError(StateNew, EventBeginStart, nil)
	if err.Err != ErrInvalidTransition {
		t.Fatalf("Err = %v, want invalid transition", err.Err)
	}
}
