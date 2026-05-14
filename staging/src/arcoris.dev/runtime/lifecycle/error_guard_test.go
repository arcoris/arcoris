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
	"fmt"
	"strings"
	"testing"
)

func TestGuardErrorError(t *testing.T) {
	t.Parallel()

	// GuardError has dual semantics: it classifies as lifecycle guard rejection
	// while preserving the guard's domain-specific cause.
	domainErr := errors.New("domain blocked")
	err := &GuardError{
		Transition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart},
		Err:        domainErr,
	}
	if got, want := err.Error(), "lifecycle: guard rejected: new --begin_start--> starting: domain blocked"; got != want {
		t.Fatalf("GuardError.Error() = %q, want %q", got, want)
	}
}

func TestGuardErrorNilReceiver(t *testing.T) {
	t.Parallel()

	var err *GuardError
	if got := err.Error(); got != ErrGuardRejected.Error() {
		t.Fatalf("nil Error() = %q, want %q", got, ErrGuardRejected.Error())
	}
	mustMatch(t, err, ErrGuardRejected)
}

func TestGuardErrorNilUnderlyingError(t *testing.T) {
	t.Parallel()

	err := &GuardError{Transition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}}
	if got, want := err.Error(), "lifecycle: guard rejected: new --begin_start--> starting"; got != want {
		t.Fatalf("GuardError.Error() = %q, want %q", got, want)
	}
	if err.Unwrap() != nil {
		t.Fatalf("Unwrap = %v, want nil", err.Unwrap())
	}
}

func TestGuardErrorTrimsDuplicateGuardText(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	tests := []error{
		ErrGuardRejected,
		fmt.Errorf("%w: domain", ErrGuardRejected),
		errors.Join(ErrGuardRejected, errors.New("domain")),
	}

	for _, cause := range tests {
		err := &GuardError{Transition: transition, Err: cause}
		if strings.Count(err.Error(), ErrGuardRejected.Error()) != 1 {
			t.Fatalf("GuardError.Error() = %q, want one guard prefix", err.Error())
		}
	}
}

func TestGuardErrorUnwrapAndIs(t *testing.T) {
	t.Parallel()

	domainErr := errors.New("domain")
	err := &GuardError{Err: domainErr}
	if got := err.Unwrap(); got != domainErr {
		t.Fatalf("Unwrap = %v, want %v", got, domainErr)
	}
	mustMatch(t, err, ErrGuardRejected)
	mustMatch(t, err, domainErr)
}

func TestNewGuardErrorPreservesTransitionAndError(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	cause := errors.New("blocked")
	err := newGuardError(transition, cause)
	assertTransitionEqual(t, err.Transition, transition)
	if err.Err != cause {
		t.Fatalf("Err = %v, want %v", err.Err, cause)
	}
}
