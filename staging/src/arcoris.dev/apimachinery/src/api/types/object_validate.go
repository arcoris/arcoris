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

package types

import "fmt"

// validateObject checks fields, duplicate names, and unknown-field policy.
func validateObject(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if !t.object.unknown.IsValid() {
		return typeErrorf(
			path+".unknown",
			ErrInvalidType,
			TypeErrorReasonInvalidUnknownPolicy,
			"unknown-field policy %d is not supported",
			t.object.unknown,
		)
	}
	seen := make(map[FieldName]struct{}, len(t.object.fields))
	for _, field := range t.object.fields {
		fieldPath := fmt.Sprintf("%s.fields[%s]", path, field.name)
		if !field.name.IsValid() {
			return typeErrorf(
				fieldPath+".name",
				ErrInvalidField,
				TypeErrorReasonInvalidFieldName,
				"field name %q is not lowerCamelCase",
				field.name,
			)
		}
		if _, ok := seen[field.name]; ok {
			return typeErrorf(
				fieldPath+".name",
				ErrDuplicateField,
				TypeErrorReasonDuplicateFieldName,
				"field name %q is declared more than once",
				field.name,
			)
		}
		seen[field.name] = struct{}{}
		if !field.presence.IsValid() {
			return typeErrorf(
				fieldPath+".presence",
				ErrInvalidField,
				TypeErrorReasonInvalidPresence,
				"field %q must be required or optional, got %s",
				field.name,
				field.presence,
			)
		}
		if err := validateType(field.typ, resolver, fieldPath+".type", resolving); err != nil {
			return err
		}
	}
	return nil
}
