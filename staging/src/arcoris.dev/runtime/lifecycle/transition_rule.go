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

// TransitionRule describes one allowed lifecycle state transition.
//
// A TransitionRule is static transition-table data. It has no revision, no
// timestamp, and no failure cause because it describes what is allowed by the
// lifecycle model, not what happened at runtime.
type TransitionRule struct {
	// From is the source state accepted by this rule.
	From State

	// Event is the lifecycle input accepted by this rule.
	Event Event

	// To is the target state produced by this rule.
	To State
}

// String returns a compact diagnostic representation of r.
func (r TransitionRule) String() string {
	return r.From.String() + " --" + r.Event.String() + "--> " + r.To.String()
}

// IsValid reports whether r is a structurally valid transition-table rule.
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
