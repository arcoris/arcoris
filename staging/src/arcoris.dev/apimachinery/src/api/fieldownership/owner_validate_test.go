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

func TestOwnerValidateAcceptsSimpleName(t *testing.T) {
	requireNoError(t, Owner("user-cli").Validate())
}

func TestOwnerValidateAcceptsSlashName(t *testing.T) {
	requireNoError(t, Owner("arcoris.dev/controller").Validate())
}

func TestOwnerValidateAcceptsColonName(t *testing.T) {
	requireNoError(t, Owner("user:anton").Validate())
}

func TestOwnerValidateRejectsEmpty(t *testing.T) {
	requireErrorIs(t, Owner("").Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsWhitespaceOnly(t *testing.T) {
	requireErrorIs(t, Owner(" \t").Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsLeadingWhitespace(t *testing.T) {
	requireErrorIs(t, Owner(" user-cli").Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsTrailingWhitespace(t *testing.T) {
	requireErrorIs(t, Owner("user-cli ").Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsControlCharacter(t *testing.T) {
	requireErrorIs(t, Owner("user\ncli").Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsInvalidUTF8(t *testing.T) {
	requireErrorIs(t, Owner(string([]byte{0xff})).Validate(), ErrInvalidOwner)
}

func TestOwnerValidateRejectsTooLong(t *testing.T) {
	requireErrorIs(t, Owner(strings.Repeat("a", MaxOwnerLength+1)).Validate(), ErrInvalidOwner)
}
