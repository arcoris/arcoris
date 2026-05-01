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
	"reflect"
	"testing"
)

var allowedRules = []TransitionRule{
	{StateNew, EventBeginStart, StateStarting},
	{StateNew, EventBeginStop, StateStopped},
	{StateStarting, EventMarkRunning, StateRunning},
	{StateStarting, EventBeginStop, StateStopping},
	{StateStarting, EventMarkFailed, StateFailed},
	{StateRunning, EventBeginStop, StateStopping},
	{StateRunning, EventMarkFailed, StateFailed},
	{StateStopping, EventMarkStopped, StateStopped},
	{StateStopping, EventMarkFailed, StateFailed},
}

func TestTransitionTableAllowedRules(t *testing.T) {
	t.Parallel()

	for _, rule := range allowedRules {
		to, ok := NextState(rule.From, rule.Event)
		if !ok {
			t.Fatalf("NextState(%s, %s) ok = false, want true", rule.From, rule.Event)
		}
		if to != rule.To {
			t.Fatalf("NextState(%s, %s) = %s, want %s", rule.From, rule.Event, to, rule.To)
		}
		if !CanTransition(rule.From, rule.Event) {
			t.Fatalf("CanTransition(%s, %s) = false, want true", rule.From, rule.Event)
		}
		if !rule.IsValid() {
			t.Fatalf("rule %s IsValid = false, want true", rule)
		}
	}
}

func TestTransitionTableDisallowedPairs(t *testing.T) {
	t.Parallel()

	for state := State(0); int(state) < stateCount; state++ {
		for event := Event(0); int(event) < eventCount; event++ {
			allowed := false
			for _, rule := range allowedRules {
				if rule.Matches(state, event) {
					allowed = true
					break
				}
			}
			if allowed {
				continue
			}

			to, ok := NextState(state, event)
			if ok {
				t.Fatalf("NextState(%s, %s) ok = true, want false", state, event)
			}
			if to != state {
				t.Fatalf("NextState(%s, %s) fallback = %s, want %s", state, event, to, state)
			}
			if CanTransition(state, event) {
				t.Fatalf("CanTransition(%s, %s) = true, want false", state, event)
			}
		}
	}
}

func TestTerminalStatesHaveNoOutgoingTransitions(t *testing.T) {
	t.Parallel()

	for _, state := range []State{StateStopped, StateFailed} {
		if got := AllowedTransitions(state); len(got) != 0 {
			t.Fatalf("AllowedTransitions(%s) = %v, want empty", state, got)
		}
	}
}

func TestTransitionRulesReturnsExpectedRules(t *testing.T) {
	t.Parallel()

	got := TransitionRules()

	if len(got) != transitionRuleCount() {
		t.Fatalf("TransitionRules len = %d, want %d", len(got), transitionRuleCount())
	}
	if !reflect.DeepEqual(got, allowedRules) {
		t.Fatalf("TransitionRules = %v, want %v", got, allowedRules)
	}

	got[0] = TransitionRule{}
	if reflect.DeepEqual(TransitionRules(), got) {
		t.Fatal("mutating returned TransitionRules slice changed internal table")
	}
}

func TestAllowedTransitionsPerState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state State
		want  []TransitionRule
	}{
		{StateNew, allowedRules[0:2]},
		{StateStarting, allowedRules[2:5]},
		{StateRunning, allowedRules[5:7]},
		{StateStopping, allowedRules[7:9]},
		{StateStopped, nil},
		{StateFailed, nil},
		{State(99), nil},
	}

	for _, tt := range tests {
		got := AllowedTransitions(tt.state)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("AllowedTransitions(%s) = %v, want %v", tt.state, got, tt.want)
		}
		if len(got) > 0 {
			got[0] = TransitionRule{}
			if reflect.DeepEqual(AllowedTransitions(tt.state), got) {
				t.Fatalf("mutating AllowedTransitions(%s) result changed table", tt.state)
			}
		}
	}
}

func TestNextStateInvalidInputDoesNotPanic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from  State
		event Event
	}{
		{State(99), EventBeginStart},
		{StateNew, Event(99)},
	}

	for _, tt := range tests {
		if _, ok := NextState(tt.from, tt.event); ok {
			t.Fatalf("NextState(%s, %s) ok = true, want false", tt.from, tt.event)
		}
	}
}

func TestCanReachState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from   State
		target State
		want   bool
	}{
		{StateNew, StateNew, true},
		{StateNew, StateFailed, true},
		{StateStarting, StateStopped, true},
		{StateRunning, StateStarting, false},
		{StateRunning, StateFailed, true},
		{StateStopping, StateRunning, false},
		{StateStopped, StateStopped, true},
		{StateStopped, StateFailed, false},
		{StateFailed, StateFailed, true},
		{StateFailed, StateStopped, false},
		{State(99), StateNew, false},
		{StateNew, State(99), false},
	}

	for _, tt := range tests {
		if got := canReachState(tt.from, tt.target); got != tt.want {
			t.Fatalf("canReachState(%s, %s) = %v, want %v", tt.from, tt.target, got, tt.want)
		}
	}
}
