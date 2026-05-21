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

func TestTransitionTableEnumBoundaries(t *testing.T) {
	t.Parallel()

	if stateCount != len(allStates) {
		t.Fatalf("stateCount = %d, want %d", stateCount, len(allStates))
	}
	if eventCount != len(allEvents) {
		t.Fatalf("eventCount = %d, want %d", eventCount, len(allEvents))
	}
}

func TestTransitionRuleCount(t *testing.T) {
	t.Parallel()

	if got, want := transitionRuleCount(), len(expectedTransitionRules); got != want {
		t.Fatalf("transitionRuleCount = %d, want %d", got, want)
	}
}

func TestTransitionTableContainsAuthoritativeRules(t *testing.T) {
	t.Parallel()

	// transitionTable is the authoritative lifecycle graph. Lookup helpers may
	// expose it differently, but this test owns the exact table contents.
	for _, rule := range expectedTransitionRules {
		entry := transitionTable[int(rule.From)][int(rule.Event)]
		if !entry.ok {
			t.Fatalf("transitionTable[%s][%s].ok = false, want true", rule.From, rule.Event)
		}
		if entry.to != rule.To {
			t.Fatalf("transitionTable[%s][%s].to = %s, want %s", rule.From, rule.Event, entry.to, rule.To)
		}
	}

	for _, state := range allStates {
		for _, event := range allEvents {
			if expectedRuleExists(state, event) {
				continue
			}
			if transitionTable[int(state)][int(event)].ok {
				t.Fatalf("transitionTable[%s][%s].ok = true, want false", state, event)
			}
		}
	}
}

func TestTransitionTableTerminalStatesHaveNoEntries(t *testing.T) {
	t.Parallel()

	// Terminal states end a controller instance, so the authoritative table must
	// not contain any outgoing edge from stopped or failed.
	for _, state := range []State{StateStopped, StateFailed} {
		for _, event := range allEvents {
			if transitionTable[int(state)][int(event)].ok {
				t.Fatalf("terminal state %s has outgoing event %s", state, event)
			}
		}
	}
}

func expectedRuleExists(from State, event Event) bool {
	for _, rule := range expectedTransitionRules {
		if rule.Matches(from, event) {
			return true
		}
	}

	return false
}
