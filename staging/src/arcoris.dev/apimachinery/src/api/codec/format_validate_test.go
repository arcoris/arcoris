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

func TestFormatValidateAcceptsKnownFormats(t *testing.T) {
	for _, format := range []Format{FormatJSON, FormatYAML, FormatCBOR} {
		t.Run(format.String(), func(t *testing.T) {
			requireNoError(t, format.Validate())
		})
	}
}

func TestFormatValidateAcceptsCustomFormat(t *testing.T) {
	requireNoError(t, Format("arcoris-binary_1").Validate())
}

func TestFormatValidateRejectsZero(t *testing.T) {
	err := Format("").Validate()

	requireErrorIs(t, err, ErrInvalidFormat)
	requireCodecError(t, err, pathCodecFormat, ErrorReasonInvalidFormat)
}

func TestFormatValidateRejectsWhitespace(t *testing.T) {
	err := Format(" json ").Validate()

	requireErrorIs(t, err, ErrInvalidFormat)
}

func TestFormatValidateRejectsUppercaseNonCanonical(t *testing.T) {
	err := Format("JSON").Validate()

	requireErrorIs(t, err, ErrInvalidFormat)
}

func TestFormatValidateRejectsSlash(t *testing.T) {
	err := Format("application/json").Validate()

	requireErrorIs(t, err, ErrInvalidFormat)
}
