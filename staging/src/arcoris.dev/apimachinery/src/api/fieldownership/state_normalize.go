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

import (
	"cmp"
	"slices"

	"arcoris.dev/apimachinery/api/fieldpath"
)

// normalizeEntries validates, sorts, merges, and prunes ownership entries.
func normalizeEntries(entries []Entry) (State, error) {
	normalized := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if err := entry.Validate(); err != nil {
			return State{}, err
		}
		if entry.IsEmpty() {
			continue
		}

		normalized = append(normalized, entry)
	}

	slices.SortFunc(normalized, compareEntriesByOwner)
	return State{entries: mergeEntriesByOwner(normalized)}, nil
}

// compareEntriesByOwner orders entries for canonical State storage.
func compareEntriesByOwner(left Entry, right Entry) int {
	return cmp.Compare(left.owner.String(), right.owner.String())
}

// mergeEntriesByOwner unions adjacent entries with the same owner.
func mergeEntriesByOwner(entries []Entry) []Entry {
	if len(entries) == 0 {
		return nil
	}

	merged := make([]Entry, 0, len(entries))
	current := entries[0]
	for _, entry := range entries[1:] {
		if entry.owner == current.owner {
			current.fields = current.fields.Union(entry.fields)
			continue
		}

		merged = append(merged, current)
		current = entry
	}

	return append(merged, current)
}

// validateFields defensively validates every path contained in fields.
//
// Public fieldpath.Set constructors already validate paths. This extra check
// keeps State robust if a future internal caller receives a malformed set value.
func validateFields(fields fieldpath.Set, detail string) error {
	for _, p := range fields.Paths() {
		if err := p.ValidateStructure(); err != nil {
			return wrapPathError(
				p,
				detail,
				err,
			)
		}
	}

	return nil
}
