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

// canReachState reports whether target is reachable from from through the static
// lifecycle transition graph.
//
// The function ignores guards, context cancellation, observers, and controller
// concurrency. It answers only a graph question:
//
//	Can any sequence of allowed lifecycle transitions move from this state to
//	the target state?
//
// canReachState is used by WaitState to fail early when a requested state is
// already impossible to observe. For example, StateRunning cannot reach
// StateStarting, and StateFailed cannot reach StateStopped.
func canReachState(from State, target State) bool {
	if !from.IsValid() || !target.IsValid() {
		return false
	}

	if from == target {
		return true
	}

	visited := make(map[State]bool, 6)
	queue := []State{from}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, rule := range AllowedTransitions(current) {
			if rule.To == target {
				return true
			}

			if visited[rule.To] {
				continue
			}

			visited[rule.To] = true
			queue = append(queue, rule.To)
		}
	}

	return false
}
