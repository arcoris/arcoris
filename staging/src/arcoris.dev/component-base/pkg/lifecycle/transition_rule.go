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
	// stateCount is the number of states in the lifecycle model. It is used to dimension the transition table.
	//
	// stateCount is computed from the State enum instead of hardcoded so additions
	// to the lifecycle graph do not require updating a separate constant.
	stateCount = int(StateFailed) + 1

	// eventCount is the number of events in the lifecycle model. It is used to dimension the transition table.
	//
	// eventCount is computed from the Event enum instead of hardcoded so additions
	// to the lifecycle graph do not require updating a separate constant.
	eventCount = int(EventMarkFailed) + 1
)

// TransitionRule describes one allowed lifecycle state transition.
//
// A TransitionRule is static transition-table data. It has no revision, no
// timestamp, and no failure cause because it describes what is allowed by the
// lifecycle model, not what actually happened at runtime.
//
// TransitionRule exists so the transition graph can be inspected, tested, and
// rendered without exposing the internal lookup table.
type TransitionRule struct {
	// From is the source state accepted by this rule.
	From State

	// Event is the lifecycle input accepted by this rule.
	Event Event

	// To is the target state produced by this rule.
	To State
}

// transitionTarget is one cell in the lifecycle transition table.
//
// ok is stored separately from to because StateNew is a valid zero-value state.
// Without an explicit ok flag, a missing transition could be confused with a
// transition to StateNew.
type transitionTarget struct {
	to State
	ok bool
}

// transitionTable is the authoritative lifecycle transition table.
//
// The table is indexed by State and Event. This makes NextState a constant-time
// lookup while keeping the lifecycle graph explicit and easy to audit.
//
// The table relies on State and Event being compact iota-based uint8 enums. That
// invariant is local to this package and MUST be preserved when adding new states
// or events.
//
// The default lifecycle graph is:
//
//	New       --EventBeginStart-->   Starting
//	New       --EventBeginStop-->    Stopped
//
//	Starting  --EventMarkRunning-->  Running
//	Starting  --EventBeginStop-->    Stopping
//	Starting  --EventMarkFailed-->   Failed
//
//	Running   --EventBeginStop-->    Stopping
//	Running   --EventMarkFailed-->   Failed
//
//	Stopping  --EventMarkStopped-->  Stopped
//	Stopping  --EventMarkFailed-->   Failed
//
//	Stopped   --terminal-->         no outgoing transitions
//	Failed    --terminal-->         no outgoing transitions
//
// The table intentionally does not model restart, pause, reload, drain,
// supervisor retry, or health degradation. Those are higher-level runtime or
// domain concerns, not base lifecycle transitions.
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

// String returns a compact diagnostic representation of r.
//
// The returned value is intended for diagnostics and tests. It is not a stable
// serialization format.
func (r TransitionRule) String() string {
	return r.From.String() + " --" + r.Event.String() + "--> " + r.To.String()
}

// IsValid reports whether r is a structurally valid transition-table rule.
//
// IsValid validates only the static rule shape. It does not check failure causes
// because a rule describes the lifecycle graph, not a runtime transition.
func (r TransitionRule) IsValid() bool {
	if !r.From.IsValid() || !r.To.IsValid() || !r.Event.IsValid() {
		return false
	}

	to, ok := NextState(r.From, r.Event)
	return ok && to == r.To
}

// Matches reports whether r accepts the given state/event pair.
func (r TransitionRule) Matches(from State, event Event) bool {
	return r.From == from && r.Event == event
}

// IsTerminal reports whether r moves the lifecycle into a terminal state.
func (r TransitionRule) IsTerminal() bool {
	return r.To.IsTerminal()
}

// IsFailure reports whether r is the static failure rule for its source state.
func (r TransitionRule) IsFailure() bool {
	return r.Event == EventMarkFailed && r.To == StateFailed
}

// NextState returns the state produced by applying event to from.
//
// The returned ok value is false when the pair is not allowed by the lifecycle
// transition table. When ok is false, the returned State is the original from
// value and MUST be ignored by callers that need a valid target state.
//
// NextState is a pure table lookup. It does not mutate state, run guards, check
// failure causes, assign transition revisions, assign timestamps, notify
// observers, or enforce concurrency. Controller remains the authoritative owner
// of committed lifecycle state.
//
// NextState is not a synchronization primitive. A true result only means the
// pair is allowed by the static transition table. Concurrent callers must still
// use Controller to serialize and commit real lifecycle transitions.
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

// CanTransition reports whether event can be applied to from according to the
// lifecycle transition table.
//
// CanTransition is a convenience wrapper around NextState. It is useful for
// tests, diagnostics, guards, and callers that need to check table-level
// possibility before attempting a controller-owned transition.
//
// CanTransition is not a synchronization primitive.
func CanTransition(from State, event Event) bool {
	_, ok := NextState(from, event)
	return ok
}

// AllowedTransitions returns all transition-table rules that can be applied from
// from.
//
// The returned slice is newly allocated and may be modified by the caller.
// Mutating it does not affect the package-level transition table.
//
// For an invalid state or a terminal state with no outgoing transitions,
// AllowedTransitions returns nil.
func AllowedTransitions(from State) []TransitionRule {
	return AppendAllowedTransitions(nil, from)
}

// AppendAllowedTransitions appends all transition-table rules that can be
// applied from from to dst and returns the extended slice.
//
// AppendAllowedTransitions is the allocation-conscious form of
// AllowedTransitions. It is useful in tests, diagnostics, or controller code that
// wants to reuse a caller-owned buffer.
//
// For an invalid state or a terminal state with no outgoing transitions, dst is
// returned unchanged.
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
// Mutating it does not affect the package-level transition table.
func TransitionRules() []TransitionRule {
	rules := make([]TransitionRule, 0, transitionRuleCount())

	for state := State(0); int(state) < stateCount; state++ {
		rules = AppendAllowedTransitions(rules, state)
	}

	return rules
}

// reduceTransition builds a candidate transition from a state, an event, and an
// optional failure cause.
//
// This is the package-local pure transition reducer. It maps the current
// lifecycle state and an input event to a candidate Transition without consulting
// external state, running guards, assigning commit metadata, notifying observers,
// or performing side effects.
//
// The returned transition is not committed. Controller is responsible for:
//
//   - rejecting missing failure causes with the proper error;
//   - running transition guards before commit;
//   - assigning Revision and At;
//   - publishing the committed state;
//   - notifying waiters and observers.
//
// When ok is false, the returned Transition is diagnostic only and MUST NOT be
// committed.
func reduceTransition(from State, event Event, cause error) (transition Transition, ok bool) {
	to, ok := NextState(from, event)
	if !ok {
		return Transition{
			From:  from,
			To:    from,
			Event: event,
			Cause: cause,
		}, false
	}

	return Transition{
		From:  from,
		To:    to,
		Event: event,
		Cause: cause,
	}, true
}

// transitionRuleCount returns the number of allowed transitions in the static
// transition table.
//
// The count is computed from the table instead of hardcoded so additions to the
// transition graph do not require updating a separate constant.
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
