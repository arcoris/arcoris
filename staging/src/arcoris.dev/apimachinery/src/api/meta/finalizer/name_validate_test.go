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

package finalizer

import (
	"errors"
	"testing"
)

func TestNameValidateLexical(t *testing.T) {
	requireNoError(t, Name("cleanup").ValidateLexical())
	requireNoError(t, Name("control.arcoris.dev/cleanup").ValidateLexical())

	requireErrorIs(t, Name("").ValidateLexical(), ErrInvalidName)
	requireErrorIs(t, Name("cleanup_name").ValidateLexical(), ErrInvalidName)
}

func TestNameValidateLexicalStructuredError(t *testing.T) {
	err := Name("cleanup_name").ValidateLexical()
	requireErrorIs(t, err, ErrInvalidName)

	var finalizerErr *Error
	if !errors.As(err, &finalizerErr) {
		t.Fatalf("errors.As(%T) = false", finalizerErr)
	}
	if finalizerErr.Path != "finalizer.name" {
		t.Fatalf("Path = %q", finalizerErr.Path)
	}
	if finalizerErr.Reason != ErrorReasonInvalidCharacter {
		t.Fatalf("Reason = %q", finalizerErr.Reason)
	}
	if finalizerErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
