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

// State stores canonical object-level field ownership.
//
// Each fieldownership.State is scoped to one object surface. Paths are relative
// to that surface root, so Desired $.image, Observed $.ready, and metadata label
// $["app"] never share one global path namespace. State is immutable by
// convention: callers update it by creating replacement fieldownership.State
// values and installing them through With* methods.
type State struct {
	// desired stores declarative user/manager intent ownership.
	//
	// objectapply is intentionally allowed to replace only this surface.
	desired fieldownership.State
	// observed stores runtime/controller-reported ownership.
	//
	// Observed ownership is reserved for observed/status update operations. It
	// must not be modified by Desired apply.
	observed fieldownership.State
	// metadata stores the supported ObjectMeta map ownership surfaces.
	//
	// Identity and system metadata fields are not modeled here.
	metadata MetadataState
}

// MetadataState stores ownable ObjectMeta map surfaces.
//
// Only labels and annotations are modeled. Object identity/system fields,
// finalizers, and owner references have separate lifecycle/governance semantics
// and are intentionally excluded from this generic ownership state.
type MetadataState struct {
	// labels stores ownership of individual metadata.labels keys.
	labels fieldownership.State
	// annotations stores ownership of individual metadata.annotations keys.
	annotations fieldownership.State
}

// IsEmpty reports whether no modeled object surface has ownership state.
func (s State) IsEmpty() bool {
	return s.desired.IsEmpty() &&
		s.observed.IsEmpty() &&
		s.metadata.IsEmpty()
}

// Desired returns Desired-surface ownership.
//
// The returned fieldownership.State remains immutable by convention. Callers
// transform it with fieldownership APIs and store the replacement with
// WithDesired.
func (s State) Desired() fieldownership.State {
	return s.desired
}

// Observed returns Observed-surface ownership.
//
// Paths in this state are relative to the Observed payload root. They must not
// be interpreted as Desired paths or as whole-object paths.
func (s State) Observed() fieldownership.State {
	return s.observed
}

// Metadata returns metadata map ownership state.
//
// The returned MetadataState contains only labels and annotations. It does not
// include identity/system metadata fields.
func (s State) Metadata() MetadataState {
	return s.metadata
}
