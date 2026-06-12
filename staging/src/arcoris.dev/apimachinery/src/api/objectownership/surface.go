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

// Surface stores ownership entries for one object surface.
//
// v1 uses only the Desired surface. Entries are intentionally stored as a slice
// so Document can preserve raw input before Normalize sorts, merges, and prunes
// it.
type Surface struct {
	// Entries lists owner/path records for this surface. Raw documents may
	// contain duplicate owners or empty field lists; Normalize canonicalizes
	// those cases.
	Entries []Entry
}

// Clone returns a detached copy of s without validating or normalizing it.
//
// Clone preserves raw entry order, duplicate owners, empty entries, and
// nil-vs-empty Entries shape. Nested Entry.Fields slices are detached.
func (s Surface) Clone() Surface {
	return Surface{Entries: s.EntriesCopy()}
}

// EntriesCopy returns detached entries with detached nested field slices.
//
// EntriesCopy preserves nil-vs-empty slice shape and raw entry order. It does
// not validate, normalize, merge, sort, or prune entries.
func (s Surface) EntriesCopy() []Entry {
	if s.Entries == nil {
		return nil
	}

	out := make([]Entry, len(s.Entries))
	for i, entry := range s.Entries {
		out[i] = entry.Clone()
	}

	return out
}

// IsEmpty reports whether the surface contains no owned fields.
//
// IsEmpty is not a validity check. It ignores owner validity, path validity,
// duplicate entries, and ordering, and only answers whether any entry mentions
// at least one field.
//
// This is semantic emptiness, not raw slice emptiness: entries with no fields
// are ignored because they would be pruned from normalized documents.
func (s Surface) IsEmpty() bool {
	for _, entry := range s.Entries {
		if !entry.IsEmpty() {
			return false
		}
	}

	return true
}
