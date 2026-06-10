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

import "cmp"

// Equal reports whether e and other are the same selector entry.
func (e SelectorEntry) Equal(other SelectorEntry) bool {
	return e.Compare(other) == 0
}

// Compare imposes deterministic ordering on selector entries.
func (e SelectorEntry) Compare(other SelectorEntry) int {
	if order := cmp.Compare(e.field, other.field); order != 0 {
		return order
	}

	return e.value.Compare(other.value)
}

// Equal reports whether s and other contain the same canonical selector entries.
func (s Selector) Equal(other Selector) bool {
	return s.Compare(other) == 0
}

// Compare imposes deterministic total ordering on selectors.
func (s Selector) Compare(other Selector) int {
	for i := 0; i < len(s.entries) && i < len(other.entries); i++ {
		if cmp := s.entries[i].Compare(other.entries[i]); cmp != 0 {
			return cmp
		}
	}

	switch {
	case len(s.entries) < len(other.entries):
		return -1
	case len(s.entries) > len(other.entries):
		return 1
	default:
		return 0
	}
}
