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

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()

			if got := tc.event.String(); got != tc.want {
				t.Fatalf("Event(%d).String() = %q, want %q", tc.event, got, tc.want)
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

	for _, tc := range tests {
		if got := tc.event.IsValid(); got != tc.want {
			t.Fatalf("%s IsValid = %v, want %v", tc.event, got, tc.want)
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

	for _, tc := range tests {
		// Events are lifecycle inputs, not target states: failure is terminal but
		// not a stop event, and BeginStop may target different states by source.
		if got := tc.event.IsStartEvent(); got != tc.start {
			t.Fatalf("%s IsStartEvent = %v, want %v", tc.event, got, tc.start)
		}
		if got := tc.event.IsStopEvent(); got != tc.stop {
			t.Fatalf("%s IsStopEvent = %v, want %v", tc.event, got, tc.stop)
		}
		if got := tc.event.IsTerminalEvent(); got != tc.terminal {
			t.Fatalf("%s IsTerminalEvent = %v, want %v", tc.event, got, tc.terminal)
		}
		if got := tc.event.RequiresCause(); got != tc.requiresCause {
			t.Fatalf("%s RequiresCause = %v, want %v", tc.event, got, tc.requiresCause)
		}
	}
}
