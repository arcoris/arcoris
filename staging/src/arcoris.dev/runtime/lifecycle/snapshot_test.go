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

func TestSnapshotString(t *testing.T) {
	t.Parallel()

	if got, want := (Snapshot{State: StateRunning, Revision: 7}).String(), "running@7"; got != want {
		t.Fatalf("Snapshot.String() = %q, want %q", got, want)
	}
}

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	validTransition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart, Revision: 1, At: testTime}
	failureTransition := Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: cause, Revision: 2, At: testTime}

	tests := []struct {
		name     string
		snapshot Snapshot
		want     bool
	}{
		// Snapshot is an immutable read model; the zero value represents the
		// initial StateNew observation before any transition has committed.
		{"zero snapshot", Snapshot{}, true},
		{"invalid state", Snapshot{State: State(99)}, false},
		{"revision zero non-new state", Snapshot{State: StateStarting}, false},
		{"revision zero with transition", Snapshot{LastTransition: validTransition}, false},
		{"revision zero with failure cause", Snapshot{FailureCause: cause}, false},
		{"revision non-zero zero transition", Snapshot{State: StateStarting, Revision: 1}, false},
		{"uncommitted transition", Snapshot{State: StateStarting, Revision: 1, LastTransition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}}, false},
		{"revision mismatch", Snapshot{State: StateStarting, Revision: 2, LastTransition: validTransition}, false},
		{"target mismatch", Snapshot{State: StateRunning, Revision: 1, LastTransition: validTransition}, false},
		{"valid active", Snapshot{State: StateStarting, Revision: 1, LastTransition: validTransition}, true},
		{"failed without failure cause", Snapshot{State: StateFailed, Revision: 2, LastTransition: failureTransition}, false},
		{"failed with non-failure transition", Snapshot{State: StateFailed, Revision: 2, LastTransition: Transition{From: StateStarting, To: StateRunning, Event: EventMarkRunning, Revision: 2, At: testTime}, FailureCause: cause}, false},
		{"failed transition without cause", Snapshot{State: StateFailed, Revision: 2, LastTransition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Revision: 2, At: testTime}, FailureCause: cause}, false},
		{"non-failed with failure cause", Snapshot{State: StateRunning, Revision: 1, LastTransition: Transition{From: StateStarting, To: StateRunning, Event: EventMarkRunning, Revision: 1, At: testTime}, FailureCause: cause}, false},
		{"valid failed", Snapshot{State: StateFailed, Revision: 2, LastTransition: failureTransition, FailureCause: cause}, true},
		{"non-comparable failure cause", Snapshot{State: StateFailed, Revision: 2, LastTransition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: nonComparableError{values: []string{"x"}}, Revision: 2, At: testTime}, FailureCause: nonComparableError{values: []string{"x"}}}, true},
	}

	for _, tc := range tests {
		if got := tc.snapshot.IsValid(); got != tc.want {
			t.Fatalf("%s IsValid = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestSnapshotPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		snapshot      Snapshot
		hasTransition bool
		terminal      bool
		active        bool
		acceptsWork   bool
		failed        bool
		stopped       bool
	}{
		{"new", Snapshot{State: StateNew}, false, false, false, false, false, false},
		{"running", Snapshot{State: StateRunning, Revision: 2}, true, false, true, true, false, false},
		{"stopped", Snapshot{State: StateStopped, Revision: 3}, true, true, false, false, false, true},
		{"failed", Snapshot{State: StateFailed, Revision: 3}, true, true, false, false, true, false},
	}

	for _, tc := range tests {
		if got := tc.snapshot.HasTransition(); got != tc.hasTransition {
			t.Fatalf("%s HasTransition = %v, want %v", tc.name, got, tc.hasTransition)
		}
		if got := tc.snapshot.IsTerminal(); got != tc.terminal {
			t.Fatalf("%s IsTerminal = %v, want %v", tc.name, got, tc.terminal)
		}
		if got := tc.snapshot.IsActive(); got != tc.active {
			t.Fatalf("%s IsActive = %v, want %v", tc.name, got, tc.active)
		}
		if got := tc.snapshot.AcceptsWork(); got != tc.acceptsWork {
			t.Fatalf("%s AcceptsWork = %v, want %v", tc.name, got, tc.acceptsWork)
		}
		if got := tc.snapshot.IsFailed(); got != tc.failed {
			t.Fatalf("%s IsFailed = %v, want %v", tc.name, got, tc.failed)
		}
		if got := tc.snapshot.IsStopped(); got != tc.stopped {
			t.Fatalf("%s IsStopped = %v, want %v", tc.name, got, tc.stopped)
		}
	}
}
