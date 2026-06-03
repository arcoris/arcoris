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

// Entries returns detached entries in deterministic owner order.
func (s State) Entries() []Entry {
	if len(s.entries) == 0 {
		return nil
	}

	entries := make([]Entry, len(s.entries))
	copy(entries, s.entries)

	return entries
}

// Owners returns sorted owners present in s.
func (s State) Owners() []Owner {
	if len(s.entries) == 0 {
		return nil
	}

	owners := make([]Owner, 0, len(s.entries))
	for _, entry := range s.entries {
		owners = append(owners, entry.owner)
	}

	return owners
}

// Fields returns the union of every owned field in s.
func (s State) Fields() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, entry := range s.entries {
		fields = unionTransform(fields, entry.fields)
	}

	return fields
}

// FieldsFor returns fields owned by owner, or an empty set when owner is absent.
//
// Invalid owner values behave as absent because this query API has no error return.
func (s State) FieldsFor(owner Owner) fieldpath.Set {
	for _, entry := range s.entries {
		if entry.owner == owner {
			return entry.fields
		}
	}

	return fieldpath.EmptySet()
}
