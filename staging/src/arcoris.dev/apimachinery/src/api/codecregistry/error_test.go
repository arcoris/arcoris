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

package codecregistry

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorStringNil(t *testing.T) {
	var registryError *Error

	if got := registryError.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestErrorStringIncludesPackageName(t *testing.T) {
	err := errorAt("codecs[0]", ErrInvalidCodec, ErrorReasonInvalidCodec, "codec must be non-nil")

	if got := err.Error(); got == "" || !strings.Contains(got, "codecregistry") {
		t.Fatalf("Error() = %q; want package name", got)
	}
}

func TestErrorUnwrapNil(t *testing.T) {
	var registryError *Error

	if got := registryError.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}

func TestErrorUnwrapExposesSentinel(t *testing.T) {
	err := errorAt("codecs[0]", ErrInvalidCodec, ErrorReasonInvalidCodec, "codec must be non-nil")

	requireErrorIs(t, err, ErrInvalidCodec)
}

func TestErrorAsRegistryError(t *testing.T) {
	err := errorAt("codecs[0]", ErrInvalidCodec, ErrorReasonInvalidCodec, "codec must be non-nil")

	var registryError *Error
	if !errors.As(err, &registryError) {
		t.Fatalf("errors.As(..., *Error) = false")
	}
}
