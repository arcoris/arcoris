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

func TestTransitionPredicates(t *testing.T) {
	t.Parallel()

	cause := errors.New("failed")
	tests := []struct {
		name       string
		transition Transition
		valid      bool
		tableValid bool
		committed  bool
		terminal   bool
		failure    bool
		zero       bool
	}{
		{name: "zero", zero: true},
		{
			name:       "candidate",
			transition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart},
			valid:      true,
			tableValid: true,
		},
		{
			name:       "committed",
			transition: Transition{From: StateNew, To: StateStarting, Event: EventBeginStart, Revision: 1, At: testTime},
			valid:      true,
			tableValid: true,
			committed:  true,
		},
		{
			name:       "terminal failure",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: cause, Revision: 1, At: testTime},
			valid:      true,
			tableValid: true,
			committed:  true,
			terminal:   true,
			failure:    true,
		},
		{
			name:       "non-comparable cause",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed, Cause: nonComparableError{values: []string{"x", "y"}}},
			valid:      true,
			tableValid: true,
			terminal:   true,
			failure:    true,
		},
		{
			name:       "missing cause",
			transition: Transition{From: StateRunning, To: StateFailed, Event: EventMarkFailed},
			tableValid: true,
			terminal:   true,
			failure:    true,
		},
		{
			name:       "invalid table",
			transition: Transition{From: StateRunning, To: StateStarting, Event: EventBeginStart},
		},
	}

	for _, tt := range tests {
		if got := tt.transition.IsValid(); got != tt.valid {
			t.Fatalf("%s IsValid = %v, want %v", tt.name, got, tt.valid)
		}
		if got := tt.transition.IsTableValid(); got != tt.tableValid {
			t.Fatalf("%s IsTableValid = %v, want %v", tt.name, got, tt.tableValid)
		}
		if got := tt.transition.IsCommitted(); got != tt.committed {
			t.Fatalf("%s IsCommitted = %v, want %v", tt.name, got, tt.committed)
		}
		if got := tt.transition.IsTerminal(); got != tt.terminal {
			t.Fatalf("%s IsTerminal = %v, want %v", tt.name, got, tt.terminal)
		}
		if got := tt.transition.IsFailure(); got != tt.failure {
			t.Fatalf("%s IsFailure = %v, want %v", tt.name, got, tt.failure)
		}
		if got := tt.transition.IsZero(); got != tt.zero {
			t.Fatalf("%s IsZero = %v, want %v", tt.name, got, tt.zero)
		}
	}
}
