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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/lexical"
)

// maxOpaqueSegmentLength bounds opaque metadata identifiers in bytes.
const maxOpaqueSegmentLength = 128

// ValidateSegment validates opaque metadata identity segments such as UID.
func ValidateSegment(s string) *Violation {
	if s == "" {
		return violation(ReasonEmptyValue, "segment must be non-empty")
	}
	if len(s) > maxOpaqueSegmentLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("segment length must be <= %d bytes", maxOpaqueSegmentLength),
		)
	}
	if HasWhitespace(s) {
		return violation(ReasonInvalidCharacter, "segment must not contain whitespace")
	}
	if HasUnsafeScalarChar(s) {
		return violation(ReasonInvalidCharacter, "segment must not contain control bytes or path separators")
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if lexical.IsASCIIAlnum(b) || b == '-' || b == '_' || b == '.' || b == ':' {
			continue
		}
		return violation(
			ReasonInvalidCharacter,
			fmt.Sprintf("segment contains invalid byte %q at index %d", b, i),
		)
	}
	return nil
}
