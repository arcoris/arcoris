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
	"time"
)

func TestTransitionString(t *testing.T) {
	t.Parallel()

	transition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	if got, want := transition.String(), "new --begin_start--> starting"; got != want {
		t.Fatalf("Transition.String() = %q, want %q", got, want)
	}
}

func TestTransitionIsValid(t *testing.T) {
	t.Parallel()

	// Candidate transitions are table-valid runtime facts before Controller adds
	// commit metadata; failure candidates additionally need a cause.
	cause := errors.New("failed")
	tests := []struct {
		name       string
		transition Transition
		want       bool
	}{
		{
			name:       "candidate allowed",
			transition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart},
			want:       true,
		},
		{
			name:       "failure with cause",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: cause},
			want:       true,
		},
		{
			name:       "missing failure cause",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed},
			want:       false,
		},
		{
			name:       "non-comparable failure cause",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: nonComparableError{values: []string{"x", "y"}}},
			want:       true,
		},
		{
			name:       "invalid from",
			transition: Transition{From: State(99), To: StateStarting, Event: EventBeginStart},
		},
		{
			name:       "invalid to",
			transition: Transition{From: StateNew, To: State(99), Event: EventBeginStart},
		},
		{
			name:       "invalid event",
			transition: Transition{From: StateNew, To: StateStarting, Event: Event(99)},
		},
		{
			name:       "invalid combination",
			transition: Transition{From: StateRunning, To: StateStarting, Event: EventBeginStart},
		},
	}

	for _, tc := range tests {
		if got := tc.transition.IsValid(); got != tc.want {
			t.Fatalf("%s IsValid = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestTransitionIsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		transition Transition
		want       bool
	}{
		{"zero", Transition{}, true},
		{"state set", Transition{From: StateStarting}, false},
		{"non-comparable cause", Transition{Cause: nonComparableError{values: []string{"x"}}}, false},
	}

	for _, tc := range tests {
		if got := tc.transition.IsZero(); got != tc.want {
			t.Fatalf("%s IsZero = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestTransitionIsTableValid(t *testing.T) {
	t.Parallel()

	// Table validity describes the static graph shape. Runtime payload rules,
	// such as MarkFailed requiring a cause, are checked by IsValid.
	transition := Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed}
	if !transition.IsTableValid() {
		t.Fatal("failure transition without cause IsTableValid = false, want true")
	}
	if transition.IsValid() {
		t.Fatal("failure transition without cause IsValid = true, want false")
	}
}

func TestTransitionIsCommitted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		transition Transition
		want       bool
	}{
		{Transition{Revision: 1, At: testTime}, true},
		{Transition{Revision: 1}, false},
		{Transition{At: testTime}, false},
		{Transition{}, false},
	}

	for _, tc := range tests {
		if got := tc.transition.IsCommitted(); got != tc.want {
			t.Fatalf("%+v IsCommitted = %v, want %v", tc.transition, got, tc.want)
		}
	}
}

func TestTransitionClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		transition Transition
		terminal   bool
		failure    bool
		startup    bool
		shutdown   bool
	}{
		{"start", Transition{To: StateStarting, Event: EventBeginStart}, false, false, true, false},
		{"running", Transition{To: StateRunning, Event: EventMarkRunning}, false, false, true, false},
		{"begin stop", Transition{To: StateStopping, Event: EventBeginStop}, false, false, false, true},
		{"stopped", Transition{To: StateStopped, Event: EventMarkStopped}, true, false, false, true},
		{"failed", Transition{To: StateFailed, Event: EventMarkFailed}, true, true, false, false},
		{"inconsistent failure", Transition{To: StateRunning, Event: EventMarkFailed}, false, false, false, false},
	}

	for _, tc := range tests {
		if got := tc.transition.IsTerminal(); got != tc.terminal {
			t.Fatalf("%s IsTerminal = %v, want %v", tc.name, got, tc.terminal)
		}
		if got := tc.transition.IsFailure(); got != tc.failure {
			t.Fatalf("%s IsFailure = %v, want %v", tc.name, got, tc.failure)
		}
		if got := tc.transition.IsStartup(); got != tc.startup {
			t.Fatalf("%s IsStartup = %v, want %v", tc.name, got, tc.startup)
		}
		if got := tc.transition.IsShutdown(); got != tc.shutdown {
			t.Fatalf("%s IsShutdown = %v, want %v", tc.name, got, tc.shutdown)
		}
	}
}

func TestTransitionWithCommitMetadata(t *testing.T) {
	t.Parallel()

	// Committing returns a copy with metadata; the candidate remains suitable for
	// guard evaluation without Revision or At.
	candidate := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}
	committed := candidate.withCommitMetadata(3, testTime)

	if !committed.IsCommitted() {
		t.Fatal("committed transition IsCommitted = false, want true")
	}
	if committed.Revision != 3 || !committed.At.Equal(testTime) {
		t.Fatalf("commit metadata = revision %d at %v", committed.Revision, committed.At)
	}
	if candidate.Revision != 0 || !candidate.At.Equal(time.Time{}) {
		t.Fatalf("candidate mutated to %+v", candidate)
	}
}
