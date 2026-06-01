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

package fieldpath

// NewSelector constructs a canonical selector from one or more entries.
//
// Input order is accepted for ergonomics, but the stored selector is
// canonicalized into deterministic field order before validation. This makes
// selector equality and formatting independent from construction order.
func NewSelector(entries ...SelectorEntry) (Selector, error) {
	canonical := Selector{
		entries: cloneEntries(entries),
	}

	sortSelectorEntries(canonical.entries)

	if err := canonical.Validate(); err != nil {
		return Selector{}, err
	}

	return canonical, nil
}

// MustSelector constructs a selector or panics when entries are malformed.
//
// It is intended for static fixtures, tests, and package-level declarations
// where malformed selector identity should fail immediately.
func MustSelector(entries ...SelectorEntry) Selector {
	selector, err := NewSelector(entries...)
	if err != nil {
		panic(err)
	}

	return selector
}
