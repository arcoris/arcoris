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

// FieldName is the canonical lowerCamelCase name of an object field.
//
// Field names are API contract names, not Go struct field names, JSON tags, or
// storage column names. The grammar is deliberately small so descriptors remain
// portable across codecs and languages.
type FieldName string

// ParseFieldName validates s and returns it as a FieldName.
func ParseFieldName(s string) (FieldName, error) {
	name := FieldName(s)

	if !name.IsValid() {
		return "", typeErrorf(
			"field.name",
			ErrInvalidField,
			TypeErrorReasonInvalidFieldName,
			"field name %q is not lowerCamelCase",
			s,
		)
	}

	return name, nil
}

// IsValid reports whether n is non-empty ASCII lowerCamelCase.
func (n FieldName) IsValid() bool {
	s := string(n)

	if s == "" || !isLower(s[0]) {
		return false
	}

	for i := 1; i < len(s); i++ {
		if !isLower(s[i]) && !isUpper(s[i]) && !isDigit(s[i]) {
			return false
		}
	}

	return true
}

// String returns the field name text.
func (n FieldName) String() string {
	return string(n)
}
