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

package labels

import "sort"

// Set stores label key/value metadata.
type Set map[Key]Value

// IsZero reports whether the set has no entries.
func (s Set) IsZero() bool {
	return len(s) == 0
}

// Len returns the number of label entries.
func (s Set) Len() int {
	return len(s)
}

// Get returns the value for key.
func (s Set) Get(key Key) (Value, bool) {
	value, ok := s[key]
	return value, ok
}

// Has reports whether key is present.
func (s Set) Has(key Key) bool {
	_, ok := s[key]
	return ok
}

// Keys returns label keys in deterministic lexical order.
func (s Set) Keys() []Key {
	keys := make([]Key, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// FromStrings validates and converts raw string label metadata into a typed set.
func FromStrings(values map[string]string) (Set, error) {
	if values == nil {
		return nil, nil
	}

	set := make(Set, len(values))
	for key, value := range values {
		set[Key(key)] = Value(value)
	}
	if err := set.ValidateLexical(); err != nil {
		return nil, err
	}
	return set, nil
}

// Strings returns a detached raw string map for codec and interop boundaries.
func (s Set) Strings() map[string]string {
	if s == nil {
		return nil
	}

	values := make(map[string]string, len(s))
	for key, value := range s {
		values[key.String()] = value.String()
	}
	return values
}
