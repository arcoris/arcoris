// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capacity

// State is pure multi-resource accounting state.
//
// State contains source facts only: configured limits and currently reserved
// amounts. It owns no locks, reservations, waiters, metrics, or policy.
type State struct {
	// Limits is the configured capacity by resource.
	Limits Vector

	// Reserved is the live reserved capacity by resource.
	Reserved Vector
}

// NewState returns canonical source accounting state.
func NewState(limits Vector, reserved Vector) (State, error) {
	if !limits.IsValid() {
		return State{}, errorAt(
			"limits",
			ErrInvalidState,
			ErrorReasonInvalidState,
			"state limits vector must be canonical",
		)
	}
	if !reserved.IsValid() {
		return State{}, errorAt(
			"reserved",
			ErrInvalidState,
			ErrorReasonInvalidState,
			"state reserved vector must be canonical",
		)
	}

	return State{Limits: limits, Reserved: reserved}, nil
}

// MustState returns NewState(limits, reserved) or panics.
func MustState(limits Vector, reserved Vector) State {
	state, err := NewState(limits, reserved)
	if err != nil {
		panic(err)
	}

	return state
}

// IsValid reports whether s contains canonical vectors.
func (s State) IsValid() bool {
	return s.Limits.IsValid() && s.Reserved.IsValid()
}

// Snapshot derives the immutable read model for s.
func (s State) Snapshot() Snapshot {
	return NewSnapshot(s.Limits, s.Reserved)
}

// Check evaluates demand against s without mutating state.
func (s State) Check(demand Demand) CheckResult {
	return s.Snapshot().Check(demand)
}

// Reserve returns state with demand reserved when it fully fits.
//
// Reserve is all-or-nothing. On refusal, the returned State is s and no partial
// resource reservation is recorded.
func (s State) Reserve(demand Demand) (State, CheckResult) {
	result := s.Check(demand)
	if result.Denied() {
		return s, result
	}

	reserved, ok := s.Reserved.CheckedAdd(demand.Vector())
	if !ok {
		return s, CheckResult{
			Status:  ReserveStatusInsufficient,
			Missing: demand.Vector(),
		}
	}

	return State{Limits: s.Limits, Reserved: reserved}, result
}

// Release returns state with demand subtracted from Reserved.
//
// Release returns false when s does not currently hold the full demand.
func (s State) Release(demand Demand) (State, bool) {
	if !demand.IsValid() {
		panicAt(
			"demand",
			ErrInvalidDemand,
			ErrorReasonInvalidDemand,
			"demand must be non-empty and canonical",
		)
	}

	reserved, ok := s.Reserved.CheckedSub(demand.Vector())
	if !ok {
		return s, false
	}

	return State{Limits: s.Limits, Reserved: reserved}, true
}
