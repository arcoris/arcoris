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

import "testing"

func TestStateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state State
		want  string
	}{
		{StateNew, "new"},
		{StateStarting, "starting"},
		{StateRunning, "running"},
		{StateStopping, "stopping"},
		{StateStopped, "stopped"},
		{StateFailed, "failed"},
		{State(99), "invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()

			if got := tc.state.String(); got != tc.want {
				t.Fatalf("State(%d).String() = %q, want %q", tc.state, got, tc.want)
			}
		})
	}
}

func TestStateIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state State
		want  bool
	}{
		{"new", StateNew, true},
		{"starting", StateStarting, true},
		{"running", StateRunning, true},
		{"stopping", StateStopping, true},
		{"stopped", StateStopped, true},
		{"failed", StateFailed, true},
		{"invalid", State(99), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.state.IsValid(); got != tc.want {
				t.Fatalf("%s IsValid = %v, want %v", tc.state, got, tc.want)
			}
		})
	}
}

func TestStateIsTerminal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state State
		want  bool
	}{
		{StateNew, false},
		{StateStarting, false},
		{StateRunning, false},
		{StateStopping, false},
		{StateStopped, true},
		{StateFailed, true},
		{State(99), false},
	}

	for _, tc := range tests {
		if got := tc.state.IsTerminal(); got != tc.want {
			t.Fatalf("%s IsTerminal = %v, want %v", tc.state, got, tc.want)
		}
	}
}

func TestStateIsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state State
		want  bool
	}{
		{StateNew, false},
		{StateStarting, true},
		{StateRunning, true},
		{StateStopping, true},
		{StateStopped, false},
		{StateFailed, false},
		{State(99), false},
	}

	for _, tc := range tests {
		if got := tc.state.IsActive(); got != tc.want {
			t.Fatalf("%s IsActive = %v, want %v", tc.state, got, tc.want)
		}
	}
}

func TestStateAcceptsWork(t *testing.T) {
	t.Parallel()

	// Only StateRunning accepts normal workload; transitional states may be busy
	// starting or draining and terminal states cannot accept new lifecycle work.
	for _, state := range append(append([]State(nil), allStates...), State(99)) {
		want := state == StateRunning
		if got := state.AcceptsWork(); got != want {
			t.Fatalf("%s AcceptsWork = %v, want %v", state, got, want)
		}
	}
}

func TestStateZeroValueInvariant(t *testing.T) {
	t.Parallel()

	// StateNew is intentionally the valid zero value so embedding structs do not
	// need constructor-only initialization before exposing lifecycle state.
	var state State
	if state != StateNew {
		t.Fatalf("zero State = %s, want new", state)
	}
	if !state.IsValid() {
		t.Fatal("zero State IsValid = false, want true")
	}
}
