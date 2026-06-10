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

// Selector identifies one associative-list element by stable key fields.
//
// It is not a query language, wildcard, or arbitrary predicate. Higher layers
// decide when semantic selector addressing is appropriate instead of positional
// index addressing.
type Selector struct {
	entries []SelectorEntry
}

// Len returns the number of canonical selector entries stored in s.
//
// Canonical selectors are sorted by field name and contain no duplicate
// identity fields.
func (s Selector) Len() int {
	return len(s.entries)
}

// IsZero reports whether s is the uninitialized zero selector.
//
// The zero selector is not valid as a semantic list-identity selector. Callers
// typically obtain valid selectors through NewSelector or MustSelector.
func (s Selector) IsZero() bool {
	return len(s.entries) == 0
}

// Entries returns a detached canonical slice of selector entries.
func (s Selector) Entries() []SelectorEntry {
	return cloneEntries(s.entries)
}

// Entry returns one canonical selector entry by index.
func (s Selector) Entry(index int) (SelectorEntry, bool) {
	if index < 0 || index >= len(s.entries) {
		return SelectorEntry{}, false
	}

	return s.entries[index], true
}

// Fields returns detached selector field names in canonical entry order.
func (s Selector) Fields() []FieldName {
	if s.entries == nil {
		return nil
	}

	fields := make([]FieldName, len(s.entries))
	for i, entry := range s.entries {
		fields[i] = entry.field
	}

	return fields
}

// ForEach visits canonical selector entries in order until fn returns false.
func (s Selector) ForEach(fn func(index int, entry SelectorEntry) bool) {
	for i, entry := range s.entries {
		if !fn(i, entry) {
			return
		}
	}
}

// Get returns the selector literal associated with field.
//
// Selectors model exact associative-list identity, so field lookup is a simple
// exact-name match and does not interpret wildcards or query expressions.
func (s Selector) Get(field FieldName) (Literal, bool) {
	for _, entry := range s.entries {
		if entry.field == field {
			return entry.value, true
		}
	}

	return Literal{}, false
}

// Has reports whether s contains field.
func (s Selector) Has(field FieldName) bool {
	_, ok := s.Get(field)
	return ok
}
