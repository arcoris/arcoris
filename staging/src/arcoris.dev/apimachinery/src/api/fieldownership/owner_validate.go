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

package fieldownership

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// validateNewOwner checks the controlled owner value before construction returns it.
func validateNewOwner(owner Owner) error {
	return owner.ValidateLexical()
}

// ValidateLexical checks whether o can be stored as a field ownership identity.
//
// It does not check admission, authorization, request identity, RBAC, runtime
// component identity, storage presence, or ownership policy.
func (o Owner) ValidateLexical() error {
	text := o.text
	switch {
	case text == "":
		return errorAt(
			"owner",
			ErrInvalidOwner,
			ErrorReasonEmptyOwner,
			"owner is empty",
		)
	case !utf8.ValidString(text):
		return errorAt(
			"owner",
			ErrInvalidOwner,
			ErrorReasonInvalidOwnerUTF8,
			"owner is not valid UTF-8",
		)
	case strings.TrimSpace(text) == "":
		return errorAt(
			"owner",
			ErrInvalidOwner,
			ErrorReasonWhitespaceOwner,
			"owner is only whitespace",
		)
	case strings.TrimSpace(text) != text:
		return errorAt(
			"owner",
			ErrInvalidOwner,
			ErrorReasonOwnerBoundaryWhitespace,
			"owner has leading or trailing whitespace",
		)
	case len(text) > MaxOwnerLength:
		return errorfAt(
			"owner",
			ErrInvalidOwner,
			ErrorReasonOwnerTooLong,
			"owner exceeds %d bytes",
			MaxOwnerLength,
		)
	}

	for _, r := range text {
		if unicode.IsControl(r) {
			return errorAt(
				"owner",
				ErrInvalidOwner,
				ErrorReasonOwnerControlCharacter,
				"owner contains a control character",
			)
		}
	}

	return nil
}
