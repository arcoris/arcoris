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

package objectapply

import "arcoris.dev/apimachinery/api/fieldownership"

// State stores object-level field ownership.
//
// v1 owns only the Desired surface. The private shape leaves room for observed
// or metadata ownership without exposing public struct fields.
type State struct {
	// desired is the fieldownership.State for the object's Desired surface.
	//
	// It is private so future observed/metadata ownership can be added without
	// making callers construct State literals.
	desired fieldownership.State
}

// EmptyState returns an empty object ownership state.
//
// It is equivalent to the zero State but makes call sites explicit.
func EmptyState() State {
	return State{}
}

// NewState constructs object ownership state from Desired ownership.
//
// The supplied fieldownership.State is already immutable-by-convention, so it is
// stored by value.
func NewState(desired fieldownership.State) State {
	return State{desired: desired}
}

// IsEmpty reports whether no object surface has ownership state.
//
// Today this means only Desired ownership is empty. The method is intentionally
// future-proof for additional surfaces.
func (s State) IsEmpty() bool {
	return s.desired.IsEmpty()
}

// Desired returns Desired-surface ownership.
//
// The returned fieldownership.State follows fieldownership's
// immutable-by-convention model.
func (s State) Desired() fieldownership.State {
	return s.desired
}

// WithDesired returns a copy of s with replacement Desired ownership.
//
// Other future object surfaces should be preserved by this method.
func (s State) WithDesired(desired fieldownership.State) State {
	s.desired = desired

	return s
}
