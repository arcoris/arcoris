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

func TestNextStateAllowedRules(t *testing.T) {
	t.Parallel()

	// NextState is a pure static lookup: it reports the table target without
	// running guards, assigning commit metadata, or mutating controller state.
	for _, rule := range expectedTransitionRules {
		to, ok := NextState(rule.From, rule.Event)
		if !ok {
			t.Fatalf("NextState(%s, %s) ok = false, want true", rule.From, rule.Event)
		}
		if to != rule.To {
			t.Fatalf("NextState(%s, %s) = %s, want %s", rule.From, rule.Event, to, rule.To)
		}
	}
}

func TestNextStateDisallowedValidPairs(t *testing.T) {
	t.Parallel()

	for _, state := range allStates {
		for _, event := range allEvents {
			if expectedRuleExists(state, event) {
				continue
			}
			to, ok := NextState(state, event)
			if ok {
				t.Fatalf("NextState(%s, %s) ok = true, want false", state, event)
			}
			if to != state {
				t.Fatalf("NextState(%s, %s) fallback = %s, want %s", state, event, to, state)
			}
		}
	}
}

func TestNextStateRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		from  State
		event Event
	}{
		{"invalid source", State(99), EventBeginStart},
		{"invalid event", StateNew, Event(99)},
	}

	for _, tc := range tests {
		to, ok := NextState(tc.from, tc.event)
		if ok {
			t.Fatalf("%s ok = true, want false", tc.name)
		}
		if to != tc.from {
			t.Fatalf("%s fallback = %s, want %s", tc.name, to, tc.from)
		}
	}
}

func TestCanTransitionMirrorsNextState(t *testing.T) {
	t.Parallel()

	for _, state := range append(append([]State(nil), allStates...), State(99)) {
		for _, event := range append(append([]Event(nil), allEvents...), Event(99)) {
			_, ok := NextState(state, event)
			if got := CanTransition(state, event); got != ok {
				t.Fatalf("CanTransition(%s, %s) = %v, want %v", state, event, got, ok)
			}
		}
	}
}

func TestAllowedTransitionsPerState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state State
		want  []TransitionRule
	}{
		{StateNew, expectedTransitionRules[0:2]},
		{StateStarting, expectedTransitionRules[2:5]},
		{StateRunning, expectedTransitionRules[5:7]},
		{StateStopping, expectedTransitionRules[7:9]},
		{StateStopped, nil},
		{StateFailed, nil},
		{State(99), nil},
	}

	for _, tc := range tests {
		assertDeepEqual(t, AllowedTransitions(tc.state), tc.want)
	}
}

func TestAllowedTransitionsReturnsCallerOwnedSlice(t *testing.T) {
	t.Parallel()

	// Returned slices are caller-owned so tests and integrations can mutate or
	// append to them without corrupting the static transition table.
	got := AllowedTransitions(StateNew)
	got[0] = TransitionRule{}
	assertDeepEqual(t, AllowedTransitions(StateNew), expectedTransitionRules[0:2])
}

func TestAppendAllowedTransitionsAppendsToExistingSlice(t *testing.T) {
	t.Parallel()

	existing := []TransitionRule{{From: StateFailed, Event: EventMarkFailed, To: StateFailed}}
	got := AppendAllowedTransitions(existing, StateNew)
	want := append([]TransitionRule(nil), existing...)
	want = append(want, expectedTransitionRules[0:2]...)
	assertDeepEqual(t, got, want)
}

func TestAppendAllowedTransitionsIgnoresInvalidState(t *testing.T) {
	t.Parallel()

	existing := []TransitionRule{{From: StateNew, Event: EventBeginStart, To: StateStarting}}
	got := AppendAllowedTransitions(existing, State(99))
	assertDeepEqual(t, got, existing)
}

func TestTransitionRulesReturnsCompleteTableInOrder(t *testing.T) {
	t.Parallel()

	assertDeepEqual(t, TransitionRules(), expectedTransitionRules)
}

func TestTransitionRulesReturnsCallerOwnedSlice(t *testing.T) {
	t.Parallel()

	got := TransitionRules()
	got[0] = TransitionRule{}
	assertDeepEqual(t, TransitionRules(), expectedTransitionRules)
}
