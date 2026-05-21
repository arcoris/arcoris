// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admission

import "testing"

func TestOutcomeIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		outcome Outcome
		want    bool
	}{
		{name: "unknown", outcome: OutcomeUnknown, want: false},
		{name: "admitted", outcome: OutcomeAdmitted, want: true},
		{name: "denied", outcome: OutcomeDenied, want: true},
		{name: "queued", outcome: OutcomeQueued, want: true},
		{name: "deferred", outcome: OutcomeDeferred, want: true},
		{name: "undefined", outcome: Outcome(99), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.outcome.IsValid(); got != tt.want {
				t.Fatalf("%v IsValid = %v, want %v", tt.outcome, got, tt.want)
			}
		})
	}
}

func TestOutcomeHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		outcome  Outcome
		admitted bool
		denied   bool
		queued   bool
		deferred bool
		terminal bool
	}{
		{name: "admitted", outcome: OutcomeAdmitted, admitted: true, terminal: true},
		{name: "denied", outcome: OutcomeDenied, denied: true, terminal: true},
		{name: "queued", outcome: OutcomeQueued, queued: true},
		{name: "deferred", outcome: OutcomeDeferred, deferred: true, terminal: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.outcome.IsAdmitted(); got != tt.admitted {
				t.Fatalf("IsAdmitted = %v, want %v", got, tt.admitted)
			}
			if got := tt.outcome.IsDenied(); got != tt.denied {
				t.Fatalf("IsDenied = %v, want %v", got, tt.denied)
			}
			if got := tt.outcome.IsQueued(); got != tt.queued {
				t.Fatalf("IsQueued = %v, want %v", got, tt.queued)
			}
			if got := tt.outcome.IsDeferred(); got != tt.deferred {
				t.Fatalf("IsDeferred = %v, want %v", got, tt.deferred)
			}
			if got := tt.outcome.IsTerminal(); got != tt.terminal {
				t.Fatalf("IsTerminal = %v, want %v", got, tt.terminal)
			}
		})
	}
}

func TestOutcomeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		outcome Outcome
		want    string
	}{
		{name: "unknown", outcome: OutcomeUnknown, want: "unknown"},
		{name: "admitted", outcome: OutcomeAdmitted, want: "admitted"},
		{name: "denied", outcome: OutcomeDenied, want: "denied"},
		{name: "queued", outcome: OutcomeQueued, want: "queued"},
		{name: "deferred", outcome: OutcomeDeferred, want: "deferred"},
		{name: "undefined", outcome: Outcome(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.outcome.String(); got != tt.want {
				t.Fatalf("String = %q, want %q", got, tt.want)
			}
		})
	}
}
