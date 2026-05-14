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

// reachabilityTable describes static graph reachability between states.
//
// It ignores guards, observer behavior, context cancellation, and controller
// concurrency. The table answers only whether any sequence of allowed lifecycle
// transitions can move from one state to another.
var reachabilityTable = [stateCount][stateCount]bool{
	StateNew: {
		StateNew:      true,
		StateStarting: true,
		StateRunning:  true,
		StateStopping: true,
		StateStopped:  true,
		StateFailed:   true,
	},
	StateStarting: {
		StateStarting: true,
		StateRunning:  true,
		StateStopping: true,
		StateStopped:  true,
		StateFailed:   true,
	},
	StateRunning: {
		StateRunning:  true,
		StateStopping: true,
		StateStopped:  true,
		StateFailed:   true,
	},
	StateStopping: {
		StateStopping: true,
		StateStopped:  true,
		StateFailed:   true,
	},
	StateStopped: {
		StateStopped: true,
	},
	StateFailed: {
		StateFailed: true,
	},
}

// canReachState reports whether target is statically reachable from from.
func canReachState(from State, target State) bool {
	if !from.IsValid() || !target.IsValid() {
		return false
	}

	return reachabilityTable[int(from)][int(target)]
}
