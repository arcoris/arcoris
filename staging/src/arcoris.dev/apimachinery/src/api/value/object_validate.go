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

package value

// validateObjectField checks one object field before payload insertion.
//
// The existing slice contains only already validated and cloned fields. Passing
// it in keeps duplicate-name detection local to object construction without
// storing an index in the final payload.
func validateObjectField(index int, field Field, existing []Field) error {
	if field.Name == "" {
		return errorf(
			objectFieldNamePath(index),
			ErrEmptyName,
			ErrorReasonEmptyName,
			"object field name is empty",
		)
	}

	if field.Value.IsZero() {
		return errorf(
			objectFieldValuePath(index),
			ErrInvalidField,
			ErrorReasonInvalidValue,
			"object field %q has an invalid zero value",
			field.Name,
		)
	}

	if hasObjectFieldName(existing, field.Name) {
		return errorf(
			objectFieldNamePath(index),
			ErrDuplicateName,
			ErrorReasonDuplicateName,
			"object field name %q is duplicated",
			field.Name,
		)
	}

	return nil
}

// hasObjectFieldName performs the intentionally small linear duplicate check.
//
// It trades O(n) lookup for lower allocation and simpler payload invariants,
// which is the better default for short API objects.
func hasObjectFieldName(fields []Field, name string) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}

	return false
}
