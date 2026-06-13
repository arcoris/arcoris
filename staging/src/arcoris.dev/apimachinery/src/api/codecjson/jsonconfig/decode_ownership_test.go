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

func TestDefaultDecodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := defaultDecodeOwnershipConfig()

	if config.UnknownFields != UnknownFieldReject {
		t.Fatalf("unknown fields = %d; want reject", config.UnknownFields)
	}
	if config.Validation != OwnershipValidationEnable {
		t.Fatalf("validation = %d; want enable", config.Validation)
	}
}

func TestResolveDecodeOwnershipConfig(t *testing.T) {
	t.Parallel()

	config := DecodeOwnershipConfig{}
	resolveDecodeOwnershipConfig(&config)

	if config.UnknownFields != UnknownFieldReject {
		t.Fatalf("unknown fields = %d; want reject", config.UnknownFields)
	}
	if config.Validation != OwnershipValidationEnable {
		t.Fatalf("validation = %d; want enable", config.Validation)
	}
}

func TestResolveDecodeOwnershipConfigKeepsDefaultValidationEnabledWithExplicitUnknownFields(t *testing.T) {
	t.Parallel()

	testCases := map[string]UnknownFieldMode{
		"reject": UnknownFieldReject,
		"ignore": UnknownFieldIgnore,
	}

	for name, unknownFields := range testCases {
		t.Run(name, func(t *testing.T) {
			config := DecodeOwnershipConfig{UnknownFields: unknownFields}
			resolveDecodeOwnershipConfig(&config)

			if config.UnknownFields != unknownFields {
				t.Fatalf("unknown fields = %d; want %d", config.UnknownFields, unknownFields)
			}
			if config.Validation != OwnershipValidationEnable {
				t.Fatalf("validation = %d; want enable", config.Validation)
			}
		})
	}
}

func TestResolveDecodeOwnershipConfigKeepsExplicitValidationDisable(t *testing.T) {
	t.Parallel()

	config := DecodeOwnershipConfig{Validation: OwnershipValidationDisable}
	resolveDecodeOwnershipConfig(&config)

	if config.Validation != OwnershipValidationDisable {
		t.Fatalf("validation = %d; want disable", config.Validation)
	}
}

func TestValidateDecodeOwnershipConfigRejectsUnknownModes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config DecodeOwnershipConfig
		path   string
	}{
		"unknown fields": {
			config: DecodeOwnershipConfig{UnknownFields: UnknownFieldMode(99), Validation: OwnershipValidationEnable},
			path:   "decode.ownership.unknown_fields",
		},
		"validation": {
			config: DecodeOwnershipConfig{UnknownFields: UnknownFieldReject, Validation: OwnershipValidationMode(99)},
			path:   "decode.ownership.validation",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateDecodeOwnershipConfig(testCase.config)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
