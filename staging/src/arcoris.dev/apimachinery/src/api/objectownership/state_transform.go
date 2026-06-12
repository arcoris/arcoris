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

// WithDesired returns a copy of s with replacement Desired ownership.
//
// Observed and metadata ownership are preserved. This is the stable update
// point for Desired-only apply orchestration.
func (s State) WithDesired(desired fieldownership.State) State {
	s.desired = desired

	return s
}

// WithObserved returns a copy of s with replacement Observed ownership.
//
// Desired and metadata ownership are preserved so observed/status updates cannot
// accidentally erase user intent or metadata ownership.
func (s State) WithObserved(observed fieldownership.State) State {
	s.observed = observed

	return s
}

// WithMetadata returns a copy of s with replacement metadata ownership.
//
// Desired and Observed ownership are preserved. Metadata patch operations should
// use this method instead of reconstructing State directly.
func (s State) WithMetadata(metadata MetadataState) State {
	s.metadata = metadata

	return s
}
