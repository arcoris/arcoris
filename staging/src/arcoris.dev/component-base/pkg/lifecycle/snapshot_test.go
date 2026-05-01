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
	"strings"
	"testing"
)

type nonComparableError struct {
	values []string
}

func (e nonComparableError) Error() string {
	return strings.Join(e.values, ":")
}

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	validTransition := Transition{From: StateNew, To: StateStarting, Event: EventBeginStart, Revision: 1, At: testTime}
	failedTransition := Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: cause, Revision: 2, At: testTime}

	tests := []struct {
		name     string
		snapshot Snapshot
		want     bool
	}{
		{name: "zero", want: true},
		{name: "invalid state", snapshot: Snapshot{State: State(99)}, want: false},
		{name: "revision zero non-zero transition", snapshot: Snapshot{LastTransition: validTransition}, want: false},
		{name: "revision zero non-comparable transition cause", snapshot: Snapshot{LastTransition: Transition{Cause: nonComparableError{values: []string{"x", "y"}}}}, want: false},
		{name: "revision non-zero zero transition", snapshot: Snapshot{State: StateStarting, Revision: 1}, want: false},
		{name: "uncommitted transition", snapshot: Snapshot{State: StateStarting, Revision: 1, LastTransition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart}}, want: false},
		{name: "revision mismatch", snapshot: Snapshot{State: StateStarting, Revision: 2, LastTransition: validTransition}, want: false},
		{name: "target mismatch", snapshot: Snapshot{State: StateRunning, Revision: 1, LastTransition: validTransition}, want: false},
		{name: "valid active", snapshot: Snapshot{State: StateStarting, Revision: 1, LastTransition: validTransition}, want: true},
		{name: "failed without cause", snapshot: Snapshot{State: StateFailed, Revision: 2, LastTransition: failedTransition}, want: false},
		{name: "non-failed with cause", snapshot: Snapshot{State: StateRunning, Revision: 1, LastTransition: Transition{From: StateStarting, To: StateRunning, Event: EventMarkRunning, Revision: 1, At: testTime}, FailureCause: cause}, want: false},
		{name: "non-failed with non-comparable cause", snapshot: Snapshot{State: StateRunning, Revision: 1, LastTransition: Transition{From: StateStarting, To: StateRunning, Event: EventMarkRunning, Revision: 1, At: testTime}, FailureCause: nonComparableError{values: []string{"x", "y"}}}, want: false},
		{name: "failed with cause", snapshot: Snapshot{State: StateFailed, Revision: 2, LastTransition: failedTransition, FailureCause: cause}, want: true},
		{name: "failed with non-comparable cause", snapshot: Snapshot{State: StateFailed, Revision: 2, LastTransition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: nonComparableError{values: []string{"x", "y"}}, Revision: 2, At: testTime}, FailureCause: nonComparableError{values: []string{"x", "y"}}}, want: true},
	}

	for _, tt := range tests {
		if got := tt.snapshot.IsValid(); got != tt.want {
			t.Fatalf("%s IsValid = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSnapshotIsValidDoesNotPanicWithNonComparableTransitionCause(t *testing.T) {
	t.Parallel()

	snapshot := Snapshot{
		State:    StateNew,
		Revision: 0,
		LastTransition: Transition{
			Cause: nonComparableError{values: []string{"x"}},
		},
	}

	if snapshot.IsValid() {
		t.Fatal("Snapshot with non-zero transition at revision zero IsValid = true, want false")
	}
}
