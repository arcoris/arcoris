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

package objectownership

import "arcoris.dev/apimachinery/api/fieldownership"

// State stores object-level field ownership.
//
// v1 owns only the Desired surface. The private shape leaves room for observed
// or metadata ownership without exposing public struct fields.
type State struct {
	// desired is the fieldownership.State for the object's Desired surface.
	//
	// It stays private so callers cannot construct brittle State literals that
	// would block future Observed or metadata ownership surfaces.
	desired fieldownership.State
}

// EmptyState returns an empty object ownership state.
//
// It is equivalent to the zero State but keeps call sites explicit when they are
// constructing object ownership rather than ordinary field ownership.
func EmptyState() State {
	return State{}
}

// NewState constructs object ownership state from Desired ownership.
//
// The supplied fieldownership.State is stored by value and follows
// fieldownership's immutable-by-convention contract.
func NewState(desired fieldownership.State) State {
	return State{desired: desired}
}

// IsEmpty reports whether no object surface has ownership state.
//
// Today this means Desired ownership is empty. The method intentionally hides
// that detail so adding future surfaces does not change callers.
func (s State) IsEmpty() bool {
	return s.desired.IsEmpty()
}

// Desired returns Desired-surface ownership.
//
// The returned fieldownership.State remains immutable-by-convention; callers
// transform it through fieldownership APIs and store the result with WithDesired.
func (s State) Desired() fieldownership.State {
	return s.desired
}

// WithDesired returns a copy of s with replacement Desired ownership.
//
// Future object surfaces should be preserved by this method, making it the
// stable update point for Desired ownership.
func (s State) WithDesired(desired fieldownership.State) State {
	s.desired = desired

	return s
}
