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

func TestErrorAsCodecError(t *testing.T) {
	err := ErrorAt("$", ErrEncodeFailed, ErrorReasonEncodeFailed, "encode failed")

	var codecErr *Error
	if !errors.As(err, &codecErr) {
		t.Fatalf("errors.As(%T) = false", codecErr)
	}
}

func TestInvalidInfoWrapsSpecificCause(t *testing.T) {
	info := validInfo()
	info.MediaTypes = []MediaType{"application/json; charset=utf-8"}

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestNilError(t *testing.T) {
	var err *Error
	if err.Error() != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", err.Error())
	}
	if err.Unwrap() != nil {
		t.Fatalf("Unwrap() = %v; want nil", err.Unwrap())
	}
}
