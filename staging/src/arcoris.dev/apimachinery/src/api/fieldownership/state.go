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

// State is deterministic immutable-by-convention field ownership state.
//
// Entries are sorted by Owner. Multiple owners may own the same path; conflict
// detection, not state normalization, decides when another attempted field set
// overlaps existing ownership.
type State struct {
	entries []Entry
}

// IsEmpty reports whether s contains no non-empty owner entries.
func (s State) IsEmpty() bool {
	return len(s.entries) == 0
}
