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

package metagrammar

import "arcoris.dev/apimachinery/api/internal/lexical"

// maxQualifiedNameSegmentLength bounds the unqualified name segment in bytes.
const maxQualifiedNameSegmentLength = 63

// ValidateQualifiedName validates metadata keys and finalizer names.
//
// The grammar is "[qualified-dns-prefix/]name". The name segment is lowercase
// ASCII, may contain digits, hyphen, and dot, and must start and end with an
// alphanumeric byte.
func ValidateQualifiedName(s string) *Violation {
	if _, _, ok := SplitQualifiedName(s); !ok {
		return violation(ReasonInvalidForm, "qualified name must be name or prefix/name")
	}

	return fromLexical(lexical.ValidateQualifiedName(s, lexical.QualifiedNameOptions{
		AllowPrefix:   true,
		MaxNameLength: maxQualifiedNameSegmentLength,
		AllowNameDot:  true,
		RequirePrefix: false,
	}))
}
