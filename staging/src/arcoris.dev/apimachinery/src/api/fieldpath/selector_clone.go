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
