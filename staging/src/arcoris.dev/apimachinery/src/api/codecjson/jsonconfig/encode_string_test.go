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

func TestDefaultEncodeStringConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeStringConfig()

	if config.EscapeHTML {
		t.Fatalf("escape HTML = true; want false")
	}
	if config.InvalidUTF8 != InvalidUTF8Reject {
		t.Fatalf("invalid UTF-8 = %d; want reject", config.InvalidUTF8)
	}
}

func TestResolveEncodeStringConfig(t *testing.T) {
	t.Parallel()

	config := EncodeStringConfig{}
	resolveEncodeStringConfig(&config)

	if config.InvalidUTF8 != InvalidUTF8Reject {
		t.Fatalf("invalid UTF-8 = %d; want reject", config.InvalidUTF8)
	}
}

func TestValidateEncodeStringConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeStringConfig
		target error
		path   string
	}{
		"unknown": {
			config: EncodeStringConfig{InvalidUTF8: InvalidUTF8Mode(99)},
			target: ErrInvalidConfig,
			path:   "encode.strings.invalid_utf8",
		},
		"replace unsupported": {
			config: EncodeStringConfig{InvalidUTF8: InvalidUTF8Replace},
			target: ErrUnsupportedConfig,
			path:   "encode.strings.invalid_utf8",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeStringConfig(testCase.config)
			requireConfigErrorIs(t, err, testCase.target)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
