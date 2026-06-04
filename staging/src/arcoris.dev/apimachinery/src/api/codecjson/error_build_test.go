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

package codecjson

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

// TestErrorAtBuildsStructuredDiagnostic covers direct error construction.
func TestErrorAtBuildsStructuredDiagnostic(t *testing.T) {
	err := errorAt(rootPath().Member("x"), ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "bad JSON")

	requireErrorIs(t, err, ErrInvalidJSON)
	requireErrorIs(t, err, codec.ErrDecodeFailed)
	requireCodecJSONError(t, err, "$.x", ErrorReasonInvalidJSON)
}

// TestErrorfAtFormatsDetail covers formatted diagnostic detail.
func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(rootPath(), ErrTrailingData, codec.ErrDecodeFailed, ErrorReasonTrailingData, "token %q", "x")

	requireDetailContains(t, err, `token "x"`)
}

// TestWrapAtPreservesCause covers lower-level cause preservation.
func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("reader failed")
	err := wrapAt(rootPath(), ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "read failed", cause)

	requireErrorIs(t, err, cause)
	requireErrorIs(t, err, ErrInvalidJSON)
}
