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

package annotations

import (
	"errors"
	"testing"
)

func TestValueValidateLexical(t *testing.T) {
	requireNoError(t, Value("").ValidateLexical())
	requireNoError(t, Value("human readable note").ValidateLexical())

	requireErrorIs(t, Value("bad\nnote").ValidateLexical(), ErrInvalidValue)
}

func TestValueValidateLexicalStructuredError(t *testing.T) {
	err := Value("bad\nnote").ValidateLexical()
	requireErrorIs(t, err, ErrInvalidValue)

	var annotationErr *Error
	if !errors.As(err, &annotationErr) {
		t.Fatalf("errors.As(%T) = false", annotationErr)
	}
	if annotationErr.Path != "annotation.value" {
		t.Fatalf("Path = %q", annotationErr.Path)
	}
	if annotationErr.Reason != ErrorReasonInvalidCharacter {
		t.Fatalf("Reason = %q", annotationErr.Reason)
	}
	if annotationErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
