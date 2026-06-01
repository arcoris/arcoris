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

// sortSelectorEntries canonicalizes selector entries in-place by field name and
// then literal value.
//
// Selectors are intentionally tiny. A simple insertion sort keeps the code
// explicit, allocation-free, and free from reflection-heavy helpers.
func sortSelectorEntries(entries []SelectorEntry) {
	for i := 1; i < len(entries); i++ {
		current := entries[i]
		j := i - 1

		for ; j >= 0; j-- {
			if entries[j].Compare(current) <= 0 {
				break
			}

			entries[j+1] = entries[j]
		}

		entries[j+1] = current
	}
}
