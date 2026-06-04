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

func TestResolveDecodeObjectConfig(t *testing.T) {
	t.Parallel()

	config := DecodeObjectConfig{}
	resolveDecodeObjectConfig(&config)

	if config.DuplicateKeys != DuplicateKeyReject {
		t.Fatalf("duplicate keys = %d; want reject", config.DuplicateKeys)
	}
	if config.TrailingData != TrailingDataReject {
		t.Fatalf("trailing data = %d; want reject", config.TrailingData)
	}
	if config.UnknownEnvelopeFields != UnknownFieldReject {
		t.Fatalf("unknown envelope fields = %d; want reject", config.UnknownEnvelopeFields)
	}
}

func TestValidateDecodeObjectConfigRejectsUnknownModes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config DecodeObjectConfig
		path   string
	}{
		"duplicate keys": {
			config: DecodeObjectConfig{DuplicateKeys: DuplicateKeyMode(99), TrailingData: TrailingDataReject, UnknownEnvelopeFields: UnknownFieldReject},
			path:   "decode.objects.duplicate_keys",
		},
		"trailing data": {
			config: DecodeObjectConfig{DuplicateKeys: DuplicateKeyReject, TrailingData: TrailingDataMode(99), UnknownEnvelopeFields: UnknownFieldReject},
			path:   "decode.objects.trailing_data",
		},
		"unknown fields": {
			config: DecodeObjectConfig{DuplicateKeys: DuplicateKeyReject, TrailingData: TrailingDataReject, UnknownEnvelopeFields: UnknownFieldMode(99)},
			path:   "decode.objects.unknown_envelope_fields",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateDecodeObjectConfig(testCase.config)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
