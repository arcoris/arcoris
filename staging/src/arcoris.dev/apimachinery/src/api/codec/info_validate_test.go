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

func TestInfoValidateAcceptsValidInfo(t *testing.T) {
	requireNoError(t, validInfo().Validate())
	requireNoError(t, ValidateInfo(validInfo()))
}

func TestInfoValidateRejectsZeroFormat(t *testing.T) {
	info := validInfo()
	info.Format = ""

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, ErrInvalidFormat)
}

func TestInfoValidateRejectsInvalidFormat(t *testing.T) {
	info := validInfo()
	info.Format = "JSON"

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, ErrInvalidFormat)
	requireCodecError(t, err, "codec.info.format", ErrorReasonInvalidFormat)
}

func TestInfoValidateRejectsNoMediaTypes(t *testing.T) {
	info := validInfo()
	info.MediaTypes = nil

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireCodecError(t, err, "codec.info.mediaTypes", ErrorReasonInvalidMediaType)
}

func TestInfoValidateRejectsInvalidMediaType(t *testing.T) {
	info := validInfo()
	info.MediaTypes = []MediaType{"Application/JSON"}

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, ErrInvalidMediaType)
	requireCodecError(t, err, "codec.info.mediaTypes[0]", ErrorReasonInvalidMediaType)
}

func TestInfoValidateRejectsDuplicateMediaTypes(t *testing.T) {
	info := validInfo()
	info.MediaTypes = []MediaType{MediaTypeJSON, MediaTypeJSON}

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireCodecError(t, err, "codec.info.mediaTypes[1]", ErrorReasonInvalidMediaType)
}

func TestInfoValidateRejectsNoTargets(t *testing.T) {
	info := validInfo()
	info.Targets = nil

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireCodecError(t, err, "codec.info.targets", ErrorReasonInvalidTarget)
}

func TestInfoValidateRejectsInvalidTarget(t *testing.T) {
	info := validInfo()
	info.Targets = []Target{"unknown"}

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireErrorIs(t, err, ErrInvalidTarget)
	requireCodecError(t, err, "codec.info.targets[0]", ErrorReasonInvalidTarget)
}

func TestInfoValidateRejectsDuplicateTargets(t *testing.T) {
	info := validInfo()
	info.Targets = []Target{TargetValue, TargetValue}

	err := info.Validate()

	requireErrorIs(t, err, ErrInvalidInfo)
	requireCodecError(t, err, "codec.info.targets[1]", ErrorReasonInvalidTarget)
}
