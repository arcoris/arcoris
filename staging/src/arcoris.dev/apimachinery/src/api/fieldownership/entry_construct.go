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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// NewEntry constructs one owner/field-set record.
//
// Empty field sets are allowed at Entry level. State construction and
// transformation helpers prune empty entries when normalizing full ownership
// state.
func NewEntry(owner Owner, fields fieldpath.Set) (Entry, error) {
	entry := Entry{owner: owner, fields: fields}
	if err := entry.Validate(); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

// MustEntry constructs an entry or panics when owner or fields are invalid.
func MustEntry(owner Owner, fields fieldpath.Set) Entry {
	entry, err := NewEntry(owner, fields)
	if err != nil {
		panic(err)
	}

	return entry
}
