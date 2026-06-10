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

package labels

import (
	"errors"
	"testing"
)

func TestKeyValidateLexical(t *testing.T) {
	requireNoError(t, Key("role").ValidateLexical())
	requireNoError(t, Key("control.arcoris.dev/role").ValidateLexical())

	requireErrorIs(t, Key("").ValidateLexical(), ErrInvalidKey)
	requireErrorIs(t, Key("role_name").ValidateLexical(), ErrInvalidKey)
}

func TestKeyValidateLexicalStructuredError(t *testing.T) {
	err := Key("Role").ValidateLexical()
	requireErrorIs(t, err, ErrInvalidKey)

	var labelErr *Error
	if !errors.As(err, &labelErr) {
		t.Fatalf("errors.As(%T) = false", labelErr)
	}
	if labelErr.Path != "label.key" {
		t.Fatalf("Path = %q", labelErr.Path)
	}
	if labelErr.Reason != ErrorReasonInvalidCharacter {
		t.Fatalf("Reason = %q", labelErr.Reason)
	}
	if labelErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
