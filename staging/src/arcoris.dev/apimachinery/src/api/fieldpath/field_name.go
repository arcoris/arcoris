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

// FieldName is one fixed field name in a semantic path.
//
// Field names are path tokens only. They must be non-empty, but descriptor
// field-name grammar and resource-specific rules belong to higher layers.
type FieldName string

// NewFieldName validates name as a fixed field-name token.
func NewFieldName(name string) (FieldName, error) {
	fieldName := FieldName(name)
	if err := fieldName.ValidateStructure(); err != nil {
		return "", err
	}

	return fieldName, nil
}

// MustFieldName validates name or panics.
//
// It is intended for tests and static semantic-path declarations. Runtime
// callers that receive untrusted text should use NewFieldName.
func MustFieldName(name string) FieldName {
	fieldName, err := NewFieldName(name)
	if err != nil {
		panic(err)
	}

	return fieldName
}

// String returns the field name text.
func (n FieldName) String() string {
	return string(n)
}

// IsZero reports whether n is the absent field-name value.
func (n FieldName) IsZero() bool {
	return n == ""
}

// ValidateStructure checks the base fieldpath field-name invariant.
//
// It does not validate descriptor field-name grammar or resource-specific
// semantics.
func (n FieldName) ValidateStructure() error {
	if !n.IsZero() {
		return nil
	}

	return nested(
		ErrInvalidElement,
		ErrorReasonEmptyFieldName,
		"field name is empty",
		ErrEmptyFieldName,
	)
}
