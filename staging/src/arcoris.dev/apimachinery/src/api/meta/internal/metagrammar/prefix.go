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

// maxNamePrefixLength bounds server-side generated-name prefixes in bytes.
const maxNamePrefixLength = 63

// ValidateNamePrefix validates a server-side name generation prefix.
//
// Prefixes use the same safe byte set as object names, but they may end in a
// hyphen so callers can use values such as "worker-".
func ValidateNamePrefix(s string) *Violation {
	if s == "" {
		return violation(ReasonEmptyValue, "name prefix must be non-empty")
	}
	if len(s) > maxNamePrefixLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("name prefix length must be <= %d bytes", maxNamePrefixLength),
		)
	}
	if !lexical.IsDNS1123LabelEdge(s[0]) {
		return violation(ReasonInvalidEdge, "name prefix must start with a lowercase letter or digit")
	}
	for i := 0; i < len(s); i++ {
		if !lexical.IsDNS1123LabelChar(s[i]) {
			return violation(
				ReasonInvalidCharacter,
				fmt.Sprintf("name prefix contains invalid byte %q at index %d", s[i], i),
			)
		}
	}
	return nil
}
