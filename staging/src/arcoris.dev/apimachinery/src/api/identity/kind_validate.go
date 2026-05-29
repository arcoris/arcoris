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

package identity

// Validate checks the strict kind grammar.
//
// The zero value is invalid as a complete kind identity.
func (k Kind) Validate() error {
	return validateKindValue(string(k))
}

// validateKindValue checks the API kind grammar.
//
// Kind validation is intentionally narrower than arbitrary Go type names:
// separators, Unicode, whitespace, and package-qualified names are not API
// identity tokens.
func validateKindValue(value string) error {
	if value == "" {
		return invalid(identityNameKind, value, ErrorReasonEmptyValue, detailKindNonEmpty)
	}

	if !isUpper(value[0]) {
		return invalid(identityNameKind, value, ErrorReasonInvalidEdge, detailKindUppercaseStart)
	}

	for i := 1; i < len(value); i++ {
		if !isUpper(value[i]) && !isLower(value[i]) && !isDigit(value[i]) {
			return invalidf(
				identityNameKind,
				value,
				ErrorReasonInvalidCharacter,
				"kind contains invalid byte %q at index %d",
				value[i],
				i,
			)
		}
	}

	return nil
}
