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

	"arcoris.dev/apimachinery/api/codec"
)

func TestErrorAtBuildsStructuredError(t *testing.T) {
	err := errorAt("codecs[0]", ErrInvalidCodec, ErrorReasonInvalidCodec, "codec must be non-nil")

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "codecs[0]", ErrorReasonInvalidCodec)
}

func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(
		"codecs[0].info.mediaTypes[0]",
		ErrDuplicateMediaType,
		ErrorReasonDuplicateMediaType,
		"codec media type %q duplicates codecs[%d]",
		codec.MediaTypeJSON,
		1,
	)

	if !strings.Contains(err.Error(), "application/json") || !strings.Contains(err.Error(), "1") {
		t.Fatalf("Error() = %q; want formatted detail", err.Error())
	}
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := codec.ErrorAt("codec.info", codec.ErrInvalidInfo, codec.ErrorReasonInvalidInfo, "bad info")

	err := wrapAt("codecs[0].info", ErrInvalidInfo, ErrorReasonInvalidInfo, "codec info is invalid", cause)

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(..., cause) = false")
	}
}

func TestInvalidInfoWrapsCodecInvalidInfo(t *testing.T) {
	_, err := New(fakeBaseCodec{info: codec.Info{}})

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
}

func TestNilCodecError(t *testing.T) {
	_, err := New(nil)

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "codecs[0]", ErrorReasonInvalidCodec)
}

func TestTypedNilCodecError(t *testing.T) {
	var c *fakeValueByteCodec

	_, err := New(c)

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "codecs[0]", ErrorReasonInvalidCodec)
}

func TestDuplicateMediaTypeErrorIs(t *testing.T) {
	_, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateMediaType)
}

func TestErrorDiagnosticPath(t *testing.T) {
	_, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireRegistryError(t, err, "codecs[1].info.mediaTypes[0]", ErrorReasonDuplicateMediaType)
}

func TestDuplicateFormatNoLongerErrors(t *testing.T) {
	_, err := New(
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeYAML),
	)

	requireNoError(t, err)
}
