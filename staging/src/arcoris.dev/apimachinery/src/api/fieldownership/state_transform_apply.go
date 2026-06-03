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

// fieldTransform changes an owner's current fields using caller-provided fields.
type fieldTransform func(current fieldpath.Set, fields fieldpath.Set) fieldpath.Set

// transformOwnerFields applies transform to one owner's current field set.
func (s State) transformOwnerFields(
	owner Owner,
	fields fieldpath.Set,
	detail string,
	transform fieldTransform,
) (State, error) {
	if err := validateOwnerFields(owner, fields, detail); err != nil {
		return State{}, err
	}

	return s.replaceOwnerFields(owner, transform(s.FieldsFor(owner), fields))
}

// transformOtherOwnerFields applies transform to every owner except owner.
func (s State) transformOtherOwnerFields(
	owner Owner,
	fields fieldpath.Set,
	detail string,
	transform fieldTransform,
) (State, error) {
	if err := validateOwnerFields(owner, fields, detail); err != nil {
		return State{}, err
	}

	entries := make([]Entry, 0, len(s.entries))
	for _, entry := range s.entries {
		if entry.owner == owner {
			entries = append(entries, entry)
			continue
		}

		entries = append(entries, Entry{
			owner:  entry.owner,
			fields: transform(entry.fields, fields),
		})
	}

	return normalizeEntries(entries)
}

// replaceOwnerFields returns s with owner set to fields.
func (s State) replaceOwnerFields(owner Owner, fields fieldpath.Set) (State, error) {
	entries := make([]Entry, 0, len(s.entries)+1)
	for _, entry := range s.entries {
		if entry.owner == owner {
			continue
		}

		entries = append(entries, entry)
	}

	if !fields.IsEmpty() {
		entries = append(entries, Entry{owner: owner, fields: fields})
	}

	return normalizeEntries(entries)
}

// validateOwnerFields validates the common owner plus field-set input shape.
func validateOwnerFields(owner Owner, fields fieldpath.Set, detail string) error {
	if err := owner.Validate(); err != nil {
		return err
	}

	return validateFields(fields, detail)
}
