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

func TestTransitionRuleString(t *testing.T) {
	t.Parallel()

	rule := TransitionRule{From: StateNew, Event: EventBeginStart, To: StateStarting}
	if got, want := rule.String(), "new --begin_start--> starting"; got != want {
		t.Fatalf("TransitionRule.String() = %q, want %q", got, want)
	}
}

func TestTransitionRuleIsValid(t *testing.T) {
	t.Parallel()

	// TransitionRule is static table data: no revision, timestamp, or runtime
	// failure cause participates in validity.
	tests := []struct {
		name string
		rule TransitionRule
		want bool
	}{
		{"valid", TransitionRule{StateNew, EventBeginStart, StateStarting}, true},
		{"invalid from", TransitionRule{State(99), EventBeginStart, StateStarting}, false},
		{"invalid event", TransitionRule{StateNew, Event(99), StateStarting}, false},
		{"invalid target", TransitionRule{StateNew, EventBeginStart, State(99)}, false},
		{"invalid combination", TransitionRule{StateRunning, EventBeginStart, StateStarting}, false},
	}

	for _, tt := range tests {
		if got := tt.rule.IsValid(); got != tt.want {
			t.Fatalf("%s IsValid = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTransitionRuleMatches(t *testing.T) {
	t.Parallel()

	rule := TransitionRule{From: StateStarting, Event: EventBeginStop, To: StateStopping}
	if !rule.Matches(StateStarting, EventBeginStop) {
		t.Fatal("rule did not match its source event pair")
	}
	if rule.Matches(StateRunning, EventBeginStop) {
		t.Fatal("rule matched a different source state")
	}
}

func TestTransitionRuleClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		rule     TransitionRule
		terminal bool
		failure  bool
	}{
		{"active", TransitionRule{StateNew, EventBeginStart, StateStarting}, false, false},
		{"stopped", TransitionRule{StateStopping, EventMarkStopped, StateStopped}, true, false},
		{"failed", TransitionRule{StateRunning, EventMarkFailed, StateFailed}, true, true},
	}

	for _, tt := range tests {
		if got := tt.rule.IsTerminal(); got != tt.terminal {
			t.Fatalf("%s IsTerminal = %v, want %v", tt.name, got, tt.terminal)
		}
		if got := tt.rule.IsFailure(); got != tt.failure {
			t.Fatalf("%s IsFailure = %v, want %v", tt.name, got, tt.failure)
		}
	}
}

func TestTransitionRuleStaticRulesAreValid(t *testing.T) {
	t.Parallel()

	for _, rule := range expectedTransitionRules {
		if !rule.IsValid() {
			t.Fatalf("static rule %s IsValid = false, want true", rule)
		}
	}
}
