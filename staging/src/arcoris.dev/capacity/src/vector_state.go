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

// VectorState is pure multi-resource accounting state.
//
// VectorState contains source facts only: configured limits and currently
// reserved amounts. It owns no locks, reservations, waiters, metrics, admission
// semantics, or scheduling policy.
type VectorState struct {
	// Limits is the configured capacity by resource.
	Limits Vector

	// Reserved is the live reserved capacity by resource.
	Reserved Vector
}

// NewVectorState returns canonical source accounting state.
func NewVectorState(limits Vector, reserved Vector) (VectorState, error) {
	if !limits.IsValid() {
		return VectorState{}, errorAt(
			"limits",
			ErrInvalidVectorState,
			"limits vector must be canonical",
		)
	}
	if !reserved.IsValid() {
		return VectorState{}, errorAt(
			"reserved",
			ErrInvalidVectorState,
			"reserved vector must be canonical",
		)
	}

	return VectorState{Limits: limits, Reserved: reserved}, nil
}

// MustVectorState returns NewVectorState(limits, reserved) or panics.
func MustVectorState(limits Vector, reserved Vector) VectorState {
	state, err := NewVectorState(limits, reserved)
	if err != nil {
		panic(err)
	}

	return state
}

// IsValid reports whether s contains canonical vectors.
func (s VectorState) IsValid() bool {
	return s.Limits.IsValid() && s.Reserved.IsValid()
}

// Snapshot derives the immutable vector read model for s.
func (s VectorState) Snapshot() VectorSnapshot {
	return NewVectorSnapshot(s.Limits, s.Reserved)
}

// Fit evaluates demand against s without mutating state.
func (s VectorState) Fit(demand Demand) Fit {
	return s.Snapshot().Fit(demand)
}

// WithReserved returns state with demand added to Reserved when it fully fits.
//
// The transformation is all-or-nothing. On refusal, the returned VectorState is
// s and no partial resource reservation is recorded.
func (s VectorState) WithReserved(demand Demand) (VectorState, Fit) {
	fit := s.Fit(demand)
	if fit.Refused() {
		return s, fit
	}

	reserved, ok := s.Reserved.CheckedAdd(demand.Vector())
	if !ok {
		return s, Fit{
			Refusal: RefusalInsufficient,
			Missing: demand.Vector(),
		}
	}

	return VectorState{Limits: s.Limits, Reserved: reserved}, fit
}

// WithoutReserved returns state with demand subtracted from Reserved.
//
// The method returns false and leaves state unchanged when s does not currently
// hold the full demand.
func (s VectorState) WithoutReserved(demand Demand) (VectorState, bool) {
	requireValidDemand("demand", demand)

	reserved, ok := s.Reserved.CheckedSub(demand.Vector())
	if !ok {
		return s, false
	}

	return VectorState{Limits: s.Limits, Reserved: reserved}, true
}
