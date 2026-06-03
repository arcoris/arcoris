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

package fieldownership

// EmptyState returns a canonical empty ownership state.
func EmptyState() State {
	return State{}
}

// NewState constructs normalized deterministic ownership state.
//
// Normalization sorts entries by owner, merges duplicate owner entries by field
// union, and prunes owners whose field sets are empty. It does not compact
// ancestor and descendant paths, and it does not enforce unique owner per path.
func NewState(entries ...Entry) (State, error) {
	return normalizeEntries(entries)
}

// MustState constructs ownership state or panics when any entry is invalid.
func MustState(entries ...Entry) State {
	state, err := NewState(entries...)
	if err != nil {
		panic(err)
	}

	return state
}
