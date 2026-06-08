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

// PreferenceSet is an immutable deterministic set of encode preferences.
//
// Preferences are ordered by descending weight. Original caller order breaks
// ties, which keeps selection deterministic without interpreting protocol
// q-values directly.
type PreferenceSet struct {
	// items stores normalized preferences in deterministic selection order.
	items []Preference
}

// NewPreferenceSet validates and orders encode preferences.
func NewPreferenceSet(items ...Preference) (PreferenceSet, error) {
	return normalizePreferenceSetAt(pathPreferences, items)
}

// MustPreferenceSet returns a normalized PreferenceSet or panics on invalid input.
func MustPreferenceSet(items ...Preference) PreferenceSet {
	preferences, err := NewPreferenceSet(items...)
	if err != nil {
		panic(err)
	}

	return preferences
}

// IsZero reports whether p contains no preferences.
func (p PreferenceSet) IsZero() bool {
	return len(p.items) == 0
}

// Len returns the preference count.
func (p PreferenceSet) Len() int {
	return len(p.items)
}

// Preferences returns detached normalized preferences in selection order.
func (p PreferenceSet) Preferences() []Preference {
	return slices.Clone(p.items)
}
