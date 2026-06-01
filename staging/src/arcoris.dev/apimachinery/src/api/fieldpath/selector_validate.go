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

// Validate checks whether s is a canonical selector with unique field names.
func (s Selector) Validate() error {
	if len(s.entries) == 0 {
		return nested(
			ErrInvalidSelector,
			ErrorReasonEmptySelector,
			"selector has no entries",
			ErrEmptySelector,
		)
	}

	for i, entry := range s.entries {
		if err := entry.Validate(); err != nil {
			return nested(
				ErrInvalidSelector,
				ErrorReasonInvalidEntry,
				"selector entry is invalid",
				err,
			)
		}

		if i == 0 {
			continue
		}

		if s.entries[i-1].field == entry.field {
			return nested(
				ErrInvalidSelector,
				ErrorReasonDuplicateSelectorField,
				"selector field name is duplicated",
				ErrDuplicateField,
			)
		}

		if s.entries[i-1].Compare(entry) > 0 {
			return newError(
				ErrInvalidSelector,
				ErrorReasonInvalidSelector,
				"selector entries are not in canonical order",
			)
		}
	}

	return nil
}
