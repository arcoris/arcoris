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

package jsonconfig

import "testing"

func TestResolveDecodeStringConfig(t *testing.T) {
	t.Parallel()

	config := DecodeStringConfig{}
	resolveDecodeStringConfig(&config)

	if config.InvalidUTF8 != InvalidUTF8Reject {
		t.Fatalf("invalid UTF-8 mode = %d; want reject", config.InvalidUTF8)
	}
}

func TestValidateDecodeStringConfigRejectsUnsupportedReplace(t *testing.T) {
	t.Parallel()

	err := validateDecodeStringConfig(DecodeStringConfig{InvalidUTF8: InvalidUTF8Replace})
	requireConfigErrorIs(t, err, ErrUnsupportedConfig)
	requireErrorTextContains(t, err, "decode.strings.invalid_utf8")
}

func TestValidateDecodeStringConfigRejectsUnknownMode(t *testing.T) {
	t.Parallel()

	err := validateDecodeStringConfig(DecodeStringConfig{InvalidUTF8: InvalidUTF8Mode(99)})
	requireConfigErrorIs(t, err, ErrInvalidConfig)
	requireErrorTextContains(t, err, "decode.strings.invalid_utf8")
}
