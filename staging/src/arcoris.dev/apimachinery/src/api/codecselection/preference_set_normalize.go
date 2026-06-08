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

package codecselection

import "slices"

// normalizePreferenceSetAt validates and orders preferences at path.
func normalizePreferenceSetAt(path string, items []Preference) (PreferenceSet, error) {
	if len(items) == 0 {
		return PreferenceSet{}, nil
	}

	normalized := make([]Preference, 0, len(items))
	seen := make(map[string]int, len(items))
	for i, item := range items {
		item.order = i
		preference, err := normalizePreferenceAt(preferencePath(path, i), item)
		if err != nil {
			return PreferenceSet{}, err
		}
		key := preference.contentType.key()
		if previous, ok := seen[key]; ok {
			return PreferenceSet{}, errorfAt(
				preferencePath(path, i),
				ErrInvalidPreference,
				ErrorReasonInvalidPreference,
				"preference content type duplicates preferences[%d]",
				previous,
			)
		}
		seen[key] = i
		normalized = append(normalized, preference)
	}

	slices.SortStableFunc(normalized, comparePreferences)

	return PreferenceSet{items: normalized}, nil
}
