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

// Normalize canonicalizes every modeled surface without changing ownership
// semantics.
//
// Each surface is normalized independently. Desired ownership is never merged
// with Observed ownership, and metadata.labels ownership is never merged with
// metadata.annotations ownership.
func Normalize(state State) State {
	return NewStateWithSurfaces(
		normalizeFieldOwnership(state.Desired()),
		normalizeFieldOwnership(state.Observed()),
		NewMetadataState(
			normalizeFieldOwnership(state.Metadata().Labels()),
			normalizeFieldOwnership(state.Metadata().Annotations()),
		),
	)
}

// normalizeFieldOwnership re-runs fieldownership normalization for one surface.
//
// Public objectownership.State values should already contain structurally valid
// fieldownership.State values. The panic protects against an internal invariant
// break rather than user input.
func normalizeFieldOwnership(state fieldownership.State) fieldownership.State {
	normalized, err := fieldownership.NewState(state.Entries()...)
	if err != nil {
		panic(err)
	}

	return normalized
}
