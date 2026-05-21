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

package retry

import "testing"

func TestEventKindString(t *testing.T) {
	tests := []struct {
		name string
		kind EventKind
		want string
	}{
		{
			name: "attempt start",
			kind: EventAttemptStart,
			want: "attempt_start",
		},
		{
			name: "attempt failure",
			kind: EventAttemptFailure,
			want: "attempt_failure",
		},
		{
			name: "retry delay",
			kind: EventRetryDelay,
			want: "retry_delay",
		},
		{
			name: "retry stop",
			kind: EventRetryStop,
			want: "retry_stop",
		},
		{
			name: "zero",
			kind: 0,
			want: "invalid",
		},
		{
			name: "unknown",
			kind: EventKind(255),
			want: "invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.String(); got != tc.want {
				t.Fatalf("EventKind.String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEventKindIsValid(t *testing.T) {
	valid := []EventKind{
		EventAttemptStart,
		EventAttemptFailure,
		EventRetryDelay,
		EventRetryStop,
	}

	for _, kind := range valid {
		t.Run(kind.String(), func(t *testing.T) {
			if !kind.IsValid() {
				t.Fatalf("%s IsValid() = false, want true", kind)
			}
		})
	}

	invalid := []EventKind{
		0,
		EventKind(255),
	}

	for _, kind := range invalid {
		t.Run(kind.String(), func(t *testing.T) {
			if kind.IsValid() {
				t.Fatalf("%s IsValid() = true, want false", kind)
			}
		})
	}
}

func TestEventKindIsAttemptScoped(t *testing.T) {
	tests := []struct {
		name string
		kind EventKind
		want bool
	}{
		{
			name: "attempt start",
			kind: EventAttemptStart,
			want: true,
		},
		{
			name: "attempt failure",
			kind: EventAttemptFailure,
			want: true,
		},
		{
			name: "retry delay",
			kind: EventRetryDelay,
			want: true,
		},
		{
			name: "retry stop",
			kind: EventRetryStop,
			want: false,
		},
		{
			name: "invalid",
			kind: 0,
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.IsAttemptScoped(); got != tc.want {
				t.Fatalf("EventKind.IsAttemptScoped() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEventKindIsTerminal(t *testing.T) {
	tests := []struct {
		name string
		kind EventKind
		want bool
	}{
		{
			name: "attempt start",
			kind: EventAttemptStart,
			want: false,
		},
		{
			name: "attempt failure",
			kind: EventAttemptFailure,
			want: false,
		},
		{
			name: "retry delay",
			kind: EventRetryDelay,
			want: false,
		},
		{
			name: "retry stop",
			kind: EventRetryStop,
			want: true,
		},
		{
			name: "invalid",
			kind: 0,
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.IsTerminal(); got != tc.want {
				t.Fatalf("EventKind.IsTerminal() = %v, want %v", got, tc.want)
			}
		})
	}
}
