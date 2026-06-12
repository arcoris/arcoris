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

// EmptyState returns an empty object ownership state.
//
// The zero State is already empty and valid. This helper exists for call sites
// that benefit from saying explicitly that they are constructing object
// ownership rather than ordinary field ownership.
func EmptyState() State {
	return State{}
}

// NewState constructs object ownership state from Desired ownership.
//
// This is the narrow constructor for Desired-only callers such as objectapply.
// Observed and metadata ownership are left empty. The supplied
// fieldownership.State is stored by value and follows fieldownership's
// immutable-by-convention contract.
func NewState(desired fieldownership.State) State {
	return State{desired: desired}
}

// NewStateWithSurfaces constructs ownership state for every modeled surface.
//
// The function performs no cross-surface merging. Each supplied
// fieldownership.State is already scoped to its own surface root and remains
// independent after construction.
func NewStateWithSurfaces(
	desired fieldownership.State,
	observed fieldownership.State,
	metadata MetadataState,
) State {
	return State{
		desired:  desired,
		observed: observed,
		metadata: metadata,
	}
}
