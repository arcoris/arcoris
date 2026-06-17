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

package objectquery

// StaticFieldSet is an immutable in-memory SelectableFieldSet for tests and
// callers that already have resolved field definitions.
type StaticFieldSet struct {
	// fields is keyed by a private structural FieldRef key, not diagnostic text.
	fields map[fieldRefKey]SelectableField
}

// NewStaticFieldSet validates fields and returns an immutable resolver. Later
// duplicate refs replace earlier definitions, which keeps construction simple
// for generated field sets.
func NewStaticFieldSet(fields ...SelectableField) (StaticFieldSet, error) {
	out := StaticFieldSet{fields: map[fieldRefKey]SelectableField{}}
	for _, field := range fields {
		if err := field.Validate(); err != nil {
			return StaticFieldSet{}, invalidFieldError(err, "invalid selectable field")
		}
		out.fields[field.Ref.key()] = field
	}

	return out, nil
}

// ResolveSelectableField returns a registered field by exact FieldRef.
func (s StaticFieldSet) ResolveSelectableField(ref FieldRef) (SelectableField, bool) {
	field, ok := s.fields[ref.key()]
	return field, ok
}
