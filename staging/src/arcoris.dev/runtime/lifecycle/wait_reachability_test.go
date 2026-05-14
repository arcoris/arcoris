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

func TestCanReachStateMatrix(t *testing.T) {
	t.Parallel()

	// Reachability is the static transition graph closure. Guards and runtime
	// progress are intentionally ignored by this table.
	want := map[State]map[State]bool{
		StateNew: {
			StateNew: true, StateStarting: true, StateRunning: true, StateStopping: true, StateStopped: true, StateFailed: true,
		},
		StateStarting: {
			StateStarting: true, StateRunning: true, StateStopping: true, StateStopped: true, StateFailed: true,
		},
		StateRunning: {
			StateRunning: true, StateStopping: true, StateStopped: true, StateFailed: true,
		},
		StateStopping: {
			StateStopping: true, StateStopped: true, StateFailed: true,
		},
		StateStopped: {
			StateStopped: true,
		},
		StateFailed: {
			StateFailed: true,
		},
	}

	for _, from := range allStates {
		for _, target := range allStates {
			if got := canReachState(from, target); got != want[from][target] {
				t.Fatalf("canReachState(%s, %s) = %v, want %v", from, target, got, want[from][target])
			}
		}
	}
}

func TestCanReachStateRejectsInvalidStates(t *testing.T) {
	t.Parallel()

	if canReachState(State(99), StateNew) {
		t.Fatal("canReachState invalid source = true, want false")
	}
	if canReachState(StateNew, State(99)) {
		t.Fatal("canReachState invalid target = true, want false")
	}
}

func TestReachabilityTableMatchesTransitionRulesClosure(t *testing.T) {
	t.Parallel()

	for _, from := range allStates {
		reachable := reachableFromTransitionRules(from)
		for _, target := range allStates {
			if got, want := reachabilityTable[int(from)][int(target)], reachable[target]; got != want {
				t.Fatalf("reachabilityTable[%s][%s] = %v, want %v", from, target, got, want)
			}
		}
	}
}

func reachableFromTransitionRules(from State) map[State]bool {
	visited := map[State]bool{from: true}
	queue := []State{from}
	adjacency := make(map[State][]State)

	for _, rule := range TransitionRules() {
		adjacency[rule.From] = append(adjacency[rule.From], rule.To)
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, target := range adjacency[cur] {
			if visited[target] {
				continue
			}
			visited[target] = true
			queue = append(queue, target)
		}
	}

	return visited
}
