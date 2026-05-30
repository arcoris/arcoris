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

// Set stores label key/value metadata.
type Set map[string]string

// IsZero reports whether the set has no entries.
func (s Set) IsZero() bool { return len(s) == 0 }

// Len returns the number of label entries.
func (s Set) Len() int { return len(s) }

// Get returns the value for key.
func (s Set) Get(key Key) (Value, bool) {
	value, ok := s[key.String()]
	return Value(value), ok
}

// Has reports whether key is present.
func (s Set) Has(key Key) bool {
	_, ok := s[key.String()]
	return ok
}
