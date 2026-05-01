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

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Fatalf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestStatePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state       State
		valid       bool
		terminal    bool
		active      bool
		acceptsWork bool
	}{
		{StateNew, true, false, false, false},
		{StateStarting, true, false, true, false},
		{StateRunning, true, false, true, true},
		{StateStopping, true, false, true, false},
		{StateStopped, true, true, false, false},
		{StateFailed, true, true, false, false},
		{State(99), false, false, false, false},
	}

	for _, tt := range tests {
		if got := tt.state.IsValid(); got != tt.valid {
			t.Fatalf("%s IsValid = %v, want %v", tt.state, got, tt.valid)
		}
		if got := tt.state.IsTerminal(); got != tt.terminal {
			t.Fatalf("%s IsTerminal = %v, want %v", tt.state, got, tt.terminal)
		}
		if got := tt.state.IsActive(); got != tt.active {
			t.Fatalf("%s IsActive = %v, want %v", tt.state, got, tt.active)
		}
		if got := tt.state.AcceptsWork(); got != tt.acceptsWork {
			t.Fatalf("%s AcceptsWork = %v, want %v", tt.state, got, tt.acceptsWork)
		}
	}
}
