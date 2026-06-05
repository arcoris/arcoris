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
	err := errorAt("registrations[0].codec", ErrInvalidCodec, ErrorReasonInvalidCodec, "codec must be non-nil")

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(
		"registrations[1].id",
		ErrDuplicateEntryID,
		ErrorReasonDuplicateEntryID,
		"entry ID %q duplicates registrations[%d]",
		MustEntryID("json.public"),
		1,
	)

	if !strings.Contains(err.Error(), "json.public") || !strings.Contains(err.Error(), "1") {
		t.Fatalf("Error() = %q; want formatted detail", err.Error())
	}
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := codec.ErrorAt("codec.info", codec.ErrInvalidInfo, codec.ErrorReasonInvalidInfo, "bad info")

	err := wrapAt("registrations[0].info", ErrInvalidInfo, ErrorReasonInvalidInfo, "codec info is invalid", cause)

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(..., cause) = false")
	}
}

func TestInvalidInfoWrapsCodecInvalidInfo(t *testing.T) {
	_, err := New(testRegistration("json.public", fakeBaseCodec{info: codec.Info{}}))

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, codec.ErrInvalidInfo)
}

func TestNilCodecError(t *testing.T) {
	_, err := New(Register(MustEntryID("json.public"), nil))

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestTypedNilCodecError(t *testing.T) {
	var c *fakeValueByteCodec

	_, err := New(testRegistration("json.public", c))

	requireErrorIs(t, err, ErrInvalidCodec)
	requireRegistryError(t, err, "registrations[0].codec", ErrorReasonInvalidCodec)
}

func TestDuplicateEntryIDErrorIs(t *testing.T) {
	_, err := New(
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.public", codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireErrorIs(t, err, ErrDuplicateEntryID)
}

func TestErrorDiagnosticPath(t *testing.T) {
	_, err := New(
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.public", codec.FormatYAML, codec.MediaTypeJSON),
	)

	requireRegistryError(t, err, "registrations[1].id", ErrorReasonDuplicateEntryID)
}

func TestDuplicateFormatNoLongerErrors(t *testing.T) {
	_, err := New(
		testValueByteRegistration("json.public", codec.FormatJSON, codec.MediaTypeJSON),
		testValueByteRegistration("json.storage", codec.FormatJSON, codec.MediaTypeYAML),
	)

	requireNoError(t, err)
}
