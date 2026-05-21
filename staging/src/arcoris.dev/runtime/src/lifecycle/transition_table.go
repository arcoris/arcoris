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

const (
	// stateCount is the number of states in the lifecycle model.
	//
	// State and Event values must remain compact iota-based enums because the
	// transition and reachability tables index directly by their numeric values.
	stateCount = int(StateFailed) + 1

	// eventCount is the number of events in the lifecycle model.
	//
	// State and Event values must remain compact iota-based enums because the
	// transition table indexes directly by their numeric values.
	eventCount = int(EventMarkFailed) + 1
)

// transitionTarget is one cell in the lifecycle transition table.
//
// ok is separate from to because StateNew is a valid zero-value state.
type transitionTarget struct {
	to State
	ok bool
}

// transitionTable is the authoritative lifecycle transition table.
//
// The table is indexed by compact State and Event enum values. Adding a new
// state or event requires preserving the contiguous iota layout and auditing this
// table, reachabilityTable, and tests together.
var transitionTable = [stateCount][eventCount]transitionTarget{
	StateNew: {
		EventBeginStart: {to: StateStarting, ok: true},
		EventBeginStop:  {to: StateStopped, ok: true},
	},

	StateStarting: {
		EventMarkRunning: {to: StateRunning, ok: true},
		EventBeginStop:   {to: StateStopping, ok: true},
		EventMarkFailed:  {to: StateFailed, ok: true},
	},

	StateRunning: {
		EventBeginStop:  {to: StateStopping, ok: true},
		EventMarkFailed: {to: StateFailed, ok: true},
	},

	StateStopping: {
		EventMarkStopped: {to: StateStopped, ok: true},
		EventMarkFailed:  {to: StateFailed, ok: true},
	},
}

// transitionRuleCount returns the number of allowed transitions in the table.
func transitionRuleCount() int {
	count := 0

	for state := State(0); int(state) < stateCount; state++ {
		for event := Event(0); int(event) < eventCount; event++ {
			if transitionTable[int(state)][int(event)].ok {
				count++
			}
		}
	}

	return count
}
