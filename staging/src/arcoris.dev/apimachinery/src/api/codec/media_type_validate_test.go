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

func TestMediaTypeValidateAcceptsKnownMediaTypes(t *testing.T) {
	for _, mediaType := range []MediaType{MediaTypeJSON, MediaTypeYAML, MediaTypeCBOR} {
		t.Run(mediaType.String(), func(t *testing.T) {
			requireNoError(t, mediaType.Validate())
		})
	}
}

func TestMediaTypeValidateAcceptsVendorMediaType(t *testing.T) {
	requireNoError(t, MediaType("application/vnd.arcoris.object+json").Validate())
}

func TestMediaTypeValidateRejectsZero(t *testing.T) {
	err := MediaType("").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
	requireCodecError(t, err, pathCodecMediaType, ErrorReasonInvalidMediaType)
}

func TestMediaTypeValidateRejectsNoSlash(t *testing.T) {
	err := MediaType("application").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestMediaTypeValidateRejectsEmptyType(t *testing.T) {
	err := MediaType("/json").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestMediaTypeValidateRejectsEmptySubtype(t *testing.T) {
	err := MediaType("application/").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestMediaTypeValidateRejectsParametersInV1(t *testing.T) {
	err := MediaType("application/json; charset=utf-8").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestMediaTypeValidateRejectsWhitespace(t *testing.T) {
	err := MediaType(" application/json ").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}

func TestMediaTypeValidateRejectsUppercaseNonCanonical(t *testing.T) {
	err := MediaType("Application/JSON").Validate()

	requireErrorIs(t, err, ErrInvalidMediaType)
}
