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

import "testing"

func TestNormalizeCodecToken(t *testing.T) {
	if got := normalizeCodecToken(" JSON "); got != "json" {
		t.Fatalf("normalizeCodecToken() = %q", got)
	}
}

func TestValidateCodecTokenAcceptsFormatToken(t *testing.T) {
	err := validateCodecToken(
		"codec.test",
		"arcoris-binary_1",
		formatTokenOptions(),
		ErrInvalidFormat,
		ErrorReasonInvalidFormat,
		"test token",
	)

	requireNoError(t, err)
}

func TestValidateCodecTokenAcceptsMediaToken(t *testing.T) {
	err := validateCodecToken(
		"codec.test",
		"vnd.arcoris.object+json",
		mediaTypeTokenOptions(),
		ErrInvalidMediaType,
		ErrorReasonInvalidMediaType,
		"test token",
	)

	requireNoError(t, err)
}

func TestValidateCodecTokenWrapsLexicalViolation(t *testing.T) {
	err := validateCodecToken(
		"codec.test",
		"bad token",
		formatTokenOptions(),
		ErrInvalidFormat,
		ErrorReasonInvalidFormat,
		"test token",
	)

	requireErrorIs(t, err, ErrInvalidFormat)
	requireCodecError(t, err, "codec.test", ErrorReasonInvalidFormat)
}

func TestValidateCodecTokenRejectsInvalidUTF8(t *testing.T) {
	err := validateCodecToken(
		"codec.test",
		string([]byte{0xff}),
		formatTokenOptions(),
		ErrInvalidFormat,
		ErrorReasonInvalidFormat,
		"test token",
	)

	requireErrorIs(t, err, ErrInvalidFormat)
}
