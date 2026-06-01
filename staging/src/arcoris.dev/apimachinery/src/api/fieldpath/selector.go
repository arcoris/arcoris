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
//
// The returned slice belongs to the caller. Mutating it does not affect the
// selector.
func (s Selector) Entries() []SelectorEntry {
	return cloneEntries(s.entries)
}

// Get returns the selector literal associated with field.
//
// Selectors model exact associative-list identity, so field lookup is a simple
// exact-name match and does not interpret wildcards or query expressions.
func (s Selector) Get(field string) (Literal, bool) {
	for _, entry := range s.entries {
		if entry.field == field {
			return entry.value, true
		}
	}

	return Literal{}, false
}

// Has reports whether s contains field.
func (s Selector) Has(field string) bool {
	_, ok := s.Get(field)
	return ok
}

// clone returns a detached selector copy.
func (s Selector) clone() Selector {
	return Selector{entries: cloneEntries(s.entries)}
}

// cloneEntries returns a caller-owned entry slice copy.
func cloneEntries(entries []SelectorEntry) []SelectorEntry {
	if entries == nil {
		return nil
	}

	cloned := make([]SelectorEntry, len(entries))
	copy(cloned, entries)
	return cloned
}
