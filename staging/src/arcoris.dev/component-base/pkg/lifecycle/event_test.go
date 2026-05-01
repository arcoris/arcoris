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

func TestEventString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		event Event
		want  string
	}{
		{EventBeginStart, "begin_start"},
		{EventMarkRunning, "mark_running"},
		{EventBeginStop, "begin_stop"},
		{EventMarkStopped, "mark_stopped"},
		{EventMarkFailed, "mark_failed"},
		{Event(99), "invalid"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			if got := tt.event.String(); got != tt.want {
				t.Fatalf("Event(%d).String() = %q, want %q", tt.event, got, tt.want)
			}
		})
	}
}

func TestEventIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		event Event
		want  bool
	}{
		{EventBeginStart, true},
		{EventMarkRunning, true},
		{EventBeginStop, true},
		{EventMarkStopped, true},
		{EventMarkFailed, true},
		{Event(99), false},
	}

	for _, tt := range tests {
		if got := tt.event.IsValid(); got != tt.want {
			t.Fatalf("%s IsValid = %v, want %v", tt.event, got, tt.want)
		}
	}
}

func TestEventPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		event         Event
		start         bool
		stop          bool
		terminal      bool
		requiresCause bool
	}{
		{EventBeginStart, true, false, false, false},
		{EventMarkRunning, true, false, false, false},
		{EventBeginStop, false, true, false, false},
		{EventMarkStopped, false, true, true, false},
		{EventMarkFailed, false, false, true, true},
		{Event(99), false, false, false, false},
	}

	for _, tt := range tests {
		// Events are lifecycle inputs, not target states: failure is terminal but
		// not a stop event, and BeginStop may target different states by source.
		if got := tt.event.IsStartEvent(); got != tt.start {
			t.Fatalf("%s IsStartEvent = %v, want %v", tt.event, got, tt.start)
		}
		if got := tt.event.IsStopEvent(); got != tt.stop {
			t.Fatalf("%s IsStopEvent = %v, want %v", tt.event, got, tt.stop)
		}
		if got := tt.event.IsTerminalEvent(); got != tt.terminal {
			t.Fatalf("%s IsTerminalEvent = %v, want %v", tt.event, got, tt.terminal)
		}
		if got := tt.event.RequiresCause(); got != tt.requiresCause {
			t.Fatalf("%s RequiresCause = %v, want %v", tt.event, got, tt.requiresCause)
		}
	}
}
