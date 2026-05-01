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

// NextState returns the state produced by applying event to from.
//
// NextState is a pure table lookup. It does not mutate state, run guards, assign
// revisions, assign timestamps, notify observers, or provide synchronization.
func NextState(from State, event Event) (State, bool) {
	if !from.IsValid() || !event.IsValid() {
		return from, false
	}

	target := transitionTable[int(from)][int(event)]
	if !target.ok {
		return from, false
	}

	return target.to, true
}

// CanTransition reports whether event can be applied to from.
//
// CanTransition is not a synchronization primitive. It reports only static
// table-level possibility.
func CanTransition(from State, event Event) bool {
	_, ok := NextState(from, event)
	return ok
}

// AllowedTransitions returns all rules that can be applied from from.
//
// The returned slice is newly allocated and may be modified by the caller.
func AllowedTransitions(from State) []TransitionRule {
	return AppendAllowedTransitions(nil, from)
}

// AppendAllowedTransitions appends rules that can be applied from from to dst.
func AppendAllowedTransitions(dst []TransitionRule, from State) []TransitionRule {
	if !from.IsValid() {
		return dst
	}

	for event := Event(0); int(event) < eventCount; event++ {
		target := transitionTable[int(from)][int(event)]
		if !target.ok {
			continue
		}

		dst = append(dst, TransitionRule{
			From:  from,
			Event: event,
			To:    target.to,
		})
	}

	return dst
}

// TransitionRules returns the complete lifecycle transition table as rules.
//
// The returned slice is newly allocated and may be modified by the caller.
func TransitionRules() []TransitionRule {
	rules := make([]TransitionRule, 0, transitionRuleCount())

	for state := State(0); int(state) < stateCount; state++ {
		rules = AppendAllowedTransitions(rules, state)
	}

	return rules
}
