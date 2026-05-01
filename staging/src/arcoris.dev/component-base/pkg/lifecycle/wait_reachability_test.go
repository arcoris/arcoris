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

func TestReachabilityTableMatchesTransitionGraph(t *testing.T) {
	t.Parallel()

	states := []State{
		StateNew,
		StateStarting,
		StateRunning,
		StateStopping,
		StateStopped,
		StateFailed,
	}

	expected := make(map[State]map[State]bool, len(states))
	for _, from := range states {
		expected[from] = reachableFromTransitionRules(from)
	}

	for _, from := range states {
		for _, target := range states {
			got := canReachState(from, target)
			want := expected[from][target]
			if got != want {
				t.Fatalf("canReachState(%s, %s) = %v, want %v", from, target, got, want)
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
		current := queue[0]
		queue = queue[1:]

		for _, target := range adjacency[current] {
			if visited[target] {
				continue
			}
			visited[target] = true
			queue = append(queue, target)
		}
	}

	return visited
}
