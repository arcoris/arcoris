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

func TestDefaultDecodeLimitsConfig(t *testing.T) {
	t.Parallel()

	limits := defaultDecodeLimitsConfig()

	if limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", limits.MaxDepth, DefaultMaxDepth)
	}
	if limits.MaxNumberDigits != DefaultMaxNumberDigits {
		t.Fatalf("max number digits = %d; want %d", limits.MaxNumberDigits, DefaultMaxNumberDigits)
	}
	if limits.MaxDocumentBytes != 0 {
		t.Fatalf("max document bytes = %d; want unlimited", limits.MaxDocumentBytes)
	}
	if limits.MaxStringBytes != 0 {
		t.Fatalf("max string bytes = %d; want unlimited", limits.MaxStringBytes)
	}
}

func TestResolveDecodeLimitsConfig(t *testing.T) {
	t.Parallel()

	limits := DecodeLimitsConfig{}
	resolveDecodeLimitsConfig(&limits)

	if limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("max depth = %d; want %d", limits.MaxDepth, DefaultMaxDepth)
	}
	if limits.MaxNumberDigits != DefaultMaxNumberDigits {
		t.Fatalf("max number digits = %d; want %d", limits.MaxNumberDigits, DefaultMaxNumberDigits)
	}
}

func TestValidateDecodeLimitsConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		limits DecodeLimitsConfig
		path   string
	}{
		"max depth": {
			limits: DecodeLimitsConfig{MaxDepth: -1, MaxNumberDigits: DefaultMaxNumberDigits},
			path:   "decode.limits.max_depth",
		},
		"document bytes": {
			limits: DecodeLimitsConfig{MaxDepth: DefaultMaxDepth, MaxDocumentBytes: -1, MaxNumberDigits: DefaultMaxNumberDigits},
			path:   "decode.limits.max_document_bytes",
		},
		"string bytes": {
			limits: DecodeLimitsConfig{MaxDepth: DefaultMaxDepth, MaxStringBytes: -1, MaxNumberDigits: DefaultMaxNumberDigits},
			path:   "decode.limits.max_string_bytes",
		},
		"max number digits": {
			limits: DecodeLimitsConfig{MaxDepth: DefaultMaxDepth, MaxNumberDigits: -1},
			path:   "decode.limits.max_number_digits",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateDecodeLimitsConfig(testCase.limits)
			requireConfigErrorIs(t, err, ErrInvalidConfig)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
