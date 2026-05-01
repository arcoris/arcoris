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
	"fmt"
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
	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "domain error only",
			err:  newGuardError(transition, domainErr),
		},
		{
			name: "guard sentinel only",
			err:  newGuardError(transition, ErrGuardRejected),
		},
		{
			name: "join guard first",
			err:  newGuardError(transition, errors.Join(ErrGuardRejected, domainErr)),
		},
		{
			name: "join domain first",
			err:  newGuardError(transition, errors.Join(domainErr, ErrGuardRejected)),
		},
		{
			name: "wrapped guard sentinel",
			err:  newGuardError(transition, fmt.Errorf("wrapped: %w", ErrGuardRejected)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "guard sentinel only" {
				mustMatch(t, tt.err, ErrGuardRejected)
			}
			if tt.name != "guard sentinel only" && tt.name != "wrapped guard sentinel" {
				mustMatch(t, tt.err, domainErr)
			}
			if !strings.HasPrefix(tt.err.Error(), ErrGuardRejected.Error()) {
				t.Fatalf("GuardError message %q does not start with %q", tt.err.Error(), ErrGuardRejected.Error())
			}
			if strings.Count(tt.err.Error(), ErrGuardRejected.Error()) > 1 && tt.name != "join domain first" {
				t.Fatalf("GuardError message %q repeats guard sentinel", tt.err.Error())
			}
		})
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
