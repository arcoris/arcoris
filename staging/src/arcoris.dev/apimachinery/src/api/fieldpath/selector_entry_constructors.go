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

// NewSelectorEntry constructs one selector field/value pair.
func NewSelectorEntry(field FieldName, value Literal) SelectorEntry {
	return SelectorEntry{
		field: field,
		value: value,
	}
}

// SelectorEntryFromString constructs one selector entry from a raw field name.
func SelectorEntryFromString(field string, value Literal) (SelectorEntry, error) {
	fieldName, err := NewFieldName(field)
	if err != nil {
		return SelectorEntry{}, err
	}

	entry := NewSelectorEntry(fieldName, value)
	if err := entry.ValidateStructure(); err != nil {
		return SelectorEntry{}, err
	}

	return entry, nil
}

// MustSelectorEntry constructs one selector entry or panics on invalid input.
func MustSelectorEntry(field string, value Literal) SelectorEntry {
	entry, err := SelectorEntryFromString(field, value)
	if err != nil {
		panic(err)
	}

	return entry
}
