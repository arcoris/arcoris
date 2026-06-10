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
	"testing"
)

func TestOwnerValidateLexicalAcceptsSimpleName(t *testing.T) {
	requireNoError(t, owner("user-cli").ValidateLexical())
}

func TestOwnerValidateLexicalAcceptsSlashName(t *testing.T) {
	requireNoError(t, owner("arcoris.dev/controller").ValidateLexical())
}

func TestOwnerValidateLexicalAcceptsColonName(t *testing.T) {
	requireNoError(t, owner("user:anton").ValidateLexical())
}

func TestValidateNewOwnerRejectsInvalidOwner(t *testing.T) {
	err := validateNewOwner(Owner{})

	requireErrorIs(t, err, ErrInvalidOwner)
}

func TestOwnerValidateLexicalRejectsEmpty(t *testing.T) {
	requireOwnerReason(t, invalidOwner(""), ErrorReasonEmptyOwner)
}

func TestOwnerValidateLexicalRejectsWhitespaceOnly(t *testing.T) {
	requireOwnerReason(t, invalidOwner(" \t"), ErrorReasonWhitespaceOwner)
}

func TestOwnerValidateLexicalRejectsLeadingWhitespace(t *testing.T) {
	requireOwnerReason(t, invalidOwner(" user-cli"), ErrorReasonOwnerBoundaryWhitespace)
}

func TestOwnerValidateLexicalRejectsTrailingWhitespace(t *testing.T) {
	requireOwnerReason(t, invalidOwner("user-cli "), ErrorReasonOwnerBoundaryWhitespace)
}

func TestOwnerValidateLexicalRejectsControlCharacter(t *testing.T) {
	requireOwnerReason(t, invalidOwner("user\ncli"), ErrorReasonOwnerControlCharacter)
}

func TestOwnerValidateLexicalRejectsInvalidUTF8(t *testing.T) {
	requireOwnerReason(t, invalidOwner(string([]byte{0xff})), ErrorReasonInvalidOwnerUTF8)
}

func TestOwnerValidateLexicalRejectsTooLong(t *testing.T) {
	requireOwnerReason(t, invalidOwner(strings.Repeat("a", MaxOwnerLength+1)), ErrorReasonOwnerTooLong)
}

func requireOwnerReason(t *testing.T, owner Owner, reason ErrorReason) {
	t.Helper()

	err := owner.ValidateLexical()
	requireErrorIs(t, err, ErrInvalidOwner)
	requireFieldOwnershipError(t, err, "owner", reason)
}
