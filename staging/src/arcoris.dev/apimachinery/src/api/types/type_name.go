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

import (
	"strings"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// TypeName is a dot-separated name for an owner-defined reusable structural descriptor.
//
// Descriptor names identify descriptor definitions, not Go packages or concrete Go
// types. Prefix segments are lower-case API namespaces. The final segment
// starts with an upper-case ASCII letter so named semantic types stand out from
// field names and namespace segments.
type TypeName string

// ParseTypeName validates s and returns it as a TypeName.
func ParseTypeName(s string) (TypeName, error) {
	name := TypeName(s)

	if !name.IsValid() {
		return "", descriptorErrorf(
			"type.name",
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonInvalidReferenceName,
			"type name %q does not match the dot-separated TypeName grammar",
			s,
		)
	}

	return name, nil
}

// IsValid reports whether n satisfies the dot-separated type-name grammar.
func (n TypeName) IsValid() bool {
	s := string(n)

	if s == "" || strings.Contains(s, "..") {
		return false
	}

	segments := strings.Split(s, ".")

	if len(segments) < 2 {
		return false
	}

	for i, segment := range segments {
		if segment == "" || !isASCII(segment) {
			return false
		}

		if i == len(segments)-1 {
			if !isUpper(segment[0]) {
				return false
			}

			for j := 1; j < len(segment); j++ {
				if !isLower(segment[j]) && !isUpper(segment[j]) && !isDigit(segment[j]) {
					return false
				}
			}

			continue
		}

		if !isLower(segment[0]) {
			return false
		}

		for j := 1; j < len(segment); j++ {
			if !isLower(segment[j]) && !isDigit(segment[j]) && segment[j] != '-' {
				return false
			}
		}
	}

	return true
}

// String returns the type name text.
func (n TypeName) String() string {
	return string(n)
}

// isASCII reports whether s contains only ASCII bytes.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}

	return true
}

// isLower reports whether b is an ASCII lower-case letter.
func isLower(b byte) bool {
	return lexical.IsASCIILower(b)
}

// isUpper reports whether b is an ASCII upper-case letter.
func isUpper(b byte) bool {
	return lexical.IsASCIIUpper(b)
}

// isDigit reports whether b is an ASCII digit.
func isDigit(b byte) bool {
	return lexical.IsASCIIDigit(b)
}
