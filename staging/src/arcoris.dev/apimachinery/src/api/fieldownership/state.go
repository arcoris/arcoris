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
// Entries are sorted by Owner. Duplicate owners are merged and empty entries are
// pruned by NewState. Multiple owners may own exact or overlapping paths because
// conflict checks are contextual to an acting owner and attempted field set.
type State struct {
	entries []Entry
}

// Len returns the number of non-empty owner entries in s.
func (s State) Len() int {
	return len(s.entries)
}

// IsEmpty reports whether s contains no non-empty owner entries.
func (s State) IsEmpty() bool {
	return len(s.entries) == 0
}
