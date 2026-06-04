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

package codec

import (
	"errors"
	"testing"
)

func TestErrorAtReturnsNilForNilError(t *testing.T) {
	if err := ErrorAt("$", nil, ErrorReasonDecodeFailed, "decode failed"); err != nil {
		t.Fatalf("ErrorAt nil sentinel = %v; want nil", err)
	}
}

func TestErrorAtWrapsSentinel(t *testing.T) {
	err := ErrorAt("", ErrDecodeFailed, ErrorReasonDecodeFailed, "bad input")

	requireErrorIs(t, err, ErrDecodeFailed)
	requireCodecError(t, err, "$", ErrorReasonDecodeFailed)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("reader failed")

	err := WrapAt("$.spec", ErrDecodeFailed, ErrorReasonDecodeFailed, "decode failed", cause)

	requireErrorIs(t, err, ErrDecodeFailed)
	requireErrorIs(t, err, cause)
}

func TestErrorDiagnosticPath(t *testing.T) {
	err := ErrorAt("$.desired", ErrInvalidDocument, ErrorReasonInvalidDocument, "invalid field")

	requireCodecError(t, err, "$.desired", ErrorReasonInvalidDocument)
}

func TestNormalizeDiagnosticPath(t *testing.T) {
	if got := normalizeDiagnosticPath(""); got != "$" {
		t.Fatalf("normalizeDiagnosticPath(empty) = %q", got)
	}
	if got := normalizeDiagnosticPath("$.desired"); got != "$.desired" {
		t.Fatalf("normalizeDiagnosticPath(path) = %q", got)
	}
}
