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
	"context"
	"errors"
	"strings"
	"testing"
)

func TestTransitionErrorMatching(t *testing.T) {
	t.Parallel()

	invalid := newTransitionError(StateRunning, EventBeginStart, ErrInvalidTransition)
	mustMatch(t, invalid, ErrInvalidTransition)
	mustNotMatch(t, invalid, ErrTerminalState)

	terminal := newTransitionError(StateStopped, EventBeginStart, ErrTerminalState)
	mustMatch(t, terminal, ErrTerminalState)
	mustMatch(t, terminal, ErrInvalidTransition)
}

func TestGuardErrorMatchingAndMessage(t *testing.T) {
	t.Parallel()

	domainErr := errors.New("domain guard")
	err := newGuardError(Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}, domainErr)

	mustMatch(t, err, ErrGuardRejected)
	mustMatch(t, err, domainErr)
	if strings.Count(err.Error(), ErrGuardRejected.Error()) != 1 {
		t.Fatalf("GuardError message %q repeats guard sentinel", err.Error())
	}

	wrapped := newGuardError(
		Transition{From: StateNew, To: StateStarting, Event: EventBeginStart},
		errors.Join(ErrGuardRejected, domainErr),
	)
	mustMatch(t, wrapped, ErrGuardRejected)
	mustMatch(t, wrapped, domainErr)
	if strings.Count(wrapped.Error(), ErrGuardRejected.Error()) != 1 {
		t.Fatalf("wrapped GuardError message %q repeats guard sentinel", wrapped.Error())
	}
}

func TestWaitErrorMatchingAndMessage(t *testing.T) {
	t.Parallel()

	snapshot := Snapshot{State: StateStarting, Revision: 1, LastTransition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart, Revision: 1, At: testTime}}
	tests := []error{
		context.Canceled,
		context.DeadlineExceeded,
		ErrWaitTargetUnreachable,
		ErrInvalidWaitPredicate,
		ErrInvalidWaitTarget,
	}

	for _, target := range tests {
		err := newWaitStateError(snapshot, StateRunning, target)
		mustMatch(t, err, target)
		if !strings.HasPrefix(err.Error(), "lifecycle:") {
			t.Fatalf("WaitError message %q does not start with lifecycle:", err.Error())
		}
	}
}
